package activite

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
)

type Handler struct {
	store types.ActiviteStore
}

func NewHandler(s types.ActiviteStore) *Handler {
	return &Handler{store: s}
}

func (h *Handler) RegisterRouter(r *mux.Router) {
	r.HandleFunc("/activite", h.GetActivite).Methods("GET")
}

func (h *Handler) GetActivite(w http.ResponseWriter, r *http.Request) {
	activite, err := h.store.GetActivite()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, activite)
}


func (h *Handler) GetActiviteById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	activite, err := h.store.GetActiviteById(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, activite)
}

func (h *Handler) GetActiviteByType(w http.ResponseWriter, r *http.Request) {
    typeActivite := mux.Vars(r)["type"]
	activite, err := h.store.GetActiviteTypes(typeActivite)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, activite)
}
