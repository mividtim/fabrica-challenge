package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
)

type order struct {
	Id     string   `json:"id"`
	UserId string   `json:"userId"`
	Items  []uint64 `json:"items"`
}

var orders []order

func postOrder(c *gin.Context) {
	var newOrder order
	if err := c.BindJSON(&newOrder); err != nil {
		log.Printf("Error binding JSON: %v\n", err)
		return
	}
	newOrder.Id = uuid.New().String()
	orders = append(orders, newOrder)
	c.IndentedJSON(http.StatusCreated, newOrder)
}

// getEnv get key environment variable if exist otherwise return defalutValue
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	router := gin.Default()
	router.POST("/orders", postOrder)
	host := getEnv("SERVER_ADDRESS", "localhost")
	port := getEnv("SERVER_PORT", "8080")
	if err := router.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Printf("Error starting server: %v\n", err)
		return
	}
}
