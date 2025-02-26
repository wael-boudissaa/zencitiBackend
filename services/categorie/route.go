package categorie

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/marquinoBackend/types"
	"github.com/wael-boudissaa/marquinoBackend/utils"
)

type Handler struct {
	store types.CategorieStore
}

func NewHandler(store types.CategorieStore) *Handler {
	return &Handler{store: store}
}

func (s *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/categories", s.getCategorie).Methods("GET")
}

func (s *Handler) getCategorie(w http.ResponseWriter, r *http.Request) {
	categories, err := s.store.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, http.StatusOK, categories)

}
