package main

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(envOr("JWT_SECRET", "yuanfu-ricefield-dev-secret-change-me"))

type employeeClaims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	RoleCode string `json:"roleCode"`
	ParkId   int64  `json:"parkId"`
	ShopId   int64  `json:"shopId"`
	jwt.RegisteredClaims
}

func signToken(c *employeeClaims) (string, error) {
	c.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(jwtSecret)
}

func parseToken(tokenStr string) (*employeeClaims, error) {
	claims := &employeeClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名方法错误")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token 无效")
	}
	return claims, nil
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if h == "" {
		return r.URL.Query().Get("token")
	}
	return strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
}

func hashPassword(pw string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b)
}

func checkPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
