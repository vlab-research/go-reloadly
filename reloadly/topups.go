package reloadly

import (
	"errors"
	"fmt"
	"sort"
)


type TopupResponse struct {
	TransactionID               int64       `json:"transactionId,omitempty"`
	OperatorTransactionID       string      `json:"operatorTransactionId,omitempty"`
	CustomIdentifier            string      `json:"customIdentifier,omitempty"`
	RecipientPhone              string      `json:"recipientPhone,omitempty"`
	RecipientEmail              string      `json:"recipientEmail,omitempty"`
	SenderPhone                 string      `json:"senderPhone,omitempty"`
	CountryCode                 string      `json:"countryCode,omitempty"`
	OperatorID                  int64       `json:"operatorId,omitempty"`
	OperatorName                string      `json:"operatorName,omitempty"`
	Discount                    float64     `json:"discount,omitempty"`
	DiscountCurrencyCode        string      `json:"discountCurrencyCode,omitempty"`
	RequestedAmount             float64     `json:"requestedAmount,omitempty"`
	RequestedAmountCurrencyCode string      `json:"requestedAmountCurrencyCode,omitempty"`
	DeliveredAmount             float64     `json:"deliveredAmount,omitempty"`
	DeliveredAmountCurrencyCode string      `json:"deliveredAmountCurrencyCode,omitempty"`
	TransactionDate             string      `json:"transactionDate,omitempty"`
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


func pickAmount(amounts []SuggestedAmount, min float64, tolerance float64) (*SuggestedAmount, error) {
	sort.Slice(amounts, func(i, j int) bool { return amounts[i].Sent < amounts[j].Sent})

	for _, a := range amounts {
		if a.Sent >= min && a.Sent <= min + tolerance {
			return &a, nil
		}
	}

	return nil, errors.New("no amount found")
}


func (s *Service) Topup(id, mobile string, operator *Operator, amount float64) (*TopupResponse, error) {
	req := &TopupRequest{
		RecipientPhone: &RecipientPhone{operator.Country.IsoName, mobile},
		OperatorID: operator.OperatorID,
		Amount: amount,
		CustomIdentifier: id,
	}

	resp := new(TopupResponse)
	_, err := s.Request("POST", "/topups", req, resp)
	return resp, err
}

func (s *Service) TopupBySuggestedAmount(
	id, mobile string,
	operator *Operator,
	amount float64,
	tolerance float64) (*TopupResponse, error) {

	amounts := operator.SuggestedAmountsMap
	amt, err := pickAmount(amounts, amount, tolerance)

	if err != nil {
		err = ReloadlyError{
			ErrorCode: "IMPOSSIBLE_AMOUNT",
			Message: fmt.Sprintf("Could not manage to find an amount of at least %v for operator %v with suggested amounts %v", amount, operator.Name, amounts),
		}
		return nil, err
	}

	return s.Topup(id, mobile, operator, amt.Pay)
}
