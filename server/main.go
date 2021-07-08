package main

import (
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"github.com/joeqian10/neo3-gogogo/sc"
	"github.com/joeqian10/neo3-gogogo/tx"
	"github.com/joeqian10/neo3-gogogo/wallet"
	"github.com/gin-gonic/gin"
	"errors"
	"fmt"
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
		return "", err
	}
	err = w.Unlock("qwerty")
	if err != nil {
		return "", err
	}

	// create a WalletHelper
	wh := wallet.NewWalletHelperFromWallet(client, w)

	// build script
	scriptHash, err := helper.UInt160FromString("0x9b437260ae6a5938a858f66dd802fc399ec128df")
	if err != nil {
		return "", err
	}
	// if your contract method has parameters
	cp1 := sc.ContractParameter{
		Type:  sc.ByteArray,
		Value: []byte{},
	}
	script, err := sc.MakeScript(scriptHash, "mint", []interface{}{cp1})
	if err != nil {
		return "", err
	}

	// get balance of gas in your account
	balancesGas, err := wh.GetAccountAndBalance(tx.GasToken)
	if err != nil {
		return "", err
	}

	// make transaction
	trx, err := wh.MakeTransaction(script, nil, []tx.ITransactionAttribute{}, balancesGas)
	if err != nil {
		return "", err
	}

	// sign transaction
	trx, err = wh.SignTransaction(trx, magic)
	if err != nil {
		return "", err
	}

	// send the transaction
	rawTxString := crypto.Base64Encode(trx.ToByteArray())
	response := wh.Client.SendRawTransaction(rawTxString)
	if response.HasError() {
		// do something
		return "", errors.New("Send transaction error")
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
		tx_hash, err := mint()
		c.JSON(200, gin.H{
			"tx_hash": tx_hash,
			"error": err,
		})
	})

	r.Run()
}
