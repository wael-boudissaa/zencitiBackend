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

func (s *store) PostOrderList(orderId string, foods []types.FoodItem) error {
	var totalPrice float64
	for _, food := range foods {
		_, err := s.db.Exec(`Insert INTO orderFood (idOrder, idFood, quantity, createdAt) VALUES (?, ?, ?, ?)`, orderId, food.IdFood, food.Quantity, time.Now())
		totalPrice += food.PriceSingle * float64(food.Quantity)
		if err != nil {
			return err
		}

	}
	query := `UPDATE orderList SET totalPrice = ? WHERE idOrder = ?`
	_, err := s.db.Exec(query, totalPrice, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) CreateReservation(idReservation string, reservation types.ReservationCreation) error {
	query := `INSERT INTO reservation (idReservation, idClient, idRestaurant, idTable, status, createdAt, numberOfPeople, timeFrom, timeTo) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, idReservation, reservation.IdClient, reservation.IdRestaurant, reservation.TableId, "pending", time.Now(), reservation.NumberOfPeople, reservation.TimeFrom, time.Now().Add(time.Hour*2))
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

//
// //NOTE: GET the reservations by the restaurants
// func (s *store) getReservationByRestaurantId(id string) (*[]types.Reservation, error) {
// 	query := `SELECT * FROM reservation WHERE idRestaurant = ?`
// 	rows, err := s.db.Query(query, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close() // Ensure rows are closed to avoid memory leaks
// 	var reservations []types.Reservation
//
// 	for rows.Next() {
// 		var reservation types.Reservation
//
// 		err = rows.Scan(
// 			&reservation.IdReservation,
// 			&reservation.IdClient,
// 			&reservation.IdRestaurant,
// 			&reservation.Status,
// 			&reservation.Price,
// 			&reservation.TimeReservation,
// 			&reservation.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		reservations = append(reservations, reservation)
//
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return &reservations, nil
// }
//
// // !NOTE: WE GONNE NEED ORDER LIST FOR ALL THE RESTAURANT , ORDER FOR EVERY RESERVATION , ORDER FOOD list FOR THE CLIENT , ORDER LIST FOOD
// func (s *store) getOrderlistForRestaurant(idRestaurant string) (*[]types.Order, error) {
// 	query := `select * from orderlist where idRestaurant = ?`
//
// 	rows, err := s.db.Query(query, idRestaurant)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	defer rows.Close()
//
// 	var orders []types.Order
// 	for rows.Next() {
// 		var order types.Order
// 		err = rows.Scan(
// 			&order.IdOrder,
// 			&order.CreatedAt,
// 			&order.Status,
// 			&order.TotalPrice,
// 		)
//
// 		if err != nil {
// 			return nil, err
// 		}
// 		orders = append(orders, order)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return &orders, nil
// }
//
// //!NOTE: Get order for each client history
// func (s *store) getOrderByClient(idClient string) (*[]types.Order, error) {
// 	query := `select * from orderlist join restaurant where idClient = ?`
//
// 	rows, err := s.db.Query(query, idClient)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	var orders []types.Order
// 	for rows.Next() {
// 		var order types.Order
// 		err = rows.Scan(
// 			&order.IdOrder,
// 			&order.CreatedAt,
// 			&order.Status,
// 			&order.TotalPrice,
// 		)
//
// 		if err != nil {
// 			return nil, err
// 		}
// 		orders = append(orders, order)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return &orders, nil
// }
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
