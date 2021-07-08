package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"github.com/joeqian10/neo3-gogogo/sc"
	"github.com/joeqian10/neo3-gogogo/tx"
	"github.com/joeqian10/neo3-gogogo/wallet"
)

func mint() (hash string, err error) {
	port := "http://seed1t.neo.org:20332"
	client := rpc.NewClient(port)

	var magic uint32 = 844378958 // change to your network magic number
	ps := helper.ProtocolSettings{
		Magic:          magic,
		AddressVersion: helper.DefaultAddressVersion,
	}
	w, err := wallet.NewNEP6Wallet("./dv.neo-wallet.json", &ps, nil, nil)
	if err != nil {
		return "Wallet file error", err
	}
	err = w.Unlock("qwerty")
	if err != nil {
		return "Wallet password error", err
	}

	// create a WalletHelper
	wh := wallet.NewWalletHelperFromWallet(client, w)

	// build script
	scriptHash, err := helper.UInt160FromString("0x9b851e83c1d46172fea6298be92a276b3cc784c6")
	if err != nil {
		return "Build script error", err
	}
	// if your contract method has parameters
	cp1 := sc.ContractParameter{
		Type:  sc.ByteArray,
		Value: []byte("name"),
	}
	cp2 := sc.ContractParameter{
		Type:  sc.ByteArray,
		Value: []byte("description"),
	}
	cp3 := sc.ContractParameter{
		Type:  sc.ByteArray,
		Value: []byte("image"),
	}
	script, err := sc.MakeScript(scriptHash, "mint", []interface{}{cp1, cp2, cp3})
	if err != nil {
		return "Mint error", err
	}

	// get balance of gas in your account
	balancesGas, err := wh.GetAccountAndBalance(tx.GasToken)
	if err != nil {
		return "GetAccountAndBalance error", err
	}

	// make transaction
	trx, err := wh.MakeTransaction(script, nil, []tx.ITransactionAttribute{}, balancesGas)
	if err != nil {
		return "Make transaction error", err
	}

	// sign transaction
	trx, err = wh.SignTransaction(trx, magic)
	if err != nil {
		return "Sign transaction error", err
	}

	// send the transaction
	rawTxString := crypto.Base64Encode(trx.ToByteArray())
	response := wh.Client.SendRawTransaction(rawTxString)
	if response.HasError() {
		// do something
		return "Send transaction error", errors.New("Send transaction error")
	}

	// hash is the transaction hash
	hash = trx.GetHash().String()
	return hash, nil
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})

	r.GET("/mint", func(c *gin.Context) {
		txHash, err := mint()
		fmt.Println(err)
		c.JSON(200, gin.H{
			"tx_hash": txHash,
			"error":   err,
		})
	})

	r.Run()
}
