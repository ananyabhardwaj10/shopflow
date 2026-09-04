package main 
import(
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/shopflow/internal/auth"
	"github.com/ananyabhardwaj10/shopflow/internal/database"
)

func (cfg *apiConfig) handlerPlaceOrder(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		DeliveryAddress string `json:"delivery_address"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Please provide valid delivery address", err)
		return 
	}

	if params.DeliveryAddress == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide valid delivery address", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect User Information", err)
		return 
	}

	cartItems, err := cfg.db.GetAllCartItems(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get cart at the moment", err)
		return
	}

	if len(cartItems) == 0 {
		respondWithError(w, http.StatusBadRequest, "Cart is Empty. Please add items first", err)
		return 
	}

	transaction, err := cfg.sqlDB.BeginTx(req.Context(), nil) 
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to start the order process. Please try again later", err)
		return 
	}

	defer transaction.Rollback()

	querytx := cfg.db.WithTx(transaction)

	var totalAmount float64

	for _, cartItem := range cartItems {
		cartItemPrice, err := strconv.ParseFloat(cartItem.Price, 64)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Issue processing the order. Please try again later", err)
			return 
		}

		totalAmount += cartItemPrice * float64(cartItem.Quantity)
	}

	totalAmtStr := fmt.Sprintf("%.2f", totalAmount)

	order, err := querytx.CreateOrder(req.Context(), database.CreateOrderParams{
		UserID: userID,
		TotalAmount: totalAmtStr,
		DeliveryAddress: params.DeliveryAddress,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to process order at the moment. Please try again later", err)
		return 
	}

	for _, cartItem := range cartItems {
		orderItem, err := querytx.CreateOrderItem(req.Context(), database.CreateOrderItemParams{
			OrderID: order.ID,
			ProductID: cartItem.ProductID,
			Quantity: cartItem.Quantity,
			PriceAtPurchase: cartItem.Price,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to process order at the moment. Please try again later", err)
			return 
		}

		_, err = querytx.ReduceProductStock(req.Context(), database.ReduceProductStockParams{
			StockQuantity: orderItem.Quantity,
			ID: orderItem.ProductID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to process order at the moment. Please try again later", err)
			return 
		}
	}

	err = querytx.ClearCartAfterOrder(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to process order at the moment. Please try again later", err)
		return 
	}

	err = transaction.Commit()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to process order at the moment. Please try again later", err)
		return 
	}

	type response struct {
		OrderID uuid.UUID `json:"order_id"`
		UserID uuid.UUID `json:"user_id"`
		Status string `json:"status"`
		TotalAmount string `json:"total_amount"`
	}

	respondWithJSON(w, http.StatusCreated, response{
		OrderID: order.ID, 
		UserID: userID, 
		Status: order.Status,
		TotalAmount: order.TotalAmount,
	})
}

func (cfg *apiConfig) handlerGetOrderHistory(w http.ResponseWriter, req *http.Request) {
	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect User Information", err)
		return 
	}

	orderList, err := cfg.db.GetOrderHistory(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get order history", err)
		return 
	}

	type response struct {
		OrderID uuid.UUID `json:"order_id"`
		UserID uuid.UUID `json:"user_id"`
		Status string `json:"status"`
		TotalAmount string `json:"total_amount"` 
	}

	allOrders := []response{}

	for _, o := range orderList {
		allOrders = append(allOrders, response{
			OrderID: o.ID, 
			UserID: o.UserID, 
			Status: o.Status, 
			TotalAmount: o.TotalAmount,
		})
	}

	respondWithJSON(w, http.StatusOK, allOrders)
}

func (cfg *apiConfig) handlerUpdateOrderStatus(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		OrderItemStatus string `json:"order_item_status"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Improper Information", err)
		return 
	}

	userID, err := auth.GetUserIDFromContext(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect User Information", err)
		return 
	}

	orderIdStr := req.PathValue("order_id")
	orderID, err := uuid.Parse(orderIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get order information", err)
		return 
	}

	orderItemIDStr := req.PathValue("item_id")
	orderItemID, err := uuid.Parse(orderItemIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get order information", err)
		return 
	}

	if params.OrderItemStatus != "confirmed" {
		respondWithError(w, http.StatusBadRequest, "Improper Information", err)
		return 
	}

	transaction, err := cfg.sqlDB.BeginTx(req.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to process order status change", err)
		return 
	}

	defer transaction.Rollback()

	querytx := cfg.db.WithTx(transaction)

	orderItem, err := querytx.UpdateOrderItemStatus(req.Context(), database.UpdateOrderItemStatusParams{
		Status: params.OrderItemStatus,
		ID: orderItemID,
		OrderID: orderID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to process order status change", err)
		return 
	}

	allUpdated, err := querytx.CheckAllOrderItemsConfirmed(req.Context(), orderID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to process order status change", err)
		return 
	}

	var order database.Order

	if allUpdated == 0 {
		order, err = querytx.UpdateOrderStatus(req.Context(), database.UpdateOrderStatusParams{
			Status: params.OrderItemStatus, 
			ID: orderID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to process order status change", err)
			return
		} 
	} 

	err = transaction.Commit()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to process order status change", err)
		return 
	}

	type response struct {
		OrderID uuid.UUID `json:"order_id"`
		OrderItemID uuid.UUID `json:"order_item_id"`
		UserID uuid.UUID `json:"user_id"`
		OrderStatus string `json:"order_status,omitempty"`
		OrderItemStatus string `json:"order_item_status"`
	}

	if allUpdated == 0 {
		respondWithJSON(w, http.StatusOK, response{
			OrderID: order.ID, 
			OrderItemID: orderItem.ID,
			UserID: userID, 
			OrderStatus: order.Status,
			OrderItemStatus: orderItem.Status,
		})
	} else {
		respondWithJSON(w, http.StatusOK, response{
			OrderID: orderID, 
			OrderItemID: orderItem.ID,
			UserID: userID, 
			OrderItemStatus: orderItem.Status,
		})
	}
}
