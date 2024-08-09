package eth

type ETH_CurGasPrice struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		LastBlock       string `json:"LastBlock"`
		SafeGasPrice    string `json:"SafeGasPrice"`
		ProposeGasPrice string `json:"ProposeGasPrice"`
		FastGasPrice    string `json:"FastGasPrice"`
		SuggestBaseFee  string `json:"suggestBaseFee"`
		GasUsedRatio    string `json:"gasUsedRatio"`
	} `json:"result"`
}
type ETH_AvgGasPrice struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		UTCDate        string `json:"UTCDate"`
		UnixTimeStamp  string `json:"unixTimeStamp"`
		MaxGasPriceWei string `json:"maxGasPrice_Wei"`
		MinGasPriceWei string `json:"minGasPrice_Wei"`
		AvgGasPriceWei string `json:"avgGasPrice_Wei"`
	} `json:"result"`
}
