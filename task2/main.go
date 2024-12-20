package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDatabaseClient() {
	dsn := "user=username password=password dbname=clinic_db sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

}

type Response struct {
	Code    int      `json:"code"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Content *Patient `json:"content,omitempty"`
}

type Patient struct {
	gorm.Model
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Contact string `json:"contact"`
	Address string `json:"address"`
}

func operateUnsuccessfulResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Code:    statusCode,
		Status:  "error",
		Message: message,
	})
}

func createPatient(w http.ResponseWriter, r *http.Request) {
	var patient Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		operateUnsuccessfulResponse(w, "Bad request: Invalid JSON received", http.StatusBadRequest)
		return
	}

	if result := db.Create(&patient); result.Error != nil {
		operateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(Response{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Patient created",
		Content: &patient,
	})
	if err != nil {
		operateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Patient created:", patient)
}

func getPatientByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var patient Patient

	result := db.First(&patient, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			operateUnsuccessfulResponse(w, "Patient is not found", http.StatusNotFound)
		} else {
			operateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Patient found",
		Content: &patient,
	})
	if err != nil {
		operateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Patient received:", patient)
}

func updatePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updateData Patient
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		operateUnsuccessfulResponse(w, "Bad request: Invalid JSON received", http.StatusBadRequest)
		return
	}

	var patient Patient
	if err := db.First(&patient, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			operateUnsuccessfulResponse(w, "Patient is not found", http.StatusNotFound)
		} else {
			operateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if err := db.Model(&patient).Updates(updateData).Error; err != nil {
		operateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Patient updated",
		Content: &patient,
	})
	if err != nil {
		operateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Patient updated:", patient)
}

func deletePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result := db.Delete(&Patient{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			operateUnsuccessfulResponse(w, "Patient is not found", http.StatusNotFound)
		} else {
			operateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if result.RowsAffected == 0 {
		operateUnsuccessfulResponse(w, "Patient is not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Patient deleted",
	})
	if err != nil {
		operateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Patient deleted with id: ", id)
}

func main() {
	initDatabaseClient()

	sqlDB, err := db.DB()
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
	r.HandleFunc("/patients", createPatient).Methods("POST")        // Create patient
	r.HandleFunc("/patients/{id}", getPatientByID).Methods("GET")   // Get patient by ID
	r.HandleFunc("/patients/{id}", updatePatient).Methods("PUT")    // Update patient by ID
	r.HandleFunc("/patients/{id}", deletePatient).Methods("DELETE") // Delete patient by ID

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
