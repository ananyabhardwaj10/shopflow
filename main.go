package main
import (
	"net/http"
	"os"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/ananyabhardwaj10/shopflow/internal/database"
)

type apiConfig struct {
	db *database.Queries
	jwtSecretKey string 
}
func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening the database: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	mux := http.NewServeMux() 

	server := &http.Server{
		Addr: ":8086",
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Working Just Fine!"))
	})

	apiCfg := apiConfig{
		db: dbQueries,
		jwtSecretKey: jwtSecretKey,
	}

	//PUBLIC ROUTES
	mux.HandleFunc("POST /api/auth/register", apiCfg.handlerRegisterUser)
	mux.HandleFunc("POST /api/auth/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/auth/logout", apiCfg.handlerLogout)
	mux.HandleFunc("POST /api/auth/refresh", apiCfg.handlerRefreshTokens)

	//Protected Routes
	mux.Handle("PATCH /api/customer/profile", chain(http.HandlerFunc(apiCfg.handlerUpdateUserProfile), apiCfg.authMiddleware, roleMiddleware("customer", "seller"),))
	mux.Handle("PATCH /api/customer/password", chain(http.HandlerFunc(apiCfg.handlerChangePassword), apiCfg.authMiddleware, roleMiddleware("customer", "seller"),))

	server.ListenAndServe()
}