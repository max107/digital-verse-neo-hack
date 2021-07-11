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
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["001e84bced320c8e46afe613c89ad5b3a174260f74b5517ff1b26c1ded159dff"] }' http://seed1t.neo.org:20332 | json_pp


https://github.com/nspcc-dev/neofs-node
1. 
./neofs-cli -r st1.storage.fs.neo.org:8080 -w /Users/maxim/Desktop/Golang/src/digital-verse-neo-hack/wallets/dv.neo-wallet.json --address NSadsnNbMKd5DDdLESNFUZwxr9Zmyc9wBJ accounting balance
2. 
./neofs-cli -r st1.storage.fs.neo.org:8080 -w /Users/maxim/Desktop/Golang/src/digital-verse-neo-hack/wallets/dv.neo-wallet.json --address NSadsnNbMKd5DDdLESNFUZwxr9Zmyc9wBJ container create --policy "REP 3" --basic-acl public-read --await
./neofs-cli -r st1.storage.fs.neo.org:8080 -k KzbbA7tBNoSQHQiigtVSjcbX17R5p89Hb3LTCBhYvF85mZjHWj6n container create --policy "REP 3" --basic-acl public-read --await
container ID: 9i3ihnXrbHdN5f5TeG6BAgBi4uPmSeCKNSZsjmsHMvjE
3. 
./neofs-cli -r st1.storage.fs.neo.org:8080 -k KzbbA7tBNoSQHQiigtVSjcbX17R5p89Hb3LTCBhYvF85mZjHWj6n object put --file /Users/maxim/Desktop/Golang/src/digital-verse-neo-hack/tests/test.mov --cid 9i3ihnXrbHdN5f5TeG6BAgBi4uPmSeCKNSZsjmsHMvjE
./neofs-cli -r st1.storage.fs.neo.org:8080 -w /Users/maxim/Desktop/Golang/src/digital-verse-neo-hack/wallets/dv.neo-wallet.json --address NSadsnNbMKd5DDdLESNFUZwxr9Zmyc9wBJ object put --file /Users/maxim/Desktop/Golang/src/digital-verse-neo-hack/tests/test.mov --cid 9i3ihnXrbHdN5f5TeG6BAgBi4uPmSeCKNSZsjmsHMvjE
ID: G7HDtS5XUBnRMYjqcv7bqdb2XYsynp9mQDGdEgZruoVi

https://http.fs.neo.org/9i3ihnXrbHdN5f5TeG6BAgBi4uPmSeCKNSZsjmsHMvjE/G7HDtS5XUBnRMYjqcv7bqdb2XYsynp9mQDGdEgZruoVi

https://github.com/nspcc-dev/neo-go
1. Transfer tokens
   ./neo-go wallet nep17 transfer -w /Users/maxim/Desktop/Golang/src/digital-verse-neo-hack/wallets/dv.neo-wallet.json -r https://rpc1.n3.nspcc.ru:20331 --from NSadsnNbMKd5DDdLESNFUZwxr9Zmyc9wBJ --to NSEawP75SPnnH9sRtk18xJbjYGHu2q5m1W --token GAS --amount 75 hash160:NSadsnNbMKd5DDdLESNFUZwxr9Zmyc9wBJ
   ce61fb389691061c6139890b7cf3f440ecb74e0d9af5e1f59331efa00060149b
2. ./neo-go wallet export -w /Users/maxim/Desktop/Golang/src/digital-verse-neo-hack/wallets/dv.neo-wallet.json --decrypt NSadsnNbMKd5DDdLESNFUZwxr9Zmyc9wBJ
    WIF - KzbbA7tBNoSQHQiigtVSjcbX17R5p89Hb3LTCBhYvF85mZjHWj6n
