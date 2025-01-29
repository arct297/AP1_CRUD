package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"task3/handlers"
	"task3/logger"
	"task3/models"
	"task3/tools"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	// load PostgreSQL connection parameters (adjust these for your environment)
	logger.Log = logrus.New()
	dsn := "host=localhost user=postgres password=1111 dbname=clinic_db port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get database instance: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// replace the global DB instance
	tools.DB = db
	return db
}

func CleanupTestDB(t *testing.T, db *gorm.DB) {
	// wrap in a transaction that we'll roll back
	tx := db.Begin()
	defer tx.Rollback()

	// clean up test data
	if err := tx.Exec("DELETE FROM doctors").Error; err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}
}

func seedDoctors(db *gorm.DB, count int, specialization string) {
	for i := 1; i <= count; i++ {
		db.Create(&models.Doctor{
			FirstName:       "Doctor",
			LastName:        fmt.Sprintf("Test%d", i),
			Gender:          "M",
			Birthdate:       time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
			Email:          fmt.Sprintf("doctor%d@test.com", i),
			PhoneNumber:    fmt.Sprintf("123-456-78%02d", i),
			ExperienceYears: i,
			Specialization:  specialization,
		})
	}
}

func TestGetDoctorsList_DefaultParams(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)

	// start transaction
	tx := db.Begin()
	defer tx.Rollback()

	// seed test data
	seedDoctors(tx, 15, "General")

	// execute request
	req := httptest.NewRequest("GET", "/doctors", nil)
	w := httptest.NewRecorder()
	handlers.GetDoctorsList(w, req)

	// validate response
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response models.ListResponse[[]models.Doctor]
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	// check default behavior
	assert.Equal(t, 10, len(*response.Content), "Expected 10 doctors by default")
	assert.Equal(t, 10, response.Meta.Limit)
	assert.Equal(t, 0, response.Meta.Offset)
	assert.Equal(t, 1, response.Meta.Page)

	// ensure sorting by ID ascending
	doctors := *response.Content
	for i := 1; i < len(doctors); i++ {
		assert.True(t, doctors[i-1].ID < doctors[i].ID, "Doctors not sorted by ID ascending")
	}
}

func TestGetDoctorsList_PaginationAndSorting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)

	// start transaction
	tx := db.Begin()
	defer tx.Rollback()

	// seed test data
	for i := 1; i <= 12; i++ {
		tx.Create(&models.Doctor{
			ExperienceYears: 13 - i, // Create descending years for sorting test
			Specialization:  "Cardiology",
		})
	}

	// execute request
	req := httptest.NewRequest("GET", "/doctors?sort=experience_years&order=desc&page=2&limit=5", nil)
	w := httptest.NewRecorder()
	handlers.GetDoctorsList(w, req)

	// validate response
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response models.ListResponse[[]models.Doctor]
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	// check pagination and sorting
	assert.Equal(t, 5, len(*response.Content), "Expected 5 doctors on page 2")
	assert.Equal(t, 5, response.Meta.Limit)
	assert.Equal(t, 5, response.Meta.Offset) // (page=2-1) * limit=5 = 5
	assert.Equal(t, 2, response.Meta.Page)

	// verify descending order of experience_years (17, 16, 16, 15, 15)
	doctors := *response.Content
	expectedYears := []int{17, 16, 16, 15, 15}
	for i, doctor := range doctors {
		assert.Equal(t, expectedYears[i], doctor.ExperienceYears, "Incorrect sorting order")
	}
}

func TestGetDoctorsList_FilterBySpecialization(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)

	// start transaction
	tx := db.Begin()
	defer tx.Rollback()

	// seed test data
	seedDoctors(tx, 3, "Pediatrics")
	seedDoctors(tx, 7, "Cardiology")

	// execute request
	req := httptest.NewRequest("GET", "/doctors?filter=specialization&filter_value=Pediatrics", nil)
	w := httptest.NewRecorder()
	handlers.GetDoctorsList(w, req)

	// validate response
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response models.ListResponse[[]models.Doctor]
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	// check filtering
	assert.Equal(t, 3, len(*response.Content), "Expected 3 Pediatric doctors")
	
	// ensure all returned doctors have the correct specialization
	for _, doctor := range *response.Content {
		assert.Equal(t, "Pediatrics", doctor.Specialization, "Incorrect specialization filter")
	}
}