package main

import (
	"bytes"
	"encoding/json"
	"maps"
	"net/http"
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
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostReceiptInvalidRetailer(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["retailer"] = "+"
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostReceiptInvalidPurchaseDate(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["purchaseDate"] = "123"
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostReceiptInvalidPurchaseTime(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["purchaseDate"] = "123"
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostReceiptEmptyItems(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["items"] = []map[string]string{}
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
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
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostReceiptInvalidTotal(t *testing.T) {
	payload := maps.Clone(successPayload)
	payload["total"] = "-123.45"
	jsonPayload, _ := json.Marshal(payload)
	url := "http://localhost:8080/receipts/process"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
