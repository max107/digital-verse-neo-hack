# TestNet Wallet:
1. 
wallets/dv.neo-wallet.json
NSadsnNbMKd5DDdLESNFUZwxr9Zmyc9wBJ
qwerty

# Deployed smart contracts
1.
NSadsnNbMKd5DDdLESNFUZwxr9Zmyc9wBJ - creator
NEO N3 TestNet:0 - network
0x97cee8dd846752091815245cdffa7ab209557f5692aa07ade8ed7f80ce949e02 - transaction
https://neo3.neotube.io/contract/0x146013865adf3f4d74e26aa16148badc879b6882 - info

# TestNet Tools

https://dora.coz.io
https://neo3.neotube.io/
https://neowish.ngd.network/#/

# API

1.
**{server}/create_nft**
POST
Params: 
name - string
description - string
url - string
Return:
tx_hash - string
url - string
error - string
2.
{server}/


# Tools
https://github.com/nspcc-dev/neo-go-sc-wrkshp
1. Tx result
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["146b3bcb2ce8c5d0b94d81bbf401342bd4d54be159f38be7ecdc1897bc834f04"] }' http://seed1t.neo.org:20332 | json_pp