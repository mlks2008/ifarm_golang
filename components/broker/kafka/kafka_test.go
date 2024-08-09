package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"testing"
	"time"
)

func Test_kafka(t *testing.T) {
	//client := NewClient("10.40.10.2:9092")
	client := NewClient("127.0.0.1:9092")

	go func() {
		consumer := client.NewKafkaReaderStartFirst("qa_equipment_log_event", "group1")
		defer consumer.Close()

		for {
			msg, err := consumer.ReadMessage(context.Background())
			if err != nil {
				t.Error(err)
				return
			} else {
				t.Logf("consumer %s %+v", msg.Value, msg.Offset)
			}
			//consumer.CommitMessages(context.Background(), msg)
		}
	}()

	time.Sleep(time.Second * 3)

	var done = make(chan struct{}, 0)
	go func() {
		producer := client.NewKafkaWriter("qa_equipment_log_event")
		defer producer.Close()

		for i := 0; i < 2; i++ {
			msg := kafka.Message{
				//Key:   []byte(fmt.Sprintf("key-%v", "key")),
				Value: []byte(fmt.Sprintf("消息%v", i)),
			}
			err := producer.WriteMessages(context.Background(), msg)
			if err != nil {
				t.Error(err)
				return
			} else {
				t.Logf("producer %s", msg.Value)
			}
			time.Sleep(time.Second)
		}
		done <- struct{}{}
	}()

	<-done
	time.Sleep(time.Second * 2)
}

func Test_data(t *testing.T) {
	type Data struct {
		Action int64
		Value  interface{}
	}

	var data = Data{
		Action: 10000,
		Value: map[string]interface{}{
			"key1": 1,
			"key2": "a",
		},
	}

	bytes, err := json.Marshal(data)
	fmt.Println(fmt.Sprintf("%s %v", bytes, err))
}
