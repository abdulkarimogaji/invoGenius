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
	"github.com/abdulkarimogaji/invoGenius/middleware"
	"github.com/abdulkarimogaji/invoGenius/server/helpers"
	"github.com/abdulkarimogaji/invoGenius/utils/types"
	"github.com/go-playground/validator/v10"
)

type createInvoiceRequest struct {
	UserID    int              `json:"user_id" validate:"omitempty,required,number"`
	Amount    int              `json:"amount" validate:"required,number"`
	Type      string           `json:"type" validate:"required"`
	IssuedAt  types.CustomDate `json:"issued_at" validate:"required"`
	FromDate  types.CustomDate `json:"from_date" validate:"required"`
	UntilDate types.CustomDate `json:"until_date" validate:"required"`
	Currency  string           `json:"currency"`
}
type createInvoiceResponse struct {
	Error     bool   `json:"error"`
	Message   string `json:"message"`
	InvoiceID int    `json:"invoice_id"`
}

type getInvoicesResponse struct {
	Error    bool                `json:"error"`
	Message  string              `json:"message"`
	Invoices []db.GetInvoicesRow `json:"invoices"`
}

func (h *Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var requestBody createInvoiceRequest

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

	_, err = db.DB.GetUserByID(r.Context(), int32(requestBody.UserID))
	if err == sql.ErrNoRows {
		helpers.ErrorResponse(w, fmt.Errorf("user not found"), http.StatusNotFound)
		return
	}
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	settingsResult, err := db.DB.GetInvoiceSettings(r.Context())
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	invoiceSettings := make(map[string]string)
	for _, value := range settingsResult {
		invoiceSettings[value.SettingKey] = value.SettingValue
	}

	if requestBody.Currency == "" {
		requestBody.Currency = invoiceSettings["currency"]
	}

	vat, err := strconv.ParseFloat(invoiceSettings["vat"], 64)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	deadlineDays, err := strconv.Atoi(invoiceSettings["deadline_days"])
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	result, err := db.DB.CreateInvoice(r.Context(), db.CreateInvoiceParams{
		UserID:    int32(requestBody.UserID),
		Amount:    float64(requestBody.Amount),
		Vat:       vat,
		Type:      requestBody.Type,
		IssuedAt:  requestBody.IssuedAt.ToTime(),
		FromDate:  requestBody.FromDate.ToTime(),
		UntilDate: requestBody.UntilDate.ToTime(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Currency:  requestBody.Currency,
		Deadline:  time.Now().Add(time.Hour * 24 * time.Duration(deadlineDays)),
	})

	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	invoiceID, err := result.LastInsertId()
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		helpers.ErrorResponse(w, fmt.Errorf("invalid user id"), http.StatusInternalServerError)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(w, fmt.Errorf("invalid user id"), http.StatusInternalServerError)
		return
	}

	// add invoice activity;
	db.DB.CreateInvoiceActivity(r.Context(), db.CreateInvoiceActivityParams{
		UserID:     int32(userID),
		InvoiceID:  int32(invoiceID),
		ActionType: "create_invoice",
		ResourceID: int32(invoiceID),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Attachment: types.JSONNullString{Valid: false},
	})

	response := createInvoiceResponse{
		Error:     false,
		Message:   "invoice created successfully",
		InvoiceID: int(invoiceID),
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

func (h *Handler) GetInvoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := db.DB.GetInvoices(r.Context())
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := getInvoicesResponse{
		Error:    false,
		Message:  "invoice created successfully",
		Invoices: invoices,
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
