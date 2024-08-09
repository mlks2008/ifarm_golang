package kafka

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	kafkago "github.com/segmentio/kafka-go"
)

type Data struct {
	Action int64
	Value  interface{}
}

func (k *Kafka) getMsg(action int64, data interface{}) kafkago.Message {
	var value = Data{
		Action: action,
		Value:  data,
	}

	var bvalue, err = json.Marshal(value)
	if err != nil {
		log.Warnf("getMsg err %v %v %v", action, data, err)
		return kafkago.Message{}
	}

	msg := kafkago.Message{
		Value: bvalue,
	}
	return msg
}
