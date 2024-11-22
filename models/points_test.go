package models

import (
	"testing"
)

func TestDateRule_Score(t *testing.T) {
	rule := DateRule{}
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{
			"Day of purchase is odd - award points",
			Receipt{
				PurchaseTime: "9:00",
				PurchaseDate: "2022-01-01",
			},
			6,
		},
		{
			"Day of purchase is even - no points awarded",
			Receipt{
				PurchaseTime: "09:00",
				PurchaseDate: "2022-01-02",
			},
			0,
		},
		{
			"Time of purchase is 3:00PM - points awarded",
			Receipt{
				PurchaseTime: "15:00",
				PurchaseDate: "2022-01-02",
			},
			10,
		},
		{
			"Time of purchase is 10:00AM and Even date - no points awarded",
			Receipt{
				PurchaseTime: "10:00",
				PurchaseDate: "2022-01-02",
			},
			0,
		},
		{
			"Day is Odd and Time of purchase is 3:00PM - points awarded",
			Receipt{
				PurchaseDate: "2022-01-01",
				PurchaseTime: "15:00",
			},
			16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Score(&tt.receipt); got != tt.want {
				t.Errorf("Total Rule returned incorrect score = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDescriptionLengthRule_Score(t *testing.T) {
	rule := DescriptionLengthRule{}
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{
			"Description Length is a multiple of 3 - award points",
			Receipt{Items: []Item{
				Item{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				},
			}},
			3,
		},
		{
			"Description Length is not multiple of 3 - award points",
			Receipt{Items: []Item{
				Item{
					ShortDescription: "Mountain Dew 12PK",
					Price:            "6.49",
				},
			}},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Score(&tt.receipt); got != tt.want {
				t.Errorf("Total Rule returned incorrect score = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemCountRule_Score(t *testing.T) {
	rule := ItemCountRule{}
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{
			"No Items in receipt",
			Receipt{Items: []Item{}},
			0,
		},
		{
			"2 Items - Points awarded once",
			Receipt{Items: []Item{
				{
					ShortDescription: "Mountain Dew 12PK",
					Price:            "6.49",
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				},
			}},
			5,
		},
		{
			"3 Items - Points awarded once",
			Receipt{Items: []Item{
				{
					ShortDescription: "Mountain Dew 12PK",
					Price:            "6.49",
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				},
			}},
			5,
		},
		{
			"4 Items - Points awarded twice (every 2 items)",
			Receipt{Items: []Item{
				{
					ShortDescription: "Mountain Dew 12PK",
					Price:            "6.49",
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				},
			}},
			10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Score(&tt.receipt); got != tt.want {
				t.Errorf("Total Rule returned incorrect score = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetailerRule_Score(t *testing.T) {
	rule := RetailerRule{}
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{
			"Retailer only letters - all characters are counted",
			Receipt{Retailer: "Target"},
			6,
		},
		{
			"Retailer with numbers in name - all characters are counted",
			Receipt{Retailer: "23andMe"},
			7,
		},
		{
			"Retailer with special characters - only alphanumeric are counter",
			Receipt{Retailer: "7-Eleven"},
			7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Score(&tt.receipt); got != tt.want {
				t.Errorf("Total Rule returned incorrect score = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTotalsRule_Score(t *testing.T) {
	rule := TotalsRule{}
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{
			"Total is round",
			Receipt{Total: "25.00"},
			50,
		},
		{
			"Total is multiple of 0.25",
			Receipt{Total: "25.50"},
			25,
		},
		{
			"Total is not multiple of 0.25 nor round",
			Receipt{Total: "25.77"},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Score(&tt.receipt); got != tt.want {
				t.Errorf("Total Rule returned incorrect score = %v, want %v", got, tt.want)
			}
		})
	}
}
