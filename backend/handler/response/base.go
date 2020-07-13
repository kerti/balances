package response

import (
	"encoding/json"
	"net/http"
	"strings"
)

// BaseResponse is the base object of all responses
type BaseResponse struct {
	Data    *interface{} `json:"data,omitempty"`
	Error   *interface{} `json:"error,omitempty"`
	Message *string      `json:"message,omitempty"`
}

// RespondWithNoContent sends a response without any content
func RespondWithNoContent(w http.ResponseWriter, code int) {
	respond(w, code, BaseResponse{})
}

// RespondWithMessage sends a response with a simple text message
func RespondWithMessage(w http.ResponseWriter, code int, message string) {
	respond(w, code, BaseResponse{Message: &message})
}

// RespondWithJSON sends a response containing a JSON object
func RespondWithJSON(w http.ResponseWriter, code int, jsonPayload interface{}) {
	respond(w, code, BaseResponse{Data: &jsonPayload})
}

// RespondWithError sends a response with an error message
func RespondWithError(w http.ResponseWriter, code int, errorMessage interface{}) {
	strErrorMessage, ok := errorMessage.(string)
	if ok {
		if strings.Contains(strErrorMessage, "not found") {
			code = http.StatusNotFound
		}
	}
	respond(w, code, BaseResponse{Error: &errorMessage})
}

func respond(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET PUT POST PATCH DELETE OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.WriteHeader(code)
	w.Write(response)
}
