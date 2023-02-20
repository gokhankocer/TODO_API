package kafka

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
		Topic:    "mail",
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
