package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"homework-dontpanicw/app/api/http/types"
	"homework-dontpanicw/app/domain"
	"homework-dontpanicw/app/repository"
	"homework-dontpanicw/app/usecases"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Task struct {
	service usecases.Task
}

// NewTaskHandler creates a new instance of Task.
func NewTaskHandler(service usecases.Task) *Task {
	return &Task{
		service: service,
	}
}

// @Summary Create task
// @Description Create a new task
// @Tags object
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body domain.Task true "Task data"
// @Success 200 {string} string types.PostTaskHandlerResponse "Task successfully created"
// @Failure 400 {object} types.Error "Invalid request data"
// @Failure 401 {object} types.Error "Unauthorized"
// @Failure 500 {object} types.Error "Internal server error"
// @Router /task [post]
// @Security BearerAuth
func (s *Task) newTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metadataStr := r.FormValue("task_metadata")
	var task domain.Task
	if err := json.Unmarshal([]byte(metadataStr), &task); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	id := uuid.New()
	status := "in progress"

	newTask := domain.Task{
		PhotoId:   id,
		Parameter: task.Parameter,
		Filter:    task.Filter,
		Status:    status,
	}

	uploadPath, err := filepath.Abs("repository/upload_photo/photo_storage")

	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		log.Printf("Ошибка создания папки: %v", err)
		http.Error(w, "Ошибка создания директории: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filename := id.String() + ".png"
	filePath := filepath.Join(uploadPath, filename)
	log.Printf("Сохранение файла по пути: %s", filePath)

	outfile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer outfile.Close()

	if err := png.Encode(outfile, img); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.service.PostTask(ctx, id, newTask)
	types.TaskProcessError(w, err, id)

	//go s.service.DoingTask(id)

}

// @Summary GetTask an object
// @Description GetTask status by id
// @Tags object
// @Accept  json
// @Produce json
// @Param id path string true "task id"
// @Success 200 {object} types.GetTaskStatusHandlerResponse
// @Failure 400 {string} types.Error
// @Failure 404 {string} types.Error
// @Failure 500 {string} types.Error
// @Router /status/{id} [get]
func (s *Task) taskStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := types.CreateGetTaskHandlerRequest(r)
	if err != nil {
		types.TaskProcessError(w, repository.TaskNotFound, nil)
		return
	}
	task, err := s.service.GetTask(ctx, req.Id)
	if err != nil {
		types.TaskProcessError(w, err, nil)
		return
	}
	types.TaskProcessError(w, nil, &types.GetTaskStatusHandlerResponse{Status: task.Status})
}

// @Summary GetTask a task result
// @Description GetTask task result by id
// @Tags object
// @Accept  json
// @Produce json
// @Param id path string true "task id"
// @Success 200 {object} types.GetTaskStatusHandlerResponse
// @Failure 400 {string} types.Error
// @Failure 404 {string} types.Error
// @Failure 500 {string} types.Error
// @Router /result/{id} [get]
func (s *Task) taskResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := types.CreateGetTaskHandlerRequest(r)
	if err != nil {
		types.TaskProcessError(w, repository.TaskNotFound, nil)
		return
	}
	task, err := s.service.GetTask(ctx, req.Id)
	if err != nil {
		types.TaskProcessError(w, err, nil)
		return
	}
	if task.Status == "ready" {
		photoName := fmt.Sprintf("repository/upload_photo/photo_storage/%s.png", task.PhotoId.String())
		http.ServeFile(w, r, photoName)
	} else {
		result := "task still in process"
		types.TaskProcessError(w, nil, &types.GetTaskResultHandlerResponse{Result: result})
	}
	//res := types.CreateResultTaskResponse(task.Status)
	//types.TaskProcessError(w, nil, res)
}

// @Summary All tasks
// @Description Supporting handler for get all tasks
// @Tags object
// @Accept  json
// @Produce json
// @Success 200 {object} ram_storage.TaskRepository
// @Failure 400 {string} types.Error
// @Failure 404 {string} types.Error
// @Router /tasks [get]
func (s *Task) allTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := s.service.GetAllTasks(ctx)
	if err != nil {
		types.TaskProcessError(w, err, nil)
	}
	json.NewEncoder(w).Encode(tasks)
}
