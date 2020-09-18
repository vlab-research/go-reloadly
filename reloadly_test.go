package reloadly

import (
	"fmt"
	// "io/ioutil"
	"net/http"
	// "strings"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestRequestGetNoQueryParams(t *testing.T) {

	_, testSling := TestServer(func(w http.ResponseWriter, r *http.Request){
		vals := r.URL.RawQuery
		assert.Equal(t, "", vals)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"bar": "baz"}`)
	})

	svc := &Service{}

	resp := new(struct{Bar string})
	err := svc.request(testSling, "GET", "/foo", new(struct{}), resp)

	assert.Nil(t, err)
	assert.Equal(t, "baz", resp.Bar)
}


func TestRequestGetReturnsErrors(t *testing.T) {

	_, testSling := TestServer(func(w http.ResponseWriter, r *http.Request){

		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"timeStamp": "2020-09-18T08:26:27.577+0000", "message": "Access Denied", "path": "/oauth/token", "errorCode": "INVALID_CREDENTIALS", "infoLink": null, "details": []}`)
	})

	svc := &Service{}

	resp := new(struct{Bar string})
	err := svc.request(testSling, "GET", "/foo", new(struct{}), resp)

	assert.NotNil(t, err)
	e, ok := err.(*ErrorResponse)

	assert.True(t, ok)
	assert.NotNil(t, e)
	assert.Equal(t, e.ErrorCode, "INVALID_CREDENTIALS")
}
