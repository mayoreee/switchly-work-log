# Switchly

## Work Log (Phase 2)

### Summary

```
08.02.2025 Saturday   4h 30m
09.02.2025 Sunday   2h 00m
15.03.2025 Saturday   8h 30m
16.03.2025 Sunday  9h 00m


Total                24h 00m 
```
### 15.03.2025 to 16.03.3035 
PR - https://github.com/SwitchlyProtocol/switchlynode/pull/2
Enabling Stellar Chain Connection to THORChain

Description
Added support for Stellar (XLM) blockchain integration in THORChain's Bifrost component, enabling native XLM asset transfers and chain observation.

Key Components Added

1. Stellar Bifrost Client (`bifrost/pkg/chainclients/stellar/stellar.go`)
- Implemented Stellar blockchain client with Horizon API integration
- Added support for mainnet and testnet networks
- Implemented transaction signing using TSS (Threshold Signature Scheme)
- Added account management and balance checking
- Implemented transaction broadcasting with retry mechanism
- Added solvency reporting for Stellar vaults

2. Block Scanner (`bifrost/pkg/chainclients/stellar/blockscanner.go`)
- Added Stellar block scanner for transaction observation
- Implemented block fetching and parsing
- Added support for memo filtering
- Added transaction confirmation handling

3. Chain Support (`common/chain.go`)
- Added Stellar chain constants and configurations
- Implemented gas asset decimal handling (7 decimals for XLM)
- Added network parameter configurations

4. Address Handling (`common/address.go`)
- Added Stellar address validation using strkey package
- Implemented address format checking for both public and test networks
- Added support for Stellar's Ed25519 public key format


Testing
Added test cases for:
- Address validation and conversion
- Transaction signing and broadcasting
- Block scanning and observation
- Public key handling
- Chain-specific configurations

Dependencies Added
- `github.com/stellar/go` - Official Stellar SDK
- `github.com/stellar/go/clients/horizonclient` - Horizon API client
- `github.com/stellar/go/txnbuild` - Transaction building
- `github.com/stellar/go/strkey` - Address encoding/decoding

Configuration
- Added Stellar network configuration (mainnet/testnet)
- Added default parameters for fees and minimum balances
- Added network passphrase handling

Security Considerations
- Implemented proper signature verification
- Added transaction validation checks
- Added proper error handling for network operations
- Implemented retry mechanism for failed broadcasts

Documentation
- Added inline documentation for new functions
- Updated configuration documentation
- Added examples for transaction handling


### 09.02.2025 Sunday

Worked on the second part of intergrating Stellar to the network which is described in the [Thorchain docs](https://gitlab.com/thorchain/thornode/-/blob/develop/docs/newchain.md?ref_type=heads).

Forked this - https://gitlab.com/thorchain/devops/node-launcher into https://github.com/SwitchlyProtocol/node-launcher


### 08.02.2025 Saturday
Focus for today was implementing Thorchain changes as the first part of intergrating Stellar to the network which is described in the [Thorchain docs](https://gitlab.com/thorchain/thornode/-/blob/develop/docs/newchain.md?ref_type=heads).

Made changes to go modules suchs as asset, chain, pubkey, address to support Stella/XLM. See changes:

https://github.com/SwitchlyProtocol/switchlynode/commit/33166dae76e355207f22dbff1f208f9bf3695c35
https://github.com/SwitchlyProtocol/switchlynode/commit/26cd0bd2327c8381ffcb1027b03cf69522f83714
https://github.com/SwitchlyProtocol/switchlynode/commit/d68908d0df9c4fb8580a5c9ba11d658bd02b5ecc
https://github.com/SwitchlyProtocol/switchlynode/commit/91f588871043e66b3982855db89529e84879c6b5
https://github.com/SwitchlyProtocol/switchlynode/commit/64435652b958479534cc1cacaea2e39f486cdd48
https://github.com/SwitchlyProtocol/switchlynode/commit/1e27e07999fa25f17f2f25e00eadbf94db444d39
https://github.com/SwitchlyProtocol/switchlynode/commit/540156ff8da355577d41b60947f6255aed34b2bc
https://github.com/SwitchlyProtocol/switchlynode/commit/347cf33a1c30c94282dbc94c7f4f786a05c92987
https://github.com/SwitchlyProtocol/switchlynode/commit/7012c05f38897fc3bea73c78382fcfab5be969d4
https://github.com/SwitchlyProtocol/switchlynode/commit/7012c05f38897fc3bea73c78382fcfab5be969d4


Also ran into a few issues getting thorchain and docker setup on a new M2 Mac evironment! Thanks, Apple!


Next, will focus on Node luahcer changes - https://gitlab.com/thorchain/devops/node-launcher


