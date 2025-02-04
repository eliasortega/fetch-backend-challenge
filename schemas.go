package main

type ID struct {
	Id string `json:"id"`
}

type PointsObject struct {
	Points int `json:"points"`
}

// Reciept represents data about a receipt
type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required,isValidDescription"`
	Price            string `json:"price" binding:"required,isValidPrice"`
}

// Receipt represents data about a receipt
type Receipt struct {
	Retailer     string `json:"retailer" binding:"required,isValidRetailer"`
	PurchaseDate string `json:"purchaseDate" binding:"required,datetime=2006-01-02"`
	PurchaseTime string `json:"purchaseTime" binding:"required,datetime=15:04"`
	Items        []Item `json:"items" binding:"gt=0,dive,required,required"`
	Total        string `json:"total" binding:"required,isValidPrice"`
}
