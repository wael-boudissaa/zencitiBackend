package types

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
    CreateClient(idUser, idClient string) error
}
type ActiviteStore interface {
	// GetActivite() (*[]Activite, error)
	GetActiviteById(id string) (*Activity, error)
	GetActivityByTypes(typeActivite string) (*[]Activity, error)
	GetActiviteTypes() (*[]ActivitetType, error)
	GetPopularActivities() (*[]Activity, error)
}

type RestaurantStore interface {
    GetRestaurantTables(restaurantId string) (*[]RestaurantTable, error)
	GetRestaurant() (*[]Restaurant, error)
	GetRestaurantById(id string) (*Restaurant, error)
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
	// GetFoodByMenu() (*[]Food, error)
	// GetFoodById(id string) (*Food, error)
	// GetWorkerFeedBack() (*[]RestaurantWorkerFeedBack, error)
	// GetWorkerRestaurantFeedBackBy(id string) (*RestaurantWorkerFeedBack, error)
}
