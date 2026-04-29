package jwt

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":	userID,
		"exp":		time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":	userID,
		"exp":		time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}

func ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string), nil
}

func ParseRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string), nil
}
