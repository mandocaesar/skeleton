package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type Hash interface {
	IsPasswordHash(password, hash string) bool
	HashPassword(password string) (string, error)
}

type hash struct{}

func NewHash() Hash {
	return &hash{}
}

func (h *hash) IsPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *hash) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	return string(bytes), err
}
