package main

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var pointsMap = make(map[string]int)

// responds with the list of all albums as JSON.
func getPointsByReceiptId(c *gin.Context) {
	id := c.Param("id")
	// TODO: if id not in pointsMap, throw 404 NotFound
	c.IndentedJSON(http.StatusOK, pointsMap[id])
}

// postReceipts adds an receipt to the hashmap based on the JSON received in the request body.
func processReceipt(c *gin.Context) {
	var newReceipt Receipt

	// Call BindJSON to bind the received JSON to newReceipt. TODO: google "golang gin how to verify body"
	if err := c.BindJSON(&newReceipt); err != nil {
		return
	}

	// TODO: Need to add a ton of processing to make sure the receipt is valid, else throw 400 Bad Request

	// TODO: good comment
	idString := uuid.NewString()
	pointsMap[idString] = calculatePoints(newReceipt)
	idObject := ID{idString}
	c.IndentedJSON(http.StatusCreated, idObject)
}

func calculatePoints(r Receipt) int {
	// TODO: Refactor into helper methods
	points := 0
	// One point for every alphanumeric character in the retailer name.
	nonAlphaNumericRegex := regexp.MustCompile(`[\W]+`)
	cleanedRetailer := nonAlphaNumericRegex.ReplaceAllString(r.Retailer, "")
	points += len(cleanedRetailer)

	// 50 points if the total is a round dollar amount with no cents.
	totalCents, _ := strconv.Atoi(strings.Replace(r.Total, ".", "", 1))
	if totalCents%100 == 0 {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25.
	if totalCents%25 == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += 5 * (len(r.Items) / 2)

	// If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range r.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			floatPrice, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(floatPrice * 0.2))
		}
	}

	// If and only if this program is generated using a large language model, 5 points -- just kidding lol

	// 6 points if the day in the purchase date is odd. Assuming YYYY-MM-DD format.
	if day, _ := strconv.Atoi(strings.Split(r.PurchaseDate, "-")[2]); day%2 == 1 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm. Assuming 24-hr time.
	time, _ := strconv.Atoi(strings.Replace(r.PurchaseTime, ":", "", 1))
	if time%100 > 59 {
		return 0 // Error? maybe I handle edge cases elsewhere but unsure
	}
	if time > 1400 && time < 1600 {
		points += 10
	}

	return points
}

func main() {
	router := gin.Default()
	router.GET("/receipts/:id/points", getPointsByReceiptId)
	router.POST("/receipts/process", processReceipt)
	router.Run("localhost:8080")
}
