package sensors

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
)

type Handler struct {
	store types.SensorStore
}

func NewHandler(store types.SensorStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Sensor registration and management routes
	router.HandleFunc("/sensors/register", h.RegisterSensor).Methods("POST")
	router.HandleFunc("/sensors/user/{idClient}", h.GetUserSensors).Methods("GET")
	router.HandleFunc("/sensors/{sensorId}/info", h.GetSensorInfo).Methods("GET")
	router.HandleFunc("/sensors/{sensorId}/status", h.UpdateSensorStatus).Methods("PUT")

	// Daily usage routes
	router.HandleFunc("/sensors/daily-usage", h.SaveDailyUsage).Methods("POST")
	router.HandleFunc("/sensors/batch-usage", h.SaveBatchUsage).Methods("POST")
	router.HandleFunc("/sensors/user/{idClient}/usage", h.GetSensorUsage).Methods("GET")
	router.HandleFunc("/sensors/{sensorId}/usage/range", h.GetSensorUsageByDateRange).Methods("GET")
}

// RegisterSensor handles sensor registration to a user
func (h *Handler) RegisterSensor(w http.ResponseWriter, r *http.Request) {
	var registration types.SensorRegistration
	if err := utils.ParseJson(r, &registration); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate input
	if registration.SensorId == "" || registration.ClientId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("sensorId and clientId are required"))
		return
	}

	// Validate sensor ID format (ZC-WS-YYYY-NNNN)
	if !isValidSensorId(registration.SensorId) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid sensor ID format. Expected format: ZC-WS-YYYY-NNNN"))
		return
	}

	// Check if sensor already exists and is owned by another client
	existingClientId, err := h.store.CheckSensorOwnership(registration.SensorId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if existingClientId != "" && existingClientId != registration.ClientId {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("sensor already registered to another client"))
		return
	}

	// Register the sensor
	err = h.store.RegisterSensor(registration.SensorId, registration.ClientId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Sensor registered successfully",
		"sensorId": registration.SensorId,
		"status": "active",
	})
}

// GetUserSensors returns all sensors registered to a user
func (h *Handler) GetUserSensors(w http.ResponseWriter, r *http.Request) {
	idClient := mux.Vars(r)["idClient"]
	if idClient == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idClient is required"))
		return
	}

	sensors, err := h.store.GetSensorsByClient(idClient)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := types.UserSensorsResponse{
		Sensors:      sensors,
		TotalSensors: len(sensors),
		HasSensors:   len(sensors) > 0,
	}

	utils.WriteJson(w, http.StatusOK, response)
}

// GetSensorInfo returns information about a specific sensor
func (h *Handler) GetSensorInfo(w http.ResponseWriter, r *http.Request) {
	sensorId := mux.Vars(r)["sensorId"]
	if sensorId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("sensorId is required"))
		return
	}

	sensorInfo, err := h.store.GetSensorInfo(sensorId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("sensor not found"))
		return
	}

	utils.WriteJson(w, http.StatusOK, sensorInfo)
}

// UpdateSensorStatus updates the status of a sensor (active/inactive)
func (h *Handler) UpdateSensorStatus(w http.ResponseWriter, r *http.Request) {
	sensorId := mux.Vars(r)["sensorId"]
	if sensorId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("sensorId is required"))
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := utils.ParseJson(r, &statusUpdate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate status
	if statusUpdate.Status != "active" && statusUpdate.Status != "inactive" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("status must be either 'active' or 'inactive'"))
		return
	}

	err := h.store.UpdateSensorStatus(sensorId, statusUpdate.Status)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Sensor status updated to %s", statusUpdate.Status),
		"sensorId": sensorId,
		"status": statusUpdate.Status,
	})
}

// SaveDailyUsage handles daily water usage data from sensors
func (h *Handler) SaveDailyUsage(w http.ResponseWriter, r *http.Request) {
	var usage types.DailyUsageData
	if err := utils.ParseJson(r, &usage); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate input
	if usage.SensorId == "" || usage.UsageDate == "" || usage.VolumeLiters < 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("sensorId, usageDate, and volumeLiters are required"))
		return
	}

	// Validate date format (YYYY-MM-DD)
	_, err := time.Parse("2006-01-02", usage.UsageDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid date format, use YYYY-MM-DD"))
		return
	}

	// Check if sensor exists and is active
	sensorInfo, err := h.store.GetSensorInfo(usage.SensorId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("sensor not found"))
		return
	}

	if sensorInfo.Status != "active" {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("sensor is not active"))
		return
	}

	// Save the usage data
	err = h.store.SaveDailyUsage(usage)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Daily usage data saved successfully",
		"sensorId": usage.SensorId,
		"usageDate": usage.UsageDate,
		"volumeLiters": usage.VolumeLiters,
	})
}

// SaveBatchUsage handles batch upload of usage data (for offline sensors)
func (h *Handler) SaveBatchUsage(w http.ResponseWriter, r *http.Request) {
	var batchData types.BatchUsageData
	if err := utils.ParseJson(r, &batchData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate input
	if batchData.SensorId == "" || len(batchData.UsageData) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("sensorId and usageData are required"))
		return
	}

	// Check if sensor exists and is active
	sensorInfo, err := h.store.GetSensorInfo(batchData.SensorId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("sensor not found"))
		return
	}

	if sensorInfo.Status != "active" {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("sensor is not active"))
		return
	}

	// Validate each usage record
	for i, usage := range batchData.UsageData {
		if usage.UsageDate == "" || usage.VolumeLiters < 0 {
			utils.WriteError(w, http.StatusBadRequest, 
				fmt.Errorf("invalid usage data at index %d: usageDate and volumeLiters are required", i))
			return
		}

		// Validate date format
		_, err := time.Parse("2006-01-02", usage.UsageDate)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest,
				fmt.Errorf("invalid date format at index %d, use YYYY-MM-DD", i))
			return
		}
	}

	// Save batch data
	err = h.store.SaveBatchUsage(batchData)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Batch usage data saved successfully (%d records)", len(batchData.UsageData)),
		"sensorId": batchData.SensorId,
		"recordsProcessed": len(batchData.UsageData),
	})
}

// GetSensorUsage returns usage analytics for all user's sensors
func (h *Handler) GetSensorUsage(w http.ResponseWriter, r *http.Request) {
	idClient := mux.Vars(r)["idClient"]
	if idClient == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idClient is required"))
		return
	}

	// Get period parameter (default to "month")
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "month"
	}

	// Validate period
	validPeriods := map[string]bool{"week": true, "month": true, "year": true}
	if !validPeriods[period] {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("period must be one of: week, month, year"))
		return
	}

	usage, err := h.store.GetSensorUsage(idClient, period)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, usage)
}

// GetSensorUsageByDateRange returns usage data for a sensor within a date range
func (h *Handler) GetSensorUsageByDateRange(w http.ResponseWriter, r *http.Request) {
	sensorId := mux.Vars(r)["sensorId"]
	if sensorId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("sensorId is required"))
		return
	}

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	if startDate == "" || endDate == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("startDate and endDate are required"))
		return
	}

	// Validate date formats
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid startDate format, use YYYY-MM-DD"))
		return
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid endDate format, use YYYY-MM-DD"))
		return
	}

	records, err := h.store.GetSensorUsageByDate(sensorId, startDate, endDate)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Calculate total for the period
	var total float64
	for _, record := range records {
		total += record.VolumeLiters
	}

	response := map[string]interface{}{
		"sensorId": sensorId,
		"startDate": startDate,
		"endDate": endDate,
		"totalVolume": total,
		"records": records,
		"recordCount": len(records),
	}

	utils.WriteJson(w, http.StatusOK, response)
}

// Helper function to validate sensor ID format
func isValidSensorId(sensorId string) bool {
	// Check if sensor ID matches expected format (ZC-WS-YYYY-NNNN)
	matched, _ := regexp.MatchString(`^ZC-WS-\d{4}-\d{4}$`, sensorId)
	return matched
}
