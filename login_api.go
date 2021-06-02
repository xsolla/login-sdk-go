package login_sdk_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Timeout = 5 * time.Second
)

type RSAKey struct {
	Alg      string `json:"alg"`
	Exponent string `json:"e"`
	Kid      string `json:"kid"`
	Kty      string `json:"kty"`
	Modulus  string `json:"n"`
	Use      string `json:"use"`
}

type LoginToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type LoginApi interface {
	GetProjectKeysForLoginProject(projectID string) ([]RSAKey, error)
	RefreshToken(token string, clientId int, clientSecret string) (*LoginToken, error)
}

func NewHttpLoginApi(baseUrl string) LoginApi {
	return httpLoginApi{&http.Client{Timeout: Timeout}, baseUrl}
}

type httpLoginApi struct {
	Client  *http.Client
	baseUrl string
}

func (api httpLoginApi) GetProjectKeysForLoginProject(projectID string) ([]RSAKey, error) {
	endpoint := api.baseUrl + "/api/projects/" + projectID + "/keys"
	res, _ := api.Client.Get(endpoint)
	defer res.Body.Close()
	var keysResp []RSAKey

	if err := json.NewDecoder(res.Body).Decode(&keysResp); err != nil {
		return nil, err
	}

	return keysResp, nil
}

func (api httpLoginApi) RefreshToken(refreshToken string, clientId int, clientSecret string) (*LoginToken, error) {
	endpoint := api.baseUrl + "/api/oauth2/token"

	data := url.Values{}
	data.Set("client_id", fmt.Sprint(clientId))
	data.Add("client_secret", clientSecret)
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, errors.New("http request error: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := api.Client.Do(req)
	if err != nil {
		return nil, errors.New("http request error: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("http request error: " + res.Status)
	}

	var loginToken LoginToken

	if err := json.NewDecoder(res.Body).Decode(&loginToken); err != nil {
		return nil, err
	}

	return &loginToken, nil
}
