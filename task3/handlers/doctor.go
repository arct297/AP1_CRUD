package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"strconv"

	"gorm.io/gorm"


	"task3/models"
	"task3/tools"
)

func GetDoctorsList(w http.ResponseWriter, r *http.Request) {
	var doctors []models.Doctor

	query := r.URL.Query()

	// Sorting parameters
	sortField := query.Get("sort")
	sortOrder := query.Get("order")

	allowedSortingFields := map[string]bool{
		"id":               true,
		"first_name":       true,
		"last_name":        true,
		"gender":           true,
		"birthdate":        true,
		"email":            true,
		"phone_number":     true,
		"experience_years": true,
		"specialization":   true,
	}

	if sortField == "" || !allowedSortingFields[sortField] {
		sortField = "id"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// Limit parameter
	limitStr := query.Get("limit")
	limit, limitErr := strconv.Atoi(limitStr)
	if limitErr != nil || limit <= 0 {
		limit = 10
	}

	// Offset and Page parameter
	offsetStr := query.Get("offset")
	offset, offsetErr := strconv.Atoi(offsetStr)
	if offsetErr != nil || offset < 0 {
		offset = 0
	}

	pageStr := query.Get("page")
	page, pageErr := strconv.Atoi(pageStr)
	if pageErr == nil && page > 1 {
		offset = (page - 1) * limit
	}

	// Log request parameters
	log.Printf("GetDoctorsList called with parameters: sort=%s, order=%s, limit=%d, offset=%d, page=%d", sortField, sortOrder, limit, offset, page)

	// Query the database
	result := tools.DB.Order(fmt.Sprintf("%s %s", sortField, sortOrder)).Offset(offset).Limit(limit).Find(&doctors)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tools.OperateUnsuccessfulResponse(w, "No doctors found", http.StatusNotFound)
		} else {
			tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the list of doctors
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	meta := models.ListResponseMeta{
		Limit:  limit,
		Offset: offset,
		Page:   page,
		Total:  len(doctors),
	}

	err := json.NewEncoder(w).Encode(models.ListResponse[[]models.Doctor]{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Doctors list",
		Content: &doctors,
		Meta:    meta,
	})
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}