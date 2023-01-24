package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/gokhankocer/TODO-API/entities"
	"github.com/segmentio/kafka-go"
)

func SendEmail(message string, toAddress string, activationLink string) (response bool, err error) {
	fromAddress := os.Getenv("EMAİL")
	fromEmailPassword := os.Getenv("PASSWORD")
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	body := fmt.Sprintf("%s\nActivation Link: %s", message, activationLink)
	msg := "From: " + fromAddress + "\n" +
		"To: " + toAddress + "\n" +
		"Subject: Account created!\n\n" +
		body

	var auth = smtp.PlainAuth("", fromAddress, fromEmailPassword, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, fromAddress, []string{toAddress}, []byte(msg))
	if err == nil {
		return true, nil
	}
	return false, err
}

func Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   "new_user",
	})
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				log.Println("Error reading message:", err)
				continue
			}
			userData := msg.Value
			fmt.Println("User Data Received: ", string(userData))
			var user entities.User
			err = json.Unmarshal(userData, &user)
			if err != nil {
				log.Println("Error parsing user data:", err)
				continue
			}
			activationLink := fmt.Sprintf("http://localhost:8080/api/activate/%d", user.ID)
			body := fmt.Sprintf("You account is now active and your ID is %d. Congrats!", user.ID)
			message := strings.Join([]string{body}, " ")
			fmt.Println("Message ready")
			go SendEmail(message, user.Email, activationLink)
			fmt.Println("Message Sent")
		}
	}
}
