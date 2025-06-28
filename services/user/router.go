package user

import (
	"fmt"
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
	router.HandleFunc("/admin/login", h.loginRestaurant).Methods("POST")
	router.HandleFunc("/admin/create", h.CreateAdmin).Methods("POST")
}

const (
	key    = "key"
	MaxAge = 86400 * 30
	isProd = false
)

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
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("User", u)

	//!NOTE: compare the password
	if !utils.ComparePasswords([]byte(user.Password), []byte(u.Password)) {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Invalid Password"))
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
