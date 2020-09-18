package reloadly

import "fmt"

type Country struct {
	IsoName string `json:"isoName,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Fx struct {
	Rate         float64 `json:"rate,omitempty"`
	CurrencyCode string  `json:"currencyCode,omitempty"`
}

type Operator struct {
	OperatorID                int64       `json:"operatorId,omitempty"`
	Name                      string      `json:"name,omitempty"`
	Bundle                    bool        `json:"bundle,omitempty"`
	Data                      bool        `json:"data,omitempty"`
	Pin                       bool        `json:"pin,omitempty"`
	SupportsLocalAmounts      bool        `json:"supportsLocalAmounts,omitempty"`
	DenominationType          string      `json:"denominationType,omitempty"`
	SenderCurrencyCode        string      `json:"senderCurrencyCode,omitempty"`
	SenderCurrencySymbol      string      `json:"senderCurrencySymbol,omitempty"`
	DestinationCurrencyCode   string      `json:"destinationCurrencyCode,omitempty"`
	DestinationCurrencySymbol string      `json:"destinationCurrencySymbol,omitempty"`
	Commission                float64     `json:"commission,omitempty"`
	InternationalDiscount     float64     `json:"internationalDiscount,omitempty"`
	LocalDiscount             float64     `json:"localDiscount,omitempty"`
	MostPopularAmount         float64     `json:"mostPopularAmount,omitempty"`
	MostPopularLocalAmount    float64     `json:"mostPopularLocalAmount,omitempty"`
	MinAmount                 float64     `json:"minAmount,omitempty"`
	MaxAmount                 float64     `json:"maxAmount,omitempty"`
	LocalMinAmount            float64     `json:"localMinAmount,omitempty"`
	LocalMaxAmount            float64     `json:"localMaxAmount,omitempty"`
	Country                   Country     `json:"country,omitempty"`
	Fx                        Fx          `json:"fx,omitempty"`
	LogoUrls                  []string    `json:"logoUrls,omitempty"`
	FixedAmounts              []float64   `json:"fixedAmounts,omitempty"`
	// FixedAmountsDescriptions  {
	// } `json:"fixedAmountsDescriptions,omitempty"`
	LocalFixedAmounts         []float64 `json:"localFixedAmounts,omitempty"`
	// LocalFixedAmountsDescriptions struct {
	// } `json:"localFixedAmountsDescriptions,omitempty"`
	SuggestedAmounts          []float64 `json:"suggestedAmounts,omitempty"`
	// SuggestedAmountsMap struct {
	// } `json:"suggestedAmountsMap,omitempty"`
	// Promotions []interface{} `json:"promotions,omitempty"`
}


func (s *Service) AutoDetect(mobile, country string) (*Operator, error) {
	path := fmt.Sprintf("/operators/auto-detect/phone/%v/countries/%v", mobile, country)

	resp := new(Operator)
	err := s.Request("GET", path, new(struct{}), resp)
	return resp, err
}
