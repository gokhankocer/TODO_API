package helper

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("SECRET"))

func GenerateJwt(user entities.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Hour * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(jwtKey)
}

func ValidateJWT(c *gin.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func CurrentUser(c *gin.Context) (entities.User, error) {
	err := ValidateJWT(c)
	if err != nil {
		return entities.User{}, err
	}
	token, _ := getToken(c)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))
	user, err := entities.FindUserById(userId)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func getToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})
	return token, err
}

func getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
