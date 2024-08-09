package kafka

import (
	"github.com/segmentio/kafka-go"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var client *Client

type Client struct {
	once    sync.Once
	brokers []string
}

func NewClient(kafkaURL string) *Client {
	once.Do(func() {
		client = &Client{brokers: strings.Split(kafkaURL, ",")}
	})

	return client
}

func (c *Client) NewKafkaWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(c.brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		MaxAttempts:            3,               // 最大重试次数
		WriteTimeout:           3 * time.Second, // 写入超时
		AllowAutoTopicCreation: true,            // 自动创建topic
	}
}

// 适用场景：新订阅组重启后从头读取数据，旧订阅组重启后接着上次读取
func (c *Client) NewKafkaReaderStartFirst(topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		GroupID: groupID,
		Topic:   topic,
		//MinBytes:       10e3,              // 10KB
		MaxBytes: 10e6, // 10MB
		//CommitInterval: time.Second,      //多久自动commit一次offset
		StartOffset: kafka.FirstOffset, // 从哪个offset开始消费，当GroupID有值时生效
	})
}

// 适用场景：新旧订阅组重启后只读取最新数据
func (c *Client) NewKafkaReaderStartLast(topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		GroupID: groupID,
		Topic:   topic,
		//MinBytes:       10e3,              // 10KB
		MaxBytes: 10e6, // 10MB
		//CommitInterval: time.Second,      //多久自动commit一次offset
		StartOffset: kafka.LastOffset, // 从哪个offset开始消费，当GroupID有值时生效
	})
}
