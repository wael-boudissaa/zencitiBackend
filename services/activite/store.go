package activite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/wael-boudissaa/zencitiBackend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetAllLocationsWithDistances(clientLat, clientLng float64) (*[]types.LocationItemWithDistance, error) {
	query := `
        (SELECT 
            r.idRestaurant as id,
            r.name,
            'Restaurant' as type,
            IFNULL(r.latitude, 0) as latitude,
            IFNULL(r.longitude, 0) as longitude,
            r.location as address,
            r.image as imageUrl,
            IFNULL(p.phoneNumber, '') as phoneNumber,
            (6371 * acos(cos(radians(?)) * cos(radians(r.latitude)) * 
            cos(radians(r.longitude) - radians(?)) + 
            sin(radians(?)) * sin(radians(r.latitude)))) AS distance
        FROM restaurant r
        LEFT JOIN adminRestaurant ar ON r.idAdminRestaurant = ar.idAdminRestaurant
        LEFT JOIN profile p ON ar.idProfile = p.idProfile
        WHERE r.latitude IS NOT NULL AND r.longitude IS NOT NULL 
        AND r.latitude != 0 AND r.longitude != 0)
        
        UNION ALL
        
        (SELECT 
            a.idActivity as id,
            a.nameActivity as name,
            'Activity' as type,
            IFNULL(a.latitude, 0) as latitude,
            IFNULL(a.longitude, 0) as longitude,
'No address available' as address,
            a.imageActivity as imageUrl,
            IFNULL(p.phoneNumber, 'No phone available') as phoneNumber,
            (6371 * acos(cos(radians(?)) * cos(radians(a.latitude)) * 
            cos(radians(a.longitude) - radians(?)) + 
            sin(radians(?)) * sin(radians(a.latitude)))) AS distance
        FROM activity a
        LEFT JOIN adminActivity aa ON a.idAdminActivity = aa.idAdminActivity
        LEFT JOIN profile p ON aa.idProfile = p.idProfile
        WHERE a.latitude IS NOT NULL AND a.longitude IS NOT NULL 
        AND a.latitude != 0 AND a.longitude != 0)
        
        ORDER BY distance ASC
        LIMIT 50
    `

	rows, err := s.db.Query(query, clientLat, clientLng, clientLat, clientLat, clientLng, clientLat)
	if err != nil {
		return nil, fmt.Errorf("error retrieving locations: %v", err)
	}
	defer rows.Close()

	var locations []types.LocationItemWithDistance
	for rows.Next() {
		var location types.LocationItemWithDistance
		var distance float64

		err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.Type,
			&location.Latitude,
			&location.Longitude,
			&location.Address,
			&location.ImageURL,
			&location.PhoneNumber,
			&distance,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning location row: %v", err)
		}

		location.Distance = distance
		location.DistanceFormatted = formatDistance(distance)
		locations = append(locations, location)
	}

	return &locations, nil
}

// Helper function to format distance
func formatDistance(distance float64) string {
	if distance < 1 {
		return fmt.Sprintf("%.0f m", distance*1000)
	}
	return fmt.Sprintf("%.1f km", distance)
}

//	func (s *Store) GetActivite() (*[]types.Activite, error) {
//		query := `SELECT * FROM activite`
//		rows, err := s.db.Query(query)
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close() // Ensure rows are closed to avoid memory leaks
//		var activite []types.Activite
//
//		for rows.Next() {
//			var act types.Activite
//			err = rows.Scan(
//				&act.IdActivite,
//				&act.NameActivite,
//				&act.Description,
//			)
//			if err != nil {
//				return nil, err
//			}
//			activite = append(activite, act)
//		}
//		if err := rows.Err(); err != nil {
//			return nil, err
//		}
//		return &activite, nil
//	}
func (s *Store) UpdateClientActivityStatus(idClientActivity string) error {
	// First, check if the client activity exists and get its details
	var timeActivity time.Time
	var currentStatus string
	checkQuery := `SELECT timeActivity, status FROM clientActivity WHERE idClientActivity = ?`
	err := s.db.QueryRow(checkQuery, idClientActivity).Scan(&timeActivity, &currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("client activity with ID %s not found", idClientActivity)
		}
		return fmt.Errorf("error checking client activity: %v", err)
	}

	// Check if already completed
	if currentStatus == "completed" {
		return fmt.Errorf("activity is already completed")
	}

	// Check if the current time is within 2 hours before or after the scheduled time
	now := time.Now()
	timeDiff := now.Sub(timeActivity)

	// Allow check-in 2 hours before or 2 hours after the scheduled time
	if timeDiff < -2*time.Hour || timeDiff > 2*time.Hour {
		return fmt.Errorf("activity can only be completed within 2 hours of the scheduled time")
	}

	// Update the status to completed
	updateQuery := `UPDATE clientActivity SET status = 'completed' WHERE idClientActivity = ?`
	result, err := s.db.Exec(updateQuery, idClientActivity)
	if err != nil {
		return fmt.Errorf("error updating client activity status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}

func (s *Store) GetAllClientActivities(idClient string) ([]types.ClientActivityInfo, error) {
	query := `
        SELECT 
            ca.idClientActivity,
            ca.timeActivity,
            ca.status,
            a.nameActivity,
            a.imageActivity,
            a.descriptionActivity
        FROM clientActivity ca
        JOIN activity a ON ca.idActivity = a.idActivity
        WHERE ca.idClient = ?
        ORDER BY ca.timeActivity DESC
    `
	rows, err := s.db.Query(query, idClient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []types.ClientActivityInfo
	for rows.Next() {
		var activity types.ClientActivityInfo
		if err := rows.Scan(
			&activity.IdClientActivity,
			&activity.TimeActivity,
			&activity.Status,
			&activity.ActivityName,
			&activity.ActivityImage,
			&activity.ActivityDescription,
		); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func (s *Store) CreateActivityClient(idClientActivity string, act types.ActivityCreation) error {
	query := `INSERT INTO clientActivity (idClientActivity,idClient, idActivity, timeActivity,status) VALUES (?,?, ?, ?,?)`
	_, err := s.db.Exec(query, idClientActivity, act.IdClient, act.IdActivity, act.TimeActivity, "pending")
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetActivityNotAvaialableAtday(day time.Time, idActivity string) ([]string, error) {
	query := `
        SELECT 
            TIME(timeActivity) as reservedTime
        FROM 
            clientActivity
        WHERE
            status = 'pending'
            AND DATE(timeActivity) = ?
            AND idActivity = ?
    ;`
	rows, err := s.db.Query(query, day.Format("2006-01-02"), idActivity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservedTimes []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		reservedTimes = append(reservedTimes, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reservedTimes, nil
}

func (s *Store) GetRecentActivities(idClient string) (*[]types.ActivityProfile, error) {
	query := `SELECT
    activity.idActivity,activity.nameActivity,activity.descriptionActivity,activity.imageActivity,activity.popularity,clientActivity.timeActivity
    FROM clientActivity join activity on clientActivity.idActivity=
    activity.idActivity where clientActivity.idClient=? ORDER BY
    clientActivity.timeActivity DESC LIMIT 5 `

	rows, err := s.db.Query(query, idClient)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var activite []types.ActivityProfile
	for rows.Next() {
		var act types.ActivityProfile
		err = rows.Scan(
			&act.IdActivity,
			&act.NameActivity,
			&act.Description,
			&act.ImageActivite,
			&act.Popularity,
			&act.TimeActivity,
		)
		if err != nil {
			return nil, err
		}
		activite = append(activite, act)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &activite, nil
}

func (s *Store) GetActiviteById(id string) (*types.Activity, error) {
	query := `SELECT * FROM activity WHERE idActivity = ?`
	row := s.db.QueryRow(query, id)
	var act types.Activity
	err := row.Scan(
		&act.IdActivity,
		&act.NameActivity,
		&act.Description,
		&act.ImageActivite,
		&act.Langitude,
		&act.Latitude,
		&act.IdTypeActivity,
		&act.Popularity,
	)
	if err != nil {
		return nil, err
	}
	return &act, nil
}

func (s *Store) GetActiviteTypes() (*[]types.ActivitetType, error) {
	query := `SELECT * FROM typeActivity`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activite []types.ActivitetType

	for rows.Next() {
		var act types.ActivitetType
		err = rows.Scan(
			&act.IdActiviteType,
			&act.NameActiviteType,
			&act.ImageActivity,
		)
		if err != nil {
			return nil, err
		}
		activite = append(activite, act)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &activite, nil
}

func (s *Store) GetActivityByTypes(id string) (*[]types.Activity, error) {
	query := `SELECT * FROM activity WHERE idTypeActivity = ?`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activity []types.Activity

	for rows.Next() {
		var act types.Activity
		err = rows.Scan(
			&act.IdActivity,
			&act.NameActivity,
			&act.Description,
			&act.ImageActivite,
			&act.Langitude,
			&act.Latitude,

			&act.IdTypeActivity,
			&act.Popularity,
			&act.IdAdminActivity,
		)
		if err != nil {
			return nil, err
		}
		activity = append(activity, act)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s *Store) GetPopularActivities() (*[]types.Activity, error) {
	query := `SELECT * FROM activity ORDER BY popularity DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activite []types.Activity

	for rows.Next() {
		var act types.Activity
		err = rows.Scan(
			&act.IdActivity,
			&act.NameActivity,
			&act.Description,
			&act.ImageActivite,
			&act.IdTypeActivity,
			&act.Popularity,
		)
		if err != nil {
			return nil, err
		}
		activite = append(activite, act)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &activite, nil
}
