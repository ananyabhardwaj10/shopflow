package main
import(
	"time"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/shopflow/internal/auth"
	"github.com/ananyabhardwaj10/shopflow/internal/database"
)

func (cfg *apiConfig) handlerRegisterUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
		ContactNumber string `json:"contact_number"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incomplete Information. Please try again.", err)
		return 
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing the password", err)
		return 
	}

	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		FirstName: params.FirstName,
		LastName: params.LastName,
		ContactNumber: params.ContactNumber,
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user at the moment. Please try again", err)
		return 
	}

	type response struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		FirstName  string `json:"first_name"`
		LastName string `json:"last_name"`
		Email string `json:"email"`
		Password string `json:"-"`
	}

	respondWithJSON(w, http.StatusCreated, response{
		ID: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
	})

}