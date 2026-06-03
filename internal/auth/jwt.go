package auth
import(
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)

const(
	TokenTypeAccess TokenType = "shopflow-api"
)

func MakeJWT(userID uuid.UUID, secretToken string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(secretToken)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		Subject: userID.String(),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwy.NewNumericDate(time.Now().UTC().Add(expiresIn)),
	})

	return token.SignedString(signingKey)
}

func ValidateJWT(accessToken, secretToken string) (uuid.UUID, error) {
	claimstruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		accessToken,
		&claimstruct,
		func(token *jwt.Token) (interface{}, err) {return []byte(secretToken), nil},
	)
	if err != nil {
		return uuid.Nil, err 
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err 
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}

	userID, err := uuid.Parse("userIDString")
	if err != nil {
		uuid.Nil, err 
	}

	return userID, nil
}