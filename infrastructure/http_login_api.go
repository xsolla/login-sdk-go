package infrastructure

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.loc/sdk-login/login-sdk-go/interfaces"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

const (
	Timeout              = 3 * time.Second
	ValidateTokenAPIPATH = "/api/token/validate"
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

func (api HttpLoginApi) ValidateHS256Token(token string) error {
	endpoint := api.baseUrl + ValidateTokenAPIPATH

	values := map[string]string{"token": token}
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed marshal data %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.New("http request error: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := api.Client.Do(req)
	if err != nil {
		return errors.New("http request error: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return errors.New("http request error: " + res.Status)
	}

	return nil
}
