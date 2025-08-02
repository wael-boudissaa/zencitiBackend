package activite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
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
func (s *Store) UpdateClientActivityStatus(idClientActivity string, idAdminActivity string) error {
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
	updateQuery := `UPDATE clientActivity SET status = 'completed',idAdminActivity=? WHERE idClientActivity = ?`
	result, err := s.db.Exec(updateQuery, idAdminActivity, idClientActivity)
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
	// First get the activity's capacity
	var capacity int
	capacityQuery := `SELECT capacity FROM activity WHERE idActivity = ?`
	err := s.db.QueryRow(capacityQuery, idActivity).Scan(&capacity)
	if err != nil {
		return nil, fmt.Errorf("error getting activity capacity: %v", err)
	}

	// Get count of reservations grouped by time for the given day
	query := `
        SELECT 
            TIME(timeActivity) as reservedTime,
            COUNT(*) as reservationCount
        FROM 
            clientActivity
        WHERE
            status IN ('pending', 'completed')
            AND DATE(timeActivity) = ?
            AND idActivity = ?
        GROUP BY TIME(timeActivity)
        HAVING COUNT(*) >= ?
    `
	rows, err := s.db.Query(query, day.Format("2006-01-02"), idActivity, capacity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unavailableTimes []string
	for rows.Next() {
		var timeStr string
		var count int
		if err := rows.Scan(&timeStr, &count); err != nil {
			return nil, err
		}
		unavailableTimes = append(unavailableTimes, timeStr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return unavailableTimes, nil
}

func (s *Store) GetRecentActivities(idClient string) (*[]types.ActivityProfile, error) {
	query := `SELECT
    activity.idActivity,activity.nameActivity,activity.descriptionActivity,activity.imageActivity,activity.capacity,clientActivity.timeActivity
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
			&act.Capacity,
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

func (s *Store) GetActivityFullDetails(id string) (*types.ActivityDetails, error) {
	// 1. Get activity + admin info (admin can be null)
	query := `
        SELECT a.idActivity, a.nameActivity, a.descriptionActivity, a.imageActivity, a.longitude, a.latitude, a.idTypeActivity, a.capacity,
               p.firstName, p.lastName,aa.idAdminActivity, p.email, p.phoneNumber
        FROM activity a
        LEFT JOIN adminActivity aa ON a.idAdminActivity = aa.idAdminActivity
        LEFT JOIN profile p ON aa.idProfile = p.idProfile
        WHERE a.idActivity = ?
    `
	row := s.db.QueryRow(query, id)
	var act types.ActivityDetails
	var idAdmin, adminFirst, adminLast, adminEmail, adminPhone sql.NullString
	err := row.Scan(
		&act.IdActivity, &act.NameActivity, &act.Description, &act.ImageActivite, &act.Langitude, &act.Latitude, &act.IdTypeActivity, &act.Capacity,
		&adminFirst, &adminLast, &idAdmin, &adminEmail, &adminPhone,
	)
	if err != nil {
		return nil, err
	}
	// If admin is null, return empty strings (not null)
	if adminFirst.Valid || adminLast.Valid {
		act.AdminName = strings.TrimSpace(adminFirst.String + " " + adminLast.String)
	} else {
		act.AdminName = ""
	}
	if adminEmail.Valid {
		act.AdminEmail = adminEmail.String
	} else {
		act.AdminEmail = ""
	}
	if adminPhone.Valid {
		act.AdminPhone = adminPhone.String
	} else {
		act.AdminPhone = ""
	}
	if idAdmin.Valid {
		act.IdAdminActivity = adminPhone.String
	} else {
		act.IdAdminActivity = ""
	}

	// 2. Get rating breakdown
	ratingCounts := make(map[int]int)
	ratingQuery := `SELECT rating, COUNT(*) FROM rating WHERE idActivity = ? GROUP BY rating`
	rows, err := s.db.Query(ratingQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rating, count int
			if err := rows.Scan(&rating, &count); err == nil {
				ratingCounts[rating] = count
			}
		}
	}
	// Ensure all 1-5 are present (even if 0)
	for i := 1; i <= 5; i++ {
		if _, ok := ratingCounts[i]; !ok {
			ratingCounts[i] = 0
		}
	}
	act.RatingCounts = ratingCounts

	// 3. Get recent reviews (limit 5)
	reviewQuery := `
        SELECT p.firstName, p.lastName, r.rating, r.comment, r.createdAt
        FROM rating r
        JOIN client c ON r.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        WHERE r.idActivity = ?
        ORDER BY r.createdAt DESC
        LIMIT 5
    `
	reviewRows, err := s.db.Query(reviewQuery, id)
	if err == nil {
		defer reviewRows.Close()
		for reviewRows.Next() {
			var rev types.ActivityReviewDetail
			var first, last string
			if err := reviewRows.Scan(&first, &last, &rev.Rating, &rev.Comment, &rev.CreatedAt); err == nil {
				rev.ReviewerName = strings.TrimSpace(first + " " + last)
				act.RecentReviews = append(act.RecentReviews, rev)
			}
		}
	}
	return &act, nil
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
		&act.Capacity,
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

// CreateActivityCategory creates a new activity category and returns the generated ID
func (s *Store) CreateActivityCategory(category types.ActivityCategoryCreation) (string, error) {
	// Generate a unique ID using the standard project UUID generator
	newID, err := utils.CreateAnId()
	if err != nil {
		return "", fmt.Errorf("error generating ID for activity category: %v", err)
	}

	query := `INSERT INTO typeActivity (idTypeActivity, nameTypeActivity, imageActivity) VALUES (?, ?, ?)`
	_, err = s.db.Exec(query, newID, category.NameTypeActivity, category.ImageActivity)
	if err != nil {
		return "", fmt.Errorf("error creating activity category: %v", err)
	}

	return newID, nil
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
			&act.Capacity,
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
	query := `
		SELECT 
			a.idActivity, 
			a.nameActivity, 
			a.descriptionActivity, 
			a.imageActivity, 
			a.longitude, 
			a.latitude, 
			a.idAdminActivity, 
			a.idTypeActivity, 
			a.capacity
		FROM activity a
		LEFT JOIN rating r ON a.idActivity = r.idActivity AND r.ratingType = 'activity'
		GROUP BY a.idActivity, a.nameActivity, a.descriptionActivity, a.imageActivity, 
		         a.longitude, a.latitude, a.idAdminActivity, a.idTypeActivity, a.capacity
		ORDER BY IFNULL(AVG(r.rating), 0) DESC, a.nameActivity ASC
	`
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
			&act.Langitude,
			&act.Latitude,
			&act.IdAdminActivity,
			&act.IdTypeActivity,
			&act.Capacity,
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

func (s *Store) GetActivityStats(idActivity string) (*types.ActivityStats, error) {
	stats := &types.ActivityStats{}

	// Get total bookings
	totalBookingsQuery := `SELECT COUNT(*) FROM clientActivity WHERE idActivity = ?`
	err := s.db.QueryRow(totalBookingsQuery, idActivity).Scan(&stats.TotalBookings)
	if err != nil {
		return nil, fmt.Errorf("error getting total bookings: %v", err)
	}

	// Get bookings today
	todayQuery := `SELECT COUNT(*) FROM clientActivity WHERE idActivity = ? AND DATE(timeActivity) = CURDATE()`
	err = s.db.QueryRow(todayQuery, idActivity).Scan(&stats.BookingsToday)
	if err != nil {
		return nil, fmt.Errorf("error getting today's bookings: %v", err)
	}

	// Get bookings this week
	weekQuery := `SELECT COUNT(*) FROM clientActivity WHERE idActivity = ? AND YEARWEEK(timeActivity, 1) = YEARWEEK(CURDATE(), 1)`
	err = s.db.QueryRow(weekQuery, idActivity).Scan(&stats.BookingsThisWeek)
	if err != nil {
		return nil, fmt.Errorf("error getting this week's bookings: %v", err)
	}

	// Get bookings this month
	monthQuery := `SELECT COUNT(*) FROM clientActivity WHERE idActivity = ? AND MONTH(timeActivity) = MONTH(CURDATE()) AND YEAR(timeActivity) = YEAR(CURDATE())`
	err = s.db.QueryRow(monthQuery, idActivity).Scan(&stats.BookingsThisMonth)
	if err != nil {
		return nil, fmt.Errorf("error getting this month's bookings: %v", err)
	}

	// Get total reviews and average rating
	reviewQuery := `SELECT COUNT(*), IFNULL(AVG(rating), 0) FROM rating WHERE idActivity = ?`
	err = s.db.QueryRow(reviewQuery, idActivity).Scan(&stats.TotalReviews, &stats.AverageRating)
	if err != nil {
		return nil, fmt.Errorf("error getting reviews stats: %v", err)
	}

	// Calculate average engagement (completed bookings / total bookings * 100)
	if stats.TotalBookings > 0 {
		completedQuery := `SELECT COUNT(*) FROM clientActivity WHERE idActivity = ? AND status = 'completed'`
		var completedBookings int
		err = s.db.QueryRow(completedQuery, idActivity).Scan(&completedBookings)
		if err != nil {
			return nil, fmt.Errorf("error getting completed bookings: %v", err)
		}
		stats.AvgEngagement = (float64(completedBookings) / float64(stats.TotalBookings)) * 100
	}

	// Get daily trends (last 30 days)
	dailyTrendsQuery := `
        SELECT DATE(timeActivity) as date, COUNT(*) as bookings
        FROM clientActivity 
        WHERE idActivity = ? 
        AND timeActivity >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
        GROUP BY DATE(timeActivity)
        ORDER BY date ASC
    `
	dailyRows, err := s.db.Query(dailyTrendsQuery, idActivity)
	if err == nil {
		defer dailyRows.Close()
		for dailyRows.Next() {
			var trend types.ActivityDailyStats
			err := dailyRows.Scan(&trend.Date, &trend.Bookings)
			if err == nil {
				stats.DailyTrends = append(stats.DailyTrends, trend)
			}
		}
	}

	// Get weekly trends (last 12 weeks) - strict SQL mode compatible
	weeklyTrendsQuery := `
    SELECT 
        YEAR(timeActivity) as year,
        WEEK(timeActivity, 1) as week_num,
        COUNT(*) as bookings
    FROM clientActivity 
    WHERE idActivity = ?
    AND timeActivity >= DATE_SUB(CURDATE(), INTERVAL 12 WEEK)
    GROUP BY YEAR(timeActivity), WEEK(timeActivity, 1)
    ORDER BY year, week_num ASC
`

	weeklyRows, err := s.db.Query(weeklyTrendsQuery, idActivity)
	if err != nil {
		fmt.Printf("Weekly trends query error: %v\n", err)
	} else {
		defer weeklyRows.Close()
		for weeklyRows.Next() {
			var year, weekNum, bookings int
			err := weeklyRows.Scan(&year, &weekNum, &bookings)
			if err != nil {
				fmt.Printf("Weekly trends scan error: %v\n", err)
			} else {
				trend := types.ActivityWeeklyStats{
					Week:     fmt.Sprintf("%d-W%02d", year, weekNum),
					Bookings: bookings,
				}
				stats.WeeklyTrends = append(stats.WeeklyTrends, trend)
			}
		}
	}

	// Get monthly trends (last 12 months) - consistent with weekly trends
	monthlyTrendsQuery := `
        SELECT 
            MONTHNAME(timeActivity) as month, 
            YEAR(timeActivity) as year, 
            COUNT(*) as bookings
        FROM clientActivity 
        WHERE idActivity = ?
        AND timeActivity >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
        GROUP BY YEAR(timeActivity), MONTH(timeActivity), MONTHNAME(timeActivity)
        ORDER BY YEAR(timeActivity), MONTH(timeActivity) ASC
    `

	monthlyRows, err := s.db.Query(monthlyTrendsQuery, idActivity)
	if err != nil {
		fmt.Printf("Monthly trends query error: %v\n", err)
	} else {
		defer monthlyRows.Close()
		for monthlyRows.Next() {
			var trend types.ActivityMonthlyStats
			err := monthlyRows.Scan(&trend.Month, &trend.Year, &trend.Bookings)
			if err != nil {
				fmt.Printf("Monthly trends scan error: %v\n", err)
			} else {
				stats.MonthlyTrends = append(stats.MonthlyTrends, trend)
			}
		}
	}

	// Initialize empty slices if no data found
	if stats.DailyTrends == nil {
		stats.DailyTrends = []types.ActivityDailyStats{}
	}
	if stats.WeeklyTrends == nil {
		stats.WeeklyTrends = []types.ActivityWeeklyStats{}
	}
	if stats.MonthlyTrends == nil {
		stats.MonthlyTrends = []types.ActivityMonthlyStats{}
	}

	// Get recent bookings (last 10)
	recentBookingsQuery := `
        SELECT CONCAT(p.firstName, ' ', p.lastName) as clientName, ca.timeActivity, ca.status, ca.timeActivity as createdAt
        FROM clientActivity ca
        JOIN client c ON ca.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        WHERE ca.idActivity = ?
        ORDER BY ca.timeActivity DESC
        LIMIT 10
    `
	recentRows, err := s.db.Query(recentBookingsQuery, idActivity)
	if err == nil {
		defer recentRows.Close()
		for recentRows.Next() {
			var booking types.ActivityBookingInfo
			err := recentRows.Scan(&booking.ClientName, &booking.BookingTime, &booking.Status, &booking.CreatedAt)
			if err == nil {
				stats.RecentBookings = append(stats.RecentBookings, booking)
			}
		}
	}

	// Get top rated reviews (5 star reviews, limit 5)
	topReviewsQuery := `
        SELECT CONCAT(p.firstName, ' ', p.lastName) as reviewerName, r.rating, r.comment, r.createdAt
        FROM rating r
        JOIN client c ON r.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        WHERE r.idActivity = ? AND r.rating = 5
        ORDER BY r.createdAt DESC
        LIMIT 5
    `
	reviewRows, err := s.db.Query(topReviewsQuery, idActivity)
	if err == nil {
		defer reviewRows.Close()
		for reviewRows.Next() {
			var review types.ActivityReviewDetail
			err := reviewRows.Scan(&review.ReviewerName, &review.Rating, &review.Comment, &review.CreatedAt)
			if err == nil {
				stats.TopRatedReviews = append(stats.TopRatedReviews, review)
			}
		}
	}

	return stats, nil
}

func (s *Store) GetActivitiesByAdminActivity(idAdminActivity string) ([]types.Activity, error) {
	query := `SELECT idActivity, nameActivity, descriptionActivity, imageActivity, longitude, latitude, idAdminActivity, idTypeActivity, capacity FROM activity WHERE idAdminActivity = ?`
	rows, err := s.db.Query(query, idAdminActivity)
	if err != nil {
		return nil, fmt.Errorf("error retrieving activities: %v", err)
	}
	defer rows.Close()

	var activities []types.Activity
	for rows.Next() {
		var activity types.Activity
		err := rows.Scan(
			&activity.IdActivity,
			&activity.NameActivity,
			&activity.Description,
			&activity.ImageActivite,
			&activity.Langitude,
			&activity.Latitude,
			&activity.IdAdminActivity,
			&activity.IdTypeActivity,
			&activity.Capacity,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning activity row: %v", err)
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (s *Store) UpdateActivityStatus(idClientActivity string, status string) error {
	// First, get current activity status and time
	var currentStatus string
	var timeActivity time.Time
	query := `SELECT status, timeActivity FROM clientActivity WHERE idClientActivity = ?`
	err := s.db.QueryRow(query, idClientActivity).Scan(&currentStatus, &timeActivity)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no activity booking found with ID %s", idClientActivity)
		}
		return fmt.Errorf("error fetching activity: %v", err)
	}

	// Validate status transition
	if !isValidActivityStatusTransition(currentStatus, status) {
		return fmt.Errorf("invalid status transition from %s to %s", currentStatus, status)
	}

	// Validate time window (2 hours before or after activity time)
	now := time.Now()
	timeDiff := timeActivity.Sub(now)
	absTimeDiff := timeDiff
	if absTimeDiff < 0 {
		absTimeDiff = -absTimeDiff
	}

	twoHours := 2 * time.Hour
	if absTimeDiff > twoHours {
		return fmt.Errorf("status can only be changed within 2 hours before or after the activity time")
	}

	// If validation passes, update the status
	updateQuery := `UPDATE clientActivity SET status = ? WHERE idClientActivity = ?`
	result, err := s.db.Exec(updateQuery, status, idClientActivity)
	if err != nil {
		return fmt.Errorf("error updating activity status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no activity booking found with ID %s", idClientActivity)
	}

	return nil
}

// Helper function to validate activity status transitions
func isValidActivityStatusTransition(currentStatus, newStatus string) bool {
	// Valid transitions for activities (pending/cancelled/completed):
	// pending → cancelled ✅
	// pending → completed ✅
	// cancelled → any ❌ (blocked)
	// completed → any ❌ (blocked)

	switch currentStatus {
	case "pending":
		return newStatus == "cancelled" || newStatus == "completed"
	case "cancelled":
		return false // cancelled cannot change to any other status
	case "completed":
		return false // completed cannot change to any other status
	default:
		return false
	}
}

func (s *Store) GetActivityStatsAdmin(idAdminActivity string) (*types.ActivityStats, error) {
	// First, get all activities managed by this admin
	activities, err := s.GetActivitiesByAdminActivity(idAdminActivity)
	if err != nil {
		return nil, fmt.Errorf("error getting admin activities: %v", err)
	}

	if len(activities) == 0 {
		return &types.ActivityStats{
			TotalBookings:     0,
			CompletedBookings: 0,
			PendingBookings:   0,
			CancelledBookings: 0,
			AvgEngagement:     0,
			TotalReviews:      0,
			AverageRating:     0,
			BookingsToday:     0,
			BookingsThisWeek:  0,
			BookingsThisMonth: 0,
			DailyTrends:       []types.ActivityDailyStats{},
			WeeklyTrends:      []types.ActivityWeeklyStats{},
			MonthlyTrends:     []types.ActivityMonthlyStats{},
			RecentBookings:    []types.ActivityBookingInfo{},
			TopRatedReviews:   []types.ActivityReviewDetail{},
		}, nil
	}

	stats := &types.ActivityStats{}

	// Create a comma-separated list of activity IDs for queries
	activityIDs := make([]string, len(activities))
	for i, activity := range activities {
		activityIDs[i] = activity.IdActivity
	}

	// Create placeholders for IN clause
	placeholders := strings.Repeat("?,", len(activityIDs))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

	// Convert to interface{} slice for query parameters
	args := make([]interface{}, len(activityIDs))
	for i, id := range activityIDs {
		args[i] = id
	}

	// Get total bookings (ALL bookings regardless of status)
	totalBookingsQuery := fmt.Sprintf(`SELECT COUNT(*) FROM clientActivity WHERE idActivity IN (%s)`, placeholders)
	err = s.db.QueryRow(totalBookingsQuery, args...).Scan(&stats.TotalBookings)
	if err != nil {
		return nil, fmt.Errorf("error getting total bookings: %v", err)
	}

	// Get completed bookings
	var completedBookings int
	completedQuery := fmt.Sprintf(`SELECT COUNT(*) FROM clientActivity WHERE idActivity IN (%s) AND status = 'completed'`, placeholders)
	err = s.db.QueryRow(completedQuery, args...).Scan(&completedBookings)
	if err != nil {
		return nil, fmt.Errorf("error getting completed bookings: %v", err)
	}

	// Get pending bookings
	var pendingBookings int
	pendingQuery := fmt.Sprintf(`SELECT COUNT(*) FROM clientActivity WHERE idActivity IN (%s) AND status = 'pending'`, placeholders)
	err = s.db.QueryRow(pendingQuery, args...).Scan(&pendingBookings)
	if err != nil {
		return nil, fmt.Errorf("error getting pending bookings: %v", err)
	}

	// Get cancelled bookings
	var cancelledBookings int
	cancelledQuery := fmt.Sprintf(`SELECT COUNT(*) FROM clientActivity WHERE idActivity IN (%s) AND status = 'cancelled'`, placeholders)
	err = s.db.QueryRow(cancelledQuery, args...).Scan(&cancelledBookings)
	if err != nil {
		return nil, fmt.Errorf("error getting cancelled bookings: %v", err)
	}

	// Get bookings today (all statuses)
	todayQuery := fmt.Sprintf(`SELECT COUNT(*) FROM clientActivity WHERE idActivity IN (%s) AND DATE(timeActivity) = CURDATE()`, placeholders)
	err = s.db.QueryRow(todayQuery, args...).Scan(&stats.BookingsToday)
	if err != nil {
		return nil, fmt.Errorf("error getting today's bookings: %v", err)
	}

	// Get bookings this week (all statuses)
	weekQuery := fmt.Sprintf(`SELECT COUNT(*) FROM clientActivity WHERE idActivity IN (%s) AND YEARWEEK(timeActivity, 1) = YEARWEEK(CURDATE(), 1)`, placeholders)
	err = s.db.QueryRow(weekQuery, args...).Scan(&stats.BookingsThisWeek)
	if err != nil {
		return nil, fmt.Errorf("error getting this week's bookings: %v", err)
	}

	// Get bookings this month (all statuses)
	monthQuery := fmt.Sprintf(`SELECT COUNT(*) FROM clientActivity WHERE idActivity IN (%s) AND MONTH(timeActivity) = MONTH(CURDATE()) AND YEAR(timeActivity) = YEAR(CURDATE())`, placeholders)
	err = s.db.QueryRow(monthQuery, args...).Scan(&stats.BookingsThisMonth)
	if err != nil {
		return nil, fmt.Errorf("error getting this month's bookings: %v", err)
	}

	// Get total reviews and average rating (adding ratingType filter)
	reviewQuery := fmt.Sprintf(`SELECT COUNT(*), IFNULL(AVG(rating), 0) FROM rating WHERE idActivity IN (%s) AND ratingType = 'activity'`, placeholders)
	err = s.db.QueryRow(reviewQuery, args...).Scan(&stats.TotalReviews, &stats.AverageRating)
	if err != nil {
		return nil, fmt.Errorf("error getting reviews stats: %v", err)
	}

	// Calculate average engagement (completed bookings / total bookings * 100)
	if stats.TotalBookings > 0 {
		stats.AvgEngagement = (float64(completedBookings) / float64(stats.TotalBookings)) * 100
	}

	// Get daily trends (last 30 days) - all statuses
	dailyTrendsQuery := fmt.Sprintf(`
        SELECT DATE(timeActivity) as date, COUNT(*) as bookings
        FROM clientActivity 
        WHERE idActivity IN (%s)
        AND timeActivity >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
        GROUP BY DATE(timeActivity)
        ORDER BY date ASC
    `, placeholders)
	dailyRows, err := s.db.Query(dailyTrendsQuery, args...)
	if err == nil {
		defer dailyRows.Close()
		for dailyRows.Next() {
			var trend types.ActivityDailyStats
			err := dailyRows.Scan(&trend.Date, &trend.Bookings)
			if err == nil {
				stats.DailyTrends = append(stats.DailyTrends, trend)
			}
		}
	}

	// Get weekly trends (last 12 weeks) - strict SQL mode compatible
	weeklyTrendsQuery := fmt.Sprintf(`
    SELECT 
        YEAR(timeActivity) as year,
        WEEK(timeActivity, 1) as week_num,
        COUNT(*) as bookings
    FROM clientActivity 
    WHERE idActivity IN (%s)
    AND timeActivity >= DATE_SUB(CURDATE(), INTERVAL 12 WEEK)
    GROUP BY YEAR(timeActivity), WEEK(timeActivity, 1)
    ORDER BY year, week_num ASC
`, placeholders)

	weeklyRows, err := s.db.Query(weeklyTrendsQuery, args...)
	if err != nil {
		fmt.Printf("Weekly trends query error: %v\n", err)
	} else {
		defer weeklyRows.Close()
		for weeklyRows.Next() {
			var year, weekNum, bookings int
			err := weeklyRows.Scan(&year, &weekNum, &bookings)
			if err != nil {
				fmt.Printf("Weekly trends scan error: %v\n", err)
			} else {
				trend := types.ActivityWeeklyStats{
					Week:     fmt.Sprintf("%d-W%02d", year, weekNum),
					Bookings: bookings,
				}
				stats.WeeklyTrends = append(stats.WeeklyTrends, trend)
			}
		}
	}

	// Get monthly trends (last 12 months) - consistent with weekly trends
	monthlyTrendsQuery := fmt.Sprintf(`
    SELECT 
        MONTHNAME(DATE(timeActivity)) as month, 
        YEAR(DATE(timeActivity)) as year, 
        COUNT(*) as bookings
    FROM clientActivity 
    WHERE idActivity IN (%s)
    AND DATE(timeActivity) >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
    GROUP BY YEAR(DATE(timeActivity)), MONTH(DATE(timeActivity)), MONTHNAME(DATE(timeActivity))
    ORDER BY YEAR(DATE(timeActivity)), MONTH(DATE(timeActivity)) ASC
`, placeholders)

	monthlyRows, err := s.db.Query(monthlyTrendsQuery, args...)
	if err != nil {
		fmt.Printf("Monthly trends query error: %v\n", err)
	} else {
		defer monthlyRows.Close()
		for monthlyRows.Next() {
			var trend types.ActivityMonthlyStats
			err := monthlyRows.Scan(&trend.Month, &trend.Year, &trend.Bookings)
			if err != nil {
				fmt.Printf("Monthly trends scan error: %v\n", err)
			} else {
				stats.MonthlyTrends = append(stats.MonthlyTrends, trend)
			}
		}
	}

	// Initialize empty slices if no data found
	if stats.WeeklyTrends == nil {
		stats.WeeklyTrends = []types.ActivityWeeklyStats{}
	}
	if stats.MonthlyTrends == nil {
		stats.MonthlyTrends = []types.ActivityMonthlyStats{}
	}

	// Get recent bookings (last 10) with activity names and all statuses
	recentBookingsQuery := fmt.Sprintf(`
        SELECT CONCAT(p.firstName, ' ', p.lastName) as clientName, ca.timeActivity, ca.status, ca.timeActivity as createdAt, a.nameActivity
        FROM clientActivity ca
        JOIN client c ON ca.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        JOIN activity a ON ca.idActivity = a.idActivity
        WHERE ca.idActivity IN (%s)
        ORDER BY ca.timeActivity DESC
        LIMIT 10
    `, placeholders)
	recentRows, err := s.db.Query(recentBookingsQuery, args...)
	if err == nil {
		defer recentRows.Close()
		for recentRows.Next() {
			var booking types.ActivityBookingInfo
			var activityName string
			err := recentRows.Scan(&booking.ClientName, &booking.BookingTime, &booking.Status, &booking.CreatedAt, &activityName)
			if err == nil {
				// Add activity name and status info for context
				booking.ClientName = fmt.Sprintf("%s (%s) - %s", booking.ClientName, activityName, strings.ToUpper(booking.Status))
				stats.RecentBookings = append(stats.RecentBookings, booking)
			}
		}
	}

	// Get top rated reviews (4 and 5 star reviews, limit 5) with activity names
	topReviewsQuery := fmt.Sprintf(`
        SELECT CONCAT(p.firstName, ' ', p.lastName) as reviewerName, r.rating, r.comment, r.createdAt, a.nameActivity
        FROM rating r
        JOIN client c ON r.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        JOIN activity a ON r.idActivity = a.idActivity
        WHERE r.idActivity IN (%s) AND r.rating >= 4 AND r.ratingType = 'activity'
        ORDER BY r.rating DESC, r.createdAt DESC
        LIMIT 5
    `, placeholders)
	reviewRows, err := s.db.Query(topReviewsQuery, args...)
	if err == nil {
		defer reviewRows.Close()
		for reviewRows.Next() {
			var review types.ActivityReviewDetail
			var activityName string
			err := reviewRows.Scan(&review.ReviewerName, &review.Rating, &review.Comment, &review.CreatedAt, &activityName)
			if err == nil {
				// Add activity name to comment for context
				review.Comment = fmt.Sprintf("[%s] %s", activityName, review.Comment)
				stats.TopRatedReviews = append(stats.TopRatedReviews, review)
			}
		}
	}

	// Add the status breakdown to stats (no confirmed bookings)
	stats.CompletedBookings = completedBookings
	stats.PendingBookings = pendingBookings
	stats.CancelledBookings = cancelledBookings

	return stats, nil
}

func (s *Store) PostRatingActivity(rating types.PostRatingActivity) error {
	query := `INSERT INTO rating (idRating, idClient, idActivity, ratingType, rating, comment, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, rating.IdRating, rating.IdClient, rating.IdActivity, "activity", rating.RatingValue, rating.Comment, time.Now())
	if err != nil {
		return fmt.Errorf("error inserting activity rating: %v", err)
	}
	return nil
}

func (s *Store) GetAllCampusFacilities() (*types.CampusFacilitiesResponse, error) {
	activities, err := s.getAllActivitiesWithDetails()
	if err != nil {
		return nil, fmt.Errorf("error fetching activities: %v", err)
	}

	restaurants, err := s.getAllRestaurantsWithDetails()
	if err != nil {
		return nil, fmt.Errorf("error fetching restaurants: %v", err)
	}

	return &types.CampusFacilitiesResponse{
		Activities:      activities,
		Restaurants:     restaurants,
		Total:           len(activities) + len(restaurants),
		ActivityCount:   len(activities),
		RestaurantCount: len(restaurants),
	}, nil
}

func (s *Store) getAllActivitiesWithDetails() ([]types.CampusFacilityItem, error) {
	query := `
		SELECT 
			a.idActivity,
			a.nameActivity,
			a.descriptionActivity,
			a.imageActivity,
			a.latitude,
			a.longitude,
			a.capacity,
			a.idTypeActivity,
			ta.nameTypeActivity,
			aa.idAdminActivity,
			CONCAT(p.firstName, ' ', p.lastName) as adminName,
			p.email as adminEmail,
			CASE 
				WHEN aa.idAdminActivity IS NULL THEN 'inactive'
				ELSE 'active'
			END as adminStatus
		FROM activity a
		LEFT JOIN adminActivity aa ON a.idAdminActivity = aa.idAdminActivity
		LEFT JOIN profile p ON aa.idProfile = p.idProfile
		LEFT JOIN typeActivity ta ON a.idTypeActivity = ta.idTypeActivity
		ORDER BY a.nameActivity
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying activities: %v", err)
	}
	defer rows.Close()

	var activities []types.CampusFacilityItem
	for rows.Next() {
		var activity types.CampusFacilityItem
		var description sql.NullString
		var image sql.NullString
		var latitude float64
		var longitude float64
		var capacity int
		var categoryID sql.NullString
		var categoryName sql.NullString
		var adminID sql.NullString
		var adminName sql.NullString
		var adminEmail sql.NullString
		var adminStatus string

		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&description,
			&image,
			&latitude,
			&longitude,
			&capacity,
			&categoryID,
			&categoryName,
			&adminID,
			&adminName,
			&adminEmail,
			&adminStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning activity row: %v", err)
		}

		activity.Type = "activity"
		
		if description.Valid {
			activity.Description = &description.String
		}
		if image.Valid {
			activity.Image = &image.String
		}
		// Always populate coordinates and capacity since they have default values
		activity.Latitude = &latitude
		activity.Longitude = &longitude
		activity.Capacity = &capacity
		
		if categoryID.Valid {
			activity.CategoryID = &categoryID.String
		}
		if categoryName.Valid {
			activity.CategoryName = &categoryName.String
		}
		if adminID.Valid {
			activity.AdminID = &adminID.String
		}
		if adminName.Valid {
			activity.AdminName = &adminName.String
		}
		if adminEmail.Valid {
			activity.AdminEmail = &adminEmail.String
		}
		activity.AdminStatus = &adminStatus

		activities = append(activities, activity)
	}

	return activities, nil
}

func (s *Store) getAllRestaurantsWithDetails() ([]types.CampusFacilityItem, error) {
	query := `
		SELECT 
			r.idRestaurant,
			r.name,
			r.description,
			r.image,
			r.latitude,
			r.longitude,
			r.capacity,
			r.location,
			ar.idAdminRestaurant,
			CONCAT(p.firstName, ' ', p.lastName) as adminName,
			p.email as adminEmail,
			CASE 
				WHEN ar.idAdminRestaurant IS NULL THEN 'inactive'
				ELSE 'active'
			END as adminStatus
		FROM restaurant r
		LEFT JOIN adminRestaurant ar ON r.idAdminRestaurant = ar.idAdminRestaurant
		LEFT JOIN profile p ON ar.idProfile = p.idProfile
		ORDER BY r.name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying restaurants: %v", err)
	}
	defer rows.Close()

	var restaurants []types.CampusFacilityItem
	for rows.Next() {
		var restaurant types.CampusFacilityItem
		var description sql.NullString
		var image sql.NullString
		var latitude float64
		var longitude float64
		var capacity int
		var location sql.NullString
		var adminID sql.NullString
		var adminName sql.NullString
		var adminEmail sql.NullString
		var adminStatus string

		err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&description,
			&image,
			&latitude,
			&longitude,
			&capacity,
			&location,
			&adminID,
			&adminName,
			&adminEmail,
			&adminStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning restaurant row: %v", err)
		}

		restaurant.Type = "restaurant"
		
		if description.Valid {
			restaurant.Description = &description.String
		}
		if image.Valid {
			restaurant.Image = &image.String
		}
		// Always populate coordinates and capacity since they have default values
		restaurant.Latitude = &latitude
		restaurant.Longitude = &longitude
		restaurant.Capacity = &capacity
		
		if location.Valid {
			restaurant.Location = &location.String
		}
		if adminID.Valid {
			restaurant.AdminID = &adminID.String
		}
		if adminName.Valid {
			restaurant.AdminName = &adminName.String
		}
		if adminEmail.Valid {
			restaurant.AdminEmail = &adminEmail.String
		}
		restaurant.AdminStatus = &adminStatus

		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

// GetActivityBookings returns all bookings for a specific activity with full client details
func (s *Store) GetActivityBookings(idActivity string) ([]types.ActivityBookingDetail, error) {
	query := `
		SELECT 
			ca.idClientActivity,
			CONCAT(p.firstName, ' ', p.lastName) as clientName,
			p.email as clientEmail,
			p.phoneNumber as clientPhone,
			c.username as clientUsername,
			ca.timeActivity as bookingTime,
			ca.status,
			ca.timeActivity as createdAt
		FROM clientActivity ca
		JOIN client c ON ca.idClient = c.idClient  
		JOIN profile p ON c.idProfile = p.idProfile
		WHERE ca.idActivity = ?
		ORDER BY ca.timeActivity DESC
	`
	
	rows, err := s.db.Query(query, idActivity)
	if err != nil {
		return nil, fmt.Errorf("error querying activity bookings: %v", err)
	}
	defer rows.Close()

	var bookings []types.ActivityBookingDetail
	for rows.Next() {
		var booking types.ActivityBookingDetail
		var clientPhone sql.NullString
		
		err := rows.Scan(
			&booking.IdClientActivity,
			&booking.ClientName,
			&booking.ClientEmail,
			&clientPhone,
			&booking.ClientUsername,
			&booking.BookingTime,
			&booking.Status,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning booking row: %v", err)
		}
		
		if clientPhone.Valid {
			booking.ClientPhone = clientPhone.String
		}
		
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// GetAdminActivityBookings returns all bookings for activities managed by an admin
func (s *Store) GetAdminActivityBookings(idAdminActivity string) ([]types.ActivityBookingDetail, error) {
	query := `
		SELECT 
			ca.idClientActivity,
			CONCAT(p.firstName, ' ', p.lastName) as clientName,
			p.email as clientEmail,
			p.phoneNumber as clientPhone,
			c.username as clientUsername,
			ca.timeActivity as bookingTime,
			ca.status,
			ca.timeActivity as createdAt,
			a.nameActivity
		FROM clientActivity ca
		JOIN client c ON ca.idClient = c.idClient  
		JOIN profile p ON c.idProfile = p.idProfile
		JOIN activity a ON ca.idActivity = a.idActivity
		WHERE a.idAdminActivity = ?
		ORDER BY ca.timeActivity DESC
	`
	
	rows, err := s.db.Query(query, idAdminActivity)
	if err != nil {
		return nil, fmt.Errorf("error querying admin activity bookings: %v", err)
	}
	defer rows.Close()

	var bookings []types.ActivityBookingDetail
	for rows.Next() {
		var booking types.ActivityBookingDetail
		var clientPhone sql.NullString
		var activityName string
		
		err := rows.Scan(
			&booking.IdClientActivity,
			&booking.ClientName,
			&booking.ClientEmail,
			&clientPhone,
			&booking.ClientUsername,
			&booking.BookingTime,
			&booking.Status,
			&booking.CreatedAt,
			&activityName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning admin booking row: %v", err)
		}
		
		if clientPhone.Valid {
			booking.ClientPhone = clientPhone.String
		}
		
		// Add activity name to client name for context
		booking.ClientName = fmt.Sprintf("%s (%s)", booking.ClientName, activityName)
		
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// GetActivityDetailedAnalytics returns comprehensive analytics for a specific activity
func (s *Store) GetActivityDetailedAnalytics(idActivity string) (*types.ActivityDetailedAnalytics, error) {
	analytics := &types.ActivityDetailedAnalytics{}

	// Get basic stats
	totalBookingsQuery := `SELECT COUNT(*) FROM clientActivity WHERE idActivity = ?`
	err := s.db.QueryRow(totalBookingsQuery, idActivity).Scan(&analytics.TotalBookings)
	if err != nil {
		return nil, fmt.Errorf("error getting total bookings: %v", err)
	}

	// Get bookings by status
	completedQuery := `SELECT COUNT(*) FROM clientActivity WHERE idActivity = ? AND status = 'completed'`
	err = s.db.QueryRow(completedQuery, idActivity).Scan(&analytics.CompletedBookings)
	if err != nil {
		return nil, fmt.Errorf("error getting completed bookings: %v", err)
	}

	pendingQuery := `SELECT COUNT(*) FROM clientActivity WHERE idActivity = ? AND status = 'pending'`
	err = s.db.QueryRow(pendingQuery, idActivity).Scan(&analytics.PendingBookings)
	if err != nil {
		return nil, fmt.Errorf("error getting pending bookings: %v", err)
	}

	cancelledQuery := `SELECT COUNT(*) FROM clientActivity WHERE idActivity = ? AND status = 'cancelled'`
	err = s.db.QueryRow(cancelledQuery, idActivity).Scan(&analytics.CancelledBookings)
	if err != nil {
		return nil, fmt.Errorf("error getting cancelled bookings: %v", err)
	}

	// Calculate completion rate
	if analytics.TotalBookings > 0 {
		analytics.CompletionRate = (float64(analytics.CompletedBookings) / float64(analytics.TotalBookings)) * 100
	}

	// Get ratings stats
	reviewQuery := `SELECT COUNT(*), IFNULL(AVG(rating), 0) FROM rating WHERE idActivity = ? AND ratingType = 'activity'`
	err = s.db.QueryRow(reviewQuery, idActivity).Scan(&analytics.TotalReviews, &analytics.AverageRating)
	if err != nil {
		return nil, fmt.Errorf("error getting reviews stats: %v", err)
	}

	// Bookings by status map
	analytics.BookingsByStatus = map[string]int{
		"completed": analytics.CompletedBookings,
		"pending":   analytics.PendingBookings,
		"cancelled": analytics.CancelledBookings,
	}

	// Get peak hours analysis
	peakHoursQuery := `
		SELECT HOUR(timeActivity) as hour, COUNT(*) as bookings
		FROM clientActivity 
		WHERE idActivity = ?
		GROUP BY HOUR(timeActivity)
		ORDER BY hour ASC
	`
	peakRows, err := s.db.Query(peakHoursQuery, idActivity)
	if err == nil {
		defer peakRows.Close()
		for peakRows.Next() {
			var hourStat types.HourlyBookingStats
			err := peakRows.Scan(&hourStat.Hour, &hourStat.Bookings)
			if err == nil {
				analytics.PeakHours = append(analytics.PeakHours, hourStat)
			}
		}
	}

	// Get daily trends (last 30 days)
	dailyTrendsQuery := `
		SELECT DATE(timeActivity) as date, COUNT(*) as bookings
		FROM clientActivity 
		WHERE idActivity = ? 
		AND timeActivity >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
		GROUP BY DATE(timeActivity)
		ORDER BY date ASC
	`
	dailyRows, err := s.db.Query(dailyTrendsQuery, idActivity)
	if err == nil {
		defer dailyRows.Close()
		for dailyRows.Next() {
			var trend types.ActivityDailyStats
			err := dailyRows.Scan(&trend.Date, &trend.Bookings)
			if err == nil {
				analytics.DailyTrends = append(analytics.DailyTrends, trend)
			}
		}
	}

	// Get weekly trends (last 12 weeks)
	weeklyTrendsQuery := `
		SELECT 
			YEAR(timeActivity) as year,
			WEEK(timeActivity, 1) as week_num,
			COUNT(*) as bookings
		FROM clientActivity 
		WHERE idActivity = ?
		AND timeActivity >= DATE_SUB(CURDATE(), INTERVAL 12 WEEK)
		GROUP BY YEAR(timeActivity), WEEK(timeActivity, 1)
		ORDER BY year, week_num ASC
	`
	weeklyRows, err := s.db.Query(weeklyTrendsQuery, idActivity)
	if err == nil {
		defer weeklyRows.Close()
		for weeklyRows.Next() {
			var year, weekNum, bookings int
			err := weeklyRows.Scan(&year, &weekNum, &bookings)
			if err == nil {
				trend := types.ActivityWeeklyStats{
					Week:     fmt.Sprintf("%d-W%02d", year, weekNum),
					Bookings: bookings,
				}
				analytics.WeeklyTrends = append(analytics.WeeklyTrends, trend)
			}
		}
	}

	// Get monthly trends (last 12 months)
	monthlyTrendsQuery := `
		SELECT 
			MONTHNAME(timeActivity) as month, 
			YEAR(timeActivity) as year, 
			COUNT(*) as bookings
		FROM clientActivity 
		WHERE idActivity = ?
		AND timeActivity >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
		GROUP BY YEAR(timeActivity), MONTH(timeActivity), MONTHNAME(timeActivity)
		ORDER BY YEAR(timeActivity), MONTH(timeActivity) ASC
	`
	monthlyRows, err := s.db.Query(monthlyTrendsQuery, idActivity)
	if err == nil {
		defer monthlyRows.Close()
		for monthlyRows.Next() {
			var trend types.ActivityMonthlyStats
			err := monthlyRows.Scan(&trend.Month, &trend.Year, &trend.Bookings)
			if err == nil {
				analytics.MonthlyTrends = append(analytics.MonthlyTrends, trend)
			}
		}
	}

	// Calculate client return rate
	uniqueClientsQuery := `SELECT COUNT(DISTINCT idClient) FROM clientActivity WHERE idActivity = ?`
	var uniqueClients int
	err = s.db.QueryRow(uniqueClientsQuery, idActivity).Scan(&uniqueClients)
	if err == nil && analytics.TotalBookings > 0 && uniqueClients > 0 {
		analytics.ClientReturnRate = (float64(analytics.TotalBookings) / float64(uniqueClients)) - 1.0
		if analytics.ClientReturnRate < 0 {
			analytics.ClientReturnRate = 0
		}
	}

	// Calculate capacity utilization (need activity capacity)
	var capacity int
	capacityQuery := `SELECT capacity FROM activity WHERE idActivity = ?`
	err = s.db.QueryRow(capacityQuery, idActivity).Scan(&capacity)
	if err == nil && capacity > 0 {
		// This is a simplified calculation - in reality you'd need to consider time slots
		analytics.CapacityUtilization = (float64(analytics.CompletedBookings) / float64(capacity)) * 100
		if analytics.CapacityUtilization > 100 {
			analytics.CapacityUtilization = 100
		}
	}

	// Get recent bookings (last 10)
	recentBookings, err := s.GetActivityBookings(idActivity)
	if err == nil && len(recentBookings) > 10 {
		analytics.RecentBookings = recentBookings[:10]
	} else if err == nil {
		analytics.RecentBookings = recentBookings
	}

	// Get top rated reviews (4 and 5 star reviews, limit 5)
	topReviewsQuery := `
		SELECT CONCAT(p.firstName, ' ', p.lastName) as reviewerName, r.rating, r.comment, r.createdAt
		FROM rating r
		JOIN client c ON r.idClient = c.idClient
		JOIN profile p ON c.idProfile = p.idProfile
		WHERE r.idActivity = ? AND r.rating >= 4 AND r.ratingType = 'activity'
		ORDER BY r.rating DESC, r.createdAt DESC
		LIMIT 5
	`
	reviewRows, err := s.db.Query(topReviewsQuery, idActivity)
	if err == nil {
		defer reviewRows.Close()
		for reviewRows.Next() {
			var review types.ActivityReviewDetail
			err := reviewRows.Scan(&review.ReviewerName, &review.Rating, &review.Comment, &review.CreatedAt)
			if err == nil {
				analytics.TopRatedReviews = append(analytics.TopRatedReviews, review)
			}
		}
	}

	// Initialize empty slices if no data found
	if analytics.PeakHours == nil {
		analytics.PeakHours = []types.HourlyBookingStats{}
	}
	if analytics.DailyTrends == nil {
		analytics.DailyTrends = []types.ActivityDailyStats{}
	}
	if analytics.WeeklyTrends == nil {
		analytics.WeeklyTrends = []types.ActivityWeeklyStats{}
	}
	if analytics.MonthlyTrends == nil {
		analytics.MonthlyTrends = []types.ActivityMonthlyStats{}
	}
	if analytics.RecentBookings == nil {
		analytics.RecentBookings = []types.ActivityBookingDetail{}
	}
	if analytics.TopRatedReviews == nil {
		analytics.TopRatedReviews = []types.ActivityReviewDetail{}
	}

	return analytics, nil
}
