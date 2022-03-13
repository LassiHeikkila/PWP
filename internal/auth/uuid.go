package auth

import (
	"github.com/google/uuid"
)

func (a *authController) GenerateUUID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
