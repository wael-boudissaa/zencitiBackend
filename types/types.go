package types

import "time"
// Add these types to your existing types.go file

type ActivityAdminCreation struct {
    FirstName   string  `json:"admin_firstName"`
    LastName    string  `json:"admin_lastName"`
    Email       string  `json:"admin_email"`
    Phone       string  `json:"admin_phone"`
    Address     string  `json:"admin_address"`
    Password    string  `json:"admin_password"`
    Type        string  `json:"admin_type"` // should be "adminActivity"
}

type ActivityCreationWithAdmin struct {
    Name            string  `json:"name"`
    Description     string  `json:"description"`
    Image           string  `json:"image"`
    Longitude       float64 `json:"longitude"`
    Latitude        float64 `json:"latitude"`
    IdTypeActivity  string  `json:"idTypeActivity"`
    Capacity        int     `json:"capacity"`
}

type ActivityStats struct {
    TotalBookings      int                    `json:"totalBookings"`
    CompletedBookings  int                    `json:"completedBookings"`
    PendingBookings    int                    `json:"pendingBookings"`
    CancelledBookings  int                    `json:"cancelledBookings"`
    AvgEngagement      float64               `json:"avgEngagement"`
    TotalReviews       int                    `json:"totalReviews"`
    AverageRating      float64               `json:"averageRating"`
    BookingsToday      int                    `json:"bookingsToday"`
    BookingsThisWeek   int                    `json:"bookingsThisWeek"`
    BookingsThisMonth  int                    `json:"bookingsThisMonth"`
    DailyTrends        []ActivityDailyStats   `json:"dailyTrends"`
    WeeklyTrends       []ActivityWeeklyStats  `json:"weeklyTrends"`
    MonthlyTrends      []ActivityMonthlyStats `json:"monthlyTrends"`
    RecentBookings     []ActivityBookingInfo  `json:"recentBookings"`
    TopRatedReviews    []ActivityReviewDetail `json:"topRatedReviews"`
}

type ActivityDailyStats struct {
    Date     string `json:"date"`
    Bookings int    `json:"bookings"`
}

type ActivityWeeklyStats struct {
    Week     string `json:"week"`
    Bookings int    `json:"bookings"`
}

type ActivityMonthlyStats struct {
    Month    string `json:"month"`
    Year     int    `json:"year"`
    Bookings int    `json:"bookings"`
}

type ActivityBookingInfo struct {
    ClientName    string    `json:"clientName"`
    BookingTime   time.Time `json:"bookingTime"`
    Status        string    `json:"status"`
    CreatedAt     time.Time `json:"createdAt"`
}

type AdminActivityProfile struct {
    IdAdminActivity string    `json:"idAdminActivity"`
    IdProfile       string    `json:"idProfile"`
    FirstName       string    `json:"firstName"`
    LastName        string    `json:"lastName"`
    Email           string    `json:"email"`
    Phone           string    `json:"phone"`
    Address         string    `json:"address"`
    CreatedAt       time.Time `json:"createdAt"`
    Activities      []Activity `json:"activities"`
}

type ActivityDetails struct {
	IdActivity     string                 `json:"idActivity"`
	NameActivity   string                 `json:"nameActivity"`
	Description    string                 `json:"description"`
	ImageActivite  string                 `json:"imageActivite"`
	Langitude      float64                `json:"langitude"`
	Latitude       float64                `json:"latitude"`
	IdTypeActivity string                 `json:"idTypeActivity"`
	Capacity       int                    `json:"capacity"`
    IdAdminActivity string                `json:"idAdminActivity"` // Optional for public activities
	AdminName      string                 `json:"adminName"`
	AdminEmail     string                 `json:"adminEmail"`
	AdminPhone     string                 `json:"adminPhone"`
	RatingCounts   map[int]int            `json:"ratingCounts"` // e.g. {5: 10, 4: 3, ...}
	RecentReviews  []ActivityReviewDetail `json:"recentReviews"`
}

type ActivityReviewDetail struct {
	ReviewerName string `json:"reviewerName"`
	Rating       int    `json:"rating"`
	Comment      string `json:"comment"`
	CreatedAt    string `json:"createdAt"`
}

type MonthlyUserStats struct {
	Month       int `json:"month"`
	Year        int `json:"year"`
	NewUsers    int `json:"newUsers"`
	ActiveUsers int `json:"activeUsers"`
}

// Add to types/types.go

type UniversalReservationDetails struct {
	// Common fields
	ReservationType string    `json:"reservationType"` // "restaurant" or "activity"
	ReservationID   string    `json:"reservationId"`
	ClientID        string    `json:"clientId"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	ReservationTime time.Time `json:"reservationTime"`

	// Client information
	ClientFirstName string `json:"clientFirstName"`
	ClientLastName  string `json:"clientLastName"`
	ClientEmail     string `json:"clientEmail"`
	ClientPhone     string `json:"clientPhone"`
	ClientUsername  string `json:"clientUsername"`

	// Restaurant-specific fields (nil for activities)
	RestaurantInfo *RestaurantReservationInfo `json:"restaurantInfo,omitempty"`

	// Activity-specific fields (nil for restaurants)
	ActivityInfo *ActivityReservationInfo `json:"activityInfo,omitempty"`
}

type RestaurantReservationInfo struct {
	RestaurantID       string  `json:"restaurantId"`
	RestaurantName     string  `json:"restaurantName"`
	RestaurantImage    string  `json:"restaurantImage"`
	RestaurantLocation string  `json:"restaurantLocation"`
	Description        string  `json:"description"`
	Capacity           int     `json:"capacity"`
	Longitude          float64 `json:"longitude"`
	Latitude           float64 `json:"latitude"`
	NumberOfPeople     int     `json:"numberOfPeople"`
	TableID            string  `json:"tableId,omitempty"`

	// Admin information
	AdminFirstName string `json:"adminFirstName"`
	AdminLastName  string `json:"adminLastName"`
	AdminEmail     string `json:"adminEmail"`
	AdminPhone     string `json:"adminPhone"`
}

type ActivityReservationInfo struct {
	ActivityID          string  `json:"activityId"`
	ActivityName        string  `json:"activityName"`
	ActivityDescription string  `json:"activityDescription"`
	ActivityImage       string  `json:"activityImage"`
	ActivityType        string  `json:"activityType"`
	Capacity            int     `json:"capacity"`
	Longitude           float64 `json:"longitude"`
	Latitude            float64 `json:"latitude"`

	// Admin information
	AdminFirstName string `json:"adminFirstName"`
	AdminLastName  string `json:"adminLastName"`
	AdminEmail     string `json:"adminEmail"`
	AdminPhone     string `json:"adminPhone"`
}

type UserStats struct {
	TotalUsers        int                `json:"totalUsers"`
	ActiveUsersToday  int                `json:"activeUsersToday"`
	NewUsersThisMonth int                `json:"newUsersThisMonth"`
	MonthlyStats      []MonthlyUserStats `json:"monthlyStats"`
}

type CampusUser struct {
	IdProfile     string   `json:"idProfile"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	Email         string   `json:"email"`
	Type          string   `json:"type"`
	Address       *string  `json:"address"`
	PhoneNumber   *string  `json:"phoneNumber"`
	CreatedAt     string   `json:"createdAt"`
	Username      *string  `json:"username,omitempty"`      // only for clients
	Roles         []string `json:"roles"`                   // can be multiple: ["client", "adminActivity"]
	IdClient      *string  `json:"idClient,omitempty"`
	IdAdmin       *string  `json:"idAdmin,omitempty"`
	IdAdminActivity *string `json:"idAdminActivity,omitempty"`
	IdAdminRestaurant *string `json:"idAdminRestaurant,omitempty"`
	AdminActivityStatus string `json:"adminActivityStatus,omitempty"` // "active", "inactive", or empty
	AdminRestaurantStatus string `json:"adminRestaurantStatus,omitempty"` // "active", "inactive", or empty
	AssignedActivityId *string `json:"assignedActivityId,omitempty"`
	AssignedActivityName *string `json:"assignedActivityName,omitempty"`
	AssignedRestaurantId *string `json:"assignedRestaurantId,omitempty"`
	AssignedRestaurantName *string `json:"assignedRestaurantName,omitempty"`
}
type LocationItemWithDistance struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Type              string  `json:"type"` // "Restaurant" or "Activity"
	Address           *string `json:"address"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	ImageURL          string  `json:"imageUrl"`
	PhoneNumber       string  `json:"phoneNumber"`
	Distance          float64 `json:"distance"`          // in kilometers
	DistanceFormatted string  `json:"distanceFormatted"` // "1.2 km" or "500 m"
}
type LocationItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"` // "Restaurant" or "Activity"
	Address     *string `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ImageURL    string  `json:"imageUrl"`
	PhoneNumber string  `json:"phoneNumber"`
}

type RestaurantCreation struct {
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
	Location    string  `json:"location"`
}
type AdminLocation struct {
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	HasLocation bool     `json:"hasLocation"`
}

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
	IdActivity      string   `json:"idActivity"`
	IdAdminActivity *string  `json:"idAdminActivity"`
	NameActivity    string   `json:"nameActivity"`
	Description     string   `json:"descriptionActivity"`
	Langitude       *float64 `json:"longitude" db:"longitude"`
	Latitude        *float64 `json:"latitude" db:"latitude"`

	IdTypeActivity string `json:"idTypeActivity"`
	ImageActivite  string `json:"imageActivity"`
	Capacity       int    `json:"capacity"`
}

type ActivitetType struct {
	IdActiviteType   string `json:"idTypeActivity"`
	NameActiviteType string `json:"nameTypeActivity"`
	ImageActivity    string `json:"imageActivity"`
}
// ActivityCategoryCreation represents the data needed to create a new activity category
type ActivityCategoryCreation struct {
	NameTypeActivity string `json:"nameTypeActivity"`
	ImageActivity    string `json:"imageActivity"`
}

// Notification represents a notification in the system
type Notification struct {
	IdNotification string `json:"idNotification"`
	IdAdmin        string `json:"idAdmin"`
	Titre          string `json:"titre"`
	Type           string `json:"type"`
	Description    string `json:"description"`
}

// NotificationCreation represents the data needed to create a new notification
type NotificationCreation struct {
	IdAdmin     string `json:"idAdmin"`
	Titre       string `json:"titre"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Feedback represents feedback from a client
type Feedback struct {
	IdFeedback int    `json:"idFeedback"`
	IdClient   string `json:"idClient"`
	Comment    string `json:"comment"`
	CreatedAt  string `json:"createdAt"`
	// Client information
	ClientFirstName   string `json:"clientFirstName,omitempty"`
	ClientLastName    string `json:"clientLastName,omitempty"`
	ClientUsername    string `json:"clientUsername,omitempty"`
	ClientEmail       string `json:"clientEmail,omitempty"`
	ClientPhoneNumber string `json:"clientPhoneNumber,omitempty"`
}

// FeedbackCreation represents the data needed to create new feedback
type FeedbackCreation struct {
	IdClient string `json:"idClient"`
	Comment  string `json:"comment"`
}

// Missing types for interface compatibility
type ReservationStatsAndList struct {
	TotalToday           int                     `json:"totalToday"`
	UpcomingReservation  int                     `json:"upcomingReservation"`
	ConfirmedRate        float64                 `json:"confirmedRate"`
	TodayReservations    []ReservationDetailsR   `json:"todayReservations"`
	UpcomingReservations []ReservationDetailsR   `json:"upcomingReservations"`
}

type ReservationDetailsR struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	TimeFrom       string `json:"timeFrom"`
	NumberOfPeople int    `json:"numberOfPeople"`
}

type Rating struct {
	IdRating     string `json:"idRating"`
	RatingValue  int    `json:"rating"`
	Comment      string `json:"comment"`
	CreatedAt    string `json:"createdAt"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

type UpcomingReservationInfo struct {
	IdReservation  string `json:"idReservation"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	NumberPeople   int    `json:"numberPeople"`
	Date           string `json:"date"`
	Day            string `json:"day"`
	Time           string `json:"time"`
	IdTable        string `json:"idTable"`
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
	AverageRating     *float64 `json:"averageRating,omitempty"`
	IsActive          *bool    `json:"isActive,omitempty"`
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

type PostRatingActivity struct {
	IdRating    string `json:"idRating"`
	IdClient    string `json:"idClient"`
	IdActivity  string `json:"idActivity"`
	RatingValue int    `json:"rating"`
	Comment     string `json:"comment"`
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
	HasSensors   bool      `json:"hasSensors"`
	SensorCount  int       `json:"sensorCount"`
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

type CampusFacilityItem struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Type        string   `json:"type"` // "activity" or "restaurant"
	Image       *string  `json:"image"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Capacity    *int     `json:"capacity"`
	Location    *string  `json:"location,omitempty"`    // Only for restaurants
	AdminID     *string  `json:"adminId,omitempty"`     // Admin managing this facility
	AdminName   *string  `json:"adminName,omitempty"`   // Admin's full name
	AdminEmail  *string  `json:"adminEmail,omitempty"`  // Admin's email
	AdminStatus *string  `json:"adminStatus,omitempty"` // "active" or "inactive"
	CategoryID  *string  `json:"categoryId,omitempty"`  // Only for activities
	CategoryName *string `json:"categoryName,omitempty"` // Only for activities
}

type CampusFacilitiesResponse struct {
	Activities  []CampusFacilityItem `json:"activities"`
	Restaurants []CampusFacilityItem `json:"restaurants"`
	Total       int                  `json:"total"`
	ActivityCount int                `json:"activityCount"`
	RestaurantCount int              `json:"restaurantCount"`
}

type FollowingFollowerInfo struct {
	ClientId  string `json:"clientId"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type FollowListResponse struct {
	Following []FollowingFollowerInfo `json:"following"`
	Followers []FollowingFollowerInfo `json:"followers"`
}

type AvailabilityCheckRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

type AvailabilityCheckResponse struct {
	EmailExists    bool `json:"emailExists"`
	UsernameExists bool `json:"usernameExists"`
	Available      bool `json:"available"`
}

// Sensor-related types for water consumption tracking
type SensorRegistration struct {
	SensorId string `json:"sensorId"`
	ClientId string `json:"clientId"`
}

type DailyUsageData struct {
	SensorId     string  `json:"sensorId"`
	UsageDate    string  `json:"usageDate"`
	VolumeLiters float64 `json:"volumeLiters"`
}

type BatchUsageData struct {
	SensorId   string           `json:"sensorId"`
	UsageData  []DailyUsageData `json:"usageData"`
}

type SensorInfo struct {
	IdSensor     string    `json:"idSensor"`
	IdClient     string    `json:"idClient"`
	Status       string    `json:"status"`
}

type SensorUsageResponse struct {
	Sensors []SensorUsageDetails `json:"sensors"`
}

type SensorUsageDetails struct {
	IdSensor      string              `json:"idSensor"`
	Status        string              `json:"status"`
	DailyUsage    []DailyUsageRecord  `json:"dailyUsage"`
	WeeklyTotal   float64             `json:"weeklyTotal"`
	MonthlyTotal  float64             `json:"monthlyTotal"`
	AverageDaily  float64             `json:"averageDaily"`
}

type DailyUsageRecord struct {
	Date         string  `json:"date"`
	VolumeLiters float64 `json:"volumeLiters"`
}

type UserSensorsResponse struct {
	Sensors      []SensorInfo `json:"sensors"`
	TotalSensors int          `json:"totalSensors"`
	HasSensors   bool         `json:"hasSensors"`
}

// Activity booking detail with full client information
type ActivityBookingDetail struct {
	IdClientActivity string    `json:"idClientActivity"`
	ClientName       string    `json:"clientName"`
	ClientEmail      string    `json:"clientEmail"`
	ClientPhone      string    `json:"clientPhone"`
	ClientUsername   string    `json:"clientUsername"`
	BookingTime      time.Time `json:"bookingTime"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
}

// Enhanced analytics with chart-ready data
type ActivityDetailedAnalytics struct {
	// Basic stats (reuse existing)
	TotalBookings      int                    `json:"totalBookings"`
	CompletedBookings  int                    `json:"completedBookings"`
	PendingBookings    int                    `json:"pendingBookings"`
	CancelledBookings  int                    `json:"cancelledBookings"`
	CompletionRate     float64               `json:"completionRate"`
	TotalReviews       int                    `json:"totalReviews"`
	AverageRating      float64               `json:"averageRating"`
	
	// Enhanced analytics
	BookingsByStatus   map[string]int         `json:"bookingsByStatus"`
	PeakHours          []HourlyBookingStats   `json:"peakHours"`
	DailyTrends        []ActivityDailyStats   `json:"dailyTrends"`
	WeeklyTrends       []ActivityWeeklyStats  `json:"weeklyTrends"`
	MonthlyTrends      []ActivityMonthlyStats `json:"monthlyTrends"`
	ClientReturnRate   float64               `json:"clientReturnRate"`
	CapacityUtilization float64              `json:"capacityUtilization"`
	RecentBookings     []ActivityBookingDetail `json:"recentBookings"`
	TopRatedReviews    []ActivityReviewDetail  `json:"topRatedReviews"`
}

// Hourly booking statistics
type HourlyBookingStats struct {
	Hour     int `json:"hour"`
	Bookings int `json:"bookings"`
}

// Client demographics for analytics
type ClientDemographicsStats struct {
	UniqueClients  int `json:"uniqueClients"`
	ReturningRate  float64 `json:"returningRate"`
	AverageBookingsPerClient float64 `json:"averageBookingsPerClient"`
}

// Admin Restaurant Statistics - aggregated data for admin panel
type AdminRestaurantStats struct {
	TotalRestaurants    int     `json:"totalRestaurants"`
	ActiveRestaurants   int     `json:"activeRestaurants"`
	AverageRating       float64 `json:"averageRating"`
	TotalBookingsLastMonth int `json:"totalBookingsLastMonth"`
}

// Restaurant Today Summary - daily metrics for a specific restaurant
type RestaurantTodaySummary struct {
	TotalReservationsToday int     `json:"totalReservationsToday"`
	CurrentOccupancy       float64 `json:"currentOccupancy"`
	ConfirmedReservations  int     `json:"confirmedReservations"`
	PendingReservations    int     `json:"pendingReservations"`
}
