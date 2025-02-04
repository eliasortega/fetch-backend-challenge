package main

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// IsValidPrice checks if the input is a valid price string
func IsValidPrice(fl validator.FieldLevel) bool {
	price := fl.Field().String()
	out, _ := regexp.MatchString(`^\d+\.\d{2}$`, price)
	return out
}

func IsValidRetailer(fl validator.FieldLevel) bool {
	retailer := fl.Field().String()
	out, _ := regexp.MatchString(`^[\w\s\-&]+$`, retailer)
	return out
}

func IsValidDescription(fl validator.FieldLevel) bool {
	description := fl.Field().String()
	out, _ := regexp.MatchString(`^[\w\s\-]+$`, description)
	return out
}
