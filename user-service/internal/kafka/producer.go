package kafka

import (
	"log"
	"user-service/config"

	"github.com/IBM/sarama"
)

type Producer struct {
	AsyncProducer sarama.AsyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	p, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for err := range p.Errors() {
			log.Printf("Failed to produce message: %v", err)
		}
	}()

	return &Producer{AsyncProducer: p}, nil
}

func (kp *Producer) Publish(message []byte) {
	kp.AsyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: config.AppConfig.Kafka.Topic,
		Value: sarama.ByteEncoder(message),
	}
}