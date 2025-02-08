# Switchly

## Work Log

### Summary

```
08.02.2024 Saturday   4h 30m


Total                04h 00m 
```

### 08.02.2024 Saturday

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
