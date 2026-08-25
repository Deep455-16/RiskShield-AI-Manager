package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/middleware"
	"github.com/riskshield-ai/backend/internal/store"
)

// getOrgID extracts the org UUID from the request context safely.
func getOrgID(r *http.Request) uuid.UUID {
	return middleware.OrgIDFromContext(r.Context())
}

// getUserID extracts the user UUID from the request context. Falls back to a zero UUID.
func getUserID(r *http.Request) uuid.UUID {
	user := middleware.GetUser(r.Context())
	if user == nil {
		return uuid.Nil
	}
	return user.UserID
}

// ─── Agent Extended Handlers ──────────────────────────────────────────────────

func handleGetAgentBehaviorLogs(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)
		agentID := chi.URLParam(r, "id")

		rows, err := db.Query(r.Context(), `
			SELECT id, agent_id, outcome, payload, ts
			FROM agent_behavior_logs
			WHERE org_id = $1 AND agent_id = $2
			ORDER BY ts DESC LIMIT 100
		`, orgID, agentID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		type Log struct {
			ID      string    `json:"id"`
			AgentID string    `json:"agent_id"`
			Outcome string    `json:"outcome"`
			Payload any       `json:"payload"`
			Ts      time.Time `json:"ts"`
		}

		var logs []Log
		for rows.Next() {
			var l Log
			var payload []byte
			rows.Scan(&l.ID, &l.AgentID, &l.Outcome, &payload, &l.Ts)
			json.Unmarshal(payload, &l.Payload)
			logs = append(logs, l)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": logs})
	}
}

func handleAgentKillSwitch(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)
		agentID := chi.URLParam(r, "id")

		_, err := db.Query(r.Context(), `
			UPDATE agents SET status = 'suspended', updated_at = NOW()
			WHERE id = $1 AND org_id = $2
		`, agentID, orgID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		db.Query(r.Context(), `
			INSERT INTO audit_logs (id, org_id, actor_id, action, resource_type, resource_id, created_at)
			VALUES (gen_random_uuid(), $1, $2, 'KILL_SWITCH', 'agent', $3, NOW())
		`, orgID, getUserID(r), agentID)

		json.NewEncoder(w).Encode(map[string]any{"status": "suspended", "agent_id": agentID})
	}
}

// ─── Shadow AI Handlers ───────────────────────────────────────────────────────

func handleListShadowAI(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)

		rows, err := db.Query(r.Context(), `
			SELECT id, source, external_id, name, details, status, discovered_at
			FROM shadow_ai_inbox WHERE org_id = $1
			ORDER BY discovered_at DESC
		`, orgID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		type Item struct {
			ID           string    `json:"id"`
			Source       string    `json:"source"`
			ExternalID   string    `json:"external_id"`
			Name         string    `json:"name"`
			Details      any       `json:"details"`
			Status       string    `json:"status"`
			DiscoveredAt time.Time `json:"discovered_at"`
		}
		var items []Item
		for rows.Next() {
			var it Item
			var details []byte
			rows.Scan(&it.ID, &it.Source, &it.ExternalID, &it.Name, &details, &it.Status, &it.DiscoveredAt)
			json.Unmarshal(details, &it.Details)
			items = append(items, it)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": items})
	}
}

func handleShadowAIDiscover(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)

		discoveries := []struct {
			Source     string
			ExternalID string
			Name       string
			Details    map[string]any
		}{
			{"github", "repo/ml-fraud-detector", "ML Fraud Detector (GitHub)", map[string]any{"lang": "Python", "keywords": []string{"ml", "inference", "fraud"}}},
			{"huggingface", "org/bert-risk-model", "BERT Risk Classifier (HuggingFace)", map[string]any{"downloads": 4200, "task": "text-classification"}},
			{"aws_bedrock", "us-east-1/claude-v3", "AWS Claude v3 (Bedrock)", map[string]any{"region": "us-east-1", "model": "anthropic.claude-3"}},
			{"azure_ai", "azure/gpt4o-banking", "GPT-4o Deployment (Azure AI)", map[string]any{"endpoint": "azure.openai.com", "model": "gpt-4o"}},
		}

		inserted := 0
		for _, d := range discoveries {
			details, _ := json.Marshal(d.Details)
			_, err := db.Query(r.Context(), `
				INSERT INTO shadow_ai_inbox (id, org_id, source, external_id, name, details)
				VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
				ON CONFLICT (content_hash) DO NOTHING
			`, orgID, d.Source, d.ExternalID, d.Name, details)
			if err == nil {
				inserted++
			}
		}

		json.NewEncoder(w).Encode(map[string]any{"inserted": inserted, "total_scanned": len(discoveries)})
	}
}

// ─── Vendor Handlers ──────────────────────────────────────────────────────────

func handleListVendors(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)

		rows, err := db.Query(r.Context(), `
			SELECT id, name, service, risk_tier, renewal_date, created_at
			FROM vendors WHERE org_id = $1
			ORDER BY created_at DESC
		`, orgID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		type Vendor struct {
			ID          string     `json:"id"`
			Name        string     `json:"name"`
			Type        string     `json:"type"`
			RiskTier    string     `json:"risk_tier"`
			RenewalDate *time.Time `json:"renewal_date"`
			CreatedAt   time.Time  `json:"created_at"`
		}
		var vendors []Vendor
		for rows.Next() {
			var v Vendor
			rows.Scan(&v.ID, &v.Name, &v.Type, &v.RiskTier, &v.RenewalDate, &v.CreatedAt)
			vendors = append(vendors, v)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": vendors})
	}
}

// ─── Tasks Handlers ───────────────────────────────────────────────────────────

func handleListTasks(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)

		rows, err := db.Query(r.Context(), `
			SELECT id, title, priority, status, due_date, assignee_id, created_at
			FROM tasks WHERE org_id = $1
			ORDER BY
				CASE priority WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 ELSE 4 END,
				due_date ASC NULLS LAST
		`, orgID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		type Task struct {
			ID         string     `json:"id"`
			Title      string     `json:"title"`
			Priority   string     `json:"priority"`
			Status     string     `json:"status"`
			DueDate    *time.Time `json:"due_date"`
			AssigneeID *string    `json:"assignee_id"`
			CreatedAt  time.Time  `json:"created_at"`
			Overdue    bool       `json:"overdue"`
		}
		var tasks []Task
		for rows.Next() {
			var t Task
			rows.Scan(&t.ID, &t.Title, &t.Priority, &t.Status, &t.DueDate, &t.AssigneeID, &t.CreatedAt)
			if t.DueDate != nil && t.DueDate.Before(time.Now()) && t.Status != "done" && t.Status != "cancelled" {
				t.Overdue = true
			}
			tasks = append(tasks, t)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": tasks})
	}
}

func handleCreateTask(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)

		var req struct {
			Title    string     `json:"title"`
			Priority string     `json:"priority"`
			DueDate  *time.Time `json:"due_date"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Priority == "" {
			req.Priority = "medium"
		}

		var id string
		db.QueryRow(r.Context(), `
			INSERT INTO tasks (id, org_id, title, priority, due_date, status)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, 'open')
			RETURNING id
		`, orgID, req.Title, req.Priority, req.DueDate).Scan(&id)

		json.NewEncoder(w).Encode(map[string]any{"id": id, "status": "created"})
	}
}

func handleUpdateTaskStatus(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)
		taskID := chi.URLParam(r, "id")

		var req struct {
			Status string `json:"status"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		db.Query(r.Context(), `
			UPDATE tasks SET status = $1, updated_at = NOW()
			WHERE id = $2 AND org_id = $3
		`, req.Status, taskID, orgID)

		json.NewEncoder(w).Encode(map[string]any{"status": "updated"})
	}
}

// ─── Approval Workflow Handlers ───────────────────────────────────────────────

func handleListApprovalWorkflows(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)

		rows, err := db.Query(r.Context(), `
			SELECT id, resource_type, resource_id, min_approvals, status, created_at
			FROM approval_workflows WHERE org_id = $1
			ORDER BY created_at DESC
		`, orgID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		type Workflow struct {
			ID           string    `json:"id"`
			ResourceType string    `json:"resource_type"`
			ResourceID   string    `json:"resource_id"`
			MinApprovals int       `json:"min_approvals"`
			Status       string    `json:"status"`
			CreatedAt    time.Time `json:"created_at"`
		}
		var workflows []Workflow
		for rows.Next() {
			var wf Workflow
			rows.Scan(&wf.ID, &wf.ResourceType, &wf.ResourceID, &wf.MinApprovals, &wf.Status, &wf.CreatedAt)
			workflows = append(workflows, wf)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": workflows})
	}
}

func handleApprovalVote(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)
		userID := getUserID(r)
		workflowID := chi.URLParam(r, "id")

		var req struct {
			Decision string `json:"decision"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		db.Query(r.Context(), `
			INSERT INTO approval_votes (id, workflow_id, user_id, decision)
			VALUES (gen_random_uuid(), $1, $2, $3)
			ON CONFLICT (workflow_id, user_id) DO NOTHING
		`, workflowID, userID, req.Decision)

		if req.Decision == "reject" {
			db.Query(r.Context(), `UPDATE approval_workflows SET status = 'rejected' WHERE id = $1`, workflowID)
		} else {
			db.Query(r.Context(), `
				UPDATE approval_workflows SET status = 'approved'
				WHERE id = $1 AND (
					SELECT COUNT(*) FROM approval_votes WHERE workflow_id = $1 AND decision = 'approve'
				) >= min_approvals
			`, workflowID)
		}

		db.Query(r.Context(), `
			INSERT INTO audit_logs (id, org_id, actor_id, action, resource_type, resource_id, created_at)
			VALUES (gen_random_uuid(), $1, $2, $3, 'approval_workflow', $4, NOW())
		`, orgID, userID, "VOTE_"+req.Decision, workflowID)

		json.NewEncoder(w).Encode(map[string]any{"status": "vote_recorded"})
	}
}

// ─── Reports Handler ──────────────────────────────────────────────────────────

func handleGenerateReport(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)
		userID := getUserID(r)

		var req struct {
			Scope  string `json:"scope"`
			Format string `json:"format"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Scope == "" {
			req.Scope = "org"
		}
		if req.Format == "" {
			req.Format = "json"
		}

		content := map[string]any{
			"org_id":       orgID,
			"scope":        req.Scope,
			"generated_at": time.Now(),
			"format":       req.Format,
		}
		contentBytes, _ := json.Marshal(content)
		hash := fmt.Sprintf("%x", sha256.Sum256(contentBytes))

		var reportID string
		db.QueryRow(r.Context(), `
			INSERT INTO generated_reports (id, org_id, scope, format, content_hash, generated_by)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
			RETURNING id
		`, orgID, req.Scope, req.Format, hash, userID).Scan(&reportID)

		json.NewEncoder(w).Encode(map[string]any{
			"id":           reportID,
			"scope":        req.Scope,
			"format":       req.Format,
			"sha256":       hash,
			"download_url": fmt.Sprintf("/api/v1/reports/%s/download", reportID),
			"content":      content,
		})
	}
}

// ─── Dashboard Snapshot Handler ───────────────────────────────────────────────

func handleDashboardSnapshot(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)

		type snapshot struct {
			IncidentsOpen    int              `json:"incidents_open"`
			TasksOpen        int              `json:"tasks_open"`
			TasksOverdue     int              `json:"tasks_overdue"`
			AgentsActive     int              `json:"agents_active"`
			ShadowAIOpen     int              `json:"shadow_ai_open"`
			VendorsActive    int              `json:"vendors_active"`
			PoliciesApproved int              `json:"policies_approved"`
			ComplianceScore  float64          `json:"compliance_score"`
			IncidentTrend    []map[string]any `json:"incident_trend"`
		}

		s := snapshot{}

		db.QueryRow(r.Context(), `SELECT COUNT(*) FROM incidents WHERE org_id = $1 AND status = 'OPEN'`, orgID).Scan(&s.IncidentsOpen)
		db.QueryRow(r.Context(), `SELECT COUNT(*) FROM tasks WHERE org_id = $1 AND status NOT IN ('done','cancelled')`, orgID).Scan(&s.TasksOpen)
		db.QueryRow(r.Context(), `SELECT COUNT(*) FROM tasks WHERE org_id = $1 AND status NOT IN ('done','cancelled') AND due_date < NOW()`, orgID).Scan(&s.TasksOverdue)
		db.QueryRow(r.Context(), `SELECT COUNT(*) FROM agents WHERE org_id = $1 AND (status IS NULL OR status = 'active')`, orgID).Scan(&s.AgentsActive)
		db.QueryRow(r.Context(), `SELECT COUNT(*) FROM shadow_ai_inbox WHERE org_id = $1 AND status = 'unreviewed'`, orgID).Scan(&s.ShadowAIOpen)
		db.QueryRow(r.Context(), `SELECT COUNT(*) FROM vendors WHERE org_id = $1`, orgID).Scan(&s.VendorsActive)
		db.QueryRow(r.Context(), `SELECT COUNT(*) FROM policies WHERE org_id = $1 AND status = 'approved'`, orgID).Scan(&s.PoliciesApproved)

		var total, implemented, partial float64
		db.QueryRow(r.Context(), `
			SELECT COUNT(*),
				SUM(CASE WHEN status='implemented' THEN 1 ELSE 0 END),
				SUM(CASE WHEN status='partial' THEN 0.5 ELSE 0 END)
			FROM compliance_controls WHERE org_id = $1
		`, orgID).Scan(&total, &implemented, &partial)
		if total > 0 {
			s.ComplianceScore = (implemented + partial) / total * 100
		}

		trendRows, _ := db.Query(r.Context(), `
			SELECT TO_CHAR(created_at, 'Mon YYYY'), COUNT(*)
			FROM incidents WHERE org_id = $1 AND created_at > NOW() - INTERVAL '6 months'
			GROUP BY 1 ORDER BY MIN(created_at) ASC
		`, orgID)
		if trendRows != nil {
			defer trendRows.Close()
			for trendRows.Next() {
				var month string
				var count int
				trendRows.Scan(&month, &count)
				s.IncidentTrend = append(s.IncidentTrend, map[string]any{"month": month, "count": count})
			}
		}

		json.NewEncoder(w).Encode(map[string]any{"data": s})
	}
}

// ─── Audit Log Export Handler ─────────────────────────────────────────────────

func handleAuditExport(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := getOrgID(r)

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="audit_log.csv"`)
		fmt.Fprintln(w, "id,actor_id,action,resource_type,resource_id,timestamp")

		rows, err := db.Query(r.Context(), `
			SELECT id, actor_id, action, resource_type, resource_id, created_at
			FROM audit_logs WHERE org_id = $1
			ORDER BY created_at DESC LIMIT 10000
		`, orgID)
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			var id, actorID, action, resType, resID, createdAt string
			rows.Scan(&id, &actorID, &action, &resType, &resID, &createdAt)
			fmt.Fprintf(w, "%s,%s,%s,%s,%s,%s\n", id, actorID, action, resType, resID, createdAt)
		}
	}
}
