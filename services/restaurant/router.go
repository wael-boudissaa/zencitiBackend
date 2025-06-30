package restaurant

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	//!NOTE: RESTAURANT INFORMATIONS
	r.HandleFunc("/restaurant", h.GetRestaurant).Methods("GET")
	r.HandleFunc("/restaurant", h.CreateRestaurant).Methods("POST")
	r.HandleFunc("/restaurant/count/{restaurantId}", h.RestaurantCountInformation).Methods("GET")
	r.HandleFunc("/restaurant/{id}", h.GetRestaurantById).Methods("GET")
	r.HandleFunc("/menu/actif/{restaurantId}", h.GetAvailableMenuInformation).Methods("GET")
	r.HandleFunc("/restaurant/workers/{idRestaurant}", h.GetRestaurantWorker).Methods("GET")
	r.HandleFunc("/restaurant/token/{token}", h.GetRestaurantByToken).Methods("GET")
	r.HandleFunc("/restaurant/stats/{restaurantId}", h.RestaurtantStats).Methods("GET")
	r.HandleFunc("/restaurant/tables/occupany/today/{restaurantId}", h.GetTableOccupationToday).Methods("GET")
	r.HandleFunc("/restaurant/food/populair/{restaurantId}", h.GetTopFoodsThisWeek).Methods("GET")
	r.HandleFunc("/restaurant/tables", h.GetRestaurantTables).Methods("POST")
	r.HandleFunc("/restaurant/worker/{idRestaurant}", h.CreateRestaurantWorker).Methods("POST")
	r.HandleFunc("/restaurant/worker/fire/{idRestaurantWorker}", h.FireRestaurantWorker).Methods("POST")
	r.HandleFunc("/food/unavailable/{idFood}", h.SetFoodUnavailable).Methods("POST")
	r.HandleFunc("/food", h.CreateFood).Methods("POST")
	r.HandleFunc("/food/category", h.CreateFoodCategory).Methods("POST")
	r.HandleFunc("/food/category/{idRestaurant}", h.GetFoodCategoriesByRestaurant).Methods("GET")
	r.HandleFunc("/food/{idFood}", h.DeleteFood).Methods("DELETE")
	r.HandleFunc("/table/{idTable}", h.UpdateTable).Methods("PUT")
	r.HandleFunc("/table/{idTable}", h.DeleteTable).Methods("DELETE")
	r.HandleFunc("/tables/{restaurantId}", h.GetTablesByRestaurant).Methods("GET")
	r.HandleFunc("/restaurant/worker/{idRestaurantWorker}", h.UpdateRestaurantWorker).Methods("PUT")
	r.HandleFunc("/menu/restaurant/{idRestaurant}", h.GetMenusByRestaurant).Methods("GET")
	r.HandleFunc("/food/category/{idRestaurant}", h.GetFoodCategoriesByRestaurant).Methods("GET")
	r.HandleFunc("/food/active/{idRestaurant}", h.GetFoodsOfActiveMenu).Methods("GET")
	r.HandleFunc("/restaurant/menu/stats/{restaurantId}", h.GetRestaurantMenuStats).Methods("GET")
	r.HandleFunc("/restaurant/food/{restaurantId}", h.GetFoodRestaurant).Methods("GET")
	r.HandleFunc("/restaurant/addfood/{idMenu}", h.AddFoodToMenu).Methods("POST")
	//!NOTE: REVIEWS
	r.HandleFunc("/reviews/{idRestaurant}", h.GetRecentReviewsRestaurant).Methods("GET")
	r.HandleFunc("/friends/reviews", h.GetFriendsReviewsRestaurant).Methods("POST")
	r.HandleFunc("/restaurant/rating", h.PostReviewRestaurant).Methods("POST")
	//!NOTE: RESERVATION
	r.HandleFunc("/reservation/month/{restaurantId}", h.GetReservationMonthStats).Methods("GET")
	r.HandleFunc("/reservation", h.CreateReservation).Methods("POST")
	r.HandleFunc("/reservation/stats/{restaurantId}", h.GetReservationStats).Methods("GET")
	r.HandleFunc("/reservation/today/{restaurantId}", h.GetReservationToday).Methods("GET")
	r.HandleFunc("/reservation/{idReservation}/status", h.UpdateReservationStatus).Methods("PUT")
	r.HandleFunc("/reservation/upcoming/{restaurantId}", h.GetUpcomingReservations).Methods("GET")
	r.HandleFunc("/restaurant/{idRestaurant}/reservations", h.GetAllRestaurantReservations).Methods("GET")
	r.HandleFunc("/reservation/{idReservation}/details", h.GetReservationDetails).Methods("GET")

	//!NOTE: ORDER
	r.HandleFunc("/order", h.CreateOrder).Methods("POST")
	r.HandleFunc("/order/{idOrder}", h.GetOrderInformation).Methods("GET")
	r.HandleFunc("/wael/{restaurantId}", h.GetOrderStats).Methods("GET")
	r.HandleFunc("/waela/{clientId}", h.GetClientInf).Methods("GET")
	r.HandleFunc("/order/place", h.PostOrderClient).Methods("POST")
	r.HandleFunc("/food/{menuId}", h.GetFoodByMenu).Methods("GET")
	r.HandleFunc("/order/{idOrder}/status", h.UpdateOrderStatus).Methods("PUT")

	r.HandleFunc("/menu", h.CreateMenu).Methods("POST")
	r.HandleFunc("/food/{idFood}", h.GetFoodById).Methods("GET")
	r.HandleFunc("/food/{idFood}", h.UpdateFood).Methods("PUT")
	r.HandleFunc("/menu/{idMenu}", h.GetMenuWithFoods).Methods("GET")
	r.HandleFunc("/food/{idFood}/status", h.SetFoodStatusInMenu).Methods("PUT")

	r.HandleFunc("/client/{idClient}/reservations", h.GetAllClientReservations).Methods("GET")
	//!NOTE:NOTIFICATIONI NOT THIS PLACE
	r.HandleFunc("/notification", h.CreateNotification).Methods("POST")
	r.HandleFunc("/notification", h.GetNotifications).Methods("GET")

	//!NOTE: Tables
	r.HandleFunc("/restaurant/{idRestaurant}/tables/bulk", h.BulkUpdateRestaurantTables).Methods("PUT")
}

func (h *Handler) BulkUpdateRestaurantTables(w http.ResponseWriter, r *http.Request) {
	idRestaurant := mux.Vars(r)["idRestaurant"]
	if idRestaurant == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idRestaurant is required"))
		return
	}

	var req struct {
		Tables []types.Table `json:"data"`
	}

	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate table data
	for i, table := range req.Tables {
		if table.Shape == "" {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("table at index %d is missing shape", i))
			return
		}
		if table.PosX < 0 || table.PosY < 0 {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("table at index %d has invalid position coordinates", i))
			return
		}
	}

	err := h.store.BulkUpdateRestaurantTables(idRestaurant, req.Tables)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return the updated tables
	updatedTables, err := h.store.GetTablesByRestaurant(idRestaurant)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("tables updated but failed to retrieve: %v", err))
		return
	}

	utils.WriteJson(w, http.StatusOK, updatedTables)
}

func (h *Handler) GetReservationDetails(w http.ResponseWriter, r *http.Request) {
	idReservation := mux.Vars(r)["idReservation"]
	if idReservation == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idReservation is required"))
		return
	}

	details, err := h.store.GetReservationDetails(idReservation)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, details)
}

func (h *Handler) GetAllRestaurantReservations(w http.ResponseWriter, r *http.Request) {
	idRestaurant := mux.Vars(r)["idRestaurant"]
	if idRestaurant == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idRestaurant is required"))
		return
	}

	// Get page parameter (default to 1)
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Fixed limit of 6 per page
	limit := 6

	reservations, err := h.store.GetAllRestaurantReservations(idRestaurant, page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, reservations)
}

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idOrder := mux.Vars(r)["idOrder"]
	if idOrder == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idOrder is required"))
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate status
	validStatuses := []string{"pending", "completed", "cancelled"}
	isValid := false
	for _, validStatus := range validStatuses {
		if req.Status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid status. Valid statuses are: pending, completed, cancelled"))
		return
	}

	err := h.store.UpdateOrderStatus(idOrder, req.Status)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		if strings.Contains(err.Error(), "already completed") {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Order status updated to %s successfully", req.Status),
	})
}

func (h *Handler) GetOrderInformation(w http.ResponseWriter, r *http.Request) {
	idOrder := mux.Vars(r)["idOrder"]
	if idOrder == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("idOrder is required"))
		return
	}

	orderInfo, err := h.store.GetOrderInformation(idOrder)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, orderInfo)
}

func (h *Handler) GetAllClientReservations(w http.ResponseWriter, r *http.Request) {
	idClient := mux.Vars(r)["idClient"]
	if idClient == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("idClient is required"))
		return
	}

	reservations, err := h.store.GetAllClientReservations(idClient)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, reservations)
}

func (h *Handler) SetFoodStatusInMenu(w http.ResponseWriter, r *http.Request) {
	idFood := mux.Vars(r)["idFood"]
	var req struct {
		Status string `json:"status"` // "available" or "unavailable"
	}
	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, 400, err)
		return
	}
	err := h.store.SetFoodStatusInMenu(idFood, req.Status)
	if err != nil {
		utils.WriteError(w, 500, err)
		return
	}
	utils.WriteJson(w, 200, map[string]string{"message": "Status updated"})
}

func (h *Handler) AddFoodToMenu(w http.ResponseWriter, r *http.Request) {
	idMenu := mux.Vars(r)["idMenu"]
	var req struct {
		IdFood string `json:"idFood"`
	}
	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, 400, err)
		return
	}
	idMenuFood, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, 500, err)
		return
	}
	err = h.store.AddFoodToMenu(idMenuFood, idMenu, req.IdFood)
	if err != nil {
		utils.WriteError(w, 500, err)
		return
	}
	utils.WriteJson(w, 201, map[string]string{"idMenuFood": idMenuFood})
}

func (h *Handler) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, 400, err)
		return
	}

	idRestaurant, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, 500, err)
		return
	}
	idAdminRestaurant := r.FormValue("idAdminRestaurant")
	name := r.FormValue("name")
	description := r.FormValue("description")
	location := r.FormValue("location")
	capacity, _ := strconv.Atoi(r.FormValue("capacity"))
	longitude, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	latitude, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)

	file, _, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, 400, err)
		return
	}
	defer file.Close()
	imageURL, err := utils.UploadImageToCloudinary(file)
	if err != nil {
		utils.WriteError(w, 500, err)
		return
	}

	err = h.store.CreateRestaurant(idRestaurant, idAdminRestaurant, name, imageURL, longitude, latitude, description, capacity, location)
	if err != nil {
		utils.WriteError(w, 500, err)
		return
	}
	utils.WriteJson(w, 201, map[string]string{"idRestaurant": idRestaurant, "image": imageURL})
}

func (h *Handler) GetFoodRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("restaurantId is required"))
		return
	}
	reservations, err := h.store.GetFoodRestaurant(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, reservations)
}

func (h *Handler) GetUpcomingReservations(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("restaurantId is required"))
		return
	}
	reservations, err := h.store.GetUpcomingReservations(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, reservations)
}

func (h *Handler) GetRestaurantMenuStats(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("restaurantId is required"))
		return
	}
	stats, err := h.store.GetRestaurantMenuStats(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, stats)
}

func (h *Handler) GetFoodsOfActiveMenu(w http.ResponseWriter, r *http.Request) {
	idRestaurant := mux.Vars(r)["idRestaurant"]
	if idRestaurant == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idRestaurant is required"))
		return
	}
	foods, err := h.store.GetFoodsOfActiveMenu(idRestaurant)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, foods)
}

func (h *Handler) GetMenuWithFoods(w http.ResponseWriter, r *http.Request) {
	idMenu := mux.Vars(r)["idMenu"]
	menu, foods, err := h.store.GetMenuWithFoods(idMenu)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"menu":  menu,
		"foods": foods,
	})
}

func (h *Handler) GetMenusByRestaurant(w http.ResponseWriter, r *http.Request) {
	idRestaurant := mux.Vars(r)["idRestaurant"]
	if idRestaurant == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idRestaurant is required"))
		return
	}
	menus, err := h.store.GetMenusByRestaurant(idRestaurant)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, menus)
}

func (h *Handler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	idTable := mux.Vars(r)["idTable"]
	var table types.Table
	if err := utils.ParseJson(r, &table); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.store.UpdateTable(idTable, table); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Table updated"})
}

func (h *Handler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	idTable := mux.Vars(r)["idTable"]
	if err := h.store.DeleteTable(idTable); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Table deleted"})
}

func (h *Handler) GetTablesByRestaurant(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["restaurantId"]
	tables, err := h.store.GetTablesByRestaurant(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, tables)
}

func (h *Handler) UpdateReservationStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["idReservation"]
	var req struct {
		Status string `json:"status"`
	}
	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.store.UpdateReservationStatus(id, req.Status); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Reservation status updated"})
}

func (h *Handler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var notif types.Notification
	if err := utils.ParseJson(r, &notif); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if notif.IdNotification == "" {
		id, _ := utils.CreateAnId()
		notif.IdNotification = id
	}
	if err := h.store.CreateNotification(notif); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, notif)
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	notifs, err := h.store.GetNotifications()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, notifs)
}

func (h *Handler) CreateTable(w http.ResponseWriter, r *http.Request) {
	var table types.Table
	if err := utils.ParseJson(r, &table); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if table.IdTable == "" {
		id, _ := utils.CreateAnId()
		table.IdTable = id
	}
	if err := h.store.CreateTable(table); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, table)
}

func (h *Handler) CreateRestaurantWorker(w http.ResponseWriter, r *http.Request) {
	idRestaurant := mux.Vars(r)["idRestaurant"]

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var worker types.RestaurantWorkerCreation

	worker.FirstName = r.FormValue("firstName")
	worker.LastName = r.FormValue("lastName")
	worker.Email = r.FormValue("email")
	worker.PhoneNumber = r.FormValue("phoneNumber")
	worker.Quote = r.FormValue("quote")
	worker.Nationnallity = r.FormValue("nationnallity")
	worker.NativeLanguage = r.FormValue("nativeLanguage")
	worker.Address = r.FormValue("address")

	id, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageURL, err := utils.UploadImageToCloudinary(file)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		worker.Image = imageURL
	}

	if err := h.store.CreateRestaurantWorker(id, idRestaurant, worker); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{
		"message": "Worker created",
		"id":      id,
	})
}

func (h *Handler) UpdateRestaurantWorker(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["idRestaurantWorker"]
	var worker types.RestaurantWorker
	if err := utils.ParseJson(r, &worker); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.store.UpdateRestaurantWorker(id, worker); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Worker updated"})
}

func (h *Handler) FireRestaurantWorker(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["idRestaurantWorker"]
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idRestaurantWorker is required"))
		return
	}
	if err := h.store.SetRestaurantWorkerStatus(id, "inactive"); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Worker status set to inactive"})
}

func (h *Handler) GetFoodById(w http.ResponseWriter, r *http.Request) {
	idFood := mux.Vars(r)["idFood"]
	food, err := h.store.GetFoodById(idFood)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, food)
}

func (h *Handler) DeleteFood(w http.ResponseWriter, r *http.Request) {
	idFood := mux.Vars(r)["idFood"]
	if err := h.store.DeleteFood(idFood); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Food deleted"})
}

func (h *Handler) GetFoodCategoriesByRestaurant(w http.ResponseWriter, r *http.Request) {
	idRestaurant := mux.Vars(r)["idRestaurant"]
	if idRestaurant == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idRestaurant is required"))
		return
	}
	categories, err := h.store.GetFoodCategoriesByRestaurant(idRestaurant)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, categories)
}

func (h *Handler) CreateFoodCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NameCategorie string `json:"nameCategorie"`
	}
	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	idCategory, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.store.CreateFoodCategory(idCategory, req.NameCategorie); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"idCategory": idCategory})
}

func (h *Handler) SetFoodUnavailable(w http.ResponseWriter, r *http.Request) {
	idFood := mux.Vars(r)["idFood"]
	if idFood == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idFood is required"))
		return
	}
	if err := h.store.SetFoodUnavailable(idFood); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Food set to unavailable"})
}

func (h *Handler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IdRestaurant string `json:"idRestaurant"`
		Name         string `json:"name"`
	}
	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	idMenu, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.store.CreateMenu(idMenu, req.IdRestaurant, req.Name); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"idMenu": idMenu})
}

func (h *Handler) CreateFood(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	idFood, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	idCategory := r.FormValue("idCategory")
	idRestaurant := r.FormValue("idRestaurant") // Add this field
	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	status := r.FormValue("status")

	// Convert price string to float64
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid price format"))
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()
	imageURL, err := utils.UploadImageToCloudinary(file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.store.CreateFood(idFood, idCategory, idRestaurant, name, description, imageURL, price, status); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"idFood": idFood, "image": imageURL})
}

func (h *Handler) UpdateFood(w http.ResponseWriter, r *http.Request) {
	idFood := mux.Vars(r)["idFood"]
	var food types.Food
	if err := utils.ParseJson(r, &food); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.store.UpdateFood(idFood, food); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Food updated"})
}

func (h *Handler) GetTopFoodsThisWeek(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	foods, err := h.store.GetTopFoodsThisWeek(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, foods)
}

func (h *Handler) GetTableOccupationToday(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	tables, err := h.store.GetTableOccupationToday(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, tables)
}

func (h *Handler) GetReservationStats(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]

	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	stats, err := h.store.GetReservationStatsAndList(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if stats == nil {
		utils.WriteJson(w, http.StatusOK, map[string]interface{}{
			"totalReservations":    0,
			"upcomingReservations": []types.ReservationListInformation{},
		})
		return
	}
	utils.WriteJson(w, http.StatusOK, stats)
}

func (h *Handler) RestaurtantStats(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	stats, err := h.store.GetRestaurantRatingStats(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if stats == nil {
		utils.WriteJson(w, http.StatusOK, map[string]interface{}{
			"averageRating": 0.0,
			"totalReviews":  0,
		})
		return
	}
	utils.WriteJson(w, http.StatusOK, stats)
}

func (h *Handler) GetRestaurantByToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, ok := vars["token"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("token is required"))
		return
	}
	userInfo, err := utils.DecodeToken(token)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
	}

	if userInfo == nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
		return
	}
	roleInterface := userInfo["role"]
	role, ok := roleInterface.(string)
	if !ok || role == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("role is required"))
		return
	}
	idInterface := userInfo["id"]
	id, ok := idInterface.(string)
	if !ok || id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is required"))
		return
	}
	fmt.Println("role:", role)
	fmt.Println("id:", id)
	if role == "adminRestaurant" {
		restaurant, err := h.store.GetRestaurantByIdProfile(id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if restaurant == nil {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("restaurant not found"))
			return
		}
		utils.WriteJson(w, http.StatusOK, restaurant)
		return
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("access denied for adminRestaurant role"))
	}
}

// GetRecentReviewsRestaurant retrieves recent reviews for a specific restaurant.
func (h *Handler) GetRestaurantWorker(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["idRestaurant"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	workers, err := h.store.GetRestaurantWorker(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if workers == nil {
		utils.WriteJson(w, http.StatusOK, []types.RestaurantWorker{})
		return
	}
	utils.WriteJson(w, http.StatusOK, workers)
}

func (h *Handler) GetRecentReviewsRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["idRestaurant"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	recentReviews, err := h.store.GetRecentReviews(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if recentReviews == nil {
		utils.WriteJson(w, http.StatusOK, []types.Rating{})
		return
	}
	utils.WriteJson(w, http.StatusOK, recentReviews)
}

func (h *Handler) GetClientInf(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	if clientId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("clientId is required"))
		return
	}
	clientDetails, err := h.store.GetClientReservationAndOrderDetails(clientId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if clientDetails == nil {
		utils.WriteJson(w, http.StatusOK, map[string]interface{}{
			"profile":      types.Profile{},
			"reservations": []types.ReservationDetails{},
			"orders":       []types.OrderDetails{},
			"totalSpent":   0.0,
		})
		return
	}
	utils.WriteJson(w, http.StatusOK, clientDetails)
}

func (h *Handler) GetOrderStats(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}

	// Fetch order stats
	orderStatsByHour, orderStatsByStatus, err := h.store.GetOrderStatsByHourAndStatus(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Fetch recent orders
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			utils.WriteError(w, http.StatusBadRequest, errors.New("invalid limit parameter"))
			return
		}
		limit = parsedLimit
	}

	recentOrders, err := h.store.GetRecentOrders(restaurantId, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Combine stats and recent orders in the response
	utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"hourlyStats":  orderStatsByHour,
		"statusStats":  orderStatsByStatus,
		"recentOrders": recentOrders,
	})
}

func (h *Handler) GetReservationMonthStats(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	countNumberOfReservation, err := h.store.CountReservationLastMonth(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if countNumberOfReservation == nil {
		utils.WriteJson(w, http.StatusOK, []types.ReservationStats{})
		return
	}
	utils.WriteJson(w, http.StatusOK, countNumberOfReservation)
}

func (h *Handler) RestaurantCountInformation(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	countNumberofReservation, err := h.store.CountReservationReceivedToday(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	countNumberOfOrders, err := h.store.CountOrderReceivedToday(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	countFirstTimeUsers, err := h.store.CountFirstTimeReservers(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	result := map[string]int{
		"numberReservation": countNumberofReservation,
		"numberOrders":      countNumberOfOrders,
		"firstTimeUsers":    countFirstTimeUsers,
	}
	utils.WriteJson(w, http.StatusOK, result)
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

	log.Println("Order ID:", idOrder)
	err = h.store.PostOrderList(idOrder, orderCreation.Foods)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	log.Println("Order list posted successfully")

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

func (h *Handler) GetReservationToday(w http.ResponseWriter, r *http.Request) {
	restaurantId := mux.Vars(r)["restaurantId"]
	if restaurantId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("restaurantId is required"))
		return
	}
	reservations, err := h.store.GetReservationTodayByRestaurantId(restaurantId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, reservations)
}

// func (h *Handler) GetOrderInformation(w http.ResponseWriter, r *http.Request) {
// 	orderId := mux.Vars(r)["orderId"]
// 	if orderId == "" {
// 		utils.WriteError(w, http.StatusBadRequest, errors.New("orderId is required"))
// 		return
// 	}
// 	orderInfo, err := h.store.GetOrderInformation(orderId)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	utils.WriteJson(w, http.StatusOK, orderInfo)
// }

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
