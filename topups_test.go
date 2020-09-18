package reloadly

import (
	"testing"
	// "io/ioutil"
//	"net/http"
	// "fmt"

	// "github.com/stretchr/testify/assert"
)


func TestAutoTopupUsesDetectedOperator(t *testing.T) {

	// _, svc := NewTestService(NewService, func(w http.ResponseWriter, r *http.Request){

	// 	expected := fmt.Sprintf(`{"client_id":"id","client_secret":"secret","audience":"http://%v","grant_type":"client_credentials"}`, r.Host)

	// 	data, err := ioutil.ReadAll(r.Body)

	// 	dat := strings.TrimSpace(string(data))
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, "/oauth/token", r.URL.Path)
	// 	assert.Equal(t, expected, dat)

	// 	w.Header().Set("Content-Type", "application/json")
	// 	fmt.Fprintf(w, `{"token_type": "Bearer", "access_token": "foobarbaz", "expires_in": 86400, "scope": "foo bar baz"}`)
	// })

	// resp := new(TopupResponse)
	// svc.Post("/topups", &TopupRequest{}, resp)

	// token, err := svc.GetOAuthToken("id", "secret")
	// assert.Nil(t, err)
	// assert.Equal(t, "foobarbaz", token.AccessToken)
}
