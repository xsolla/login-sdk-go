package model

import (
	"crypto/rsa"
	"math/big"
)

type ProjectPublicKey struct {
	Alg      string `json:"alg"`
	Exponent string `json:"e"`
	Kid      string `json:"kid"`
	Kty      string `json:"kty"`
	Modulus  string `json:"n"`
	Use      string `json:"use"`
}

func (k ProjectPublicKey) CreateRSAPublicKey() *rsa.PublicKey {
	modulusBigInt, mOk := new(big.Int).SetString(k.Modulus, 16)
	exponentBigInt, eOk := new(big.Int).SetString(k.Exponent, 16)

	if mOk == false || eOk == false {
		return &rsa.PublicKey{}
	}

	return &rsa.PublicKey{
		N: modulusBigInt,
		E: int(exponentBigInt.Int64()),
	}
}
