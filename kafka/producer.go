package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/luopengift/golibs/channel"
	"github.com/luopengift/golibs/logger"
)

type Producer struct {
	addrs   []string
	topic   string
	message chan []byte      //将数据写入这个管道中
	channel *channel.Channel //并发写topic的协程控制
}

func NewProducer(addrs []string, topic string, maxThreads int) *Producer {
	return &Producer{
		addrs:   addrs,
		topic:   topic,
		message: make(chan []byte, 1000),
		channel: channel.NewChannel(maxThreads),
	}
}

func (self *Producer) ChanInfo() string {
	return self.channel.String()
}

func (self *Producer) Write(msg []byte) (int, error) {
	self.message <- msg
	return len(msg), nil
}

func (self *Producer) WriteToTopic() error {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	if err := config.Validate(); err != nil {
		logger.Error("<config error> %v", err)
		return err
	}

	producer, err := sarama.NewSyncProducer(self.addrs, config)
	if err != nil {
		logger.Error("<Failed to produce message> %v", err)
		return err
	}
	defer producer.Close()

	for {
		select {
		case message := <-self.message:
			self.channel.Add()
			go func(message []byte) {
				msg := &sarama.ProducerMessage{
					Topic:     self.topic,
					Partition: int32(-1),
					Key:       sarama.StringEncoder("key"),
					Value:     sarama.ByteEncoder(message),
				}
				if partition, offset, err := producer.SendMessage(msg); err != nil {
					logger.Error("<write to kafka error,partition=%v,offset=%v> %v", partition, offset, err)
				}
				self.channel.Done()
			}(message)
		}
	}
	return nil
}
