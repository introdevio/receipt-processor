package api

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/introdevio/receipt_processor/models"
	"github.com/introdevio/receipt_processor/store"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ReceiptApi defines and registers the routes for the /receipts API
type ReceiptApi struct {
	service *ReceiptService
}
type RuleEvaluator struct {
	rules []models.ScoringRule
}

type ReceiptService struct {
	store         store.ReceiptStore
	ruleEvaluator *RuleEvaluator
}

func NewReceiptsApi(echo *echo.Echo) *ReceiptApi {
	receiptStore := store.NewInMemoryReceiptStore()
	service := NewReceiptService(receiptStore)
	api := &ReceiptApi{service: service}
	api.registerRoutes(echo)
	api.registerValidators(echo)
	return api
}

func NewReceiptService(receiptStore store.ReceiptStore) *ReceiptService {
	return &ReceiptService{
		store:         receiptStore,
		ruleEvaluator: NewRuleEvaluator(),
	}
}

func (s *ReceiptService) Save(receipt *models.Receipt) *models.ReceiptResponse {
	accruedPoints := s.ruleEvaluator.CalculateScore(receipt)
	receipt.AccruedPoints = accruedPoints
	return &models.ReceiptResponse{Id: s.store.Save(receipt)}
}

func (s *ReceiptService) GetReceiptPoints(receiptId string) (*models.PointsResponse, error) {
	receipt, err := s.store.Retrieve(receiptId)
	if err != nil {
		return nil, err
	}
	return &models.PointsResponse{Points: receipt.AccruedPoints}, nil
}

func NewRuleEvaluator() *RuleEvaluator {
	return &RuleEvaluator{rules: []models.ScoringRule{
		models.RetailerRule{},
		models.TotalsRule{},
		models.ItemCountRule{},
		models.DescriptionLengthRule{},
		models.DateRule{},
	}}
}

func (re *RuleEvaluator) CalculateScore(receipt *models.Receipt) int {
	score := 0
	for _, rule := range re.rules {
		score += rule.Score(receipt)
	}
	return score
}

func (api *ReceiptApi) registerRoutes(echo *echo.Echo) {
	g := echo.Group("/receipts")
	g.POST("/process", api.postReceipt)
	g.GET("/:id/points", api.getPoints)
}

func (api *ReceiptApi) registerValidators(e *echo.Echo) {
	v := validator.New()
	err := v.RegisterValidation("date", models.DateValidator)
	handleValidationRegistrationError(e, err)
	err = v.RegisterValidation("time", models.TimeValidator)
	handleValidationRegistrationError(e, err)
	err = v.RegisterValidation("text", models.StringValidator)
	handleValidationRegistrationError(e, err)
	err = v.RegisterValidation("amount", models.AmountValidator)
	handleValidationRegistrationError(e, err)
	e.Validator = &models.ReceiptValidator{Validator: v}
}

func handleValidationRegistrationError(e *echo.Echo, err error) {
	if err != nil {
		e.Logger.Fatal("Could not register a validator", err)
	}
}

func (api *ReceiptApi) postReceipt(c echo.Context) error {
	r := new(models.Receipt)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, api.service.Save(r))
}

func (api *ReceiptApi) getPoints(c echo.Context) error {
	id := c.Param("id")
	result, err := api.service.GetReceiptPoints(id)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, result)
}
