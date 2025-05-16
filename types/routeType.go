package types

import "time"

// !TODO: PHONE NUMBER AS LONGIN INFORMATION
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Role      string `json:"role"`
	Address   string `json:"address"`
	Type      string `json:"type"`
	Phone     string `json:"phone_number"`
}

type ReservationCreation struct {
	IdClient     string    `json:"idClient"`
	IdRestaurant string    `json:"idRestaurant"`
	TimeSlot     time.Time `json:"timeSlot"`
}

type OrderCreation struct {
	IdReservation string `json:"idReservation"`
}

type FoodItem struct {
	IdFood   string `json:"idFood"`
	Quantity int    `json:"quantity"`
}

type OrderFinalization struct {
	IdOrder string     `json:"idOrder"`
	Foods   []FoodItem `json:"foods"`
}



type AddFoodToOrder struct {
	IdOrder  string `json:"idOrder"`
	IdFood   string `json:"idFood"`
	Quantity int    `json:"quantity"`
}

type GetStatusTables struct {
	RestaurantId string `json:"restaurantId"`
	TimeSlot     string `json:"timeSlot"`
}
type RequestCreate struct {
	ClientId string `json:"client_id"`
	// Status string `json:"status"`
	//!NOTE: FOR NOW THE PRESTATIRE SHOULD BE FOUND IN THE SAME ROUTE FUNCTION THE FRONT SHOULD SEND ME ONLY THE CLIENT ID I THINK FOR NOW
	// PrestataireId string `json:"prestataire_id"`
	// Price float64 `json:"price"`
}

type ServicesAssignPrestataire struct {
	Services []string `json:"services"`
}

type RequestConfirmationRoute struct {
	ClientId      string `json:"client_id"`
	PrestataireId string `json:"prestataire_id"`
}
