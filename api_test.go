package main

import (
	"bytes"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var successPayload = map[string]interface{}{
	"retailer":     "Target",
	"purchaseDate": "2022-01-01",
	"purchaseTime": "13:01",
	"items": []map[string]string{
		{
			"shortDescription": "Mountain Dew 12PK",
			"price":            "6.49",
		}, {
			"shortDescription": "Emils Cheese Pizza",
			"price":            "12.25",
		}, {
			"shortDescription": "Knorr Creamy Chicken",
			"price":            "1.26",
		}, {
			"shortDescription": "Doritos Nacho Cheese",
			"price":            "3.35",
		}, {
			"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
			"price":            "12.00",
		},
	},
	"total": "35.35",
}

func TestPostReceiptSuccess(t *testing.T) {
	jsonPayload, _ := json.Marshal(successPayload)
	url := "http://localhost:8080/receipts/process"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal()
	}
	defer resp.Body.Close()
	// verify 200 status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var output ID
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &output); err != nil {
		t.Fatal()
	}
	match, _ := regexp.MatchString(`^\S+$`, output.Id)
	// verify returns valid ID
	assert.Equal(t, true, match)
}

func TestPostReceiptInvalidRetailer(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["retailer"] = "+"
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	// verify 400 status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostReceiptInvalidPurchaseDate(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["purchaseDate"] = "123"
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	// verify 400 status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostReceiptInvalidPurchaseTime(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["purchaseTime"] = "123"
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	// verify 400 status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostReceiptEmptyItems(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["items"] = []map[string]string{}
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	// verify 400 status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostReceiptInvalidItems(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["items"] = []map[string]string{
		{
			"shortDescription": "abc",
			"price":            "0000000000",
		},
	}
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	// verify 400 status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostReceiptInvalidTotal(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["total"] = "-123.45"
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	// verify 400 status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetValidReceipt(t *testing.T) {
	jsonPayload, _ := json.Marshal(successPayload)
	postUrl := "http://localhost:8080/receipts/process"
	resp1, err := http.Post(postUrl, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal()
	}
	defer resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		t.Fatal()
	}

	var idObj ID
	body, _ := io.ReadAll(resp1.Body)
	if err := json.Unmarshal(body, &idObj); err != nil {
		t.Fatal()
	}

	getUrl := "http://localhost:8080/receipts/" + idObj.Id + "/points"
	resp2, err := http.Get(getUrl)
	if err != nil {
		t.Fatal()
	}
	defer resp2.Body.Close()

	var output PointsObject
	body2, _ := io.ReadAll(resp2.Body)
	if err := json.Unmarshal(body2, &output); err != nil {
		t.Fatal()
	}

	expectedPoints := PointsObject{Points: 28}

	assert.Equal(t, expectedPoints, output)
}

func TestGetInvalidReceipt(t *testing.T) {
	getUrl := "http://localhost:8080/receipts/2/points"
	resp, err := http.Get(getUrl)
	if err != nil {
		t.Fatal()
	}
	defer resp.Body.Close()
	// verify 400 status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TODO: add unit tests for processing functions
