package main

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	_ "github.com/lib/pq"

	"task3/handlers"
	"task3/logger"
	"task3/tools"
)

// Define a global rate limiter
var limiter = rate.NewLimiter(1, 3) // 1 request per second with a burst of 3 requests

// Rate-limiting middleware
func rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			// If rate limit is exceeded, respond with 429 Too Many Requests
			w.Header().Set("Retry-After", time.Now().Add(limiter.Reserve().Delay()).Format(time.RFC1123))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize database client
	tools.InitDatabaseClient()

	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Logger initialization failed: %v", err)
	}

	// Get SQL database instance from GORM
	sqlDB, err := tools.DB.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB from GORM DB:", err)
	}
	defer sqlDB.Close()

	// Initialize router
	r := mux.NewRouter()

	// Serve static files like HTML, CSS, JS
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Return index.html on root path
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	}).Methods("GET")

	// Patient API routes
	r.HandleFunc("/patients", handlers.CreatePatient).Methods("POST")        // Create patient
	r.HandleFunc("/patients/{id}", handlers.GetPatientByID).Methods("GET")   // Get patient by ID
	r.HandleFunc("/patients", handlers.GetPatientsList).Methods("GET")       // Get patients list
	r.HandleFunc("/patients/{id}", handlers.UpdatePatient).Methods("PUT")    // Update patient by ID
	r.HandleFunc("/patients/{id}", handlers.DeletePatient).Methods("DELETE") // Delete patient by ID

	// Doctor API route
	r.HandleFunc("/doctors", handlers.GetDoctorsList).Methods("GET")

	// Mailing API route
	r.HandleFunc("/mailing", handlers.MakeMailing).Methods("POST")

	// Apply rate-limiting middleware
	rateLimitedRouter := rateLimiterMiddleware(r)

	// Start the HTTP server
	http.Handle("/", rateLimitedRouter)
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
