package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/repository"
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

func Consume(ctx context.Context, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   "mail",
		GroupID: "mail",
	})
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("Error reading message: ", err)
			continue
		}
		//log.Println("Successfully read message: ", string(msg.Value))
		userData := (msg.Value)
		var user entities.User
		err = json.Unmarshal(userData, &user)
		if err != nil {
			log.Println("Error parsing user data: ", err)
			continue
		}
		//log.Println("message", userData)
		activationLink := fmt.Sprintf("http://localhost:8080/api/activate/%d", user.ID)
		body := fmt.Sprintf("You account is now active and your ID is %d. Congrats!", user.ID)
		message := strings.Join([]string{body}, " ")
		//database.DB.Save(&user)
		go SendEmail(message, user.Email, activationLink)
		//log.Println(body, message, activationLink)

	}
}

func Activate(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User Id"})
		return
	}
	user, err := repository.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user.IsActive = true
	if err := repository.UpdateUser(uint(userID), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User successfully activated"})

}

func SendResetPasswordEmail(email string, resetPasswordToken string) {
	resetPasswordLink := fmt.Sprintf("http://localhost:8080/reset_password/%s", resetPasswordToken)
	message := fmt.Sprintf("Subject: Reset Password\n\nTo reset your password, please follow this link: %s", resetPasswordLink)
	response, err := SendEmail(message, email, "")
	if err != nil {
		log.Println("Error sending email:", err)
	}
	log.Println(response)
}
