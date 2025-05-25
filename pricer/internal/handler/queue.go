package handler

import (
	"encoding/json"
	"net/http"

	"github.com/chiyonn/vendiq2/pricer/internal/service"
)

type QueueHandler struct {
	srv *service.QueueService
}

func NewQueueHandler(srv *service.QueueService) *QueueHandler {
	return &QueueHandler{srv: srv}
}

func (h *QueueHandler) PostQueue(w http.ResponseWriter, r *http.Request) {
	h.srv.AddQueue()

	w.Header().Set("Content-Type", "application/json")
}

func (h *QueueHandler) GetQueues(w http.ResponseWriter, r *http.Request) {
	queues, err := h.srv.GetAllQueues()
	if err != nil {
		http.Error(w, "failed to fetch queues data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(queues); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
