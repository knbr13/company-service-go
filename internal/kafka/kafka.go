package kafka

import (
	"github.com/IBM/sarama"
)

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
