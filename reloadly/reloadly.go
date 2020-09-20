package reloadly

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
)

type Service struct {
	Client  *http.Client
	BaseUrl string
	AuthUrl string
	Token   *Token
}

func New() *Service {
	return &Service{
		http.DefaultClient,
		"https://topups.reloadly.com",
		"https://auth.reloadly.com",
		nil,
	}
}

func (s *Service) Sandbox() {
	s.BaseUrl = "https://topups-sandbox.reloadly.com"
}

func (s *Service) request(sli *sling.Sling, method, path string, params interface{}, resp interface{}) (*http.Response, error) {
	switch strings.ToUpper(method) {
	case "GET":
		sli = sli.Get(path).QueryStruct(params)
	case "POST":
		sli = sli.Post(path).BodyJSON(params)
	}

	apiError := APIError{}
	httpResponse, err := sli.Receive(resp, &apiError)
	if err != nil {
		return nil, err
	}

	status := httpResponse.StatusCode

	if !apiError.Empty() {
		apiError.StatusCode = status
		return httpResponse, apiError
	}

	// Reloadly will send an erro response without
	// a body, sometimes, so we just create our
	// own "APIError" from the status.
	if status < 200 || status > 299 {
		return httpResponse, APIError{
			Message:    httpResponse.Status,
			ErrorCode:  fmt.Sprint(status),
			StatusCode: status,
		}
	}
	return httpResponse, nil
}

func (s *Service) Request(method, path string, params interface{}, resp interface{}) (*http.Response, error) {

	sli := sling.New().Client(s.Client).Base(s.BaseUrl).Set("Accept", "application/com.reloadly.topups-v1+json")

	if s.Token != nil {
		auth := fmt.Sprintf("%v %v", s.Token.TokenType, s.Token.AccessToken)
		sli = sli.Set("Authorization", auth)
	}
	return s.request(sli, method, path, params, resp)
}
