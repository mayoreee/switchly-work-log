# Switchly

## Work Log

### Summary

```
24.08.2024 Saturday   2h 30m
25.08.2024 Sunday   4h
30.08.2024 Friday  2h
07.09.2024 Saturday 4h 30m
08.09.2024 Sunday 6h
14.09.2024 Saturday 9h
17.09.2024 Tuesday 5h 30m

Total                33h 30m
```
### 17.09.2024 Tuesday

Today, we are learing about the Soroban Multisig Smart contract.

The multisig smart contract vault on Soroban [sample code here](https://github.com/mayoreee/switchly-work-log/blob/main/demo/soroban/src/lib.rs) would have custom logic that enforces signature thresholds, meaning a set number of validator approvals are required for any transaction to proceed. However, this requires the Switchly each Switchly validators maintaining their own private/public key pair on Stellar - does this complicate implementation on Stellar?

Practically, how would this work. Vadlidator A launches the chain and depoloys the smart contract. Additional validators choin the network and present their Stellar public key to Validator A, which will then deploy the smart contract on Soroban, initializing it with all Validator public keys as signers. 

To transfer funds, each validator would independently sign an approval transaction calling the smart contract and the funds will then be trasferred after the threshold approval is met.

On churn event, each validator would independently sign a transaction calling the smart contract and then folks will be churned in and out  after the threshold approval is met. Hmmm


I have also published final report on work done here: https://docs.google.com/document/d/1DWyCTxvt5cG1OPxHXjC8UFNyH8KqZdFjd-yZuJoxouk/edit#heading=h.c3uvvzqtoxph

### 14.09.2024 Saturday

Last time, I made my initial implementation of the TSS for Stellar, but couldn't get a valid signature. Today's focus was on understanding this problem better and fixing the code for a valid Stellar signature. 

Finally, after debugging and studying more on the edwards/ed25519 curve, I succeffully signed and broadcasted transactions on the Stellar testnet with a new approach. Haha I can't believe how I implemented this before!

Below is my revised implementation summary and full code [here](https://github.com/mayoreee/switchly-work-log/blob/main/demo/stellar/vault.go)

**Approach Summary:**

**1. Key Generation:** Generated TSS keys for multiple parties.

**2. Address Conversion:** Converted an EDDSA public key to a Stellar address.

**3. Transaction Creation:** Created a Stellar payment transaction using the Stellar Go SDK.

**4. TSS Signing:** Signed the transaction using TSS, involving multiple parties for enhanced security.

**5. Broadcasting:** Broadcasted the signed transaction to the Stellar network.


**Output:**
- Stellar Address: `GAQVFGRFNC6RNZ65JMPOMZKKZ66B3GCBO6RPUFFGSX3WCHCN5OXU74Z3`
- Transaction Hash: `6495b6142ccc57ee08dbdc9842c4a1ef686e34e33da4b8452de7b06bd074daf5` [View in explorer](https://testnet.stellarchain.io/transactions/6495b6142ccc57ee08dbdc9842c4a1ef686e34e33da4b8452de7b06bd074daf5)


***Note: This is a transaction that transfers 2 XLM from the TSS-controlled account `GAQVFGRFNC6RNZ65JMPOMZKKZ66B3GCBO6RPUFFGSX3WCHCN5OXU74Z3` to my personal account `GBZFRQE42G2ULRFFITXP2UZAXRBYKQM7R7LZ3QS7YHDUUI5QQRHGBZCY`***

**Steps to Reproduce:**
```cmd
git clone https://github.com/mayoreee/switchly-work-log.git
```
```cmd
cd demo/stellar && go mod tidy
```
```cmd
go run .
```


Sweet! I also heard of the Soroban Multisig. Will look into that next



### 07.09.2024 Saturday - 08.09.2024 Sunday

Bulding the threshold signature scheme (TSS) with Stellar integration.

Okay, back from hiatus. The objective this time is to implement TSS using the Binance TSS library [tss-lib](https://github.com/bnb-chain/tss-lib) and apply it to sign a Stellar transaction with a group of participants.


**Step 1: Key Generation**

We start by generating the cryptographic keys for the participants in the Threshold Signature Scheme (TSS). In TSS, the private key is split into multiple parts, with each participant holding a share. Here’s how we do that:

```go
keys, signPIDs, err := keygen.LoadKeygenTestFixturesRandomSet(threshold+1, participants)
```
- Threshold: The minimum number of participants required to sign (in this case, 3).
- Participants: The total number of participants (set to 5 here).

I then convert the generated keys into pointers to use them in other functions:

```go
keyPointers := make([]*keygen.LocalPartySaveData, len(keys))
for i := range keys {
    keyPointers[i] = &keys[i]
}
```

Next, we extract the X-coordinate from the first public key (TSS keys are on an elliptic curve, so they have X and Y coordinates). We need this X-coordinate to convert it into a Stellar address format:

```go
x := keys[0].EDDSAPub.X()
pubKeyBytes := x.Bytes()
```

**Step 2: Convert Public Key to Stellar Address**

Stellar uses a unique address format, so I had to convert the public key (32 bytes) to a Stellar-compatible address using this function:

```go
addrHex, err := PublicKeyToStellarAddress(pubKeyBytes)

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
```

This allows the multi-party generated key to be used as a valid Stellar address.

**Step 3: Create a Stellar Transaction**

Once we have the Stellar address, we can create a transaction. The goal here is to send 2 XLM to another address on the Stellar testnet:

```go
tx, err := createStellarTransaction(addrHex)

func createStellarTransaction(sourceAddress string) (*txnbuild.Transaction, error) {
	client := horizonclient.DefaultTestNetClient

	accountRequest := horizonclient.AccountRequest{AccountID: sourceAddress}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to load source account: %v", err)
	}

	// Define timebounds
	timebounds := txnbuild.NewTimeout(43200) // 12 hours before expiry

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
```
The function createStellarTransaction does the following:

- Loads the source account details from the Stellar testnet.
- Builds a payment operation to send XLM from the source address to the recipient.
- Prepares a transaction with some basic options like a sequence number, base fee, and time bounds.


Now, the most important part: signing the Stellar transaction using TSS. In TSS, no single party holds the full private key. Instead, each party holds a "share" of the key. To sign a transaction, we need a minimum number of parties (the threshold) to collaborate.

**Step 4: Prepare the signing parties**
Each party is initialized with the necessary key share and communication setup:

```go
p2pCtx := tss.NewPeerContext(signPIDs)
parties := make([]*signing.LocalParty, 0, len(signPIDs))
```
We also set up communication channels to handle messages and errors:

```go
errCh := make(chan *tss.Error, len(signPIDs))
outCh := make(chan tss.Message, len(signPIDs))
endCh := make(chan *common.SignatureData, len(signPIDs))
```
Next, we generate the hash of the Stellar transaction, which is what we will sign:

```go
msgData, err := tx.Hash(network.TestNetworkPassphrase)
```

**Step 5: Initialize the signing parties**
We loop through each party and initialize them with their key shares and parameters (like the number of parties, threshold, etc.):

```go
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
```
This starts the signing process for each party. Each party will generate a partial signature using its key share.

**Step 6: Message passing**
In TSS, parties need to communicate with each other to combine their partial signatures. The communication happens via message passing through the channels:

```go
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
            go updater(parties[dest[0].Index], msg, errCh)
        }
    }
}
```
- outCh: Sends messages between the parties to coordinate the signing process.
- errCh: Captures any errors.
- updater: A helper function that ensures messages get sent to the right parties.


**Step 7: Collect signatures**
When all participants have signed, the partial signatures are combined into a final signature:

```go
case sigData := <-endCh:
    atomic.AddInt32(&ended, 1)
    if atomic.LoadInt32(&ended) == int32(len(signPIDs)) {
        signature := append(sigData.R, sigData.S...)
        ok := ed25519.Verify(pubKey, msgData[:], signature)
        fmt.Println("Signature verification is: ", ok)
```
Here, the signature parts (R and S) are combined, and we verify the signature using the public key and transaction hash.


**Step 8: Attach Signature to Transaction**

Once the signature is ready, it gets attached to the transaction:

```go
tx, err := appendSignatureToTransaction(tx, decoratedSig)
```
This function appends the signature to the Stellar transaction’s XDR (Stellar’s format for encoding transactions).


Okay, it's been a long week, all seems to be working except the signature verification step whtih isn't passing as expected. Next time, I will inviestigned this issue. Here the full code for now (also in the `demo/stellar/vault.go`):

```go
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

```


### 30.08.2024 

Reading the revised version of the Gennaro-Goldfeder paper [here](https://eprint.iacr.org/2019/114.pdf)

Piecing together Thorchain's implementation and usage of the GG TSS algo


### 25.08.2024 Sunday

**Focusing on the understanding the practial implementation of the vault acount creation and churn event**

During a key-generation/churn event, do we create a new vault with the new validator set added as signers then transfer the assets in? Or do we keep the existing vault and simply remove the validators that are being churned out from the signers list and add the churned in validators to the signers list?

Let's see how Thorchain does this!

Okay, [thorchain](https://gitlab.com/thorchain/thornode#churn) uses the former approach - create a new vault and move assets in there! Also uses the Genarro-Goldfeder Threshold Signature Scheme, which allows signing with no trusted party. Let's [read this paper](https://eprint.iacr.org/2020/540.pdf), should be insightful!

The GG TSS used in thorchain key-generation ceremony allows the nominated committee to construct the parameters for a new vault and the output is a public key from which the vault addresses for each chain is derived. Both secp256k1 and ed25519 chains are supported (Perfect! [Stellar is a ed25519 chain!](https://developers.stellar.org/docs/learn/encyclopedia/security/signatures-multisig)), and the vault address derivation is dependent on the chain.

Alright! that was a good read. I now have clarity of the possible approaches:
```txt
Approach 1 - Gennaro-Goldfeder TSS

Use Gennaro-Goldfeder TSS key-generation process on the Switchly protocol level. During each churn event, the nominated validators form a committee to construct a public key using this process. The constructed public key can then be used to derive vault addresses for each chain (including Stellar). Assets are moved into this new vault address.

The vault account on Stellar will be controlled by only one signer (the commitee of validators in the set) and a single signature generated by the Switchly validators in the GG TSS process will be the only signature that can facilitate a transaction from the vault. 

```

```txt
Approach 2 - Stellar Multisig

Use the Stellar native multisig mechanism to create an account. Account will be created by a leader (centralization worries) who becomes the master key by default. The leader adds on other signers (each validator in the committee). The leader removes and adds signers on each churn event.

```

Let's compare both approaches side-by-side

| Criteria | Gennaro-Goldfeder TSS | Stellar Multisig |
| :-------------------------------------- | :-------------: | :---------: |
| Trusted leader required to control master key | No | Yes |
| Number of signers added to the Stellar account | 1 | N - number of validators in the nominated committee |
| Upper limit on  possible threshold | No. [Supports any number of signers and any threshold value](https://eprint.iacr.org/2021/060) | Yes. Only a maximun of 20 signers can sign a transaction, otherwise it fails |
| Complexity | Less complex. The validators generate a single public key as a committee, which is then used to generate vault addresses across all supported chains | More complex and peculiar to Stellar only. Determining leadership can be a hassle with centralization concerns; and every validator on Switchly must have their own public key as signer on Stellar
| Cost efficiency | Highly cost efficient. Only one signer controls the Stellar account, so just a minimum balance of 0.5XLM must be maintained. Only pay for moving funds to a new vault address on each churn event | Less cost efficient. N number of validators will cost (0.5XLM * N) balance requirement. Also, on each churn event, there is the cost of adding N validators to the list of signers on the account and removing some in a transaction with multiple `Set Options` operations - more expensive!
| Computational  efficiency | Higher, as the vault public key is generated only once and used across multiple chains | Lower, as there's the need for determining leadership, adding and removing signers only on Stellar |
| Privacy | Yes. Partipating signers in a signing event are not revealed | No. Signers in a signing event are known publicly


Okay, I think the Gennaro-Goldfeder TSS approach is the right path forward. I will look into testing the vault creation and churn event next time, using this approach. I found a good developer library to work with [here](https://github.com/bnb-chain/tss-lib).


### 24.08.2024 Saturday

**Diving into Stellar-native signature options & appliacability to switchly**

**Summary:** ```Multisig appears to be ideal for the switchly use case. While Pre-authorized transaction and Hash(x) are great signature types, they seem to be unsuitable/inefficient/expensive for this use case - I think they are best suited for use cases requiring smaller quoroms and fewer outbound txs per hour, such as joint accounts & escrows.```


Stellar accounts can set their own signature weight, threshold values, and additional signing keys with the `Set Options` operation. This would be ideal for the switchly validators to control the vault account. Each validator will hold a signing key to the vault account. 

The threshold  on the vault account determines the requiered signature weight to authorize a transaction/operation. (e.g. payment operation).

The account will be created with a master key. Who controls the master key? should the master key weight be set to 0 weight to avoid centralization? if yes, what should be the order of event from account creation?

Okay, for a transaction/operation from the account to be successful, the sum of the weights of all signatures in the transaction must be greater than the threshold for that operation. Sounds like m of n threshold signature check going on here.


- Multisig

Interestingly, only 20 signatures can be attached to one transactions. If you have more signatures, it will fail even if they are all valid signatures! So we have a hard upper limit of `20 of n` signers in the validator set.

Another constraint is for each additional signer (other than the master key) the required minimun balance of the account increases by 0.5 XLM currently. This implies that if we have 50 validators for example, the vault account must hold at least 0.5 XLM * 50 = 25 XLM. Hmmm who pays for this? should each validator entering the set pay their share of total cost i.e. 0.5 XLM?

- Pre-authorized Transaction

A transaction that has been pre-constructed and validly pre-signed by the account signer(s), allowing any holder of the transaction to publish it to the network whenever they want - the sequence number must still be valid at that time tho.

You add the hash of a future transaction as a signer on the account. Once a matching transaction is sent to the network in future, the hash signer is removed from the account (regardless of whether the transaction succeeds or fails)... 

How will this play our practically? 
 1. Validators witness inbound txs on blockchain network X
 2. Validators come to a consesus on the ordering of the txs and the structure of equivalent txs to be settled on Stellar
 3. Any of the validator adds the hashes of the agreed settlement txs as signers to the vault account using the `Set Options` (expensive operation in time and $ - cost grows with the number of txs, and that's just at a point in time)
 4. Any of the validator broadcasts the txs to the validator in future. 

 Hmmm why do we need to presign a tx (by a comittee of validator) and broadcast same tx in future when we can sign (by the committee) and broadcast now? Sounds like ours isn't the right use case for pre-authorized transactions, which is more suitable to escrow type use cases.


- Hash(x) signature

Create x randomly, and generate hash(x) using SHA256. Then, add the hash as a signer to the vault account.

In a separate transaction (e.g. payment) you will add x as one of the signatures and submit to the network. At this point, x is known to the world and anyone who sees x on the network can sign tx for this vault account (wow!). Solution is to have more signers to ensure anyone who tries to use x independently will not have sufficeint weight to make a successful tx - haha sounds like we are doing multsig but in a hard way!


Okay! tomorrow, I will dive into multisig code implementation properly.
   