package restaurant

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// "log"

	"github.com/wael-boudissaa/zencitiBackend/types"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{db: db}
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
SELECT food.*,menu.name as menuName
 FROM menu
JOIN food ON food.idMenu = menu.idMenu
where menu.active = 1 and menu.idRestaurant = ?;
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
			&menu.IdCategory,
			&menu.IdMenu,
			&menu.Name,
			&menu.Description,
			&menu.Image,
			&menu.Price,
			&menu.Status,
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
			status, createdAt, numberOfPeople, timeFrom, timeTo
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = s.db.Exec(query,
		idReservation,
		reservation.IdClient,
		reservation.IdRestaurant,
		reservation.TableId,
		"pending",
		time.Now(),
		reservation.NumberOfPeople,
		reservation.TimeFrom,
		reservation.TimeFrom.Add(time.Hour*2),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) ReserveTable(idReservation string, reservation types.ReservationCreation) error {
	query := `INSERT INTO table_reservation (idTable, idReservation, numberOfPeople, timeFrom, timeTo) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, reservation.TableId, idReservation, reservation.NumberOfPeople, reservation.TimeFrom, reservation.TimeTo)
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

func (s *store) GetRestaurantTables(restaurantId string, timeReserved time.Time) (*[]types.RestaurantTableStatus, error) {
	query := `SELECT tr.idTable, tr.idRestaurant, r.idReservation, tr.posX, tr.posY, r.timeFrom, r.timeTo, r.numberOfPeople,
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
			&table.IdReservation,
			&table.PosX,
			&table.PosY,
			&table.TimeFrom,
			&table.TimeTo,
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
	query := `SELECT * FROM restaurant`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var restaurant []types.Restaurant

	for rows.Next() {
		var rest types.Restaurant

		err = rows.Scan(
			&rest.IdRestaurant,
			&rest.IdAdminRestaurant,
			&rest.NameRestaurant,
			&rest.Image,
			&rest.Description,
			&rest.Capacity,
			&rest.Location,
		)
		if err != nil {
			return nil, err
		}
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
func (s *store) GetFoodByMenu(id string) (*[]types.Food, error) {
	query := `SELECT food.* 
          FROM food 
             JOIN menu ON food.idMenu = menu.idMenu 
                WHERE menu.idMenu = ?`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var foods []types.Food
	for rows.Next() {
		var food types.Food

		err = rows.Scan(
			&food.IdFood,
			&food.IdCategory,
			&food.IdMenu,
			&food.Name,
			&food.Description,
			&food.Image,
			&food.Price,
			&food.Status,
		)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	if err := rows.Err(); err != nil {
		return nil, err
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
	query := `SELECT idRestaurantWorker,firstName,lastName,email,phoneNumber,quote,startWorking,nationnallity ,nativeLanguage,rating , address ,status FROM restaurantWorkers WHERE idRestaurant = ?`
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
		WHERE reservation.idRestaurant = ? 
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
			totalSpent += totalPrice
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
	var totalToday, upcomingReservations, confirmedRate int
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
