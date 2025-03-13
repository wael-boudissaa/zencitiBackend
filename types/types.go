package types

import "time"

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(user User) (*User, error)
	CreateUser(user RegisterUser, idUser string, token string, hashedPassword string) error
}
type ActiviteStore interface { 
    GetActivite() (*[]Activite, error)
    GetActiviteById(id string) (*Activite, error)
    GetActiviteTypes() (*[]Activite, error)
}
type RestaurantStore interface {
    GetRestaurant() (*[]Restaurant, error)
    GetRestaurantById(id string) (*Restaurant, error)
}

type Activite struct {
    IdActivite string `json:"idActivite"`
    NameActivite string `json:"nameActivite"`
    Description string `json:"descriptionActivity"`
}

type ActivitetType struct {
    IdActiviteType string `json:"idActiviteType"`
    NameActiviteType string `json:"nameActiviteType"`
}
type Restaurant struct {
    IdRestaurant string `json:"idRestaurant"`
    NameRestaurant string `json:"nameRestaurant"`
    Description string `json:"description"`
    IdActivite string `json:"idActivite"`
    Image string `json:"image"` 
    Location string `json:"location"`
    Capacity int `json:"capacity"`

}
type RestaurantWorker struct {
    IdRestaurantWorker string `json:"idRestaurantWorker"`
    IdRestaurant string `json:"idRestaurant"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    Email string `json:"email"`
    PhoneNumber string `json:"phoneNumber"`
    Quote string `json:"quote"`
    StartWorking time.Time `json:"startWorking"`
    Nationnallity string `json:"nationnallity"`
    NativeLanguage string `json:"nativeLanguage"`
    Rating float64 `json:"rating"`
    Address string `json:"address"`
    Status string `json:"status"`
}

type RestaurantWorkerFeedBack struct {
    IdRestaurantWorkerFeedBack string `json:"idRestaurantWorkerFeedBack"`
    IdRestaurantWorker string `json:"idRestaurantWorker"`
    IdClient string `json:"idClient"`
    Comment string `json:"comment"`
    CreatedAt time.Time `json:"createdAt"`
}

type Friendship struct { 
    IdFriendship string `json:"idFriendship"`
    IdClient1 string `json:"idClient1"`
    IdClient2 string `json:"idClient2"`
    Status string `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
}
type Menu struct {
    IdMenu string `json:"idMenu"`
    IdRestaurant string `json:"idRestaurant"`
    Name string `json:"name"`
    CreatedAt time.Time `json:"createdAt"`

}
type Food struct {
    IdFood string `json:"idFood"`
    IdCategory string `json:"idCategory"`
    IdMenu string `json:"idMenu"`
    Name string `json:"name"`
    Description string `json:"description"`
    Image string `json:"image"`
    Price float64 `json:"price"`
    Status string `json:"status"`
}
//!WARNING:: THERE SHOULD BE A GENEARL THING ON THE RESERVATION FOR THE RESTAURANT AND THE ACITIVITE AND ALSO FOR THE RATING AND FEEDBACK
type Reservation struct {
    IdReservation string `json:"idReservation"`
    IdClient string `json:"idClient"`
    IdRestaurant string `json:"idRestaurant"`
    Status string `json:"status"`
    Price float64 `json:"price"`
    TimeReservation time.Time `json:"timeReservation"`
    CreatedAt time.Time `json:"createdAt"`
}
type Order struct {
    //!NOTE: I think in this place im gonna fetch all the information about the order the quantity and the food and all
    IdOrder string `json:"idOrder"`
    IdReservation string `json:"idReservation"`
    IdRestaurant string `json:"idRestaurant"`
    TotalPrice float64 `json:"totalPrice"`
    Status string `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
}


// type CategorieStore interface {
// 	GetCategories() (*[]Categorie, error)
// }
// type CommandeStore interface {
// 	GetAllCommandes() (*[]Commande, error)
// 	GetCommandeById(id string) (*Commande, error)
// 	CreateCommande(idCommande, idCustomer string, price int) error
// 	// UpdateCommande(commande Commande) error
// 	// DeleteCommande(commande Commande) error
// 	// GetCommandeByUser(idUser string) (*[]Commande, error)
// 	InsertProductINCommande(product ProductBought, idCommande string) (*CommandeProduct, error)
// }
//
// type CommandeProduct struct {
// 	IdCommande string `json:"idCommande"`
// 	IdProduct  string `json:"idProduct"`
// }
//
// type Commande struct {
// 	IdCommande string `json:"idCommande"`
// 	IdCustomer string `json:"idCustomer"`
// }
// type CommandeCreate struct {
// 	IdCustomer string `json:"idCustomer"`
// }
//
type User struct {
	Id           string `json:"idProfile"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Type         string `json:"type"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	LastLogin    time.Time `json:"lastLogin"`
	CreatedAt    time.Time `json:"createdAt"`
	Refreshtoken string `json:"refreshToken"`
}
//

type RegisterUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Phone     string `json:"phoneNumber"`
	Password  string `json:"password"`
    Type     string `json:"type"`
}

type ProductStore interface {
	GetProductById(id string) (*Product, error)
	GetAllProducts() (*[]Product, error)
	CreateProduct(product ProductCreate, idProduct string) error
	// UpdateProduct(product Product) error
	// DeleteProduct(product Product) error
	GetProductByCategorie(idCategorie string) (*[]Product, error)
}

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
type FeedBackStore interface {
	GetAllFeedBack() (*FeedBack, error)
	CreateFeedBack(idFeedBack, idCustomer, comment string) error
}
