package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
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

	// Filter parameter
	filter := query.Get("filter")
	filterValue := query.Get("filter_value")
	filterFrom := query.Get("filter_from")
	filterTo := query.Get("filter_to")

	db := tools.DB // Use local db variable to avoid overwriting global DB

	if filter != "" {
		allowedFilters := map[string]bool{
			"specialization":   true,
			"experience_years": true,
			"gender":           true,
			"birthdate":        true,
		}

		if !allowedFilters[filter] {
			tools.OperateUnsuccessfulResponse(w, "Invalid filter field", http.StatusBadRequest)
			return
		}

		if filterValue != "" {
			log.Printf("Filtering by field: %s with value: %s", filter, filterValue)
			db = db.Where(fmt.Sprintf("%s = ?", filter), filterValue)
		} else if filterFrom != "" || filterTo != "" {
			if filterFrom != "" && filterTo != "" {
				log.Printf("Filtering by field: %s with range: %s to %s", filter, filterFrom, filterTo)
				db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter), filterFrom, filterTo)
			} else if filterFrom != "" {
				log.Printf("Filtering by field: %s with minimum value: %s", filter, filterFrom)
				db = db.Where(fmt.Sprintf("%s >= ?", filter), filterFrom)
			} else if filterTo != "" {
				log.Printf("Filtering by field: %s with maximum value: %s", filter, filterTo)
				db = db.Where(fmt.Sprintf("%s <= ?", filter), filterTo)
			}
		}
	}

	// Log request parameters
	log.Printf("GetDoctorsList called with parameters: sort=%s, order=%s, limit=%d, offset=%d, page=%d", sortField, sortOrder, limit, offset, page)

	// Query the database
	result := db.Order(fmt.Sprintf("%s %s", sortField, sortOrder)).
		Offset(offset).
		Limit(limit).
		Find(&doctors)

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