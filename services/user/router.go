package user

import (
	"fmt"
	// "log"
	"net/http"

	// "os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"

	// "github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
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
	router.HandleFunc("/login", h.loginUser).Methods("POST")
	router.HandleFunc("/getfriendship/{idClient}", h.GetFriendshipClient).Methods("GET")
	router.HandleFunc("/acceptfriendship", h.AcceptFrienship).Methods("POST")
	router.HandleFunc("/signup", h.signUpUser).Methods("POST")
	router.HandleFunc("/sendrequest", h.SendRequestFriend).Methods("POST")
	router.HandleFunc("/auth/{provider}", beginAuth).Methods("GET")
	router.HandleFunc("/auth/{provider}/callback", completeAuth).Methods("GET")
	router.HandleFunc("/auth/logout", logout).Methods("GET")
	router.HandleFunc("/clientinformation/{idClient}", h.ClientInformation).Methods("GET")
	// admin
}

const (
	key    = "key"
	MaxAge = 86400 * 30
	isProd = false
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// googleClientId := configs.Env.GoogleClientId
	// googleClientSecret := configs.Env.GoogleClientSecret

	// fmt.Println("Google Client ID:", googleClientId)     // Debugging
	// fmt.Println("Google Client ID:", googleClientSecret) // Debugging

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	// goth.UseProviders(
	// 	google.New(googleClientId, googleClientSecret,
	// 		"http://localhost:8080/auth/google/callback"),
	// )

	fmt.Println("Google Provider Registered")
}

func beginAuth(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

// Step 2: Handle Google Callback
func completeAuth(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// You can now use user.Email, user.Name, etc.
	fmt.Fprintf(w, "Welcome, %s! Your email is %s.", user.Name, user.Email)
}

// Step 3: Logout Handler
func logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	fmt.Fprintln(w, "Logged out successfully!")
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
	utils.WriteJson(w, http.StatusOK, map[string]interface{}{"token": token, "user": u})
}

func (h *Handler) signUpUser(w http.ResponseWriter, r *http.Request) {
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
