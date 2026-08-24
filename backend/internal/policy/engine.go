package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/audit"
	"github.com/riskshield-ai/backend/internal/risk"
	"github.com/riskshield-ai/backend/internal/store"
)

type Service struct {
	db       *store.DB
	riskSvc  *risk.Service
	auditSvc *audit.Service
}

func NewService(db *store.DB, riskSvc *risk.Service, auditSvc *audit.Service) *Service {
	return &Service{db: db, riskSvc: riskSvc, auditSvc: auditSvc}
}

type Policy struct {
	ID          uuid.UUID       `json:"id"`
	OrgID       uuid.UUID       `json:"org_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Condition   json.RawMessage `json:"condition"`
	Action      string          `json:"action"` // allow, review, block
	Priority    int             `json:"priority"`
	Version     int             `json:"version"`
	ApprovedBy  *uuid.UUID      `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time      `json:"approved_at,omitempty"`
	Status      string          `json:"status"` // draft, active, archived
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type PolicyCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // eq, gt, gte, lt, lte, contains
	Value    interface{} `json:"value"`
}

type PolicyDecision struct {
	Decision    string   `json:"decision"`    // ALLOW, REVIEW, BLOCK
	PolicyID    *uuid.UUID `json:"policy_id,omitempty"`
	PolicyName  string   `json:"policy_name,omitempty"`
	Reason      string   `json:"reason"`
	RiskScore   float64  `json:"risk_score"`
	RequiredAction string `json:"required_action,omitempty"`
}

func (s *Service) Evaluate(ctx context.Context, orgID uuid.UUID, riskScore *risk.RiskScore, context map[string]interface{}) (*PolicyDecision, error) {
	// Load active policies for org, ordered by priority
	rows, err := s.db.Query(ctx, `
		SELECT id, org_id, name, condition, action, priority
		FROM policies WHERE org_id = $1 AND status = 'active' ORDER BY priority DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Policy
		rows.Scan(&p.ID, &p.OrgID, &p.Name, &p.Condition, &p.Action, &p.Priority)
		if s.matchesCondition(p.Condition, riskScore, context) {
			decision := &PolicyDecision{
				Decision:   s.actionToDecision(p.Action),
				PolicyID:   &p.ID,
				PolicyName: p.Name,
				Reason:     fmt.Sprintf("Policy '%s' triggered: %s", p.Name, p.Action),
				RiskScore:  riskScore.OverallScore,
			}
			if decision.Decision == "BLOCK" {
				decision.RequiredAction = "Action blocked by policy enforcement"
			} else if decision.Decision == "REVIEW" {
				decision.RequiredAction = "Human review required"
			}
			return decision, nil
		}
	}

	// Default decision based on risk level
	decision := &PolicyDecision{
		Decision:  "ALLOW",
		Reason:    "No policies triggered; default allow",
		RiskScore: riskScore.OverallScore,
	}
	if riskScore.RiskLevel == "CRITICAL" {
		decision.Decision = "BLOCK"
		decision.Reason = "Critical risk level - default block"
		decision.RequiredAction = "Critical risk requires explicit approval"
	} else if riskScore.RiskLevel == "HIGH" {
		decision.Decision = "REVIEW"
		decision.Reason = "High risk level - default review"
		decision.RequiredAction = "Human review required"
	}

	return decision, nil
}

func (s *Service) matchesCondition(conditionJSON json.RawMessage, riskScore *risk.RiskScore, ctx map[string]interface{}) bool {
	var conds []PolicyCondition
	if err := json.Unmarshal(conditionJSON, &conds); err != nil {
		return false
	}
	for _, c := range conds {
		var val float64
		switch c.Field {
		case "risk_score":
			val = riskScore.OverallScore
		case "fraud_probability":
			if v, ok := ctx["fraud_probability"].(float64); ok {
				val = v
			}
		case "prompt_injection_detected":
			if v, ok := ctx["prompt_injection_detected"].(bool); ok && v {
				return c.Operator == "eq" && c.Value == true
			}
		default:
			continue
		}

		target, _ := c.Value.(float64)
		switch c.Operator {
		case "gte":
			if val < target {
				return false
			}
		case "gt":
			if val <= target {
				return false
			}
		case "eq":
			if val != target {
				return false
			}
		}
	}
	return true
}

func (s *Service) actionToDecision(action string) string {
	switch action {
	case "block":
		return "BLOCK"
	case "review":
		return "REVIEW"
	default:
		return "ALLOW"
	}
}

func (s *Service) CreatePolicy(ctx context.Context, orgID uuid.UUID, name, description string, condition json.RawMessage, action string, priority int) (*Policy, error) {
	var p Policy
	err := s.db.QueryRow(ctx, `
		INSERT INTO policies (id, org_id, name, description, condition, action, priority, version, status, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, 1, 'draft', NOW(), NOW())
		RETURNING id, org_id, name, description, condition, action, priority, version, status, created_at, updated_at
	`, orgID, name, description, condition, action, priority).Scan(
		&p.ID, &p.OrgID, &p.Name, &p.Description, &p.Condition, &p.Action, &p.Priority,
		&p.Version, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (s *Service) ListPolicies(ctx context.Context, orgID uuid.UUID) ([]*Policy, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, org_id, name, description, action, priority, version, status, created_at
		FROM policies WHERE org_id = $1 ORDER BY priority DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []*Policy
	for rows.Next() {
		var p Policy
		rows.Scan(&p.ID, &p.OrgID, &p.Name, &p.Description, &p.Action, &p.Priority,
			&p.Version, &p.Status, &p.CreatedAt)
		policies = append(policies, &p)
	}
	return policies, rows.Err()
}
