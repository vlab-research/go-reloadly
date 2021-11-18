package reloadly

import (
	"net/http"
	"net/http/httptest"

	"github.com/dghubble/sling"
)

type TestTransport func(req *http.Request) (*http.Response, error)

func (r TestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

func TestServer(handler func(http.ResponseWriter, *http.Request)) (*httptest.Server, *sling.Sling) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	sli := sling.New().Client(&http.Client{}).Base(ts.URL)
	return ts, sli
}

func TestServerMux() (*httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	return ts, mux
}
