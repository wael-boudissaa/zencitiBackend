package sensors

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wael-boudissaa/zencitiBackend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// RegisterSensor registers a new sensor to a client or updates existing one to active
func (s *Store) RegisterSensor(sensorId, clientId string) error {
	// Check if sensor already exists and is linked to another client
	var existingClientId string
	err := s.db.QueryRow("SELECT idClient FROM waterSensor WHERE idSensor = ?", sensorId).Scan(&existingClientId)

	if err == nil && existingClientId != clientId {
		return fmt.Errorf("sensor already registered to another client")
	}

	if err == sql.ErrNoRows {
		// Sensor doesn't exist, create new record with ACTIVE status
		_, err = s.db.Exec(
			"INSERT INTO waterSensor (idSensor, idClient, status) VALUES (?, ?, 'active')",
			sensorId, clientId,
		)
		return err
	}

	// If sensor exists but same client, update to active (reactivation)
	_, err = s.db.Exec(
		"UPDATE waterSensor SET status = 'active' WHERE idSensor = ? AND idClient = ?",
		sensorId, clientId,
	)
	return err
}

// GetSensorsByClient returns all sensors registered to a client
func (s *Store) GetSensorsByClient(clientId string) ([]types.SensorInfo, error) {
	query := `
		SELECT idSensor, idClient, status
		FROM waterSensor 
		WHERE idClient = ? AND status = 'active'
	`

	rows, err := s.db.Query(query, clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []types.SensorInfo
	for rows.Next() {
		var sensor types.SensorInfo

		err := rows.Scan(&sensor.IdSensor, &sensor.IdClient, &sensor.Status)
		if err != nil {
			log.Printf("Error scanning sensor row: %v", err)
			continue
		}

		sensors = append(sensors, sensor)
	}

	return sensors, nil
}

// CheckSensorOwnership checks if a sensor exists and returns the owner's clientId
func (s *Store) CheckSensorOwnership(sensorId string) (string, error) {
	var clientId string
	err := s.db.QueryRow("SELECT idClient FROM waterSensor WHERE idSensor = ?", sensorId).Scan(&clientId)

	if err == sql.ErrNoRows {
		return "", nil // Sensor not registered
	}

	return clientId, err
}

// SaveDailyUsage saves daily water usage data
func (s *Store) SaveDailyUsage(usage types.DailyUsageData) error {
	// Use INSERT ... ON DUPLICATE KEY UPDATE to handle duplicates
	query := `
		INSERT INTO dailyWaterUsage (idSensor, usageDate, volumeLiters) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE volumeLiters = VALUES(volumeLiters)
	`

	_, err := s.db.Exec(query, usage.SensorId, usage.UsageDate, usage.VolumeLiters)
	if err != nil {
		log.Printf("Error saving daily usage for sensor %s: %v", usage.SensorId, err)
	}

	return err
}

// SaveBatchUsage saves multiple daily usage records at once
func (s *Store) SaveBatchUsage(batchData types.BatchUsageData) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO dailyWaterUsage (idSensor, usageDate, volumeLiters) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE volumeLiters = VALUES(volumeLiters)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, usage := range batchData.UsageData {
		_, err := stmt.Exec(batchData.SensorId, usage.UsageDate, usage.VolumeLiters)
		if err != nil {
			log.Printf("Error saving batch usage for sensor %s, date %s: %v",
				batchData.SensorId, usage.UsageDate, err)
			return err
		}
	}

	return tx.Commit()
}

// GetSensorUsage returns usage data for all sensors of a client with analytics
func (s *Store) GetSensorUsage(clientId, period string) (*types.SensorUsageResponse, error) {
	// Get client's sensors
	sensors, err := s.GetSensorsByClient(clientId)
	if err != nil {
		return nil, err
	}

	var response types.SensorUsageResponse

	for _, sensor := range sensors {
		usageDetails, err := s.getSensorUsageDetails(sensor.IdSensor, period)
		if err != nil {
			log.Printf("Error getting usage details for sensor %s: %v", sensor.IdSensor, err)
			continue
		}

		response.Sensors = append(response.Sensors, *usageDetails)
	}

	return &response, nil
}

// getSensorUsageDetails gets detailed usage data for a specific sensor
func (s *Store) getSensorUsageDetails(sensorId, period string) (*types.SensorUsageDetails, error) {
	var dateFilter string
	switch period {
	case "week":
		dateFilter = "WHERE usageDate >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)"
	case "month":
		dateFilter = "WHERE usageDate >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)"
	case "year":
		dateFilter = "WHERE usageDate >= DATE_SUB(CURDATE(), INTERVAL 365 DAY)"
	default:
		dateFilter = "WHERE usageDate >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)" // Default to month
	}

	// Get daily usage records
	query := fmt.Sprintf(`
		SELECT usageDate, volumeLiters 
		FROM dailyWaterUsage 
		WHERE idSensor = ? %s
		ORDER BY usageDate DESC
	`, dateFilter)

	rows, err := s.db.Query(query, sensorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dailyUsage []types.DailyUsageRecord
	var totalVolume float64

	for rows.Next() {
		var record types.DailyUsageRecord
		err := rows.Scan(&record.Date, &record.VolumeLiters)
		if err != nil {
			continue
		}

		dailyUsage = append(dailyUsage, record)
		totalVolume += record.VolumeLiters
	}

	// Calculate analytics
	avgDaily := float64(0)
	if len(dailyUsage) > 0 {
		avgDaily = totalVolume / float64(len(dailyUsage))
	}

	weeklyTotal := s.calculateWeeklyTotal(sensorId)
	monthlyTotal := s.calculateMonthlyTotal(sensorId)

	return &types.SensorUsageDetails{
		IdSensor:     sensorId,
		Status:       "active",
		DailyUsage:   dailyUsage,
		WeeklyTotal:  weeklyTotal,
		MonthlyTotal: monthlyTotal,
		AverageDaily: avgDaily,
	}, nil
}

// GetSensorUsageByDate returns usage data for a sensor within a date range
func (s *Store) GetSensorUsageByDate(sensorId, startDate, endDate string) ([]types.DailyUsageRecord, error) {
	query := `
		SELECT usageDate, volumeLiters 
		FROM dailyWaterUsage 
		WHERE idSensor = ? AND usageDate BETWEEN ? AND ?
		ORDER BY usageDate ASC
	`

	rows, err := s.db.Query(query, sensorId, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []types.DailyUsageRecord
	for rows.Next() {
		var record types.DailyUsageRecord
		err := rows.Scan(&record.Date, &record.VolumeLiters)
		if err != nil {
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

// UpdateSensorStatus updates the status of a sensor
func (s *Store) UpdateSensorStatus(sensorId, status string) error {
	_, err := s.db.Exec("UPDATE waterSensor SET status = ? WHERE idSensor = ?", status, sensorId)
	return err
}

// GetSensorInfo returns information about a specific sensor
func (s *Store) GetSensorInfo(sensorId string) (*types.SensorInfo, error) {
	var sensor types.SensorInfo

	err := s.db.QueryRow(`
		SELECT idSensor, idClient, status  
		FROM waterSensor WHERE idSensor = ?
	`, sensorId).Scan(&sensor.IdSensor, &sensor.IdClient, &sensor.Status)
	if err != nil {
		return nil, err
	}

	return &sensor, nil
}

// Helper functions for analytics
func (s *Store) calculateWeeklyTotal(sensorId string) float64 {
	var total float64
	err := s.db.QueryRow(`
		SELECT COALESCE(SUM(volumeLiters), 0) 
		FROM dailyWaterUsage 
		WHERE idSensor = ? AND usageDate >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
	`, sensorId).Scan(&total)
	if err != nil {
		log.Printf("Error calculating weekly total for sensor %s: %v", sensorId, err)
		return 0
	}
	return total
}

func (s *Store) calculateMonthlyTotal(sensorId string) float64 {
	var total float64
	err := s.db.QueryRow(`
		SELECT COALESCE(SUM(volumeLiters), 0) 
		FROM dailyWaterUsage 
		WHERE idSensor = ? AND usageDate >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
	`, sensorId).Scan(&total)
	if err != nil {
		log.Printf("Error calculating monthly total for sensor %s: %v", sensorId, err)
		return 0
	}
	return total
}

// CheckClientHasSensors checks if a client has any registered sensors
func (s *Store) CheckClientHasSensors(clientId string) (bool, int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM waterSensor WHERE idClient = ? AND status = 'active'", clientId).Scan(&count)
	if err != nil {
		return false, 0, err
	}
	return count > 0, count, nil
}
