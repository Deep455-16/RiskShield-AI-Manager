package api

import (
	"encoding/json"
	"net/http"

	"github.com/riskshield-ai/backend/internal/audit"
	"github.com/riskshield-ai/backend/internal/policy"
	"github.com/riskshield-ai/backend/internal/risk"
	"github.com/riskshield-ai/backend/internal/store"
)

// Helper for simple JSON responses
func returnMockData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
}

func returnEmptyData(w http.ResponseWriter) {
	returnMockData(w, []interface{}{})
}

// ----------------------------------------------------------------------------
// Agents
// ----------------------------------------------------------------------------

func handleListAgents(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnMockData(w, []map[string]interface{}{
			{"id": "1", "name": "CreditScorer-v3", "purpose": "Automated credit limit adjustments", "model": "gpt-4", "environment": "production", "risk_level": "high", "approval_status": "approved"},
			{"id": "2", "name": "SupportAgent-Gen2", "purpose": "Customer email resolution", "model": "claude-3-sonnet", "environment": "production", "risk_level": "medium", "approval_status": "approved"},
		})
	}
}

func handleCreateAgent(db *store.DB, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "new-agent"})
	}
}

func handleAgentAction(db *store.DB, riskSvc *risk.Service, policySvc *policy.Service, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "action_recorded"})
	}
}

func handleListAgentActions(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnEmptyData(w)
	}
}

// ----------------------------------------------------------------------------
// Incidents
// ----------------------------------------------------------------------------

func handleListIncidents(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnMockData(w, []map[string]interface{}{
			{"id": "1", "title": "Prompt Injection in SupportAgent", "description": "User attempted to bypass system prompt", "severity": "HIGH", "status": "OPEN"},
			{"id": "2", "title": "Model Drift: FraudDetect-v2", "description": "Accuracy dropped by 4% in last 24h", "severity": "MEDIUM", "status": "RESOLVED"},
		})
	}
}

func handleCreateIncident(db *store.DB, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "new-incident"})
	}
}

func handleUpdateIncidentStatus(db *store.DB, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}
}

// ----------------------------------------------------------------------------
// Compliance
// ----------------------------------------------------------------------------

func handleListFrameworks(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnEmptyData(w)
	}
}

func handleListControls(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnMockData(w, []map[string]interface{}{
			{"id": "1", "control_id": "SEC-1", "title": "Model Security", "description": "Protect AI models from attacks", "status": "IMPLEMENTED"},
			{"id": "2", "control_id": "GOV-2", "title": "Risk Assessment", "description": "Conduct regular risk assessments", "status": "PARTIALLY_IMPLEMENTED"},
			{"id": "3", "control_id": "FAIR-1", "title": "Bias Detection", "description": "Detect and mitigate algorithmic bias", "status": "NOT_IMPLEMENTED"},
		})
	}
}

func handleUpdateControlStatus(db *store.DB, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}
}

// ----------------------------------------------------------------------------
// Audit Logs
// ----------------------------------------------------------------------------

func handleListAuditLogs(auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnMockData(w, []map[string]interface{}{
			{"id": "1", "action": "policy_updated", "resource_type": "policy", "timestamp": "2026-08-25T00:00:00Z"},
			{"id": "2", "action": "agent_action_blocked", "resource_type": "agent_action", "timestamp": "2026-08-25T00:05:00Z"},
		})
	}
}

// ----------------------------------------------------------------------------
// Copilot
// ----------------------------------------------------------------------------

func handleCopilotExplain(riskSvc *risk.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"explanation": "This was flagged by AI."})
	}
}

// ----------------------------------------------------------------------------
// Reports
// ----------------------------------------------------------------------------

func handleReportSummary(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"summary": "All systems nominal."})
	}
}

// ----------------------------------------------------------------------------
// Approvals
// ----------------------------------------------------------------------------

func handleListApprovals(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnEmptyData(w)
	}
}

func handleDecideApproval(db *store.DB, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "decided"})
	}
}

// ----------------------------------------------------------------------------
// Dashboard
// ----------------------------------------------------------------------------

func handleDashboard(db *store.DB, riskSvc *risk.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"active_incidents": 2,
			"critical_risks": 5,
			"policy_violations_24h": 18,
			"ai_systems_monitored": 12,
			"compliance_score": 92,
			"transactions_today": 3450,
			"agents_active": 8,
			"audit_events_24h": 890,
		})
	}
}
