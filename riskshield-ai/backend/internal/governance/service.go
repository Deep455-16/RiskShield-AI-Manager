package governance

import (
	"context"
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

func (s *Service) ListAISystems(ctx context.Context, orgID uuid.UUID) ([]AISystem, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, org_id, name, version, description, owner_id, deployment_env, purpose, risk_tier, created_at, updated_at
		FROM ai_systems
		WHERE org_id = $1
		ORDER BY created_at DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var systems []AISystem
	for rows.Next() {
		var sys AISystem
		var desc, purpose string
		var ownerID uuid.UUID
		err := rows.Scan(&sys.ID, &sys.OrgID, &sys.Name, &sys.Version, &desc, &ownerID, &sys.DeploymentEnv, &purpose, &sys.RiskTier, &sys.CreatedAt, &sys.UpdatedAt)
		if err != nil {
			return nil, err
		}
		sys.Description = desc
		sys.Purpose = purpose
		sys.OwnerID = ownerID
		systems = append(systems, sys)
	}
	return systems, nil
}

func (s *Service) ListComplianceControls(ctx context.Context, orgID uuid.UUID) ([]ComplianceControl, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, org_id, framework, control_ref, description, status, ai_system_id, evidence_url, updated_at
		FROM compliance_controls
		WHERE org_id = $1
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var controls []ComplianceControl
	for rows.Next() {
		var c ComplianceControl
		var desc, evidence string
		var sysID uuid.UUID
		err := rows.Scan(&c.ID, &c.OrgID, &c.Framework, &c.ControlRef, &desc, &c.Status, &sysID, &evidence, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		c.Description = desc
		c.EvidenceURL = evidence
		c.AISystemID = sysID
		controls = append(controls, c)
	}
	return controls, nil
}

func (s *Service) CreateAISystem(ctx context.Context, orgID, actorID uuid.UUID, sys AISystem) (AISystem, error) {
	sys.ID = uuid.New()
	sys.OrgID = orgID
	sys.CreatedAt = time.Now()
	sys.UpdatedAt = time.Now()

	_, err := s.db.Query(ctx, `
		INSERT INTO ai_systems (id, org_id, name, version, description, owner_id, deployment_env, purpose, risk_tier, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, sys.ID, sys.OrgID, sys.Name, sys.Version, sys.Description, sys.OwnerID, sys.DeploymentEnv, sys.Purpose, sys.RiskTier, sys.CreatedAt, sys.UpdatedAt)
	
	if err != nil {
		return sys, err
	}

	// Audit Logging
	s.auditSvc.Log(ctx, &audit.AuditEvent{
		OrgID:        orgID,
		ActorID:      actorID,
		Action:       "CREATE",
		ResourceType: "ai_system",
		ResourceID:   sys.ID,
		Timestamp:    time.Now(),
	})

	return sys, nil
}
