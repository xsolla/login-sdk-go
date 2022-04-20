package infrastructure

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestHttpLoginApi_ValidateHS256Token(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), ValidateTokenAPIPATH)

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	api := HttpLoginApi{server.Client(), server.URL}
	err := api.ValidateHS256Token("some_token")

	assert.NoError(t, err)
}
