package auth
import(
	"time"
	"fmt"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string 

const(
	TokenTypeAccess TokenType = "shopflow-api"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func MakeJWT(userID uuid.UUID, role, secretToken string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(secretToken)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		Subject: userID.String(),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
	}})

	return token.SignedString(signingKey)
}

func ValidateJWT(accessToken, secretToken string) (uuid.UUID, string, error) {
	claimstruct := CustomClaims{}
	token, err := jwt.ParseWithClaims(
		accessToken,
		&claimstruct,
		func(token *jwt.Token) (interface{}, error) {return []byte(secretToken), nil},
	)
	if err != nil {
		return uuid.Nil, "", err 
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, "", err 
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, "", err
	}

	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, "", fmt.Errorf("invalid issuer")
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, "", err 
	}

	role := claimstruct.Role
	
	return userID, role, nil
}