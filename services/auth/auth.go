package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wael-boudissaa/marquinoBackend/configs"
	"github.com/wael-boudissaa/marquinoBackend/types"
	"golang.org/x/crypto/bcrypt"
)

func ComparePasswords(password []byte, hashedPassword []byte) bool {
	// fmt.Println(password)
	// fmt.Println(hashedPassword)
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}

func HashedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

var secretKey = []byte(configs.Env.TokenSecretWord)

func CreateRefreshToken(user types.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  user.Id,
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
