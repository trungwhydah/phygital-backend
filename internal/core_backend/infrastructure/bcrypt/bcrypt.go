package bcrypt

import (
	"backend-service/internal/core_backend/common"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHash struct{}

func NewBcryptHash() *BcryptHash {
	return &BcryptHash{}
}

func (h *BcryptHash) HashPassword(data string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(data), common.HashCost)

	return string(bytes)
}

func (h *BcryptHash) VerifyPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
