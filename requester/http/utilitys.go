package httpkit

import (
	"api_journal/core/util"
	dtoRequest "api_journal/requester/dto"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func GetBearerToken(auth string) string {
	result := strings.Replace(auth, "Bearer", "", -1)
	result = strings.TrimSpace(result)
	return result
}

// Returns a struct with a count of the params and a map[string]string to get the param
func GetUrlParams(request *http.Request) (dtoRequest.Params, error) {
	paramsInterface := request.Context().Value("params")
	params, test := paramsInterface.(dtoRequest.Params)
	if !test {
		return params, errors.New("erro ao obter parametros")
	}
	return params, nil
}

func GetDataToken(request *http.Request) any {
	tokenInterface := request.Context().Value("token")
	return tokenInterface
}

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
		AppUnauthorized("Token inválido ou expirado", response)
	}

	data, _ := claims["data"].(map[string]interface{})
	return data
}
