package login_sdk_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	RefreshToken(token string, clientId int, clientSecret string) (LoginToken, error)
}

func NewHttpLoginApi(baseUrl string) LoginApi {
	return httpLoginApi{baseUrl}
}

type httpLoginApi struct {
	baseUrl string
}

func (l httpLoginApi) GetProjectKeysForLoginProject(projectID string) ([]RSAKey, error) {
	endpoint := l.baseUrl + "/api/projects/" + projectID + "/keys"
	res, _ := http.Get(endpoint)
	defer res.Body.Close()
	var keysResp []RSAKey

	if err := json.NewDecoder(res.Body).Decode(&keysResp); err != nil {
		return nil, err
	}

	return keysResp, nil
}

func (l httpLoginApi) RefreshToken(refreshToken string, clientId int, clientSecret string) (LoginToken, error) {
	client := &http.Client{}
	endpoint := l.baseUrl + "/api/oauth2/token"

	data := url.Values{}
	data.Set("client_id", fmt.Sprint(clientId))
	data.Add("client_secret", clientSecret)
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return LoginToken{}, errors.New("http request error: " + res.Status)
	}

	var loginToken LoginToken

	if err := json.NewDecoder(res.Body).Decode(&loginToken); err != nil {
		return LoginToken{}, err
	}

	return loginToken, nil
}
