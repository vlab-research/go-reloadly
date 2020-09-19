package reloadly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	// "strings"
	"testing"

	// "io/ioutil"
	//	"net/http"
	// "fmt"

	"github.com/stretchr/testify/assert"
)


func getOperators() []Operator {
	dat, _ := ioutil.ReadFile("test/operators.json")
	var ops []Operator
    json.Unmarshal(dat, &ops)

	return ops
}

func TestPickAmountGetsEnough(t *testing.T) {
	ops := getOperators()

	for i, op := range ops {
		amt, err := pickAmount(op.SuggestedAmountsMap, 100, 50)

		if err != nil {
			assert.Equal(t, i , 3)
			continue
		}

		assert.True(t, amt.Sent >= 100)
		assert.True(t, amt.Pay <= 2.0)
	}
}


func TestTopupBySuggestedAmountReturnsErrorOnEmptyAmounts(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request){})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.TopupBySuggestedAmount("foo", "+123", &Operator{Name: "Foodafone"}, 100, 50)

	assert.NotNil(t, err)
	assert.Equal(t, "IMPOSSIBLE_AMOUNT", err.(ReloadlyError).ErrorCode)
	assert.Contains(t, err.(ReloadlyError).Message, "Foodafone")
}

func TestTopupBySuggestedAmountSendsRequestForGoodAmount(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request){

		expected := `{"recipientPhone":{"countryCode":"IN","number":"+123"},"operatorId":211,"amount":1.82,"customIdentifier":"foo"}`

		data, _ := ioutil.ReadAll(r.Body)
		dat := strings.TrimSpace(string(data))
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"hey": "yeah"}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}

	op := getOperators()[5]
	_, err := svc.TopupBySuggestedAmount("foo", "+123", &op, 100, 50)

	assert.Nil(t, err)
}
