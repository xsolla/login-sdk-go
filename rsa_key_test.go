package login_sdk_go

import (
	"context"
	"crypto/rsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.loc/sdk-login/login-sdk-go/model"
)

type testKeysStorage struct {
	result *[]model.ProjectPublicKey
	err    error
}

func (s testKeysStorage) GetProjectKeysForLoginProject(ctx context.Context, projectId string) ([]model.ProjectPublicKey, error) {
	return *s.result, s.err
}

func TestGetPublicKey(t *testing.T) {
	t.Run("should return err", func(t *testing.T) {
		pk := RSAPublicKeyGetter{projectId: "project-1", storage: testKeysStorage{result: &[]model.ProjectPublicKey{}, err: fmt.Errorf("an error occurred")}}
		_, err := pk.getPublicKey(context.Background(), "42")
		assert.Error(t, err)
	})

	t.Run("should return error if there is no keys for project", func(t *testing.T) {
		pk := RSAPublicKeyGetter{projectId: "project-1", storage: testKeysStorage{result: &[]model.ProjectPublicKey{}}}
		_, err := pk.getPublicKey(context.Background(), "42")
		assert.Error(t, err)
	})

	t.Run("should return error if there is no key with proper Kid", func(t *testing.T) {
		pk := RSAPublicKeyGetter{projectId: "project-1", storage: testKeysStorage{result: &[]model.ProjectPublicKey{{Kid: "12"}}}}
		_, err := pk.getPublicKey(context.Background(), "42")
		assert.Error(t, err)
	})

	t.Run("should return public key", func(t *testing.T) {
		modulus := "b24f209563937f253a57adf222822a89f1ee0a33d25826925ff8214868490fd3b2f6ff5dd9d422f412904fc1539c1e84f6dbce648cb28db6b3f2640a1c7cc3092066db5ba2ab38cab5618b58c6eef0994070bee055c521bdd43eea93c2146a07634276ca9b9ff7a1cdf20fe3e2f65e2719cb367b6c08072cfc0ec8caa6fbcfbfbafed23b70f050b827b0d2b3d1fd0dbd0b9adc7596021796336fae611d599827edca0c9d34f998db6e336fb522438bf7e24223eee09338034f5c9930d8edb27a84bf044d96c3709e601554af16b8d57c74333bd0690a40c5b615d84c348b9c3eb30dd0c94af1b0bc7d8ea6c8910fdaaeea1c7ae85f4b7f2b1b4576a50fb2a545"
		modulusBigInt, _ := new(big.Int).SetString(modulus, 16)
		expected := &rsa.PublicKey{
			E: 65537,
			N: modulusBigInt,
		}

		pk := RSAPublicKeyGetter{projectId: "project-1", storage: testKeysStorage{result: &[]model.ProjectPublicKey{{
			Alg:      "RS256",
			Exponent: "10001",
			Kid:      "42",
			Kty:      "RSA",
			Modulus:  modulus,
			Use:      "sig",
		}}}}

		key, err := pk.getPublicKey(context.Background(), "42")
		assert.NoError(t, err)
		assert.Equal(t, key, expected)
	})
}
