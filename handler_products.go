package main
import(
	"fmt"
	"time"
	"net/http"
	"strconv"
	"encoding/json"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/shopflow/internal/auth"
	"github.com/ananyabhardwaj10/shopflow/internal/database"
)

func (cfg *apiConfig) handlerCreateProduct(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Name string `json:"product_name"`
		Description string `json:"product_description"`
		Price float64 `json:"price"`
		StockQuantity int32 `json:"stock_quantity"`
		CategoryID string `json:"category_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Improper Request. Please try again.", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get user id", err)
		return 
	}

	seller, err := cfg.db.GetSellerByUserID(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get seller id", err)
		return 
	}

	priceStr := fmt.Sprintf("%.2f", params.Price)

	var categoryID uuid.NullUUID
	if params.CategoryID != "" {
		parsed, err := uuid.Parse(params.CategoryID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid category ID", err)
			return 
		}

		categoryID = uuid.NullUUID{UUID: parsed, Valid: true}
	}

	product, err := cfg.db.CreateProduct(req.Context(), database.CreateProductParams{
		Name: params.Name,
		Description: params.Description,
		Price: priceStr,
		StockQuantity: params.StockQuantity,
		SellerID: seller.ID,
		CategoryID: categoryID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create product. Please try again.", err)
		return
	}

	type response struct {
		SellerID uuid.UUID `json:"seller_id"`
		ProductID uuid.UUID `json:"product_id"`
		Name string `json:"product_name"`
		Description string `json:"product_description"`
		Price string `json:"price"`
		StockQuantity int32 `json:"stock_quantity"`
		CategoryID uuid.NullUUID `json:"category_id"`
		CreatedAt time.Time `json:"created_at"`
	}

	respondWithJSON(w, http.StatusCreated, response{
		SellerID: seller.ID,
		ProductID: product.ID,
		Name: product.Name,
		Description: product.Description,
		Price: product.Price,
		StockQuantity: product.StockQuantity,
		CategoryID: product.CategoryID,
		CreatedAt: product.CreatedAt,
	})
}

func (cfg *apiConfig) handlerGetAllProductsBySeller(w http.ResponseWriter, req *http.Request) {
	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get user ID", err)
		return 
	}

	seller, err := cfg.db.GetSellerByUserID(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get seller information at the moment. Please try again.", err)
		return 
	}

	pageStr := req.URL.Query().Get("page")
	limitStr := req.URL.Query().Get("limit")
	var page int
	var limit int 

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to get page number from query", err)
			return 
		}
	} else {
		page = 1
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to extract the limit from query", err)
			return 
		}
	} else {
		limit = 10
	}

	offSet := (page - 1) * limit

	productList, err := cfg.db.GetAllProductsBySeller(req.Context(), database.GetAllProductsBySellerParams{
		SellerID: seller.ID,
		Limit: int32(limit),
		Offset: int32(offSet), 
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get all products by a seller. Please try again later.", err)
		return 
	}

	type response struct {
		SellerID uuid.UUID `json:"seller_id"`
		ProductID uuid.UUID `json:"product_id"`
		Name string `json:"product_name"`
		Description string `json:"product_description"`
		Price string `json:"price"`
		StockQuantity int32 `json:"stock_quantity"`
		CategoryID uuid.NullUUID `json:"category_id"`
		CreatedAt time.Time `json:"created_at"`
	}

	allProducts := []response{}

	for _, product := range productList {
		allProducts = append(allProducts, response{
			SellerID: seller.ID,
			ProductID: product.ID,
			Name: product.Name,
			Description: product.Description,
			Price: product.Price,
			StockQuantity: product.StockQuantity,
			CategoryID: product.CategoryID,
			CreatedAt: product.CreatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, allProducts)
}

func (cfg *apiConfig) handlerUpdateProduct(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Name *string `json:"product_name"`
		Description *string `json:"product_description"`
		Price *float64 `json:"price"`
		StockQuantity *int32 `json:"stock_quantity"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Improper information. Please try again.", err)
		return 
	}

	productIDStr := req.PathValue("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get product ID", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get User ID", err)
		return 
	}

	seller, err := cfg.db.GetSellerByUserID(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get seller information at the moment. Please try again.", err)
		return
	}

	name := sql.NullString{}
	description := sql.NullString{}
	price := sql.NullString{}
	stockQuantity := sql.NullInt32{}

	if params.Name != nil {
		name.String = *params.Name
		name.Valid = true
	}

	if params.Description != nil {
		description.String = *params.Description
		description.Valid = true
	}

	if params.Price != nil {
		price.String = fmt.Sprintf("%.2f", *params.Price)
		price.Valid = true 
	}

	if params.StockQuantity != nil {
		stockQuantity.Int32 = *params.StockQuantity
		stockQuantity.Valid = true 
	}

	product, err := cfg.db.UpdateProductDetails(req.Context(), database.UpdateProductDetailsParams{
		ID: productID, 
		SellerID: seller.ID,
		Name: name,
		Description: description,
		Price: price,
		StockQuantity: stockQuantity,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update product details at the moment. Please try again.", err)
		return 
	}

	type response struct {
		ID uuid.UUID `json:"product_id"`
		SellerID uuid.UUID `json:"seller_id"`
		CategoryID uuid.NullUUID `json:"category_id"`
		UpdatedAt time.Time `json:"updated_at"`
		Name string `json:"product_name"`
		Description string `json:"product_description"`
		Price string `json:"price"`
		StockQuantity int32 `json:"stock_quantity"`
	}

	respondWithJSON(w, http.StatusOK, response{
		ID: product.ID,
		SellerID: seller.ID,
		CategoryID: product.CategoryID,
		UpdatedAt: product.UpdatedAt,
		Name: product.Name,
		Description: product.Description,
		Price: product.Price,
		StockQuantity: product.StockQuantity,
	})
}

func (cfg *apiConfig) handlerDeleteProduct(w http.ResponseWriter, req *http.Request) {
	productIDStr := req.PathValue("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get product id", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get user id", err)
		return 
	}

	seller, err := cfg.db.GetSellerByUserID(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get seller information.", err)
		return 
	}

	err = cfg.db.DeleteProduct(req.Context(), database.DeleteProductParams{
		ID: productID, 
		SellerID: seller.ID, 
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete product. Please try again later", err)
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}