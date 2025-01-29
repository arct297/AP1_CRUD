package int_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Mailing struct {
	ID             uint   `gorm:"primaryKey"`
	ReceivingGroup string `gorm:"size:255;not null"`
	Topic          string `gorm:"size:255;not null"`
	Message        string `gorm:"type:text;not null"`
}

type MailingRequest struct {
	ReceivingGroup string `json:"receiving_group"`
	Topic          string `json:"topic"`
	Message        string `json:"message"`
}

func MakeMailing(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request MailingRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.ReceivingGroup != "debarbiest@gmail.com" {
			http.Error(w, "Invalid receiving group", http.StatusBadRequest)
			return
		}

		mailing := Mailing{
			ReceivingGroup: request.ReceivingGroup,
			Topic:          request.Topic,
			Message:        request.Message,
		}
		if err := db.Create(&mailing).Error; err != nil {
			http.Error(w, "Failed to save mailing", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

type MailingIntegrationTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *MailingIntegrationTestSuite) SetupTest() {
	dsn := "host=localhost user=postgres password=1111 dbname=test_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		suite.T().Fatal("Не удалось подключиться к тестовой базе данных:", err)
	}
	suite.db = db

	if err := suite.db.AutoMigrate(&Mailing{}); err != nil {
		suite.T().Fatal("Не удалось выполнить миграцию:", err)
	}

	suite.db.Exec("TRUNCATE TABLE mailings RESTART IDENTITY")
}

func (suite *MailingIntegrationTestSuite) TestMailingAPI() {
	mailingData := MailingRequest{
		ReceivingGroup: "debarbiest@gmail.com",
		Topic:          "Test Topic1",
		Message:        "Test Message1",
	}
	body, _ := json.Marshal(mailingData)

	req := httptest.NewRequest(http.MethodPost, "/mailing", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := MakeMailing(suite.db)
	handler(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var count int64
	suite.db.Table("mailings").Where("receiving_group = ?", "debarbiest@gmail.com").Count(&count)
	assert.Equal(suite.T(), int64(1), count, "Запись в базе данных отсутствует")
}

func TestMailingIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(MailingIntegrationTestSuite))
}

func main() {
	dsn := "host=localhost user=postgres password=1111 dbname=real_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Не удалось подключиться к базе данных:", err)
		return
	}

	db.AutoMigrate(&Mailing{})

	http.HandleFunc("/mailing", MakeMailing(db))
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
