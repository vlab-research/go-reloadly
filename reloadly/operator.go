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

type Fees struct {
	International           float64 `json:"international,omitempty"`
	Local                   float64 `json:"local,omitempty"`
	LocalPercentage         float64 `json:"localPercentage,omitempty"`
	InternationalPercentage float64 `json:"internationalPercentage,omitempty"`
}

type GeographicalRechargePlan struct {
	LocationCode                  string            `json:"locationCode,omitempty"`
	LocationName                  string            `json:"locationName,omitempty"`
	FixedAmounts                  []float64         `json:"fixedAmounts,omitempty"`
	LocalAmounts                  []float64         `json:"localAmounts,omitempty"`
	FixedAmountsPlanNames         map[string]string `json:"fixedAmountsPlanNames,omitempty"`
	FixedAmountsDescriptions      map[string]string `json:"fixedAmountsDescriptions,omitempty"`
	LocalFixedAmountsPlanNames    map[string]string `json:"localFixedAmountsPlanNames,omitempty"`
	LocalFixedAmountsDescriptions map[string]string `json:"localFixedAmountsDescriptions,omitempty"`
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
	ID                                int64                      `json:"id,omitempty"`
	OperatorID                        int64                      `json:"operatorId,omitempty"`
	Name                              string                     `json:"name,omitempty"`
	Bundle                            bool                       `json:"bundle,omitempty"`
	Data                              bool                       `json:"data,omitempty"`
	Pin                               bool                       `json:"pin,omitempty"`
	ComboProduct                      bool                       `json:"comboProduct,omitempty"`
	SupportsLocalAmounts              bool                       `json:"supportsLocalAmounts,omitempty"`
	SupportsGeographicalRechargePlans bool                       `json:"supportsGeographicalRechargePlans,omitempty"`
	DenominationType                  string                     `json:"denominationType,omitempty"`
	SenderCurrencyCode                string                     `json:"senderCurrencyCode,omitempty"`
	SenderCurrencySymbol              string                     `json:"senderCurrencySymbol,omitempty"`
	DestinationCurrencyCode           string                     `json:"destinationCurrencyCode,omitempty"`
	DestinationCurrencySymbol         string                     `json:"destinationCurrencySymbol,omitempty"`
	Commission                        float64                    `json:"commission,omitempty"`
	InternationalDiscount             float64                    `json:"internationalDiscount,omitempty"`
	LocalDiscount                     float64                    `json:"localDiscount,omitempty"`
	MostPopularAmount                 *float64                   `json:"mostPopularAmount,omitempty"`
	MostPopularLocalAmount            *float64                   `json:"mostPopularLocalAmount,omitempty"`
	MinAmount                         *float64                   `json:"minAmount,omitempty"`
	MaxAmount                         *float64                   `json:"maxAmount,omitempty"`
	LocalMinAmount                    *float64                   `json:"localMinAmount,omitempty"`
	LocalMaxAmount                    *float64                   `json:"localMaxAmount,omitempty"`
	Country                           Country                    `json:"country,omitempty"`
	Fx                                Fx                         `json:"fx,omitempty"`
	LogoUrls                          []string                   `json:"logoUrls,omitempty"`
	FixedAmounts                      []float64                  `json:"fixedAmounts,omitempty"`
	FixedAmountsDescriptions          map[string]string          `json:"fixedAmountsDescriptions,omitempty"`
	LocalFixedAmounts                 []float64                  `json:"localFixedAmounts,omitempty"`
	LocalFixedAmountsDescriptions     map[string]string          `json:"localFixedAmountsDescriptions,omitempty"`
	SuggestedAmounts                  []float64                  `json:"suggestedAmounts,omitempty"`
	SuggestedAmountsMap               SuggestedAmountsMap        `json:"suggestedAmountsMap,omitempty"`
	Fees                              Fees                       `json:"fees,omitempty"`
	GeographicalRechargePlans         []GeographicalRechargePlan `json:"geographicalRechargePlans,omitempty"`
	// Promotions []interface{} `json:"promotions,omitempty"`
}

// GetFixedAmounts returns the appropriate fixed amounts based on whether geographical recharge plans are supported
// If geographical plans are supported, it returns amounts from the first available plan
// Otherwise, it returns the legacy fixed amounts
func (o *Operator) GetFixedAmounts() []float64 {
	if o.SupportsGeographicalRechargePlans && len(o.GeographicalRechargePlans) > 0 {
		// Use the first geographical plan as default when we don't know user's location
		return o.GeographicalRechargePlans[0].FixedAmounts
	}
	return o.FixedAmounts
}

// GetLocalFixedAmounts returns the appropriate local fixed amounts based on whether geographical recharge plans are supported
// If geographical plans are supported, it returns local amounts from the first available plan
// Otherwise, it returns the legacy local fixed amounts
func (o *Operator) GetLocalFixedAmounts() []float64 {
	if o.SupportsGeographicalRechargePlans && len(o.GeographicalRechargePlans) > 0 {
		// Use the first geographical plan as default when we don't know user's location
		return o.GeographicalRechargePlans[0].LocalAmounts
	}
	return o.LocalFixedAmounts
}

// GetFixedAmountsDescriptions returns the appropriate fixed amounts descriptions
// If geographical plans are supported, it returns descriptions from the first available plan
// Otherwise, it returns the legacy descriptions
func (o *Operator) GetFixedAmountsDescriptions() map[string]string {
	if o.SupportsGeographicalRechargePlans && len(o.GeographicalRechargePlans) > 0 {
		return o.GeographicalRechargePlans[0].FixedAmountsDescriptions
	}
	return o.FixedAmountsDescriptions
}

// GetLocalFixedAmountsDescriptions returns the appropriate local fixed amounts descriptions
// If geographical plans are supported, it returns descriptions from the first available plan
// Otherwise, it returns the legacy descriptions
func (o *Operator) GetLocalFixedAmountsDescriptions() map[string]string {
	if o.SupportsGeographicalRechargePlans && len(o.GeographicalRechargePlans) > 0 {
		return o.GeographicalRechargePlans[0].LocalFixedAmountsDescriptions
	}
	return o.LocalFixedAmountsDescriptions
}

// GetGeographicalPlanByLocationCode returns a specific geographical plan by location code
// Returns nil if not found or if geographical plans are not supported
func (o *Operator) GetGeographicalPlanByLocationCode(locationCode string) *GeographicalRechargePlan {
	if !o.SupportsGeographicalRechargePlans {
		return nil
	}

	for _, plan := range o.GeographicalRechargePlans {
		if plan.LocationCode == locationCode {
			return &plan
		}
	}
	return nil
}

// GetGeographicalPlanByLocationName returns a specific geographical plan by location name
// Returns nil if not found or if geographical plans are not supported
func (o *Operator) GetGeographicalPlanByLocationName(locationName string) *GeographicalRechargePlan {
	if !o.SupportsGeographicalRechargePlans {
		return nil
	}

	for _, plan := range o.GeographicalRechargePlans {
		if plan.LocationName == locationName {
			return &plan
		}
	}
	return nil
}

// GetDefaultGeographicalPlan returns the first available geographical plan
// Returns nil if geographical plans are not supported or empty
func (o *Operator) GetDefaultGeographicalPlan() *GeographicalRechargePlan {
	if !o.SupportsGeographicalRechargePlans || len(o.GeographicalRechargePlans) == 0 {
		return nil
	}
	return &o.GeographicalRechargePlans[0]
}

type OperatorsParams struct {
	SuggestedAmounts    bool `url:"suggestedAmounts,omitempty"`
	SuggestedAmountsMap bool `url:"suggestedAmountsMap,omitempty"`
	IncludeBundles      bool `url:"includeBundles,omitempty"`
	IncludeData         bool `url:"includeData,omitempty"`
	IncludePin          bool `url:"includePin,omitempty"`
}

func (s *TopupsService) OperatorsAutoDetect(mobile, country string) (*Operator, error) {
	path := fmt.Sprintf("/operators/auto-detect/phone/%v/countries/%v", mobile, country)
	params := &OperatorsParams{SuggestedAmountsMap: true, SuggestedAmounts: true}
	resp := new(Operator)
	_, err := s.Request("GET", path, params, resp)
	return resp, err
}

func (s *TopupsService) OperatorsByCountry(country string) ([]Operator, error) {
	path := fmt.Sprintf("/operators/countries/%v", country)
	params := &OperatorsParams{SuggestedAmountsMap: true, SuggestedAmounts: true}
	resp := new([]Operator)
	_, err := s.Request("GET", path, params, resp)
	return *resp, err
}

func (s *TopupsService) SearchOperator(country, name string) (*Operator, error) {
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

func (s *TopupsService) GetOperatorByID(operatorID int64) (*Operator, error) {
	path := fmt.Sprintf("/operators/%v", operatorID)
	params := &OperatorsParams{SuggestedAmountsMap: true, SuggestedAmounts: true}
	resp := new(Operator)
	_, err := s.Request("GET", path, params, resp)
	return resp, err
}
