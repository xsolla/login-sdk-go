module example.com/test-login-sdk

go 1.16

replace gitlab.loc/sdk-login/login-sdk-go => ../

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/onsi/gomega v1.14.0 // indirect
	gitlab.loc/sdk-login/login-sdk-go v0.0.0-00010101000000-000000000000
)
