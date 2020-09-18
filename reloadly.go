package reloadly

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/sling"
)

type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	format := "2006-01-02T15:04:05.000-0700"
	s := strings.Trim(string(b), "\"")
	parsed, err := time.Parse(format, s)
	if err != nil {
		return err
	}
	*t = Timestamp(parsed)
	return nil
}

type ErrorResponse struct {
	ErrorCode string `json:"errorCode,omitempty"`
	Message string `json:"message,omitempty"`
	TimeStamp *Timestamp `json:"timeStamp,omitempty"`
	InfoLink string `json:"infoLink,omitempty"`
	Path string `json:"path,omitempty"`
	Details []map[string]string `json:"details,omitempty"`
}

func (e *ErrorResponse) AsError() error {
	if e.ErrorCode == "" {
		return nil
	}
	return e
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

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

func (s *Service) request(sli *sling.Sling, method, path string, params interface{}, resp interface{}) error {
	switch strings.ToUpper(method) {
	case "GET":
		sli = sli.Get(path).QueryStruct(params)
	case "POST":
		sli = sli.Post(path).BodyJSON(params)
	}

	reloadlyError := new(ErrorResponse)
	_, err := sli.Receive(resp, reloadlyError)
	if err == nil {
		err = reloadlyError.AsError()
	}
	return err
}

func (s *Service) Request(method, path string, params interface{}, resp interface{}) error {
	sli := sling.New().Client(s.Client).Base(s.BaseUrl).Set("Accept", "application/com.reloadly.topups-v1+json")
	if s.Token != nil {
		sli = sli.Set("Authorization", fmt.Sprintf("%v %v", s.Token.TokenType, s.Token.AccessToken))
	}
	return s.request(sli, method, path, params, resp)
}
