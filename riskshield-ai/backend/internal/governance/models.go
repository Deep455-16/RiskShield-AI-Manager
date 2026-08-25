package governance

import (
	"time"

	"github.com/google/uuid"
)

type AISystem struct {
	ID             uuid.UUID `json:"id"`
	OrgID          uuid.UUID `json:"org_id"`
	Name           string    `json:"name"`
	Version        string    `json:"version"`
	Description    string    `json:"description"`
	OwnerID        uuid.UUID `json:"owner_id"`
	DeploymentEnv  string    `json:"deployment_env"`
	Purpose        string    `json:"purpose"`
	RiskTier       string    `json:"risk_tier"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ComplianceControl struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Framework   string    `json:"framework"`
	ControlRef  string    `json:"control_ref"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AISystemID  uuid.UUID `json:"ai_system_id"`
	EvidenceURL string    `json:"evidence_url"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Assessment struct {
	ID               uuid.UUID `json:"id"`
	OrgID            uuid.UUID `json:"org_id"`
	AISystemID       uuid.UUID `json:"ai_system_id"`
	Framework        string    `json:"framework"`
	FrameworkVersion string    `json:"framework_version"`
	Status           string    `json:"status"`
	AssessorID       uuid.UUID `json:"assessor_id"`
	CreatedAt        time.Time `json:"created_at"`
}
