package restaurant

import (
	"database/sql"
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

// !NOTE: GET all restaurant

func (s *store) PostOrderList(order types.OrderFinalization) error {
	var totalPrice float64
	for _, food := range order.Foods {
		_, err := s.db.Exec(`Insert INTO orderFood (idOrder, idFood, quantity, createdAt) VALUES (?, ?, ?, ?)`, order.IdOrder, food.IdFood, food.Quantity, time.Now())
		totalPrice += food.PriceSingle * float64(food.Quantity)
		if err != nil {
			return err
		}

	}
	query := `UPDATE orderList SET totalPrice = ? WHERE idOrder = ?`
	_, err := s.db.Exec(query, totalPrice, order.IdOrder)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) CreateReservation(idReservation string, reservation types.ReservationCreation) error {
	query := `INSERT INTO reservation (idReservation, idClient, idRestaurant, status, price, timeReservation, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, idReservation, reservation.IdClient, reservation.IdRestaurant, "pending", 0, reservation.TimeSlot, time.Now())
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

func (s *store) GetRestaurantTables(restaurantId string) (*[]types.RestaurantTable, error) {
	query := `SELECT * FROM table_restaurant WHERE idRestaurant = ?`
	rows, err := s.db.Query(query, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var tables []types.RestaurantTable
	for rows.Next() {
		var table types.RestaurantTable
		err = rows.Scan(
			&table.IdTable,
			&table.IdRestaurant,
			&table.ReservationTime,
			&table.PosX,
			&table.PosY,
			&table.Duration_minutes,
			&table.Is_available,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
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
