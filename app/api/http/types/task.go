package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"homework-dontpanicw/app/domain"
	"homework-dontpanicw/app/repository"
	"net/http"
)

// POST
type PostTaskHandlerRequest struct {
	domain.Task
}

//type NewTaskHandlerRequest struct {
//	Filter string `json:"filter"`
//	Status string `json:"status"`
//}

func CreatePostTaskHandlerRequest(r *http.Request) (*PostTaskHandlerRequest, error) {
	var req PostTaskHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

type PostTaskHandlerResponse struct {
	Id uuid.UUID `json:"id"`
}

//GET

type GetTaskHandlerRequest struct {
	Id uuid.UUID `json:"id"`
}

func CreateGetTaskHandlerRequest(r *http.Request) (*GetTaskHandlerRequest, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return nil, fmt.Errorf("missing id parameter")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("id is not string")
	}
	return &GetTaskHandlerRequest{Id: id}, nil
}

type GetTaskStatusHandlerResponse struct {
	Status string `json:"status"`
}

type GetTaskHandlerResponse struct {
	Parameter float64 `json:"parameter"`
	Filter    string  `json:"filter"`
	Status    string  `json:"status"`
}

type GetTaskResultHandlerResponse struct {
	Result string `json:"result"`
}

func CreateResultTaskResponse(status string) *GetTaskResultHandlerResponse {
	var result string
	if status == "ready" {
		//http.ServeFile(w, r, "image.jpg")
		result = "photo was redacted successfully"
	} else {
		result = "task still in process"
	}
	return &GetTaskResultHandlerResponse{Result: result}
}

//ERRORS

type ErrorResponse struct {
	Error string `json:"error"`
}

func TaskProcessError(w http.ResponseWriter, err error, resp any) {
	if err != nil {
		var statusCode int
		var errMessage string
		switch {
		case errors.Is(err, repository.BadRequest):
			statusCode = http.StatusBadRequest
			errMessage = "Bad request"
		case errors.Is(err, repository.TaskNotFound):
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
	w.WriteHeader(http.StatusOK)
	if resp != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}
}
