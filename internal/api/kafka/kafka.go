package kafka

import (
	"github.com/IBM/sarama"
)

// type KafkaService struct {
// 	log *slog.Logger
// }


// func NewKafka(log *slog.Logger) KafkaService {
// 	return KafkaService{
// 		log: log,
// 	}
// }

type Service struct {
	consumer sarama.Consumer
	producer sarama.SyncProducer
}


func New() (*Service, error){
	brokers := []string {"localhost:9092"}
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return nil, err
	}
	consumer, err := ConnectConsumer(brokers)
	if err != nil {
		return nil, err
	}

	return &Service{
		consumer: consumer,
		producer: producer,
	}, nil
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}


func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}


func (k *Service)PushUserToQueue(topic string, message []byte) error {
	defer k.producer.Close()

	// New message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// Send message
	_, _, err := k.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}