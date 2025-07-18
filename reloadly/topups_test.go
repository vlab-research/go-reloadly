package reloadly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getOperators() []Operator {
	dat, _ := ioutil.ReadFile("test/operators.json")
	var ops []Operator
	json.Unmarshal(dat, &ops)

	return ops
}

func TestPickAmountGetsEnough(t *testing.T) {

	for i, op := range getOperators() {
		amt, err := pickAmount(op.SuggestedAmountsMap, 100, 50)
		if err != nil {
			assert.Equal(t, i, 3)
			continue
		}

		assert.True(t, amt.Sent >= 100)
		assert.True(t, amt.Pay <= 2.0)
	}
}

func TestTopupReturnsErrorWithoutOperator(t *testing.T) {
	svc := &Service{}
	_, err := svc.Topups().Topup("+123", 100)
	assert.NotNil(t, err)
	assert.Equal(t, "INVALID_CALL", err.(ReloadlyError).ErrorCode)
}

func TestTopupReturnsErrorIfFindOperatorFails(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "/operators/countries/IN", r.URL.Path)

		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"timeStamp": "2020-09-18T08:26:27.577+0000", "message": "Not Found", "path": "/oauth/token", "errorCode": "OPERATOR_NOT_FOUND", "infoLink": null, "details": []}`)

	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.Topups().FindOperator("IN", "foo").Topup("+123", 100)

	assert.NotNil(t, err)
	assert.Equal(t, "OPERATOR_NOT_FOUND", err.(APIError).ErrorCode)
}

func TestTopupCallsWithOperatorIfFindOperatorSucceeds(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/operators.json")
	operators := string(dat)
	done := make(chan bool)
	ts, mux := TestServerMux()

	mux.HandleFunc("/operators/countries/IN", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, operators)
	})

	mux.HandleFunc("/topups", func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":200,"amount":100}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
		close(done)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.Topups().FindOperator("IN", "Airtel India").Topup("+123", 100)
	assert.Nil(t, err)

	<-done
}

func TestTopupReturnsErrorIfAutoDetectFails(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/operators/auto-detect/phone/+123/countries/IN", r.URL.Path)

		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"timeStamp": "2020-09-18T08:26:27.577+0000", "message": "Could not auto detect operator for given phone number", "path": "/oauth/token", "errorCode": "COULD_NOT_AUTO_DETECT_OPERATOR", "infoLink": null, "details": []}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.Topups().AutoDetect("IN").Topup("+123", 100)

	assert.NotNil(t, err)
	assert.Equal(t, "COULD_NOT_AUTO_DETECT_OPERATOR", err.(APIError).ErrorCode)
}

func TestTopupCallsWithOperatorIfAutoDetectSucceeds(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/airtel.json")
	airtel := string(dat)
	done := make(chan bool)
	ts, mux := TestServerMux()

	mux.HandleFunc("/operators/auto-detect/phone/+123/countries/IN", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, airtel)
	})

	mux.HandleFunc("/topups", func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":200,"amount":100}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
		close(done)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.Topups().AutoDetect("IN").Topup("+123", 100)
	assert.Nil(t, err)
}

func TestTopupBySuggestedAmountReturnsErrorOnEmptyAmounts(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {})
	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	op := Operator{Name: "Foodafone", DenominationType: "FIXED"}
	_, err := svc.Topups().SuggestedAmount(50).Operator(&op).Topup("+123", 100)

	assert.NotNil(t, err)
	assert.Equal(t, "IMPOSSIBLE_AMOUNT", err.(ReloadlyError).ErrorCode)
	assert.Contains(t, err.(ReloadlyError).Message, "Foodafone")
}

func TestTopupBySuggestedAmountReturnsErrorOnAmountOutOfRange(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":211,"amount":1.82}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	minAmount := float64(0)
	maxAmount := float64(50)
	op := Operator{
		Name:                 "Foodafone",
		DenominationType:     "RANGE",
		SupportsLocalAmounts: true,
		LocalMinAmount:       &minAmount,
		LocalMaxAmount:       &maxAmount,
	}
	_, err := svc.Topups().SuggestedAmount(50).Operator(&op).Topup("+123", 100)

	assert.NotNil(t, err)
	assert.Equal(t, "IMPOSSIBLE_AMOUNT", err.(ReloadlyError).ErrorCode)
	assert.Contains(t, err.(ReloadlyError).Message, "Foodafone")
	assert.Contains(t, err.(ReloadlyError).Message, "50")
}

func TestTopupBySuggestedAmountReturnsErrorOnAmountOutOfRangeNonLocalAmount(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":211,"amount":1.82}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	minAmount := float64(0)
	maxAmount := float64(50)
	op := Operator{
		Name:                 "Foodafone",
		DenominationType:     "RANGE",
		SupportsLocalAmounts: false,
		MinAmount:            &minAmount,
		MaxAmount:            &maxAmount,
	}
	_, err := svc.Topups().SuggestedAmount(50).Operator(&op).Topup("+123", 100)

	assert.NotNil(t, err)
	assert.Equal(t, "IMPOSSIBLE_AMOUNT", err.(ReloadlyError).ErrorCode)
	assert.Contains(t, err.(ReloadlyError).Message, "Foodafone")
	assert.Contains(t, err.(ReloadlyError).Message, "50")
}

func TestTopupBySuggestedAmountSendsRequestForGoodAmount(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":211,"amount":1.82}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	op := getOperators()[5]
	_, err := svc.Topups().SuggestedAmount(50).Operator(&op).Topup("+123", 100)

	assert.Nil(t, err)
}

func TestTopupBySuggestedAmountSendsAmountIfInRange(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":211,"amount":0.48}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	minAmount := float64(0)
	maxAmount := float64(50)
	op := Operator{
		OperatorID:           211,
		Name:                 "Foodafone",
		DenominationType:     "RANGE",
		Country:              Country{"IN", "India"},
		Fx:                   Fx{52.63, "INR"},
		SupportsLocalAmounts: true,
		LocalMinAmount:       &minAmount,
		LocalMaxAmount:       &maxAmount,
	}
	_, err := svc.Topups().SuggestedAmount(5).Operator(&op).Topup("+123", 25)

	assert.Nil(t, err)
}

func TestTopupBySuggestedAmountSendsAmountIfInRangeWithTolerance(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":211,"amount":0.58}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	minAmount := float64(30)
	maxAmount := float64(50)
	op := Operator{
		OperatorID:           211,
		Name:                 "Foodafone",
		DenominationType:     "RANGE",
		Country:              Country{"IN", "India"},
		Fx:                   Fx{52.63, "INR"},
		SupportsLocalAmounts: true,
		LocalMinAmount:       &minAmount,
		LocalMaxAmount:       &maxAmount,
	}
	_, err := svc.Topups().SuggestedAmount(5).Operator(&op).Topup("+123", 25)

	assert.Nil(t, err)
}

func TestTopupBySuggestedAmountSendsAmountIfInRangeWithNonLocalAmount(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":211,"amount":0.16}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	minAmount := float64(0)
	maxAmount := float64(10)
	op := Operator{
		OperatorID:           211,
		Name:                 "Foodafone",
		DenominationType:     "RANGE",
		Country:              Country{"IN", "India"},
		Fx:                   Fx{50, "INR"},
		SupportsLocalAmounts: false,
		MinAmount:            &minAmount,
		MaxAmount:            &maxAmount,
	}
	_, err := svc.Topups().SuggestedAmount(5).Operator(&op).Topup("+123", 8)

	assert.Nil(t, err)
}

func TestTopupBySuggestedAmountSendsAmountIfInRangeWithNonLocalAmountAndTolerance(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":211,"amount":2}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	minAmount := float64(2)
	maxAmount := float64(5)
	op := Operator{
		OperatorID:           211,
		Name:                 "Foodafone",
		DenominationType:     "RANGE",
		Country:              Country{"IN", "India"},
		Fx:                   Fx{50, "INR"},
		SupportsLocalAmounts: false,
		MinAmount:            &minAmount,
		MaxAmount:            &maxAmount,
	}
	_, err := svc.Topups().SuggestedAmount(75).Operator(&op).Topup("+123", 50)

	assert.Nil(t, err)
}

func TestTopupWithAutoFallbackReCallsWithNewOperatorId(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/airtel.json")
	airtel := string(dat)
	done := make(chan bool)
	ts, mux := TestServerMux()

	mux.HandleFunc("/operators/auto-detect/phone/+123/countries/IN", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, airtel)
	})

	count := 0
	mux.HandleFunc("/topups", func(w http.ResponseWriter, r *http.Request) {

		if count == 0 {
			expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":1,"amount":100}`

			data, _ := ioutil.ReadAll(r.Body)
			dat := strings.TrimSpace(string(data))
			assert.Equal(t, expected, dat)

			w.WriteHeader(404)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"errorCode": "INVALID_RECIPIENT_PHONE"}`)
		}

		if count == 1 {
			expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":200,"amount":100}`

			data, _ := ioutil.ReadAll(r.Body)
			dat := strings.TrimSpace(string(data))
			assert.Equal(t, expected, dat)

			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"hey": "yes"}`)

			close(done)
		}
		count++
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	op := Operator{Name: "Foodafone", OperatorID: 1, Country: Country{"IN", "India"}}
	_, err := svc.Topups().Operator(&op).AutoFallback().Topup("+123", 100)

	assert.Nil(t, err)
}

func TestTopupCustomIdentifierAddsIdentifier(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":211,"amount":1.82,"customIdentifier":"foobar"}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	op := getOperators()[5]
	_, err := svc.Topups().SuggestedAmount(50).Operator(&op).CustomIdentifier("foobar").Topup("+123", 100)

	assert.Nil(t, err)
}
