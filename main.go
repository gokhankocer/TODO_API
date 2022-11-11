package main

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"net/http"


)

type Todo struct{
	ID string `json:"id"`
	Status string `json:"status"`
	Description string `json:"description"`
}

func main() {

}