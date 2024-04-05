package util

import (
	dto "github.com/greetinc/greet-auth-srv/dto/auth"

	"golang.org/x/crypto/bcrypt"
)

func GenerateFromPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func EncryptPassword(req *dto.AuthRegisterRequest) (err error) {
	hashedPassword, err := GenerateFromPassword(req.Password)

	if err != nil {
		return err
	}
	req.Password = string(hashedPassword)
	return nil
}

func EncryptPasswordAfterReset(req *dto.Reset) (err error) {
	hashedPassword, err := GenerateFromPassword(req.NewPassword)

	if err != nil {
		return err
	}
	req.NewPassword = string(hashedPassword)
	return nil
}
