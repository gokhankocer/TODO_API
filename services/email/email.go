package email

import (
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
	fromAddress := os.Getenv("EMAIL")
	if fromAddress == "" {
		log.Println("Error: EMAIL environment variable is not set.")
		return
	}
	fromEmailPassword := os.Getenv("PASSWORD")
	if fromEmailPassword == "" {
		log.Println("Error: PASSWORD environment variable is not set.")
		return
	}
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	body := fmt.Sprintf("%s\nActivation Link: %s", message, activationLink)
	msg := "From: " + fromAddress + "\n" +
		"To: " + toAddress + "\n" +
		"Subject: Account created!\n\n" +
		body
	log.Println("SMTP server:", smtpServer)
	log.Println("SMTP port:", smtpPort)
	log.Println("From address:", fromAddress)
	log.Println("To address:", toAddress)
	log.Println("Message:", msg)
	var auth = smtp.PlainAuth("", fromAddress, fromEmailPassword, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, fromAddress, []string{toAddress}, []byte(msg))
	if err == nil {
		log.Println("Email sent successfully")
		return true, nil
	}
	log.Printf("Error sending email: %v", err)
	return false, err
}

func SendResetPasswordEmail(emailAddress string, resetPasswordToken string) {
	resetPasswordLink := fmt.Sprintf("http://localhost:3000/reset_password/%s", resetPasswordToken)
	message := fmt.Sprintf("Subject: Reset Password\n\nTo reset your password, please follow this link: %s", resetPasswordLink)
	response, err := SendEmail(message, emailAddress, "")
	if err != nil {
		log.Println("Error sending email:", err)
	}
	log.Println(response)
}

func MailCallback(msg kafka.Message) {
	//log.Println("Successfully read message: ", string(msg.Value))
	userData := (msg.Value)
	var user entities.User
	err := json.Unmarshal(userData, &user)
	if err != nil {
		log.Println("Error parsing user data: ", err)
	}
	//log.Println("message", userData)
	activationLink := fmt.Sprintf("http://localhost:3000/api/activate/%d", user.ID)
	body := fmt.Sprintf("You account is now active and your ID is %d. Congrats!", user.ID)
	message := strings.Join([]string{body}, " ")
	//database.DB.Save(&user)
	go SendEmail(message, user.Email, activationLink)
	//log.Println(body, message, activationLink)
}
