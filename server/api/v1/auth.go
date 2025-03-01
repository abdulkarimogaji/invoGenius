package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/abdulkarimogaji/invoGenius/db"
	"github.com/abdulkarimogaji/invoGenius/server/helpers"
	"github.com/abdulkarimogaji/invoGenius/services/password"
	"github.com/abdulkarimogaji/invoGenius/services/token"
	"github.com/go-playground/validator/v10"
)

type loginResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Token   string `json:"token"`
	Role    string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var requestBody loginRequest

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

	user, err := db.DB.GetUserByEmail(r.Context(), requestBody.Email)
	if err == sql.ErrNoRows {
		helpers.ErrorResponse(w, fmt.Errorf("email and password does not match"), http.StatusUnauthorized)
		return
	}
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	valid := password.CheckPasswordHash(requestBody.Password, user.Password.String)
	if !valid {
		helpers.ErrorResponse(w, fmt.Errorf("email and password does not match"), http.StatusUnauthorized)
		return
	}

	tokenString, err := token.CreateToken(string(user.ID))
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := loginResponse{
		Error:   false,
		Message: "login successful",
		Token:   tokenString,
		Role:    user.Role,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Error {
		w.WriteHeader(http.StatusInternalServerError)
	} else {

		w.WriteHeader(http.StatusOK)
	}
	w.Write(jsonResponse)
}
