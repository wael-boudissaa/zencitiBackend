package user

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
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
