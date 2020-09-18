package reloadly

import "github.com/dghubble/sling"

type OAuthRequest struct {
	ID string `json:"client_id,omitempty"`
	Secret string `json:"client_secret,omitempty"`
	Audience string `json:"audience,omitempty"`
	GrantType string `json:"grant_type,omitempty"`
}

type Token struct {
	TokenType string `json:"token_type,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	Scope string `json:"scope,omitempty"`
	ExpiresIn int64 `json:"expires_in,omitempty"`
}

func (s *Service) GetOAuthToken(clientId, clientSecret string) (*Token, error){
	body := &OAuthRequest{clientId, clientSecret, s.AuthUrl, "client_credentials"}
	token := new(Token)

	sli := sling.New().Client(s.Client).Base(s.AuthUrl)

	err := s.request(sli, "POST", "oauth/token", body, token)
	return token, err
}

func (s *Service) Auth(clientId, clientSecret string) error {
	res, err := s.GetOAuthToken(clientId, clientSecret)
	if err != nil {
		return err
	}
	s.Token = res
	return nil
}
