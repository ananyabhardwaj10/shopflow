package main
import(
	"net/http"
)

func (cfg *apiConfig) handlerRegisterUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
		Email string `json:"email"`
		
	}
}