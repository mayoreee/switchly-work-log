# Maya Cardano Integration

## Work Log

### Batch One

```
20.03.2024 Wednesday   2h
08.04.2024 Monday      5h 20m
09.04.2024 Tuesday     1h
10.04.2024 Wednesday   2h
11.04.2024 Thursday    4h 40m
16.04.2024 Tuesday     1h
17.04.2024 Wednesday   6h
15.05.2024 Wednesday   9h
22.05.2024 Wednesday   6h
Total                 37h
```

#### 20.03.2024 Wednesday 2h

Right. Here we go. New things new times new vibes. Fresh.  
Working on Cardano now with the Maya crew (Hi guys)  

Reviewing the roadmap...  

EdDSA  
Edwards-curve Digital Signature Algorithm  
https://pkg.go.dev/github.com/katzenpost/core/crypto/eddsa  

Where are the TSS libraries?  
gitlab.com/thorchain/tss/go-tss/tss  

Ash said someone else is on bifrost/tss?  

My main job for this is to implement the bifrost cardano chain client.  

Might send this in a bit, double checking the timeline first:  
Good call, I'm sure it will be. Need to ask about the timeline. I find just
blitzing it and tackling whatever seems most efficient for the entire project
to be my preferred approach. Sometimes that means diving straight into the most
technical parts to gain a better understanding, sometimes it means getting some
of the tooling done, low hanging tedious fruit that will speed everything else
up once done. So splitting it all out into smaller milestones and then doing
those in a pre-determined order isn't what I'd advice.

Ash's Timeline:
- edsa maya (doesn't really make sense?)
- bifrost observer
- bifrost signer
  - thorchain go-tss
- smoke tests (1 month?!)

There's nothing on the timeline for setting up a testing environment.  
That took me months for dash to work out the masternode quorum and package
it in a developer/project friendly way.

Is EsDSA relevant in maya? Isn't it just listening to generic bifrost messages?  
Yeah a search for `dash` in `gitlab.com/mayachain/mayanode/x/mayachain/`  
returned just common stuff.

Smoke tests are harder to get right than the bifrost chain clients.

--------------------------------------------------

- node-launcher
  - cardano docker wrapped binary
- edsa maya (doesn't really make sense? it's just adding the generic/common chain info stuff right?)
- bifrost (itzamna/bitol/maximon -> me?)
  - observer
  - signer
    - thorchain go-tss (itzamna?)
- smoke tests (1 month?! seems tight)
- dex integration

--------------------------------------------------

Itzamna: Co-Founder and Lead Development.   
Responsible for Cardano Bitfröst development & TSS modifications.  

Kukulkan: Dev & Node ops. Node Launcher and DevOps (make a start and keep him updated)  
Responsible for Cardano fullnode in Kubernetes.  

Bitol: Full stack blockchain ops.  
Responsible for Cardano Bitfröst development.  

Maximon: Full stack blockchain ops.  
Responsible for Cardano Bitfröst development.  

Hunahpu: Full stack blockchain ops.  
Responsible for Cardano Dex aggregation development  

--------------------------------------------------

heimdall edsa / jp from thor (heimdall's working on another wallet)  
if not him then heimdall  
cadana foundation on telegram  
or reach out directly  

#### 08.04.2024 Monday 5h 20m

Well, that was quite a long diversion. Sorry for the wait guys.

Time for a strong maya day.

Official docker image?

```
docker pull ghcr.io/intersectmbo/cardano-node:8.9.1

docker run --rm ghcr.io/intersectmbo/cardano-node:8.9.1 run --help
```

https://developers.cardano.org/docs/get-started/installing-cardano-node

Who are these dudes and what do they want?

- Bryon
- Shelly 
- Goguen 
- Basho 
- Voltaire 

https://roadmap.cardano.org/en/

Okay then who are these other guys who are mentioned in the code/scripts?

- Alonzo
- Conway
- Allegra
- Mary

```
CARDANO_CONFIG=/opt/cardano/config/mainnet-config.json
CARDANO_DATABASE_PATH=/opt/cardano/data
```

```
--config /opt/cardano/config/mainnet-config.json
--topology /opt/cardano/config/mainnet-topology.json
--database-path /opt/cardano/data
```

Firstly, how do we select testnet?  
An environment variable: `NETWORK=testnet`  

> [Error] Managed configuration for network testnet does not exist

How do we generate these network config/topology files?

https://github.com/IntersectMBO/cardano-node/tree/master/scripts

```
./scripts/byron-to-alonzo/mkfiles.sh alonzo
```

They're using a linux package manager called nix.

There's a `start-cluster` command after running `nix flake lock` I'd like to get
at but I don't want to actually use nix on my mac.

Digging further...

```
docker run \
  --rm \
  --name cardano \
  -it \
  -v $(pwd):/mnt/cardano \
  -w /mnt/cardano \
  ubuntu

apt update && apt upgrade -y
apt install -y curl xz-utils sudo git

sh <(curl -L https://nixos.org/nix/install) --daemon
```

> Nix won't work in active shell sessions until you restart them.

> I don't support your init system yet; you may want to add nix-daemon manually.

/nix/store/v5m3612zrcgi0xl2477kkg4i74vdg6av-nix-2.21.1/bin/nix-daemon

```
docker exec -it cardano bash
nix-daemon
```

```
docker exec -it cardano bash
nix flake lock
nix --extra-experimental-features nix-command flake lock
nix --extra-experimental-features nix-command --extra-experimental-features flakes flake lock
nix --extra-experimental-features nix-command --extra-experimental-features flakes flake lock --update-input iohkNix
nix --extra-experimental-features nix-command --extra-experimental-features flakes flake update
```

> Nix feature 'nix-command' is disabled; add '--extra-experimental-features nix-command' to enable it

> Nix feature 'flakes' is disabled; add '--extra-experimental-features flakes' to enable it

> error: executing 'git': No such file or directory

```
vi /etc/nix/nix.conf

extra-experimental-features = flakes nix-command
```

> error: flake 'git+file:///mnt/cardano' does not provide attribute 'apps.aarch64-linux.default', 'defaultApp.aarch64-linux', 'packages.aarch64-linux.default' or 'defaultPackage.aarch64-linux'

Damn, that's the final nail in that approach I think.

Sure linux is the best os blah blah, but what about the cross-platform one  
command and you're up and running approach to development?

Back to working out what `start-cluster` would actually do...

`backend_nomad start-nomad-job "${dir}"`

What's `nomad`?  
Hashicorp container manager. Lite version of kubernetes.  
Great.  

Even if I run an entire ubuntu vm it'll get caught up on nix/flake not having a  
configuration for `aarch64`. Do I try and make that?

I'm leaning towards just finding a way to create the simplest config/topology  
file possible by hand.

Okay what are the specs for those files? Maybe I should start there.

`cardano-node/nix/workbench/backend/nomad-job.nix`

https://cardano.stackexchange.com/questions/8112/how-to-create-a-cardano-private-network

Config files:

- `Node Config`
- `DB Sync Config`
- `Submit API Config`
- `Node Topology`
- `Byron Genesis`
- `Shelley Genesis`
- `Alonzo Genesis`
- `Conway Genesis`

Plutip

```
nix run github:mlabs-haskell/plutip#plutip-core:exe:local-cluster -- -n 2
```

> error: flake 'git+file:///mnt/cardano' does not provide attribute 'apps.aarch64-linux.default', 'defaultApp.aarch64-linux', 'packages.aarch64-linux.default' or 'defaultPackage.aarch64-linux'

`vbars` recommended using one of the docker-compose files created as part of the
hydra project:
https://github.com/input-output-hk/hydra/blob/cc8f40fcd21261329ca36a3f7c03981721977631/demo/docker-compose.yaml

https://book.world.dev.cardano.org/environments/preprod/topology.json

```
./prepare-devnet
docker-compose up -d
watch -n 1 'docker-compose exec cardano-node cardano-cli query utxo --testnet-magic 42 --whole-utxo'
./seed-devnet.sh
```

> cardano-node The requested image's platform (linux/amd64) does not match the
  detected host platform (linux/arm64/v8) and no specific platform was
  requested

`export DOCKER_DEFAULT_PLATFORM=linux/amd64`

> Conway related : There was an error parsing the genesis
  file: "/devnet/genesis-conway.json" Error: "Error in $.poolVotingThresholds:
  parsing Cardano.Ledger.Conway.Core.PoolVotingThresholds
  (PoolVotingThresholds) failed, key \"pvtMotionNoConfidence\" not found"

Same with:
- poolVotingThresholds
  - pvtCommitteeNormal
  - pvtCommitteeNoConfidence
  - pvtHardForkInitiation
  - dvtMotionNoConfidence
- dRepVotingThresholds
  - dvtMotionNoConfidence
  - dvtCommitteeNormal
  - dvtCommitteeNoConfidence
  - dvtUpdateToConstitution
  - dvtHardForkInitiation
  - dvtPPNetworkGroup
  - dvtPPEconomicGroup
  - dvtPPTechnicalGroup
  - dvtPPGovGroup
  - dvtTreasuryWithdrawal

Okay with all those added, the node seems to stay up.

There are 2 utxos on the chain visible with that `query utxo` command above.

`seed-devnet` is hanging on the first tx even though at this stage I don't think  
hydra is relevant.

```
docker-compose exec cardano-node cardano-cli transaction submit --tx-file /devnet/seed-alice.signed --testnet-magic 42
```

> Command failed: transaction submit  Error: Error while submitting tx:
  ShelleyTxValidationError ShelleyBasedEraBabbage (ApplyTxError [UtxowFailure
  (UtxoFailure (AlonzoInBabbageUtxoPredFailure (ValueNotConservedUTxO
  (MaryValue (Coin 0) (MultiAsset (fromList []))) (MaryValue
  (Coin 900000000000) (MultiAsset (fromList [])))))),UtxowFailure (UtxoFailure
  (AlonzoInBabbageUtxoPredFailure (BadInputsUTxO (fromList [TxIn (TxId{unTxId =
  SafeHash "8c78893911a35d7c52104c98e8497a14d7295b4d9bf7811fc1d4e9f449884284"})
  (TxIx 0)]))))])

Need to switch to setting up maya node 3 + 4.


#### 09.04.2024 Tuesday 1h

So. Cardano private net. Kinda running. Kinda failing to send transations.

What do these errors mean.

Let's break it down.

```
ShelleyTxValidationError
  ShelleyBasedEraBabbage
    (ApplyTxError
      [UtxowFailure
        (UtxoFailure
          (AlonzoInBabbageUtxoPredFailure (ValueNotConservedUTxO
  (MaryValue (Coin 0) (MultiAsset (fromList []))) (MaryValue
  (Coin 900000000000) (MultiAsset (fromList [])))))),UtxowFailure (UtxoFailure
  (AlonzoInBabbageUtxoPredFailure (BadInputsUTxO (fromList [TxIn (TxId{unTxId =
  SafeHash "8c78893911a35d7c52104c98e8497a14d7295b4d9bf7811fc1d4e9f449884284"})
  (TxIx 0)]))))])
```

hmm maybe that's nothing to worry about tbh, looks like sending from a pool that  
doesn't exist. i.e. hydra. I was using a script there.

Let's see what we can do with the node.

Get balance would be nice:

```
cardano-cli shelley query tx-mempool --testnet-magic 42

docker-compose exec cardano-node cardano-cli query utxo --testnet-magic 42 --whole-utxo
docker-compose exec cardano-node cardano-cli query utxo --testnet-magic 42
docker-compose exec cardano-node cardano-cli query protocol-state --testnet-magic 42
docker-compose exec cardano-node cardano-cli query stake-snapshot --testnet-magic 42
docker-compose exec cardano-node cardano-cli shelley address key-gen --normal-key --key-output-format text-envelope

docker-compose exec cardano-node cardano-cli query tip --testnet-magic 42
```

Haha 0.87 sync progress with my own private node?? Something feels off.

```
docker-compose exec cardano-node cardano-cli query ledger-state --testnet-magic 42

docker exec -it demo-cardano-node-1 bash
```

```
cardano-cli address key-gen \
  --verification-key-file /tmp/payment.vkey \
  --signing-key-file /tmp/payment.skey

cat payment.skey
{
    "type": "PaymentSigningKeyShelley_ed25519",
    "description": "Payment Signing Key",
    "cborHex": "5820b8d5d7ad28c8099e21cee310dc56b77fdc8226d97a8dfb3178849ff8a59a729d"
}
cat payment.vkey
{
    "type": "PaymentVerificationKeyShelley_ed25519",
    "description": "Payment Verification Key",
    "cborHex": "58207bdf3b8f18bf27bc635a22ac5b8af9d313b81af355a316714bde52de17d603db"
}
```

Well that's address generation sorted.

What about sending a payment...

Oh stumbled on some good shit:  
https://cardano-course.gitbook.io/cardano-course/handbook/setting-up-a-local-cluster/create-a-local-cluster

```
cardano-cli keygen --secret utxo-keys/payment.000.key
```

```
cardano-cli transaction build-raw \
--shelley-era \
--invalid-hereafter $(expr $(cardano-cli query tip --testnet-magic 42 | jq .slot) + 1000) \
--fee 1000000 \
--tx-in $(cardano-cli byron transaction txid --tx transactions/tx0.tx)#0 \
--tx-out $(cat utxo-keys/user1.payment.addr)+29999999998000000 \
--out-file transactions/tx1.raw
```

https://docs.cardano.org/learn/cardano-keys/

tbh I think I need to go back a step or two to really understand what's going on  
here.

so I need to migrate from byron to shelley.

I haven't. Let's do that.

Hydra's great and all that but I need to understand this, I'd like to go from the  
ground up.

```
docker run \
  --rm \
  --name cardano \
  -it \
  -p 3001:3001 \
  -p 12798:12798 \
  --entrypoint bash \
  ghcr.io/intersectmbo/cardano-node:8.9.1

cardano-cli genesis create-cardano
```

```
{
   "Producers": [
     {
       "addr": "127.0.0.1",
       "port": 3001,
       "valency": 1
     }
   ]
 }
```


```
wget -P template/ https://raw.githubusercontent.com/input-output-hk/iohk-nix/master/cardano-lib/testnet-template/alonzo.json
wget -P template/ https://raw.githubusercontent.com/input-output-hk/iohk-nix/master/cardano-lib/testnet-template/byron.json
wget -P template/ https://raw.githubusercontent.com/input-output-hk/iohk-nix/master/cardano-lib/testnet-template/config.json
wget -P template/ https://raw.githubusercontent.com/input-output-hk/iohk-nix/master/cardano-lib/testnet-template/shelley.json
wget -P template/ https://raw.githubusercontent.com/input-output-hk/iohk-nix/master/cardano-lib/testnet-template/conway.json
```

#### 10.04.2024 Wednesday 2h

Messing with the hydra network configurations...

```
docker network create cardano

docker run \
  --rm \
  --name c1 \
  -it \
  --network cardano \
  -v /Users/adc/Desktop/maya-cardano/network/template:/mnt/template:ro \
  -v /Users/adc/Desktop/maya-cardano/network/mount:/mnt/cardano \
  --entrypoint /mnt/template/c1.sh \
  ghcr.io/intersectmbo/cardano-node:8.9.1
```

```
docker run \
  --rm \
  --name c2 \
  -it \
  --network cardano \
  -v /Users/adc/Desktop/maya-cardano/network/template:/mnt/template:ro \
  -v /Users/adc/Desktop/maya-cardano/network/mount:/mnt/cardano \
  -w /mnt/cardano/c2 \
  --entrypoint /mnt/template/c2.sh \
  ghcr.io/intersectmbo/cardano-node:8.9.1
```

```
docker exec \
-it \
-w /mnt/cardano/c1 \
-e CARDANO_NODE_SOCKET_PATH=/mnt/cardano/c1/node.socket \
c1 bash

cardano-cli query tip --testnet-magic 42
```

```
docker exec \
-it \
-w /mnt/cardano/c2 \
-e CARDANO_NODE_SOCKET_PATH=/mnt/cardano/c2/node.socket \
c2 bash

cardano-cli query tip --testnet-magic 42
```

#### 11.04.2024 Thursday 4h 40m

Domain: "c2" Starting Subscription Worker, valency 1  
Domain: "c2" Failed to start all required subscriptions  

I can curl both nodes on port 3000.  

They're not http protocol so it throws a suspend peer warning on close, but  
that proves it's nothing to do with the network setup, they can both talk to  
each other.  

Failed to start all required subscriptions  

Migrating from hydra until I can get it working...  

- node-config
  - RequiresNetworkMagic: "RequiresNoMagic"
  - ApplicationVersion: 1
  - LastKnownBlockVersion-Major: 6
  - Enable* <delete>

- condrad
  - committee
    - quorum: 0
  - committeeMinSize: 0

- shelly
  - updateQuorum: 1

- topology
  - producers: []

> [eecfc63a:cardano.node.Forge:Error:37] [2024-04-12 00:58:37.00 UTC] fromList [
  ("credentials",String "Cardano"),("val",Object (fromList [
  ("kind",String "TraceNoLedgerView"),("slot",Number 2700.0)]))]
  
At block 2700 it craps out, probably because c1 is not the leader.

`CARDANO_BLOCK_PRODUCER=true`

Node config  
- `getLast = Just 0.0.0.0`

It looks like the shelly epoch length is being overridden after the create  
command?  

They've also wiped the hashes for the config files, might make them a bit easier  
to tweak.  

Well, the hydra node is running in the babbage era:  

```
{
    "block": 736,
    "epoch": 152,
    "era": "Babbage",
    "hash": "7af7de436d2d5aba77232fab86f07b1d947630d676c4e741ddce6311d0ef5643",
    "slot": 762,
    "slotInEpoch": 2,
    "slotsToEpochEnd": 3,
    "syncProgress": "100.00"
}
```

Is that good?

Shelley -> Alonzo -> Babbage

So yeah babbage is great.

The only issue is, they're using an older docker image for hydra and I'm using  
the latest:

older:   ghcr.io/input-output-hk/cardano-node:8.7.3  
latest:  ghcr.io/intersectmbo/cardano-node:8.9.1  

input-output don't have the 8.9.1 version.

Also different authors. Interesting.

They've changed the protocol:

old:   `Protocol: Cardano`  
new:   `Protocol: RealPBFT`  

When they referenced `hydra-cluster/config/protocol-parameters` they were  
talking about the `alonso.json`

> Error in $.Protocol: Parsing of Protocol failed. "RealPBFT" is not a valid protocol

Even though they say that on their readme. Cool.

> TraceNodeNotLeader

This is my main issue atm.

My node is spinning up fine, but the epoch is 0, the slots aren't progressing,  
no leader.

Interestingly, the slot counts are still continuing from where I left off even  
if I delete the entire node database. This must be time based?  

Had to rewrite the sed time replacements to work better on my mac. Also removed  
the genesis hashes to enable more flexibility.  

Now the slot numbers are restarting, so that confirms they are indeed time based.  

Right now I have a one-line uber fast restart script for the hydra node.  
Going to hack away config variables until I break it, then try introduce those  
variables in my other config.

The hydra node immediately pushed a block using the genesis hash as the previous  
hash. My node doesn't do that.

Right, brilliant. Now I have a node up and running that I can restart with a   
single command back to block 0 running on the latest version.  

Would be nice to know how the shelly `staking` and `initialFunds` fields were  
generated. They seem to be critical.

For now at least I have enough to move on to generating new accounts, sending  
some of those initial funds, inspecting the transaction format etc.


#### 16.04.2024 Tuesday 1h

Been tough to get to this lately but gotta start somewhere.

Right, network up? Yes. Time to learn cardano properly.

```
cardano-cli ping --host localhost
```

Nope.

```
cardano-cli ping --unixsock node.socket
```

> node.socket Protocol error: Refused 32784 "version data mismatch


Whoa really?

Oh now it's `--magic` not `--testnet-magic`.

```
cardano-cli ping --unixsock node.socket --magic 42
```

#### 17.04.2024 Wednesday 6h

The plan: create a new address and send funds from faucet / genesis account.

> Genesis UTxO has no address

My restart script causes all the genesis addresses to mismatch because the data  
directory is deleted? Think restart is a no-go.  

> The era of the node and the tx do not match. The node is running in the Babbage era, but the transaction is for the Byron era.

INTERESTING.

What's the progression again here? I need a visual map.

Byron ->  
Shelley -> Alonzo -> Babbage  

Oh, after a few restarts I'm now getting this again:  
> TraceNodeNotLeader

What did I do?? How did I break it?  

Don't forget - protocol parameters are codenamed `alonzo`.

These are the differences between hydra and my genesis generated files:

- alonzo IDENTICAL (nice!!!)  
- byron  
  - blockVersionData.softforkRule.minThd  
    mine:    1000000000000000  
    theirs:   600000000000000  
  - blockVersionData.txFeePolicy.multiplier  
    mine:    43946000000  
    theirs:  43000000000  
  - blockVersionData.updateImplicit  
    mine:    21600  
    theirs:  10000  
  - bootStakeholders
    mine:    "931b71ae59f0f4c36674ab14f660f62e8266e81fb52ff86947d94efc": 1  
    theirs:  "7a4519c93d7be4577dd85bd524c644e6b809e44eae0457b43128c1c7": 1  
  - heavyDelegation  
    mine: loads for 931b...  
    theirs: empty object  
  - nonAvvmBalances:  
    mine: lots for 2657...  
    theirs: empty object  
  - startTime  
    set via script on both  
- conway  
  - dRepDeposit  
    mine:    500000000  
    theirs:    2000000  
  - govActionDeposit:  
    mine:   50000000000  
    theirs:  1000000000  
  - dRepVotingThresholds  
    everything has dvt prefix on older ver  
  - poolVotingThresholds  
    everything has pvt prefix on older ver  
- shelly  
  - activeSlotsCoeff  
    mine    0.05  
    theirs  1  
  - epochLength  
    mine    432000  
    theirs  5  
  - genDelegs  
    mine has data theirs is empty  
  - initialFunds  
    mine is empty thiers has data  
  - maxLovelaceSuypply  
    mine    45000000000000000   
    theirs      2000000000000  
  - protocolPrams.protocolVersion.major  
    mine    2  
    theirs  7  
  - staking.pools  
    they have an entry, I don't  
  - staking.stake  
    they have an entry, I don't  
  - systemStart  
    both set via script  
  - updateQuorum  
    mine    1  
    theirs  2  
- node-config  
  - \*GenesisHash  
    they've removed all these keys  
  - Protocol  
    mine    not set  
    theirs  `Cardano`
  
Right, now I remember what I did!  

I copied the hydra keys and stake configuration before to get it working, and  
my `start.sh` script re-generated the keys so everything broke.

```
cp ../network-hydra/config/byron-delegate.key ./data/genesis/delegate-keys/byron.000.key
cp ../network-hydra/config/byron-delegation.cert ./data/genesis/delegate-keys/byron.000.cert.json
cp ../network-hydra/config/kes.skey ./data/genesis/delegate-keys/shelley.000.kes.skey
cp ../network-hydra/config/opcert.cert ./data/genesis/delegate-keys/shelley.000.opcert.json
cp ../network-hydra/config/vrf.skey ./data/genesis/delegate-keys/shelley.000.vrf.skey
rm ./data/genesis/delegate-keys/shelley.000.counter.json
rm ./data/genesis/delegate-keys/*.vkey
```

Should leave just the 5 keys we need.  
I assume the vkey can just be generated from the skey.  

```
jq -n 'input | . + (input | pick(.activeSlotsCoeff, .genDelegs, .initialFunds, .maxLovelaceSupply, .updateQuorum, .epochLength, .protocolParams, .staking))' \
  ./data/genesis/shelley-genesis.json \
  ../network-hydra/config/genesis-shelley.json > \
  ./data/genesis/shelley-genesis.json.bak

mv ./data/genesis/shelley-genesis.json.bak ./data/genesis/shelley-genesis.json

jq -n 'input | . + (input | pick(.bootStakeholders, .heavyDelegation, .nonAvvmBalances, .blockVersionData))' \
  ./data/genesis/byron-genesis.json \
  ../network-hydra/config/genesis-byron.json > \
  ./data/genesis/byron-genesis.json.bak

mv ./data/genesis/byron-genesis.json.bak ./data/genesis/byron-genesis.json

jq '. |= with_entries(select(.key|test("Hash")|not)) | .Protocol="Cardano"' \
  ./data/genesis/node-config.json > \
  ./data/genesis/node-config.json.bak

mv ./data/genesis/node-config.json.bak ./data/genesis/node-config.json
```

There we go, a better reproducable quick and easy restart.

```
./genesis.sh && sleep 1 && ./start.sh
```

Okay back to spending the genesis tx.

Getting:

> Invalid argument `("2657WMsDfac6ij1TdtVjh68jLPCeuVTMnTgX3kciEYni3iTedEVbCCnTVAs5HyjKJ", 29999999999000000)'

when trying to spend. That's not very detailed.

maxLovelaceSupply is 2000000000000  
we're spending:  29999999999000000  
that could be the issue.  

> Invalid argument `("2657WMsDfac6ij1TdtVjh68jLPCeuVTMnTgX3kciEYni3iTedEVbCCnTVAs5HyjKJ", 1000000000000)'

> Genesis UTxO has no address

Ok...

What am I missing here?  
I was using the delegate byron key instead of the utxo byron key.  
Same error though.  

```
# cardano-cli address key-gen \
  # --verification-key-file alex.vkey \
  # --signing-key-file alex.skey

# cardano-cli address build \
  # --payment-verification-key-file alex.vkey \
  # --out-file alex.addr \
  # --testnet-magic 42

# alexAddr=$(cardano-cli address build --payment-verification-key-file alex.vkey --testnet-magic 42)

# cardano-cli query utxo --address $(cat alex.addr) --testnet-magic 42
# cardano-cli query protocol-parameters --testnet-magic 42 --out-file pparams.json


# cardano-cli query utxo --address $alexAddr --testnet-magic 42

genesisDir=/mount/data/genesis

cardano-cli keygen --secret $genesisDir/utxo-keys/payment.000.key

cardano-cli signing-key-address \
  --testnet-magic 42 \
  --secret payment.000.key > $genesisDir/utxo-keys/payment.000.addr

cardano-cli signing-key-address \
  --testnet-magic 42 \
  --secret $genesisDir/utxo-keys/byron.000.key > $genesisDir/utxo-keys/byron.000.addr

mkdir txs

cardano-cli issue-genesis-utxo-expenditure \
  --genesis-json $genesisDir/byron-genesis.json \
  --testnet-magic 42 \
  --tx txs/tx0.tx \
  --wallet-key $genesisDir/utxo-keys/byron.000.key \
  --rich-addr-from $(head -n 1 $genesisDir/utxo-keys/byron.000.addr) \
  --txout "(\"$(head -n 1 $genesisDir/utxo-keys/payment.000.addr)\", 1000000000000)"

cardano-cli submit-tx \
  --testnet-magic 42 \
  --tx txs/tx0.tx
```

Trying the same stuff with the hydra node.

```
cd /devnet

cardano-cli keygen --secret payment.000.key

cardano-cli signing-key-address \
  --testnet-magic 42 \
  --secret payment.000.key > payment.000.addr

cardano-cli signing-key-address \
  --testnet-magic 42 \
  --secret byron-delegate.key > byron-delegate.addr

mkdir txs

cardano-cli issue-genesis-utxo-expenditure \
  --genesis-json genesis-byron.json \
  --testnet-magic 42 \
  --tx txs/tx0.tx \
  --wallet-key byron-delegate.key \
  --rich-addr-from $(head -n 1 byron-delegate.addr) \
  --txout "(\"$(head -n 1 payment.000.addr)\", 1000000000000)"

cardano-cli submit-tx \
  --testnet-magic 42 \
  --tx txs/tx0.tx
```

Same error.

Oh the `cardano-cli` is a completely separate repo.  
That error is coming from `cardano-cli/src/Cardano/CLI/Byron/Tx.hs:133`.  
No

There's a `SpendGenesisUTxO` command.

Let's try checking the utxos for every address we can generate.

```
addr() {
  cardano-cli signing-key-address \
    --testnet-magic 42 \
    --secret $1 | head -n 1
}

cardano-cli query utxo --address $(addr /mount/data/genesis/genesis-keys/byron.000.key) --testnet-magic 42
cardano-cli query utxo --address $(addr /mount/data/genesis/utxo-keys/byron.000.key) --testnet-magic 42
cardano-cli query utxo --address $(addr /mount/data/node1/byron.000.key) --testnet-magic 42
```

Nothing. Okay well maybe the genesis utxo is a bit different.

How does this command work? `cardano-cli genesis initial-addr`

Watching the hydra demo again.

```
d=/mount/data
t=/mount/template

ccli() {
  cardano-cli ${@} --testnet-magic 42
}

ccli query utxo --whole-utxo

faucetAddr=$(ccli address build --payment-verification-key-file $t/creds/faucet.vk)
faucetTxin=$(ccli query utxo --address $faucetAddr --out-file /dev/stdout | jq -r 'keys[0]')
```

Argh no jq on box and no package manager to install.  
Can I find a prepackaged binary?

```
wget https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-arm64 -P ~/Downloads
chmod +x ~/Downloads/jq-linux-arm64
docker cp ~/Downloads/jq-linux-arm64 node1:/usr/local/bin/jq
```

```
paymentAddr=$(head -n 1 payment.000.addr)

ccli transaction build \
  --babbage-era \
  --cardano-mode \
  --change-address $faucetAddr \
  --tx-in $faucetTxin \
  --tx-out $paymentAddr+30000000 \
  --out-file seed.draft

ccli transaction sign \
    --tx-body-file seed.draft \
    --signing-key-file $t/creds/faucet.sk \
    --out-file seed.signed

ccli transaction submit --tx-file seed.signed

ccli query utxo --address $paymentAddr --out-file /dev/stdout
```

So you don't use the genesis utxo, but the faucet instead. Makes sense.  
2mins left on the clock, what a great place to pause for today.  


#### 15.05.2024 Wednesday 9h

Cardano bifrost

https://github.com/echovl/cardano-go  
https://developers.cardano.org/docs/get-started/dandelion-apis/  

cardano-cli query

Installed ghcup and haskell support for visual studio code.

Surely we can connect directly to the node socket.

Right there was an update last week to `gitlab.com/mayachain/devops` on the  
`cardano` branch. Will try and get mainnet and testnet running on my local  
machine as well as the regtest.  

Installing the mainnet config files like with `init.sh` into a `cardano-mainnet`  
docker volume:

```
docker run --rm -it -v ~/Desktop/maya-cardano-network/mainnet:/root/cardano alpine sh
```

```
mkdir -p /root/cardano/config
cd $_
apk add curl

curl -O -J https://book.world.dev.cardano.org/environments/mainnet/config.json
curl -O -J https://book.world.dev.cardano.org/environments/mainnet/db-sync-config.json
curl -O -J https://book.world.dev.cardano.org/environments/mainnet/submit-api-config.json
curl -O -J https://book.world.dev.cardano.org/environments/mainnet/topology.json
curl -O -J https://book.world.dev.cardano.org/environments/mainnet/byron-genesis.json
curl -O -J https://book.world.dev.cardano.org/environments/mainnet/shelley-genesis.json
curl -O -J https://book.world.dev.cardano.org/environments/mainnet/alonzo-genesis.json
curl -O -J https://book.world.dev.cardano.org/environments/mainnet/conway-genesis.json
```

```
docker run \
  --name cardano-mainnet \
  -d \
  -p 7731:3001 \
  -v ~/Desktop/maya-cardano-network/mainnet:/root/cardano \
  ghcr.io/intersectmbo/cardano-node:8.9.2@sha256:afa42f9ecdfe254b6ae8a999ce83912395caf83936adfa5e5589495a83e73022 \
    run \
    --topology /root/cardano/config/topology.json \
    --config /root/cardano/config/config.json \
    --database-path /root/cardano/db/ \
    --socket-path /root/cardano/node.socket \
    --host-addr "0.0.0.0" \
    --port "3001"
```

I want to expose the `node.socket` and see if I can query that directly...

```
~/Desktop/maya-cardano-network/mainnet/node.socket
```

Okay now let's do a little mitm logging with golang...

```
nc -vlU ~/Desktop/maya-cardano-network/mainnet/node.socket
nc -vlU /tmp/sock0120
```

```
docker exec -it -w /root/cardano cardano-mainnet bash

export CARDANO_NODE_NETWORK_ID=764824073
export CARDANO_NODE_SOCKET_PATH=node.socket
```

Googling this:

> goland dial unix domain socket connection refused

```
{
  "block": 2258175,
  "epoch": 104,
  "era": "Byron",
  "hash": "52d74cf2b8757dd67088aeb34bb4936d0fec66f3cd45ff0f8bf6704cc2542bce",
  "slot": 2259664,
  "slotInEpoch": 13264,
  "slotsToEpochEnd": 8336,
  "syncProgress": "21.56"
}
```

query tip  
queryChainPoint  
`runQueryTipCmd :: ()`  
executeLocalStateQueryExpr  

pprof

"unix"  
"unixgram"  
"unixpacket"  

Here's the docker image / details:  
https://github.com/IntersectMBO/cardano-node/pkgs/container/cardano-node


There are `Docker image builder` files `*.nix` all over the `cardano-node` repo.

I just want to know what base image they're using.  
`dockerTools.buildImage`

It builds a docker compatible image without actually using docker. That's  
pretty cool. I'd still prefer it actually use docker so it's cross-compatible  
without the dependency on nix, but still, interesteting that you can do that.  

https://ryantm.github.io/nixpkgs/builders/images/dockertools/

```
cd /usr/local/bin
ls -la
```

```
cardano-cli -> /nix/store/9188hm14aqkmg7flk8qi9dsfiikndwd9-cardano-cli-exe-cardano-cli-8.20.3.0/bin/cardano-cli
cardano-node -> /nix/store/9ph9pp8d21y17k95mhgk8jk2z2036dhj-cardano-node-exe-cardano-node-8.9.2/bin/cardano-node
```

Inside the docker container is an `entrypoint` script at `/usr/local/bin/entrypoint`  
which will start one of:
- `run-network`
- `run-node`
- `run-client`

In that order depending on if `NETWORK` env variable has been set, or if the first  
argument to the container is either `run` or `cli`.  

I just want to test that it's not something to do with container volume mounts  
that's stopping my golang script from connecting to the `node.socket`.  

```
docker build -t adc/cardano -f - . <<EOF
FROM ghcr.io/intersectmbo/cardano-node:8.9.2 AS source
FROM alpine AS runner
COPY --from=source /nix/store/9188hm14aqkmg7flk8qi9dsfiikndwd9-cardano-cli-exe-cardano-cli-8.20.3.0/bin/cardano-cli /usr/local/bin
COPY --from=source /nix/store/9ph9pp8d21y17k95mhgk8jk2z2036dhj-cardano-node-exe-cardano-node-8.9.2/bin/cardano-node /usr/local/bin
RUN apk add netcat-openbsd
EOF

docker build -t adc/cardano -f - . <<EOF
FROM ghcr.io/intersectmbo/cardano-node:8.9.2 AS source
FROM alpine AS runner
RUN apk add nc
COPY --from source:/
EOF
```

```
docker run --rm -it adc/cardano sh

/ # cardano-node version
rosetta error: failed to open elf at /nix/store/aw2fw9ag10wr9pf0qk4nk5sxi0q0bn56-glibc-2.37-8/lib/ld-linux-x86-64.so.2
 Trace/breakpoint trap
```

Okay, looks like I need ALL the nix store shit.

```
docker build -t adc/cardano -f - . <<EOF
FROM ghcr.io/intersectmbo/cardano-node:8.9.2 AS source
FROM alpine AS runner
COPY --from=source /nix/store /nix/store
RUN ln -s /nix/store/9188hm14aqkmg7flk8qi9dsfiikndwd9-cardano-cli-exe-cardano-cli-8.20.3.0/bin/cardano-cli /usr/local/bin/
RUN ln -s /nix/store/9ph9pp8d21y17k95mhgk8jk2z2036dhj-cardano-node-exe-cardano-node-8.9.2/bin/cardano-node /usr/local/bin/
RUN apk add netcat-openbsd socat
EOF
```

Beautiful! Now I have a full alpine container with the ability to install other  
packages with `apk add ...` but also with cardano-node/cli available.

```
docker run \
  --name cardano-mainnet \
  -d \
  -p 7731:3001 \
  -v ~/Desktop/maya-cardano-network/mainnet:/root/cardano \
  adc/cardano \
    cardano-node run \
    --topology /root/cardano/config/topology.json \
    --config /root/cardano/config/config.json \
    --database-path /root/cardano/db/ \
    --socket-path /root/cardano/node.socket \
    --host-addr "0.0.0.0" \
    --port "3001"
```

```
docker run \
  --rm \
  -it \
  -e CARDANO_NODE_NETWORK_ID=764824073 \
  -e CARDANO_NODE_SOCKET_PATH=node.socket \
  -w /root/cardano \
  -v ~/Desktop/maya-cardano-network/mainnet:/root/cardano \
  adc/cardano \
    cardano-cli query tip
```

Hmm works fine.

```
socat - UNIX-CONNECT:node.socket
socat - UNIX-CONNECT:/Users/adc/Desktop/maya-cardano-network/mainnet/node.socket
```

connection refused. same as golang. well I be.

```
docker run \
  --name tt \
  --rm \
  -p 0.0.0.0:7777:7777 \
  -it \
  -e CARDANO_NODE_NETWORK_ID=764824073 \
  -e CARDANO_NODE_SOCKET_PATH=node.socket \
  -w /root/cardano \
  -v ~/Desktop/maya-cardano-network/mainnet:/root/cardano \
  adc/cardano \
    sh

docker run \
  --name ttt \
  --rm \
  -it \
  -e CARDANO_NODE_NETWORK_ID=764824073 \
  -e CARDANO_NODE_SOCKET_PATH=node.socket \
  -w /root/cardano \
  -v ~/Desktop/maya-cardano-network/mainnet:/root/cardano \
  adc/cardano \
    sh
```

```
GOOS=linux GARCH=arm64 CGO_ENABLED=0 go build -o runme
docker cp ./runme tt:/tmp/
```

This link explains how to proxy the IPC over TCP.

https://medium.com/neoncat-io/how-to-communicate-with-the-cardano-node-on-a-remote-host-fe05dfd1bb94

I'd quite like to do that without

```
export CARDANO_NODE_NETWORK_ID=764824073
export CARDANO_NODE_SOCKET_PATH=node.socket
export NETWORK=mainnet

cardano-cli query tip
```

This was the query:  
> 4fcc6820000000478200a81980091a2d964a0919800a1a2d964a0919800b1a2d964a0919800c1a2d964a0919800d1a2d964a0919800e1a2d964a0919800f821a2d964a09f4198010821a2d964a09f4

Now just have to work out how that's encoded, and send over TCP instead of IPC.

```
GOOS=linux GARCH=arm64 CGO_ENABLED=0 go build -o cardano-cli-tcp && \
docker cp ./cardano-cli-tcp tt:/tmp/
```

For dash:

- getNetworkInfo

This is to retrieve the relay fee

- getBlockHash / getBlockVerboseTx

Used to retrieve all txs for a block

- getBestChainlock

Current height of the node

- getRawTransactionVerbose

- importAddressRescan


For cardano:

```
cardano-cli query tip

cardano-cli latest query slot-number 

cardano-cli latest query
cardano-cli latest query tx-mempool --mainnet info

cardano-cli latest query utxo --output-json --whole-utxo

cardano-cli latest query utxo --address 

cardano-cli latest query pool-state --all-stake-pools
```

First things first, wait for `syncProgress` in the tip query to be 100%.

To iterate the ledger for cardano:  
https://forum.cardano.org/t/how-to-query-based-on-block-height-and-listen-to-new-block-transaction-events-also-sockets-api/12472

Suggests using this DEPRECATED api:  
https://github.com/input-output-hk/cardano-rest  
/api/blocks/pages/total 

That repo points to:  
graphql, rosetta, submit-api or wallet.  

How does this site work? Is it open source?  
https://beta.explorer.cardano.org/en/  

Transactions in the last 24 hours

https://github.com/cardano-foundation/cf-explore

Oh god, the can of worms goes deeper.

Might need to understand how the ledger-sync works.

https://github.com/cardano-foundation/cf-ledger-sync

> The supported events are:
> 
> - **blockEvent** - Shelley and Post Shelley Blocks data (Includes everything transactions, witnesses..)
> - **rollbackEvent** - Rollback event with rollback point
> - **byronMainBlockEvent** - Byron Main Block data
> - **byronEbBlockEvent** - Byron Epoch Boundary Block data
> - **blockHeaderEvent** - Shelley and Post Shelley Block Header data
> - **genesisBlockEvent** - Genesis Block data
> - **certificateEvent** - Certificate data in a block
> - **transactionEvent** - Transaction data. One transactionEvent with all transactions in a block
> - **auxDataEvent** - Auxiliary data in a block
> - **mintBurnEvent** - Mint and Burn data in a block
> - **scriptEvent** - Script data of all transactions in a block

Right so the ledger connects directly to a node.

`STORE_CARDANO_HOST`

sanchonet-node.play.dev.cardano.org  
preview-node.play.dev.cardano.org  
backbone.cardano-mainnet.iohk.io  

all port 3001

`LedgerSyncApplication` has the `main` entrypoint.

```
./gradlew clean build -x test
```

> No matching variant of nu.studer:gradle-jooq-plugin:8.2.3 was found.

Oh Java. How I'd like to never work with you again.

Ah, a `Dockerfile`...

```
docker build -t ledgersync .
```

```
docker run \
  -p 8080:8080 \
  -e STORE_CARDANO_HOST=preprod-node.world.dev.cardano.org \
  -e STORE_CARDANO_PORT=30000 \
  -e STORE_CARDANO_PROTOCOLMAGIC=1 \
  -e NETWORK=preprod \
  ledgersync
```

Starts with `Aggregation App Ledger Sync` in overly massive lettering.

Then continues like this presumably until 10 milli blocks.

```
2024-05-16T06:19:18.390Z  INFO 1 --- [Ledger Sync Aggregation App] [ntLoopGroup-4-1] c.b.c.y.s.c.service.CursorServiceImpl    : # of blocks written: 100, Time taken: 141 ms
2024-05-16T06:19:18.390Z  INFO 1 --- [Ledger Sync Aggregation App] [ntLoopGroup-4-1] c.b.c.y.s.c.service.CursorServiceImpl    : Block No: 314648  , Era: Babbage
2024-05-16T06:19:18.541Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AddressTxAmountProcessor   : Time taken to save additional address_tx_amounts records : 1004, time: 73 ms
2024-05-16T06:19:18.546Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    : ### Starting account balance calculation upto block: 314748 ###
2024-05-16T06:19:18.568Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    :   Total Stake Address Balance records 112, Time taken to save: 6
2024-05-16T06:19:18.568Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    :   Time taken to delete stake address balance history: 0
2024-05-16T06:19:18.617Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    :   Total Address Balance records 797, Time taken to save: 54
2024-05-16T06:19:18.618Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    :   Time taken to delete address balance history: 1
2024-05-16T06:19:18.621Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    : ### Total balance processing and saving time 74 ###

2024-05-16T06:19:18.725Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.u.processor.AddressProcessor   : Address save size : 48, time: 257 ms
2024-05-16T06:19:18.725Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.u.processor.AddressProcessor   : Address Cache Size: 33137
2024-05-16T06:19:18.726Z  INFO 1 --- [Ledger Sync Aggregation App] [ntLoopGroup-4-1] c.b.c.y.s.c.service.CursorServiceImpl    : # of blocks written: 100, Time taken: 336 ms
2024-05-16T06:19:18.726Z  INFO 1 --- [Ledger Sync Aggregation App] [ntLoopGroup-4-1] c.b.c.y.s.c.service.CursorServiceImpl    : Block No: 314748  , Era: Babbage
2024-05-16T06:19:18.861Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AddressTxAmountProcessor   : Time taken to save additional address_tx_amounts records : 1066, time: 79 ms
2024-05-16T06:19:18.867Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    : ### Starting account balance calculation upto block: 314848 ###
2024-05-16T06:19:18.880Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    :   Total Stake Address Balance records 104, Time taken to save: 5
2024-05-16T06:19:18.881Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    :   Time taken to delete stake address balance history: 1
2024-05-16T06:19:18.939Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    :   Total Address Balance records 811, Time taken to save: 56
2024-05-16T06:19:18.939Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    :   Time taken to delete address balance history: 0
2024-05-16T06:19:18.943Z  INFO 1 --- [Ledger Sync Aggregation App] [               ] c.b.c.y.s.a.p.AccountBalanceProcessor    : ### Total balance processing and saving time 75 ###```
```

Just noticed the eras progressing:
- Shelly
- Allegra
- Mary
- Alonzo
- Babbage

Highest block on mainnet is currently `10321501`.

`BlockAggregatorServiceImpl.java:54`

https://store.yaci.xyz/  
Just want to see what this is:  
`com.bloxbean.cardano.yaci.store.events.BlockEvent`  

This looks like where all the useful data structures are:

`cf-ledger-sync/application/src/main/java/org/cardanofoundation/ledgersync/aggregate`

There's a `handleBlockEvent` function within `BlockEventListener` but no other  
mention of it in the entire project. Typical java magic.  

Yeah it's the `streamer-app` that connects to the node.

All relies on the yaci store / protocol.

> Events are published as Spring events

> error: invalid source release: 21
> Deprecated Gradle features were used in this build, making it incompatible with Gradle 9.0.

> Error: LinkageError occurred while loading main class
  org.cardanofoundation.ledgersync.streamer.LedgerSyncStreamerApplication
  java.lang.UnsupportedClassVersionError:
  org/cardanofoundation/ledgersync/streamer/LedgerSyncStreamerApplication has
  been compiled by a more recent version of the Java Runtime (class file
  version 65.0), this version of the Java Runtime only recognizes class file
  versions up to 61.0

Managed to get it running by going into `Project Structure` in Intellij IDEA and  
setting everything to use java 22.

> ERROR 28759 --- [main] com.zaxxer.hikari.pool.HikariPool : HikariPool-1 - Exception during pool initialization.

> org.springframework.beans.factory.BeanCreationException: Error creating bean
  with name 'flyway' defined in class path resource
  [org/springframework/boot/autoconfigure/flyway/FlywayAutoConfiguration$FlywayConfiguration.class]:
  Failed to instantiate [org.flywaydb.core.Flyway]: Factory method 'flyway'
  threw exception with message:
  org.springframework.jdbc.support.MetaDataAccessException: Could not get
  Connection for extracting meta-data

Riight.

It's got to be one of these 3 that I need to dive deeper into.

- AggregationApplication -> `aggregation-app`
- LedgerSyncApplication -> `application`
- LedgerSyncStreamerApplication -> `streamer-app`

Think this one can be ignored:

- SchedulerAppApplication -> `scheduler-app`

The two errors above come from both `LedgerSync.*` projects.

Lots of good structs and default network config files here:  
`bloxbean/cardano/yaci/store/common`

postgresql://localhost:5432  
rabbitmq on port 5672  

Could just run the docker-compose and see what happens

Why is it so difficult to find the section of code responsible for connecting  
to the actual node?

It's using netty. Might be worth trying netty-go.

Looking into `yaci-core`

Everything in the `network/` directory seems critical.

Especially `NodeClient` and `Session`.

https://start.spring.io/  
https://central.sonatype.com/artifact/com.bloxbean.cardano/yaci-store-spring-boot-starter/overview  

Created a little throwaway java spring app to see if I can get yaci running with  
some kind of minimal config.

> 2024-05-16T01:47:02.659-07:00  WARN 30106 --- [demo]
  [main] s.c.a.AnnotationConfigApplicationContext : Exception encountered
  during context initialization - cancelling refresh attempt:
  org.springframework.beans.factory.UnsatisfiedDependencyException: Error
  creating bean with name 'appcationEventListener' defined in URL
  [jar:file:/Users/adc/.gradle/caches/modules-2/files-2.1/com.bloxbean.cardano/yaci-store-core/0.1.0-rc3/228dd44bf0a90d6a7edead0ec7a418828ee2cd5e/yaci-store-core-0.1.0-rc3.jar!/com/bloxbean/cardano/yaci/store/core/service/AppcationEventListener.class]:
  Unsatisfied dependency expressed through constructor parameter 0: Error
  creating bean with name 'startService': Injection of autowired dependencies
  failed

It connected. Then exploded with some weird bean shit I don't understand.

Did they mis-spell 'application' or is that indended?

> Error creating bean with name 'startService': Injection of autowired dependencies failed

Which one?

```
mitmdump --help
mitmdump --raw-tcp -p 5656 
mitmdump -p 7732 -R localhost:7731 --tcp .*
mitmdump -p 7732 --rawtcp --rfile localhost:7731
mitmdump -p 7732 --rawtcp --rfile http://localhost:7731
mitmdump -p 7732 --mode reverse:http://localhost:7731
```

Need to look into `chainsync` and `blockfetch` folders.

#### 22.05.2024 Wednesday 6h

Right...

The issue at hand is getting ledger data out of the cardano node when the only  
way to read it is via the raw node to node protocol, the original code being in  
haskell, which I find extremely hard to read (especially with no  
intellisense).

Following on to another project, this time in java, which can index utxo data  
but which comes with with unavoidable message queue and database dependencies  
which make it difficult to run, understand, and obscures the underlying protocol.  

I was in the middle of digging deeper into the dependencies of this project with  
my own test demo.  

So `github.com/cardano-foundation/cf-ledger-sync` led to  
`https://store.yaci.xyz/docker`.  

I started by making a spring boot app in the same vein as the ledger sync. Then  
I deleted the spring boot entrypoint `runApplication<DemoApplication> (*args)`  
because it does a load of classic over-engineered design-principally-focused java  
scaffolding magic, reverse injecting a load of dependencies including one instance  
of the `TCPNodeClient` which I'm trying to test.

So I disabled that entrypoint, and am now trying to manually and declaritively  
wire things up to a point where I can listen to utxos. Then perhaps I can start  
reverse engineering things.  

My cardano mainnet node is running in docker and exposed on port `7731`.

I tried viewing the raw data with something like this:

```
mitmdump -vvv -p 7732 --mode reverse:http://localhost:7731 --raw --flow-detail=3
```

But it just hangs, where a direct connection results logs like:

```
15:23:09.970 [nioEventLoopGroup-2-1] INFO com.bloxbean.cardano.yaci.core.network.Session -- Connection established
15:23:09.972 [nioEventLoopGroup-2-1] INFO com.bloxbean.cardano.yaci.core.network.NodeClient -- Connected !!!
15:23:09.989 [nioEventLoopGroup-2-1] INFO com.bloxbean.cardano.yaci.core.protocol.handshake.HandshakeAgent -- Handshake Ok!!! AcceptVersion(versionNumber=13, versionData=N2NVersionData(initiatorOnlyDiffusionMode=true, peerSharing=0, query=false))
done15:23:15.001 [nioEventLoopGroup-2-1] INFO com.bloxbean.cardano.yaci.core.network.NodeClient -- Connection closed !!!
15:23:15.001 [nioEventLoopGroup-2-1] INFO com.bloxbean.cardano.yaci.core.network.Session -- Disposing the session !!!
15:23:15.001 [nioEventLoopGroup-2-1] WARN com.bloxbean.cardano.yaci.core.network.NodeClient -- Trying to reconnect !!!
15:23:15.003 [nioEventLoopGroup-2-2] INFO com.bloxbean.cardano.yaci.core.network.Session -- Connection established
15:23:15.003 [nioEventLoopGroup-2-2] INFO com.bloxbean.cardano.yaci.core.network.NodeClient -- Connected !!!
15:23:15.022 [nioEventLoopGroup-2-2] INFO com.bloxbean.cardano.yaci.core.protocol.handshake.HandshakeAgent -- Handshake Ok!!! AcceptVersion(versionNumber=13, versionData=N2NVersionData(initiatorOnlyDiffusionMode=true, peerSharing=0, query=false))
15:23:20.031 [nioEventLoopGroup-2-2] INFO com.bloxbean.cardano.yaci.core.network.NodeClient -- Connection closed !!!
15:23:20.031 [nioEventLoopGroup-2-2] INFO com.bloxbean.cardano.yaci.core.network.Session -- Disposing the session !!!
15:23:20.031 [nioEventLoopGroup-2-2] WARN com.bloxbean.cardano.yaci.core.network.NodeClient -- Trying to reconnect !!!
15:23:20.032 [nioEventLoopGroup-2-3] INFO com.bloxbean.cardano.yaci.core.network.Session -- Connection established
15:23:20.033 [nioEventLoopGroup-2-3] INFO com.bloxbean.cardano.yaci.core.network.NodeClient -- Connected !!!
15:23:20.040 [nioEventLoopGroup-2-3] INFO com.bloxbean.cardano.yaci.core.protocol.handshake.HandshakeAgent -- Handshake Ok!!! AcceptVersion(versionNumber=13, versionData=N2NVersionData(initiatorOnlyDiffusionMode=true, peerSharing=0, query=false))
```

Ad infinitum...

Oh yes, another issue was connecting to the `node.socket` as unix domain sockets  
are not accessible on the host side of the docker container. This is how the  
`cardano-cli` connects to the node but is only useful for sending and signing  
transactions - not for reading all ledger data. I suspect n2n will make this  
all redundant, still, I've created a little go binary that can expose the  
`cardano-cli` over a http rpc interface.  

Now back to n2n...

```
mitmdump -vvv -p 7732 --raw-tcp 7731 --flow-detail=3
mitmdump -vvv -p 7732 --mode transparent --flow-detail=3 localhost:7731
sudo mitmdump -vvv -p 7732 --mode transparent --flow-detail=3 localhost:7731
sudo mitmdump -vvv -p 7732 --mode reverse:tcp://localhost:7731 --flow-detail=3
sudo mitmdump -vvv -p 7732 --mode reverse:tcp://localhost:7731 --flow-detail=3 --set console_default_contentview="raw hex stream"
sudo mitmdump -vvv -p 7732 --mode reverse:tcp://localhost:7731 --flow-detail=3 --set dumper_default_contentview=hex
sudo mitmdump -vvv -p 7732 --mode reverse:tcp://localhost:7731 --set dumper_default_contentview=hex
sudo mitmdump -vvv -p 7732 --mode reverse:tcp://localhost:7731
sudo mitmdump -vvv -p 7732 --mode reverse:tcp://localhost:7731 --insecure
sudo mitmdump -vvv -p 7732 --mode reverse:tcp://localhost:7731 --ssl-insecure
```

> Insufficient privileges to access pfctl.

> Could not resolve original destination.

Fuck it. Just going to write my own proxy. Not sure why I'm not seeing ANY data  
from mitmdump.

Nice, now I'm getting what I need. This is what the handshake looks like:

```
<-- 000e266c000000598200aa04821a2d964a09f505821a2d964a09f506821a2d964a09f507821a2d964a09f508821a2d964a09f509821a2d964a09f50a821a2d964a09f50b841a2d964a09f500f40c841a2d964a09f500f40d841a2d964a09f500f4
--> 9e6451288000000c83010d841a2d964a09f500f4
```

That's all the client does at the moment. Is it consistent?

```
<-- 000df349000000598200aa04821a2d964a09f505821a2d964a09f506821a2d964a09f507821a2d964a09f508821a2d964a09f509821a2d964a09f50a821a2d964a09f50b841a2d964a09f500f40c841a2d964a09f500f40d841a2d964a09f500f4
--> a59938ce8000000c83010d841a2d964a09f500f4
```

The first 4 bytes of both request/response look like a time-based hash or maybe  
a message index/count. The rest show consistencies.

This was shown in the log:

```
AcceptVersion(versionNumber=13, versionData=N2NVersionData(initiatorOnlyDiffusionMode=true, peerSharing=0, query=false))
```

What patterns are there?

```
000df349
000000
598200aa
04
  821a2d964a09f5
05
  821a2d964a09f5
06
  821a2d964a09f5
07
  821a2d964a09f5
08
  821a2d964a09f5
09
  821a2d964a09f5
0a
  821a2d964a09f5
0b
  841a2d964a09f5
  00f4
0c
  841a2d964a09f5
  00f4
0d
  841a2d964a09f5
  00f4

-----

a59938ce
800000
0c
  8301
0d
  841a2d964a09f5
  00f4
```

```java
public class AcceptVersion implements Message {
    private long versionNumber;
    private VersionData versionData;
}

public class VersionData {
    protected long networkMagic;
}
```

```java
public class N2NVersionData extends VersionData {
    private Boolean initiatorOnlyDiffusionMode;
    private Integer peerSharing = 0;
    private Boolean query = Boolean.FALSE;
  ...
}

public class VersionData {
    protected long networkMagic;
}
```

What's the network magic for mainnet? Maybe I can isolate that.


`Message` led to `DefaultSerializer` which led to:

```java
    default byte[] serialize(T object) {
        DataItem di = serializeDI(object);
        return CborSerializationUtil.serialize(di);
    }
```

Oh yeahhhh, cbor. Okay running that through cyber chef.

If you take the first 8 bytes off the front, cbor deserialize works:

```
[
    0,
    {}
]
```

But isn't very promising.

`a59938ce8000000c83010d841a2d964a09f500f4`

```
[
    1,
    13,
    [
        764824073,
        true,
        0,
        false
    ]
]
```

Well, the response looks very promising. Why is the request so borked?

Okay well `cbor.me` handles that just fine:

```
[
  0,
  {
    4: [764824073, true],
    5: [764824073, true],
    6: [764824073, true],
    7: [764824073, true],
    8: [764824073, true],
    9: [764824073, true],
    10: [764824073, true],
    11: [764824073, true, 0, false],
    12: [764824073, true, 0, false],
    13: [764824073, true, 0, false]
  }
]
```

Okay could the first 4 bytes be a hash, and the next 4 bytes be a message discriminator?  
Or the first integer of the cbor message could be a discrim.  
We're getting somewhere here.  

https://github.com/fxamacker/cbor

For later: look at `TxUtil::calculateTxHash` for `blake2bHash256` usage and format.

Okay now I'm intercepting requests using the proxy and am able to parse  
arbitrary structs. Making sense of them will require a discriminator map, and  
I also need to trigger the java client to begin the block download / sync.  

`cd ~/code/yaci/core/src/main/java/com/bloxbean/cardano/yaci/core/protocol/localstate/queries`

`grep -ir "array.add(new UnsignedInteger(" .`

```
./StakeDistributionQuery.java:        array.add(new UnsignedInteger(5));
./GenesisConfigQuery.java:        queryArray.add(new UnsignedInteger(11));
./PoolDistrQuery.java:        array.add(new UnsignedInteger(21));
./EpochStateQuery.java:        queryArray.add(new UnsignedInteger(8));
./UtxoByAddressQuery.java:        array.add(new UnsignedInteger(6));
./EpochNoQuery.java:        queryArray.add(new UnsignedInteger(1));
./SystemStartQuery.java:        array.add(new UnsignedInteger(1));
./CurrentProtocolParamsQuery.java:        queryArray.add(new UnsignedInteger(3));
./DelegationsAndRewardAccountsQuery.java:        array.add(new UnsignedInteger(10)); //tag
./StakePoolParamsQuery.java:        array.add(new UnsignedInteger(17));
./ChainPointQuery.java:        array.add(new UnsignedInteger(3));
./StakeSnapshotQuery.java:        array.add(new UnsignedInteger(20));
./BlockHeightQuery.java:        array.add(new UnsignedInteger(2));
```

Interesting. Missing numbers and discriminator overlap on `1` and `3`.

Okay the ones that overlap have `wrapWithOuterArray(...)` in the `serialize`  
method which I can only assume solves that issue.

Going to need some trial and error here no doubt.

Back to java for a second, how do I get the block scanning going?

How does the `BlockEventListener` actually receive it's events?  
Spring framework event listener  
Should be a publisher  

> Ledger Sync Streamer app reads data from a Cardano node and publishes
  blockchain data to a messaging middleware like Kafka or RabbitMQ.

> To run the streamer app, you need to have following components:
> 1. Cardano Node or connect to a remote Cardano node
> 2. Database (PostgreSQL, MySQL, H2) : To store cursor/checkpoint data
> 3. Messaging middleware (Kafka, RabbitMQ) : To publish events

The streamer app itself only has one file and that just starts a spring app.  
So it's the yaci dependency that is magically all wired up via spring. The  
main takeaways I think are probably:  

`build.gradle`

```
implementation(libs.yaci.store.starter)
implementation(libs.yaci.store.remote.starter)
```

and the main file / entrypoint:

```java
public class LedgerSyncStreamerApplication   {

    //TODO -- Move these properties to Yaci Store
    @Value("${store.include-block-cbor:false}")
    private boolean includeBlockCbor;

    @Value("${store.include-txbody-cbor:false}")
    private boolean includeTxCbor;

    public static void main(String[] args) {
        SpringApplication.run(LedgerSyncStreamerApplication.class, args);
    }

    @PostConstruct
    public void init() {
        if (includeBlockCbor)
            YaciConfig.INSTANCE.setReturnBlockCbor(true);

        if (includeTxCbor)
            YaciConfig.INSTANCE.setReturnTxBodyCbor(true);
    }
}
```

So back to the yaci project and those `set...Cbor(...` methods.  
Just basic setter/getter.  
Dead end.

```
Gradle: com.bloxbean.cardano:yaci:0.3.0-beta13
Gradle: com.bloxbean.cardano:yaci-core:0.3.0-beta13
Gradle: com.bloxbean.cardano:yaci-helper:0.3.0-beta13
Gradle: com.bloxbean.cardano:yaci-store-client:0.1.0-rc3
Gradle: com.bloxbean.cardano:yaci-store-common:0.1.0-rc3
Gradle: com.bloxbean.cardano:yaci-store-core:0.1.0-rc3
Gradle: com.bloxbean.cardano:yaci-store-events:0.1.0-rc3
Gradle: com.bloxbean.cardano:yaci-store-spring-boot-starter:0.1.0-rc3
```

> Yaci is a Java-based Cardano mini-protocol implementation that allows users to
  connect to a remote or local Cardano node and interact with it in a variety
  of ways. With Yaci's simple APIs, you can listen to incoming blocks in
  real-time, fetch previous blocks, query information from a local node,
  monitor the local mempool, and submit transactions to a local node.

Well, at least I know I'm in the right place.

The store isn't part of the core project though.

git@github.com:bloxbean/yaci-core.git  
git@github.com:bloxbean/yaci-store.git  

Found the publisher! `ShelleyBlockEventPublisher`. There's also a  
`ByronBlockEventPublisher`.  

What is a spring framework `@Transactional` label? Relevant?  

Think it's part of the orm shenanigans. The annoying thing for me is finding the  
start point of where this crazy event / db write / node message is.  

Oh sweet Jesus could this be the clue I need:

```java
BlockSync blockSync = new BlockSync(node, nodePort, protocolMagic, Constants.WELL_KNOWN_PREPROD_POINT);
BlockFinder blockFinder = new BlockFinder(blockSync);
```

This class `com.bloxbean.cardano.yaci.core.common.Constants` has:
- `..._PROTOCOL_MAGIC`
- `WELL_KNOWN_..._POINT`

Well that all went brilliantly. Apart from the fact my proxy is no longer  
actually decoding the cbor messages correctly, my terminal was lit up with block  
data.  

It just clicked what the second part of the prefix is. It's the length of the  
next cbor message. I was thinking there'd be other delimiters for that like  
`\r\n` in http but who knows.  

Still not sure what the first bit is, can't be a discriminator because the  
same version response has a different prefix, must be either random based or  
time based. Could be a timestamp. Or time based sequence.  

Increased the buffer size on proxy, capped message throughput to 1/second  
so my eyes don't melt.  

```

<-- 0002e94b00000059
8200aa04821a2d964a09f505821a2d964a09f506821a2d964a09f507821a2d964a09f508821a2d964a09f509821a2d964a09f50a821a2d964a09f50b841a2d964a09f500f40c841a2d964a09f500f40d841a2d964a09f500f4

[0, {4: [764824073, true], 5: [764824073, true], 6: [764824073, true], 7: [764824073, true], 8: [764824073, true], 9: [764824073, true], 10: [764824073, true], 11: [764824073, true, 0, false], 12: [764824073, true, 0, false], 13: [764824073, true, 0, false]}]

--> a03a4bac8000000c
83010d841a2d964a09f500f4

[1, 13, [764824073, true, 0, false]]

<-- 0002aea000080005
82001904d20002ad650002002b820481821a00fd1fc158204e9bbbb67e3ae262133d94c3da5bffce7b1127fc436e7433b87668dba34c354a

[0, 1234]

--> a06824ad80020058
8305821a00fd1fc158204e9bbbb67e3ae262133d94c3da5bffce7b1127fc436e7433b87668dba34c354a82821a053c22005820bb1f48efacd09174088ab314633f46d84b830c73c2a78e0de94a40ef174208e21a00826365

[5, [16588737, h'4e9bbbb67e3ae262133d94c3da5bffce7b1127fc436e7433b87668dba34c354a'], [[87826944, h'bb1f48efacd09174088ab314633f46d84b830c73c2a78e0de94a40ef174208e2'], 8545125]]

--> a068261480080005
82011904d2

[1, 1234]

<-- 000025ff0002002b
820481821a053c22005820bb1f48efacd09174088ab314633f46d84b830c73c2a78e0de94a40ef174208e2

[4, [[87826944, h'bb1f48efacd09174088ab314633f46d84b830c73c2a78e0de94a40ef174208e2']]]

--> a095eb6d80020058
8305821a053c22005820bb1f48efacd09174088ab314633f46d84b830c73c2a78e0de94a40ef174208e282821a053c2f325820667dfdedeebab4cd862d98accbeeadaf187ca38e773372f94ea4cf51b8586e941a00826403

[5, [87826944, h'bb1f48efacd09174088ab314633f46d84b830c73c2a78e0de94a40ef174208e2'], [[87830322, h'667dfdedeebab4cd862d98accbeeadaf187ca38e773372f94ea4cf51b8586e94'], 8545283]]

<-- 000f366300020002
8100

[0]

--> a0c3acbc80020058
8303821a053c22005820bb1f48efacd09174088ab314633f46d84b830c73c2a78e0de94a40ef174208e282821a053c3b345820d126943fc7d05eb7ea3e1b72f611f2035f13d7a0663379b0ba2865cdfdd7ad891a0082649b

[3, [87826944, h'bb1f48efacd09174088ab314633f46d84b830c73c2a78e0de94a40ef174208e2'], [[87833396, h'd126943fc7d05eb7ea3e1b72f611f2035f13d7a0663379b0ba2865cdfdd7ad89'], 8545435]]

<-- 000f40a900020002
8100

[0]

--> a0f174ec80020392
83028205d81859035b828a1a008263661a053c22265820bb1f48efacd09174088ab314633f46d84b830c73c2a78e0de94a40ef174208e258205efcfe47ba8846f293938fc3f44cb842ace6412c9f740ad91abfec3873505fb1582063f33c336233060d2fc23db0c2f849fe850f61c5dc61ad1c6e8d58eb58a405ea82584088c22d4795589692287e6a40857191fead67833e14dbe319e8399fd459e352d06a17c6ddae035f3197aa751fbb3125d50f0bcfde1c6f54d707711db0385d26dc5850d372fdafeb92791e0ea155d737f5e647e82690d73c32ec0b6a0ff1fdac5c0b02fdb9d3f4423ac9ecae2835684a50080a13407d85584f1f7f94e98ec87417f02b65afcb199faeae37f6184f2cdf83230f1992f55820a22f8e617b9bf4c25c8f2731e434904a8ae6d313ff4caec5529576692d6832408458203f78bc2833a71365020a34a6c7ad3a474336da52388b40374782eaaed43e6d7510190275584099e60c4802818bde870cb21c3d1857118dc20c7a6af52e9161f76c456601e40fbb82f3178b8a16f163b7a937

--> 712790e21e46ef88
1425b6b430f92a6a22c414088208005901c095841f5a7656973a78a3ca2383d671358eb0608c775f877d1bcd42cc69f94ae0d202c36091e02fc11316658fd372f6ecc10ee19e8745ab74979ed7b37099ab00af81e0f46cff773a0f37e398873a0aa52a55a10b23feadf1ccab0a97644f078d4bbb45583efd934a9fd52c5bf88878dcde3d5390db577562b37c8ab17005e95c5694597ba274f288fc54abdf3458dce28396ffcd7a77e55e9c09726df7a0bc83f07df4037be4401c93ec34423e7446b41277cdea4c96eb55ea22546f0ebab92a261bff9ca0b51e30cdbf569d26100198370de36778c3c241e8fa978f377d38f62f5e2b77e81fc99e326c97ea01e4d9557afe226ad7a8666bca56463ea1532e18332344c7fa5441b12e2a2e25e670501b55986f5cbed5fe7ce8ed400176414fa86c49b6400b24376f435066e81be1882f7b4305e84c6e2d6b45429b5cdc52e33d6dd1b19ccad60a5f4f519efe4d46f0baf97230d148684ea1db9bdb99c3706287297db22679a69a8b27c4e330d2e97ab3ae65de1daf12ae7f

20

--> ea87caa7c807ffdc
245834f386d878dc3a473cd1a053baa2eedc7e40d568a52a5788aacbd7378c77a70dd15cc472dfc12473c4d5147a6bc26d6976867884ee72f9a18201e505d35e82821a053c45a758208b5b1a16b12fe6da7584ec5ecb0a521635d3f20a22e0cf85bad97935a9daf7b91a00826523

-5
^C
Process finished with the exit code 130 (interrupted by signal 2:SIGINT)

```

Need to work out what's going on here.

Probably need to update my buffer logic so that it will append multiple smaller  
messages together, I think the n2n protocol has quite a small limit.  

That's a tomorrow problem.



