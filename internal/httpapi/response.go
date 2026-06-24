package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"project-budget-service/internal/domain"
)

type apiResponse struct {
	Data  any            `json:"data"`
	Error *errorPayload  `json:"error"`
	Meta  map[string]any `json:"meta"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiResponse{
		Data:  data,
		Error: nil,
		Meta:  map[string]any{"request_id": uuid.NewString()},
	})
}

func writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	payload := errorPayload{Code: string(domain.CodeInternal), Message: "internal server error"}
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		payload.Code = string(appErr.Code)
		payload.Message = appErr.Message
		switch appErr.Code {
		case domain.CodeValidation:
			status = http.StatusBadRequest
		case domain.CodeUnauthorized:
			status = http.StatusUnauthorized
		case domain.CodeForbidden:
			status = http.StatusForbidden
		case domain.CodeNotFound:
			status = http.StatusNotFound
		case domain.CodeConflict:
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}
	}
	slog.Warn("request failed", "status", status, "code", payload.Code, "message", payload.Message, "error", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiResponse{
		Data:  nil,
		Error: &payload,
		Meta:  map[string]any{"request_id": uuid.NewString()},
	})
}
