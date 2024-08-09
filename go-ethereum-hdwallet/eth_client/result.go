package eth

type Transaction struct {
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	Hash        string `json:"hash"`
	Input       string `json:"input"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Gas         string `json:"gas"`
	GasPrice    string `json:"gasPrice"`
	Type        string `json:"type"`
}

type BlockInfo struct {
	Hash         string        `json:"hash"`
	BlockNumber  string        `json:"number"`
	BaseFee      string        `json:"baseFeePerGas"`
	Transactions []Transaction `json:"transactions"`
}

type EthLog struct {
	Address         string   `json:"address"`
	Topics          []string `json:"topics"`
	Data            string   `json:"data"`
	BlockNumber     string   `json:"blockNumber"`
	TransactionHash string   `json:"transactionHash"`
}

type TransactionReceipt struct {
	To      string   `json:"to"`
	Logs    []EthLog `json:"logs"`
	Status  string   `json:"status"`
	GasUsed string   `json:"gasUsed"`
}
