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
	r.HandleFunc("/activity/type/create", h.CreateActivityCategory).Methods("POST")
	r.HandleFunc("/activity/notAvailable", h.GetActivityNotAvailableAtDay).Methods("POST")
	r.HandleFunc("/client/{idClient}/activities", h.GetAllClientActivities).Methods("GET")
	r.HandleFunc("/activity/complete", h.CompleteClientActivity).Methods("POST")
	r.HandleFunc("/locations", h.GetAllLocationsWithDistances).Methods("POST")
	r.HandleFunc("/admin/{idAdminActivity}/activities", h.GetActivitiesByAdmin).Methods("GET")
	// By Activity id
	// r.HandleFunc("/admin/{idAdminActivity}/stats", h.GetActivityStats).Methods("GET")

	r.HandleFunc("/admin/{idAdminActivity}/stats", h.GetActivityStats).Methods("GET")
	r.HandleFunc("/booking/{idClientActivity}/status", h.UpdateActivityBookingStatus).Methods("PUT")
	r.HandleFunc("/activity/rating", h.PostReviewActivity).Methods("POST")
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
		IdAdminActivity  string `json:"idAdminActivity,omitempty"` // Optional, can be used for admin-specific logic
	}

	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if req.IdClientActivity == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idClientActivity is required"))
		return
	}

	err := h.store.UpdateClientActivityStatus(req.IdClientActivity, req.IdAdminActivity)
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

func (h *Handler) GetActivityNotAvailableAtDay(w http.ResponseWriter, r *http.Request) {
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
	unavailableTimes, err := h.store.GetActivityNotAvaialableAtday(day, req.IdActivity)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return times when activity is at full capacity (unavailable)

	utils.WriteJson(w, http.StatusOK, unavailableTimes)
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
	activite, err := h.store.GetActivityFullDetails(id)
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

func (h *Handler) GetActivityStats(w http.ResponseWriter, r *http.Request) {
	idAdminActivity := mux.Vars(r)["idAdminActivity"]
	if idAdminActivity == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idAdminActivity is required"))
		return
	}

	stats, err := h.store.GetActivityStatsAdmin(idAdminActivity)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, stats)
}

func (h *Handler) GetActivitiesByAdmin(w http.ResponseWriter, r *http.Request) {
	idAdminActivity := mux.Vars(r)["idAdminActivity"]
	if idAdminActivity == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idAdminActivity is required"))
		return
	}

	activities, err := h.store.GetActivitiesByAdminActivity(idAdminActivity)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, activities)
}

func (h *Handler) UpdateActivityBookingStatus(w http.ResponseWriter, r *http.Request) {
	idClientActivity := mux.Vars(r)["idClientActivity"]
	if idClientActivity == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idClientActivity is required"))
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
	validStatuses := []string{"pending", "confirmed", "cancelled", "completed"}
	isValid := false
	for _, validStatus := range validStatuses {
		if req.Status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid status. Valid statuses are: pending, confirmed, cancelled, completed"))
		return
	}

	err := h.store.UpdateActivityStatus(idClientActivity, req.Status)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Activity booking status updated to %s successfully", req.Status),
	})
}

func (h *Handler) PostReviewActivity(w http.ResponseWriter, r *http.Request) {
	var review types.PostRatingActivity
	if err := utils.ParseJson(r, &review); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if review.IdClient == "" || review.IdActivity == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("idClient and idActivity are required"))
		return
	}
	idReview, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	review.IdRating = idReview
	err = h.store.PostRatingActivity(review)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "Review posted successfully"})
}

// CreateActivityCategory creates a new activity category
func (h *Handler) CreateActivityCategory(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error parsing form: %v", err))
		return
	}

	// Get form values
	nameTypeActivity := r.FormValue("nameTypeActivity")
	if nameTypeActivity == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("nameTypeActivity is required"))
		return
	}

	// Handle image upload
	file, _, err := r.FormFile("imageActivity")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("imageActivity is required"))
		return
	}
	defer file.Close()

	// Upload image to Cloudinary
	imageURL, err := utils.UploadImageToCloudinary(file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error uploading image: %v", err))
		return
	}

	// Create the category data
	categoryData := types.ActivityCategoryCreation{
		NameTypeActivity: nameTypeActivity,
		ImageActivity:    imageURL,
	}

	// Create the category in the database
	categoryID, err := h.store.CreateActivityCategory(categoryData)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create activity category: %v", err))
		return
	}

	// Return success response with the new category ID
	response := map[string]interface{}{
		"message":        "Activity category created successfully",
		"idTypeActivity": categoryID,
		"imageURL":       imageURL,
	}
	utils.WriteJson(w, http.StatusCreated, response)
}
