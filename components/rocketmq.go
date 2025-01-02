package components

import (
	"context"
	"errors"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/sirupsen/logrus"
	"time"
)

var RocketMQProducer rocketmq.Producer
var RocketMQConsumer rocketmq.PushConsumer

type RocketMQConf struct {
	Enabled      bool                  `json:"enabled"`
	NameServer   []string              `json:"nameServer"`
	MessageModel consumer.MessageModel `json:"messageModel"`
	GroupName    string                `json:"groupName"`
	Retry        int                   `json:"retry"`
}

type InfraRocketMQProducer struct {
	Enabled  bool
	Producer rocketmq.Producer
}

func (p *InfraRocketMQProducer) SendSync(msg *primitive.Message) (*primitive.SendResult, error) {
	if !p.Enabled {
		return nil, errors.New("RocketMQ not enabled")
	}
	return RocketMQProducer.SendSync(context.Background(), msg)
}

func (p *InfraRocketMQProducer) SendAsync(mq func(ctx context.Context, result *primitive.SendResult, err error), msg ...*primitive.Message) error {
	if !p.Enabled {
		return errors.New("RocketMQ not enabled")
	}
	return RocketMQProducer.SendAsync(context.Background(), mq, msg...)
}

type InfraRocketMQConsumer struct {
	Enabled  bool
	Consumer rocketmq.PushConsumer
}

func (p *InfraRocketMQConsumer) RegConsumer(topic string, f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) {
	if !p.Enabled {
		return
	}
	c := RocketMQConsumer
	err := c.Subscribe(topic, consumer.MessageSelector{}, f)
	if err != nil {
		logrus.Errorf("注册topic: %s 失败, 1分钟后将重试, 错误信息: %s", topic, err)
		time.AfterFunc(1*time.Minute, func() {
			p.RegConsumer(topic, f)
		})
	}
}

func InitProducer(config RocketMQConf) InfraRocketMQProducer {
	prod, err := rocketmq.NewProducer(
		producer.WithNameServer(config.NameServer),
		producer.WithRetry(config.Retry),
		producer.WithGroupName(config.GroupName),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	if config.Enabled {
		err = prod.Start()
		if err != nil {
			logrus.Fatal(err)
		}
	}
	RocketMQProducer = prod
	return InfraRocketMQProducer{
		Enabled:  config.Enabled,
		Producer: prod,
	}
}

func InitConsumer(config RocketMQConf) InfraRocketMQConsumer {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(config.NameServer),
		consumer.WithRetry(config.Retry),
		consumer.WithConsumerModel(config.MessageModel),
		consumer.WithGroupName(config.GroupName),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	if config.Enabled {
		err = c.Start()
		if err != nil {
			logrus.Fatal(err)
		}
	}
	RocketMQConsumer = c
	return InfraRocketMQConsumer{
		Enabled:  config.Enabled,
		Consumer: c,
	}
}

func createTopic(config RocketMQConf, topic string) {
	h, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(config.NameServer)))
	if err != nil {
		logrus.Fatal(err)
	}
	err = h.CreateTopic(context.Background(), admin.WithTopicCreate(topic))
	if err != nil {
		logrus.Println(err)
	}
}
