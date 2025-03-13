package user

import (
	"fmt"
	"net/http"
	// "os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/wael-boudissaa/zencitiBackend/configs"
	"github.com/wael-boudissaa/zencitiBackend/services/auth"
	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
)

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

	googleClientId := configs.Env.GoogleClientId
	googleClientSecret := configs.Env.GoogleClientSecret

	fmt.Println("Google Client ID:", googleClientId)     // Debugging
	fmt.Println("Google Client ID:", googleClientSecret) // Debugging

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret,
			"http://localhost:8080/auth/google/callback"),
	)

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

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.loginUser).Methods("POST")
	router.HandleFunc("/signup", h.signUpUser).Methods("POST")
	router.HandleFunc("/auth/{provider}", beginAuth).Methods("GET")
	router.HandleFunc("/auth/{provider}/callback", completeAuth).Methods("GET")
	router.HandleFunc("/auth/logout", logout).Methods("GET")
	//admin
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var user types.UserLogin

	// fmt.Println("user called")

	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !auth.ComparePasswords([]byte(user.Password), []byte(u.Password)) {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Invalid Password"))
		return
	}

	token, err := auth.CreateRefreshToken(*&u.Id)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) signUpUser(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUser

	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	//!NOTE: CREATE AND ID
	idUser, err := auth.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//!NOTE: CREATE REFRESH TOKEN
	token, err := auth.CreateRefreshToken(idUser)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//!NOTE: HASH PASSWORD
	hashedPassword, err := auth.HashedPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	//!NOTE: CREATE USER
	if err := h.store.CreateUser(user, idUser, token, string(hashedPassword)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token, "message": "User created successfully"})

}
