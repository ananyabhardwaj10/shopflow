package auth 
import(
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"github.com/alexedwards/argon2id"
)

type ContextKey string 
const (
	ContextKeyUserID ContextKey = "userID"
	ContextKeyRole ContextKey = "role"
)

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")

	if token == "" {
		return "", fmt.Errorf("No authorization information found. Please try again.")
	}

	splitToken := strings.Split(token, " ")
	if len(splitToken) < 2 || splitToken[0] != "Bearer" {
		return "", fmt.Errorf("malformed authorization header")
	}

	return splitToken[1], nil
}

func HashPassword(password string) (string, error) {
	hashed_password, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err 
	}

	return hashed_password, nil 
}

func CheckHashedPassword(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	return match, err 
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err 
	}

	return hex.EncodeToString(token), nil 
}