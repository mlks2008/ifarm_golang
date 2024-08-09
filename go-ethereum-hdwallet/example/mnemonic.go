package main

import (
	"fmt"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/miguelmota/go-ethereum-hdwallet/eth_client"
	"github.com/shopspring/decimal"
	"github.com/tyler-smith/go-bip39"
	"log"
	"math/big"
)

func main() {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(128) //128,256
	mnemonic, _ := bip39.NewMnemonic(entropy)
	//mnemonic = ""

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	client := eth.NewClient(1, "https://mainnet.infura.io/v3/246a4f7d74db49dea312e204835eda06", "", "", "", false)

	for i := 0; i < 24; i++ {

		//"m/44'/60'/0'/0/0" "m/44'/60'/0'/0/1"
		var addrpath = fmt.Sprintf("m/44'/60'/0'/0/%v", i)

		path := hdwallet.MustParseDerivationPath(addrpath)
		account, err := wallet.Derive(path, false)
		if err != nil {
			log.Fatal(err)
			return
		}

		privateKey, err := wallet.PrivateKeyHex(account)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Printf("%v Address in hex: %s\n", i, account.Address.Hex()) // 0xC49926C4124cEe1cbA0Ea94Ea31a6c12318df947
		fmt.Printf("%v Private key in hex: %s\n", i, privateKey)
		continue

		ethBalance, err := client.GetBalance(nil, account.Address.Hex())
		if err != nil {
			fmt.Println("err %v", err)
			return
		}
		usdtBalance, err := client.BalanceOfERC20(nil, "0xdac17f958d2ee523a2206206994597c13d831ec7", account.Address.Hex())
		if err != nil {
			fmt.Println("err %v", err)
			return
		}

		if ethBalance.Cmp(big.NewInt(0)) > 0 || usdtBalance.Cmp(big.NewInt(0)) > 0 {
			fmt.Printf("%v Address in hex: %s\n", i, account.Address.Hex()) // 0xC49926C4124cEe1cbA0Ea94Ea31a6c12318df947
			fmt.Printf("%v Private key in hex: %s\n", i, privateKey)
			fmt.Println("ethBalance", decimal.NewFromBigInt(ethBalance, 0).Div(decimal.New(1, 18)), "usdtBalance", decimal.NewFromBigInt(usdtBalance, 0).Div(decimal.New(1, 6)))
			fmt.Println("")
		}
	}
}
