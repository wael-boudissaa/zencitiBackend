package user

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	// "log"
	"net/http"

	// "os"

	"github.com/gorilla/mux"
	// "github.com/gorilla/sessions"
	// "github.com/joho/godotenv"
	// "github.com/markbates/goth"
	// "github.com/markbates/goth/gothic"
	// "github.com/markbates/goth/providers/google"
	// "github.com/wael-boudissaa/zencitiBackend/configs"
	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// router.HandleFunc("/auth/{provider}", beginAuth).Methods("GET")
	// router.HandleFunc("/auth/{provider}/callback", completeAuth).Methods("GET")
	router.HandleFunc("/ws/client/location", h.ClientLocationWS)
	// router.HandleFunc("/auth/logout", logout).Methods("GET")
	//!NOTE: Client
	router.HandleFunc("/login", h.loginUser).Methods("POST")
	router.HandleFunc("/getfriendship/{idClient}", h.GetFriendshipClient).Methods("GET")
	router.HandleFunc("/acceptfriendship", h.AcceptFrienship).Methods("POST")
	router.HandleFunc("/signup", h.signUpUser).Methods("POST")
	router.HandleFunc("/sendrequest", h.SendRequestFriend).Methods("POST")
	router.HandleFunc("/clientinformation/{idClient}", h.ClientInformation).Methods("GET")
	router.HandleFunc("/usernameinformation/{username}", h.ClientInformationUsername).Methods("GET")
	router.HandleFunc("/username", h.GetUsername).Methods("GET")
	router.HandleFunc("/admin/assignactivity", h.AssignClientToAdminActivity).Methods("POST")
	router.HandleFunc("/admin/clients", h.GetAllClients).Methods("GET")

	//!NOTE: admin
	router.HandleFunc("/admin/{idAdmin}/location", h.SetAdminLocation).Methods("PUT")
	router.HandleFunc("/api/admin/{idAdmin}/location", h.GetAdminLocation).Methods("GET")

	router.HandleFunc("/admin/login", h.loginRestaurant).Methods("POST")
	router.HandleFunc("/admin/create", h.CreateAdmin).Methods("POST")
	router.HandleFunc("/restaurant/create-with-admin", h.CreateRestaurantWithAdmin).Methods("POST")
	router.HandleFunc("/activity/create-with-admin", h.CreateActivityWithAdmin).Methods("POST")

	router.HandleFunc("/users/stats", h.GetUserStats).Methods("GET")
}

const (
	key    = "key"
	MaxAge = 86400 * 30
	isProd = false
)

func (h *Handler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetUserStats()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, stats)
}

func (h *Handler) SetAdminLocation(w http.ResponseWriter, r *http.Request) {
	idAdmin := mux.Vars(r)["idAdmin"]
	if idAdmin == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idAdmin is required"))
		return
	}

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("latitude must be between -90 and 90"))
		return
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("longitude must be between -180 and 180"))
		return
	}

	err := h.store.SetAdminLocation(idAdmin, req.Latitude, req.Longitude)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"latitude":  req.Latitude,
		"longitude": req.Longitude,
	})
}

func (h *Handler) GetAdminLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adminID := vars["idAdmin"]

	location, err := h.store.GetAdminLocation(adminID)
	if err != nil {
		http.Error(w, "Location not found", http.StatusNotFound)
		return
	}

	response := map[string]float64{
		"latitude":  *location.Latitude,
		"longitude": *location.Longitude,
	}

	utils.WriteJson(w, http.StatusOK, response)
}

func (h *Handler) CreateRestaurantWithAdmin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error parsing form: %v", err))
		return
	}

	// Parse profile data
	var profileData types.RegisterAdmin
	profileData.Email = r.FormValue("email")
	profileData.Password = r.FormValue("password")
	profileData.FirstName = r.FormValue("first_name")
	profileData.LastName = r.FormValue("last_name")
	profileData.Address = r.FormValue("address")
	profileData.Type = r.FormValue("type")
	profileData.Phone = r.FormValue("phone_number")

	// Validate required profile fields
	if profileData.Email == "" || profileData.Password == "" || profileData.FirstName == "" || profileData.LastName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("email, password, first_name, and last_name are required"))
		return
	}

	if profileData.Type != "adminRestaurant" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("type must be adminRestaurant"))
		return
	}

	// Parse restaurant data
	var restaurantData types.RestaurantCreation
	restaurantData.Name = r.FormValue("restaurant_name")
	restaurantData.Description = r.FormValue("restaurant_description")
	restaurantData.Location = r.FormValue("restaurant_location")

	// Parse numeric fields
	if capacityStr := r.FormValue("restaurant_capacity"); capacityStr != "" {
		capacity, err := strconv.Atoi(capacityStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid capacity format"))
			return
		}
		restaurantData.Capacity = capacity
	}

	if longitudeStr := r.FormValue("restaurant_longitude"); longitudeStr != "" {
		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid longitude format"))
			return
		}
		restaurantData.Longitude = longitude
	}

	if latitudeStr := r.FormValue("restaurant_latitude"); latitudeStr != "" {
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid latitude format"))
			return
		}
		restaurantData.Latitude = latitude
	}

	// Validate required restaurant fields
	if restaurantData.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("restaurant_name is required"))
		return
	}

	// Handle image upload
	file, _, err := r.FormFile("restaurant_image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("restaurant_image is required"))
		return
	}
	defer file.Close()

	imageURL, err := utils.UploadImageToCloudinary(file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error uploading image: %v", err))
		return
	}
	restaurantData.Image = imageURL

	// Create restaurant with admin
	idRestaurant, token, err := h.store.CreateRestaurantWithAdmin(restaurantData, profileData)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.SendRestaurantAdminWelcomeEmail(profileData.Email, profileData.FirstName, profileData.LastName, profileData.Password, restaurantData.Name)
	if err != nil {
		log.Printf("Failed to send welcome email to %s: %v", profileData.Email, err)
	}
	utils.WriteJson(w, http.StatusCreated, map[string]interface{}{
		"message":      "Restaurant and admin created successfully",
		"idRestaurant": idRestaurant,
		"token":        token,
		"image":        imageURL,
	})
}

func (h *Handler) loginRestaurant(w http.ResponseWriter, r *http.Request) {
	var user types.UserLogin
	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.store.GetAdminByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("User", u)

	//!NOTE: compare the password
	if !utils.ComparePasswords([]byte(user.Password), []byte(u.Password)) {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Invalid Password"))
		return
	}
	isAssigned, _, err := h.store.VerifyAdminRestaurantAssignment(u.IdAdminRestaurant)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if !isAssigned {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("admin is no longer assigned to any restaurant"))
		return
	}

	//!NOTE: create a token
	token, err := utils.CreateRefreshToken(u.Id, u.Type)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]interface{}{"token": token, "user": u})
}

func (h *Handler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterAdmin
	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	idUser, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//!NOTE: Create refresht token

	token, err := utils.CreateRefreshToken(idUser, user.Type)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//!NOTE: Hash the password
	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if u != nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("user already exists"))
		return
	}

	//!NOTE: Create the user
	if err := h.store.CreateUser(user, idUser, token, string(hashedPassword)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	idAdmin, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	//!NOTE::Asign the role to the admin
	if user.Type == "adminRestaurant" {
		if err := h.store.CreateAdminRestaurant(idUser, idAdmin); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		//!NOTE : WHAT HAS BEEN DONE THE ADMIN CAN BE ADMIN FOR MANY RESTAURANTS BUT THE RESTAURANT CAN HAVE ONLY ONE ADMIN
		if user.IdActivitie != "" {
			if err := h.store.UpdateRestaurantAdmin(user.IdActivitie, idAdmin); err != nil {
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}
		}
	} else if user.Type == "adminActivity" {
		if err := h.store.CreateAdminActivity(idUser, idAdmin); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	} else {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid role"))
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token, "message": "Admin created successfully"})
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var user types.UserLogin
	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	// Add this check to prevent nil pointer dereference
	if u == nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}
	
	fmt.Println("User", u)

	//!NOTE: compare the password
	if !utils.ComparePasswords([]byte(user.Password), []byte(u.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	//!NOTE: create a token
	token, err := utils.CreateRefreshToken(u.Id, u.Type)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	isAdmin, idAdminActivity, err := h.store.IsClientAdminActivity(u.Id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"token":           token,
		"user":            u,
		"isAdminActivity": isAdmin,
	}

	if isAdmin {
		response["idAdminActivity"] = idAdminActivity
	}

	utils.WriteJson(w, http.StatusOK, response)
}

func (h *Handler) GetAllClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.store.GetAllClients()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, clients)
}

func (h *Handler) AssignClientToAdminActivity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IdClient string `json:"idClient"`
	}
	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.store.AssignClientToAdminActivity(req.IdClient)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Client assigned to admin activity successfully",
	})
}

func (h *Handler) signUpUser(w http.ResponseWriter, r *http.Request) {
	//!TODO: CHANGE THE TYPE SHOULD BE FIXED IN THE BACKEND NOT THE FRONT
	var user types.RegisterUser

	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	//!NOTE: Create an id
	idUser, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//!NOTE: Create refresht token

	token, err := utils.CreateRefreshToken(idUser, user.Role)
	// utils.MailSend()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//!NOTE: Hash the password
	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if u != nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("user already exists"))
		return
	}

	//!NOTE: Create the user
	if err := h.store.CreateUser(user, idUser, token, string(hashedPassword)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	idClient, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//!NOTE::Asign the role to the client

	if err := h.store.CreateClient(idUser, idClient, user.UserName); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.SendEmail(user.Email)

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token, "message": "Client created successfully"})
	// utils.SendSms()
}

func (h *Handler) SendRequestFriend(w http.ResponseWriter, r *http.Request) {
	var request types.SendRequestFriend
	if err := utils.ParseJson(r, &request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate the request
	if request.FromClient == "" || request.ToClient == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("sender and receiver IDs are required"))
		return
	}
	senderId, err := h.store.GetClientIdByUsername(request.FromClient)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("sender does not exist"))
		return

	}
	receiverId, err := h.store.GetClientIdByUsername(request.ToClient)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("receiver does not exist"))
		return
	}
	idFriendShip, err := utils.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.store.SendRequestFriend(idFriendShip, senderId, receiverId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Friend request sent successfully"})
}

func (h *Handler) AcceptFrienship(w http.ResponseWriter, r *http.Request) {
	var request types.AcceptFriendRequest
	if err := utils.ParseJson(r, &request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.AcceptRequestFriend(request.IdFriendship); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Friend Accepted successfully"})
}

func (h *Handler) GetFriendshipClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idClient, ok := vars["idClient"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idClient is required"))
		return
	}

	friendships, err := h.store.GetFriendshipRequested(idClient)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, friendships)
}

func (h *Handler) ClientInformation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idClient, ok := vars["idClient"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idClient is required"))
		return
	}

	u, err := h.store.GetClientInformation(idClient)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, u)
}

func (h *Handler) ClientInformationUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("idClient is required"))
		return
	}

	u, err := h.store.GetClientInformationUsername(username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, u)
}

func (h *Handler) GetUsername(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	if prefix == "" {
		http.Error(w, "prefix is required", http.StatusBadRequest)
		return
	}
	if prefix == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("username is required"))
		return
	}

	username, err := h.store.SearchUsersByUsernamePrefix(prefix)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]interface{}{"usernames": username})
}

func (h *Handler) CreateActivityWithAdmin(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form
    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, err)
        return
    }

    // Parse profile data
    var profileData types.ActivityAdminCreation
    profileData.FirstName = r.FormValue("firstName")
    profileData.LastName = r.FormValue("lastName")
    profileData.Email = r.FormValue("email")
    profileData.Phone = r.FormValue("phone")
    profileData.Address = r.FormValue("address")
    profileData.Password = r.FormValue("password")
    profileData.Type = r.FormValue("type")

    // Validate type
    if profileData.Type != "adminActivity" {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("type must be adminActivity"))
        return
    }

    // Parse activity data
    var activityData types.ActivityCreationWithAdmin
    activityData.Name = r.FormValue("activity_name")
    activityData.Description = r.FormValue("activity_description")
    activityData.IdTypeActivity = r.FormValue("activity_type")

    // Parse numeric fields
    if longitudeStr := r.FormValue("activity_longitude"); longitudeStr != "" {
        longitude, err := strconv.ParseFloat(longitudeStr, 64)
        if err != nil {
            utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid longitude format"))
            return
        }
        activityData.Longitude = longitude
    }

    if latitudeStr := r.FormValue("activity_latitude"); latitudeStr != "" {
        latitude, err := strconv.ParseFloat(latitudeStr, 64)
        if err != nil {
            utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid latitude format"))
            return
        }
        activityData.Latitude = latitude
    }

    // Handle image upload
    file, _, err := r.FormFile("activity_image")
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("activity image is required"))
        return
    }
    defer file.Close()

    imageURL, err := utils.UploadImageToCloudinary(file)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
        return
    }
    activityData.Image = imageURL

    // Create activity with admin
    idActivity, token, err := h.store.CreateActivityWithAdmin(activityData, profileData)
    if err != nil {
        if strings.Contains(err.Error(), "already exists") {
            utils.WriteError(w, http.StatusConflict, err)
            return
        }
        utils.WriteError(w, http.StatusInternalServerError, err)
        return
    }

    err = utils.SendActivityAdminWelcomeEmail(profileData.Email, profileData.FirstName, profileData.LastName, profileData.Password, activityData.Name)
    if err != nil {
        log.Printf("Failed to send welcome email to %s: %v", profileData.Email, err)
    }
    utils.WriteJson(w, http.StatusCreated, map[string]interface{}{
        "message":      "Activity and admin created successfully",
        "idActivity": idActivity,
        "token":        token,
        "image":        imageURL,
    })
}
