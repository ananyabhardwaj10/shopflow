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
		respondWithError(w, http.StatusBadRequest, "Incorrect User Information", err)
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

func (cfg *apiConfig) handlerGetCartItems(w http.ResponseWriter, req *http.Request) {
	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect User Information", err)
		return 
	}

	cartItems, err := cfg.db.GetAllCartItems(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get cart items at the moment", err)
		return 
	}

	type response struct {
		UserID uuid.UUID `json:"user_id"`
		ProductID uuid.UUID `json:"product_id"`
		CartItemID uuid.UUID `json:"cart_item_id"`
		ProductName string `json:"product_name"`
		Price string `json:"price"`
		Quantity int32 `json:"quantity"`
	}

	allItems := []response{}

	for _, cartItem := range cartItems {
		allItems = append(allItems, response {
			UserID: userID,
			ProductID: cartItem.ProductID,
			CartItemID: cartItem.ID,
			ProductName: cartItem.Name,
			Price: cartItem.Price,
			Quantity: cartItem.Quantity,
		})
	}

	respondWithJSON(w, http.StatusOK, allItems)
}

func (cfg *apiConfig) handlerUpdateCartItemQuantity(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Quantity int32 `json:"quantity"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect update information", err)
		return 
	}

	cartItemIDStr := req.PathValue("id")
	cartItemID, err := uuid.Parse(cartItemIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get cart item", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect User Information", err)
		return 
	}

	cartItem, err := cfg.db.UpdateCartItemQuantity(req.Context(), database.UpdateCartItemQuantityParams{
		Quantity: params.Quantity,
		ID: cartItemID,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update item quantity", err)
		return 
	}

	type response struct {
		UserID uuid.UUID `json:"user_id"`
		CartItemID uuid.UUID `json:"cart_item_id"`
		UpdatedQuantity int32 `json:"updated_quantity"`
	}

	respondWithJSON(w, http.StatusOK, response{
		UserID: cartItem.UserID,
		CartItemID: cartItem.ID,
		UpdatedQuantity: cartItem.Quantity,
	})
}

func (cfg *apiConfig) handlerDeleteItemFromCart(w http.ResponseWriter, req *http.Request) {
	cartItemIDStr := req.PathValue("id")
	cartItemID, err := uuid.Parse(cartItemIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get cart item", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect User Information", err)
		return 
	}

	err = cfg.db.DeleteItemFromCart(req.Context(), database.DeleteItemFromCartParams{
		ID: cartItemID, 
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete item from cart", err)
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}