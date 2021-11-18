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
	res, err := svc.OperatorsByCountry("IN")

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
	res, err := svc.SearchOperator("IN", "Reliance Jio India Bundles")

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
	_, err := svc.SearchOperator("IN", "fooo")

	assert.NotNil(t, err)
	assert.Equal(t, err.(ReloadlyError).ErrorCode, "OPERATOR_NOT_FOUND")
}
