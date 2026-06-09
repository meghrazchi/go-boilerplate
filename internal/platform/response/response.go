package response

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, statusCode int, message string, data any) {
	JSON(w, statusCode, Envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, statusCode int, message string, errs any) {
	JSON(w, statusCode, Envelope{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

func Paginated(w http.ResponseWriter, statusCode int, message string, data any, meta PaginationMeta) {
	JSON(w, statusCode, Envelope{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
