package main

import (

	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"sync/atomic"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/eddsa/signing"
	"github.com/bnb-chain/tss-lib/v2/test"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

const (
	participants = 5
	threshold    = 3
)

func main() {
	// Step 1: Keygen

	// Load keygen fixtures
	keys, signPIDs, err := keygen.LoadKeygenTestFixturesRandomSet(threshold+1, participants)
	if err != nil {
		log.Fatalf("Failed to load keygen fixtures: %v", err)
	}

	// Convert the keys slice to a slice of pointers
	keyPointers := make([]*keygen.LocalPartySaveData, len(keys))
	for i := range keys {
		keyPointers[i] = &keys[i]
	}

	// Extract the X coordinate of the public key
	x := keys[0].EDDSAPub.X()

	// Convert X to a 32-byte array (Ed25519 public keys are 32 bytes)
	pubKeyBytes := x.Bytes()

	// Convert the 32-byte public key to a Stellar address
	addrHex, err := PublicKeyToStellarAddress(pubKeyBytes)
	if err != nil {
		fmt.Printf("Failed to convert public key to Stellar address: %v\n", err)
		return
	}

	fmt.Printf("Pub Key (X coordinate, hex): %s\n", hex.EncodeToString(pubKeyBytes))
	fmt.Printf("Stellar Address: %s\n", addrHex)

	// Step 2: Create a Stellar transaction
	tx, err := createStellarTransaction(addrHex)
	if err != nil {
		log.Fatalf("Failed to create Stellar transaction: %v", err)
	}

	// Step 3: TSS Signing
	decoratedSig, err := signTransactionWithTSS(tx, keyPointers, signPIDs, pubKeyBytes)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Step 4: Attach the signature and broadcast the transaction
	tx, err = appendSignatureToTransaction(tx, decoratedSig)
	if err != nil {
		log.Fatalf("Failed to attach signature to transaction: %v", err)
	}

	// client := horizonclient.DefaultTestNetClient
	// resp, err := client.SubmitTransaction(tx)
	// if err != nil {
	// 	log.Fatalf("Failed to broadcast transaction: %v", err)
	// }

	// fmt.Printf("Transaction successful! Hash: %s\n", resp.Hash)
}

func PublicKeyToStellarAddress(pubKeyBytes []byte) (string, error) {
	if len(pubKeyBytes) != 32 {
		return "", fmt.Errorf("invalid public key length: expected 32 bytes, got %d bytes")
	}
	stellarAddress, err := strkey.Encode(strkey.VersionByteAccountID, pubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to encode public key to Stellar address: %v", err)
	}
	return stellarAddress, nil
}

func createStellarTransaction(sourceAddress string) (*txnbuild.Transaction, error) {
	client := horizonclient.DefaultTestNetClient

	accountRequest := horizonclient.AccountRequest{AccountID: sourceAddress}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to load source account: %v", err)
	}

	// Define timebounds
	timebounds := txnbuild.NewTimeout(43200)

	// Define preconditions
	preconditions := txnbuild.Preconditions{
		TimeBounds: timebounds,
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
		Preconditions: preconditions,
	}

	tx, err := txnbuild.NewTransaction(txParams)
	fmt.Println(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: %v", err)
	}

	return tx, nil
}

func signTransactionWithTSS(tx *txnbuild.Transaction, keys []*keygen.LocalPartySaveData, signPIDs tss.SortedPartyIDs, pubKeyBytes []byte) (xdr.DecoratedSignature, error) {
	p2pCtx := tss.NewPeerContext(signPIDs)
	parties := make([]*signing.LocalParty, 0, len(signPIDs))

	errCh := make(chan *tss.Error, len(signPIDs))
	outCh := make(chan tss.Message, len(signPIDs))
	endCh := make(chan *common.SignatureData, len(signPIDs))

	updater := test.SharedPartyUpdater
	msgData, err := tx.Hash(network.TestNetworkPassphrase) // The hash of the transaction to be signed
	if err != nil {
		return xdr.DecoratedSignature{}, fmt.Errorf("failed to build transaction hash: %v", err)
	}

	// Initialize the parties
	for i := 0; i < len(signPIDs); i++ {
		params := tss.NewParameters(tss.Edwards(), p2pCtx, signPIDs[i], len(signPIDs), threshold)
		P := signing.NewLocalParty(new(big.Int).SetBytes(msgData[:]), params, *keys[i], outCh, endCh, len(msgData[:])).(*signing.LocalParty)
		parties = append(parties, P)
		go func(P *signing.LocalParty) {
			if err := P.Start(); err != nil {
				errCh <- err
			}
		}(P)
	}

	var ended int32
	for {
		select {
		case err := <-errCh:
			fmt.Printf("Error: %s\n", err)
			return xdr.DecoratedSignature{}, err

		case msg := <-outCh:
			dest := msg.GetTo()
			if dest == nil {
				for _, P := range parties {
					if P.PartyID().Index == msg.GetFrom().Index {
						continue
					}
					go updater(P, msg, errCh)
				}
			} else {
				if dest[0].Index == msg.GetFrom().Index {
					fmt.Printf("party %d tried to send a message to itself (%d)\n", dest[0].Index, msg.GetFrom().Index)
				}
				go updater(parties[dest[0].Index], msg, errCh)

			}

		case sigData := <-endCh:
			atomic.AddInt32(&ended, 1)
			if atomic.LoadInt32(&ended) == int32(len(signPIDs)) {
				fmt.Printf("Received signature data from %d participants\n", ended)

				// Extract the 32-byte public key
				pubKey := ed25519.PublicKey(pubKeyBytes)

				signature := append(sigData.R, sigData.S...)
				fmt.Println("signature is: ", signature)

				// Verify the signature
				ok := ed25519.Verify(pubKey, msgData[:], signature)
				fmt.Println("Signature verification is: ", ok)

				decoratedSig, err := CreateDecoratedSignature(pubKeyBytes, sigData.S)
				if err != nil {
					return xdr.DecoratedSignature{}, fmt.Errorf("failed to create decorated signature: %v", err)
				}

				return decoratedSig, nil
			}
		}
	}
}

// appendSignatureToTransaction adds a signature to the transaction.
func appendSignatureToTransaction(tx *txnbuild.Transaction, sig xdr.DecoratedSignature) (*txnbuild.Transaction, error) {
	// Marshal the transaction to XDR
	txXDR, err := tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction to XDR: %v", err)
	}

	// Unmarshal XDR to TransactionEnvelope
	var txEnvelope xdr.TransactionEnvelope
	err = xdr.SafeUnmarshal(txXDR, &txEnvelope)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction XDR: %v", err)
	}

	// Ensure that V1 is not nil and initialize the Signatures slice if necessary
	if txEnvelope.V1 == nil {
		return nil, fmt.Errorf("unexpected nil V0 in transaction envelope")
	}

	if txEnvelope.V1.Signatures == nil {
		txEnvelope.V1.Signatures = []xdr.DecoratedSignature{}
	}

	// Add the signature to the transaction envelope
	txEnvelope.V1.Signatures = append(txEnvelope.V1.Signatures, sig)

	// Marshal updated TransactionEnvelope to XDR
	updatedXDR, err := xdr.MarshalBase64(txEnvelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated transaction envelope: %v", err)
	}

	// Recreate the transaction from updated XDR
	updatedGenericTx, err := txnbuild.TransactionFromXDR(updatedXDR)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction from updated XDR: %v", err)
	}

	// Now we need to handle the specific transaction type
	updatedTx, ok := updatedGenericTx.Transaction()
	if !ok {
		return nil, fmt.Errorf("unexpected transaction type: %T", updatedGenericTx)
	}

	return updatedTx, nil
}

// CreateDecoratedSignature converts a byte slice signature to a DecoratedSignature.
func CreateDecoratedSignature(pubKeyBytes []byte, sig []byte) (xdr.DecoratedSignature, error) {
	// Convert the public key to XDR SignatureHint (last 4 bytes of the public key)
	hint := xdr.SignatureHint{}
	copy(hint[:], pubKeyBytes[len(pubKeyBytes)-4:])

	// Convert the signature bytes to XDR Signature
	xdrSig := xdr.Signature(sig)

	return xdr.DecoratedSignature{
		Hint:      hint,
		Signature: xdrSig,
	}, nil
}
