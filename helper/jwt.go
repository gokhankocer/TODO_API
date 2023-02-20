package helper

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/repository"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("SECRET"))

func GenerateJwt(user entities.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(jwtKey)
}

func ValidateJWT(c *gin.Context) error {
	token, err := GetToken(c)
	if err != nil {
		return err
	}

	if isBlacklisted(c, token) {
		return errors.New("invalid token provided")
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
	token, _ := GetToken(c)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))
	user, err := repository.FindUserById(userId)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func GetToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Println(os.Getenv("SECRET"))
		//Bu kismin ne oldugunu anlamadim.
		return []byte(os.Getenv("SECRET")), nil
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

func isBlacklisted(c *gin.Context, token *jwt.Token) bool {

	res, err := database.RDB.Exists(c, token.Raw).Result()
	if err != nil {
		panic(err)
	}
	if res == 1 {
		return true
	}
	return false

}
