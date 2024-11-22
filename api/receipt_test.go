package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/introdevio/receipt_processor/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var e = echo.New()
var receiptApi = NewReceiptsApi(e)

func TestPostReceipt(t *testing.T) {
	requests := []struct {
		name         string
		req          *http.Request
		expectedCode int
	}{
		{
			"valid receipt - 200 OK",
			createRequest(validReceipt),
			200,
		},
		{
			"missing field - 400 BadRequest",
			createRequest(missingfieldReceipt),
			400,
		},
		{
			"invalid date - 400 BadRequest",
			createRequest(invalidDateReceipt),
			400,
		},
		{
			"invalid time - 400 BadRequest",
			createRequest(invalidTimeReceipt),
			400,
		},
		{
			"invalid amount - 400 BadRequest",
			createRequest(invalidAmount),
			400,
		},
	}

	for _, tt := range requests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := e.NewContext(tt.req, rec)

			// Assertions
			err := receiptApi.postReceipt(c)
			if err != nil {
				assert.Equal(t, http.StatusOK, rec.Code)
			} else {
				var httpError *echo.HTTPError
				if errors.As(err, &httpError) {
					assert.Equal(t, tt.expectedCode, httpError.Code)
				}

			}
		})
	}
}

func TestGetPoints(t *testing.T) {
	rec := httptest.NewRecorder()
	r := createRequest(validReceipt)
	c := e.NewContext(r, rec)
	err := receiptApi.postReceipt(c)

	if err != nil {
		t.Fatal("Received Error when creating receipt")
	}

	var response models.ReceiptResponse

	if err != nil {
		t.Fatal("Received Error when reading response body")
	}

	err = json.Unmarshal(rec.Body.Bytes(), &response)

	if err != nil {
		t.Fatal("Received Error parsing response", err)
	}

	rec = httptest.NewRecorder()
	path := fmt.Sprintf("/receipts/%s/points", response.Id)
	r = httptest.NewRequest(http.MethodGet, path, nil)
	c = e.NewContext(r, rec)
	c.SetParamNames("id")
	c.SetParamValues(response.Id)

	if assert.NoError(t, receiptApi.getPoints(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var pointResponse models.PointsResponse
		err = json.Unmarshal(rec.Body.Bytes(), &pointResponse)
		if err != nil {
			t.Fatal("Could not parse response")
		}
		assert.Equal(t, 28, pointResponse.Points)
	}

}

func createRequest(payload string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req
}

var validReceipt = `
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}`

var missingfieldReceipt = `
{
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}`

var invalidDateReceipt = `
{
  "retailer": "Target"
  "purchaseDate": "invalid",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}`

var invalidTimeReceipt = `
{
  "retailer": "Target"
  "purchaseDate": "2022-01-01",
  "purchaseTime": "invalid",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}`

var invalidAmount = `
{
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35234"
}`
