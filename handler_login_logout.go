package main
import(
	"net/http"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/shopflow/internal/auth"
	"github.com/ananyabhardwaj10/shopflow/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incomplete email or password. Please try again.", err)
		return 
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email) 
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect email or password.", err)
		return 
	}

	match, err := auth.CheckHashedPassword(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to check password", err)
		return 
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password. Please try again.", err)
		return 
	}

	ref_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create refresh token", err)
		return 
	}

	accessToken, err := auth.MakeJWT(user.ID, user.Role, cfg.jwtSecretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create access token", err)
		return 
	}

	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: ref_token, 
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to store refresh token in database.", err)
		return 
	}

	type response struct {
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
		Email string `json:"email"`
		UserID uuid.UUID `json:"user_id"`
		RefreshToken string `json:"refresh_token"`
		AccessToken string `json:"access_token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		UserID: user.ID,
		RefreshToken: ref_token,
		AccessToken: accessToken,
	})
	
}

func (cfg *apiConfig) handlerLogout(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incomplete Authorization Information", err)
		return 
	}

	_, err = cfg.db.GetRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh Token mismatch", err)
		return 
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke refresh token", err)
		return 
	}

	w.WriteHeader(http.StatusOK)
}