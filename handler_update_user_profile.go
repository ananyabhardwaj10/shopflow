package main
import(
	"net/http"
	"encoding/json"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/shopflow/internal/auth"
	"github.com/ananyabhardwaj10/shopflow/internal/database"
)

func (cfg *apiConfig) handlerUpdateUserProfile(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		FirstName *string `json:"first_name"`
		LastName *string `json:"last_name"`
		ContactNumber *string `json:"contact_number"`
		Address *string `json:"address"`
		Email *string `json:"email"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Improper data. Please try again.", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get user id using context", err)
		return 
	}

	first_name := sql.NullString{}
	last_name := sql.NullString{}
	contact_number := sql.NullString{}
	address := sql.NullString{}
	email := sql.NullString{}

	if params.FirstName != nil {
		first_name.String = *params.FirstName
		first_name.Valid = true
	}

	if params.LastName != nil {
		last_name.String = *params.LastName
		last_name.Valid = true
	}

	if params.ContactNumber != nil {
		contact_number.String = *params.ContactNumber
		contact_number.Valid = true
	}

	if params.Address != nil {
		address.String = *params.Address
		address.Valid = true
	}

	if params.Email != nil {
		email.String = *params.Email
		email.Valid = true
	}

	user, err := cfg.db.UpdateUserDetails(req.Context(), database.UpdateUserDetailsParams{
		ID: userID,
		FirstName: first_name,
		LastName: last_name,
		ContactNumber: contact_number,
		Address: address,
		Email: email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update user details", err)
		return 
	}

	type response struct {
		UserID uuid.UUID `json:"user_id"`
		Email string `json:"email"`
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
	}

	respondWithJSON(w, http.StatusOK, response{
		UserID: userID,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
	})
}

func (cfg *apiConfig) handlerChangePassword(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		CurrentPassword string `json:"current_password"`
		NewPassword string `json:"new_password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incomplete Data. Please try again.", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get user id from context", err)
		return 
	}

	user, err := cfg.db.GetUserFromUserID(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get user from database using user id", err)
		return 
	}

	match, err := auth.CheckHashedPassword(params.CurrentPassword, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to compare passwords.", err)
		return 
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Mismatched credentials. Please try again.", err)
		return 
	}

	hashed_new_password, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash new password", err)
		return 
	}

	err = cfg.db.ChangeUserPassword(req.Context(), database.ChangeUserPasswordParams{
		HashedPassword: hashed_new_password,
		ID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to change password. Please try again.", err)
		return 
	}

	w.WriteHeader(http.StatusOK)
}