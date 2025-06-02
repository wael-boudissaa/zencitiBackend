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
    UserName string `json:"username"`
}

type ReservationCreation struct {
	IdClient       string     `json:"idClient"`
	IdRestaurant   string     `json:"idRestaurant"`
	NumberOfPeople int        `json:"numberOfPeople"`
	TimeFrom       time.Time  `json:"timeFrom"`
	TimeTo         *time.Time `json:"timeTo"`
	TableId        string     `json:"idTable"`
}

type GetRestaurantTable struct {
	IdRestaurant string    `json:"idRestaurant"`
	TimeSlot     time.Time `json:"timeSlot"`
}

type OrderCreation struct {
	IdReservation string     `json:"idReservation"`
	Foods         []FoodItem `json:"foods"`
}

type FoodItem struct {
	IdFood      string  `json:"idFood"`
	Quantity    int     `json:"quantity"`
	PriceSingle float64 `json:"priceSingle"`
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

type SendRequestFriend struct {
	FromClient string `json:"from_client"`
	ToClient   string `json:"to_client"`
}

type AcceptFriendRequest struct {
    IdFriendship string `json:"idFriendship"`
}
