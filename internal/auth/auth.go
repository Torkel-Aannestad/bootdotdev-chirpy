package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func GeneratePasswordHash(password string) (string, error) {
	hashedPassord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassord), nil
}

func CompareHashAndPassword(passwordHash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

func GenerateJTW(email string, id int, JWTSecret string) (signedToken string, err error) {
	const JWT_EXPIRES int64 = 3600
	jwtExpiresTime := time.Now().UTC().Unix() + JWT_EXPIRES

	claims := jwt.StandardClaims{
		Issuer:    "chirpy",
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiresAt: jwtExpiresTime,
		Subject:   fmt.Sprint(id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

func VerifyTokenAndGetUser(signedToken string, JWTSecret string) (string, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return "", errors.New("jwt.ParseWithClaims")
	}

	if !token.Valid {
		fmt.Printf("!token.Valid: %v\n", token.Valid)
		return "", errors.New("!token.Valid")
	}

	return claims.Subject, nil
}

func GenerateRefreshToken() (refreshToken string) {
	dat := make([]byte, 256)
	rand.Read(dat)
	refreshToken = hex.EncodeToString(dat)

	return refreshToken
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
