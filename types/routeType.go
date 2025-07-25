package types

import "time"

// !TODO: PHONE NUMBER AS LONGIN INFORMATION
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


type ActivityProfile struct {
    IdActivity    string    `json:"idActivity"`
    NameActivity  string    `json:"nameActivity"`
    Description   string    `json:"descriptionActivity"`
    ImageActivite string    `json:"imageActivity"`
    Capacity      int       `json:"capacity"`
    TimeActivity  time.Time `json:"timeActivity"`
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

type RegisterAdmin struct{ 
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Address   string `json:"address"`
	Type      string `json:"type"`
	Phone     string `json:"phone_number"`
    IdActivity string `json:"idActivity,omitempty"`    // For adminActivity
    IdRestaurant string `json:"idRestaurant,omitempty"` // For adminRestaurant
}

type ReservationCreation struct {
	IdClient       string     `json:"idClient"`
	IdRestaurant   string     `json:"idRestaurant"`
	NumberOfPeople int        `json:"numberOfPeople"`
	TimeFrom       time.Time  `json:"timeFrom"`
	TableId        string     `json:"idTable"`
}

type GetRestaurantTable struct {
	IdRestaurant string    `json:"idRestaurant"`
	TimeSlot     time.Time `json:"timeSlot"`
}

type OrderCreation struct {
	IdReservation string     `json:"idReservation"`
	Foods         []FoodItem `json:"food"`
}

type FoodItem struct {
	IdFood      string  `json:"idFood"`
	PriceSingle float64 `json:"priceSingle"`
    Quantity    int     `json:"quantity"`
}
type FoodItemInformation struct {
Name        string  `json:"name"`
	Price float64 `json:"priceSingle"`
    Quantity    int     `json:"quantity"`
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


type FriendsReviewsRestaruant struct{
    IdClient string `json:"idClient"`
    IdRestaurant string `json:"idRestaurant"`

}
type TimeNotAvaialable struct { 
    IdActivity    string    `json:"idActivity"`
    Day string    `json:"day"`
}

type ActivityCreation struct { 
    IdClient      string    `json:"idClient"`
    IdActivity    string    `json:"idActivity"`
    TimeActivity      time.Time `json:"timeActivity"`

}
