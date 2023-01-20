package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/gokhankocer/TODO-API/entities"
	"github.com/segmentio/kafka-go"
)

func SendEmail(message string, toAddress string) (response bool, err error) {
	fromAddress := os.Getenv("EMAİL")
	fromEmailPassword := os.Getenv("PASSWORD")
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	var auth = smtp.PlainAuth("", fromAddress, fromEmailPassword, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, fromAddress, []string{toAddress}, []byte(message))
	if err == nil {
		return true, nil
	}
	return false, err
}

func Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   "new_user",
		GroupID: "email-new-users",
	})
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		userData := msg.Value
		var user entities.User
		err = json.Unmarshal(userData, &user)
		if err != nil {
			panic("could not parse userData " + err.Error())
		}
		subject := "Subject: Account created!\n\n"
		body := fmt.Sprintf("You account is now active and your ID is %d. Congrats!", user.ID)
		message := strings.Join([]string{subject, body}, " ")
		SendEmail(message, user.Email)

	}
}
