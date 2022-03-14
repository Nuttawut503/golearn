package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	_accessSecret  = "Acc3ss"
	_refreshSecret = "l23fr3sh"
)

type Token struct {
	Token  string
	UUID   string
	Expire int64
}

func GenerateToken(username string, isRefresh bool) (Token, error) {
	randomUUID := uuid.New().String()
	claims := jwt.MapClaims{
		"iss":      "MY_APP",
		"iat":      time.Now().Unix(),
		"username": username,
	}
	var (
		ss     string
		expire int64
		err    error
		secret string
	)
	if isRefresh {
		expire = time.Now().Add(time.Minute * 10).Unix()
		claims["refresh_uuid"] = randomUUID
		secret = _refreshSecret
	} else {
		expire = time.Now().Add(time.Minute).Unix()
		claims["access_uuid"] = randomUUID
		claims["authorized"] = true
		secret = _accessSecret
	}
	claims["exp"] = expire
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err = token.SignedString([]byte(secret))
	return Token{
		Token:  ss,
		UUID:   randomUUID,
		Expire: expire,
	}, err
}

func GeneratePairToken(username string) ([]Token, error) {
	tokens := make([]Token, 2)
	var err error
	if tokens[0], err = GenerateToken(username, false); err != nil {
		return nil, err
	}
	if tokens[1], err = GenerateToken(username, true); err != nil {
		return nil, err
	}
	return tokens, nil
}

func ExtractToken(token string, isRefresh bool) (jwt.MapClaims, error) {
	secret := _accessSecret
	if isRefresh {
		secret = _refreshSecret
	}
	v, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected method (%v)", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := v.Claims.(jwt.MapClaims)
	if !ok || !v.Valid {
		return nil, errors.New("Can't parse")
	}
	return claims, nil
}
