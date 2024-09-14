package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"sync/atomic"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/eddsa/signing"
	"github.com/bnb-chain/tss-lib/v2/test"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
)

const (
	participants = 5
	threshold    = 3
)

func main() {
	// Keygen: Load keygen fixtures
	keys, signPIDs, err := keygen.LoadKeygenTestFixturesRandomSet(threshold+1, participants)
	if err != nil {
		log.Fatalf("Failed to load keygen fixtures: %v", err)
	}

	// Convert public key to Stellar address
	stellarAddr, err := getStellarAddressFromKey(keys[0].EDDSAPub.X(), keys[0].EDDSAPub.Y())
	if err != nil {
		log.Fatalf("Failed to convert public key to Stellar address: %v", err)
	}
	fmt.Printf("Stellar Address: %s\n", stellarAddr)

	// Create Stellar transaction
	tx, err := createStellarTransaction(stellarAddr)
	if err != nil {
		log.Fatalf("Failed to create Stellar transaction: %v", err)
	}

	// Sign transaction with TSS
	sigBase64, err := signTransactionWithTSS(tx, keys, signPIDs)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Attach signature and broadcast
	tx, err = appendSignatureToTransaction(stellarAddr, tx, sigBase64)
	if err != nil {
		log.Fatalf("Failed to attach signature to transaction: %v", err)
	}

	client := horizonclient.DefaultTestNetClient
	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalf("Failed to broadcast transaction: %v", err)
	}

	fmt.Printf("Transaction successful! Hash: %s\n", resp.Hash)
}

// getStellarAddressFromKey converts X and Y coordinates from EDDSA public key to a Stellar address.
func getStellarAddressFromKey(pkX, pkY *big.Int) (string, error) {
	pubKey := edwards.PublicKey{
		Curve: edwards.Edwards(),
		X:     pkX,
		Y:     pkY,
	}

	pubBytes := pubKey.Serialize()
	stellarPubKey := ed25519.PublicKey(pubBytes)

	stellarAddress, err := strkey.Encode(strkey.VersionByteAccountID, stellarPubKey)
	if err != nil {
		return "", fmt.Errorf("failed to encode public key to Stellar address: %v", err)
	}
	return stellarAddress, nil
}

// createStellarTransaction builds a Stellar transaction with a payment operation.
func createStellarTransaction(sourceAddress string) (*txnbuild.Transaction, error) {
	client := horizonclient.DefaultTestNetClient
	accountRequest := horizonclient.AccountRequest{AccountID: sourceAddress}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to load source account: %v", err)
	}

	txParams := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		Operations: []txnbuild.Operation{
			&txnbuild.Payment{
				SourceAccount: sourceAccount.GetAccountID(),
				Destination:   "GBZFRQE42G2ULRFFITXP2UZAXRBYKQM7R7LZ3QS7YHDUUI5QQRHGBZCY",
				Amount:        "2",
				Asset:         txnbuild.NativeAsset{},
			},
		},
		BaseFee:       txnbuild.MinBaseFee,
		Preconditions: txnbuild.Preconditions{TimeBounds: txnbuild.NewTimeout(43200)},
	}

	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: %v", err)
	}

	return tx, nil
}

func signTransactionWithTSS(tx *txnbuild.Transaction, keys []keygen.LocalPartySaveData, signPIDs tss.SortedPartyIDs) (string, error) {
	p2pCtx := tss.NewPeerContext(signPIDs)
	parties := make([]*signing.LocalParty, len(signPIDs))

	errCh := make(chan *tss.Error, len(signPIDs))
	outCh := make(chan tss.Message, len(signPIDs))
	endCh := make(chan *common.SignatureData, len(signPIDs))

	// Use SharedPartyUpdater (from test) as the PartyUpdater
	updater := test.SharedPartyUpdater
	msgData, err := network.HashTransactionInEnvelope(tx.ToXDR(), network.TestNetworkPassphrase)
	if err != nil {
		return "", fmt.Errorf("failed to hash transaction: %v", err)
	}

	// Initialize TSS parties
	for i := range signPIDs {
		params := tss.NewParameters(edwards.Edwards(), p2pCtx, signPIDs[i], len(signPIDs), threshold)
		localParty := signing.NewLocalParty(new(big.Int).SetBytes(msgData[:]), params, keys[i], outCh, endCh, len(msgData)).(*signing.LocalParty)
		parties[i] = localParty

		go func(p *signing.LocalParty) {
			if err := p.Start(); err != nil {
				errCh <- err
			}
		}(localParty)
	}

	// Wait for signature data
	var completed int32
	for {
		select {
		case err := <-errCh:
			return "", fmt.Errorf("TSS error: %v", err)
		case msg := <-outCh:
			dest := msg.GetTo()
			if dest == nil { // Broadcast message
				for _, party := range parties {
					if party.PartyID().Index != msg.GetFrom().Index {
						go updater(party, msg, errCh)
					}
				}
			} else { // Direct message
				go updater(parties[dest[0].Index], msg, errCh)
			}
		case sigData := <-endCh:
			if atomic.AddInt32(&completed, 1) == int32(len(signPIDs)) {
				// Convert msgData from [32]byte to []byte before passing
				return processSignatureData(sigData, keys[0], msgData[:])
			}
		}
	}
}

// processSignatureData verifies and processes the signature data.
func processSignatureData(sigData *common.SignatureData, key keygen.LocalPartySaveData, msgData []byte) (string, error) {
	r := new(big.Int).SetBytes(sigData.R)
	s := new(big.Int).SetBytes(sigData.S)

	// Verify EDDSA signature
	pubKey := edwards.PublicKey{
		Curve: edwards.Edwards(),
		X:     key.EDDSAPub.X(),
		Y:     key.EDDSAPub.Y(),
	}

	if !edwards.Verify(&pubKey, msgData, r, s) {
		return "", fmt.Errorf("EDDSA signature verification failed")
	}

	sig := edwards.Signature{R: r, S: s}
	sigBytes := sig.Serialize()

	// Encode to base64 and return
	sigBase64 := base64.StdEncoding.EncodeToString(sigBytes)
	return sigBase64, nil
}

// appendSignatureToTransaction adds a signature to the Stellar transaction.
func appendSignatureToTransaction(stellarAddress string, tx *txnbuild.Transaction, sig string) (*txnbuild.Transaction, error) {
	cleanedTx, err := tx.ClearSignatures()
	if err != nil {
		return nil, fmt.Errorf("failed to clear signatures in tx: %v", err)
	}

	updatedTx, err := cleanedTx.AddSignatureBase64(network.TestNetworkPassphrase, stellarAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("failed to append signature to tx: %v", err)
	}

	return updatedTx, nil
}
