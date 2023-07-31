package utils

import (
	"github.com/machtwatch/catalyst-go-skeleton/utils/hash"
	"github.com/machtwatch/catalyst-go-skeleton/utils/jwt"
)

type Utils struct {
	JWT  jwt.JWT
	Hash hash.Hash
}

func NewUtils() Utils {
	return Utils{
		JWT:  jwt.NewJWT(),
		Hash: hash.NewHash(),
	}
}
