package activite

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
	r.HandleFunc("/activity/single/{id}", h.GetActiviteById).Methods("GET")
	r.HandleFunc("/activity/create", h.CreateActivite).Methods("POST")
	r.HandleFunc("/activity/populaire", h.GetPopulaireActivity).Methods("GET")
	r.HandleFunc("/activity/recent/{idClient}", h.GetRecentActivite).Methods("GET")
	r.HandleFunc("/activity/type/{type}", h.GetActiviteByType).Methods("GET")
	r.HandleFunc("/activity/type", h.GetActiviteTypes).Methods("GET")
	r.HandleFunc("/activity/notAvailable", h.GetActivityNotAvaialbaleAtday).Methods("POST")
	r.HandleFunc("/client/{idClient}/activities", h.GetAllClientActivities).Methods("GET")
	r.HandleFunc("/activity/complete", h.CompleteClientActivity).Methods("POST")
	r.HandleFunc("/locations", h.GetAllLocationsWithDistances).Methods("POST")
}

func (h *Handler) GetAllLocationsWithDistances(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("latitude must be between -90 and 90"))
		return
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("longitude must be between -180 and 180"))
		return
	}

	locations, err := h.store.GetAllLocationsWithDistances(req.Latitude, req.Longitude)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if locations == nil || len(*locations) == 0 {
		utils.WriteJson(w, http.StatusOK, []types.LocationItemWithDistance{})
		return
	}

	utils.WriteJson(w, http.StatusOK, locations)
}

func (h *Handler) GetAllClientActivities(w http.ResponseWriter, r *http.Request) {
	idClient := mux.Vars(r)["idClient"]
	if idClient == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("idClient is required"))
		return
	}

	activities, err := h.store.GetAllClientActivities(idClient)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, activities)
}

func (h *Handler) CompleteClientActivity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IdClientActivity string `json:"idClientActivity"`
	}

	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if req.IdClientActivity == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idClientActivity is required"))
		return
	}

	err := h.store.UpdateClientActivityStatus(req.IdClientActivity)
	if err != nil {
		// Check for specific error types
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		if strings.Contains(err.Error(), "within 2 hours") || strings.Contains(err.Error(), "already completed") {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Activity completed successfully",
	})
}

func (h *Handler) GetActivityNotAvaialbaleAtday(w http.ResponseWriter, r *http.Request) {
	var req types.TimeNotAvaialable
	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	day, err := time.Parse("2006-01-02", req.Day)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	log.Println("Day received:", req.Day)
	reservedTimes, err := h.store.GetActivityNotAvaialableAtday(day, req.IdActivity)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Convert []*time.Time to []string in "15:04" format (HH:MM)

	utils.WriteJson(w, http.StatusOK, reservedTimes)
} //	func (h *Handler) GetActivite(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) CreateActivite(w http.ResponseWriter, r *http.Request) {
	var activity types.ActivityCreation
	if err := utils.ParseJson(r, &activity); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	idClientActivity, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	if err := h.store.CreateActivityClient(idClientActivity, activity); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, idClientActivity)
}
