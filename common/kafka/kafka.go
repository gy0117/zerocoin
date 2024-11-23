package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

const retries = 3

type KafkaConfig struct {
	Addr          string `json:"addr,optional"`
	WriteCap      int    `json:"writeCap,optional"`
	ReadCap       int    `json:"readCap,optional"`
	ConsumerGroup string `json:"ConsumerGroup,optional"`
}

type KafkaData struct {
	Topic string
	Key   []byte
	Data  []byte
}

type KafkaClient struct {
	w         *kafka.Writer
	r         *kafka.Reader
	readChan  chan KafkaData
	writeChan chan KafkaData
	conf      KafkaConfig
	closed    bool
	lock      sync.Mutex
}

func NewKafkaClient(conf KafkaConfig) *KafkaClient {
	return &KafkaClient{
		conf: conf,
	}
}

// Send 将数据发送给chan
func (k *KafkaClient) Send(data KafkaData) {
	defer func() {
		if err := recover(); err != nil {
			k.closed = true
		}
	}()
	k.writeChan <- data
	k.closed = false
}

// SendSync 同步发送
func (k *KafkaClient) SendSync(data KafkaData) error {
	k.w = &kafka.Writer{
		Addr:                   kafka.TCP(k.conf.Addr),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	msg := []kafka.Message{
		{
			Topic: data.Topic,
			Key:   data.Key,
			Value: data.Data,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return k.w.WriteMessages(ctx, msg...)
}

// Read 从chan中读数据，chan里面的数据，是从kafka中读取然后写入的
func (k *KafkaClient) Read() KafkaData {
	data := <-k.readChan
	return data
}

// StartWrite 从chan里面读数据，往kafka里写数据
func (k *KafkaClient) StartWrite() {
	k.w = &kafka.Writer{
		Addr:     kafka.TCP(k.conf.Addr),
		Balancer: &kafka.LeastBytes{},
		//Topic:                  "kline_1m_test_1",
		AllowAutoTopicCreation: true,
	}
	// 写chan的缓存大小，外界将数据写入writeChan
	k.writeChan = make(chan KafkaData, k.conf.WriteCap)
	go k.sendMsgToKafka()
}

// StartRead 从kafka中读数据，写到readChan中
func (k *KafkaClient) StartRead(topic string) *KafkaClient {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.conf.Addr},
		GroupID:  k.conf.ConsumerGroup,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		Topic:    topic,
	})
	client := NewKafkaClient(k.conf)
	client.r = r
	client.readChan = make(chan KafkaData, k.conf.ReadCap)
	go client.readMsgFromKafka()
	return client
}

func (k *KafkaClient) sendMsgToKafka() {
	for {
		select {
		case data := <-k.writeChan:
			msg := []kafka.Message{
				{
					Topic: data.Topic,
					Key:   data.Key,
					Value: data.Data,
				},
			}
			// 重试 + 发送失败后重新发送
			var success bool
			for i := 0; i < retries; i++ {
				err := k.w.WriteMessages(context.Background(), msg...)
				logx.Info("kafka | sendKafka | writeMessage first err: ", err)
				if err == nil {
					success = true
					break
				}
				if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					success = false
					continue
				}
				if err != nil {
					success = false
					logx.Info("kafka | sendKafka | writeMessage err: ", err.Error())
				}
			}
			if !success {
				// 三次WriteMessage重试失败后，重新写入writeChan
				k.Send(data)
			}
		}
	}
}

func (k *KafkaClient) readMsgFromKafka() {
	for {
		message, err := k.r.ReadMessage(context.Background())
		if err != nil {
			logx.Error(err)
			continue
		}
		data := KafkaData{
			Key:   message.Key,
			Data:  message.Value,
			Topic: message.Topic,
		}
		k.readChan <- data
	}
}

// RepeatPut 重新将数据放入channel，readChan是提供给外界读的，里面是kafka数据
func (k *KafkaClient) RepeatPut(data KafkaData) {
	k.readChan <- data
}
