package login_sdk_go

import (
	"crypto/rsa"
	"math/big"
)

func fromBase16(base16 string) (*big.Int, bool) {
	return new(big.Int).SetString(base16, 16)
}

func createPublicKey(k RSAKey) *rsa.PublicKey {
	modulusBigInt, mOk := fromBase16(k.Modulus)
	exponentBigInt, eOk := fromBase16(k.Exponent)

	if mOk == false || eOk == false {
		return &rsa.PublicKey{}
	}

	return &rsa.PublicKey{
		N: modulusBigInt,
		E: int(exponentBigInt.Int64()),
	}
}
