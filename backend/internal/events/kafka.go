package events

import (
	"context"
	"fmt"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Publisher interface {
	Publish(context.Context, string, []byte) error
	Close()
}
type KafkaPublisher struct{ client *kgo.Client }

func NewKafkaPublisher(brokers string) (*KafkaPublisher, error) {
	addresses := strings.FieldsFunc(brokers, func(r rune) bool { return r == ',' })
	if len(addresses) == 0 {
		return nil, fmt.Errorf("at least one Kafka broker is required")
	}
	client, err := kgo.NewClient(kgo.SeedBrokers(addresses...))
	if err != nil {
		return nil, fmt.Errorf("create Kafka client: %w", err)
	}
	return &KafkaPublisher{client: client}, nil
}
func (p *KafkaPublisher) Publish(ctx context.Context, topic string, value []byte) error {
	return p.client.ProduceSync(ctx, &kgo.Record{Topic: topic, Value: value}).FirstErr()
}
func (p *KafkaPublisher) Close() { p.client.Close() }
