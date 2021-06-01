package login_sdk_go

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpLoginApi_GetProjectKeysForLoginProject(t *testing.T) {
	testProjectId := "test"
	expected := []RSAKey{{Kid: "12"}, {Kid: "21"}}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/projects/"+testProjectId+"/keys")

		js, err := json.Marshal(expected)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)
	}))
	defer server.Close()

	api := httpLoginApi{server.Client(), server.URL}
	body, err := api.GetProjectKeysForLoginProject(testProjectId)

	assert.NoError(t, err)
	assert.Equal(t, expected, body)
}

func TestHttpLoginApi_RefreshToken(t *testing.T) {
	expected := LoginToken{
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
		rw.Write(js)
	}))
	defer server.Close()

	api := httpLoginApi{server.Client(), server.URL}
	body, err := api.RefreshToken("refresh_token", 1, "secret")

	assert.NoError(t, err)
	assert.Equal(t, expected, body)
}
