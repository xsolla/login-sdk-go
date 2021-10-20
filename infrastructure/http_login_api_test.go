package infrastructure

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.loc/sdk-login/login-sdk-go/model"
)

func TestHttpLoginApi_GetProjectKeysForLoginProject(t *testing.T) {
	testProjectId := "test"
	expected := &[]model.ProjectPublicKey{{Kid: "12"}, {Kid: "21"}}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/projects/"+testProjectId+"/keys")

		js, err := json.Marshal(expected)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		_, err = rw.Write(js)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	api := HttpLoginApi{server.Client(), server.URL}
	body, err := api.GetProjectKeysForLoginProject(testProjectId)

	assert.NoError(t, err)
	assert.Equal(t, expected, &body)
}

func TestHttpLoginApi_RefreshToken(t *testing.T) {
	expected := &model.LoginToken{
		AccessToken:  "test_access",
		RefreshToken: "test_refresh",
		ExpiresIn:    3600,
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/oauth2/token")

		js, err := json.Marshal(expected)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		_, err = rw.Write(js)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	api := HttpLoginApi{server.Client(), server.URL}
	body, err := api.RefreshToken("refresh_token", 1, "secret")

	assert.NoError(t, err)
	assert.Equal(t, expected, body)
}

func TestHttpLoginApi_ValidateHS256Token(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/token/validate")

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	api := HttpLoginApi{server.Client(), server.URL}
	err := api.ValidateHS256Token("some_token")

	assert.NoError(t, err)
}
