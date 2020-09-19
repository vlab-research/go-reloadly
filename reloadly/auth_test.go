package reloadly

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestGetAuthTokenMakesAudience(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request){

		expected := `{"client_id":"id","client_secret":"secret","audience":"reloadly.com","grant_type":"client_credentials"}`

		data, err := ioutil.ReadAll(r.Body)

		dat := strings.TrimSpace(string(data))
		assert.Nil(t, err)
		assert.Equal(t, "/oauth/token", r.URL.Path)
		assert.Equal(t, expected, dat)

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"token_type": "Bearer", "access_token": "foobarbaz", "expires_in": 86400, "scope": "foo bar baz"}`)
	})


	svc := &Service{AuthUrl: ts.URL, BaseUrl: "reloadly.com", Client: &http.Client{}}

	err := svc.Auth("id", "secret")
	assert.Nil(t, err)
	assert.Equal(t, "foobarbaz", svc.Token.AccessToken)
}


func TestGetAuthTokenReturnsErrors(t *testing.T) {

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request){

		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"timeStamp": "2020-09-18T08:26:27.577+0000", "message": "Access Denied", "path": "/oauth/token", "errorCode": "INVALID_CREDENTIALS", "infoLink": null, "details": []}`)
	})

	svc := &Service{AuthUrl: ts.URL, BaseUrl: "reloadly.com", Client: &http.Client{}}

	err := svc.Auth("id", "secret")
	assert.NotNil(t, err)
	e, ok := err.(APIError)

	assert.True(t, ok)
	assert.NotNil(t, e)
	assert.Equal(t, e.ErrorCode, "INVALID_CREDENTIALS")
	assert.Equal(t, e.StatusCode, 401)
}
