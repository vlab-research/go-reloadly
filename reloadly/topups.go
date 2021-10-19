package reloadly

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type TransactionDate time.Time

func (t *TransactionDate) UnmarshalJSON(b []byte) error {
	format := "2006-01-02 15:04:05"
	s := strings.Trim(string(b), "\"")
	parsed, err := time.Parse(format, s)
	if err != nil {
		return err
	}
	*t = TransactionDate(parsed)
	return nil
}

func (t *TransactionDate) MarshalJSON() ([]byte, error) {
	format := "2006-01-02 15:04:05"
	s := time.Time(*t).Format(format)
	return []byte(s), nil
}

func (t *TransactionDate) MarshalCSV() ([]byte, error) {
	return t.MarshalJSON()
}

type TopupResponse struct {
	TransactionID               int64            `csv:"transactionId" json:"transactionId,omitempty"`
	OperatorTransactionID       string           `csv:"operatorTransactionId" json:"operatorTransactionId,omitempty"`
	CustomIdentifier            string           `csv:"customIdentifier" json:"customIdentifier,omitempty"`
	RecipientPhone              string           `csv:"recipientPhone" json:"recipientPhone,omitempty"`
	RecipientEmail              string           `csv:"recipientEmail" json:"recipientEmail,omitempty"`
	SenderPhone                 string           `csv:"senderPhone" json:"senderPhone,omitempty"`
	CountryCode                 string           `csv:"countryCode" json:"countryCode,omitempty"`
	OperatorID                  int64            `csv:"operatorId" json:"operatorId,omitempty"`
	OperatorName                string           `csv:"operatorName" json:"operatorName,omitempty"`
	Discount                    float64          `csv:"discount" json:"discount,omitempty"`
	DiscountCurrencyCode        string           `csv:"discountCurrencyCode" json:"discountCurrencyCode,omitempty"`
	RequestedAmount             float64          `csv:"requestedAmount" json:"requestedAmount,omitempty"`
	RequestedAmountCurrencyCode string           `csv:"requestedAmountCurrencyCode" json:"requestedAmountCurrencyCode,omitempty"`
	DeliveredAmount             float64          `csv:"deliveredAmount" json:"deliveredAmount,omitempty"`
	DeliveredAmountCurrencyCode string           `csv:"deliveredAmountCurrencyCode" json:"deliveredAmountCurrencyCode,omitempty"`
	TransactionDate             *TransactionDate `csv:"transactionDate" json:"transactionDate,omitempty"`
}

type RecipientPhone struct {
	CountryCode string `json:"countryCode,omitempty"`
	Number      string `json:"number,omitempty"`
}

type SenderPhone struct {
	CountryCode string `json:"countryCode,omitempty"`
	Number      string `json:"number,omitempty"`
}

type TopupRequest struct {
	RecipientPhone   *RecipientPhone `json:"recipientPhone,omitempty"`
	SenderPhone      *SenderPhone    `json:"senderPhone,omitempty"`
	OperatorID       int64           `json:"operatorId,omitempty"`
	Amount           float64         `json:"amount,omitempty"`
	CustomIdentifier string          `json:"customIdentifier,omitempty"`
}

type TopupsService struct {
	*Service
	autoDetect      bool
	suggestedAmount bool
	autoFallback    bool
	operator        *Operator
	country         string
	tolerance       float64
	error           error
}

func (s *Service) Topups() *TopupsService {
	return &TopupsService{s, false, false, false, nil, "", 0.0, nil}
}

func (s *TopupsService) New() *TopupsService {
	return s.Topups()
}

func (s *TopupsService) AutoDetect(country string) *TopupsService {
	s.autoDetect = true
	s.autoFallback = false
	s.country = country
	return s
}

func (s *TopupsService) GetSetOperator() *Operator {
	return s.operator
}

func (s *TopupsService) FindOperator(country, name string) *TopupsService {
	op, err := s.SearchOperator(country, name)
	s.operator = op
	s.error = err
	return s
}

func (s *TopupsService) Operator(operator *Operator) *TopupsService {
	s.operator = operator
	return s
}

func (s *TopupsService) SuggestedAmount(tolerance float64) *TopupsService {
	s.suggestedAmount = true
	s.tolerance = tolerance
	return s
}

func (s *TopupsService) AutoFallback() *TopupsService {
	s.autoFallback = true
	return s
}

func NormalizeNumber(number string) (string, error) {
	// normalize
	return number, nil
}

func checkRangeAmount(operator *Operator, amount float64) (float64, error) {
	min := operator.LocalMinAmount
	max := operator.LocalMaxAmount

	if amount >= min && amount <= max {
		upper_lim := amount / operator.Fx.Rate
		upper_lim = math.Ceil(upper_lim*100) / 100
		return upper_lim, nil
	}

	return 0, ReloadlyError{
		ErrorCode: "IMPOSSIBLE_AMOUNT",
		Message:   fmt.Sprintf("Operator %v has a minimum amount of %v and max of %v. Amount %v requested could not be fulfilled", operator.Name, min, max, amount),
	}
}

func pickAmount(amounts []SuggestedAmount, min float64, tolerance float64) (*SuggestedAmount, error) {
	sort.Slice(amounts, func(i, j int) bool { return amounts[i].Sent < amounts[j].Sent })

	for _, a := range amounts {
		if a.Sent >= min && a.Sent <= min+tolerance {
			return &a, nil
		}
	}

	return nil, errors.New("no amount found")
}

func GetSuggestedAmount(operator *Operator, amount float64, tolerance float64) (float64, error) {
	if operator.DenominationType == "RANGE" {
		return checkRangeAmount(operator, amount)
	}

	amounts := operator.SuggestedAmountsMap
	amt, err := pickAmount(amounts, amount, tolerance)

	if err != nil {
		err = ReloadlyError{
			ErrorCode: "IMPOSSIBLE_AMOUNT",
			Message:   fmt.Sprintf("Could not manage to find an amount of at least %v for operator %v with suggested amounts %v", amount, operator.Name, amounts),
		}
		return 0, err
	}

	return amt.Pay, nil
}

// create retry timeout function
// retries every Xmin for Y mins
// for a set of ErrorCodes:
// PHONE_RECENTLY_RECHARGED
// TRANSACTION_CANNOT_BE_PROCESSED_AT_THE_MOMENT
// PROVIDER_INTERNAL_ERROR
// SERVICE_TO_OPERATOR_TEMPORARILY_UNAVAILABLE

func tryAutoFallback(err error) bool {
	if e, ok := err.(APIError); ok {
		switch e.ErrorCode {

		case "TRANSACTION_REFUSED_BY_OPERATOR",
			"INVALID_RECIPIENT_PHONE",
			"INVALID_AMOUNT_FOR_OPERATOR":

			return true

		default:
			return false
		}
	}
	return false
}

func (s *TopupsService) Topup(mobile string, requested_amount float64) (*TopupResponse, error) {
	amount := requested_amount

	if s.error != nil {
		return nil, s.error
	}

	if s.autoDetect {
		op, err := s.OperatorsAutoDetect(mobile, s.country)
		if err != nil {
			return nil, err
		}
		s.operator = op
	}

	if s.operator == nil {
		return nil, ReloadlyError{"INVALID_CALL", "You must set an operator to call Topup"}
	}

	// TODO: this is poor naming given current behavior.
	// It's confusing the way tolerance is overloaded.
	// needs some rethinking.
	if s.suggestedAmount {
		a, err := GetSuggestedAmount(s.operator, requested_amount, s.tolerance)
		if err != nil {
			return nil, err
		}
		amount = a
	}

	req := &TopupRequest{
		RecipientPhone: &RecipientPhone{s.operator.Country.IsoName, mobile},
		OperatorID:     s.operator.OperatorID,
		Amount:         amount,
	}

	resp := new(TopupResponse)
	_, err := s.Request("POST", "/topups", req, resp)

	// add retries??
	if err == nil || s.autoFallback == false || !tryAutoFallback(err) {
		return resp, err
	}

	// try with auto detect!
	return s.AutoDetect(s.operator.Country.IsoName).Topup(mobile, requested_amount)
}
