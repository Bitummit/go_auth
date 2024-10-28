package my_kafka

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

type Kafka struct {
	Conn *kafka.Conn
	Brokers []string
	Topic string
	Writer *kafka.Writer
}


func New(ctx context.Context, leaderAddress, topic string, partition int, brokers []string) (*Kafka, error){
	conn, err := kafka.DialLeader(ctx, "tcp", leaderAddress, topic, partition)
	if err != nil {
		return nil, fmt.Errorf("failed to dial leader: %w", err)
	}
	conn.SetReadDeadline(time.Now().Add(10*time.Second))

	return &Kafka{
		Conn: conn,
		Brokers: brokers,
		Topic: topic,
	}, nil

}

func (k *Kafka) InitProducer() {
	w := &kafka.Writer{
		Addr:     kafka.TCP(k.Brokers...),
		Topic:   k.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	k.Writer = w
}

func (k *Kafka) PushEmailToQueue(ctx context.Context, key string, value string) error {
	err := k.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		},
	)
	if err != nil {
		k.Writer.Close()
		return fmt.Errorf("sending message to queue: %w", err)
	}

	return nil
}