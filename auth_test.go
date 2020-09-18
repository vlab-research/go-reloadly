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

		expected := fmt.Sprintf(`{"client_id":"id","client_secret":"secret","audience":"http://%v","grant_type":"client_credentials"}`, r.Host)

		data, err := ioutil.ReadAll(r.Body)

		dat := strings.TrimSpace(string(data))
		assert.Nil(t, err)
		assert.Equal(t, "/oauth/token", r.URL.Path)
		assert.Equal(t, expected, dat)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"token_type": "Bearer", "access_token": "foobarbaz", "expires_in": 86400, "scope": "foo bar baz"}`)
	})


	svc := &Service{AuthUrl: ts.URL, Client: &http.Client{}}

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

	svc := &Service{AuthUrl: ts.URL, Client: &http.Client{}}

	err := svc.Auth("id", "secret")
	assert.NotNil(t, err)
	e, ok := err.(*ErrorResponse)

	assert.True(t, ok)
	assert.NotNil(t, e)
	assert.Equal(t, e.ErrorCode, "INVALID_CREDENTIALS")
}
