package restaurant

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
)

type Handler struct {
	store types.RestaurantStore
}

func NewHandler(s types.RestaurantStore) *Handler {
	return &Handler{store: s}
}

func (h *Handler) RegisterRouter(r *mux.Router) {
	r.HandleFunc("/restaurant", h.GetRestaurant).Methods("GET")
    r.HandleFunc("/restaurant/{id}", h.GetRestaurantById).Methods("GET")
    r.HandleFunc("/restaurantWorker", h.GetRestaurantWorker).Methods("GET")
    r.HandleFunc("/restaurantWorker/{id}", h.GetRestaurantWorkerById).Methods("GET")
    r.HandleFunc("/restaurantWorker/{id}/feedback", h.GetRestaurantWorkerFeedback).Methods("GET")
    r.HandleFunc("/reservation", h.GetReservation).Methods("GET")
    r.HandleFunc("/reservation/{id}", h.GetReservationById).Methods("GET")

}

func (h *Handler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurant, err := h.store.GetRestaurant()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, restaurant)
}

func (h *Handler) GetRestaurantById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	restaurant, err := h.store.GetRestaurantById(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, restaurant)
}

func (h *Handler) GetRestaurantWorker(w http.ResponseWriter, r *http.Request) {
	restaurant, err := h.store.GetRestaurantWorker()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, restaurant)
}

func (h *Handler) GetRestaurantWorkerById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	restaurant, err := h.store.GetRestaurantWorkerById(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, restaurant)
}

func (h *Handler) GetRestaurantWorkerFeedback(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	restaurant, err := h.store.GetRestaurantWorkerFeedBack(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, restaurant)
}


func (h *Handler) GetReservation(w http.ResponseWriter, r *http.Request) {
    restaurant, err := h.store.GetReservation()
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, err)
        return
    }
    utils.WriteJson(w, http.StatusOK, restaurant)
}

func (h *Handler) GetReservationById(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    restaurant, err := h.store.GetReservationById(id)
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, err)
        return
    }
    utils.WriteJson(w, http.StatusOK, restaurant)
}

