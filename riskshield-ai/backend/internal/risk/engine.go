package risk

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/audit"
	"github.com/riskshield-ai/backend/internal/store"
)

type Service struct {
	db       *store.DB
	auditSvc *audit.Service
}

func NewService(db *store.DB, auditSvc *audit.Service) *Service {
	return &Service{db: db, auditSvc: auditSvc}
}

// RiskDimension represents one risk scoring dimension
type RiskDimension struct {
	Name               string   `json:"name"`
	Score              float64  `json:"score"`        // 0-100
	Weight             float64  `json:"weight"`       // 0-1
	Confidence         float64  `json:"confidence"`   // 0-1
	ContributingFactors []string `json:"contributing_factors"`
	Evidence           []Evidence `json:"evidence"`
}

type Evidence struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Value       interface{} `json:"value"`
	Timestamp   time.Time   `json:"timestamp"`
}

type RiskScore struct {
	ID              uuid.UUID       `json:"id"`
	OrgID           uuid.UUID       `json:"org_id"`
	AssetID         uuid.UUID       `json:"asset_id"`
	AssetType       string          `json:"asset_type"`
	OverallScore    float64         `json:"overall_score"`
	RiskLevel       string          `json:"risk_level"`
	Confidence      float64         `json:"confidence"`
	Dimensions      json.RawMessage `json:"dimensions"`
	ContributingFactors json.RawMessage `json:"contributing_factors"`
	Evidence        json.RawMessage `json:"evidence"`
	ScoringVersion  string          `json:"scoring_version"`
	ModelVersion    string          `json:"model_version"`
	PolicyVersion   string          `json:"policy_version"`
	Timestamp       time.Time       `json:"timestamp"`
}

// Default weights for fintech payment risk
var DefaultWeights = map[string]float64{
	"fraud":       0.25,
	"security":    0.25,
	"privacy":     0.15,
	"compliance":  0.15,
	"fairness":    0.10,
	"reliability": 0.05,
	"operational": 0.05,
}

func (s *Service) CalculateRisk(ctx context.Context, orgID, assetID uuid.UUID, assetType string, dimensions []RiskDimension) (*RiskScore, error) {
	var weightedSum float64
	var totalWeight float64
	var allFactors []string
	var allEvidence []Evidence
	var minConfidence float64 = 1.0

	for _, d := range dimensions {
		weightedSum += d.Score * d.Weight
		totalWeight += d.Weight
		allFactors = append(allFactors, d.ContributingFactors...)
		allEvidence = append(allEvidence, d.Evidence...)
		if d.Confidence < minConfidence {
			minConfidence = d.Confidence
		}
	}

	if totalWeight == 0 {
		totalWeight = 1
	}

	overall := weightedSum / totalWeight
	overall = math.Min(100, math.Max(0, overall))

	level := "LOW"
	if overall >= 75 {
		level = "CRITICAL"
	} else if overall >= 50 {
		level = "HIGH"
	} else if overall >= 25 {
		level = "MEDIUM"
	}

	dimJSON, _ := json.Marshal(dimensions)
	factorJSON, _ := json.Marshal(allFactors)
	evidenceJSON, _ := json.Marshal(allEvidence)

	score := &RiskScore{
		ID:                  uuid.New(),
		OrgID:               orgID,
		AssetID:             assetID,
		AssetType:           assetType,
		OverallScore:        math.Round(overall*10) / 10,
		RiskLevel:           level,
		Confidence:          math.Round(minConfidence*100) / 100,
		Dimensions:          dimJSON,
		ContributingFactors: factorJSON,
		Evidence:            evidenceJSON,
		ScoringVersion:      "1.0.0",
		ModelVersion:        "riskshield-v1",
		PolicyVersion:       "1.0.0",
		Timestamp:           time.Now().UTC(),
	}

	_, err := s.db.Pool().Exec(ctx, `
		INSERT INTO risk_scores (id, org_id, asset_id, asset_type, overall_score, risk_level,
			confidence, dimensions, contributing_factors, evidence, scoring_version,
			model_version, policy_version, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, score.ID, score.OrgID, score.AssetID, score.AssetType, score.OverallScore,
		score.RiskLevel, score.Confidence, score.Dimensions, score.ContributingFactors,
		score.Evidence, score.ScoringVersion, score.ModelVersion, score.PolicyVersion, score.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("persist risk score: %w", err)
	}

	// Audit log
	s.auditSvc.Log(ctx, &audit.AuditEvent{
		ActorID:      orgID, // system actor
		Action:       "risk_calculated",
		ResourceType: "risk_score",
		ResourceID:   score.ID,
		OrgID:        orgID,
	})

	return score, nil
}

func (s *Service) GetLatestForAsset(ctx context.Context, orgID, assetID uuid.UUID) (*RiskScore, error) {
	var score RiskScore
	err := s.db.QueryRow(ctx, `
		SELECT id, org_id, asset_id, asset_type, overall_score, risk_level, confidence,
			dimensions, contributing_factors, evidence, scoring_version, model_version,
			policy_version, timestamp
		FROM risk_scores
		WHERE org_id = $1 AND asset_id = $2
		ORDER BY timestamp DESC LIMIT 1
	`, orgID, assetID).Scan(&score.ID, &score.OrgID, &score.AssetID, &score.AssetType,
		&score.OverallScore, &score.RiskLevel, &score.Confidence, &score.Dimensions,
		&score.ContributingFactors, &score.Evidence, &score.ScoringVersion,
		&score.ModelVersion, &score.PolicyVersion, &score.Timestamp)
	if err != nil {
		return nil, err
	}
	return &score, nil
}

func (s *Service) ListForOrg(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*RiskScore, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, org_id, asset_id, asset_type, overall_score, risk_level, confidence,
			dimensions, contributing_factors, evidence, scoring_version, model_version,
			policy_version, timestamp
		FROM risk_scores WHERE org_id = $1 ORDER BY timestamp DESC LIMIT $2 OFFSET $3
	`, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []*RiskScore
	for rows.Next() {
		var score RiskScore
		rows.Scan(&score.ID, &score.OrgID, &score.AssetID, &score.AssetType,
			&score.OverallScore, &score.RiskLevel, &score.Confidence, &score.Dimensions,
			&score.ContributingFactors, &score.Evidence, &score.ScoringVersion,
			&score.ModelVersion, &score.PolicyVersion, &score.Timestamp)
		scores = append(scores, &score)
	}
	return scores, rows.Err()
}
