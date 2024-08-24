# Switchly

## Work Log

### Summary

```
24.08.2024 Saturday   2h 30m


Total                2h 30m
```


#### 24.08.2024 Saturday

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
   