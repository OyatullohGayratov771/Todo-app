package kafka

import (
	"log"
	"notification-service/config"

	"github.com/IBM/sarama"
)

func ConnKafka(cfg config.Config) sarama.Consumer {
	consumer, err := sarama.NewConsumer([]string{cfg.Kafka.Host + ":" + cfg.Kafka.Port}, nil)
	if err != nil {
		log.Fatalf("Kafka consumer yaratishda xatolik: %v", err)
	}
	return consumer
}
