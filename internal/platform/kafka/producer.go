package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	producer *kafka.Writer
}

func NewProducer(addr string, topic string) *Producer {
	return &Producer{
		producer: &kafka.Writer{
			Addr: kafka.TCP(addr),
			Topic: topic,
			Balancer: &kafka.LeastBytes{},
			
		},
	}
}

func (p *Producer) Publish(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return p.producer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: data,
	})
}