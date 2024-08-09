package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
	"net/http"
)

func (self *Client) BalanceOfERC20(r *http.Request, toContractAddr string, toAddr string) (*big.Int, error) {
	if common.IsHexAddress(toAddr) == false {
		return big.NewInt(0), fmt.Errorf("toaddr is not valid %v", toAddr)
	}

	if common.IsHexAddress(toContractAddr) == false {
		return big.NewInt(0), fmt.Errorf("contractAddr is not valid %v", toContractAddr)
	}

	balanceAddr := "0x70a08231000000000000000000000000" //erc20 合约方法 balanceof 地址
	destAddr := fmt.Sprintf("%s%s", balanceAddr, toAddr[2:])

	type Params struct {
		Result string      `json:"result"`
		Error  interface{} `json:"error"`
	}

	var params Params
	var pm = make(map[string]string)
	pm["data"] = destAddr
	pm["to"] = toContractAddr
	err := self.callNode(r, "eth_call", 0, []interface{}{pm, "latest"}, &params)
	if err != nil {
		return nil, err
	}
	if params.Error != nil {
		return nil, fmt.Errorf("%v", params.Error)
	}

	if params.Result == "" {
		return big.NewInt(0), fmt.Errorf(`Result is ""`)
	}

	//余额为0x时,ParseBig256会报错
	if params.Result == "0x" {
		params.Result = "0x0"
	}

	balance, b := math.ParseBig256(params.Result)
	if b == false {
		return big.NewInt(0), fmt.Errorf(`ParseBig256 faild`)
	}
	return balance, nil
}
