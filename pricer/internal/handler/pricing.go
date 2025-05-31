package handler

import (
	"encoding/json"
	"net/http"

	"github.com/chiyonn/vendiq2/pricer/internal/service"
)

type PricingHandler struct {
	srv *service.PricingService
}

func NewPricingHandler(srv *service.PricingService) *PricingHandler {
	return &PricingHandler{srv: srv}
}

func (h *PricingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	pricings, err := h.srv.GetAll()
	if err != nil {
		http.Error(w, "failed to fetch pricing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pricings); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PricingHandler) SyncAll(w http.ResponseWriter, r *http.Request) {

	err := h.srv.SyncAll(); if err != nil {
		http.Error(w, "failed to sync pricings", http.StatusInternalServerError)
	}

	pricings, err := h.srv.GetAll()
	if err != nil {
		http.Error(w, "failed to fetch pricing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pricings); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
