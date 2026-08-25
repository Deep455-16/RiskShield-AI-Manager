package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/governance"
)

type GovernanceHandlers struct {
	govSvc *governance.Service
}

func NewGovernanceHandlers(govSvc *governance.Service) *GovernanceHandlers {
	return &GovernanceHandlers{govSvc: govSvc}
}

func (h *GovernanceHandlers) ListAISystems(w http.ResponseWriter, r *http.Request) {
	orgID := r.Context().Value("orgID").(uuid.UUID)
	
	systems, err := h.govSvc.ListAISystems(r.Context(), orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{"data": systems})
}

func (h *GovernanceHandlers) CreateAISystem(w http.ResponseWriter, r *http.Request) {
	orgID := r.Context().Value("orgID").(uuid.UUID)
	actorID := r.Context().Value("userID").(uuid.UUID)
	
	var sys governance.AISystem
	if err := json.NewDecoder(r.Body).Decode(&sys); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	created, err := h.govSvc.CreateAISystem(r.Context(), orgID, actorID, sys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{"data": created})
}

func (h *GovernanceHandlers) ListComplianceControls(w http.ResponseWriter, r *http.Request) {
	orgID := r.Context().Value("orgID").(uuid.UUID)
	
	controls, err := h.govSvc.ListComplianceControls(r.Context(), orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{"data": controls})
}

func (h *GovernanceHandlers) RegisterRoutes(r chi.Router) {
	r.Get("/systems", h.ListAISystems)
	r.Post("/systems", h.CreateAISystem)
	r.Get("/compliance/controls", h.ListComplianceControls)
}
