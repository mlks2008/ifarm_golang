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

func Test_erc20(t *testing.T) {
	var toAddress = "0x78df9620d1e00231ae4ed7ef50ba00445c1f30e9"
	var fromAddrs = make(map[string]string)
	//erc20
	{
		url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=tokentx&address=%v&page=1&offset=10000&startblock=0&endblock=99999999&sort=asc&apikey=7ECWSPTS2BIKJK2GXI8R9Y7BJT469Q3V69", toAddress)
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
	//eth
	{
		for i := 0; i < 2; i++ {
			var url string
			if i == 0 {
				url = fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%v&startblock=0&endblock=99999999&page=1&offset=10000&sort=asc&apikey=7ECWSPTS2BIKJK2GXI8R9Y7BJT469Q3V69", toAddress)
			} else {
				url = fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%v&startblock=0&endblock=99999999&page=1&offset=10000&sort=desc&apikey=7ECWSPTS2BIKJK2GXI8R9Y7BJT469Q3V69", toAddress)
			}
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
	}

	fmt.Println(len(fromAddrs))
	//写入文件
	{
		values := make([]string, 0, len(fromAddrs))
		for key, _ := range fromAddrs {
			values = append(values, key)
		}
		joinedValues := strings.Join(values, ",")

		file, err := os.Create("addr_erc20.txt")
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

// 排除地址
func Test_addr(t *testing.T) {
	//要排除的地址
	var addr1 = make(map[string]string)
	{
		data, _ := os.ReadFile("2021_addr_eth.txt")
		list := strings.Split(string(data), ",")
		for _, addr := range list {
			addr1[strings.ToLower(addr)] = ""
		}
	}

	//上次导出地址
	var addr2 = make(map[string]string)
	{
		data, _ := os.ReadFile("2021-erc20.csv")
		list := strings.Split(string(data), "\n")
		for _, addr := range list {
			if addr != "" {
				addr2[strings.ToLower(strings.Split(addr, ",")[1])] = ""
			}
		}
	}

	//最终排除地址：【排除的地址】去掉【上次已经导出的地址】
	var result = make(map[string]string)
	for addr, _ := range addr1 {
		//导出库中不存在，本次排除
		if _, ok := addr2[addr]; !ok {
			result[addr] = ""
		}
	}
	//写入文件
	{
		values := make([]string, 0, len(result))
		for key, _ := range result {
			values = append(values, key)
		}
		joinedValues := strings.Join(values, ",")

		file, err := os.Create("2021_eth_filter.txt")
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
	fmt.Println("ok", len(result))
}
