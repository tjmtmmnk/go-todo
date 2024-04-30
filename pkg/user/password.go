package user

import "golang.org/x/crypto/bcrypt"

const (
	Cost = 10
)

type RawPassword string
type HashedPassword string

func ToHash(password RawPassword) (HashedPassword, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return HashedPassword(""), err
	}
	return HashedPassword(hashedPassword), nil
}
