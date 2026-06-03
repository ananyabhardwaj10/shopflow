package auth 
import(
	"net/http"
	"fmt"
	"strings"
)

func GetBearerToken(headers http.Header) (string, err) {
	token := headers.Get("Authorization")

	if token = "" {
		"", fmt.Errorf("No authorization information found. Please try again.")
	}

	splitToken := strings.Split(" ")
	if len(splitToken) < 2 || splitToken[0] != "Bearer" {
		return "", fmt.Errorf("malformed authorization header")
	}

	return splitToken[1], nil
}