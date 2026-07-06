package main
import(
	"net/http"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/shopflow/internal/auth"
	"github.com/ananyabhardwaj10/shopflow/internal/database"
)

func (cfg *apiConfig) handlerSellerOnboarding(w http.ResponseWriter, req *http.Request) {
	const(
		role = "seller"
	)
	
	type parameters struct {
		BusinessName string	`json:"business_name"`
	}
	
	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode request body", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get user ID", err)
		return 
	}

	transaction, err := cfg.sqlDB.BeginTx(req.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to begin a transaction", err)
		return 
	}

	defer transaction.Rollback()

	querytx := cfg.db.WithTx(transaction)

	seller, err := querytx.CreateSeller(req.Context(), database.CreateSellerParams{
		BusinessName: params.BusinessName,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create seller", err)
		return 
	}

	user, err := querytx.UpdateUserRole(req.Context(), database.UpdateUserRoleParams{
		Role: role,
		ID: userID, 
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update role", err)
		return 
	}

	err = transaction.Commit()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to commit the transaction", err)
		return 
	}

	type response struct {
		UserID uuid.UUID `json:"user_id"`
		SellerID uuid.UUID `json:"seller_id"`
		UpdatedAt time.Time `json:"updated_at"`
		BusinessName string `json:"business_name"`
		Role string `json:"role"`
	}

	respondWithJSON(w, http.StatusCreated, response{
		UserID: user.ID,
		SellerID: seller.ID,
		UpdatedAt: seller.UpdatedAt,
		BusinessName: seller.BusinessName,
		Role: user.Role,
	})
}