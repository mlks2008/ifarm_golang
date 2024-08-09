package kafka

import (
	"components/broker/kafka"
	"components/common/consts"
	"context"
	"gamefi_equipment/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	kafkago "github.com/segmentio/kafka-go"
)

type Kafka struct {
	logEventProducer *kafkago.Writer
	log              *log.Helper
}

func NewKafka(c *conf.Kafka, logger log.Logger) *Kafka {

	logeventCli := kafka.NewClient(c.Logevent.GetAddr())

	logEventProducer := logeventCli.NewKafkaWriter(c.Logevent.GetTopic())

	return &Kafka{
		logEventProducer: logEventProducer,
		log:              log.NewHelper(logger),
	}
}

// 掉落日志
func (k *Kafka) PubAddEquipmentLog(data interface{}) {
	msg := k.getMsg(consts.EquipmentLogEvent_Action_Drop, data)
	if len(msg.Value) == 0 {
		return
	}

	err := k.logEventProducer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Warnf("PubDropEquipmentLog err %v", err)
	}
}

// 升级日志
func (k *Kafka) PubUpgradeEquipmentLog(data interface{}) {
	msg := k.getMsg(consts.EquipmentLogEvent_Action_Upgrade, data)
	if len(msg.Value) == 0 {
		return
	}

	err := k.logEventProducer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Warnf("PubDropEquipmentLog err %v", err)
	}
}
