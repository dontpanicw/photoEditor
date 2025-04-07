package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework-dontpanicw/app/repository"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateUserRequest(r *http.Request) (*RegisterRequest, error) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || (req.Username == "" || req.Password == "") {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

type LoginResponse struct {
	SessionId int64 `json:"session_id"`
}

type RegisterResponse struct {
	Result string `json:"result"`
}

//ERRORS

func UserProcessError(w http.ResponseWriter, err error, resp any, successStatus int) {
	if err != nil {
		var statusCode int
		var errMessage string
		switch {
		case errors.Is(err, repository.BadRequest):
			statusCode = http.StatusBadRequest
			errMessage = "Bad request"
		case errors.Is(err, repository.UserNotFound):
			statusCode = http.StatusNotFound
			errMessage = "Not found"
		default:
			statusCode = http.StatusInternalServerError
			errMessage = "Internal server error"
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(ErrorResponse{Error: errMessage})
		return
	}
	if successStatus != 0 {
		w.WriteHeader(successStatus)
	}
	//w.WriteHeader(http.StatusOK)
	if resp != nil {
		json.NewEncoder(w).Encode(resp)
	}
}
