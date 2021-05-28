package login_sdk_go

import (
	"log"
	"math/big"
)

func fromBase16(base16 string) *big.Int {
	i, ok := new(big.Int).SetString(base16, 16)
	if !ok {
		log.Fatal("bad number: " + base16)
	}
	return i
}
