package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kayaramazan/insider-message/api/model"
)

func (h *handlerImpl) GetAllSentMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.messageService.GetAllSentMessages(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch messages")
		return
	}

	respondJSON(w, http.StatusOK, messages)
}

func (h *handlerImpl) CreateMessage(w http.ResponseWriter, r *http.Request) {
	message := model.Message{}
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := message.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.messageService.CreateMessage(r.Context(), &message)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create message")
		return
	}

	respondJSON(w, http.StatusOK, message)
}

func (h *handlerImpl) StartOrStop(w http.ResponseWriter, r *http.Request) {
	h.job.Toggle()

	respondJSON(w, http.StatusOK, map[string]any{
		"Status":            "Accepted",
		"Automation Status": h.job.IsRunning(),
	})
}
