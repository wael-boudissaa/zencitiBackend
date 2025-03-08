package user

import (
	"context"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/zencitiBackend/services/auth"
	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig = &oauth2.Config{
	ClientID:     "YOUR_CLIENT_ID",
	ClientSecret: "YOUR_CLIENT_SECRET",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
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
func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("randomstate")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Step 2: Handle Google Callback
func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	client := googleOAuthConfig.Client(ctx, token)

	// Fetch user info from Google API
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	fmt.Fprintln(w, "User info retrieved successfully!")
}
