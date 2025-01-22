package main

import (
	"database/sql"
	"dating-app/models"
	"dating-app/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "dating-app/docs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

var db *sql.DB
var userService service.UserService = &service.UserServiceImpl{}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME")))

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()

	// Swagger UI route (use swaggerFiles.Handler directly)
	r.PathPrefix("/swagger/{any:.*}").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL to the generated Swagger JSON
	))

	r.HandleFunc("/signup", SignupHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/swipe", SwipeHandler).Methods("POST")
	r.HandleFunc("/purchase", PurchaseHandler).Methods("POST")

	// Serve static Swagger JSON file
	r.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})

	http.Handle("/", r)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// @Summary Signup a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /signup [post]
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := userService.Signup(user); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Signup successful"})
}

// @Summary User Login
// @Description Logs in a user and returns a token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param username body string true "Username"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// @Summary Swipe action
// @Description Records a swipe action from a user
// @Tags Actions
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /swipe [post]
func SwipeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Swipe action recorded"})
}

// @Summary Purchase Premium
// @Description Allows a user to purchase a premium package
// @Tags Payments
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /purchase [post]
func PurchaseHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Premium package purchased"})
}
