package main
import(
	"net/http"
	"log"
	"context"

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

func roleMiddleware(role string) (func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		ctx := req.Context()
		roleVal, ok := ctx.Value(auth.ContextKeyRole).(string)
		if !ok {
			log.Printf("Missing Role information")
			w.WriteHeader(http.StatusUnauthorized)
			return 
		}

		if roleVal != role {
			log.Printf("Role mismatch")
			w.WriteHeader(http.StatusForbidden)
			return 
		}

		next.ServeHTTP(w, req)
	})
}}