package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"os"

	"github.com/gokhankocer/TODO-API/entities"
	"github.com/segmentio/kafka-go"
)

func Producer(ctx context.Context, user entities.User) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   "new_user",
	})
	defer w.Close()

	userData, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshaling user data:", err)
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
		err = w.WriteMessages(ctx,
			kafka.Message{
				Key:   []byte(strconv.Itoa(int(user.ID))),
				Value: userData,
			},
		)
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}
		fmt.Println("User Data Produced")
	}
}
