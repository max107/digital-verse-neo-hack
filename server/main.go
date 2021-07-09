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

const port = "http://seed1t.neo.org:20332"
const magic = 844378958
const walletPath = "./dv.neo-wallet.json"
const walletPassword = "qwerty"
const scriptHash = "0x9b851e83c1d46172fea6298be92a276b3cc784c6"
const explorerLink = "https://neo3.neotube.io"
const explorerLinkContract = explorerLink + "/contract/"
const explorerLinkAddress = explorerLink + "/address/"
const explorerLinkTx = explorerLink + "/block/transaction/"

func mint(name string, description string, url string) (hash string, err error) {
	// init arguments
	cp1 := sc.ContractParameter{
		Type:  sc.ByteArray,
		Value: []byte(name),
	}
	cp2 := sc.ContractParameter{
		Type:  sc.ByteArray,
		Value: []byte(description),
	}
	cp3 := sc.ContractParameter{
		Type:  sc.ByteArray,
		Value: []byte(url),
	}
	hash, err = invokeContract("mint", []interface{}{cp1, cp2, cp3})

	return hash, err
}

func invokeContract(methodName string, args []interface{}) (hash string, err error) {
	client := rpc.NewClient(port)

	ps := helper.ProtocolSettings{
		Magic:          magic,
		AddressVersion: helper.DefaultAddressVersion,
	}
	w, err := wallet.NewNEP6Wallet(walletPath, &ps, nil, nil)
	if err != nil {
		return "", err
	}
	err = w.Unlock(walletPassword)
	if err != nil {
		return "", err
	}

	// create a WalletHelper
	wh := wallet.NewWalletHelperFromWallet(client, w)

	// build script
	scriptHash, err := helper.UInt160FromString(scriptHash)
	if err != nil {
		return "", err
	}

	script, err := sc.MakeScript(scriptHash, methodName, args)
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
		return "", errors.New(response.Error.Message)
	}

	// transaction hash
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

	r.POST("/create_nft", func(c *gin.Context) {
		name := c.PostForm("name")
		description := c.PostForm("description")
		url := c.PostForm("url")

		txHash, err := mint(name, description, url)
		txHash = "0x" + txHash

		fmt.Println(err)

		c.JSON(200, gin.H{
			"tx_hash": txHash,
			"url": explorerLinkTx + txHash,
			"error":   err,
		})
	})

	r.Run()
}
