package reloadly

import (
	"fmt"
	"net/http"

	"testing"

	"github.com/dghubble/sling"
	"github.com/stretchr/testify/assert"
)

func TestRequestGetNoQueryParams(t *testing.T) {

	_, testSling := TestServer(func(w http.ResponseWriter, r *http.Request) {
		vals := r.URL.RawQuery
		assert.Equal(t, "", vals)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"bar": "baz"}`)
	})

	svc := &Service{}

	resp := new(struct{ Bar string })
	_, err := svc.request(testSling, "GET", "/foo", new(struct{}), resp)

	assert.Nil(t, err)
	assert.Equal(t, "baz", resp.Bar)
}

func TestRequestGetReturnsErrors(t *testing.T) {

	_, testSling := TestServer(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"timeStamp": "2020-09-18T08:26:27.577+0000", "message": "Access Denied", "path": "/oauth/token", "errorCode": "INVALID_CREDENTIALS", "infoLink": null, "details": []}`)
	})

	svc := &Service{}

	resp := new(struct{ Bar string })
	_, err := svc.request(testSling, "GET", "/foo", new(struct{}), resp)

	assert.NotNil(t, err)
	e, ok := err.(APIError)

	assert.True(t, ok)
	assert.NotNil(t, e)
	assert.Equal(t, e.ErrorCode, "INVALID_CREDENTIALS")
}

func TestRequestGetReturnsErrorsOnHttpError(t *testing.T) {

	svc := &Service{}
	resp := new(struct{ Bar string })
	_, err := svc.request(sling.New().Client(&http.Client{}).Base("http://foo"), "GET", "/foo", new(struct{}), resp)

	assert.NotNil(t, err)
	_, ok := err.(APIError)
	assert.False(t, ok)
}

func TestRequestDoesReAuthOnErrorCodeTOKEN_EXPIRED(t *testing.T) {
	done := make(chan bool)
	ts, mux := TestServerMux()

	authCount := 0
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		if authCount == 0 {
			fmt.Fprintf(w, `{"token_type": "Bearer", "access_token": "foobar", "expires_in": 86400, "scope": "foo bar baz"}`)
		}
		if authCount == 1 {
			fmt.Fprintf(w, `{"token_type": "Bearer", "access_token": "foobarbaz", "expires_in": 86400, "scope": "foo bar baz"}`)
		}
		authCount++
	})

	count := 0
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		if count == 0 {
			assert.Equal(t, "Bearer foobar", r.Header.Get("Authorization"))
			w.WriteHeader(401)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"errorCode":"TOKEN_EXPIRED"}`)
		}

		if count == 1 {
			assert.Equal(t, "Bearer foobarbaz", r.Header.Get("Authorization"))
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"Bar": "qux"}`)
			close(done)
		}
		count++
	})
	svc := &Service{BaseUrl: ts.URL, AuthUrl: ts.URL, Client: &http.Client{}}
	resp := new(struct{ Bar string })

	err := svc.Auth("id", "secret")
	assert.Nil(t, err)

	_, err = svc.Request("GET", "/foo", new(struct{}), resp)
	assert.Nil(t, err)

	assert.Equal(t, 2, authCount)
	assert.Equal(t, 2, count)
	<-done
}
