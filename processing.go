package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

func calculatePoints(r Receipt) int {
	return ScoreRetailer(r.Retailer) + ScoreTotal(r.Total) + ScoreItems(r.Items) + ScoreDate(r.PurchaseDate) + ScoreTime(r.PurchaseTime)
}

// ScoreRetailer returns one point for every alphanumeric character in the retailer name.
func ScoreRetailer(retailer string) int {
	nonAlphaNumericRegex := regexp.MustCompile(`[\W]+`)
	cleanedRetailer := nonAlphaNumericRegex.ReplaceAllString(retailer, "")
	return len(cleanedRetailer)
}

// ScoreTotal takes the total amount paid as a string,
// returns 50 points if the total is a round dollar amount with no cents plus an additional 25 points if the total is a multiple of .25
func ScoreTotal(total string) int {
	points := 0
	totalCents, _ := strconv.Atoi(strings.Replace(total, ".", "", 1))
	if totalCents%100 == 0 {
		points += 50
	}
	if totalCents%25 == 0 {
		points += 25
	}
	return points
}

// ScoreItems takes the slice of Items,
// returns 5 points for every two items on the receipt plus additional points for each item with a trimmed description length that's a multiple of 3
// equal to 0.2 times its price, rounded up to the nearest integer
func ScoreItems(items []Item) int {
	points := 0
	points += 5 * (len(items) / 2)
	for _, item := range items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			floatPrice, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(floatPrice * 0.2))
		}
	}
	return points
}

// ScoreDate takes the date as a String (Assuming YYYY-MM-DD format),
// returns 6 points if the day in the purchase date is odd, 0 otherwise.
func ScoreDate(date string) int {
	if day, _ := strconv.Atoi(strings.Split(date, "-")[2]); day%2 == 1 {
		return 6
	}
	return 0
}

// ScoreTime takes the time as a String (Assuming 15:04 24hr format),
// returns 10 points if the time of purchase is after 2:00pm and before 4:00pm, 0 otherwise.
func ScoreTime(time string) int {
	cleanedTime, _ := strconv.Atoi(strings.Replace(time, ":", "", 1))
	if cleanedTime > 1400 && cleanedTime < 1600 {
		return 10
	}
	return 0
}
