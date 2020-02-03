package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type consumer struct {
	reader *kafka.Reader
}

// Consumer read from the stream
type Consumer interface {
	// Read read from the stream
	Read(ctx context.Context, chMsg chan Message, chErr chan error)
}

func NewConsumer(brokers []string, topic string) Consumer {

	c := kafka.ReaderConfig{
		Brokers:         brokers,         // from:9092 
		Topic:           topic,           // chat
		MinBytes:        10e3,            // 10 KB
		MaxBytes:        10e6,            // 10 MB
		MaxWait:         1 * time.Second, // 1 Seconde
		ReadLagInterval: -1,
		GroupID:         Uuid(),
		StartOffset:     kafka.LastOffset,
	}

	return &consumer{kafka.NewReader(c)}
}

func (c *consumer) Read(ctx context.Context, chMsg chan Message, chErr chan error) {
	defer c.reader.Close()

	for {

		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			chErr <- errors.New(fmt.Sprintf("Error while reading a message: %v", err))
			continue
		}

		var message Message
		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			chErr <- err
		}

		chMsg <- message
	}
}
