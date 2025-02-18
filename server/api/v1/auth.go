package v1

import (
	"encoding/json"
	"net/http"

	"github.com/abdulkarimogaji/invoGenius/services/token"
)

type loginResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var response loginResponse
	tokenString, err := token.CreateToken("123")
	if err != nil {
		response.Error = true
		response.Message = err.Error()
	} else {
		response.Error = false
		response.Message = "login successful"
		response.Token = tokenString
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
