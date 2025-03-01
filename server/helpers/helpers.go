package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type errorResponseStruct struct {
	Error      bool                 `json:"error"`
	Message    string               `json:"message"`
	Validation []ValidationResponse `json:"validation,omitempty"`
}

type ValidationResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, err error, status int) {
	fmt.Println(err)
	response := errorResponseStruct{
		Error:   true,
		Message: err.Error(),
	}

	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		response.Message = "invalid validation"
		status = http.StatusInternalServerError
	}

	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			response.Validation = append(response.Validation, ValidationResponse{
				// TODO: generate proper field name and message
				Field:   e.Field(),
				Message: e.Tag(),
			})
		}
	}

	// TODO; handle common sql errors like duplicate key

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)
}
