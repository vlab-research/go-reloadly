package reloadly

import (
	"encoding/json"
)

type TopupResponse struct {
	TransactionId int64 `json:"transactionId,omitempty"`
	OperatorTransactionId int64 `json:"operatorTransactionId,omitempty"`
	CustomIdentifier string `json:"customIdentifier,omitempty"`
	RecipientPhone string `json:"recipientPhone,omitempty"`
	SenderPhone string `json:"senderPhone,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	OperatorId int64 `json:"operatorId,omitempty"`
	OperatorName string `json:"operatorName,omitempty"`
	Discount float64 `json:"discount,omitempty"`
	DiscountCurrencyCode string `json:"discountCurrencyCode,omitempty"`
	RequestedAmount float64 `json:"requestedAmount,omitempty"`
	RequestedAmountCurrencyCode string `json:"requestedAmountCurrencyCode,omitempty"`
	DeliveredAmount float64 `json:"deliveredAmount,omitempty"`
	DeliveredAmountCurrencyCode string `json:"deliveredAmountCurrencyCode,omitempty"`
	TransactionDate string `json:"transactionDate,omitempty"`
	PinDetail *json.RawMessage `json:"pinDetail,omitempty"`
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
	RecipientPhone   RecipientPhone `json:"recipientPhone,omitempty"`
	SenderPhone      SenderPhone `json:"senderPhone,omitempty"`
	OperatorID       int64    `json:"operatorId,omitempty"`
	Amount           float64    `json:"amount,omitempty"`
	CustomIdentifier string `json:"customIdentifier,omitempty"`
}

func (s *Service) Topup(id, mobile string, operator *Operator, amount float64) (*TopupResponse, error) {
	req := &TopupRequest{
		RecipientPhone: RecipientPhone{operator.Country.IsoName, mobile},
		OperatorID: operator.OperatorID,
		Amount: amount,
		CustomIdentifier: id,
	}

	resp := new(TopupResponse)
	err := s.Request("POST", "/topups", req, resp)
	return resp, err
}

func (s *Service) AutoTopup(id, mobile, country string, amount float64) (*TopupResponse, error) {
	op, err := s.AutoDetect(mobile, country)
	if err != nil {
		return nil, err
	}

	return s.Topup(id, mobile, op, amount)
}
