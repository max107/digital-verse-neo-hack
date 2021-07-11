package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"github.com/joeqian10/neo3-gogogo/sc"
	"github.com/joeqian10/neo3-gogogo/tx"
	"github.com/joeqian10/neo3-gogogo/wallet"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const port = "http://seed1t.neo.org:20332"
const magic = 844378958
const walletPath = "./dv.neo-wallet.json"
const walletPassword = "qwerty"
const scriptHash = "0x19d98abb558d15cb9b893a6c6b4f01b3aa380336"
const explorerLink = "https://neo3.neotube.io"
const explorerLinkContract = explorerLink + "/contract/"
const explorerLinkAddress = explorerLink + "/address/"
const explorerLinkTx = explorerLink + "/transaction/"



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

func getTokenProperties(tokenId string) (hash string, err error) {
	// init arguments
	cp1 := sc.ContractParameter{
		Type:  sc.ByteArray,
		Value: []byte(tokenId),
	}
	
	hash, err = invokeContract("properties", []interface{}{cp1})

	return hash, err
}

func getTokens() (hash string, err error) {

	hash, err = invokeContract("tokens", []interface{}{})

	return hash, err
}

func totalSupply() (hash string, err error) {
	hash, err = invokeContract("totalSupply", []interface{}{})
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
		return "1", err
	}
	err = w.Unlock(walletPassword)
	if err != nil {
		return "2", err
	}

	// create a WalletHelper
	wh := wallet.NewWalletHelperFromWallet(client, w)

	// build script
	scriptHash, err := helper.UInt160FromString(scriptHash)
	if err != nil {
		return "3", err
	}

	script, err := sc.MakeScript(scriptHash, methodName, args)
	if err != nil {
		return "4", err
	}

	// get balance of gas in your account
	balancesGas, err := wh.GetAccountAndBalance(tx.GasToken)
	if err != nil {
		return "5", err
	}

	// make transaction
	trx, err := wh.MakeTransaction(script, nil, []tx.ITransactionAttribute{}, balancesGas)
	if err != nil {
		return "6", err
	}

	// sign transaction
	trx, err = wh.SignTransaction(trx, magic)
	if err != nil {
		return "7", err
	}

	// send the transaction
	rawTxString := crypto.Base64Encode(trx.ToByteArray())
	response := wh.Client.SendRawTransaction(rawTxString)

	if response.HasError() {
		return "8", errors.New(response.Error.Message)
	}

	// transaction hash
	hash = trx.GetHash().String()
	return hash, nil
}

func getLogsFromTx(txHash string, wait bool) (stack string, err error) {
	if wait {
		time.Sleep(20 * time.Second) // wait until transaction is included in block...?
	}
	type RequestStruct struct {
		Jsonrpc    	string `json:"jsonrpc"`
		Id   		int `json:"id"`
		Method  	string  `json:"method"`
		Params  	[]string `json:"params"`
	}
	data := RequestStruct{
		Jsonrpc:    "2.0",
		Id:      	1,
		Method: 	"getapplicationlog",
		Params:   	[]string{ txHash },
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://seed1t.neo.org:20332", "application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func uploadFileToNeoFS(fileUrl string) (url string, err error) {
	// TODO upload from s3 and get local file path
	localFilePath := "./videos/test.mov"
	
	return "", nil
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
		fileUrl := c.PostForm("url")
		showTxLogsRequestValue := c.DefaultPostForm("show_tx_logs", "false")
		showTxLogs, err := strconv.ParseBool(showTxLogsRequestValue)
		if err != nil {
			fmt.Println(err)
		}
		txHash, err := mint(name, description, fileUrl)
		if err != nil {
			fmt.Println(err)
		}
		txLogs, err := getLogsFromTx(txHash, showTxLogs)
		if err != nil {
			fmt.Println(err)
		}

		responseTxHash := "0x" + txHash
		c.JSON(200, gin.H{
			"tx_hash": responseTxHash,
			"url": explorerLinkTx + responseTxHash,
			"stack": txLogs,
			"error":   err,
		})
	})

	r.POST("/token_properties", func(c *gin.Context) {
		tokenId := c.PostForm("tokenId")

		txHash, err := getTokenProperties(tokenId)
		if err != nil {
			fmt.Println(err)
		}

		txLogs, err := getLogsFromTx(txHash, true)
		if err != nil {
			fmt.Println(err)
		}

		responseTxHash := "0x" + txHash
		c.JSON(200, gin.H{
			"tx_hash": responseTxHash,
			"url": explorerLinkTx + responseTxHash,
			"logs": txLogs,
			"error":   err,
		})
	})

	r.POST("/upload_file_to_neofs", func(c *gin.Context) {
		fileUrl := c.PostForm("fileUrl")
		uploadedFileUrl, err := uploadFileToNeoFS(fileUrl)
		c.JSON(200, gin.H{
			"url": uploadedFileUrl,
			"error":   err,
		})
	})

	r.GET("/tokens", func(c *gin.Context) {

		txHash, err := getTokens()
		if err != nil {
			fmt.Println(err)
		}

		txLogs, err := getLogsFromTx(txHash, true)
		if err != nil {
			fmt.Println(err)
		}

		responseTxHash := "0x" + txHash
		c.JSON(200, gin.H{
			"tx_hash": responseTxHash,
			"url": explorerLinkTx + responseTxHash,
			"logs": txLogs,
			"error":   err,
		})
	})

	r.GET("/total_supply", func(c *gin.Context) {

		txHash, err := totalSupply()
		if err != nil {
			fmt.Println(err)
		}
		txLogs, err := getLogsFromTx(txHash, true)
		if err != nil {
			fmt.Println(err)
		}

		responseTxHash := "0x" + txHash
		c.JSON(200, gin.H{
			"tx_hash": responseTxHash,
			"url": explorerLinkTx + responseTxHash,
			"stack": txLogs,
			"error":   err,
		})
	})

	r.Run()
}
