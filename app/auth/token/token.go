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

func CreateExpiringToken(username, signingKey string, expirationLength time.Duration, backend string) string {
	expirationTimestamp := time.Now().Add(expirationLength).Unix()
	claims := jwt.MapClaims{
		"jti":  username,
		"iss":  backend,
		"exp":  expirationTimestamp,
		"name": username,
		//Everyone who logs in is an admin by default for now. Could check ldap groups for this.
		"role": "admin",
	}

	token, _ := CreateToken(signingKey, claims)
	return token
}

func CreateExpirationFreeAgentToken(name, signingKey string) string {
	claims := jwt.MapClaims{
		"iss": "golerta-token-tool",
		"iat": time.Now().Unix(),
		"jti": name,
	}

	token, _ := CreateToken(signingKey, claims)
	return token
}
