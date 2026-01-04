package auth

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/AaravMalani/cvwo-forum-backend/internal/env"
	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
)

var JwtTokenAuth *jwtauth.JWTAuth

func Setup() {
	JwtTokenAuth = jwtauth.New("HS256", env.JwtSecret, nil)
}

const (
	AuthGenerateToken = "utils.auth.GenerateToken"
	AuthAddJWT        = "utils.auth.AddJWT"
	JwtSignError      = "Unable to encode JWT in %s"
	JwtAddError       = "Unable to add JWT in %s"
	JWTMaxAge         = 7 * 86400 // 7 days
)

func GenerateToken(userId string) (string, error) {
	_, token, err := JwtTokenAuth.Encode(map[string]interface{}{"user": userId})
	return token, errors.Wrap(err, fmt.Sprintf(JwtSignError, AuthGenerateToken))
}

func AddJWT(w http.ResponseWriter, userId string) error {
	jwt, err := GenerateToken(userId)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(JwtAddError, AuthAddJWT))
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		MaxAge:   JWTMaxAge,
		HttpOnly: true,
	})
	return nil
}

// Todo: move out of auth
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())

		if err == nil && claims["user"] != nil {
			ctx := context.WithValue(r.Context(), "user", claims["user"].(string))
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateLogin(salt string, password string, hash string) bool {
	h := sha512.New()
	h.Write([]byte(password))
	h.Write([]byte(salt))

	password_attempt := h.Sum(nil)
	return hex.EncodeToString(password_attempt) == hash
}

func GeneratePasswordHash(salt string, password string) string {
	h := sha512.New()
	h.Write([]byte(password))
	h.Write([]byte(salt))

	return hex.EncodeToString(h.Sum(nil))
}
