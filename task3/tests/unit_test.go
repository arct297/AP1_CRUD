package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreatePatientHandler(t *testing.T) {
    tests := []struct {
        name        string
        patient     Patient
        wantStatus  int
        wantMessage string
    }{
        {
            name: "Valid Patient",
            patient: Patient{
                Name:    "John Doe",
                Age:     30,
                Gender:  "M",
                Contact: "+71234567890",
                Address: "Some Address",
            },
            wantStatus:  http.StatusCreated,
            wantMessage: "Patient created successfully!",
        },
        {
            name: "Invalid Age (negative)",
            patient: Patient{
                Name:    "John Doe",
                Age:     -5,
                Gender:  "M",
                Contact: "+71234567890",
            },
            wantStatus:  http.StatusBadRequest,
            wantMessage: "Age must be a positive number",
        },
        {
            name: "Invalid Gender",
            patient: Patient{
                Name:    "John Doe",
                Age:     30,
                Gender:  "X",
                Contact: "+71234567890",
            },
            wantStatus:  http.StatusBadRequest,
            wantMessage: "Gender must be either 'M' or 'F'",
        },
        {
            name: "Invalid Name (too short)",
            patient: Patient{
                Name:    "A",
                Age:     30,
                Gender:  "M",
                Contact: "+71234567890",
            },
            wantStatus:  http.StatusBadRequest,
            wantMessage: "Name must be between 2 and 25 characters",
        },
        {
            name: "Invalid Contact (wrong format)",
            patient: Patient{
                Name:    "John Doe",
                Age:     30,
                Gender:  "M",
                Contact: "1234567890",
            },
            wantStatus:  http.StatusBadRequest,
            wantMessage: "Contact must start with '+7' and be at least 12 characters long",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.patient)
            req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
            rr := httptest.NewRecorder()

            CreatePatientHandler(rr, req)

            if status := rr.Code; status != tt.wantStatus {
                t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
            }

            var response map[string]interface{}
            json.Unmarshal(rr.Body.Bytes(), &response)
            
            if msg := response["message"]; msg != tt.wantMessage {
                t.Errorf("handler returned unexpected message: got %v want %v", msg, tt.wantMessage)
            }
        })
    }
}

type Patient struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Gender  string `json:"gender"`
    Contact string `json:"contact"`
    Address string `json:"address"`
}

func CreatePatientHandler(w http.ResponseWriter, r *http.Request) {
    var patient Patient
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&patient); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if patient.Age <= 0 {
        respondWithError(w, http.StatusBadRequest, "Age must be a positive number")
        return
    }

    if patient.Gender != "M" && patient.Gender != "F" {
        respondWithError(w, http.StatusBadRequest, "Gender must be either 'M' or 'F'")
        return
    }

    if len(patient.Name) < 2 || len(patient.Name) > 25 {
        respondWithError(w, http.StatusBadRequest, "Name must be between 2 and 25 characters")
        return
    }

    if !validContact(patient.Contact) {
        respondWithError(w, http.StatusBadRequest, "Contact must start with '+7' and be at least 12 characters long")
        return
    }

    respondWithJSON(w, http.StatusCreated, map[string]string{
        "message": "Patient created successfully!",
    })
}

func validContact(contact string) bool {
    return len(contact) >= 12 && contact[:2] == "+7"
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}