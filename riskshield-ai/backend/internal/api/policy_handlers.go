package api

import (
	"encoding/json"
	"net/http"

	"github.com/riskshield-ai/backend/internal/policy"
	mw "github.com/riskshield-ai/backend/internal/middleware"
)

func handleListPolicies(policySvc *policy.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		policies, err := policySvc.ListPolicies(r.Context(), orgID)
		if err != nil {
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": policies})
	}
}

func handleCreatePolicy(policySvc *policy.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		var req struct {
			Name        string          `json:"name"`
			Description string          `json:"description"`
			Condition   json.RawMessage `json:"condition"`
			Action      string          `json:"action"`
			Priority    int             `json:"priority"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		p, err := policySvc.CreatePolicy(r.Context(), orgID, req.Name, req.Description, req.Condition, req.Action, req.Priority)
		if err != nil {
			http.Error(w, `{"error":"failed to create policy"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}
