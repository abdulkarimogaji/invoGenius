package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/abdulkarimogaji/invoGenius/db"
	"github.com/abdulkarimogaji/invoGenius/server/helpers"
	"github.com/abdulkarimogaji/invoGenius/services/password"
	"github.com/abdulkarimogaji/invoGenius/utils/types"
	"github.com/go-playground/validator/v10"
)

type createUserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Role      string `json:"role" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
}
type createUserResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

type getCustomersResponse struct {
	Error           bool                 `json:"error"`
	Message         string               `json:"message"`
	Customers       []db.GetCustomersRow `json:"customers"`
	DefaultCurrency string               `json:"default_currency"`
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
		FirstName: requestBody.FirstName,
		LastName:  requestBody.LastName,
		Role:      requestBody.Role,
		Email:     requestBody.Email,
		Password:  types.JSONNullString{String: hash, Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Phone:     requestBody.Phone,
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

func (h *Handler) GetCustomers(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	cId, err := strconv.Atoi(query.Get("customer_id"))
	customerId := sql.NullInt32{Int32: int32(cId), Valid: err == nil}

	filters := db.GetCustomersParams{
		CustomerID: customerId,
		FirstName:  query.Get("first_name"),
		LastName:   query.Get("last_name"),
		Email:      query.Get("email"),
		Phone:      query.Get("phone"),
		SortBy:     query.Get("sort_by"),
		SortOrder:  query.Get("sort_order"),
	}

	customers, err := db.DB.GetCustomers(r.Context(), filters)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if customers == nil {
		customers = []db.GetCustomersRow{}
	}

	// get default currency
	// TODO: map the default currency in every customer
	defaultCurrency, err := db.DB.GetDefaultCurrency(r.Context())
	if err == sql.ErrNoRows {
		helpers.ErrorResponse(w, fmt.Errorf("default currency not set"), http.StatusInternalServerError)
		return
	}
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := getCustomersResponse{
		Error:           false,
		Message:         "customers fetched successfully",
		Customers:       customers,
		DefaultCurrency: defaultCurrency,
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
