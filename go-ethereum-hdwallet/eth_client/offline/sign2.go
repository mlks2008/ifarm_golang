package ethoffline

import (
	"fmt"
	"github.com/ethereum2/go-ethereum/accounts/abi"
	"github.com/ethereum2/go-ethereum/common"
	"github.com/ethereum2/go-ethereum/common/hexutil"
	"github.com/ethereum2/go-ethereum/core/types"
	"github.com/ethereum2/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"math/big"
	"strings"
)

func Sign2(chainId *big.Int, tipCap *big.Int, from, to, prikey string, amount *big.Int, gasPrice *big.Int, gasLimit int64, nonce int64) (string, error) {
	if common.IsHexAddress(to) == false {
		return "", fmt.Errorf("to is invalid %s", to)
	}

	//最大费
	var feeCap = decimal.NewFromBigInt(gasPrice, 0).BigInt()

	if tipCap.Cmp(feeCap) >= 0 {
		feeCap = decimal.NewFromBigInt(tipCap, 0).Mul(decimal.NewFromFloat(2)).BigInt()
	}

	pKey, err := crypto.HexToECDSA(prikey)
	if err != nil {
		return "", fmt.Errorf("prikey crypto %v", err)
	}

	toAddr := common.HexToAddress(to)

	accesses := types.AccessList{types.AccessTuple{
		Address:     toAddr,
		StorageKeys: []common.Hash{{0}},
	}}
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:    chainId,
		Nonce:      uint64(nonce),
		GasTipCap:  tipCap,
		GasFeeCap:  feeCap,
		Gas:        uint64(gasLimit),
		To:         &toAddr,
		Value:      amount,
		AccessList: accesses,
		Data:       nil})

	signer := types.LatestSignerForChainID(chainId)

	signedTx, err := types.SignTx(tx, signer, pKey)
	if err != nil {
		return "", fmt.Errorf("signtx %v", err.Error())
	}

	b, err := signedTx.MarshalBinary()
	if err != nil {
		return "", err
	}

	return hexutil.Encode(b), nil
}

func Sign2ERC20(chainId *big.Int, tipCap *big.Int, from, to, contractAddr string, prikey string, amount *big.Int, gasPrice *big.Int, gasLimit int64, nonce int64) (string, error) {
	if common.IsHexAddress(to) == false {
		return "", fmt.Errorf("to is invalid %s", to)
	}

	if common.IsHexAddress(contractAddr) == false {
		return "", fmt.Errorf("contractAddr is invalid %s", contractAddr)
	}

	//最大费
	var feeCap = decimal.NewFromBigInt(gasPrice, 0).BigInt()

	if tipCap.Cmp(feeCap) >= 0 {
		feeCap = decimal.NewFromBigInt(tipCap, 0).Mul(decimal.NewFromFloat(2)).BigInt()
	}

	erc20abi, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return "", fmt.Errorf("abi.json %v", err.Error())
	}

	if data, err := erc20abi.Pack("transfer", common.HexToAddress(to), amount); err != nil {

		return "", fmt.Errorf("abi.pack %v", err.Error())

	} else {

		pKey, err := crypto.HexToECDSA(prikey)
		if err != nil {
			return "", fmt.Errorf("prikey crypto %v", err)
		}

		contAddr := common.HexToAddress(contractAddr)

		accesses := types.AccessList{types.AccessTuple{
			Address:     contAddr,
			StorageKeys: []common.Hash{{0}},
		}}

		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:    chainId,
			Nonce:      uint64(nonce),
			GasTipCap:  tipCap,
			GasFeeCap:  feeCap,
			Gas:        uint64(gasLimit),
			To:         &contAddr,
			Value:      big.NewInt(0),
			AccessList: accesses,
			Data:       data})

		signer := types.LatestSignerForChainID(chainId)

		signedTx, err := types.SignTx(tx, signer, pKey)
		if err != nil {
			return "", fmt.Errorf("signtx %v", err.Error())
		}

		b, err := signedTx.MarshalBinary()
		if err != nil {
			return "", err
		}

		return hexutil.Encode(b), nil
	}
}
