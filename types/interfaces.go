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
	DeleteFriendRequestFromDB(idFriendship string) error
    RemoveFromFollowing(currentUserId string, targetUserId string) error
    RemoveFollower(currentUserId string, targetUserId string) error
	VerifyAdminRestaurantAssignment(idAdminRestaurant string) (bool, string, error)
	IsClientAdminActivity(idProfile string) (bool, string, error)
	GetAdminLocation(idAdmin string) (*AdminLocation, error)
	SetAdminLocation(idAdmin string, latitude, longitude float64) error
	CreateRestaurantWithAdmin(restaurantData RestaurantCreation, profileData RegisterAdmin) (string, string, error)
	AssignClientToAdminActivity(idClient string) error
	GetUserByEmail(email string) (*User, error)
	GetAdminByEmail(email string) (*UserAdmin, error)
	GetGeneralAdminByEmail(email string) (*User, error)
	GetAllCampusUsers() ([]CampusUser, error)
	AssignUserToRole(idUser string, role string) error
	AssignUserToRoleWithEntity(idUser string, role string, idActivity string, idRestaurant string) error
	UpdateActivityAdmin(idActivity string, idAdminActivity string) error
	GetAllClients() ([]ClientInfo, error)
	UpdateClientLocation(idClient string, longitude, latitude float64) error
	GetUserById(user User) (*User, error)
	CreateUser(user interface{}, idUser string, token string, hashedPassword string) error
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
	CreateAdminRestaurant(idUser string, idAdminRestaurant string) error
	GetUserStats() (*UserStats, error)
	GetMonthlyUserStats() ([]MonthlyUserStats, error)
	CreateActivityWithAdmin(activityData ActivityCreationWithAdmin, profileData ActivityAdminCreation) (string, string, error)

	// Notification methods
	CreateNotification(notification NotificationCreation) (string, error)
	GetNotificationsByAdmin(idAdmin string) ([]Notification, error)
	GetAllNotifications() ([]Notification, error)

	// Feedback methods
	CreateFeedback(feedback FeedbackCreation) error
	GetAllFeedbackWithClientInfo() ([]Feedback, error)

	// Following/Followers methods
	GetClientFollowersAndFollowing(idClient string) (*FollowListResponse, error)

	// Availability check methods
	CheckEmailExists(email string) (bool, error)
	CheckUsernameExists(username string) (bool, error)
	CheckEmailAndUsernameAvailability(email, username string) (*AvailabilityCheckResponse, error)

	// Friend request status check methods
	CheckFriendRequestStatus(fromUsername, toUsername string) (*FriendRequestStatusResponse, error)

	UpdateRestaurantAdmin(idRestaurant string, idAdminRestaurant string) error
	CreateAdminActivity(idUser string, idAdminRestaurant string) error
}

type ActiviteStore interface {
	// GetActivite() (*[]Activite, error)
	GetAllLocationsWithDistances(clientLat, clientLng float64) (*[]LocationItemWithDistance, error)
	UpdateClientActivityStatus(idClientActivity string, idAdminActivity string) error
	GetActivityStatsAdmin(idAdminActivity string) (*ActivityStats, error)
	GetRecentActivities(idClient string) (*[]ActivityProfile, error)
	GetActiviteById(id string) (*Activity, error)
	GetActivityByTypes(typeActivite string) (*[]Activity, error)
	GetActiviteTypes() (*[]ActivitetType, error)
	CreateActivityCategory(category ActivityCategoryCreation) (string, error)
	CreateActivityClient(idClientActivity string, act ActivityCreation) error
	GetActivityFullDetails(id string) (*ActivityDetails, error)
	GetActivityNotAvaialableAtday(day time.Time, idActivity string) ([]string, error)
	GetAllClientActivities(idClient string) ([]ClientActivityInfo, error)
	GetPopularActivities() (*[]Activity, error)
	GetActivityStats(idAdminActivity string) (*ActivityStats, error)
	GetActivitiesByAdminActivity(idAdminActivity string) ([]Activity, error)
	UpdateActivityStatus(idClientActivity string, status string) error
	PostRatingActivity(rating PostRatingActivity) error
	GetAllCampusFacilities() (*CampusFacilitiesResponse, error)
	
	// New methods for bookings and analytics
	GetActivityBookings(idActivity string) ([]ActivityBookingDetail, error)
	GetActivityDetailedAnalytics(idActivity string) (*ActivityDetailedAnalytics, error)
	GetAdminActivityBookings(idAdminActivity string) ([]ActivityBookingDetail, error)
}

type RestaurantStore interface {
	GetRestaurantByIdProfile(idProfile string) (*UserAdmin, error)
	GetRestaurantRatingStats(idRestaurant string) (*RestaurantRatingStats, error)
	GetReservationStatsAndList(idRestaurant string) (*ReservationStatsAndList, error)
	GetFoodRestaurant(idRestaurant string) (*[]Food, error)
	GetAllClientReservations(idClient string) ([]ClientReservationInfo, error)
	GetUniversalReservationDetails(reservationId string, reservationType string) (*UniversalReservationDetails, error)
	CountFirstTimeReservers(idRestaurant string) (int, error)
	GetRestaurantTables(restaurantId string, timeSlot time.Time) (*[]RestaurantTableStatus, error)
	GetRecentOrders(idRestaurant string, limit int) ([]RecentOrder, error)
	GetRestaurant() (*[]Restaurant, error)
	GetRestaurantById(id string) (*Restaurant, error)
	BulkUpdateRestaurantTables(idRestaurant string, tables []Table) error
	SetMenuActive(idMenu, idRestaurant string) error
	GetRestaurantWorkerWithRatings(idRestaurantWorker string) (*RestaurantWorkerWithRatings, error)
	CreateReservation(idReservation string, reservation ReservationCreation) error
	GetOrderInformation(idOrder string) (*OrderInformation, error)
	UpdateOrderStatus(idOrder string, status string) error
	GetAllRestaurantReservations(idRestaurant string, page, limit int) (*PaginatedReservations, error)
	GetReservationDetails(idReservation string) (*ReservationIdDetails, error)
	GetRecentReviews(idRestaurant string) ([]*Rating, error)
	CreateOrder(idOrder string, order OrderCreation) error
	GetReservationTodayByRestaurantId(idRestaurant string) (*[]ReservationListInformation, error)
	AddFoodToOrder(food AddFoodToOrder) error
	PostOrderList(orderId string, food []FoodItem) error
	CountReservationUpcomingWeek(idRestaurant string) (int, error)
	CountReservationLastMonth(idRestaurant string) (*[]ReservationStats, error)
	GetTableOccupationToday(idRestaurant string) ([]TableOccupation, error)
	GetTopFoodsThisWeek(idRestaurant string) ([]FoodPopularity, error)
	CreateRestaurantWorker(id string, idRestaurant string, worker RestaurantWorkerCreation) error
	CreateRestaurant(idRestaurant, idAdminRestaurant, name, image string, longitude, latitude float64, description string, capacity int, location string) error
	AddFoodToMenu(idMenuFood, idMenu, idFood string) error
	SetFoodStatusInMenu(idFood, status string) error
	UpdateRestaurantWorker(id string, worker RestaurantWorker) error
	SetRestaurantWorkerStatus(idRestaurantWorker string, status string) error
	UpdateFood(idFood string, food Food) error
	GetMenusByRestaurant(idRestaurant string) ([]Menu, error)
	GetFoodsOfActiveMenu(idRestaurant string) ([]Food, error)
	GetFoodById(idFood string) (*Food, error)
	GetRestaurantMenuStats(restaurantId string) (*RestaurantMenuStats, error)
	GetUpcomingReservations(restaurantId string) ([]UpcomingReservationInfo, error)
	GetMenuWithFoods(idMenu string) (*Menu, *[]Food, error)
	CreateTable(table Table) error
	UpdateTable(idTable string, table Table) error
	DeleteTable(idTable string) error
	GetTablesByRestaurant(restaurantId string) ([]Table, error)
	UpdateReservationStatus(idReservation, status string) error
	CreateNotification(notification Notification) error
	GetNotifications() ([]Notification, error)
	CreateFood(idFood, idCategory, idRestaurant, name, description, image string, price float64, status string) error
	CreateMenu(idMenu, idRestaurant, name string) error
	// errorCreateMenu(idMenu, idRestaurant, name string) error
	SetFoodUnavailable(idFood string) error
	GetRestaurantWorker(idRestaurant string) (*[]RestaurantWorker, error)
	CreateFoodCategory(idCategory, nameCategorie string) error
	CountOrderReceivedToday(idRestaurant string) (int, error)
	CountReservationReceivedToday(idRestaurant string) (int, error)
	GetFoodCategoriesByRestaurant() ([]FoodCategory, error)
	DeleteFood(idFood string) error
	GetAvailableMenuInformation(restaurantId string) (*[]MenuInformationFood, error)
	ReserveTable(idReservation string, reservation ReservationCreation) error
	GetFriendsOfClient(idClient string) (*[]string, error)
	GetRatingOfFriendsRestaurant(friendsId []string, idRestaurant string) (*[]RatingRestaurant, error)
	PostRatingRestaurant(rating PostRatingRestaurant) error
	GetOrderStatsByHourAndStatus(idRestaurant string) (map[int]int, map[string]int, error)
	GetClientReservationAndOrderDetails(idClient string) (*ClientDetails, error)
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
	
	// New admin-level statistics and enhanced APIs
	GetAdminRestaurantStats() (*AdminRestaurantStats, error)
	GetAllRestaurantReviews(idRestaurant string) ([]*Rating, error)
	GetRestaurantTodaySummary(idRestaurant string) (*RestaurantTodaySummary, error)
}

type SensorStore interface {
	// Sensor registration methods
	RegisterSensor(sensorId, clientId string) error
	GetSensorsByClient(clientId string) ([]SensorInfo, error)
	CheckSensorOwnership(sensorId string) (string, error) // returns clientId if exists
	
	// Daily usage methods
	SaveDailyUsage(usage DailyUsageData) error
	SaveBatchUsage(batchData BatchUsageData) error
	GetSensorUsage(clientId, period string) (*SensorUsageResponse, error)
	GetSensorUsageByDate(sensorId, startDate, endDate string) ([]DailyUsageRecord, error)
	
	// Sensor management
	UpdateSensorStatus(sensorId, status string) error
	GetSensorInfo(sensorId string) (*SensorInfo, error)
}
