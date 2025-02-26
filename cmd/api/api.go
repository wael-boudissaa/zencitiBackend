package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/handlers" // Import the CORS package
	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/marquinoBackend/services/categorie"
	"github.com/wael-boudissaa/marquinoBackend/services/commande"
	"github.com/wael-boudissaa/marquinoBackend/services/product"
	"github.com/wael-boudissaa/marquinoBackend/services/user"
)

type APISERVER struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APISERVER {
	return &APISERVER{addr: addr, db: db}
}

func (s *APISERVER) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	commandeStore := commande.NewStore(s.db)
	commandeHandler := commande.NewHandler(commandeStore)
	commandeHandler.RegisterRoutes(subrouter)

	categorieStore := categorie.NewStore(s.db)
	categorieHandler := categorie.NewHandler(categorieStore)
	categorieHandler.RegisterRoutes(subrouter)
	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	// Use CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Change to your frontend's origin if needed
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(router)

	// Handle the request with the CORS handler
	corsHandler.ServeHTTP(w, r)
}

// Run starts the HTTP server
func (s *APISERVER) Run() error {
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, s)
}
