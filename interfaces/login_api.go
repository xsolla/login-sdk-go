package interfaces

import "gitlab.loc/sdk-login/login-sdk-go/model"

type LoginApi interface {
	GetProjectKeysForLoginProject(projectID string) ([]model.ProjectPublicKey, error)
	RefreshToken(token string, clientId int, clientSecret string) (*model.LoginToken, error)
}
