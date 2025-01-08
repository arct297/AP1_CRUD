package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
<<<<<<< HEAD

	"github.com/sirupsen/logrus"
=======
	"gorm.io/gorm"
>>>>>>> 37d16bdc2759e7cf5f6f06b0c06433312db102ad

	"task3/logger"
	"task3/models"
	"task3/tools"
)

func GetDoctorsList(w http.ResponseWriter, r *http.Request) {
	var doctors []models.Doctor

	query := r.URL.Query()

	logger.Log.WithFields(logrus.Fields{
		"action": "get_doctors_list",
		"method": r.Method,
		"query":  query.Encode(),
	}).Info("Received request")

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
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		logger.Log.WithFields(logrus.Fields{
			"action": "get_doctors_list",
			"error":  "Invalid limit parameter",
		}).Warn("Using default limit")
		limit = 10
	}

	// Offset and Page parameter
	offsetStr := query.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		logger.Log.WithFields(logrus.Fields{
			"action": "get_doctors_list",
			"error":  "Invalid offset parameter",
		}).Warn("Using default offset")
		offset = 0
	}

	pageStr := query.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err == nil && page > 1 {
		offset = (page - 1) * limit
	} else if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"action": "get_doctors_list",
			"error":  "Invalid page parameter",
		}).Warn("Defaulting to page 1")
		page = 1
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
<<<<<<< HEAD
			logger.Log.WithFields(logrus.Fields{
				"action": "get_doctors_list",
				"filter": filter,
			}).Error("Invalid filter field")
=======
>>>>>>> 37d16bdc2759e7cf5f6f06b0c06433312db102ad
			tools.OperateUnsuccessfulResponse(w, "Invalid filter field", http.StatusBadRequest)
			return
		}

		if filterValue != "" {
<<<<<<< HEAD
			db = db.Where(fmt.Sprintf("%s = ?", filter), filterValue)
		} else if filterFrom != "" || filterTo != "" {
			if filterFrom != "" && filterTo != "" {
				db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter), filterFrom, filterTo)
			} else if filterFrom != "" {
				db = db.Where(fmt.Sprintf("%s >= ?", filter), filterFrom)
			} else if filterTo != "" {
=======
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
>>>>>>> 37d16bdc2759e7cf5f6f06b0c06433312db102ad
				db = db.Where(fmt.Sprintf("%s <= ?", filter), filterTo)
			}
		}
	}
<<<<<<< HEAD
=======

	// Log request parameters
	log.Printf("GetDoctorsList called with parameters: sort=%s, order=%s, limit=%d, offset=%d, page=%d", sortField, sortOrder, limit, offset, page)
>>>>>>> 37d16bdc2759e7cf5f6f06b0c06433312db102ad

	// Query the database
	result := db.Order(fmt.Sprintf("%s %s", sortField, sortOrder)).
		Offset(offset).
		Limit(limit).
		Find(&doctors)

	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"action": "get_doctors_list",
			"error":  result.Error.Error(),
		}).Error("Failed to fetch doctors")
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(doctors) == 0 {
		logger.Log.WithFields(logrus.Fields{
			"action": "get_doctors_list",
			"filter": filter,
		}).Warn("No doctors match the criteria")
		tools.OperateUnsuccessfulResponse(w, "No doctors match the filter", http.StatusNotFound)
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"action": "get_doctors_list",
		"count":  len(doctors),
	}).Info("Fetched doctors successfully")

	// Respond with the list of doctors
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	meta := models.ListResponseMeta{
		Limit:  limit,
		Offset: offset,
		Page:   page,
		Total:  len(doctors),
	}

	err = json.NewEncoder(w).Encode(models.ListResponse[[]models.Doctor]{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Doctors list",
		Content: &doctors,
		Meta:    meta,
	})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"action": "get_doctors_list",
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		}).Error("Failed to encode response to JSON")
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"action": "get_doctors_list",
		"status": http.StatusOK,
		"count":  len(doctors),
	}).Info("Response sent successfully")
}
