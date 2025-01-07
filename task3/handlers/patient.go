package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"gorm.io/gorm"

	"task3/models"
	"task3/tools"
)

func CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		tools.OperateUnsuccessfulResponse(w, "Bad request: Invalid JSON received", http.StatusBadRequest)
		return
	}

	if result := tools.DB.Create(&patient); result.Error != nil {
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(models.Response{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Patient created",
		Content: &patient,
	})
	if err != nil {
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Patient created:", patient)
}

func GetPatientByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var patient models.Patient

	result := tools.DB.First(&patient, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tools.OperateUnsuccessfulResponse(w, "Patient is not found", http.StatusNotFound)
		} else {
			tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(models.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Patient found",
		Content: &patient,
	})
	if err != nil {
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Patient received:", patient)
}

func GetPatientsList(w http.ResponseWriter, _ *http.Request) {
	var patients []models.Patient

	result := tools.DB.Find(&patients)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tools.OperateUnsuccessfulResponse(w, "No patients found", http.StatusNotFound)
		} else {
			tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(models.ListResponse[[]models.Patient]{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Patients list",
		Content: &patients,
	})
	if err != nil {
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Patients list retrieved:", patients)
}

func UpdatePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updateData models.Patient
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		tools.OperateUnsuccessfulResponse(w, "Bad request: Invalid JSON received", http.StatusBadRequest)
		return
	}

	var patient models.Patient
	if err := tools.DB.First(&patient, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tools.OperateUnsuccessfulResponse(w, "Patient is not found", http.StatusNotFound)
		} else {
			tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if err := tools.DB.Model(&patient).Updates(updateData).Error; err != nil {
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(models.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Patient updated",
		Content: &patient,
	})
	if err != nil {
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Patient updated:", patient)
}

func DeletePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result := tools.DB.Delete(&models.Patient{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tools.OperateUnsuccessfulResponse(w, "Patient is not found", http.StatusNotFound)
		} else {
			tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if result.RowsAffected == 0 {
		tools.OperateUnsuccessfulResponse(w, "Patient is not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(models.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Patient deleted",
	})
	if err != nil {
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Patient deleted with id: ", id)
}
