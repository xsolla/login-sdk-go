package interfaces

import "gitlab.loc/sdk-login/login-sdk-go/model"

type LoginApi interface {
	GetProjectKeysForLoginProject(projectID string) ([]model.ProjectPublicKey, error)
	ValidateHS256Token(token string) error
}
