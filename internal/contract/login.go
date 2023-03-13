package contract

import (
	"context"

	"github.com/xsolla/login-sdk-go/model"
)

type LoginAPI interface {
	GetProjectKeysForLoginProject(ctx context.Context, projectID string) ([]model.ProjectPublicKey, error)
	ValidateHS256Token(ctx context.Context, token string) error
}
