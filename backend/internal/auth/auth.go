package auth
import (
	"time"
	"net/http"
	"github.com/go-chi/jwtauth/v5"
)

var TokenAuth *jwtauth.JWTAuth

func InitJWT(secret string) {
	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func GenerateToken(userId int) (string, error) {
	_, tokenString, err := TokenAuth.Encode(map[string]interface{}{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return tokenString, err
}

// Boas práticas para proteger rota
func VerifyToken(next http.Handler) http.Handler {
	return jwtauth.Verifier(TokenAuth)(jwtauth.Authenticator(next))
}
