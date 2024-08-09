package jsonc

import jsoniter "github.com/json-iterator/go"

func ToJsonIgnoreErr(v interface{}) string {
	bArr, _ := jsoniter.Marshal(v)
	return string(bArr)
}

func ToJson(v interface{}) (string, error) {
	bArr, err := jsoniter.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bArr), nil
}

func FromJsonIgnoreErr(jsonStr string, o interface{}) {
	_ = jsoniter.Unmarshal([]byte(jsonStr), &o)
	return
}

func FromJson(jsonStr string, o interface{}) error {
	err := jsoniter.Unmarshal([]byte(jsonStr), &o)
	if err != nil {
		return err
	}
	return nil
}
