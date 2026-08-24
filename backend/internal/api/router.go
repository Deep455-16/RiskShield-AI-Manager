package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/riskshield-ai/backend/internal/audit"
	"github.com/riskshield-ai/backend/internal/auth"
	"github.com/riskshield-ai/backend/internal/jobs"
	"github.com/riskshield-ai/backend/internal/policy"
	"github.com/riskshield-ai/backend/internal/risk"
	"github.com/riskshield-ai/backend/internal/store"
	mw "github.com/riskshield-ai/backend/internal/middleware"
	"github.com/riskshield-ai/backend/internal/datasets"
)

func RegisterRoutes(r chi.Router, db *store.DB, authSvc *auth.Service, auditSvc *audit.Service,
	riskSvc *risk.Service, policySvc *policy.Service, jobSvc *jobs.Service, jwtSecret string,
	registry *datasets.Registry, replayEngine *datasets.ReplayEngine) {

	// Public routes
	r.Post("/api/v1/auth/register", handleRegister(authSvc))
	r.Post("/api/v1/auth/login", handleLogin(authSvc))
	r.Post("/api/v1/auth/refresh", handleRefresh(authSvc))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware(authSvc))
		r.Use(mw.TenantIsolation)

		r.Post("/api/v1/auth/logout", handleLogout)

		// Datasets
		r.Get("/api/v1/datasets", handleListDatasets(registry))
		r.Post("/api/v1/datasets/{id}/replay", handleDatasetReplay(replayEngine))

		// Ethical AI
		r.Get("/api/v1/ethical-ai/evaluate", handleEthicalAIEvaluate())

		// AI Systems
		r.Get("/api/v1/ai-systems", handleListAISystems(db))
		r.Post("/api/v1/ai-systems", handleCreateAISystem(db, auditSvc))
		r.Get("/api/v1/ai-systems/{id}", handleGetAISystem(db))

		// Models
		r.Get("/api/v1/models", handleListModels(db))
		r.Post("/api/v1/models", handleCreateModel(db, auditSvc))

		// Agents
		r.Get("/api/v1/agents", handleListAgents(db))
		r.Post("/api/v1/agents", handleCreateAgent(db, auditSvc))
		r.Post("/api/v1/agents/{id}/actions", handleAgentAction(db, riskSvc, policySvc, auditSvc))
		r.Get("/api/v1/agents/{id}/actions", handleListAgentActions(db))

		// Transactions
		r.Get("/api/v1/transactions", handleListTransactions(db))
		r.Post("/api/v1/transactions", handleCreateTransaction(db))
		r.Post("/api/v1/transactions/{id}/assess", handleAssessTransaction(db, riskSvc, policySvc, auditSvc))

		// Risk
		r.Get("/api/v1/risk/scores", handleListRiskScores(riskSvc))
		r.Get("/api/v1/risk/scores/{id}", handleGetRiskScore(db))
		r.Post("/api/v1/risk/assess", handleRiskAssess(riskSvc))

		// Policies
		r.Get("/api/v1/policies", handleListPolicies(policySvc))
		r.Post("/api/v1/policies", handleCreatePolicy(policySvc))

		// Incidents
		r.Get("/api/v1/incidents", handleListIncidents(db))
		r.Post("/api/v1/incidents", handleCreateIncident(db, auditSvc))
		r.Patch("/api/v1/incidents/{id}/status", handleUpdateIncidentStatus(db, auditSvc))

		// Compliance
		r.Get("/api/v1/compliance/frameworks", handleListFrameworks(db))
		r.Get("/api/v1/compliance/controls", handleListControls(db))
		r.Patch("/api/v1/compliance/controls/{id}", handleUpdateControlStatus(db, auditSvc))

		// Audit
		r.Get("/api/v1/audit-logs", handleListAuditLogs(auditSvc))

		// Copilot
		r.Post("/api/v1/copilot/explain", handleCopilotExplain(riskSvc))

		// Reports
		r.Get("/api/v1/reports/summary", handleReportSummary(db))

		// Approval requests
		r.Get("/api/v1/approvals", handleListApprovals(db))
		r.Post("/api/v1/approvals/{id}/decide", handleDecideApproval(db, auditSvc))

		// Dashboard
		r.Get("/api/v1/dashboard", handleDashboard(db, riskSvc))

		// Simulator
		r.Post("/api/v1/simulate/payment", handleSimulatePayment(riskSvc, policySvc))
		r.Post("/api/v1/simulate/attack", handleSimulateAttack(riskSvc, policySvc, auditSvc))
	})
}
