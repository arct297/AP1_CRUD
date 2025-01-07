package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"task3/handlers"
	"task3/tools"
)

func main() {
	tools.InitDatabaseClient()

	sqlDB, err := tools.DB.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB from GORM DB:", err)
	}
	defer sqlDB.Close()

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

	r.HandleFunc("/doctors", handlers.GetDoctorsList).Methods("GET")

	r.HandleFunc("/mailing", handlers.MakeMailing).Methods("POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
