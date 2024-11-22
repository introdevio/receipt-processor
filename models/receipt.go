package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"time"
)

type ReceiptResponse struct {
	Id string `json:"id"`
}

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required,text"`
	Price            string `json:"price" validate:"required,amount"`
}

type Receipt struct {
	Retailer      string `json:"retailer" validate:"required,text"`
	PurchaseDate  string `json:"purchaseDate" validate:"required,date"`
	PurchaseTime  string `json:"purchaseTime" validate:"required,time"`
	Total         string `json:"total" validate:"required,amount"`
	Items         []Item `json:"items" validate:"required,dive"`
	AccruedPoints int    `json:-`
}

type ReceiptValidator struct {
	Validator *validator.Validate
}

func (cv *ReceiptValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

//Naming these generically as all strings and amounts
//follow the same patter in api spec

var TextRegex = regexp.MustCompile("^[\\w\\s\\-&]+$")
var AmountRegex = regexp.MustCompile("^\\d+\\.\\d{2}$")

func StringValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return TextRegex.MatchString(value)
}

func AmountValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return AmountRegex.MatchString(value)
}

func DateValidator(fl validator.FieldLevel) bool {
	purchaseDate := fl.Field().String()
	_, err := time.Parse("2006-01-02", purchaseDate)

	if err != nil {
		return false
	}
	return true
}

func TimeValidator(fl validator.FieldLevel) bool {
	purchaseTime := fl.Field().String()
	_, err := time.Parse("15:04", purchaseTime)

	if err != nil {
		return false
	}
	return true
}
