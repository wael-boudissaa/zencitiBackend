package restaurant

import (
	"database/sql"

	"github.com/wael-boudissaa/zencitiBackend/types"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{db: db}
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
			&rest.NameRestaurant,
			&rest.Description,
			&rest.IdActivite,
			&rest.Image,
			&rest.Location,
			&rest.Capacity,
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

func (s *store) GetRestaurantById(id string) (*types.Restaurant, error) {
	query := `SELECT * FROM restaurant WHERE idRestaurant = ?`
	row := s.db.QueryRow(query, id)
	var rest types.Restaurant
	err := row.Scan(
		&rest.IdRestaurant,
		&rest.NameRestaurant,
		&rest.Description,
		&rest.IdActivite,
		&rest.Image,
		&rest.Location,
		&rest.Capacity,
	)
	if err != nil {
		return nil, err
	}
	return &rest, nil
}

func (s *store) GetRestaurantWorkers() (*[]types.RestaurantWorker, error) {
	query := `SELECT * FROM restaurantWorker`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var restaurantWorker []types.RestaurantWorker

	for rows.Next() {
		var rest types.RestaurantWorker
		err = rows.Scan(
			&rest.IdRestaurantWorker,
			&rest.FirstName,
			&rest.LastName,
			&rest.Address,
			&rest.LastName,
			&rest.Email,
			&rest.PhoneNumber,
			&rest.Quote,
			&rest.StartWorking,
			&rest.Nationnallity,
			&rest.NativeLanguage,
			&rest.Rating,
			&rest.Status,
			&rest.IdRestaurant,
		)
		if err != nil {
			return nil, err
		}
		restaurantWorker = append(restaurantWorker, rest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &restaurantWorker, nil
}

//!NOTE: WE  GONNA NEED THE RESTAURANTWORKER BY ID WHEN I GET THE RATING IG AND IT WILL BE CHANGED

func (s *store) getMenueByRestaurantId(id string) (*[]types.Menu, error) {
	query := `SELECT * FROM menue WHERE idRestaurant = ?`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var menues []types.Menu

	for rows.Next() {
		var menue types.Menu

		err = rows.Scan(
			&menue.IdMenu,
			&menue.IdRestaurant,
			&menue.Name,
			&menue.Active,
			&menue.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		menues = append(menues, menue)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &menues, nil
}

func (s *store) getFoodByMenuId(id string) (*[]types.Food, error) {
	query := `SELECT * FROM food WHERE idMenu = ?`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
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

func (s *store) getReservationByRestaurantId(id string) (*[]types.Reservation, error) {
	query := `SELECT * FROM reservation WHERE idRestaurant = ?`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var reservations []types.Reservation

	for rows.Next() {
		var reservation types.Reservation

		err = rows.Scan(
			&reservation.IdReservation,
			&reservation.IdClient,
			&reservation.IdRestaurant,
			&reservation.Status,
			&reservation.Price,
			&reservation.TimeReservation,
			&reservation.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &reservations, nil
}

// !NOTE: WE GONNE NEED ORDER LIST FOR ALL THE RESTAURANT , ORDER FOR EVERY RESERVATION , ORDER FOOD list FOR THE CLIENT , ORDER LIST FOOD
func (s *store) getOrderlistForRestaurant(idRestaurant string) (*[]types.Order, error) {
	query := `select * from orderlist where idRestaurant = ?`

	rows, err := s.db.Query(query, idRestaurant)
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

func (s *store) getFoodByOrder(idOrder string) (*[]types.Food, error) {
	query := `select * from orderFood join food join categoryFood where idOrder = ?`

	rows, err := s.db.Query(query, idOrder)
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
	return &foods, nil
}
