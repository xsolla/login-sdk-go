package main

import (
	"gitlab.loc/sdk-login/login-sdk-go"
	"log"
)

func main() {
	config := login_sdk_go.Config{
		ShaSecretKey: "test",
	}
	loginSdk := login_sdk_go.New(config)
	_, err := loginSdk.Validate("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.5mhBHqs5_DTLdINd9p5m7ZJ6XD0Xc55kIaCRY5r6HRA")
	if err != nil {
		log.Fatal(err)
	}
}
