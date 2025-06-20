package utils

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

type TokenType struct {
	id   string
	exp  string
	role string
}

var secretKey = []byte(configs.Env.TokenSecretWord)

func CreateRefreshToken(id string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   id,
			"exp":  time.Now().Add(time.Hour * 24 * 2).Unix(),
			"role": role,
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateAccesToken(id string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   id,
			"exp":  time.Now().Add(time.Hour * 1).Unix(),
			"role": role,
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return "", fmt.Errorf("token expired")
			}
		}
		//Verify if the token contains the id
		//!TODO: Verify the authorization of each user to do that action if the action is allowed from that types of the user depends on the "ROLE"
		// if id, ok := claims["id"].(string); ok {
		// 	return id, nil
		// }

		if role, ok := claims["role"].(string); ok {
			return role, nil
		}

		return "", fmt.Errorf("id not found in token")
	}

	return "", fmt.Errorf("invalid token")
}

func DecodeToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

//

func CreateAnId() (string, error) {
	Id := uuid.New().String()
	if Id == "" {
		return "", fmt.Errorf("Error while creating an id")
	}
	return Id, nil
}
