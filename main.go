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
	Status string   `json:"status"`
}

type orderUpdate struct {
	OrderId string `json:"orderId"`
	Status  string `json:"status"`
}

var orders = map[string]order{}

func postOrder(c *gin.Context) {
	var newOrder order
	if err := c.BindJSON(&newOrder); err != nil {
		log.Printf("Error binding JSON: %v\n", err)
		return
	}
	newOrder.Id = uuid.New().String()
	newOrder.Status = "queued"
	orders[newOrder.Id] = newOrder
	c.IndentedJSON(http.StatusCreated, newOrder)
}

func updateOrder(c *gin.Context) {
	var update orderUpdate
	if err := c.BindJSON(&update); err != nil {
		log.Printf("Error binding JSON: %v\n", err)
		return
	}
	foundOrder, orderWasFound := orders[update.OrderId]
	if !orderWasFound {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": fmt.Sprintf("Order with id %s not found", update.OrderId)})
		return
	}
	if (foundOrder.Status == "queued" && update.Status == "en-route") ||
		(foundOrder.Status == "en-route" && update.Status == "closed") {
		foundOrder.Status = update.Status
		orders[update.OrderId] = foundOrder
		c.IndentedJSON(http.StatusOK, foundOrder)
		return
	}
	c.JSON(http.StatusUnprocessableEntity, gin.H{"code": "UNPROCESSABLE_ENTITY", "message": fmt.Sprintf("Change of status from %s to %s is not allowed", foundOrder.Status, update.Status)})
}

// getEnv get key environment variable if exist otherwise return defaultValue
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
	router.PUT("/orders", updateOrder)
	host := getEnv("SERVER_ADDRESS", "localhost")
	port := getEnv("SERVER_PORT", "8080")
	if err := router.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Printf("Error starting server: %v\n", err)
		return
	}
}
