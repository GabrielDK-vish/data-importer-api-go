package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
    if s := os.Getenv("JWT_SECRET"); s != "" {
        return s
    }
    return "sua-chave-secreta-super-segura-aqui"
}

// Claims representa as claims do JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken gera um token JWT para o usuário
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token válido por 24 horas
	
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken valida um token JWT
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, errors.New("token inválido")
	}
	
	return claims, nil
}

func ValidateCredentials(username, password string) bool {
	validUsers := map[string]string{
		"admin": "admin123",
		"user":  "user123",
		"demo":  "demo123",
	}
	
	expectedPassword, exists := validUsers[username]
	return exists && expectedPassword == password
}