package mainimport(
	"net/http"
	"log"

	"github.com/ananyabhardwaj10/shopflow/internal/auth"
)

func (cfg *apiConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		accessToken, err := auth.GetBearerToken(req.Header)
		if err != nil {
			log.Printf("Error extracting access token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, role, err := auth.ValidateJWT(accessToken, cfg.jwtSecretKey)
		if err != nil {
			log.Printf("Error validating the access token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return 
		}

		ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, userID)
		ctx = context.WithValue(ctx, auth.ContextKeyRole, role)
		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}