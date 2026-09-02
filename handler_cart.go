package main
import(
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/shopflow/internal/auth"
	"github.com/ananyabhardwaj10/shopflow/internal/database"
)

func (cfg *apiConfig) handlerAddToCart(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		ProductID uuid.UUID `json:"product_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get Product id", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get user id", err)
		return 
	}

	cartItem, err := cfg.db.AddItemToCart(req.Context(), database.AddItemToCartParams{
		UserID: userID, 
		ProductID: params.ProductID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to add item to cart", err)
		return
	}

	type response struct {
		UserID uuid.UUID `json:"user_id"`
		ProductID uuid.UUID `json:"product_id"`
		CartItemID uuid.UUID `json:"cart_item_id"`
		Quantity int32 `json:"quantity"`
	}

	respondWithJSON(w, http.StatusCreated, response{
		UserID: userID,
		ProductID: params.ProductID,
		CartItemID: cartItem.ID,
		Quantity: cartItem.Quantity,
	})

}