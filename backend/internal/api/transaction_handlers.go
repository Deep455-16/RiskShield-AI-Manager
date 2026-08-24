package api

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/audit"
	"github.com/riskshield-ai/backend/internal/policy"
	"github.com/riskshield-ai/backend/internal/risk"
	"github.com/riskshield-ai/backend/internal/store"
	mw "github.com/riskshield-ai/backend/internal/middleware"
)

type Transaction struct {
	ID               uuid.UUID `json:"id"`
	OrgID            uuid.UUID `json:"org_id"`
	TransactionID    string    `json:"transaction_id"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	MerchantID       string    `json:"merchant_id"`
	CustomerID       string    `json:"customer_id"`
	DeviceID         string    `json:"device_id"`
	IPAddress        string    `json:"ip_address"`
	Location         string    `json:"location"`
	Velocity         int       `json:"velocity"`
	ModelProbability float64   `json:"model_probability"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

func handleListTransactions(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		rows, err := db.Query(r.Context(), `
			SELECT id, transaction_id, amount, currency, merchant_id, customer_id,
				device_id, ip_address, location, velocity, model_probability, status, created_at
			FROM transactions WHERE org_id = $1 ORDER BY created_at DESC LIMIT 100
		`, orgID)
		if err != nil {
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var txs []*Transaction
		for rows.Next() {
			var t Transaction
			rows.Scan(&t.ID, &t.TransactionID, &t.Amount, &t.Currency, &t.MerchantID, &t.CustomerID,
				&t.DeviceID, &t.IPAddress, &t.Location, &t.Velocity, &t.ModelProbability, &t.Status, &t.CreatedAt)
			txs = append(txs, &t)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": txs})
	}
}

func handleCreateTransaction(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		var req Transaction
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		if req.TransactionID == "" {
			req.TransactionID = "TXN-" + uuid.New().String()[:8]
		}
		var id uuid.UUID
		err := db.QueryRow(r.Context(), `
			INSERT INTO transactions (id, org_id, transaction_id, amount, currency, merchant_id,
				customer_id, device_id, ip_address, location, velocity, model_probability, status, created_at)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'pending', NOW())
			RETURNING id
		`, orgID, req.TransactionID, req.Amount, req.Currency, req.MerchantID, req.CustomerID,
			req.DeviceID, req.IPAddress, req.Location, req.Velocity, req.ModelProbability).Scan(&id)
		if err != nil {
			http.Error(w, `{"error":"failed to create"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id.String(), "transaction_id": req.TransactionID})
	}
}

func handleAssessTransaction(db *store.DB, riskSvc *risk.Service, policySvc *policy.Service, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}

		var txn Transaction
		err = db.QueryRow(r.Context(), `
			SELECT id, transaction_id, amount, currency, merchant_id, customer_id,
				device_id, ip_address, location, velocity, model_probability, status
			FROM transactions WHERE id = $1 AND org_id = $2
		`, id, orgID).Scan(&txn.ID, &txn.TransactionID, &txn.Amount, &txn.Currency,
			&txn.MerchantID, &txn.CustomerID, &txn.DeviceID, &txn.IPAddress,
			&txn.Location, &txn.Velocity, &txn.ModelProbability, &txn.Status)
		if err != nil {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}

		// Calculate risk dimensions deterministically
		dims := calculateTransactionRisk(txn)
		score, err := riskSvc.CalculateRisk(r.Context(), orgID, id, "transaction", dims)
		if err != nil {
			http.Error(w, `{"error":"risk calculation failed"}`, http.StatusInternalServerError)
			return
		}

		// Policy evaluation
		ctx := map[string]interface{}{
			"fraud_probability": txn.ModelProbability,
		}
		decision, err := policySvc.Evaluate(r.Context(), orgID, score, ctx)
		if err != nil {
			http.Error(w, `{"error":"policy evaluation failed"}`, http.StatusInternalServerError)
			return
		}

		// Update transaction status
		status := "approved"
		if decision.Decision == "BLOCK" {
			status = "blocked"
		} else if decision.Decision == "REVIEW" {
			status = "review_required"
		}
		db.Exec(r.Context(), `UPDATE transactions SET status = $1 WHERE id = $2`, status, id)

		// Create approval request if review needed
		if decision.Decision == "REVIEW" || decision.Decision == "BLOCK" {
			var approvalID uuid.UUID
			db.QueryRow(r.Context(), `
				INSERT INTO approval_requests (id, org_id, resource_type, resource_id, risk_score_id,
					decision, status, created_at)
				VALUES (gen_random_uuid(), $1, 'transaction', $2, $3, $4, 'pending', NOW())
				RETURNING id
			`, orgID, id, score.ID, decision.Decision).Scan(&approvalID)
			decision.RequiredAction = "Approval request created: " + approvalID.String()
		}

		user := mw.GetUser(r.Context())
		auditSvc.Log(r.Context(), &audit.AuditEvent{
			ActorID:      user.UserID,
			Action:       "transaction_assessed",
			ResourceType: "transaction",
			ResourceID:   id,
			OrgID:        orgID,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"transaction_id": txn.TransactionID,
			"risk_score":     score,
			"policy_decision": decision,
			"status":         status,
		})
	}
}

func calculateTransactionRisk(txn Transaction) []risk.RiskDimension {
	var dims []risk.RiskDimension

	// Fraud dimension
	fraudScore := txn.ModelProbability * 100
	if txn.Velocity > 5 {
		fraudScore += 15
	}
	if txn.Amount > 50000 {
		fraudScore += 10
	}
	fraudScore = math.Min(100, fraudScore)
	dims = append(dims, risk.RiskDimension{
		Name:       "fraud",
		Score:      fraudScore,
		Weight:     0.25,
		Confidence: 0.85,
		ContributingFactors: []string{
			fmt.Sprintf("Model fraud probability: %.2f", txn.ModelProbability),
			fmt.Sprintf("Transaction velocity: %d", txn.Velocity),
			fmt.Sprintf("Amount: %.2f %s", txn.Amount, txn.Currency),
		},
		Evidence: []risk.Evidence{
			{ID: "fraud-1", Type: "model_score", Description: "Fraud model probability", Value: txn.ModelProbability, Timestamp: time.Now()},
		},
	})

	// Security dimension
	secScore := 30.0
	if txn.IPAddress != "" && (strings.HasPrefix(txn.IPAddress, "10.") || strings.HasPrefix(txn.IPAddress, "192.168")) {
		secScore -= 10 // Internal IP is lower risk
	}
	if txn.DeviceID == "" {
		secScore += 20
	}
	secScore = math.Min(100, math.Max(0, secScore))
	dims = append(dims, risk.RiskDimension{
		Name:       "security",
		Score:      secScore,
		Weight:     0.25,
		Confidence: 0.75,
		ContributingFactors: []string{
			fmt.Sprintf("IP Address: %s", txn.IPAddress),
			fmt.Sprintf("Device ID present: %v", txn.DeviceID != ""),
		},
		Evidence: []risk.Evidence{
			{ID: "sec-1", Type: "network", Description: "IP reputation check", Value: txn.IPAddress, Timestamp: time.Now()},
		},
	})

	// Privacy dimension
	privacyScore := 25.0
	if txn.Amount > 100000 {
		privacyScore += 15
	}
	dims = append(dims, risk.RiskDimension{
		Name:       "privacy",
		Score:      privacyScore,
		Weight:     0.15,
		Confidence: 0.70,
		ContributingFactors: []string{"High-value transaction requires PII handling review"},
		Evidence:   []risk.Evidence{{ID: "priv-1", Type: "data_classification", Description: "Transaction value classification", Value: txn.Amount, Timestamp: time.Now()}},
	})

	// Compliance
	complianceScore := 20.0
	if txn.Currency != "USD" && txn.Currency != "EUR" {
		complianceScore += 15
	}
	dims = append(dims, risk.RiskDimension{
		Name:       "compliance",
		Score:      complianceScore,
		Weight:     0.15,
		Confidence: 0.80,
		ContributingFactors: []string{fmt.Sprintf("Currency: %s requires cross-border compliance check", txn.Currency)},
		Evidence:   []risk.Evidence{{ID: "comp-1", Type: "regulatory", Description: "Currency compliance", Value: txn.Currency, Timestamp: time.Now()}},
	})

	// Fairness
	fairnessScore := 15.0
	dims = append(dims, risk.RiskDimension{
		Name:       "fairness",
		Score:      fairnessScore,
		Weight:     0.10,
		Confidence: 0.60,
		ContributingFactors: []string{"No demographic data available for fairness check"},
		Evidence:   []risk.Evidence{{ID: "fair-1", Type: "fairness", Description: "Fairness evaluation pending", Value: "N/A", Timestamp: time.Now()}},
	})

	// Reliability
	reliabilityScore := 20.0
	if txn.ModelProbability < 0.5 {
		reliabilityScore += 10 // Low confidence in model
	}
	dims = append(dims, risk.RiskDimension{
		Name:       "reliability",
		Score:      reliabilityScore,
		Weight:     0.05,
		Confidence: 0.75,
		ContributingFactors: []string{fmt.Sprintf("Model confidence: %.2f", txn.ModelProbability)},
		Evidence:   []risk.Evidence{{ID: "rel-1", Type: "model_confidence", Description: "Model reliability", Value: txn.ModelProbability, Timestamp: time.Now()}},
	})

	// Operational
	operationalScore := 15.0
	if txn.Velocity > 10 {
		operationalScore += 20
	}
	dims = append(dims, risk.RiskDimension{
		Name:       "operational",
		Score:      operationalScore,
		Weight:     0.05,
		Confidence: 0.80,
		ContributingFactors: []string{fmt.Sprintf("Velocity score: %d", txn.Velocity)},
		Evidence:   []risk.Evidence{{ID: "op-1", Type: "velocity", Description: "Transaction velocity", Value: txn.Velocity, Timestamp: time.Now()}},
	})

	return dims
}
