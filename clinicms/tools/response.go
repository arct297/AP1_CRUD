package tools

import (
	"encoding/json"
	"net/http"

	"clinicms/models"
)

func OperateUnsuccessfulResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.Response{
		Code:    statusCode,
		Status:  "error",
		Message: message,
	})
}
