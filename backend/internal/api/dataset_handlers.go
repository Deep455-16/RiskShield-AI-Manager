package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/riskshield-ai/backend/internal/datasets"
)

func handleListDatasets(registry *datasets.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		meta := registry.ListDatasets()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": meta})
	}
}

func handleDatasetReplay(engine *datasets.ReplayEngine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		datasetID := chi.URLParam(r, "id")
		
		var req struct {
			Action string  `json:"action"` // "start", "stop"
			Speed  float64 `json:"speed"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Action == "start" {
			speed := req.Speed
			if speed <= 0 {
				speed = 1.0
			}
			engine.StartReplay(datasetID, speed)
		} else if req.Action == "stop" {
			engine.StopReplay(datasetID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "action": req.Action})
	}
}

func handleEthicalAIEvaluate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mock fairness metrics for the Bank Marketing dataset
		metrics := map[string]interface{}{
			"selection_rate": 0.15,
			"true_positive_rate": 0.82,
			"false_positive_rate": 0.12,
			"demographic_parity_difference": 0.04,
			"equal_opportunity_difference": 0.02,
			"status": "compliant",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": metrics})
	}
}
