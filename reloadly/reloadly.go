package reloadly

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
)


type Service struct {
	Client *http.Client
	BaseUrl string
	AuthUrl string
	Token *Token
}

func New() *Service {
	return &Service{
		http.DefaultClient,
		"https://topups.reloadly.com",
		"https://auth.reloadly.com",
		nil,
	}
}

func NewSandbox() *Service {
	return &Service{
		http.DefaultClient,
		"https://topups-sandbox.reloadly.com",
		"https://auth.reloadly.com",
		nil,
	}
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
	if err == nil {
		err = apiError.AsError()
	}

	if err == nil && httpResponse.StatusCode >= 300 {
		return httpResponse, APIError{
			Message: fmt.Sprintf("Non-200 Status Code: %v", httpResponse.StatusCode),
			ErrorCode: fmt.Sprint(httpResponse.StatusCode),
		}
	}

	return httpResponse, err
}

func (s *Service) Request(method, path string, params interface{}, resp interface{}) (*http.Response, error) {
	sli := sling.New().Client(s.Client).Base(s.BaseUrl).Set("Accept", "application/com.reloadly.topups-v1+json")
	if s.Token != nil {
		sli = sli.Set("Authorization", fmt.Sprintf("%v %v", s.Token.TokenType, s.Token.AccessToken))
	}
	return s.request(sli, method, path, params, resp)
}
