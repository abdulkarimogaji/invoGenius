package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/abdulkarimogaji/invoGenius/db"
	"github.com/abdulkarimogaji/invoGenius/server/helpers"
	"github.com/abdulkarimogaji/invoGenius/services/password"
	"github.com/go-playground/validator/v10"
)

type createUserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}
type createUserResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

var validate *validator.Validate

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody createUserRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal([]byte(body), &requestBody)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	validate = validator.New()
	err = validate.Struct(&requestBody)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	hash, err := password.HashPassword(requestBody.Password)
	if err != nil {
		helpers.ErrorResponse(w, fmt.Errorf("failed to hash password: %v", err), http.StatusInternalServerError)
		return
	}

	result, err := db.DB.CreateUser(r.Context(), db.CreateUserParams{
		FirstName: sql.NullString{String: requestBody.FirstName, Valid: true},
		LastName:  sql.NullString{String: requestBody.LastName, Valid: true},
		Role:      "staff",
		Email:     requestBody.Email,
		Password:  sql.NullString{String: hash, Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := createUserResponse{
		Error:   false,
		Message: "user created successfully",
		UserID:  int(userID),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
