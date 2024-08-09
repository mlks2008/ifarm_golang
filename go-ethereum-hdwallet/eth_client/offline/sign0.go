package ethoffline

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"strings"
)

func Sign0(chainId *big.Int, from, to, prikey string, amount *big.Int, gasPrice *big.Int, gasLimit int64, nonce int64) (string, error) {
	if common.IsHexAddress(to) == false {
		return "", fmt.Errorf("to is invalid %s", to)
	}

	rawTx := types.NewTransaction(uint64(nonce), common.HexToAddress(to), amount, uint64(gasLimit), gasPrice, nil)

	privateKey, err := crypto.HexToECDSA(prikey)
	if err != nil {
		return "", fmt.Errorf("prikey crypto %v", err)
	}

	signedTx, err := types.SignTx(rawTx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		return "", fmt.Errorf("signtx %v", err)
	}

	data, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", fmt.Errorf("rlp.encode %v", err.Error())
	}
	return ToHex(data), nil
}

func Sign0ERC20(chainId *big.Int, from, to, contractAddr string, prikey string, amount *big.Int, gasPrice *big.Int, gasLimit int64, nonce int64) (string, error) {
	if common.IsHexAddress(to) == false {
		return "", fmt.Errorf("to is invalid %s", to)
	}

	erc20abi, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return "", fmt.Errorf("abi.json %v", err)
	}

	if data, err := erc20abi.Pack("transfer", common.HexToAddress(to), amount); err != nil {

		return "", fmt.Errorf("abi.pack %v", err)

	} else {

		contAddr := common.HexToAddress(contractAddr)

		rawTx := types.NewTransaction(uint64(nonce), contAddr, big.NewInt(0), uint64(gasLimit), gasPrice, data)

		privateKey, err := crypto.HexToECDSA(prikey)
		if err != nil {
			return "", fmt.Errorf("prikey crypto %v", err)
		}

		signedTx, err := types.SignTx(rawTx, types.NewEIP155Signer(chainId), privateKey)
		if err != nil {
			return "", fmt.Errorf("signtx %v", err)
		}

		b, err := rlp.EncodeToBytes(signedTx)
		if err != nil {
			return "", fmt.Errorf("rlp.encode %v", err.Error())
		}
		return ToHex(b), nil
	}
}
