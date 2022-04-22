package interfaces

import (
	"context"

	"gitlab.loc/sdk-login/login-sdk-go/model"
)

type LoginApi interface {
	GetProjectKeysForLoginProject(ctx context.Context, projectID string) ([]model.ProjectPublicKey, error)
	ValidateHS256Token(ctx context.Context, token string) error
}
