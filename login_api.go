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

type RSAKeyResponse struct {
	Alg      string `json:"alg"`
	Exponent string `json:"e"`
	Kid      string `json:"kid"`
	Kty      string `json:"kty"`
	Modulus  string `json:"n"`
	Use      string `json:"use"`
}

type RSAKeysResponse []RSAKeyResponse

type LoginTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type LoginApi struct {
	baseUrl string
}

func (l LoginApi) GetProjectKeysForLoginProject(projectID string) (RSAKeysResponse, error) {
	endpoint := l.baseUrl + "/api/projects/" + projectID + "/keys"
	res, _ := http.Get(endpoint)
	defer res.Body.Close()
	var keysResp RSAKeysResponse

	if err := json.NewDecoder(res.Body).Decode(&keysResp); err != nil {
		return nil, err
	}

	return keysResp, nil
}

func (l LoginApi) RefreshToken(refreshToken string, clientId int, clientSecret string) (LoginTokenResponse, error) {
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
		return LoginTokenResponse{}, errors.New("http request error: " + res.Status)
	}

	var loginToken LoginTokenResponse

	if err := json.NewDecoder(res.Body).Decode(&loginToken); err != nil {
		return LoginTokenResponse{}, err
	}

	return loginToken, nil
}
