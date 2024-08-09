package parser

import jsoniter "github.com/json-iterator/go"

type JsonParser struct {
}

func (j *JsonParser) TypeName() string {
	return "json"
}

func (j *JsonParser) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}
