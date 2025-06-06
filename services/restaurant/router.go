package restaurant

import (
	"errors"
	"log"
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
	r.HandleFunc("/reservation", h.CreateReservation).Methods("POST")
	r.HandleFunc("/order", h.CreateOrder).Methods("POST")
	r.HandleFunc("/friends/reviews", h.GetFriendsReviewsRestaurant).Methods("POST")
	r.HandleFunc("/restaurant/rating", h.PostReviewRestaurant).Methods("POST")
	r.HandleFunc("/menu/actif/{restaurantId}", h.GetAvailableMenuInformation).Methods("GET")

	// r.HandleFunc("/order/add", h.AddFoodToOrder).Methods("POST")
	r.HandleFunc("/order/place", h.PostOrderClient).Methods("POST")

	// r.HandleFunc("/restaurantWorker", h.GetRestaurantWorker).Methods("GET")
	// r.HandleFunc("/restaurantWorker/{id}", h.GetRestaurantWorkerById).Methods("GET")
	// r.HandleFunc("/restaurantWorker/{id}/feedback", h.GetRestaurantWorkerFeedback).Methods("GET")
	// r.HandleFunc("/reservation", h.GetReservation).Methods("GET")
	// r.HandleFunc("/reservation/{id}", h.GetReservationById).Methods("GET")
	r.HandleFunc("/restaurant/tables", h.GetRestaurantTables).Methods("POST")
	r.HandleFunc("/food/{menuId}", h.GetFoodByMenu).Methods("GET")
	// r.HandleFunc("/restauran/tables/status", h.GetStatusTables).Methods("GET")
}

func (h *Handler) GetAvailableMenuInformation(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	menu, err := h.store.GetAvailableMenuInformation(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, menu)
}

func (h *Handler) PostOrderClient(w http.ResponseWriter, r *http.Request) {
	var orderCreation types.OrderCreation
	if err := utils.ParseJson(r, &orderCreation); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if orderCreation.IdReservation == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("idReservation"))
		return
	}
	idOrder, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = h.store.CreateOrder(idOrder, orderCreation)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.PostOrderList(idOrder, orderCreation.Foods)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, "Success modifying price")
}

func (h *Handler) AddFoodToOrder(w http.ResponseWriter, r *http.Request) {
	var foodToOrder types.AddFoodToOrder
	if err := utils.ParseJson(r, &foodToOrder); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if foodToOrder.IdOrder == "" || foodToOrder.IdFood == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("idOrder and idFood are required"))
		return
	}
	err := h.store.AddFoodToOrder(foodToOrder)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, "Success adding food to order")
}

func (h *Handler) GetFoodByMenu(w http.ResponseWriter, r *http.Request) {
	menuId := mux.Vars(r)["menuId"]
	if menuId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("menuId is required"))
		return
	}
	food, err := h.store.GetFoodByMenu(menuId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, food)
}

func (h *Handler) GetRestaurantTables(w http.ResponseWriter, r *http.Request) {
	var tableRestaurant types.GetRestaurantTable

	if err := utils.ParseJson(r, &tableRestaurant); err != nil {

		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	log.Println("Parsed JSON:", tableRestaurant)
	if tableRestaurant.IdRestaurant == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}

	tables, err := h.store.GetRestaurantTables(tableRestaurant.IdRestaurant, tableRestaurant.TimeSlot)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	log.Println("Tables:", tables)

	utils.WriteJson(w, http.StatusOK, tables)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order types.OrderCreation
	if err := utils.ParseJson(r, &order); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if order.IdReservation == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("idReservation and idRestaurant are required"))
		return
	}
	idOrder, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = h.store.CreateOrder(idOrder, order)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, idOrder)
}

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation types.ReservationCreation
	if err := utils.ParseJson(r, &reservation); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if reservation.IdClient == "" || reservation.IdRestaurant == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("idClient and idRestaurant are required"))
		return
	}
	// if reservation.TimeFrom.In.IsZero() {
	// 	utils.WriteError(w, http.StatusBadRequest, errors.New("timeSlot is required"))
	// 	return
	// }
	idReservation, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateReservation(idReservation, reservation)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// err = h.store.ReserveTable(idReservation, reservation)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, err)
	// 	return
	// }

	utils.WriteJson(w, http.StatusCreated, idReservation)
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

// func (h *Handler) GetRestaurantWorker(w http.ResponseWriter, r *http.Request) {
// 	restaurant, err := h.store.GetRestaurantWorker()
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	utils.WriteJson(w, http.StatusOK, restaurant)
// }
//
// func (h *Handler) GetRestaurantWorkerById(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	restaurant, err := h.store.GetRestaurantWorkerById(id)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	utils.WriteJson(w, http.StatusOK, restaurant)
// }
//
// func (h *Handler) GetRestaurantWorkerFeedback(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	restaurant, err := h.store.GetRestaurantWorkerFeedBack(id)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	utils.WriteJson(w, http.StatusOK, restaurant)
// }
//
// func (h *Handler) GetReservation(w http.ResponseWriter, r *http.Request) {
// 	restaurant, err := h.store.GetReservation()
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	utils.WriteJson(w, http.StatusOK, restaurant)
// }
//
// func (h *Handler) GetReservationById(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	restaurant, err := h.store.GetReservationById(id)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	utils.WriteJson(w, http.StatusOK, restaurant)
// }

func (h *Handler) PostReviewRestaurant(w http.ResponseWriter, r *http.Request) {
	var review types.PostRatingRestaurant
	if err := utils.ParseJson(r, &review); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if review.IdClient == "" || review.IdRestaurant == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("idClient and idRestaurant are required"))
		return
	}
	idReview, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	review.IdRating = idReview
	err = h.store.PostRatingRestaurant(review)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "Review posted successfully"})
}

func (h *Handler) GetFriendsReviewsRestaurant(w http.ResponseWriter, r *http.Request) {
	var rating types.FriendsReviewsRestaruant
	if err := utils.ParseJson(r, &rating); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	friends, err := h.store.GetFriendsOfClient(rating.IdClient)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	reviews, err := h.store.GetRatingOfFriendsRestaurant(*friends, rating.IdRestaurant)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, reviews)
}
