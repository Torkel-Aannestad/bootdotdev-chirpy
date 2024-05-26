package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

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

func GenerateJTW(email string, id int, expiresInSeconds int, JWTSecret string) (signedToken string, err error) {

	var jwtExpires int64 = 86400
	if expiresInSeconds < 86400 && expiresInSeconds != 0 {
		jwtExpires = int64(expiresInSeconds)
	}
	var jwtExpiresTime int64 = time.Now().UTC().Unix() + jwtExpires

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

	fmt.Printf("claims.subject: %v\n", claims.Subject)

	return claims.Subject, nil
}
