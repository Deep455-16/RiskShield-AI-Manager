package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/audit"
	"github.com/riskshield-ai/backend/internal/store"
	mw "github.com/riskshield-ai/backend/internal/middleware"
)

type AISystem struct {
	ID              uuid.UUID `json:"id"`
	OrgID           uuid.UUID `json:"org_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Purpose         string    `json:"purpose"`
	Owner           string    `json:"owner"`
	RiskClass       string    `json:"risk_class"`
	DataClass       string    `json:"data_class"`
	DeploymentStatus string `json:"deployment_status"`
	ApprovalStatus  string    `json:"approval_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func handleListAISystems(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		rows, err := db.Query(r.Context(), `
			SELECT id, org_id, name, description, purpose, owner, risk_class, data_class,
				deployment_status, approval_status, created_at
			FROM ai_systems WHERE org_id = $1 ORDER BY created_at DESC
		`, orgID)
		if err != nil {
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var systems []*AISystem
		for rows.Next() {
			var s AISystem
			rows.Scan(&s.ID, &s.OrgID, &s.Name, &s.Description, &s.Purpose, &s.Owner,
				&s.RiskClass, &s.DataClass, &s.DeploymentStatus, &s.ApprovalStatus, &s.CreatedAt)
			systems = append(systems, &s)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": systems})
	}
}

func handleCreateAISystem(db *store.DB, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		var req AISystem
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		var id uuid.UUID
		err := db.QueryRow(r.Context(), `
			INSERT INTO ai_systems (id, org_id, name, description, purpose, owner, risk_class, data_class,
				deployment_status, approval_status, created_at, updated_at)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
			RETURNING id
		`, orgID, req.Name, req.Description, req.Purpose, req.Owner, req.RiskClass,
			req.DataClass, req.DeploymentStatus, req.ApprovalStatus).Scan(&id)
		if err != nil {
			http.Error(w, `{"error":"failed to create"}`, http.StatusInternalServerError)
			return
		}

		user := mw.GetUser(r.Context())
		auditSvc.Log(r.Context(), &audit.AuditEvent{
			ActorID:      user.UserID,
			Action:       "ai_system_created",
			ResourceType: "ai_system",
			ResourceID:   id,
			OrgID:        orgID,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id.String()})
	}
}

func handleGetAISystem(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}
		var s AISystem
		err = db.QueryRow(r.Context(), `
			SELECT id, org_id, name, description, purpose, owner, risk_class, data_class,
				deployment_status, approval_status, created_at
			FROM ai_systems WHERE id = $1 AND org_id = $2
		`, id, orgID).Scan(&s.ID, &s.OrgID, &s.Name, &s.Description, &s.Purpose, &s.Owner,
			&s.RiskClass, &s.DataClass, &s.DeploymentStatus, &s.ApprovalStatus, &s.CreatedAt)
		if err != nil {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	}
}

// Model handlers
type Model struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Type        string    `json:"type"`
	Framework   string    `json:"framework"`
	Owner       string    `json:"owner"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func handleListModels(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		rows, err := db.Query(r.Context(), `SELECT id, org_id, name, version, type, framework, owner, status, created_at FROM models WHERE org_id = $1`, orgID)
		if err != nil {
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var models []*Model
		for rows.Next() {
			var m Model
			rows.Scan(&m.ID, &m.OrgID, &m.Name, &m.Version, &m.Type, &m.Framework, &m.Owner, &m.Status, &m.CreatedAt)
			models = append(models, &m)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": models})
	}
}

func handleCreateModel(db *store.DB, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		var req Model
		json.NewDecoder(r.Body).Decode(&req)
		var id uuid.UUID
		db.QueryRow(r.Context(), `
			INSERT INTO models (id, org_id, name, version, type, framework, owner, status, created_at)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id
		`, orgID, req.Name, req.Version, req.Type, req.Framework, req.Owner, req.Status).Scan(&id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id.String()})
	}
}
