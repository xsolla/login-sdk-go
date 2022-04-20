package infrastructure

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.loc/sdk-login/login-sdk-go/interfaces"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

const (
	Timeout              = 3 * time.Second
	ValidateTokenAPIPATH = "/api/token/validate"
	ProjectsPath         = "/api/projects/"
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

func (api HttpLoginApi) makeRequest(ctx context.Context, method string, url string, body []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, 500, errors.New("failed create request:" + err.Error())
	}
	response, err := api.Client.Do(req)
	if err != nil {
		return nil, 500, errors.New("failed make request: " + err.Error())
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 500, errors.New("failed read response: " + err.Error())
	}

	return respBody, response.StatusCode, nil
}

func (api HttpLoginApi) GetProjectKeysForLoginProject(ctx context.Context, projectID string) ([]model.ProjectPublicKey, error) {
	response, _, err := api.makeRequest(ctx, "GET", fmt.Sprintf("%s%s%s%s", api.baseUrl, ProjectsPath, projectID, "/keys"), nil)
	if err != nil {
		return nil, err
	}

	var keysResp []model.ProjectPublicKey

	if err = json.Unmarshal(response, &keysResp); err != nil {
		return nil, errors.New("failed unmarshal data: " + err.Error())
	}

	return keysResp, nil
}

func (api HttpLoginApi) ValidateHS256Token(ctx context.Context, token string) error {

	values := map[string]string{"token": token}
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed marshal data %w", err)
	}

	_, statusCode, err := api.makeRequest(ctx, "GET", fmt.Sprintf("%s%s", api.baseUrl, ValidateTokenAPIPATH), data)
	if statusCode != 204 {
		return err
	}

	return nil
}
