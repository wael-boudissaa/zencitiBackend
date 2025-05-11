package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/handlers" // Import the CORS package
	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/zencitiBackend/services/activite"
	"github.com/wael-boudissaa/zencitiBackend/services/user"
)

type APISERVER struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APISERVER {
    user.NewAuth()
	return &APISERVER{addr: addr, db: db}
}

func (s *APISERVER) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
      subrouter := router.PathPrefix("/").Subrouter()
      // !NOTE : SUBROUTER FOR THE USER 

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

    activiteStore := activite.NewStore(s.db)
    activiteHandler := activite.NewHandler(activiteStore)
    activiteHandler.RegisterRouter(subrouter)


    // !NOTE : SUBROUTER FOR THE COMMANDES  

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
