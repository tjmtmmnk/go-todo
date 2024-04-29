package user

import "golang.org/x/crypto/bcrypt"

const (
	Cost = 10
)

func ToHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), Cost)
}
