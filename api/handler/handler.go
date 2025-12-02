package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kayaramazan/insider-message/api/job"
	"github.com/kayaramazan/insider-message/api/service"
)

type Handler struct {
	messageService *service.MessageService
	job            *job.Job
}

func NewHandler(messageService *service.MessageService, job *job.Job) *Handler {
	return &Handler{
		messageService: messageService,
		job:            job,
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
