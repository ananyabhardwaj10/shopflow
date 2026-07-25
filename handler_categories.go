package main
import(
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCategories(w http.ResponseWriter, req *http.Request) {
	categories, err := cfg.db.GetAllCategories(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get categories at the moment", err)
		return 
	}

	type response struct {
		CategoryID uuid.UUID `json:"category_id"`
		CategoryName string `json:"category_name"`
	}

	var allCategories []response

	for _, category := range categories {
		allCategories = append(allCategories, response {
			CategoryID: category.ID,
			CategoryName: category.Name,
		})
	}

	respondWithJSON(w, http.StatusOK, allCategories)
}