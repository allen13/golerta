package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func CreateToken(signingKey string, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func CreateExpirationFreeAgentToken(name , signingKey string)(string){
	claims := jwt.MapClaims{
		"iss": "golerta-token-tool",
		"iat": time.Now().Unix(),
		"jti": name,
	}

	token, _ := CreateToken(signingKey, claims)
	return token
}
