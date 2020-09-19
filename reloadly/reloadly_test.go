package reloadly

import (
	"fmt"

	"net/http"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/dghubble/sling"
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
	_, err := svc.request(testSling, "GET", "/foo", new(struct{}), resp)

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
	_, err := svc.request(testSling, "GET", "/foo", new(struct{}), resp)

	assert.NotNil(t, err)
	e, ok := err.(APIError)

	assert.True(t, ok)
	assert.NotNil(t, e)
	assert.Equal(t, e.ErrorCode, "INVALID_CREDENTIALS")
}


func TestRequestGetReturnsErrorsOnHttpError(t *testing.T) {

	svc := &Service{}
	resp := new(struct{Bar string})
	_, err := svc.request(sling.New().Client(&http.Client{}).Base("http://foo"), "GET", "/foo", new(struct{}), resp)

	assert.NotNil(t, err)
	_, ok := err.(APIError)

	assert.False(t, ok)
	t.Log(err)
	// assert.NotNil(t, e)
	// assert.Equal(t, e.ErrorCode, "INVALID_CREDENTIALS")
}
