package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/risk"
	"github.com/riskshield-ai/backend/internal/store"
	mw "github.com/riskshield-ai/backend/internal/middleware"
)

func handleListRiskScores(riskSvc *risk.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		scores, err := riskSvc.ListForOrg(r.Context(), orgID, 100, 0)
		if err != nil {
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": scores})
	}
}

func handleGetRiskScore(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}
		var score risk.RiskScore
		err = db.QueryRow(r.Context(), `
			SELECT id, org_id, asset_id, asset_type, overall_score, risk_level, confidence,
				dimensions, contributing_factors, evidence, scoring_version, model_version,
				policy_version, timestamp
			FROM risk_scores WHERE id = $1 AND org_id = $2
		`, id, orgID).Scan(&score.ID, &score.OrgID, &score.AssetID, &score.AssetType,
			&score.OverallScore, &score.RiskLevel, &score.Confidence, &score.Dimensions,
			&score.ContributingFactors, &score.Evidence, &score.ScoringVersion,
			&score.ModelVersion, &score.PolicyVersion, &score.Timestamp)
		if err != nil {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(score)
	}
}

func handleRiskAssess(riskSvc *risk.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		var req struct {
			AssetID     uuid.UUID           `json:"asset_id"`
			AssetType   string              `json:"asset_type"`
			Dimensions  []risk.RiskDimension `json:"dimensions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		score, err := riskSvc.CalculateRisk(r.Context(), orgID, req.AssetID, req.AssetType, req.Dimensions)
		if err != nil {
			http.Error(w, `{"error":"assessment failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(score)
	}
}
