package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Initialize the database connection
func initDB() {
	var err error
	connStr := "user=username password=password dbname=clinic_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

// Patient struct to represent the patient data
type Patient struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Contact string `json:"contact"`
	Address string `json:"address"`
}

// Create patient
func createPatient(w http.ResponseWriter, r *http.Request) {
    var patient Patient
    err := json.NewDecoder(r.Body).Decode(&patient)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Логирование полученных данных для отладки
    log.Printf("Received Patient Data: %+v", patient)

    query := `INSERT INTO patients (name, age, gender, contact, address) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    err = db.QueryRow(query, patient.Name, patient.Age, patient.Gender, patient.Contact, patient.Address).Scan(&patient.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(patient); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


// Read (get all) patients
func getPatients(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM patients")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var patients []Patient
	for rows.Next() {
		var patient Patient
		if err := rows.Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender, &patient.Contact, &patient.Address); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		patients = append(patients, patient)
	}
	json.NewEncoder(w).Encode(patients)
}

// Update patient by ID
func updatePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var patient Patient
	err = json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE patients SET name=$1, age=$2, gender=$3, contact=$4, address=$5 WHERE id=$6`
	_, err = db.Exec(query, patient.Name, patient.Age, patient.Gender, patient.Contact, patient.Address, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	patient.ID = id
	json.NewEncoder(w).Encode(patient)
}

// Delete patient by ID
func deletePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM patients WHERE id=$1`
	_, err = db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// Serve static files like HTML, CSS, JS
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Return index.html on root path
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	}).Methods("GET")

	// Patient API routes
	r.HandleFunc("/patients", createPatient).Methods("POST")        // Create patient
	r.HandleFunc("/patients", getPatients).Methods("GET")           // Get all patients
	r.HandleFunc("/patients/{id}", updatePatient).Methods("PUT")    // Update patient by ID
	r.HandleFunc("/patients/{id}", deletePatient).Methods("DELETE") // Delete patient by ID

	http.Handle("/", r)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
