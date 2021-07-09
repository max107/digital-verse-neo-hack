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
	"net/url"
	"strings"
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

func getStackFromTx2(txHash string) (stack string, err error) {

	params := url.Values{}
	// params.Add("{ jsonrpc: 2.0, id: 1, method: getapplicationlog, params: [" + txHash + "] }", ``)
	params.Add("{ jsonrpc: 2.0, id: 1, method: getapplicationlog, params: [\"72cf181c6845c879be3012af31b539d6b622f8bf508b1b8d5faee0701f49c19c\"] }", ``)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "http://seed1t.neo.org:20332", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(result))

	return "", nil
}

func getStackFromTx(txHash string) (stack string, err error) {
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
		Params:   	[]string{"72cf181c6845c879be3012af31b539d6b622f8bf508b1b8d5faee0701f49c19c"},
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

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res)

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

		txHash, err := mint(name, description, fileUrl)
		txHash = "0x" + txHash

		fmt.Println(err)

		c.JSON(200, gin.H{
			"tx_hash": txHash,
			"url": explorerLinkTx + txHash,
			"error":   err,
		})
	})

	r.GET("/token_properties", func(c *gin.Context) {
		tokenId := c.Query("tokenId")

		txHash, err := getTokenProperties(tokenId)
		txHash = "0x" + txHash

		fmt.Println(err)

		c.JSON(200, gin.H{
			"tx_hash": txHash,
			"url": explorerLinkTx + txHash,
			"error":   err,
		})
	})

	r.GET("/total_supply", func(c *gin.Context) {

		txHash, err := totalSupply()
		fmt.Println(err)
		stack, err := getStackFromTx(txHash)
		fmt.Println(err)

		responseTxHash := "0x" + txHash

		c.JSON(200, gin.H{
			"tx_hash": responseTxHash,
			"url": explorerLinkTx + responseTxHash,
			"stack": stack,
			"error":   err,
		})
	})

	r.Run()
}
