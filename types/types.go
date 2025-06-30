package types

import "time"

type OrderInformation struct {
	IdOrder         string          `json:"idOrder"`
	TotalPrice      float64         `json:"totalPrice"`
	Status          string          `json:"status"`
	CreatedAt       time.Time       `json:"createdAt"`
	ClientFirstName string          `json:"clientFirstName"`
	ClientLastName  string          `json:"clientLastName"`
	ClientEmail     string          `json:"clientEmail"`
	ClientPhone     string          `json:"clientPhone"`
	ClientAddress   string          `json:"clientAddress"`
	ClientUsername  string          `json:"clientUsername"`
	ReservationTime time.Time       `json:"reservationTime"`
	NumberOfPeople  int             `json:"numberOfPeople"`
	FoodItems       []OrderFoodItem `json:"foodItems"`
}

type OrderFoodItem struct {
	IdFood      string  `json:"idFood"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}
type RestaurantReservationDetail struct {
	IdReservation  string    `json:"idReservation"`
	TimeFrom       time.Time `json:"timeFrom"`
	FullName       string    `json:"fullName"`
	TableId        string    `json:"tableId"`
	NumberOfPeople int       `json:"numberOfPeople"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}
type ReservationIdDetails struct {
	IdReservation   string               `json:"idReservation"`
	TimeFrom        time.Time            `json:"timeFrom"`
	NumberOfPeople  int                  `json:"numberOfPeople"`
	Status          string               `json:"status"`
	CreatedAt       time.Time            `json:"createdAt"`
	FullName        string               `json:"fullName"`
	FirstName       string               `json:"firstName"`
	LastName        string               `json:"lastName"`
	Email           string               `json:"email"`
	PhoneNumber     string               `json:"phoneNumber"`
	TotalVisits     int                  `json:"totalVisits"`
	AverageSpending float64              `json:"averageSpending"`
	TotalSpent      float64              `json:"totalSpent"`
	FavoriteFood    string               `json:"favoriteFood"`
	TotalOrders     int                  `json:"totalOrders"`
	Orders          []ClientOrderSummary `json:"orders"`
}

type ClientOrderSummary struct {
	IdOrder    string    `json:"idOrder"`
	TotalPrice float64   `json:"totalPrice"`
	CreatedAt  time.Time `json:"createdAt"`
	ItemCount  int       `json:"itemCount"`
	Status     string    `json:"status"`
}

type PaginatedReservations struct {
	Reservations []RestaurantReservationDetail `json:"reservations"`
	CurrentPage  int                           `json:"currentPage"`
	TotalPages   int                           `json:"totalPages"`
	TotalCount   int                           `json:"totalCount"`
	HasNext      bool                          `json:"hasNext"`
	HasPrevious  bool                          `json:"hasPrevious"`
}
type Profile struct {
	IdProfile string `json:"idProfile"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
}
type ClientReservationInfo struct {
	IdReservation      string    `json:"idReservation"`
	TimeFrom           time.Time `json:"timeFrom"`
	NumberOfPeople     int       `json:"numberOfPeople"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"createdAt"`
	RestaurantName     string    `json:"restaurantName"`
	RestaurantImage    string    `json:"restaurantImage"`
	RestaurantLocation string    `json:"restaurantLocation"`
	IdRestaurant       string    `json:"idRestaurant"`
}
type ClientActivityInfo struct {
	IdClientActivity    string    `json:"idClientActivity"`
	TimeActivity        time.Time `json:"timeActivity"`
	Status              string    `json:"status"`
	ActivityName        string    `json:"activityName"`
	ActivityImage       string    `json:"activityImage"`
	ActivityDescription string    `json:"activityDescription"`
}
type ClientInfo struct {
	IdClient        string `json:"idClient"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	IsAdminActivity bool   `json:"isAdminActivity"`
}

type UserInformation struct {
	IdClient  string `json:"idClient"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Following int    `json:"following"`
	Followers int    `json:"followers"`
	Username  string `json:"username"`
}

type TableOccupation struct {
	IdTable   string   `json:"idTable"`
	Occupied  bool     `json:"occupied"`
	TimeSlots []string `json:"timeSlots"`
}
type FoodPopularity struct {
	IdFood string `json:"idFood"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Total  int    `json:"total"`
}
type ProfilePage struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Following int    `json:"following"`
	Followers int    `json:"followers"`
	Username  string `json:"username"`
}

type Activity struct {
	IdActivity   string   `json:"idActivity"`
	NameActivity string   `json:"nameActivity"`
	Description  string   `json:"descriptionActivity"`
	Langitude    *float64 `json:"longitude" db:"longitude"`
	Latitude     *float64 `json:"latitude" db:"latitude"`

	IdTypeActivity string `json:"idTypeActivity"`
	ImageActivite  string `json:"imageActivity"`
	Popularity     int    `json:"popularity"`
}

type ActivitetType struct {
	IdActiviteType   string `json:"idTypeActivity"`
	NameActiviteType string `json:"nameTypeActivity"`
	ImageActivity    string `json:"imageActivity"`
}
type Restaurant struct {
	IdRestaurant      *string  `json:"idRestaurant" db:"idRestaurant"`
	IdAdminRestaurant *string  `json:"idAdminRestaurant" db:"idAdminRestaurant"`
	NameRestaurant    *string  `json:"name" db:"name"`
	Description       *string  `json:"description" db:"description"`
	Langitude         *float64 `json:"longitude" db:"longitude"`
	Latitude          *float64 `json:"latitude" db:"latitude"`
	Image             *string  `json:"image" db:"image"`
	Location          *string  `json:"location" db:"location"`
	Capacity          *int     `json:"capacity" db:"capacity"`
}

type RestaurantTable struct {
	IdTable          string    `json:"idTable"`
	IdRestaurant     string    `json:"idRestaurant"`
	ReservationTime  time.Time `json:"reservation_time"`
	PosX             int       `json:"posX"` // Position on UI map[jko]type
	PosY             int       `json:"posY"`
	Duration_minutes int       `json:"duration_minutes"`
	Is_available     bool      `json:"is_available"`
}

type RestaurantTableStatus struct {
	IdTable        *string    `json:"idTable"`
	Shape          *string    `json:"shape"` // New field
	PosX           *int       `json:"posX"`
	PosY           *int       `json:"posY"`
	IdRestaurant   *string    `json:"idRestaurant"`
	IdReservation  *string    `json:"idReservation"`
	NumberOfPeople *int       `json:"numberOfPeople"`
	TimeFrom       *time.Time `json:"timeFrom"`
	Status         *string    `json:"status"`
}
type FoodCategory struct {
	IdCategory    string `json:"idCategory"`
	NameCategorie string `json:"nameCategorie"`
}
type RestaurantWorker struct {
	IdRestaurantWorker string  `json:"idRestaurantWorker"`
	IdRestaurant       string  `json:"idRestaurant"`
	FirstName          string  `json:"firstName"`
	LastName           string  `json:"lastName"`
	Image              *string `json:"image"`
	Email              string  `json:"email"`
	PhoneNumber        string  `json:"phoneNumber"`
	Quote              string  `json:"quote"`
	StartWorking       string  `json:"startWorking"`
	Nationnallity      string  `json:"nationnallity"`
	NativeLanguage     string  `json:"nativeLanguage"`
	Rating             float32 `json:"rating"`
	Address            string  `json:"address"`
	Status             string  `json:"status"`
}
type RestaurantWorkerCreation struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	Quote          string `json:"quote"`
	Image          string `json:"image"`
	Nationnallity  string `json:"nationnallity"`
	NativeLanguage string `json:"nativeLanguage"`
	Address        string `json:"address"`
}

type RestaurantWorkerFeedBack struct {
	IdRestaurantWorkerFeedBack string    `json:"idRestaurantWorkerFeedBack"`
	IdRestaurantWorker         string    `json:"idRestaurantWorker"`
	IdClient                   string    `json:"idClient"`
	Comment                    string    `json:"comment"`
	CreatedAt                  time.Time `json:"createdAt"`
}

type Friendship struct {
	IdFriendship string    `json:"idFriendship"`
	Username     string    `json:"username"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}
type Menu struct {
	IdMenu       string    `json:"idMenu"`
	IdRestaurant string    `json:"idRestaurant"`
	Name         string    `json:"name"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
}
type LayoutCreation struct {
	Shape string `json:"shape"` // New field
	PosX  int    `json:"posX"`
	PosY  int    `json:"posY"`
}
type RestaurantWorkerWithRatings struct {
	IdRestaurantWorker string            `json:"idRestaurantWorker"`
	FirstName          string            `json:"firstName"`
	LastName           string            `json:"lastName"`
	Email              string            `json:"email"`
	PhoneNumber        string            `json:"phoneNumber"`
	Quote              string            `json:"quote"`
	StartWorking       time.Time         `json:"startWorking"`
	Nationnallity      string            `json:"nationnallity"`
	NativeLanguage     string            `json:"nativeLanguage"`
	Rating             float64           `json:"rating"`
	Image              *string           `json:"image"`
	Address            string            `json:"address"`
	Status             string            `json:"status"`
	IdRestaurant       string            `json:"idRestaurant"`
	RecentRatings      []WorkerRating    `json:"recentRatings"`
	RatingStats        WorkerRatingStats `json:"ratingStats"`
}

type WorkerRating struct {
	RatingValue     int       `json:"ratingValue"`
	Comment         string    `json:"comment"`
	CreatedAt       time.Time `json:"createdAt"`
	ClientFirstName string    `json:"clientFirstName"`
	ClientLastName  string    `json:"clientLastName"`
}

type WorkerRatingStats struct {
	TotalRatings     int     `json:"totalRatings"`
	AverageRating    float64 `json:"averageRating"`
	Percentage5Stars float64 `json:"percentage5Stars"`
	Percentage4Stars float64 `json:"percentage4Stars"`
	Percentage3Stars float64 `json:"percentage3Stars"`
	Percentage2Stars float64 `json:"percentage2Stars"`
	Percentage1Star  float64 `json:"percentage1Star"`
}

type Table struct {
	IdTable      string `json:"idTable"`
	IdRestaurant string `json:"idRestaurant"`
	Shape        string `json:"shape"` // New field
	PosX         int    `json:"posX"`
	PosY         int    `json:"posY"`
	IsAvailable  bool   `json:"is_available"`
}
type MenuInformationFood struct {
	IdMenu       string  `json:"idMenu" db:"idMenu"`
	IdFood       string  `json:"idFood" db:"idFood"`
	IdCategory   string  `json:"idCategory" db:"idCategory"`
	Name         string  `json:"name" db:"name"`
	Description  *string `json:"description" db:"description"`
	IdRestaurant string  `json:"idRestaurant" db:"idRestaurant"`
	Image        *string `json:"image" db:"image"`
	Price        float64 `json:"price" db:"price"`
	Status       string  `json:"status" db:"status"`
	MenuName     string  `json:"menuName" db:"menuName"`
}
type Food struct {
	IdFood      string   `json:"idFood"`
	IdCategory  string   `json:"idCategory"`
	IdMenu      *string  `json:"idMenu"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Image       *string  `json:"image"`
	Price       *float64 `json:"price"`
	Status      *string  `json:"status"`
}
type RestaurantMenuStats struct {
	TotalMenus       int           `json:"totalMenus"`
	ActiveMenuName   string        `json:"activeMenuName"`
	TotalItems       int           `json:"totalItems"`
	TotalCategories  int           `json:"totalCategories"`
	AvailableFoods   int           `json:"availableFoods"`
	UnavailableFoods int           `json:"unavailableFoods"`
	PopularFoods     []PopularFood `json:"popularFoods"`
}

type PopularFood struct {
	FoodName   string `json:"foodName"`
	OrderCount int    `json:"orderCount"`
}

// !WARNING:: THERE SHOULD BE A GENEARL THING ON THE RESERVATION FOR THE RESTAURANT AND THE ACITIVITE AND ALSO FOR THE RATING AND FEEDBACK
type ReservationListInformation struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	NumberOfPeople int    `json:"numberOfPeople"`
	Address        string `json:"address"`
	Status         string `json:"status"`
}

type ReservationStats struct {
	Date                 time.Time `json:"day"`
	NumberOfReservations int       `json:"reservations"`
}

type OrderDetails struct {
	IdOrder    string                `json:"idOrder"`
	CreatedAt  time.Time             `json:"createdAt"`
	Status     string                `json:"status"`
	FoodItems  []FoodItemInformation `json:"foodItems"`
	TotalPrice float64               `json:"totalPrice"`
}
type RecentOrder struct {
	IdOrder    string    `json:"idOrder"`
	FirstName  string    `json:"firstName"`
	IdClient   string    `json:"idClient"`
	LastName   string    `json:"lastName"`
	CreatedAt  time.Time `json:"createdAt"`
	TimeFrom   time.Time `json:"timeFrom"`
	ItemCount  int       `json:"itemCount"`
	TotalPrice float64   `json:"totalPrice"`
	Status     string    `json:"status"`
}

type Reservation struct {
	IdReservation              string    `json:"idReservation"`
	IdClient                   string    `json:"idClient"`
	IdRestaurant               string    `json:"idRestaurant"`
	IdTable                    string    `json:"idTable"`
	Status                     string    `json:"status"`
	NumberOfPeople             int       `json:"numberOfPeople"`
	CreatedAt                  time.Time `json:"createdAt"`
	TimeFrom                   time.Time `json:"timeFrom"`
	TimeTo                     time.Time `json:"timeTo"`
	ConfirmedByAdminRestaurant *string   `json:"confirmedByAdminRestaurant"`
}

type ReservationDetails struct {
	IdReservation  string    `json:"idReservation"`
	Status         string    `json:"status"`
	IdRestaurant   string    `json:"idRestaurant"`
	CreatedAt      time.Time `json:"createdAt"`
	NumberOfPeople int       `json:"numberOfPeople"`
}
type PostRatingRestaurant struct {
	IdRating     string `json:"idRating"`
	IdClient     string `json:"idClient"`
	IdRestaurant string `json:"idRestaurant"`
	RatingValue  int    `json:"rating"`
	Comment      string `json:"comment"`
}

type RatingRestaurant struct {
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	RatingValue int       `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ClientDetails struct {
	Profile        Profile
	Orders         []OrderDetails
	TotalSpent     float64
	TotalOrders    int
	FirstOrderDate *time.Time
}

// !TODO: REMOVE THE IDRESTAURANT FROM THE ORDER
type Order struct {
	//!NOTE: I think in this place im gonna fetch all the information about the order the quantity and the food and all
	IdOrder       string    `json:"idOrder"`
	IdReservation string    `json:"idReservation"`
	TotalPrice    float64   `json:"totalPrice"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
}
type UserAdmin struct {
	Id                string    `json:"idProfile"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Type              string    `json:"type"`
	Email             string    `json:"email"`
	Address           string    `json:"address"`
	Password          string    `json:"password"`
	Phone             string    `json:"phone"`
	LastLogin         time.Time `json:"lastLogin"`
	CreatedAt         time.Time `json:"createdAt"`
	IdRestaurant      string    `json:"idRestaurant"`
	IdAdminRestaurant string    `json:"idAdminRestaurant"`
}

type RestaurantRatingStats struct {
	MonthlyStats     []MonthlyRatingStats `json:"monthlyStats"`
	OverallAverage   float64              `json:"overallAverage"`
	TotalRatings     int                  `json:"totalRatings"`
	Percentage5Stars float64              `json:"percentage5Stars"`
	Percentage4Stars float64              `json:"percentage4Stars"`
	Percentage3Stars float64              `json:"percentage3Stars"`
	Percentage2Stars float64              `json:"percentage2Stars"`
	Percentage1Star  float64              `json:"percentage1Star"`
}
type MonthlyRatingStats struct {
	Month         int     `json:"month"`
	Year          int     `json:"year"`
	AverageRating float64 `json:"averageRating"`
	TotalRatings  int     `json:"totalRatings"`
}
type User struct {
	Id           string    `json:"idProfile"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Type         string    `json:"type"`
	Email        string    `json:"email"`
	Address      string    `json:"address"`
	Phone        string    `json:"phone"`
	Password     string    `json:"password"`
	LastLogin    time.Time `json:"lastLogin"`
	CreatedAt    time.Time `json:"createdAt"`
	Refreshtoken string    `json:"refreshToken"`
	ClientId     string    `json:"idClient"`
	Username     string    `json:"username"`
}

//

type Product struct {
	IdProduct      string    `json:"idProduct"`
	NameProduct    string    `json:"nameProduct"`
	Price          int       `json:"price"`
	Description    string    `json:"description"`
	IdCategorie    string    `json:"idCategorie"`
	Stock          int       `json:"stock"`
	CreatedAt      time.Time `json:"createdAt"`
	DateExpiration time.Time `json:"dateExpiration"`
	Boosted        bool      `json:"boosted"`
}

type ProductCreate struct {
	NameProduct    string `json:"nameProduct"`
	Price          string `json:"price"`
	Description    string `json:"description"`
	IdCategorie    string `json:"idCategorie"`
	Stock          string `json:"stock"`
	DateExpiration string `json:"dateExpiration"`
	Boosted        string `json:"boosted"`
}
type ProductBought struct {
	IdProduct   string `json:"idProduct"`
	NameProduct string `json:"nameProduct"`
	Price       int    `json:"price"`
	IdCategorie string `json:"idCategorie"`
	Quantity    int    `json:"quantity"`
}

type Categorie struct {
	IdCategorie   string `json:"idCategorie"`
	NameCategorie string `json:"nameCategorie"`
}
type FeedBack struct {
	IDCustomer string `json:"idCustomer"`
	IDFeedBack string `json:"idFeedback"`
	Comment    string `json:"comment"`
	CreatedAt  string `json:"createdAt"`
}
type Notification struct {
	IdNotification string `json:"idNotification"`
	IdAdmin        string `json:"idAdmin"`
	Titre          string `json:"titre"`
	Type           string `json:"type"`
	Description    string `json:"description"`
}

type Rating struct {
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Comment     string    `json:"comment"`
	RatingValue int       `json:"rating"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ReservationDetailsR struct {
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	TimeFrom       time.Time `json:"timeFrom"`
	NumberOfPeople int       `json:"numberOfPeople"`
}
type UpcomingReservationInfo struct {
	IdReservation string `json:"idReservation"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	NumberPeople  int    `json:"numberPeople"`
	Date          string `json:"date"`
	Day           string `json:"day"`
	Time          string `json:"time"`
	IdTable       string `json:"idTable"`
}

type ReservationStatsAndList struct {
	TotalToday           int                   `json:"totalToday"`
	UpcomingReservation  int                   `json:"upcomingReservation"`
	ConfirmedRate        int                   `json:"confirmedRate"`
	TodayReservations    []ReservationDetailsR `json:"todayReservations"`
	UpcomingReservations []ReservationDetailsR `json:"upcomingReservations"`
}
