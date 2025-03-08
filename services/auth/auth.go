package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wael-boudissaa/zencitiBackend/configs"
	"golang.org/x/crypto/bcrypt"
)

func ComparePasswords(password []byte, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}

func HashedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

var secretKey = []byte(configs.Env.TokenSecretWord)

func CreateRefreshToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": time.Now().Add(time.Hour * 24 * 2).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateAnId() (string, error) {
	var Id string
	Id = uuid.New().String()
	if Id == "" {
		return "", fmt.Errorf("Error while creating an id")
	}
	return Id, nil
}
