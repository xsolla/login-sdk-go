package login_sdk_go

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LoginApiInterface interface {
	GetProjectKeysForLoginProject(projectID string) (RSAKeysResponse, error)
}

type LoginApi struct {
	baseUrl string
}

type RSAKeysResponse []struct {
	Alg      string `json:"alg"`
	Exponent string `json:"e"`
	Kid      string `json:"kid"`
	Kty      string `json:"kty"`
	Modulus  string `json:"n"`
	Use      string `json:"use"`
}

type LoginToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

func (l LoginApi) GetProjectKeysForLoginProject(projectID string) (RSAKeysResponse, error) {
	url := l.baseUrl + "/api/projects/" + projectID + "/keys"
	res, _ := http.Get(url)
	defer res.Body.Close()
	var keysResp RSAKeysResponse

	if err := json.NewDecoder(res.Body).Decode(&keysResp); err != nil {
		return nil, err
	}

	return keysResp, nil
}

func (l LoginApi) RefreshToken(refreshToken string) (LoginToken, error) {
	url := l.baseUrl + "/api/oauth2/token"
	payload := strings.NewReader("refresh_token=" + refreshToken + "&grant_type=refresh_token&client_secret=2&client_id=")

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	var loginToken LoginToken

	if err := json.NewDecoder(res.Body).Decode(&loginToken); err != nil {
		return LoginToken{}, err
	}

	return loginToken, nil
}
