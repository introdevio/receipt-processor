package models

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	DateLayout = "2006-01-02 15:04"
)

type PointsResponse struct {
	Points int `json:"points"`
}

type ScoringRule interface {
	Score(receipt *Receipt) int
}

type RetailerRule struct{}
type TotalsRule struct{}
type ItemCountRule struct{}
type DescriptionLengthRule struct{}
type DateRule struct{}

func (r RetailerRule) Score(receipt *Receipt) int {
	var score int
	retailer := receipt.Retailer

	for _, char := range retailer {
		if unicode.IsNumber(char) || unicode.IsLetter(char) {
			score++
		}
	}
	return score
}

func (r ItemCountRule) Score(receipt *Receipt) int {
	return (len(receipt.Items) / 2) * 5
}

func (r TotalsRule) Score(receipt *Receipt) int {
	total, _ := strconv.ParseFloat(receipt.Total, 64)

	//let's convert to cents
	cents := math.Round(total * 100)

	if int(cents)%100 == 0 {
		return 50
	}

	if int(cents)%25 == 0 {
		return 25
	}
	return 0
}

func (r DescriptionLengthRule) Score(receipt *Receipt) int {

	var score int
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			total, _ := strconv.ParseFloat(item.Price, 64)
			score += int(math.Ceil(total * 0.2))
		}
	}
	return score
}

func (r DateRule) Score(receipt *Receipt) int {
	var score int

	const dateFormat = "%s %s"
	beforeTime, _ := time.Parse(DateLayout, fmt.Sprintf(dateFormat, receipt.PurchaseDate, "16:00"))
	afterTime, _ := time.Parse(DateLayout, fmt.Sprintf(dateFormat, receipt.PurchaseDate, "14:00"))
	purchaseTime, _ := time.Parse(DateLayout, fmt.Sprintf(dateFormat, receipt.PurchaseDate, receipt.PurchaseTime))

	if purchaseTime.Day()%2 != 0 {
		score += 6
	}

	if purchaseTime.After(afterTime) && purchaseTime.Before(beforeTime) {
		score += 10
	}

	return score
}
