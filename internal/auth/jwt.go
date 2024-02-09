package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwtString(uid int, key string, issuer string) (string, error) {
	var expires time.Duration
	switch issuer {
	case "chirpy-access":
		expires = time.Hour * 1
	case "chirpy-refresh":
		expires = time.Hour * 24 * 60
	default:
		return "", errors.New("invalid issuer")
	}

	signingMethod := jwt.SigningMethodHS256
	registeredClaims := jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expires)),
		Subject:   fmt.Sprint(uid),
	}
	token := jwt.NewWithClaims(signingMethod, registeredClaims)
	signedString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func CheckToken(tokenString string, key string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, errors.New("no claims retreived")
	}
	if claims.Issuer != "chirpy" || claims.Subject == "" {
		return 0, errors.New("invalid issuer, or subject")
	}
	uid, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, errors.New("invalid uid")
	}
	return uid, nil
}
