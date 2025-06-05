package activite

import (
	"log"
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
	// r.HandleFunc("/activite", h.GetActivite).Methods("GET")
	r.HandleFunc("/activite/single/{id}", h.GetActiviteById).Methods("GET")
	r.HandleFunc("/activity/populaire", h.GetPopulaireActivity).Methods("GET")
	r.HandleFunc("/activity/recent/{idClient}", h.GetRecentActivite).Methods("GET")
	r.HandleFunc("/activity/type/{type}", h.GetActiviteByType).Methods("GET")
	r.HandleFunc("/activity/type", h.GetActiviteTypes).Methods("GET")
}

//	func (h *Handler) GetActivite(w http.ResponseWriter, r *http.Request) {
//		activite, err := h.store.GetActivite()
//		if err != nil {
//			utils.WriteError(w, http.StatusBadRequest, err)
//			return
//		}
//		utils.WriteJson(w, http.StatusOK, activite)
//	}
//
//	func (h *Handler) GetActiviteById(w http.ResponseWriter, r *http.Request) {
//		id := mux.Vars(r)["id"]
//		activite, err := h.store.GetActiviteById(id)
//		if err != nil {
//			utils.WriteError(w, http.StatusBadRequest, err)
//			return
//		}
//		utils.WriteJson(w, http.StatusOK, activite)

//	}

func (h *Handler) GetActiviteById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("ID of activity:", id)
	activite, err := h.store.GetActiviteById(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// response := map[string]interface{}{
	// 	"message": "Success",
	// 	"data":    activite,
	// }

	utils.WriteJson(w, http.StatusOK, activite)
}

func (h *Handler) GetRecentActivite(w http.ResponseWriter, r *http.Request) {
	idClient := mux.Vars(r)["idClient"]
	// Get the recent activities for the client
	// This function should be implemented in the store to fetch recent activities

	activite, err := h.store.GetRecentActivities(idClient)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, activite)
}

func (h *Handler) GetActiviteByType(w http.ResponseWriter, r *http.Request) {
	typeActivite := mux.Vars(r)["type"]
	log.Println("Type of activity:", typeActivite)
	activite, err := h.store.GetActivityByTypes(typeActivite)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// response := map[string]interface{}{
	// 	"message": "Success",
	// 	"data":    activite,
	// }

	utils.WriteJson(w, http.StatusOK, activite)
}

func (h *Handler) GetActiviteTypes(w http.ResponseWriter, r *http.Request) {
	activite, err := h.store.GetActiviteTypes()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// response := map[string]interface{}{
	// 	"message": "Success",
	// 	"data":    activite,
	// }

	utils.WriteJson(w, http.StatusOK, activite)
}

func (h *Handler) GetPopulaireActivity(w http.ResponseWriter, r *http.Request) {
	activite, err := h.store.GetPopularActivities()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// response := map[string]interface{}{
	//     "message": "Success",
	//     "data":    activite,
	// }

	utils.WriteJson(w, http.StatusOK, activite)
}
