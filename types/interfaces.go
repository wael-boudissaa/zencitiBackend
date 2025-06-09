package types

import "time"

type ProfileStore interface {
	GetProfileById(id string) (*User, error)
	CreateProfile(profile User, id string) error
	UpdateProfile(profile User) error
	DeleteProfile(profile User) error
}

type ProductStore interface {
	GetProductById(id string) (*Product, error)
	GetAllProducts() (*[]Product, error)
	CreateProduct(product ProductCreate, idProduct string) error
	// UpdateProduct(product Product) error
	// DeleteProduct(product Product) error
	GetProductByCategorie(idCategorie string) (*[]Product, error)
}

type FeedBackStore interface {
	GetAllFeedBack() (*FeedBack, error)
	CreateFeedBack(idFeedBack, idCustomer, comment string) error
}
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(user User) (*User, error)
	CreateUser(user RegisterUser, idUser string, token string, hashedPassword string) error
	CreateClient(idUser, idClient, username string) error
	GetClientIdByUsername(username string) (string, error)
	SearchUsersByUsernamePrefix(prefix string) (*[]string, error)
	SendRequestFriend(idFriendship string, idSender string, idReceiver string) error
	AcceptRequestFriend(idFriendship string) error
	GetFriendshipRequested(idClient string) (*[]Friendship, error)
	CountFollowers(idClient string) (int, error)
	CountFollowing(idClient string) (int, error)
	GetClientInformationUsername(username string) (*ProfilePage, error)
	GetClientInformation(idClient string) (*ProfilePage, error)
}

type ActiviteStore interface {
	// GetActivite() (*[]Activite, error)
	GetRecentActivities(idClient string) (*[]ActivityProfile, error)
	GetActiviteById(id string) (*Activity, error)
	GetActivityByTypes(typeActivite string) (*[]Activity, error)
	GetActiviteTypes() (*[]ActivitetType, error)
	GetPopularActivities() (*[]Activity, error)
}

type RestaurantStore interface {
	GetRestaurantTables(restaurantId string, timeSlot time.Time) (*[]RestaurantTableStatus, error)
	GetRestaurant() (*[]Restaurant, error)
	GetRestaurantById(id string) (*Restaurant, error)
	CreateReservation(idReservation string, reservation ReservationCreation) error
	CreateOrder(idOrder string, order OrderCreation) error
	GetReservationTodayByRestaurantId(idRestaurant string) (*[]ReservationListInformation, error)
	AddFoodToOrder(food AddFoodToOrder) error
	PostOrderList(orderId string, food []FoodItem) error
	CountReservationUpcomingWeek(idRestaurant string) (int, error)
    CountReservationLastMonth(idRestaurant string) (*[]ReservationStats, error)
	CountOrderReceivedToday(idRestaurant string) (int, error)
	CountReservationReceivedToday(idRestaurant string) (int, error)
	GetAvailableMenuInformation(restaurantId string) (*[]MenuInformationFood, error)
	ReserveTable(idReservation string, reservation ReservationCreation) error
	GetFriendsOfClient(idClient string) (*[]string, error)
	GetRatingOfFriendsRestaurant(friendsId []string, idRestaurant string) (*[]RatingRestaurant, error)
	PostRatingRestaurant(rating PostRatingRestaurant) error
	// GetRestaurantWorker() (*[]RestaurantWorker, error)
	// GetRestaurantWorkerById(id string) (*RestaurantWorker, error)
	// GetRestaurantWorkerFeedBack(id string) (*[]RestaurantWorkerFeedBack, error)
	// GetRestaurantWorkerFeedBackById(id string) (*RestaurantWorkerFeedBack, error)
	// GetReservation() (*[]Reservation, error)
	// GetReservationById(id string) (*Reservation, error)
	// PostReservation(reservation Reservation) error
	// PostOrder(order Order) error
	// PostWorkerFeedBack(workerFeedBack RestaurantWorkerFeedBack) error
	// GetOrder() (*[]Order, error)
	// GetOrderById(id string) (*Order, error)
	// GetMenu() (*[]Menu, error)
	// GetMenuById(id string) (*Menu, error)
	// getMenuByRestaurantId(id string) (*Menu, error)
	GetFoodByMenu(idMenu string) (*[]Food, error)
	// GetFoodById(id string) (*Food, error)
	// GetWorkerFeedBack() (*[]RestaurantWorkerFeedBack, error)
	// GetWorkerRestaurantFeedBackBy(id string) (*RestaurantWorkerFeedBack, error)
}
