package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jabuta/chirpy/internal/database"
)

func CreateAccessToken(uid int, key string) (string, error) {
	issuer := "chirpy-access"
	expires := time.Hour * 1
	return createJwtString(uid, key, issuer, expires)
}

func CreateRefreshToken(uid int, key string) (string, error) {
	issuer := "chirpy-refresh"
	expires := time.Hour * 24 * 60
	return createJwtString(uid, key, issuer, expires)
}

func createJwtString(uid int, key string, issuer string, expires time.Duration) (string, error) {
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

func CheckAccessToken(tokenString string, key string) (int, error) {
	uid, issuer, err := verifyToken(tokenString, key)
	if err != nil {
		return 0, err
	}
	if issuer != "chirpy-access" {
		return 0, errors.New("invalid token issuer")
	}
	return uid, nil
}

func CheckRefreshToken(tokenString string, key string, db *database.DB) (int, error) {
	uid, issuer, err := verifyToken(tokenString, key)
	if err != nil {
		return 0, err
	}
	if issuer != "chirpy-refresh" {
		return 0, errors.New("invalid token issuer")
	}
	err = db.TokenIsValid(tokenString)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func verifyToken(tokenString string, key string) (int, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, "", errors.New("no claims retreived")
	}
	if claims.Subject == "" {
		return 0, "", errors.New("invalid issuer, or subject")
	}
	uid, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, "", errors.New("invalid uid")
	}
	return uid, claims.Issuer, nil
}

func RevokeToken(token string, secret string, db *database.DB) error {
	if _, err := CheckRefreshToken(token, secret, db); err != nil {
		return err
	}
	err := db.AddRevocation(token)
	return err
}

func CheckHeaderAuth(r *http.Request, token string) (int, error) {
	if len(strings.Fields(r.Header.Get("Authorization"))) != 2 {

		return 0, errors.New("Token Invalid")
	}
	authToken := strings.Fields(r.Header.Get("Authorization"))[1]
	return CheckAccessToken(authToken, token)
}
