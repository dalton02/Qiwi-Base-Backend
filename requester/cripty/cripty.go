package jwtkit

import (
	"api_journal/core/util"
	httpkit "api_journal/requester/http"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func GenerateJwt[T any](data T) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["data"] = data
	claims["exp"] = util.DateFormatJwt()
	fmt.Println(claims["exp"])
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetJwtInfo(tokenString string, response http.ResponseWriter) map[string]interface{} {
	secretKey := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		httpkit.AppUnauthorized("Token inválido ou expirado", response)
	}

	data, _ := claims["data"].(map[string]interface{})
	return data
}
