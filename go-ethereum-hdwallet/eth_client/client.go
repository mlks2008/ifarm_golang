package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"net/http"
	"strings"
)

type Client struct {
	chainId      int64
	nodeIndex    int
	nodeUrls     []string
	otherUrl     string
	validateUrl  string
	mainName     string
	autoGasLimit bool
	statCount    int64
	statTime     string
}

func NewClient(chainId int64, nodeUrl, otherUrl, validateUrl, mainName string, auto_gaslimit bool) *Client {
	client := Client{
		chainId:      chainId,
		nodeIndex:    0,
		nodeUrls:     strings.Split(nodeUrl, ","),
		otherUrl:     otherUrl,
		validateUrl:  validateUrl,
		mainName:     mainName,
		autoGasLimit: auto_gaslimit,
	}
	return &client
}

func (self *Client) GetChainId() (int64, error) {
	type Params struct {
		Result string      `json:"result"`
		Error  interface{} `json:"error"`
	}

	var params Params
	err := self.callNode(nil, "eth_chainId", 0, []interface{}{}, &params)
	if err != nil {
		return 0, err
	}
	if params.Error != nil {
		return 0, fmt.Errorf("%v", params.Error)
	}

	intbig, err := hexutil.DecodeBig(params.Result)
	if err != nil {
		return 0, err
	}

	return intbig.Int64(), nil
}

func (self *Client) GetLastHeight(r *http.Request) (int64, error) {
	type Params struct {
		Result string      `json:"result"`
		Error  interface{} `json:"error"`
	}

	var params Params
	err := self.callNode(r, "eth_blockNumber", 0, []interface{}{}, &params)
	if err != nil {
		return 0, err
	}
	if params.Error != nil {
		return 0, fmt.Errorf("%v", params.Error)
	}

	intbig, err := hexutil.DecodeBig(params.Result)
	if err != nil {
		return 0, err
	}

	return intbig.Int64(), nil
}

func (self *Client) GetBlockLogs(r *http.Request, hexHeight string) (map[string][]EthLog, error) {
	type Params struct {
		Result []EthLog    `json:"result"`
		Error  interface{} `json:"error"`
	}
	var params Params

	type Req struct {
		//BlockHash string `json:"blockHash"`
		FromBlock string   `json:"fromBlock"`
		ToBlock   string   `json:"toBlock"`
		Topics    []string `json:"topics"`
	}
	var req = Req{FromBlock: hexHeight, ToBlock: hexHeight, Topics: []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}}

	err := self.callNode(r, "eth_getLogs", 0, []interface{}{req}, &params)
	if err != nil {
		return nil, err
	}
	if params.Error != nil {
		return nil, fmt.Errorf("%v", params.Error)
	}

	var ethlogs = make(map[string][]EthLog)
	for _, log := range params.Result {
		if _, ok := ethlogs[log.TransactionHash]; !ok {
			ethlogs[log.TransactionHash] = make([]EthLog, 0)
		}
		ethlogs[log.TransactionHash] = append(ethlogs[log.TransactionHash], log)
	}
	return ethlogs, nil
}

func (self *Client) GetBlockByNumber(r *http.Request, blockHeigth string) (*BlockInfo, error) {
	type Params struct {
		Result BlockInfo   `json:"result"`
		Error  interface{} `json:"error"`
	}
	var params Params

	err := self.callNode(r, "eth_getBlockByNumber", 0, []interface{}{blockHeigth, true}, &params)
	if err != nil {
		return nil, err
	}
	if params.Error != nil {
		return nil, fmt.Errorf("%v", params.Error)
	}

	return &params.Result, nil
}

func (self *Client) GetTransactionByHash(r *http.Request, tx_id string) (*Transaction, error) {
	type Params struct {
		Result Transaction `json:"result"`
		Error  interface{} `json:"error"`
	}

	var params Params

	err := self.callNode(r, "eth_getTransactionByHash", 0, []interface{}{tx_id}, &params)
	if err != nil {
		return nil, err
	}
	if params.Error != nil {
		return nil, fmt.Errorf("%v", params.Error)
	}

	return &params.Result, nil
}

func (self *Client) GetTransactionReceipt(r *http.Request, tx_id string) (*TransactionReceipt, error) {
	type Params struct {
		Result TransactionReceipt `json:"result"`
		Error  interface{}        `json:"error"`
	}

	var params Params

	err := self.callNode(r, "eth_getTransactionReceipt", 0, []interface{}{tx_id}, &params)
	if err != nil {
		return nil, err
	}

	if params.Error != nil {
		return nil, fmt.Errorf("%v", params.Error)
	}

	return &params.Result, nil
}

func (self *Client) GetBalance(r *http.Request, addr string) (*big.Int, error) {
	type Params struct {
		Result hexutil.Big `json:"result"`
		Error  interface{} `json:"error"`
	}
	var params Params

	err := self.callNode(r, "eth_getBalance", 0, []interface{}{addr, "latest"}, &params)
	if err != nil {
		return nil, err
	}
	if params.Error != nil {
		return nil, fmt.Errorf("%v", params.Error)
	}
	return (*big.Int)(&params.Result), nil
}

func (self *Client) GasPrice(r *http.Request) (*big.Int, error) {
	type Params struct {
		Result hexutil.Big `json:"result"`
		Error  interface{} `json:"error"`
	}
	var params Params
	err := self.callNode(r, "eth_gasPrice", 0, []interface{}{}, &params)
	if err != nil {
		return nil, err
	}
	if params.Error != nil {
		return nil, fmt.Errorf("%v", params.Error)
	}
	return (*big.Int)(&params.Result), nil
}

func (self *Client) EstimateGas(r *http.Request, msg ethereum.CallMsg) (uint64, error) {
	if msg.Data == nil {
		callMsg, err := self.parseCallArg(msg)
		if err != nil {
			return 21000, nil
		}

		type Params struct {
			Result hexutil.Uint64 `json:"result"`
			Error  error          `json:"error"`
		}
		var params Params
		err = self.callNode(r, "eth_estimateGas", 0, []interface{}{callMsg}, &params)
		if err != nil {
			return 21000, err
		}
		if params.Error != nil {
			return 21000, fmt.Errorf("%v", params.Error)
		}

		if self.autoGasLimit == false {
			//FTX 钱包安全：既没有对接收方地址为合约地址进行任何限制。也没有对ETH原生Token的转账GAS Limit 进行限制，而是采用 estimateGas 方法评估手续费，这种方法导致GAS LIMIT大部分为500,000，超出默认21,000值的24倍。
			if params.Result > 23000 {
				return 23000, nil
			} else {
				return uint64(params.Result), nil
			}
		} else {
			return uint64(params.Result), nil
		}
	} else {
		var minGasLimit uint64 = 45000

		callMsg, err := self.parseCallArg(msg)
		if err != nil {
			return minGasLimit, err
		}

		type Params struct {
			Result hexutil.Uint64 `json:"result"`
			Error  error          `json:"error"`
		}
		var params Params
		err = self.callNode(r, "eth_estimateGas", 0, []interface{}{callMsg}, &params)
		if err != nil {
			return minGasLimit, err
		}
		if params.Error != nil {
			return minGasLimit, fmt.Errorf("%v", params.Error)
		}

		gasLimit := uint64(uint64(params.Result) * 110 / 100)

		if gasLimit == 0 {
			return minGasLimit, nil
		} else {
			//节点返回的gaslimit不太准，一些ERC20币(如RSR)可能会因为gaslimit过小，导致gas out
			if gasLimit < minGasLimit {
				return minGasLimit, nil
			} else {
				return gasLimit, nil
			}
		}
	}
}

func (self *Client) PendingNonceAt(r *http.Request, addr string) (uint64, error) {
	type Params struct {
		Result hexutil.Uint64 `json:"result"`
		Error  interface{}    `json:"error"`
	}
	var params Params
	err := self.callNode(nil, "eth_getTransactionCount", 0, []interface{}{addr, "pending"}, &params)
	if err != nil {
		return 0, err
	}
	if params.Error != nil {
		return 0, fmt.Errorf("%v", params.Error)
	}

	return uint64(params.Result), nil
}

func (self *Client) PendingNonceAt2(r *http.Request, addr string) (uint64, error) {
	type Params struct {
		Result hexutil.Uint64 `json:"result"`
		Error  interface{}    `json:"error"`
	}
	var params Params
	err := self.callNode(nil, "eth_getTransactionCount", 0, []interface{}{addr, "latest"}, &params)
	if err != nil {
		return 0, err
	}
	if params.Error != nil {
		return 0, fmt.Errorf("%v", params.Error)
	}

	return uint64(params.Result), nil
}

func (self *Client) SendTransaction(r *http.Request, sigdata string) (string, error) {
	type Params struct {
		Result string      `json:"result"`
		Error  interface{} `json:"error"`
	}

	var params Params
	err := self.callNode(nil, "eth_sendRawTransaction", 0, []interface{}{sigdata}, &params)
	if err != nil {
		return "", err
	}
	if params.Error != nil {
		return "", fmt.Errorf("%v", params.Error)
	}

	return params.Result, nil
}

func (self *Client) parseCallArg(msg ethereum.CallMsg) (interface{}, error) {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}

	return arg, nil
}
