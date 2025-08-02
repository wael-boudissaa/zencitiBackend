package restaurant

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// "log"

	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{db: db}
}

// Add to services/restaurant/store.go

func (s *store) GetUniversalReservationDetails(reservationId string, reservationType string) (*types.UniversalReservationDetails, error) {
	if reservationType == "restaurant" {
		return s.getRestaurantReservationDetails(reservationId)
	} else if reservationType == "activity" {
		return s.getActivityReservationDetails(reservationId)
	} else {
		return nil, fmt.Errorf("invalid reservation type. Must be 'restaurant' or 'activity'")
	}
}

func (s *store) getRestaurantReservationDetails(reservationId string) (*types.UniversalReservationDetails, error) {
	query := `
        SELECT 
            r.idReservation,
            r.idClient,
            r.status,
            r.createdAt,
            r.timeFrom,
            r.numberOfPeople,
            r.idTable,
            
            -- Client information
            p.firstName,
            p.lastName,
            p.email,
            p.phoneNumber,
            c.username,
            
            -- Restaurant information
            rest.idRestaurant,
            rest.name,
            rest.image,
            rest.location,
            rest.description,
            rest.capacity,
            rest.longitude,
            rest.latitude,
            
            -- Restaurant admin information
            adminProfile.firstName,
            adminProfile.lastName,
            adminProfile.email,
            adminProfile.phoneNumber
            
        FROM reservation r
        JOIN client c ON r.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        JOIN restaurant rest ON r.idRestaurant = rest.idRestaurant
        JOIN adminRestaurant ar ON rest.idAdminRestaurant = ar.idAdminRestaurant
        JOIN profile adminProfile ON ar.idProfile = adminProfile.idProfile
        WHERE r.idReservation = ?
    `

	row := s.db.QueryRow(query, reservationId)

	var details types.UniversalReservationDetails
	var restaurantInfo types.RestaurantReservationInfo
	var tableID sql.NullString

	err := row.Scan(
		&details.ReservationID,
		&details.ClientID,
		&details.Status,
		&details.CreatedAt,
		&details.ReservationTime,
		&restaurantInfo.NumberOfPeople,
		&tableID,

		// Client info
		&details.ClientFirstName,
		&details.ClientLastName,
		&details.ClientEmail,
		&details.ClientPhone,
		&details.ClientUsername,

		// Restaurant info
		&restaurantInfo.RestaurantID,
		&restaurantInfo.RestaurantName,
		&restaurantInfo.RestaurantImage,
		&restaurantInfo.RestaurantLocation,
		&restaurantInfo.Description,
		&restaurantInfo.Capacity,
		&restaurantInfo.Longitude,
		&restaurantInfo.Latitude,

		// Admin info
		&restaurantInfo.AdminFirstName,
		&restaurantInfo.AdminLastName,
		&restaurantInfo.AdminEmail,
		&restaurantInfo.AdminPhone,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("restaurant reservation with ID %s not found", reservationId)
		}
		return nil, fmt.Errorf("error retrieving restaurant reservation details: %v", err)
	}

	if tableID.Valid {
		restaurantInfo.TableID = tableID.String
	}

	details.ReservationType = "restaurant"
	details.RestaurantInfo = &restaurantInfo

	return &details, nil
}

func (s *store) getActivityReservationDetails(reservationId string) (*types.UniversalReservationDetails, error) {
	query := `
        SELECT 
            ca.idClientActivity,
            ca.idClient,
            ca.status,
            ca.timeActivity,
            
            -- Client information
            p.firstName,
            p.lastName,
            p.email,
            p.phoneNumber,
            c.username,
            
            -- Activity information
            a.idActivity,
            a.nameActivity,
            a.descriptionActivity,
            a.imageActivity,
            a.capacity,
            a.longitude,
            a.latitude,
            
            -- Activity type
            ta.nameTypeActivity,
            
            -- Activity admin information
            adminProfile.firstName,
            adminProfile.lastName,
            adminProfile.email,
            adminProfile.phoneNumber
            
        FROM clientActivity ca
        JOIN client c ON ca.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        JOIN activity a ON ca.idActivity = a.idActivity
        JOIN typeActivity ta ON a.idTypeActivity = ta.idTypeActivity
        LEFT JOIN adminActivity aa ON ca.idAdminActivity = aa.idAdminActivity
        LEFT JOIN profile adminProfile ON aa.idProfile = adminProfile.idProfile
        WHERE ca.idClientActivity = ?
    `

	row := s.db.QueryRow(query, reservationId)

	var details types.UniversalReservationDetails
	var activityInfo types.ActivityReservationInfo

	err := row.Scan(
		&details.ReservationID,
		&details.ClientID,
		&details.Status,
		&details.ReservationTime,

		// Client info
		&details.ClientFirstName,
		&details.ClientLastName,
		&details.ClientEmail,
		&details.ClientPhone,
		&details.ClientUsername,

		// Activity info
		&activityInfo.ActivityID,
		&activityInfo.ActivityName,
		&activityInfo.ActivityDescription,
		&activityInfo.ActivityImage,
		&activityInfo.Capacity,
		&activityInfo.Longitude,
		&activityInfo.Latitude,

		// Activity type
		&activityInfo.ActivityType,

		// Admin info
		&activityInfo.AdminFirstName,
		&activityInfo.AdminLastName,
		&activityInfo.AdminEmail,
		&activityInfo.AdminPhone,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("activity reservation with ID %s not found", reservationId)
		}
		return nil, fmt.Errorf("error retrieving activity reservation details: %v", err)
	}

	details.ReservationType = "activity"
	details.ActivityInfo = &activityInfo
	// Set created time to reservation time for activities (since they don't have separate created time)
	details.CreatedAt = details.ReservationTime

	return &details, nil
}

func (s *store) GetReservationDetails(idReservation string) (*types.ReservationIdDetails, error) {
	reservationQuery := `
        SELECT 
            r.idReservation,
            r.idClient,
            r.timeFrom,
            r.numberOfPeople,
            r.status,
            r.createdAt,
            CONCAT(p.firstName, ' ', p.lastName) as fullName,
            p.firstName,
            p.lastName,
            p.email,
            p.phoneNumber
        FROM reservation r
        JOIN client c ON r.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        WHERE r.idReservation = ?
    `

	row := s.db.QueryRow(reservationQuery, idReservation)
	var details types.ReservationIdDetails
	var idClient string

	err := row.Scan(
		&details.IdReservation,
		&idClient,
		&details.TimeFrom,
		&details.NumberOfPeople,
		&details.Status,
		&details.CreatedAt,
		&details.FullName,
		&details.FirstName,
		&details.LastName,
		&details.Email,
		&details.PhoneNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reservation with ID %s not found", idReservation)
		}
		return nil, fmt.Errorf("error retrieving reservation details: %v", err)
	}

	// Get client statistics (overall for the client)
	statsQuery := `
        SELECT 
            COUNT(DISTINCT r.idReservation) as totalVisits,
            IFNULL(AVG(ol.totalPrice), 0) as averageSpending,
            IFNULL(SUM(ol.totalPrice), 0) as totalSpent
        FROM reservation r
        LEFT JOIN orderList ol ON r.idReservation = ol.idReservation AND ol.status = 'completed'
        WHERE r.idClient = ?
    `

	err = s.db.QueryRow(statsQuery, idClient).Scan(
		&details.TotalVisits,
		&details.AverageSpending,
		&details.TotalSpent,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving client statistics: %v", err)
	}

	// Get total orders for THIS specific reservation
	orderCountQuery := `
        SELECT COUNT(*) as totalOrders
        FROM orderList ol
        WHERE ol.idReservation = ?
    `

	err = s.db.QueryRow(orderCountQuery, idReservation).Scan(&details.TotalOrders)
	if err != nil {
		return nil, fmt.Errorf("error retrieving order count for reservation: %v", err)
	}

	// Get favorite food (for the client overall)
	favoriteQuery := `
        SELECT f.name
        FROM orderFood orderFood
        JOIN food f ON orderFood.idFood = f.idFood
        JOIN orderList ol ON orderFood.idOrder = ol.idOrder
        JOIN reservation r ON ol.idReservation = r.idReservation
        WHERE r.idClient = ?
        GROUP BY f.idFood
        ORDER BY SUM(orderFood.quantity) DESC
        LIMIT 1
    `

	err = s.db.QueryRow(favoriteQuery, idClient).Scan(&details.FavoriteFood)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error retrieving favorite food: %v", err)
	}
	if err == sql.ErrNoRows {
		details.FavoriteFood = "No orders yet"
	}

	// Get list of orders for THIS SPECIFIC RESERVATION ONLY
	ordersQuery := `
        SELECT 
            ol.idOrder,
            ol.totalPrice,
            ol.createdAt,
            ol.status,
            COUNT(orderFood.idFood) as itemCount
        FROM orderList ol
        LEFT JOIN orderFood orderFood ON ol.idOrder = orderFood.idOrder
        WHERE ol.idReservation = ?
        GROUP BY ol.idOrder
        ORDER BY ol.createdAt DESC
    `

	rows, err := s.db.Query(ordersQuery, idReservation)
	if err != nil {
		return nil, fmt.Errorf("error retrieving reservation orders: %v", err)
	}
	defer rows.Close()

	var orders []types.ClientOrderSummary
	for rows.Next() {
		var order types.ClientOrderSummary
		err := rows.Scan(
			&order.IdOrder,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.Status,
			&order.ItemCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %v", err)
		}
		orders = append(orders, order)
	}

	details.Orders = orders
	return &details, nil
}

func (s *store) GetAllRestaurantReservations(idRestaurant string, page, limit int) (*types.PaginatedReservations, error) {
	offset := (page - 1) * limit

	// Get total count for pagination
	countQuery := `
        SELECT COUNT(*) 
        FROM reservation r
        JOIN client c ON r.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        WHERE r.idRestaurant = ?
    `
	var totalCount int
	err := s.db.QueryRow(countQuery, idRestaurant).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("error counting reservations: %v", err)
	}

	// Get reservations with pagination
	query := `
        SELECT 
            r.idReservation,
            r.timeFrom,
            CONCAT(p.firstName, ' ', p.lastName) as fullName,
            r.idTable,
            r.numberOfPeople,
            r.status,
            r.createdAt
        FROM reservation r
        JOIN client c ON r.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        WHERE r.idRestaurant = ?
        ORDER BY r.timeFrom DESC
        LIMIT ? OFFSET ?
    `

	rows, err := s.db.Query(query, idRestaurant, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving reservations: %v", err)
	}
	defer rows.Close()

	var reservations []types.RestaurantReservationDetail
	for rows.Next() {
		var reservation types.RestaurantReservationDetail
		var idTable sql.NullString

		err := rows.Scan(
			&reservation.IdReservation,
			&reservation.TimeFrom,
			&reservation.FullName,
			&idTable,
			&reservation.NumberOfPeople,
			&reservation.Status,
			&reservation.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning reservation: %v", err)
		}

		if idTable.Valid {
			reservation.TableId = idTable.String
		} else {
			reservation.TableId = "Not assigned"
		}

		reservations = append(reservations, reservation)
	}

	totalPages := (totalCount + limit - 1) / limit

	return &types.PaginatedReservations{
		Reservations: reservations,
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalCount:   totalCount,
		HasNext:      page < totalPages,
		HasPrevious:  page > 1,
	}, nil
}

func (s *store) GetOrderInformation(idOrder string) (*types.OrderInformation, error) {
	profileQuery := `
        SELECT 
            ol.idOrder,
            ol.totalPrice,
            ol.status,
            ol.createdAt,
            profile.firstName,
            profile.lastName,
            profile.email,
            profile.phoneNumber,
            profile.address,
            client.username,
            r.timeFrom,
            r.numberOfPeople
        FROM orderList ol
        JOIN reservation r ON ol.idReservation = r.idReservation
        JOIN client ON r.idClient = client.idClient
        JOIN profile ON client.idProfile = profile.idProfile
        WHERE ol.idOrder = ?
    `

	row := s.db.QueryRow(profileQuery, idOrder)
	var orderInfo types.OrderInformation
	err := row.Scan(
		&orderInfo.IdOrder,
		&orderInfo.TotalPrice,
		&orderInfo.Status,
		&orderInfo.CreatedAt,
		&orderInfo.ClientFirstName,
		&orderInfo.ClientLastName,
		&orderInfo.ClientEmail,
		&orderInfo.ClientPhone,
		&orderInfo.ClientAddress,
		&orderInfo.ClientUsername,
		&orderInfo.ReservationTime,
		&orderInfo.NumberOfPeople,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order with ID %s not found", idOrder)
		}
		return nil, fmt.Errorf("error retrieving order information: %v", err)
	}

	foodQuery := `
        SELECT 
            food.idFood,
            food.name,
            food.description,
            food.image,
            food.price,
            orderFood.quantity,
            (food.price * orderFood.quantity) as subtotal
        FROM orderFood 
        JOIN food ON orderFood.idFood = food.idFood
        WHERE orderFood.idOrder = ?
    `

	rows, err := s.db.Query(foodQuery, idOrder)
	if err != nil {
		return nil, fmt.Errorf("error retrieving food items: %v", err)
	}
	defer rows.Close()

	var foodItems []types.OrderFoodItem
	for rows.Next() {
		var foodItem types.OrderFoodItem
		err := rows.Scan(
			&foodItem.IdFood,
			&foodItem.Name,
			&foodItem.Description,
			&foodItem.Image,
			&foodItem.Price,
			&foodItem.Quantity,
			&foodItem.Subtotal,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning food item: %v", err)
		}
		foodItems = append(foodItems, foodItem)
	}

	orderInfo.FoodItems = foodItems
	return &orderInfo, nil
}

func (s *store) UpdateOrderStatus(idOrder string, status string) error {
	var currentStatus string
	checkQuery := `SELECT status FROM orderList WHERE idOrder = ?`
	err := s.db.QueryRow(checkQuery, idOrder).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("order with ID %s not found", idOrder)
		}
		return fmt.Errorf("error checking order status: %v", err)
	}

	if currentStatus == "completed" && status == "completed" {
		return fmt.Errorf("order is already completed")
	}

	updateQuery := `UPDATE orderList SET status = ? WHERE idOrder = ?`
	result, err := s.db.Exec(updateQuery, status, idOrder)
	if err != nil {
		return fmt.Errorf("error updating order status: %v", err)
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

func (s *store) GetAllClientReservations(idClient string) ([]types.ClientReservationInfo, error) {
	query := `
        SELECT 
            r.idReservation,
            r.timeFrom,
            r.numberOfPeople,
            r.status,
            r.createdAt,
            rest.name,
            rest.image,
            rest.location,
            rest.idRestaurant
        FROM reservation r
        JOIN restaurant rest ON r.idRestaurant = rest.idRestaurant
        WHERE r.idClient = ?
        ORDER BY r.timeFrom DESC
    `
	rows, err := s.db.Query(query, idClient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []types.ClientReservationInfo
	for rows.Next() {
		var reservation types.ClientReservationInfo
		if err := rows.Scan(
			&reservation.IdReservation,
			&reservation.TimeFrom,
			&reservation.NumberOfPeople,
			&reservation.Status,
			&reservation.CreatedAt,
			&reservation.RestaurantName,
			&reservation.RestaurantImage,
			&reservation.RestaurantLocation,
			&reservation.IdRestaurant,
		); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (s *store) GetUpcomingReservations(restaurantId string) ([]types.UpcomingReservationInfo, error) {
	query := `
        SELECT 
            r.idReservation,
            profile.firstName,
            profile.lastName,
            r.numberOfPeople,
            DATE(r.timeFrom) as date,
            DAYNAME(r.timeFrom) as day,
            TIME_FORMAT(r.timeFrom, '%H:%i') as time,
            r.idTable
        FROM reservation r
        JOIN client c ON r.idClient = c.idClient
        JOIN profile ON c.idProfile = profile.idProfile
        WHERE r.idRestaurant = ? AND DATE(r.timeFrom) > CURDATE()
        ORDER BY r.timeFrom ASC
        LIMIT 4
    `
	rows, err := s.db.Query(query, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []types.UpcomingReservationInfo
	for rows.Next() {
		var r types.UpcomingReservationInfo
		if err := rows.Scan(&r.IdReservation, &r.FirstName, &r.LastName, &r.NumberPeople, &r.Date, &r.Day, &r.Time, &r.IdTable); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	return reservations, nil
}

func (s *store) CreateMenu(idMenu, idRestaurant, name string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE menu SET active = 0 WHERE idRestaurant = ?`, idRestaurant)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Insert the new menu as active
	_, err = tx.Exec(`INSERT INTO menu (idMenu, idRestaurant, name, active) VALUES (?, ?, ?, 1)`, idMenu, idRestaurant, name)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *store) CreateFood(idFood, idCategory, idRestaurant, name, description, image string, price float64, status string) error {
	query := `INSERT INTO food (idFood, idCategory, idRestaurant, name, description, image, price, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, idFood, idCategory, idRestaurant, name, description, image, price, status)
	return err
}

func (s *store) CreateFoodCategory(idCategory, nameCategorie string) error {
	query := `INSERT INTO foodCategory (idCategory, nameCategorie) VALUES (?, ?)`
	_, err := s.db.Exec(query, idCategory, nameCategorie)
	return err
}

func (s *store) DeleteFood(idFood string) error {
	query := `DELETE FROM food WHERE idFood = ?`
	_, err := s.db.Exec(query, idFood)
	return err
}

func (s *store) GetFoodById(idFood string) (*types.Food, error) {
	query := `SELECT * FROM food WHERE idFood = ?`
	row := s.db.QueryRow(query, idFood)
	var food types.Food
	err := row.Scan(
		&food.IdFood, &food.IdCategory, &food.Name,
		&food.Description, &food.Image, &food.Price, &food.Status,
	)
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (s *store) AddFoodToMenu(idMenuFood, idMenu, idFood string) error {
	query := `INSERT INTO menufood (idMenuFood, idMenu, idFood) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, idMenuFood, idMenu, idFood)
	return err
}

func (s *store) SetFoodStatusInMenu(idFood, status string) error {
	query := `UPDATE food SET status = ? WHERE idFood = ?`
	_, err := s.db.Exec(query, status, idFood)
	return err
}

func (s *store) SetMenuActive(idMenu, idRestaurant string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var menuExists int
	checkQuery := `SELECT COUNT(*) FROM menu WHERE idMenu = ? AND idRestaurant = ?`
	err = tx.QueryRow(checkQuery, idMenu, idRestaurant).Scan(&menuExists)
	if err != nil {
		return fmt.Errorf("error checking menu existence: %v", err)
	}
	if menuExists == 0 {
		return fmt.Errorf("menu with ID %s not found for restaurant %s", idMenu, idRestaurant)
	}

	deactivateQuery := `UPDATE menu SET active = 0 WHERE idRestaurant = ?`
	_, err = tx.Exec(deactivateQuery, idRestaurant)
	if err != nil {
		return fmt.Errorf("error deactivating existing menus: %v", err)
	}

	// Activate the selected menu
	activateQuery := `UPDATE menu SET active = 1 WHERE idMenu = ? AND idRestaurant = ?`
	result, err := tx.Exec(activateQuery, idMenu, idRestaurant)
	if err != nil {
		return fmt.Errorf("error activating menu: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no menu was activated")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *store) GetRestaurantWorkerWithRatings(idRestaurantWorker string) (*types.RestaurantWorkerWithRatings, error) {
	// Get worker information
	workerQuery := `
        SELECT 
            idRestaurantWorker,
            firstName,
            lastName,
            email,
            phoneNumber,
            IFNULL(quote, '') as quote,
            startWorking,
            IFNULL(nationnallity, '') as nationnallity,
            IFNULL(nativeLanguage, '') as nativeLanguage,
            IFNULL(rating, 0) as rating,
            IFNULL(image, '') as image,
            IFNULL(address, '') as address,
            status,
            idRestaurant
        FROM restaurantWorkers 
        WHERE idRestaurantWorker = ?
    `

	row := s.db.QueryRow(workerQuery, idRestaurantWorker)
	var worker types.RestaurantWorkerWithRatings

	err := row.Scan(
		&worker.IdRestaurantWorker,
		&worker.FirstName,
		&worker.LastName,
		&worker.Email,
		&worker.PhoneNumber,
		&worker.Quote,
		&worker.StartWorking,
		&worker.Nationnallity,
		&worker.NativeLanguage,
		&worker.Rating,
		&worker.Image,
		&worker.Address,
		&worker.Status,
		&worker.IdRestaurant,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("restaurant worker with ID %s not found", idRestaurantWorker)
		}
		return nil, fmt.Errorf("error retrieving restaurant worker: %v", err)
	}

	// Get recent ratings for this worker
	ratingsQuery := `
        SELECT 
            r.rating,
            IFNULL(r.comment, '') as comment,
            r.createdAt,
            p.firstName,
            p.lastName
        FROM rating r
        JOIN client c ON r.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        WHERE r.idRestaurantWorker = ? AND r.ratingType = 'worker'
        ORDER BY r.createdAt DESC
        LIMIT 10
    `

	rows, err := s.db.Query(ratingsQuery, idRestaurantWorker)
	if err != nil {
		return nil, fmt.Errorf("error retrieving worker ratings: %v", err)
	}
	defer rows.Close()

	var ratings []types.WorkerRating
	for rows.Next() {
		var rating types.WorkerRating
		err := rows.Scan(
			&rating.RatingValue,
			&rating.Comment,
			&rating.CreatedAt,
			&rating.ClientFirstName,
			&rating.ClientLastName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning worker rating: %v", err)
		}
		ratings = append(ratings, rating)
	}

	// Get rating statistics - handle NULL/empty case
	statsQuery := `
        SELECT 
            IFNULL(COUNT(*), 0) as totalRatings,
            IFNULL(AVG(rating), 0) as averageRating,
            IFNULL(SUM(CASE WHEN rating = 5 THEN 1 ELSE 0 END), 0) AS count5Stars,
            IFNULL(SUM(CASE WHEN rating = 4 THEN 1 ELSE 0 END), 0) AS count4Stars,
            IFNULL(SUM(CASE WHEN rating = 3 THEN 1 ELSE 0 END), 0) AS count3Stars,
            IFNULL(SUM(CASE WHEN rating = 2 THEN 1 ELSE 0 END), 0) AS count2Stars,
            IFNULL(SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END), 0) AS count1Star
        FROM rating
        WHERE idRestaurantWorker = ? AND ratingType = 'worker'
    `

	var totalRatings, count5Stars, count4Stars, count3Stars, count2Stars, count1Star int
	var averageRating float64

	err = s.db.QueryRow(statsQuery, idRestaurantWorker).Scan(
		&totalRatings,
		&averageRating,
		&count5Stars,
		&count4Stars,
		&count3Stars,
		&count2Stars,
		&count1Star,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving worker rating stats: %v", err)
	}

	// Calculate percentages
	percent := func(count int) float64 {
		if totalRatings == 0 {
			return 0
		}
		return float64(count) * 100 / float64(totalRatings)
	}

	// Initialize ratings slice if nil
	if ratings == nil {
		ratings = []types.WorkerRating{}
	}

	worker.RecentRatings = ratings
	worker.RatingStats = types.WorkerRatingStats{
		TotalRatings:     totalRatings,
		AverageRating:    averageRating,
		Percentage5Stars: percent(count5Stars),
		Percentage4Stars: percent(count4Stars),
		Percentage3Stars: percent(count3Stars),
		Percentage2Stars: percent(count2Stars),
		Percentage1Star:  percent(count1Star),
	}

	return &worker, nil
}

func (s *store) GetFoodRestaurant(idRestaurant string) (*[]types.Food, error) {
	query := `
        SELECT DISTINCT f.idFood, f.idCategory, f.name, f.description, f.image, f.price, f.status
        FROM food f
        WHERE f.idRestaurant = ?
    `
	rows, err := s.db.Query(query, idRestaurant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var foods []types.Food
	for rows.Next() {
		var food types.Food
		if err := rows.Scan(&food.IdFood, &food.IdCategory, &food.Name, &food.Description, &food.Image, &food.Price, &food.Status); err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return &foods, nil
}

func (s *store) GetMenuWithFoods(idMenu string) (*types.Menu, *[]types.Food, error) {
	menuQuery := `SELECT * FROM menu WHERE idMenu = ?`
	row := s.db.QueryRow(menuQuery, idMenu)
	var menu types.Menu
	err := row.Scan(&menu.IdMenu, &menu.IdRestaurant, &menu.Name, &menu.Active, &menu.CreatedAt)
	if err != nil {
		return nil, nil, err
	}
	foods, err := s.GetFoodByMenu(idMenu)
	if err != nil {
		return &menu, nil, err
	}
	return &menu, foods, nil
}

func (s *store) DeleteTable(idTable string) error {
	query := `DELETE FROM table_restaurant WHERE idTable = ?`
	_, err := s.db.Exec(query, idTable)
	return err
}

func (s *store) GetTablesByRestaurant(restaurantId string) ([]types.Table, error) {
	query := `SELECT idTable, idRestaurant, shape, posX, posY, is_available FROM table_restaurant WHERE idRestaurant = ?`
	rows, err := s.db.Query(query, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []types.Table
	for rows.Next() {
		var t types.Table
		if err := rows.Scan(&t.IdTable, &t.IdRestaurant, &t.Shape, &t.PosX, &t.PosY, &t.IsAvailable); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

func (s *store) UpdateReservationStatus(idReservation, status string) error {
	// First, get current reservation status and time
	var currentStatus string
	var timeFrom time.Time
	query := `SELECT status, timeFrom FROM reservation WHERE idReservation = ?`
	err := s.db.QueryRow(query, idReservation).Scan(&currentStatus, &timeFrom)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("reservation not found")
		}
		return fmt.Errorf("error fetching reservation: %v", err)
	}

	// Validate status transition
	if !isValidReservationStatusTransition(currentStatus, status) {
		return fmt.Errorf("invalid status transition from %s to %s", currentStatus, status)
	}
	if status == "confirmed" {
		now := time.Now()
		timeDiff := timeFrom.Sub(now)
		absTimeDiff := timeDiff
		if absTimeDiff < 0 {
			absTimeDiff = -absTimeDiff
		}

		twoHours := 2 * time.Hour
		if absTimeDiff > twoHours {
			return fmt.Errorf("status can only be changed within 2 hours before or after the reservation time")
		}

	}

	// If validation passes, update the status
	updateQuery := `UPDATE reservation SET status = ? WHERE idReservation = ?`
	_, err = s.db.Exec(updateQuery, status, idReservation)
	return err
}

func isValidReservationStatusTransition(currentStatus, newStatus string) bool {
	switch currentStatus {
	case "pending":
		return newStatus == "confirmed" || newStatus == "cancelled"
	case "confirmed":
		return false // confirmed cannot change to any other status
	case "cancelled":
		return false // cancelled cannot change to any other status
	default:
		return false
	}
}

func (s *store) GetRestaurantMenuStats(restaurantId string) (*types.RestaurantMenuStats, error) {
	stats := &types.RestaurantMenuStats{}

	// Total menus
	err := s.db.QueryRow(`SELECT COUNT(*) FROM menu WHERE idRestaurant = ?`, restaurantId).Scan(&stats.TotalMenus)
	if err != nil {
		return nil, err
	}

	// Active menu id and name
	var activeMenuId string
	err = s.db.QueryRow(`SELECT idMenu, name FROM menu WHERE idRestaurant = ? AND active = 1 LIMIT 1`, restaurantId).Scan(&activeMenuId, &stats.ActiveMenuName)
	if err != nil {
		// If no active menu, return stats with zeroes
		return stats, nil
	}

	// Total items in active menu
	err = s.db.QueryRow(`SELECT COUNT(*) FROM food join menufood on food.idFood=menufood.idFood WHERE menufood.idMenu = ?`, activeMenuId).Scan(&stats.TotalItems)
	if err != nil {
		return nil, err
	}

	// Number of categories in active menu
	err = s.db.QueryRow(`SELECT COUNT(DISTINCT idCategory) FROM food join menufood on food.idFood=menufood.idFood WHERE menufood.idMenu = ?`, activeMenuId).Scan(&stats.TotalCategories)
	if err != nil {
		return nil, err
	}

	// Available foods in active menu
	err = s.db.QueryRow(`SELECT COUNT(*) FROM food join menufood on food.idFood=menufood.idFood WHERE menufood.idMenu = ? AND food.status = 'available'`, activeMenuId).Scan(&stats.AvailableFoods)
	if err != nil {
		return nil, err
	}

	// Unavailable foods in active menu
	err = s.db.QueryRow(`SELECT COUNT(*) FROM food join menufood on food.idFood=menufood.idFood WHERE menufood.idMenu = ? AND food.status != 'available'`, activeMenuId).Scan(&stats.UnavailableFoods)
	if err != nil {
		return nil, err
	}

	// Top 4 popular foods of the restaurant (all time)
	rows, err := s.db.Query(`
        SELECT f.name, COUNT(ol.idFood) as orderCount
        FROM food f
        join menufood on f.idFood=menufood.idFood
        JOIN menu m ON menufood.idMenu = m.idMenu
        JOIN orderFood ol ON f.idFood = ol.idFood
        JOIN orderList o ON ol.idOrder = o.idOrder
        WHERE m.idRestaurant = ?
        GROUP BY f.idFood
        ORDER BY orderCount DESC
        LIMIT 4
    `, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pf types.PopularFood
		if err := rows.Scan(&pf.FoodName, &pf.OrderCount); err != nil {
			return nil, err
		}
		stats.PopularFoods = append(stats.PopularFoods, pf)
	}

	return stats, nil
}

func (s *store) GetFoodsOfActiveMenu(idRestaurant string) ([]types.Food, error) {
	var idMenu string
	err := s.db.QueryRow("SELECT idMenu FROM menu WHERE idRestaurant = ? AND active = 1 LIMIT 1", idRestaurant).Scan(&idMenu)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT food.idFood, idCategory, menufood.idMenu, name, description, image, price, status FROM food join menufood on food.idFood = menufood.idFood WHERE menufood.idMenu = ?", idMenu)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var foods []types.Food
	for rows.Next() {
		var food types.Food
		if err := rows.Scan(&food.IdFood, &food.IdCategory, &food.IdMenu, &food.Name, &food.Description, &food.Image, &food.Price, &food.Status); err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return foods, nil
}

func (s *store) GetMenusByRestaurant(idRestaurant string) ([]types.Menu, error) {
	query := `SELECT idMenu, idRestaurant, name, active, createdAt FROM menu WHERE idRestaurant = ?`
	rows, err := s.db.Query(query, idRestaurant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var menus []types.Menu
	for rows.Next() {
		var menu types.Menu
		if err := rows.Scan(&menu.IdMenu, &menu.IdRestaurant, &menu.Name, &menu.Active, &menu.CreatedAt); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func (s *store) UpdateFood(idFood string, food types.Food) error {
	query := `UPDATE food SET idCategory=?,  name=?, description=?, image=?, price=?, status=? WHERE idFood=?`
	_, err := s.db.Exec(query, food.IdCategory, food.Name, food.Description, food.Image, food.Price, food.Status, idFood)
	return err
}

func (s *store) UpdateRestaurantWorker(id string, worker types.RestaurantWorker) error {
	query := `UPDATE restaurantWorkers SET firstName=?, lastName=?, email=?, phoneNumber=?, quote=?, startWorking=?, nationnallity=?, nativeLanguage=?, rating=?, address=?, status=? WHERE idRestaurantWorker=?`
	_, err := s.db.Exec(query, worker.FirstName, worker.LastName, worker.Email, worker.PhoneNumber, worker.Quote, worker.StartWorking, worker.Nationnallity, worker.NativeLanguage, worker.Rating, worker.Address, worker.Status, id)
	return err
}

// func (s *store) PostFeedbackRestaurant(feedback types.FeedbackRestaurant) error {
// 	query := `INSERT INTO feedbackRestaurant (idClient, idRestaurant, comment, createdAt) VALUES (?, ?, ?, NOW())`
// 	_, err := s.db.Exec(query, feedback.IdClient, feedback.IdRestaurant, feedback.Comment)
// 	return err
// }

// func (s *store) PostFeedbackWorker(feedback types.FeedbackWorker) error {
// 	query := `INSERT INTO feedbackWorker (idClient, idRestaurantWorker, comment, createdAt) VALUES (?, ?, ?, NOW())`
// 	_, err := s.db.Exec(query, feedback.IdClient, feedback.IdRestaurantWorker, feedback.Comment)
// 	return err
// }

func (s *store) CreateNotification(notification types.Notification) error {
	query := `INSERT INTO notifications (idNotification, idAdmin, titre, type, description) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, notification.IdNotification, notification.IdAdmin, notification.Titre, notification.Type, notification.Description)
	return err
}

func (s *store) GetNotifications() ([]types.Notification, error) {
	query := `SELECT idNotification, idAdmin, titre, type, description FROM notifications  `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifs []types.Notification
	for rows.Next() {
		var n types.Notification
		if err := rows.Scan(&n.IdNotification, &n.IdAdmin, &n.Titre, &n.Type, &n.Description); err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (s *store) UpdateTable(idTable string, table types.Table) error {
	query := `UPDATE table_restaurant SET shape=?, posX=?, posY=?, is_available=? WHERE idTable=?`
	_, err := s.db.Exec(query, table.Shape, table.PosX, table.PosY, table.IsAvailable, idTable)
	return err
}

func (s *store) CreateTable(table types.Table) error {
	query := `INSERT INTO table_restaurant (idTable, idRestaurant, shape, posX, posY, is_available) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, table.IdTable, table.IdRestaurant, table.Shape, table.PosX, table.PosY, table.IsAvailable)
	return err
}

func (s *store) GetFoodCategoriesByRestaurant() ([]types.FoodCategory, error) {
	query := `
        SELECT DISTINCT fc.idCategory, fc.nameCategorie
        FROM foodCategory fc
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []types.FoodCategory
	for rows.Next() {
		var cat types.FoodCategory
		if err := rows.Scan(&cat.IdCategory, &cat.NameCategorie); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (s *store) SetFoodUnavailable(idFood string) error {
	query := `UPDATE food SET status = 'unavailable' WHERE idFood = ?`

	_, err := s.db.Exec(query, idFood)
	return err
}

func (s *store) GetTableOccupationToday(idRestaurant string) ([]types.TableOccupation, error) {
	query := `
        SELECT t.idTable, r.timeFrom
        FROM table_restaurant t
        LEFT JOIN reservation r ON t.idTable = r.idTable
            AND DATE(r.timeFrom) = CURDATE()
            AND r.idRestaurant = ?
        WHERE t.idRestaurant = ?
        ORDER BY t.idTable, r.timeFrom
    `
	rows, err := s.db.Query(query, idRestaurant, idRestaurant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tableMap := make(map[string]*types.TableOccupation)
	for rows.Next() {
		var idTable string
		var timeFrom sql.NullTime
		if err := rows.Scan(&idTable, &timeFrom); err != nil {
			return nil, err
		}
		if _, exists := tableMap[idTable]; !exists {
			tableMap[idTable] = &types.TableOccupation{
				IdTable:   idTable,
				Occupied:  false,
				TimeSlots: []string{},
			}
		}
		if timeFrom.Valid {
			tableMap[idTable].Occupied = true
			slot := fmt.Sprintf("%02d:%02d", timeFrom.Time.Hour(), timeFrom.Time.Minute())
			tableMap[idTable].TimeSlots = append(tableMap[idTable].TimeSlots, slot)
		}
	}
	result := []types.TableOccupation{}
	for _, v := range tableMap {
		result = append(result, *v)
	}
	return result, nil
}

func (s *store) CountFirstTimeReservers(idRestaurant string) (int, error) {
	query := `
        SELECT COUNT(*) FROM (
            SELECT r.idClient
            FROM reservation r
            WHERE r.idRestaurant = ?
            AND r.idReservation = (
                SELECT r2.idReservation
                FROM reservation r2
                WHERE r2.idClient = r.idClient
                ORDER BY r2.timeFrom ASC
                LIMIT 1
            )
            GROUP BY r.idClient
        ) AS first_time_users;
    `
	row := s.db.QueryRow(query, idRestaurant)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *store) GetTopFoodsThisWeek(idRestaurant string) ([]types.FoodPopularity, error) {
	query := `
        SELECT f.idFood, f.name, f.image, SUM(ofd.quantity) as total
        FROM orderFood ofd
        JOIN orderList ol ON ofd.idOrder = ol.idOrder
        JOIN reservation r ON ol.idReservation = r.idReservation
        JOIN food f ON ofd.idFood = f.idFood
        WHERE r.idRestaurant = ?
          AND YEARWEEK(ol.createdAt, 1) = YEARWEEK(CURDATE(), 1)
        GROUP BY f.idFood, f.name, f.image
        ORDER BY total DESC
        LIMIT 3
    `
	rows, err := s.db.Query(query, idRestaurant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []types.FoodPopularity
	for rows.Next() {
		var food types.FoodPopularity
		if err := rows.Scan(&food.IdFood, &food.Name, &food.Image, &food.Total); err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return foods, nil
}

func (s *store) CreateRestaurant(idRestaurant, idAdminRestaurant, name, image string, longitude, latitude float64, description string, capacity int, location string) error {
	query := `INSERT INTO restaurant (idRestaurant, idAdminRestaurant, name, image, longitude, latitude, description, capacity, location) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, idRestaurant, idAdminRestaurant, name, image, longitude, latitude, description, capacity, location)
	return err
}

func (s *store) CreateRestaurantWorker(id, idRestaurant string, worker types.RestaurantWorkerCreation) error {
	checkQuery := `SELECT COUNT(*) FROM restaurantWorkers WHERE email = ?`
	var count int
	err := s.db.QueryRow(checkQuery, worker.Email).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking email existence: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("email %s already exists", worker.Email)
	}

	query := `
        INSERT INTO restaurantWorkers (
            idRestaurantWorker, idRestaurant, firstName, lastName, email, phoneNumber, quote,
            startWorking, nationnallity, nativeLanguage, rating, address, image, status
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'active')
    `
	_, err = s.db.Exec(query,
		id, idRestaurant, worker.FirstName, worker.LastName,
		worker.Email, worker.PhoneNumber, worker.Quote, time.Now(), worker.Nationnallity,
		worker.NativeLanguage, 0, worker.Address, worker.Image,
	)
	return err
}

func (s *store) SetRestaurantWorkerStatus(idRestaurantWorker string, status string) error {
	query := `UPDATE restaurantWorkers SET status = ? WHERE idRestaurantWorker = ?`
	_, err := s.db.Exec(query, status, idRestaurantWorker)
	return err
}

func (s *store) GetRestaurantByIdProfile(idProfile string) (*types.UserAdmin, error) {
	query := `SELECT 
	  profile.idProfile AS profileId,
	  profile.firstName,
	  profile.lastName,
	  profile.email,
	  profile.createdAt,
	  profile.type,
	  profile.address,
	  profile.lastLogin,
	  profile.phoneNumber,
      adminRestaurant.idAdminRestaurant,
      restaurant.idRestaurant
	FROM profile 
    join adminRestaurant ON profile.idProfile = adminRestaurant.idProfile
    join restaurant ON adminRestaurant.idAdminRestaurant = restaurant.idAdminRestaurant
	WHERE profile.idProfile = ?`
	row := s.db.QueryRow(query, idProfile)
	var rest types.UserAdmin
	err := row.Scan(
		&rest.Id,
		&rest.FirstName,
		&rest.LastName,
		&rest.Email,
		&rest.CreatedAt,
		&rest.Type,
		&rest.Address,
		&rest.LastLogin,
		&rest.Phone,
		&rest.IdAdminRestaurant,
		&rest.IdRestaurant,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No restaurant found for profile ID %s", idProfile)
			return nil, nil // Return nil if no restaurant is found
		}
		log.Printf("Error retrieving restaurant by profile ID %s: %v", idProfile, err)
		return nil, fmt.Errorf("error retrieving restaurant by profile ID: %v", err)
	}
	return &rest, nil
}

func (s *store) GetAvailableMenuInformation(restaurantId string) (*[]types.MenuInformationFood, error) {
	query := `
SELECT food.*,menu.idMenu,menu.name as menuName
 FROM menu
 join menufood on menufood.idMenu=menu.idMenu
JOIN food ON food.idFood = menufood.idFood
where menu.active = 1 and food.status="available" and menu.idRestaurant = ?;
`
	rows, err := s.db.Query(query, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var menuInformation []types.MenuInformationFood
	for rows.Next() {

		var menu types.MenuInformationFood
		err = rows.Scan(
			&menu.IdFood,
			&menu.IdRestaurant,
			&menu.IdCategory,
			&menu.Name,
			&menu.Description,
			&menu.Image,
			&menu.Price,
			&menu.Status,
			&menu.IdMenu,
			&menu.MenuName,
		)
		if err != nil {
			return nil, err
		}

		menuInformation = append(menuInformation, menu)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &menuInformation, nil
}

// !NOTE: GET all restaurant
func (s *store) CountOrderReceivedToday(idRestaurant string) (int, error) {
	query := `SELECT COUNT(*) FROM orderList join reservation on orderList.idReservation=reservation.idReservation WHERE DATE(orderList.createdAt) = CURDATE() and reservation.idRestaurant = ?`
	row := s.db.QueryRow(query, idRestaurant)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Printf("Error counting orders received today: %v", err)
		return 0, err
	}
	return count, nil
}

func (s *store) CountReservationReceivedToday(idRestaurant string) (int, error) {
	query := `SELECT COUNT(*) FROM reservation WHERE DATE(reservation.createdAt) = CURDATE() and idRestaurant = ?`
	row := s.db.QueryRow(query, idRestaurant)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Printf("Error counting reservations received today: %v", err)
		return 0, err
	}
	return count, nil
}

func (s *store) CountReservationThisMonth() (int, error) {
	query := `SELECT COUNT(*) FROM reservation WHERE MONTH(createdAt) = MONTH(CURDATE()) AND YEAR(createdAt) = YEAR(CURDATE())`
	row := s.db.QueryRow(query)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Printf("Error counting reservations this month: %v", err)
		return 0, err
	}
	return count, nil
}

func (s *store) PostOrderList(orderId string, foods []types.FoodItem) error {
	var totalPrice float64
	if len(foods) == 0 {
		log.Println("⚠️ No foods provided for order:", orderId)
	}
	log.Printf("Inserting %d foods into order %s", len(foods), orderId)
	for _, food := range foods {
		res, err := s.db.Exec(`Insert INTO orderFood (idOrder, idFood, quantity, createdAt) VALUES (?, ?, ?, ?)`, orderId, food.IdFood, food.Quantity, time.Now())
		totalPrice += food.PriceSingle * float64(food.Quantity)
		if err != nil {
			log.Printf("Error inserting into orderFood: %v", err)

			return err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			log.Printf("Errior getting rows affected: %v", err)
			return err
		}

	}
	query := `UPDATE orderList SET totalPrice = ? WHERE idOrder = ?`
	res, err := s.db.Exec(query, totalPrice, orderId)
	if err != nil {
		log.Printf("Error inserting into orderFood: %v", err)
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("⚠️ No rows updated in orderList for idOrder: %s", orderId)
	}
	return nil
}

func (s *store) CreateReservation(idReservation string, reservation types.ReservationCreation) error {
	date := reservation.TimeFrom.Format("2006-01-02")
	checkQuery := `SELECT COUNT(*) FROM reservation WHERE idClient = ? AND DATE(timeFrom) = ?`
	var count int
	err := s.db.QueryRow(checkQuery, reservation.IdClient, date).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("You already has a reservation on %s", date)
	}

	// Insert reservation
	query := `
		INSERT INTO reservation (
			idReservation, idClient, idRestaurant, idTable,
			status, createdAt, numberOfPeople, timeFrom
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = s.db.Exec(query,
		idReservation,
		reservation.IdClient,
		reservation.IdRestaurant,
		reservation.TableId,
		"pending",
		time.Now(),
		reservation.NumberOfPeople,
		reservation.TimeFrom,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) ReserveTable(idReservation string, reservation types.ReservationCreation) error {
	query := `INSERT INTO table_reservation (idTable, idReservation, numberOfPeople, timeFrom) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, reservation.TableId, idReservation, reservation.NumberOfPeople, reservation.TimeFrom)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) CreateOrder(idOrder string, order types.OrderCreation) error {
	query := `INSERT INTO orderList (idOrder, idReservation, totalPrice, status, createdAt) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, idOrder, order.IdReservation, 0, "pending", time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (s *store) AddFoodToOrder(food types.AddFoodToOrder) error {
	query := `INSERT INTO orderFood (idOrder, idFood, quantity,createdAt) VALUES (?, ?, ?,?)`
	_, err := s.db.Exec(query, food.IdOrder, food.IdFood, food.Quantity, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (s *store) BulkUpdateRestaurantTables(idRestaurant string, tables []types.Table) error {
	// Start a transaction to ensure atomicity
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete all existing tables for this restaurant
	deleteQuery := `DELETE FROM table_restaurant WHERE idRestaurant = ?`
	_, err = tx.Exec(deleteQuery, idRestaurant)
	if err != nil {
		return fmt.Errorf("error deleting existing tables: %v", err)
	}

	// Insert new tables if any are provided
	if len(tables) > 0 {
		insertQuery := `INSERT INTO table_restaurant (idTable, idRestaurant, shape, posX, posY, is_available) VALUES (?, ?, ?, ?, ?, ?)`

		for _, table := range tables {
			// Generate ID if not provided
			if table.IdTable == "" {
				id, err := utils.CreateAnId()
				if err != nil {
					return fmt.Errorf("error generating table ID: %v", err)
				}
				table.IdTable = id
			}

			// Set restaurant ID
			table.IdRestaurant = idRestaurant

			// Set availability to true by default (since frontend doesn't send this field)
			table.IsAvailable = true

			_, err = tx.Exec(insertQuery, table.IdTable, table.IdRestaurant, table.Shape, table.PosX, table.PosY, table.IsAvailable)
			if err != nil {
				return fmt.Errorf("error inserting table %s: %v", table.IdTable, err)
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *store) GetRestaurantTables(restaurantId string, timeReserved time.Time) (*[]types.RestaurantTableStatus, error) {
	query := `SELECT tr.idTable, tr.idRestaurant, tr.shape, r.idReservation, tr.posX, tr.posY, r.timeFrom, r.numberOfPeople,
    IF(r.idReservation IS NOT NULL, 'reserved', 'available') AS status
FROM 
    table_restaurant tr
LEFT JOIN 
    reservation r 
    ON tr.idTable = r.idTable 
    AND r.timeFrom = ?
WHERE tr.idRestaurant = ?;
`

	log.Println("Restaurant ID:", restaurantId)
	log.Println("Time Reserved:", timeReserved.Format("2006-01-02 15:04:05"))
	rows, err := s.db.Query(query, timeReserved, restaurantId)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var tables []types.RestaurantTableStatus
	for rows.Next() {
		var table types.RestaurantTableStatus
		err = rows.Scan(
			&table.IdTable,
			&table.IdRestaurant,
			&table.Shape,
			&table.IdReservation,
			&table.PosX,
			&table.PosY,
			&table.TimeFrom,
			&table.NumberOfPeople,
			&table.Status,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after iterating rows:", err)
		return nil, err
	}

	return &tables, nil
}

func (s *store) GetRestaurant() (*[]types.Restaurant, error) {
	query := `
		SELECT 
			r.idRestaurant,
			r.idAdminRestaurant,
			r.name,
			r.image,
			r.longitude,
			r.latitude,
			r.description,
			r.capacity,
			r.location,
			COALESCE(AVG(rating.rating), 0) as averageRating
		FROM restaurant r
		LEFT JOIN rating ON r.idRestaurant = rating.idRestaurant
		GROUP BY r.idRestaurant, r.idAdminRestaurant, r.name, r.image, r.longitude, r.latitude, r.description, r.capacity, r.location
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var restaurant []types.Restaurant

	for rows.Next() {
		var rest types.Restaurant
		var averageRating sql.NullFloat64

		err = rows.Scan(
			&rest.IdRestaurant,
			&rest.IdAdminRestaurant,
			&rest.NameRestaurant,
			&rest.Image,
			&rest.Langitude,
			&rest.Latitude,
			&rest.Description,
			&rest.Capacity,
			&rest.Location,
			&averageRating,
		)
		if err != nil {
			return nil, err
		}

		// Set average rating if available
		if averageRating.Valid {
			rating := averageRating.Float64
			rest.AverageRating = &rating
		}

		// Determine if restaurant is active (has admin assigned)
		isActive := rest.IdAdminRestaurant != nil
		rest.IsActive = &isActive

		restaurant = append(restaurant, rest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &restaurant, nil
}

// // func (s *store) GetCategorieFoods() (*[]types.FoodCategory, error) {
// //     query := `SELECT * FROM foodCategory`
// //     rows, err := s.db.Query(query)
// //     if err != nil {
// //         return nil, err
// //     }
// //     defer rows.Close() // Ensure rows are closed to avoid memory leaks
// //     var foodCategory []types.FoodCategory
// //
// //     for rows.Next() {
// //         var foodCat types.FoodCategory
// //         err = rows.Scan(
// //             &foodCat.IdCategory,
// //             &foodCat.Name,
// //         )
// //         if err != nil {
// //             return nil, err
// //         }
// //         foodCategory = append(foodCategory, foodCat)
// //     }
// //     if err := rows.Err(); err != nil {
// //         return nil, err
// //     }
// //     return &foodCategory, nil
// // }
// //!NOTE: Get all informations aboout the restaurant
func (s *store) GetRestaurantById(id string) (*types.Restaurant, error) {
	query := `SELECT * FROM restaurant WHERE idRestaurant = ?`
	row := s.db.QueryRow(query, id)
	var rest types.Restaurant
	err := row.Scan(
		&rest.IdRestaurant,
		&rest.IdAdminRestaurant,
		&rest.NameRestaurant,
		&rest.Image,
		&rest.Langitude,
		&rest.Latitude,
		&rest.Description,
		&rest.Capacity,
		&rest.Location,
	)
	if err != nil {
		return nil, err
	}
	return &rest, nil
}

// //!NOTE : GET all restaurant workers
//
//	func (s *store) GetRestaurantWorkers() (*[]types.RestaurantWorker, error) {
//		query := `SELECT * FROM restaurantWorker`
//		rows, err := s.db.Query(query)
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close() // Ensure rows are closed to avoid memory leaks
//		var restaurantWorker []types.RestaurantWorker
//
//		for rows.Next() {
//			var rest types.RestaurantWorker
//			err = rows.Scan(
//				&rest.IdRestaurantWorker,
//				&rest.FirstName,
//				&rest.LastName,
//				&rest.Address,
//				&rest.LastName,
//				&rest.Email,
//				&rest.PhoneNumber,
//				&rest.Quote,
//				&rest.StartWorking,
//				&rest.Nationnallity,
//				&rest.NativeLanguage,
//				&rest.Rating,
//				&rest.Status,
//				&rest.IdRestaurant,
//			)
//			if err != nil {
//				return nil, err
//			}
//			restaurantWorker = append(restaurantWorker, rest)
//		}
//		if err := rows.Err(); err != nil {
//			return nil, err
//		}
//		return &restaurantWorker, nil
//	}
//
// //!NOTE: WE  GONNA NEED THE RESTAURANTWORKER BY ID WHEN I GET THE RATING IG AND IT WILL BE CHANGED
// //!NOTE : menu by restaurant
//
//	func (s *store) getMenueByRestaurantId(id string) (*[]types.Menu, error) {
//		query := `SELECT * FROM menue WHERE idRestaurant = ?`
//		rows, err := s.db.Query(query, id)
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close() // Ensure rows are closed to avoid memory leaks
//		var menues []types.Menu
//
//		for rows.Next() {
//			var menue types.Menu
//
//			err = rows.Scan(
//				&menue.IdMenu,
//				&menue.IdRestaurant,
//				&menue.Name,
//				&menue.Active,
//				&menue.CreatedAt,
//			)
//			if err != nil {
//				return nil, err
//			}
//			menues = append(menues, menue)
//		}
//		if err := rows.Err(); err != nil {
//			return nil, err
//		}
//		return &menues, nil
//	}
//
// //!NOTE: GET THE FOOD BY THE MENU ID
func (s *store) GetFoodByMenu(idMenu string) (*[]types.Food, error) {
	query := `
        SELECT f.idFood, f.idCategory, f.name, f.description, f.image, f.price, f.status
        FROM food f
        JOIN menufood mf ON f.idFood = mf.idFood
        WHERE mf.idMenu = ?
    `
	rows, err := s.db.Query(query, idMenu)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var foods []types.Food
	for rows.Next() {
		var food types.Food
		if err := rows.Scan(&food.IdFood, &food.IdCategory, &food.Name, &food.Description, &food.Image, &food.Price, &food.Status); err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return &foods, nil
}

func (s *store) GetReservationTodayByRestaurantId(idRestaurant string) (*[]types.ReservationListInformation, error) {
	query := `SELECT profile.firstName,profile.lastName,profile.email,profile.address,reservation.numberOfPeople ,reservation.status FROM
    reservation join client on reservation.idClient=client.idClient
    join profile on profile.idProfile=client.idProfile
    WHERE idRestaurant = ? AND DATE(reservation.createdAt) = CURDATE()`
	rows, err := s.db.Query(query, idRestaurant)
	if err != nil {
		return nil, fmt.Errorf("error retrieving reservations: %v", err)
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var reservations []types.ReservationListInformation

	for rows.Next() {
		var reservation types.ReservationListInformation
		err = rows.Scan(
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Address,
			&reservation.NumberOfPeople,
			&reservation.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning reservation row: %v", err)
		}
		reservations = append(reservations, reservation)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over reservation rows: %v", err)
	}
	return &reservations, nil
}

func (s *store) GetOrderListForRestaurantToday(idRestaurant string) (*[]types.Order, error) {
	query := `SELECT * FROM orderList WHERE idRestaurant = ? AND DATE(createdAt) = CURDATE()`
	rows, err := s.db.Query(query, idRestaurant)
	if err != nil {
		return nil, fmt.Errorf("error retrieving orders: %v", err)
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var orders []types.Order
	for rows.Next() {
		var order types.Order
		err = rows.Scan(
			&order.IdOrder,
			&order.CreatedAt,
			&order.Status,
			&order.TotalPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order row: %v", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over order rows: %v", err)
	}
	return &orders, nil
}

func (s *store) CountReservationUpcomingWeek(idRestaurant string) (int, error) {
	query := `SELECT COUNT(*) FROM reservation WHERE idRestaurant = ? AND timeFrom >= CURDATE() AND timeFrom < DATE_ADD(CURDATE(), INTERVAL 7 DAY)`
	row := s.db.QueryRow(query, idRestaurant)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *store) CountReservationLastMonth(idRestaurant string) (*[]types.ReservationStats, error) {
	query := `
    SELECT 
    DATE(timeFrom) AS day,
    COUNT(*) AS reservations 
FROM reservation 
WHERE 
    idRestaurant = ? 
    AND MONTH(timeFrom) = MONTH(CURDATE()) 
    AND YEAR(timeFrom) = YEAR(CURDATE()) 
GROUP BY DATE(timeFrom)
ORDER BY day ASC;

    `
	rows, err := s.db.Query(query, idRestaurant)
	if err != nil {
		log.Printf("Error counting reservations last month for restaurant %s: %v", idRestaurant, err)
		return nil, fmt.Errorf("error counting reservations: %v", err)
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var reservations []types.ReservationStats
	for rows.Next() {
		var reservation types.ReservationStats
		err = rows.Scan(
			&reservation.Date,
			&reservation.NumberOfReservations,
		)
		if err != nil {
			log.Printf("Error scanning reservation stats row for restaurant %s: %v", idRestaurant, err)
			return nil, fmt.Errorf("error scanning reservation stats row: %v", err)
		}
		reservations = append(reservations, reservation)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over reservation stats rows for restaurant %s: %v", idRestaurant, err)
		return nil, fmt.Errorf("error iterating over reservation stats rows: %v", err)
	}
	return &reservations, nil
}

func (s *store) GetOrderListOfClientInRestaurant(idRestaurant string, idClient string) (*[]types.Order, error) {
	query := `SELECT orderList.* FROM orderList join reservation on orderList.idReservation = reservation.idReservation WHERE reservation.idRestaurant = ? AND idClient = ?`
	rows, err := s.db.Query(query, idRestaurant, idClient)
	if err != nil {
		log.Printf("Error retrieving orders for client %s in restaurant %s: %v", idClient, idRestaurant, err)
		return nil, fmt.Errorf("error retrieving orders: %v", err)
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var orders []types.Order
	for rows.Next() {
		var order types.Order
		err = rows.Scan(
			&order.IdOrder,
			&order.CreatedAt,
			&order.Status,
			&order.TotalPrice,
		)
		if err != nil {
			log.Printf("Error scanning order row for client %s in restaurant %s: %v", idClient, idRestaurant, err)
			return nil, fmt.Errorf("error scanning order row: %v", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over order rows for client %s in restaurant %s: %v", idClient, idRestaurant, err)
		return nil, fmt.Errorf("error iterating over order rows: %v", err)
	}
	return &orders, nil
}

// func (s *store) GetOrderInformation(idOrder string) (*[]interface{}, err) {
// 	// query := `Select * from orderList join reservation join uj
// 	query := `
//     select profile.* , reservation.* , client.* , food.* from client join profile on client.idProfile = profile.idProfile
//  join reservation on client.idClient = reservation.idClient
// join orderList on reservation.idReservation = orderList.idReservation
// join orderFood on orderList.idOrder = orderFood.idOrder
// join food on food.idFood = orderFood.idFood
// where orderList.idOrder=?
//     `
//     rows, err := s.db.Query(query, idOrder)
//     if err != nil {
//         return nil, fmt.Errorf("error retrieving order information: %v", err)
//     }
//     defer rows.Close() // Ensure rows are closed to avoid memory leaks
//     var orderInfo []interface{}
//     for rows.Next() {
//         var profile types.Profile
//         var reservation types.Reservation
//         var client types.Client
//         var food types.Food
//
//         err = rows.Scan(
//             &profile.IdProfile,
//             &profile.FirstName,
//             &profile.LastName,
//             &profile.Email,
//             &profile.PhoneNumber,
//             &profile.Image,
//             &reservation.IdReservation,
//             &reservation.IdClient,
//             &reservation.IdRestaurant,
//             &reservation.Status,
//             &reservation.Price,
//             &reservation.TimeReservation,
//             &reservation.CreatedAt,
//             &client.IdClient,
//             &client.Username,
//             &food.IdFood,
//             &food.Name,
//             &food.Description,
//             &food.Image,
//             &food.Price,
//         )
//         if err != nil {
//             return nil, fmt.Errorf("error scanning order information row: %v", err)
//         }
//         orderInfo = append(orderInfo, profile, reservation, client, food)
//     }
// }

// //NOTE: GET the reservations by the restaurants
//
//	func (s *store) getReservationByRestaurantId(id string) (*[]types.Reservation, error) {
//		query := `SELECT * FROM reservation WHERE idRestaurant = ?`
//		rows, err := s.db.Query(query, id)
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close() // Ensure rows are closed to avoid memory leaks
//		var reservations []types.Reservation
//
//		for rows.Next() {
//			var reservation types.Reservation
//
//			err = rows.Scan(
//				&reservation.IdReservation,
//				&reservation.IdClient,
//				&reservation.IdRestaurant,
//				&reservation.Status,
//				&reservation.Price,
//				&reservation.TimeReservation,
//				&reservation.CreatedAt,
//			)
//			if err != nil {
//				return nil, err
//			}
//			reservations = append(reservations, reservation)
//
//		}
//		if err := rows.Err(); err != nil {
//			return nil, err
//		}
//		return &reservations, nil
//	}
//
// // !NOTE: WE GONNE NEED ORDER LIST FOR ALL THE RESTAURANT , ORDER FOR EVERY RESERVATION , ORDER FOOD list FOR THE CLIENT , ORDER LIST FOOD
//
//	func (s *store) getOrderlistForRestaurant(idRestaurant string) (*[]types.Order, error) {
//		query := `select * from orderlist where idRestaurant = ?`
//
//		rows, err := s.db.Query(query, idRestaurant)
//		if err != nil {
//			return nil, err
//		}
//
//		defer rows.Close()
//
//		var orders []types.Order
//		for rows.Next() {
//			var order types.Order
//			err = rows.Scan(
//				&order.IdOrder,
//				&order.CreatedAt,
//				&order.Status,
//				&order.TotalPrice,
//			)
//
//			if err != nil {
//				return nil, err
//			}
//			orders = append(orders, order)
//		}
//		if err := rows.Err(); err != nil {
//			return nil, err
//		}
//		return &orders, nil
//	}
//
// //!NOTE: Get order for each client history
func (s *store) getOrderByClient(idClient string) (*[]types.Order, error) {
	query := `select * from orderlist join restaurant where idClient = ?`

	rows, err := s.db.Query(query, idClient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var order types.Order
		err = rows.Scan(
			&order.IdOrder,
			&order.CreatedAt,
			&order.Status,
			&order.TotalPrice,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &orders, nil
}

// //!NOTE :Food list ofr the client
// func(s *store) getFoodByOrder(idOrder string) (*[]types.Food , error) {
//     query := `select * from orderFood join food join foodCategory where idOrder = ?`
//
//     rows, err := s.db.Query(query, idOrder)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()
//     var foods []types.Food
//     for rows.Next() {
//         var food types.Food
//         err = rows.Scan(
//             &food.IdFood,
//             &food.IdCategory,
//             &food.IdMenu,
//             &food.Name,
//             &food.Description,
//             &food.Image,
//             &food.Price,
//             &food.Status,
//         )
//
//         if err != nil {
//             return nil, err
//         }
//         foods = append(foods, food)
//     }
//     return &foods, nil
// }
// //!NOTE: FEEDBACK restaurant worker
// func (s *store ) GetRestaurantWorkerFeedback (id string) (*[]types.RestaurantWorkerFeedBack , error) {
//     query := `select * from restaurantWorkerFeedBack where idRestaurantWorker = ?`
//
//     rows, err := s.db.Query(query, id)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()
//     var feedbacks []types.RestaurantWorkerFeedBack
//     for rows.Next() {
//         var feedback types.RestaurantWorkerFeedBack
//         err = rows.Scan(
//             &feedback.IdRestaurantWorkerFeedBack,
//             &feedback.IdRestaurantWorker,
//             &feedback.IdClient,
//             &feedback.Comment,
//             &feedback.CreatedAt,
//         )
//
//         if err != nil {
//             return nil, err
//         }
//         feedbacks = append(feedbacks, feedback)
//     }
//     return &feedbacks, nil
// }
//
//
//
// //!NOTE: RESERVATION PARTJ
// func (s *store) GetReservation () (*[]types.Reservation , error) {
//     query := `SELECT * FROM reservation`
//     rows, err := s.db.Query(query)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close() // Ensure rows are closed to avoid memory leaks
//     var reservation []types.Reservation
//
//     for rows.Next() {
//         var res types.Reservation
//         err = rows.Scan(
//             &res.IdReservation,
//             &res.IdClient,
//             &res.IdRestaurant,
//             &res.Status,
//             &res.Price,
//             &res.TimeReservation,
//             &res.CreatedAt,
//         )
//         if err != nil {
//             return nil, err
//         }
//         reservation = append(reservation, res)
//     }
//     if err := rows.Err(); err != nil {
//         return nil, err
//     }
//     return &reservation, nil
// }
//
// func (s *store) PostReservation(reservation types.Reservation) error {
//     query := `INSERT INTO reservation (idReservation, idClient, idRestaurant, status, price, timeReservation, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?)`
//     _, err := s.db.Exec(query, reservation.IdReservation, reservation.IdClient, reservation.IdRestaurant, reservation.Status, reservation.Price, reservation.TimeReservation, reservation.CreatedAt)
//     if err != nil {
//         return err
//     }
//     return nil
// }
//

func (s *store) GetFriendsOfClient(idClient string) (*[]string, error) {
	query := `select idClient2 from friendship where idClient1=? and status ="accepted"`
	rows, err := s.db.Query(query, idClient)
	if err != nil {
		return nil, fmt.Errorf("error retrieving friends: %v", err)
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	friends := []string{}
	for rows.Next() {
		var friendId string
		err = rows.Scan(&friendId)
		if err != nil {
			return nil, fmt.Errorf("error scanning friend row: %v", err)
		}
		friends = append(friends, friendId)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over friend rows: %v", err)
	}
	return &friends, nil
}

func convertToInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

func (s *store) PostRatingRestaurant(rating types.PostRatingRestaurant) error {
	query := `INSERT INTO rating (idRating, idClient, idRestaurant, ratingType, rating, comment, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, rating.IdRating, rating.IdClient, rating.IdRestaurant, "restaurant", rating.RatingValue, rating.Comment, time.Now())
	if err != nil {
		return fmt.Errorf("error inserting rating: %v", err)
	}
	return nil
}

func (s *store) GetRatingOfFriendsRestaurant(friendsId []string, idRestaurant string) (*[]types.RatingRestaurant, error) {
	if len(friendsId) == 0 {
		return &[]types.RatingRestaurant{}, nil
	}

	query := `
SELECT firstName, lastName, rating, comment, createdAt
FROM (
    SELECT 
        profile.firstName,
        profile.lastName,
        rating.rating,
        rating.comment,
        rating.createdAt,
        ROW_NUMBER() OVER (PARTITION BY rating.idClient ORDER BY rating.createdAt DESC) as rn
    FROM rating 
    JOIN client ON rating.idClient = client.idClient 
    JOIN profile ON client.idProfile = profile.idProfile 
    WHERE rating.idRestaurant = ? AND rating.idClient IN (`
	for i := range friendsId {
		if i > 0 {
			query += ", "
		}
		query += "?"
	}
	query += `)) AS ranked WHERE rn = 1`
	args := append([]interface{}{idRestaurant}, convertToInterfaceSlice(friendsId)...)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ratings: %v", err)
	}

	defer rows.Close()

	ratings := []types.RatingRestaurant{}
	for rows.Next() {
		var rating types.RatingRestaurant
		err = rows.Scan(
			&rating.FirstName,
			&rating.LastName,
			&rating.RatingValue,
			&rating.Comment,
			&rating.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rating row: %v", err)
		}
		ratings = append(ratings, rating)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rating rows: %v", err)
	}
	return &ratings, nil
}

func (s *store) GetRestaurantWorker(idRestaurant string) (*[]types.RestaurantWorker, error) {
	query := `SELECT idRestaurantWorker,firstName,lastName,email,phoneNumber,quote,startWorking,nationnallity ,nativeLanguage,rating ,image, address ,status FROM restaurantWorkers WHERE idRestaurant = ?`
	rows, err := s.db.Query(query, idRestaurant)
	if err != nil {
		return nil, fmt.Errorf("error retrieving restaurant workers: %v", err)
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var workers []types.RestaurantWorker
	for rows.Next() {
		var worker types.RestaurantWorker
		err = rows.Scan(
			&worker.IdRestaurantWorker,
			&worker.FirstName,
			&worker.LastName,
			&worker.Email,
			&worker.PhoneNumber,
			&worker.Quote,
			&worker.StartWorking,
			&worker.Nationnallity,
			&worker.NativeLanguage,
			&worker.Rating,
			&worker.Image,
			&worker.Address,
			&worker.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning restaurant worker row: %v", err)
		}
		workers = append(workers, worker)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over restaurant worker rows: %v", err)
	}
	return &workers, nil
}

func (s *store) GetRecentReviews(idRestaurant string) ([]*types.Rating, error) {
	query := `select
    rating.comment,rating.rating,rating.createdAt,profile.firstName,profile.lastName
    from rating join client on rating.idClient = client.idClient join profile
    on profile.idProfile = client.idProfile where idRestaurant=? and
    ratingType="restaurant" order by createdAt desc limit 5`
	rows, err := s.db.Query(query, idRestaurant)
	if err != nil {
		return nil, fmt.Errorf("error retrieving recent reviews: %v", err)
	}
	defer rows.Close()
	var reviews []*types.Rating
	for rows.Next() {
		var review types.Rating
		err = rows.Scan(
			&review.Comment,
			&review.RatingValue,
			&review.CreatedAt,
			&review.FirstName,
			&review.LastName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning recent review: %v", err)
		}
		reviews = append(reviews, &review)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over recent review rows: %v", err)
	}
	return reviews, nil
}

func (s *store) GetRecentOrders(idRestaurant string, limit int) ([]types.RecentOrder, error) {
	query := `SELECT orderList.idOrder, profile.firstName, profile.lastName, orderList.createdAt,client.idClient ,reservation.timeFrom,
        COUNT(orderFood.idFood) AS itemCount, orderList.totalPrice, orderList.status 
        FROM orderList 
        JOIN reservation ON orderList.idReservation = reservation.idReservation 
        JOIN client ON reservation.idClient = client.idClient 
        JOIN profile ON client.idProfile = profile.idProfile 
        JOIN orderFood ON orderList.idOrder = orderFood.idOrder 
        WHERE reservation.idRestaurant = ? AND reservation.status = 'confirmed'
        GROUP BY orderList.idOrder 
        ORDER BY orderList.createdAt DESC 
        LIMIT ?`

	rows, err := s.db.Query(query, idRestaurant, limit)
	if err != nil {
		return nil, fmt.Errorf("error retrieving recent orders: %v", err)
	}
	defer rows.Close()

	var recentOrders []types.RecentOrder
	for rows.Next() {
		var order types.RecentOrder
		if err := rows.Scan(&order.IdOrder, &order.FirstName, &order.LastName, &order.CreatedAt, &order.IdClient, &order.TimeFrom, &order.ItemCount, &order.TotalPrice, &order.Status); err != nil {
			return nil, fmt.Errorf("error scanning recent order: %v", err)
		}
		recentOrders = append(recentOrders, order)
	}

	return recentOrders, nil
}

func (s *store) GetOrderStatsByHourAndStatus(idRestaurant string) (map[int]int, map[string]int, error) {
	queryHour := `SELECT HOUR(orderList.createdAt) AS hour, COUNT(*) AS count FROM orderList join reservation on orderList.idReservation=reservation.idReservation WHERE reservation.idRestaurant = ? GROUP BY HOUR(createdAt)`
	queryStatus := `SELECT orderList.status, COUNT(*) AS count FROM orderList join reservation on orderList.idReservation=reservation.idReservation WHERE idRestaurant = ? GROUP BY status`

	rowsHour, err := s.db.Query(queryHour, idRestaurant)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving order stats by hour: %v", err)
	}
	defer rowsHour.Close()

	orderCountByHour := make(map[int]int)
	for rowsHour.Next() {
		var hour, count int
		if err := rowsHour.Scan(&hour, &count); err != nil {
			return nil, nil, fmt.Errorf("error scanning order stats by hour: %v", err)
		}
		orderCountByHour[hour] = count
	}

	// Order count by status
	rowsStatus, err := s.db.Query(queryStatus, idRestaurant)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving order stats by status: %v", err)
	}
	defer rowsStatus.Close()

	orderCountByStatus := make(map[string]int)
	for rowsStatus.Next() {
		var status string
		var count int
		if err := rowsStatus.Scan(&status, &count); err != nil {
			return nil, nil, fmt.Errorf("error scanning order stats by status: %v", err)
		}
		orderCountByStatus[status] = count
	}

	return orderCountByHour, orderCountByStatus, nil
}

func (s *store) GetClientReservationAndOrderDetails(idClient string) (*types.ClientDetails, error) {
	queryProfile := `SELECT profile.idProfile, profile.firstName, profile.lastName, profile.email, profile.phoneNumber, profile.address 
		FROM client 
		JOIN profile ON client.idProfile = profile.idProfile 
		WHERE client.idClient = ?`
	queryOrders := `SELECT orderList.idOrder, orderList.totalPrice, orderList.createdAt, orderList.status, food.name, food.price, orderFood.quantity 
		FROM orderList 
        join reservation ON orderList.idReservation = reservation.idReservation
		JOIN orderFood ON orderList.idOrder = orderFood.idOrder 
		JOIN food ON orderFood.idFood = food.idFood 
		WHERE reservation.idClient = ? 
        ORDER BY orderList.createdAt DESC`

	row := s.db.QueryRow(queryProfile, idClient)
	var profile types.Profile
	if err := row.Scan(&profile.IdProfile, &profile.FirstName, &profile.LastName, &profile.Email, &profile.Phone, &profile.Address); err != nil {
		return nil, fmt.Errorf("error retrieving client profile: %v", err)
	}

	rowsOrders, err := s.db.Query(queryOrders, idClient)
	if err != nil {
		return nil, fmt.Errorf("error retrieving orders: %v", err)
	}
	defer rowsOrders.Close()

	var orders []types.OrderDetails
	var totalSpent float64
	var totalOrders int
	var firstOrderDate *time.Time

	orderMap := make(map[string]*types.OrderDetails)
	var orderIDs []string

	for rowsOrders.Next() {
		var idOrder string
		var totalPrice float64
		var createdAt time.Time
		var status string
		var foodName string
		var foodPrice float64
		var quantity int

		if err := rowsOrders.Scan(&idOrder, &totalPrice, &createdAt, &status, &foodName, &foodPrice, &quantity); err != nil {
			return nil, fmt.Errorf("error scanning order details: %v", err)
		}

		if firstOrderDate == nil {
			firstOrderDate = &createdAt
		}

		if _, exists := orderMap[idOrder]; !exists {
			orderMap[idOrder] = &types.OrderDetails{
				IdOrder:    idOrder,
				TotalPrice: totalPrice,
				CreatedAt:  createdAt,
				Status:     status,
				FoodItems:  []types.FoodItemInformation{},
			}
			orderIDs = append(orderIDs, idOrder)
			totalOrders++
			if status == "completed" {
				totalSpent += totalPrice
			}
		}

		orderMap[idOrder].FoodItems = append(orderMap[idOrder].FoodItems, types.FoodItemInformation{
			Name:     foodName,
			Price:    foodPrice,
			Quantity: quantity,
		})
	}

	for _, id := range orderIDs {
		orders = append(orders, *orderMap[id])
	}

	return &types.ClientDetails{
		Profile:        profile,
		Orders:         orders,
		TotalOrders:    totalOrders,
		TotalSpent:     totalSpent,
		FirstOrderDate: firstOrderDate,
	}, nil
}

func (s *store) GetRestaurantRatingStats(idRestaurant string) (*types.RestaurantRatingStats, error) {
	// 1. Monthly stats (average, count by month)
	monthlyQuery := `
        SELECT 
            MONTH(createdAt) AS month,
            YEAR(createdAt) AS year,
            AVG(rating) AS averageRating,
            COUNT(*) AS totalRatings
        FROM rating
        WHERE idRestaurant = ?
        GROUP BY YEAR(createdAt), MONTH(createdAt)
        ORDER BY year, month
    `
	rows, err := s.db.Query(monthlyQuery, idRestaurant)
	if err != nil {
		return nil, fmt.Errorf("error retrieving restaurant monthly rating stats: %v", err)
	}
	defer rows.Close()

	var stats []types.MonthlyRatingStats

	for rows.Next() {
		var stat types.MonthlyRatingStats
		err = rows.Scan(
			&stat.Month,
			&stat.Year,
			&stat.AverageRating,
			&stat.TotalRatings,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning monthly rating stats row: %v", err)
		}
		stats = append(stats, stat)
	}

	// 2. Overall stats (global averages, percentages)
	overallQuery := `
        SELECT 
            AVG(rating) AS overallAverage,
            COUNT(*) AS totalRatings,
            SUM(CASE WHEN rating = 5 THEN 1 ELSE 0 END) AS count5Stars,
            SUM(CASE WHEN rating = 4 THEN 1 ELSE 0 END) AS count4Stars,
            SUM(CASE WHEN rating = 3 THEN 1 ELSE 0 END) AS count3Stars,
            SUM(CASE WHEN rating = 2 THEN 1 ELSE 0 END) AS count2Stars,
            SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END) AS count1Star
        FROM rating
        WHERE idRestaurant = ?
    `
	var overallAverage float64
	var totalRatings, count5Stars, count4Stars, count3Stars, count2Stars, count1Star int

	err = s.db.QueryRow(overallQuery, idRestaurant).Scan(
		&overallAverage,
		&totalRatings,
		&count5Stars,
		&count4Stars,
		&count3Stars,
		&count2Stars,
		&count1Star,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving restaurant overall rating stats: %v", err)
	}

	percent := func(count int) float64 {
		if totalRatings == 0 {
			return 0
		}
		return float64(count) * 100 / float64(totalRatings)
	}

	return &types.RestaurantRatingStats{
		MonthlyStats:     stats,
		OverallAverage:   overallAverage,
		TotalRatings:     totalRatings,
		Percentage5Stars: percent(count5Stars),
		Percentage4Stars: percent(count4Stars),
		Percentage3Stars: percent(count3Stars),
		Percentage2Stars: percent(count2Stars),
		Percentage1Star:  percent(count1Star),
	}, nil
}

func (s *store) GetReservationStatsAndList(idRestaurant string) (*types.ReservationStatsAndList, error) {
	// Stats query: use ROUND and NULLIF to avoid decimal scan issues and division by zero
	queryStats := `
		SELECT 
			(SELECT COUNT(*) FROM reservation WHERE DATE(createdAt) = CURDATE() AND idRestaurant = ?) AS totalToday,
			(SELECT COUNT(*) FROM reservation WHERE timeFrom > CURDATE() AND idRestaurant = ?) AS upcomingReservations,
			(SELECT IFNULL(ROUND((SELECT COUNT(*) FROM reservation WHERE status = 'confirmed' AND idRestaurant = ?) * 100.0 / NULLIF((SELECT COUNT(*) FROM reservation WHERE idRestaurant = ?), 0)),0)) AS confirmedRate
	`
	var totalToday, upcomingReservations int
	var confirmedRate float64
	err := s.db.QueryRow(queryStats, idRestaurant, idRestaurant, idRestaurant, idRestaurant).Scan(&totalToday, &upcomingReservations, &confirmedRate)
	if err != nil {
		return nil, fmt.Errorf("error retrieving reservation stats: %v", err)
	}

	// Query for today's reservations
	queryToday := `
		SELECT profile.firstName, profile.lastName, reservation.timeFrom, reservation.numberOfPeople
		FROM reservation
		JOIN client ON reservation.idClient = client.idClient
		JOIN profile ON client.idProfile = profile.idProfile
		WHERE DATE(reservation.createdAt) = CURDATE() AND reservation.idRestaurant = ?
	`
	rowsToday, err := s.db.Query(queryToday, idRestaurant)
	if err != nil {
		return nil, fmt.Errorf("error retrieving today's reservations: %v", err)
	}
	defer rowsToday.Close()

	var todayReservations []types.ReservationDetailsR
	for rowsToday.Next() {
		var reservation types.ReservationDetailsR
		err = rowsToday.Scan(&reservation.FirstName, &reservation.LastName, &reservation.TimeFrom, &reservation.NumberOfPeople)
		if err != nil {
			return nil, fmt.Errorf("error scanning today's reservation row: %v", err)
		}
		todayReservations = append(todayReservations, reservation)
	}

	// Query for upcoming reservations (limit 4)
	queryUpcoming := `
		SELECT profile.firstName, profile.lastName, reservation.timeFrom, reservation.numberOfPeople
		FROM reservation
		JOIN client ON reservation.idClient = client.idClient
		JOIN profile ON client.idProfile = profile.idProfile
		WHERE reservation.timeFrom > CURDATE() AND reservation.idRestaurant = ?
		ORDER BY reservation.timeFrom ASC
		LIMIT 4
	`
	rowsUpcoming, err := s.db.Query(queryUpcoming, idRestaurant)
	if err != nil {
		return nil, fmt.Errorf("error retrieving upcoming reservations: %v", err)
	}
	defer rowsUpcoming.Close()

	var upcomingReservationsList []types.ReservationDetailsR
	for rowsUpcoming.Next() {
		var reservation types.ReservationDetailsR
		err = rowsUpcoming.Scan(&reservation.FirstName, &reservation.LastName, &reservation.TimeFrom, &reservation.NumberOfPeople)
		if err != nil {
			return nil, fmt.Errorf("error scanning upcoming reservation row: %v", err)
		}
		upcomingReservationsList = append(upcomingReservationsList, reservation)
	}

	return &types.ReservationStatsAndList{
		TotalToday:           totalToday,
		UpcomingReservation:  upcomingReservations,
		ConfirmedRate:        confirmedRate,
		TodayReservations:    todayReservations,
		UpcomingReservations: upcomingReservationsList,
	}, nil
}

// GetAdminRestaurantStats retrieves aggregated statistics for all restaurants
func (s *store) GetAdminRestaurantStats() (*types.AdminRestaurantStats, error) {
	var stats types.AdminRestaurantStats

	// Get total restaurants count
	err := s.db.QueryRow("SELECT COUNT(*) FROM restaurant").Scan(&stats.TotalRestaurants)
	if err != nil {
		return nil, fmt.Errorf("error getting total restaurants count: %v", err)
	}

	// Get active restaurants count (restaurants with admin assigned)
	err = s.db.QueryRow("SELECT COUNT(*) FROM restaurant WHERE idAdminRestaurant IS NOT NULL").Scan(&stats.ActiveRestaurants)
	if err != nil {
		return nil, fmt.Errorf("error getting active restaurants count: %v", err)
	}

	// Get average rating across all restaurants
	err = s.db.QueryRow(`
		SELECT COALESCE(AVG(rating), 0) 
		FROM rating 
		WHERE idRestaurant IS NOT NULL
	`).Scan(&stats.AverageRating)
	if err != nil {
		return nil, fmt.Errorf("error getting average rating: %v", err)
	}

	// Get total bookings last month across all restaurants
	err = s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM reservation 
		WHERE createdAt >= DATE_SUB(CURDATE(), INTERVAL 1 MONTH) 
		AND createdAt < CURDATE()
	`).Scan(&stats.TotalBookingsLastMonth)
	if err != nil {
		return nil, fmt.Errorf("error getting total bookings last month: %v", err)
	}

	return &stats, nil
}

// GetAllRestaurantReviews retrieves all reviews for a specific restaurant
func (s *store) GetAllRestaurantReviews(idRestaurant string) ([]*types.Rating, error) {
	query := `
		SELECT 
			r.idRating,
			r.rating,
			r.comment,
			r.createdAt,
			p.firstName,
			p.lastName
		FROM rating r
		JOIN profile p ON r.idClient = p.idProfile
		WHERE r.idRestaurant = ?
		ORDER BY r.createdAt DESC
	`

	rows, err := s.db.Query(query, idRestaurant)
	if err != nil {
		return nil, fmt.Errorf("error querying all restaurant reviews: %v", err)
	}
	defer rows.Close()

	var reviews []*types.Rating
	for rows.Next() {
		var review types.Rating
		
		err := rows.Scan(
			&review.IdRating,
			&review.RatingValue,
			&review.Comment,
			&review.CreatedAt,
			&review.FirstName,
			&review.LastName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning review row: %v", err)
		}

		reviews = append(reviews, &review)
	}

	return reviews, nil
}

// GetRestaurantTodaySummary retrieves today's summary for a specific restaurant
func (s *store) GetRestaurantTodaySummary(idRestaurant string) (*types.RestaurantTodaySummary, error) {
	var summary types.RestaurantTodaySummary

	// Get total reservations today
	err := s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM reservation 
		WHERE idRestaurant = ? 
		AND DATE(timeFrom) = CURDATE()
	`, idRestaurant).Scan(&summary.TotalReservationsToday)
	if err != nil {
		return nil, fmt.Errorf("error getting total reservations today: %v", err)
	}

	// Get confirmed reservations today
	err = s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM reservation 
		WHERE idRestaurant = ? 
		AND DATE(timeFrom) = CURDATE() 
		AND status = 'confirmed'
	`, idRestaurant).Scan(&summary.ConfirmedReservations)
	if err != nil {
		return nil, fmt.Errorf("error getting confirmed reservations: %v", err)
	}

	// Get pending reservations today
	err = s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM reservation 
		WHERE idRestaurant = ? 
		AND DATE(timeFrom) = CURDATE() 
		AND status = 'pending'
	`, idRestaurant).Scan(&summary.PendingReservations)
	if err != nil {
		return nil, fmt.Errorf("error getting pending reservations: %v", err)
	}

	// Calculate current occupancy (confirmed reservations / restaurant capacity * 100)
	var capacity int
	err = s.db.QueryRow("SELECT capacity FROM restaurant WHERE idRestaurant = ?", idRestaurant).Scan(&capacity)
	if err != nil {
		return nil, fmt.Errorf("error getting restaurant capacity: %v", err)
	}

	if capacity > 0 {
		summary.CurrentOccupancy = float64(summary.ConfirmedReservations) / float64(capacity) * 100
	}

	return &summary, nil
}
