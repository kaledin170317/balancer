package erros

import (
	"encoding/json"
	"net/http"
)

type DefaultAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(DefaultAPIError{
		Code:    code,
		Message: msg,
	})
}
