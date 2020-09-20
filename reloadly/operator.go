package reloadly

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Country struct {
	IsoName string `json:"isoName,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Fx struct {
	Rate         float64 `json:"rate,omitempty"`
	CurrencyCode string  `json:"currencyCode,omitempty"`
}

type SuggestedAmount struct {
	Pay  float64
	Sent float64
}

type SuggestedAmountsMap []SuggestedAmount

func (s *SuggestedAmountsMap) UnmarshalJSON(b []byte) error {
	var m map[string]float64
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	res := []SuggestedAmount{}

	for payStr, sent := range m {
		pay, err := strconv.ParseFloat(payStr, 64)
		if err != nil {
			return err
		}
		res = append(res, SuggestedAmount{pay, sent})
	}

	*s = SuggestedAmountsMap(res)
	return nil
}

type Operator struct {
	OperatorID                int64     `json:"operatorId,omitempty"`
	Name                      string    `json:"name,omitempty"`
	Bundle                    bool      `json:"bundle,omitempty"`
	Data                      bool      `json:"data,omitempty"`
	Pin                       bool      `json:"pin,omitempty"`
	SupportsLocalAmounts      bool      `json:"supportsLocalAmounts,omitempty"`
	DenominationType          string    `json:"denominationType,omitempty"`
	SenderCurrencyCode        string    `json:"senderCurrencyCode,omitempty"`
	SenderCurrencySymbol      string    `json:"senderCurrencySymbol,omitempty"`
	DestinationCurrencyCode   string    `json:"destinationCurrencyCode,omitempty"`
	DestinationCurrencySymbol string    `json:"destinationCurrencySymbol,omitempty"`
	Commission                float64   `json:"commission,omitempty"`
	InternationalDiscount     float64   `json:"internationalDiscount,omitempty"`
	LocalDiscount             float64   `json:"localDiscount,omitempty"`
	MostPopularAmount         float64   `json:"mostPopularAmount,omitempty"`
	MostPopularLocalAmount    float64   `json:"mostPopularLocalAmount,omitempty"`
	MinAmount                 float64   `json:"minAmount,omitempty"`
	MaxAmount                 float64   `json:"maxAmount,omitempty"`
	LocalMinAmount            float64   `json:"localMinAmount,omitempty"`
	LocalMaxAmount            float64   `json:"localMaxAmount,omitempty"`
	Country                   Country   `json:"country,omitempty"`
	Fx                        Fx        `json:"fx,omitempty"`
	LogoUrls                  []string  `json:"logoUrls,omitempty"`
	FixedAmounts              []float64 `json:"fixedAmounts,omitempty"`
	// FixedAmountsDescriptions  {
	// } `json:"fixedAmountsDescriptions,omitempty"`
	LocalFixedAmounts []float64 `json:"localFixedAmounts,omitempty"`
	// LocalFixedAmountsDescriptions struct {
	// } `json:"localFixedAmountsDescriptions,omitempty"`
	SuggestedAmounts    []float64           `json:"suggestedAmounts,omitempty"`
	SuggestedAmountsMap SuggestedAmountsMap `json:"suggestedAmountsMap,omitempty"`
	// Promotions []interface{} `json:"promotions,omitempty"`
}

type OperatorsParams struct {
	SuggestedAmounts    bool `url:"suggestedAmounts,omitempty"`
	SuggestedAmountsMap bool `url:"suggestedAmountsMap,omitempty"`
	IncludeBundles      bool `url:"includeBundles,omitempty"`
	IncludeData         bool `url:"includeData,omitempty"`
	IncludePin          bool `url:"includePin,omitempty"`
}

func (s *Service) OperatorsAutoDetect(mobile, country string) (*Operator, error) {
	path := fmt.Sprintf("/operators/auto-detect/phone/%v/countries/%v", mobile, country)

	params := &OperatorsParams{SuggestedAmountsMap: true}
	resp := new(Operator)
	_, err := s.Request("GET", path, params, resp)
	return resp, err
}

func (s *Service) OperatorsByCountry(country string) ([]Operator, error) {
	path := fmt.Sprintf("/operators/countries/%v", country)
	resp := new([]Operator)

	params := &OperatorsParams{SuggestedAmountsMap: true}

	_, err := s.Request("GET", path, params, resp)
	return *resp, err
}

func (s *Service) SearchOperator(country, name string) (*Operator, error) {
	ops, err := s.OperatorsByCountry(country)
	if err != nil {
		return nil, err
	}

	for _, op := range ops {
		if op.Name == name {
			return &op, nil
		}
	}

	err = ReloadlyError{
		"OPERATOR_NOT_FOUND",
		fmt.Sprintf("Could not find operator with name: %v in country: %v", name, country),
	}
	return nil, err
}
