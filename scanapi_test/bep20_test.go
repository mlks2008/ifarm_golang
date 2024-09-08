package htb_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func Test_bep20(t *testing.T) {
	var toAddress = "0x53844067aD2C47C2c64086634C18eAf8397CA0da"
	var fromAddrs = make(map[string]string)
	//bep20
	{
		url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=tokentx&address=%v&page=1&offset=10000&startblock=0&endblock=99999999&sort=asc&apikey=UAZ9V44Y7BYX5APWF814M8WNHNKU3YVQ7U", toAddress)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching the URL:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error fetching the URL:", err)
			return
		}

		type Respone struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Result  []struct {
				BlockNumber       string `json:"blockNumber"`
				TimeStamp         string `json:"timeStamp"`
				Hash              string `json:"hash"`
				Nonce             string `json:"nonce"`
				BlockHash         string `json:"blockHash"`
				From              string `json:"from"`
				ContractAddress   string `json:"contractAddress"`
				To                string `json:"to"`
				Value             string `json:"value"`
				TokenName         string `json:"tokenName"`
				TokenSymbol       string `json:"tokenSymbol"`
				TokenDecimal      string `json:"tokenDecimal"`
				TransactionIndex  string `json:"transactionIndex"`
				Gas               string `json:"gas"`
				GasPrice          string `json:"gasPrice"`
				GasUsed           string `json:"gasUsed"`
				CumulativeGasUsed string `json:"cumulativeGasUsed"`
				Input             string `json:"input"`
				Confirmations     string `json:"confirmations"`
			} `json:"result"`
		}

		var txs Respone
		err = json.Unmarshal(body, &txs)
		if err != nil {
			fmt.Println("Error fetching the URL:", err)
			return
		}
		for _, tx := range txs.Result {
			fromAddrs[tx.From] = ""
		}
	}
	//bnb
	{
		url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%v&startblock=0&endblock=99999999&page=1&offset=10000&sort=asc&apikey=UAZ9V44Y7BYX5APWF814M8WNHNKU3YVQ7U", toAddress)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching the URL:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error fetching the URL:", err)
			return
		}

		type Respone struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Result  []struct {
				BlockNumber       string `json:"blockNumber"`
				BlockHash         string `json:"blockHash"`
				TimeStamp         string `json:"timeStamp"`
				Hash              string `json:"hash"`
				Nonce             string `json:"nonce"`
				TransactionIndex  string `json:"transactionIndex"`
				From              string `json:"from"`
				To                string `json:"to"`
				Value             string `json:"value"`
				Gas               string `json:"gas"`
				GasPrice          string `json:"gasPrice"`
				Input             string `json:"input"`
				MethodId          string `json:"methodId"`
				FunctionName      string `json:"functionName"`
				ContractAddress   string `json:"contractAddress"`
				CumulativeGasUsed string `json:"cumulativeGasUsed"`
				TxreceiptStatus   string `json:"txreceipt_status"`
				GasUsed           string `json:"gasUsed"`
				Confirmations     string `json:"confirmations"`
				IsError           string `json:"isError"`
			} `json:"result"`
		}

		var txs Respone
		err = json.Unmarshal(body, &txs)
		if err != nil {
			fmt.Println("Error fetching the URL:", err)
			return
		}
		for _, tx := range txs.Result {
			fromAddrs[tx.From] = ""
		}
	}
	//写入文件
	{
		values := make([]string, 0, len(fromAddrs))
		for key, _ := range fromAddrs {
			values = append(values, key)
		}
		joinedValues := strings.Join(values, ",")

		file, err := os.Create("addr_bep20.txt")
		if err != nil {
			fmt.Println("Error creating the file:", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(joinedValues)
		if err != nil {
			fmt.Println("Error writing to the file:", err)
			return
		}
	}
	fmt.Println("ok", len(fromAddrs))
}
