package utils

import (
	"crypto/rand"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateReferalCode(length int) (string, error) {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		sb.WriteRune(charset[randNum.Int64()])
	}
	return sb.String(), nil
}

func HashPassword(password string) (string, error) {
	// todo add password validation
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
