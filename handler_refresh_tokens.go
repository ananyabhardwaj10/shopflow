package main
import(
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/shopflow/internal/auth"
	"github.com/ananyabhardwaj10/shopflow/internal/database"
)

func (cfg *apiConfig) handlerRefreshTokens(w http.ResponseWriter, req *http.Request) {
	refToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token mismatch", err)
		return 
	}

	user, err := cfg.db.GetUserByRefreshToken(req.Context(), refToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get user using refresh token", err)
		return 
	}

	accessToken, err := auth.MakeJWT(user.ID, user.Role, cfg.jwtSecretKey, 15 * time.Minute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create access token", err)
		return 
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), refToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke refresh token for token rotation", err)
		return 
	}

	newRefToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create new refresh token", err)
		return 
	}

	refreshToken, err := cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: newRefToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to store new refresh token in database", err)
		return 
	}

	type response struct {
		UserID uuid.UUID `json:"user_id"`
		RefreshToken string `json:"refresh_token"`
		AccessToken string `json:"access_token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		UserID: user.ID,
		RefreshToken: refreshToken.Token,
		AccessToken: accessToken,
	})
}