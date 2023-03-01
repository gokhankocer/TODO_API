package main

import (
	//"errors"

	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/kafka_service/kafka"
	"github.com/gokhankocer/TODO-API/routers"
	"github.com/joho/godotenv"
	//"github.com/golang-jwt/jwt/v4"
	//"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

func main() {

	loadEnv()
	database.ConnectPostgres()

	migrateFlag := flag.Bool("migrate", false, "migrate argument")
	flag.Parse()
	if *migrateFlag {
		fmt.Println("Migrate started!")
		database.DB.Migrator().DropTable(&entities.User{}, &entities.Todo{})
		database.DB.Migrator().CreateTable(&entities.User{}, &entities.Todo{})
		return
	}

	database.ConnectRedis()
	go kafka.Consume(context.Background(), "mail")
	router := routers.SetupRouter()
	log.Fatal(router.Run("0.0.0.0:3000"))

}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file")
	} else {
		log.Print("Env successfully loaded")
	}

}
