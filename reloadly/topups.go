package reloadly

import (
	"errors"
	"fmt"
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
	SenderPhone      *SenderPhone `json:"senderPhone,omitempty"`
	OperatorID       int64    `json:"operatorId,omitempty"`
	Amount           float64    `json:"amount,omitempty"`
	CustomIdentifier string `json:"customIdentifier,omitempty"`
}

type TopupsService struct {
	*Service
	autoDetect bool
	suggestedAmount bool
	autoFallback bool
	operator *Operator
	country string
	tolerance float64
	error error
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

func pickAmount(amounts []SuggestedAmount, min float64, tolerance float64) (*SuggestedAmount, error) {
	sort.Slice(amounts, func(i, j int) bool { return amounts[i].Sent < amounts[j].Sent})

	for _, a := range amounts {
		if a.Sent >= min && a.Sent <= min + tolerance {
			return &a, nil
		}
	}

	return nil, errors.New("no amount found")
}

func getSuggestedAmount(operator *Operator, amount float64, tolerance float64) (float64, error) {
	amounts := operator.SuggestedAmountsMap
	amt, err := pickAmount(amounts, amount, tolerance)

	if err != nil {
		err = ReloadlyError{
			ErrorCode: "IMPOSSIBLE_AMOUNT",
			Message: fmt.Sprintf("Could not manage to find an amount of at least %v for operator %v with suggested amounts %v", amount, operator.Name, amounts),
		}
		return 0, err
	}

	return amt.Pay, nil
}

func tryAutoFallback(err error) bool {
	if e, ok := err.(APIError); ok {
		switch e.ErrorCode {

		case "TRANSACTION_REFUSED_BY_OPERATOR", "INVALID_RECIPIENT_PHONE":
			return true

		default:
			return false
		}
	}
	return false
}

func (s *TopupsService) Topup(mobile string, amount float64) (*TopupResponse, error) {
	amt := amount

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

	if s.suggestedAmount {
		a, err := getSuggestedAmount(s.operator, amount, s.tolerance)
		if err != nil {
			return nil, err
		}
		amt = a
	}

	req := &TopupRequest{
		RecipientPhone: &RecipientPhone{s.operator.Country.IsoName, mobile},
		OperatorID: s.operator.OperatorID,
		Amount: amt,
	}

	resp := new(TopupResponse)
	_, err := s.Request("POST", "/topups", req, resp)

	// add retries??

	if err == nil || s.autoFallback == false || !tryAutoFallback(err) {
		return resp, err
	}

	// try with auto detect!
	return s.AutoDetect(s.operator.Country.IsoName).Topup(mobile, amount)
}
