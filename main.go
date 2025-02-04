package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var pointsMap = make(map[string]int)

// getPointsByReceiptId responds with the number of points a receipt generated based on it's ID. Returns 400 if ID doesn't exist.
func getPointsByReceiptId(c *gin.Context) {
	id := c.Param("id")
	if points, ok := pointsMap[id]; ok {
		pointsObject := PointsObject{points}
		c.JSON(http.StatusOK, pointsObject)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"description": "No receipt found for that ID."})
	}
}

// postReceipts adds an receipt to the hashmap based on the JSON received in the request body.
func processReceipt(c *gin.Context) {
	var newReceipt Receipt

	// Call BindJSON to bind the received JSON to newReceipt and throw 400 if validation doesn't pass.
	if err := c.BindJSON(&newReceipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The receipt is invalid.", "error": err.Error()})
		return
	}

	// Create new unique ID
	idString := uuid.NewString()

	// Calculate points and store in map with key = ID
	pointsMap[idString] = calculatePoints(newReceipt)
	idObject := ID{idString}
	c.JSON(http.StatusOK, idObject)
}

func main() {
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("isValidPrice", IsValidPrice)
		_ = v.RegisterValidation("isValidRetailer", IsValidRetailer)
		_ = v.RegisterValidation("isValidDescription", IsValidDescription)
	}

	router.GET("/receipts/:id/points", getPointsByReceiptId)
	router.POST("/receipts/process", processReceipt)
	router.Run("localhost:8080")
}
