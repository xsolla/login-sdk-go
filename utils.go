package login_sdk_go

import (
	"crypto/rsa"
	"math/big"
)

func fromBase16(base16 string) *big.Int {
	i, ok := new(big.Int).SetString(base16, 16)
	if !ok {
		panic("bad number: " + base16)
	}
	return i
}

func createPublicKey(k RSAKey) *rsa.PublicKey {
	return &rsa.PublicKey{
		N: fromBase16(k.Modulus),
		E: int(fromBase16(k.Exponent).Int64()),
	}
}
