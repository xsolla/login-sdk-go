package infrastructure

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.loc/sdk-login/login-sdk-go/interfaces"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

const (
	Timeout = 5 * time.Second
)

type HttpLoginApi struct {
	Client  *http.Client
	baseUrl string
}

func NewHttpLoginApi(baseUrl string, ignoreSslErrors bool) interfaces.LoginApi {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ignoreSslErrors},
	}
	return HttpLoginApi{&http.Client{Timeout: Timeout, Transport: tr}, baseUrl}
}

func (api HttpLoginApi) GetProjectKeysForLoginProject(projectID string) ([]model.ProjectPublicKey, error) {
	endpoint := api.baseUrl + "/api/projects/" + projectID + "/keys"
	res, err := api.Client.Get(endpoint)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var keysResp []model.ProjectPublicKey

	if err := json.NewDecoder(res.Body).Decode(&keysResp); err != nil {
		return nil, err
	}

	return keysResp, nil
}

func (api HttpLoginApi) RefreshToken(refreshToken string, clientId int, clientSecret string) (*model.LoginToken, error) {
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

	var loginToken model.LoginToken
	if err := json.NewDecoder(res.Body).Decode(&loginToken); err != nil {
		return nil, err
	}

	return &loginToken, nil
}
