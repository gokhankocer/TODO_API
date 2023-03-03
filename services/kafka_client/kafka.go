package kafka_client

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func Producer(topic string, message interface{}) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_BROKER")),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("error marshalling message: ", err)
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: messageBytes,
	})

	if err != nil {
		log.Println("error writing message to topic: ", err)
	} else {
		log.Println("message written to topic successfully: ", string(messageBytes))
	}
}

func Consume(ctx context.Context, topic string, callBack func(message kafka.Message)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   topic,
		GroupID: topic,
	})
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("Error reading message: ", err)
			continue
		}
		callBack(msg)
	}
}
