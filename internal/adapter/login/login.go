package login

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gitlab.loc/sdk-login/login-sdk-go/model"
)

//nolint:gosec
const (
	validateTokenPath = "/api/token/validate"
	projectsPath      = "/api/projects"
)

var ErrWrongStatusCode = errors.New("wrong status code")

func (a *Adapter) makeRequest(ctx context.Context, method string, path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", a.baseUrl, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if a.extraHeaderName != "" {
		req.Header.Set(a.extraHeaderName, a.extraHeaderValue)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed make request: %w", err)
	}

	return response, nil
}

func (a *Adapter) GetProjectKeysForLoginProject(ctx context.Context, projectID string) ([]model.ProjectPublicKey, error) {
	response, err := a.makeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s/keys", projectsPath, projectID), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %w", err)
	}

	var keysResp []model.ProjectPublicKey

	if err = json.Unmarshal(respBody, &keysResp); err != nil {
		return nil, fmt.Errorf("failed unmarshal data: %w", err)
	}

	return keysResp, nil
}

func (a *Adapter) ValidateHS256Token(ctx context.Context, token string) error {
	values := map[string]string{"token": token}
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed marshal data %w", err)
	}

	response, err := a.makeRequest(ctx, http.MethodPost, validateTokenPath, data)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%w:%d", ErrWrongStatusCode, response.StatusCode)
	}

	return nil
}
