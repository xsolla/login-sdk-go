package model

import (
	"crypto/rsa"
	"math/big"
)

const base = 16

type ProjectPublicKey struct {
	Alg      string `json:"alg"`
	Exponent string `json:"e"`
	Kid      string `json:"kid"`
	Kty      string `json:"kty"`
	Modulus  string `json:"n"`
	Use      string `json:"use"`
}

func (k ProjectPublicKey) CreateRSAPublicKey() *rsa.PublicKey {
	modulusBigInt, mOk := new(big.Int).SetString(k.Modulus, base)
	exponentBigInt, eOk := new(big.Int).SetString(k.Exponent, base)

	if !mOk || !eOk {
		return &rsa.PublicKey{}
	}

	return &rsa.PublicKey{
		N: modulusBigInt,
		E: int(exponentBigInt.Int64()),
	}
}
