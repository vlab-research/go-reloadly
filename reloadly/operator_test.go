package reloadly

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOperatorsByCountry(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/operators.json")
	operators := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/operators/countries/IN", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.topups-v1+json")
		fmt.Fprintf(w, operators)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.Topups().OperatorsByCountry("IN")

	assert.Nil(t, err)
	assert.Equal(t, "Airtel India", res[0].Name)
	assert.Equal(t, "BSNL India", res[1].Name)
}

func TestSearchOperatorWhenFoundReturnsOperator(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/operators.json")
	operators := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/operators/countries/IN", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.topups-v1+json")
		fmt.Fprintf(w, operators)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.Topups().SearchOperator("IN", "Reliance Jio India Bundles")

	assert.Nil(t, err)
	assert.Equal(t, "Reliance Jio India Bundles", res.Name)
	assert.Equal(t, int64(186), res.OperatorID)
}

func TestSearchOperatorWhenNotFoundReturnsError(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/operators.json")
	operators := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/operators/countries/IN", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.topups-v1+json")
		fmt.Fprintf(w, operators)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.Topups().SearchOperator("IN", "fooo")

	assert.NotNil(t, err)
	assert.Equal(t, err.(ReloadlyError).ErrorCode, "OPERATOR_NOT_FOUND")
}

func TestGetOperatorByID(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/operator_single.json")
	operator := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/operators/186", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.topups-v1+json")
		fmt.Fprintf(w, operator)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.Topups().GetOperatorByID(186)

	assert.Nil(t, err)
	assert.Equal(t, "Reliance Jio India Bundles", res.Name)
	assert.Equal(t, int64(186), res.OperatorID)
}

func TestGetFixedAmountsWithGeographicalPlans(t *testing.T) {
	// Test operator with geographical recharge plans
	op := &Operator{
		SupportsGeographicalRechargePlans: true,
		GeographicalRechargePlans: []GeographicalRechargePlan{
			{
				LocationCode: "HP",
				LocationName: "Himachal Pradesh",
				FixedAmounts: []float64{0.17, 2.0, 8.33, 16.67, 83.33},
			},
		},
		FixedAmounts: []float64{1.0, 2.0, 3.0}, // Legacy amounts (should be ignored)
	}

	amounts := op.GetFixedAmounts()
	expected := []float64{0.17, 2.0, 8.33, 16.67, 83.33}

	if len(amounts) != len(expected) {
		t.Errorf("Expected %d amounts, got %d", len(expected), len(amounts))
	}

	for i, amount := range amounts {
		if amount != expected[i] {
			t.Errorf("Expected amount %f at index %d, got %f", expected[i], i, amount)
		}
	}
}

func TestGetFixedAmountsWithoutGeographicalPlans(t *testing.T) {
	// Test operator without geographical recharge plans
	op := &Operator{
		SupportsGeographicalRechargePlans: false,
		FixedAmounts:                      []float64{1.0, 2.0, 3.0}, // Legacy amounts (should be used)
	}

	amounts := op.GetFixedAmounts()
	expected := []float64{1.0, 2.0, 3.0}

	if len(amounts) != len(expected) {
		t.Errorf("Expected %d amounts, got %d", len(expected), len(amounts))
	}

	for i, amount := range amounts {
		if amount != expected[i] {
			t.Errorf("Expected amount %f at index %d, got %f", expected[i], i, amount)
		}
	}
}

func TestGetLocalFixedAmountsWithGeographicalPlans(t *testing.T) {
	// Test operator with geographical recharge plans
	op := &Operator{
		SupportsGeographicalRechargePlans: true,
		GeographicalRechargePlans: []GeographicalRechargePlan{
			{
				LocationCode: "HP",
				LocationName: "Himachal Pradesh",
				LocalAmounts: []float64{10.0, 120.0, 500.0, 1000.0, 5000.0},
			},
		},
		LocalFixedAmounts: []float64{1.0, 2.0, 3.0}, // Legacy amounts (should be ignored)
	}

	amounts := op.GetLocalFixedAmounts()
	expected := []float64{10.0, 120.0, 500.0, 1000.0, 5000.0}

	if len(amounts) != len(expected) {
		t.Errorf("Expected %d amounts, got %d", len(expected), len(amounts))
	}

	for i, amount := range amounts {
		if amount != expected[i] {
			t.Errorf("Expected amount %f at index %d, got %f", expected[i], i, amount)
		}
	}
}

func TestGetGeographicalPlanByLocationCode(t *testing.T) {
	op := &Operator{
		SupportsGeographicalRechargePlans: true,
		GeographicalRechargePlans: []GeographicalRechargePlan{
			{
				LocationCode: "HP",
				LocationName: "Himachal Pradesh",
				FixedAmounts: []float64{0.17, 2.0, 8.33},
			},
			{
				LocationCode: "DL",
				LocationName: "Delhi",
				FixedAmounts: []float64{1.0, 2.0, 3.0},
			},
		},
	}

	// Test finding existing plan
	plan := op.GetGeographicalPlanByLocationCode("HP")
	if plan == nil {
		t.Error("Expected to find plan with location code HP")
	}
	if plan.LocationName != "Himachal Pradesh" {
		t.Errorf("Expected location name 'Himachal Pradesh', got '%s'", plan.LocationName)
	}

	// Test finding non-existing plan
	plan = op.GetGeographicalPlanByLocationCode("XX")
	if plan != nil {
		t.Error("Expected nil for non-existing location code")
	}

	// Test operator without geographical plans
	op2 := &Operator{
		SupportsGeographicalRechargePlans: false,
	}
	plan = op2.GetGeographicalPlanByLocationCode("HP")
	if plan != nil {
		t.Error("Expected nil for operator without geographical plans")
	}
}

func TestGetDefaultGeographicalPlan(t *testing.T) {
	op := &Operator{
		SupportsGeographicalRechargePlans: true,
		GeographicalRechargePlans: []GeographicalRechargePlan{
			{
				LocationCode: "HP",
				LocationName: "Himachal Pradesh",
				FixedAmounts: []float64{0.17, 2.0, 8.33},
			},
		},
	}

	plan := op.GetDefaultGeographicalPlan()
	if plan == nil {
		t.Error("Expected to get default plan")
	}
	if plan.LocationCode != "HP" {
		t.Errorf("Expected location code 'HP', got '%s'", plan.LocationCode)
	}

	// Test operator without geographical plans
	op2 := &Operator{
		SupportsGeographicalRechargePlans: false,
	}
	plan = op2.GetDefaultGeographicalPlan()
	if plan != nil {
		t.Error("Expected nil for operator without geographical plans")
	}

	// Test operator with empty geographical plans
	op3 := &Operator{
		SupportsGeographicalRechargePlans: true,
		GeographicalRechargePlans:         []GeographicalRechargePlan{},
	}
	plan = op3.GetDefaultGeographicalPlan()
	if plan != nil {
		t.Error("Expected nil for operator with empty geographical plans")
	}
}
