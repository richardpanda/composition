package utils

import (
	"encoding/json"
	"net/http"

	"github.com/richardpanda/composition/server/api/types"
)

func SetErrorResponse(w http.ResponseWriter, code int, message string) {
	respBody, _ := json.Marshal(types.ErrorResponseBody{Message: message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respBody)
}
