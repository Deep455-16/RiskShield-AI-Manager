package api

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/audit"
	"github.com/riskshield-ai/backend/internal/policy"
	"github.com/riskshield-ai/backend/internal/risk"
	mw "github.com/riskshield-ai/backend/internal/middleware"
)

func handleSimulatePayment(riskSvc *risk.Service, policySvc *policy.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		var req struct {
			Amount          float64 `json:"amount"`
			Currency        string  `json:"currency"`
			Velocity        int     `json:"velocity"`
			DeviceNovelty   bool    `json:"device_novelty"`
			LocationNovelty bool    `json:"location_novelty"`
			MerchantRisk    float64 `json:"merchant_risk"`
			IPReputation    float64 `json:"ip_reputation"`
			ModelConfidence float64 `json:"model_confidence"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		fraudScore := req.ModelConfidence * 100
		if req.Velocity > 3 { fraudScore += float64(req.Velocity) * 5 }
		if req.DeviceNovelty { fraudScore += 15 }
		if req.LocationNovelty { fraudScore += 10 }
		fraudScore += req.MerchantRisk * 20
		fraudScore += (1 - req.IPReputation) * 30
		fraudScore = math.Min(100, fraudScore)

		securityScore := 30.0 + (1-req.IPReputation)*50
		if req.DeviceNovelty { securityScore += 15 }
		securityScore = math.Min(100, securityScore)

		privacyScore := 20.0
		if req.Amount > 100000 { privacyScore += 20 }

		complianceScore := 25.0
		if req.Currency != "USD" { complianceScore += 15 }

		fairnessScore := 15.0
		reliabilityScore := 20.0 + (1-req.ModelConfidence)*30
		operationalScore := 15.0 + float64(req.Velocity)*3

		dims := []risk.RiskDimension{
			{Name: "fraud", Score: fraudScore, Weight: 0.25, Confidence: 0.85, ContributingFactors: []string{"Simulated fraud model"}, Evidence: []risk.Evidence{{ID: "sim-fraud", Type: "simulation", Description: "Fraud simulation", Value: fraudScore, Timestamp: time.Now()}}},
			{Name: "security", Score: securityScore, Weight: 0.25, Confidence: 0.75, ContributingFactors: []string{"IP reputation simulation"}, Evidence: []risk.Evidence{{ID: "sim-sec", Type: "simulation", Description: "Security simulation", Value: securityScore, Timestamp: time.Now()}}},
			{Name: "privacy", Score: privacyScore, Weight: 0.15, Confidence: 0.70, ContributingFactors: []string{"Amount-based privacy risk"}, Evidence: []risk.Evidence{{ID: "sim-priv", Type: "simulation", Description: "Privacy simulation", Value: privacyScore, Timestamp: time.Now()}}},
			{Name: "compliance", Score: complianceScore, Weight: 0.15, Confidence: 0.80, ContributingFactors: []string{"Currency compliance check"}, Evidence: []risk.Evidence{{ID: "sim-comp", Type: "simulation", Description: "Compliance simulation", Value: complianceScore, Timestamp: time.Now()}}},
			{Name: "fairness", Score: fairnessScore, Weight: 0.10, Confidence: 0.60, ContributingFactors: []string{"No demographic data"}, Evidence: []risk.Evidence{{ID: "sim-fair", Type: "simulation", Description: "Fairness simulation", Value: fairnessScore, Timestamp: time.Now()}}},
			{Name: "reliability", Score: reliabilityScore, Weight: 0.05, Confidence: 0.75, ContributingFactors: []string{"Model confidence"}, Evidence: []risk.Evidence{{ID: "sim-rel", Type: "simulation", Description: "Reliability simulation", Value: reliabilityScore, Timestamp: time.Now()}}},
			{Name: "operational", Score: operationalScore, Weight: 0.05, Confidence: 0.80, ContributingFactors: []string{"Velocity check"}, Evidence: []risk.Evidence{{ID: "sim-op", Type: "simulation", Description: "Operational simulation", Value: operationalScore, Timestamp: time.Now()}}},
		}

		assetID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		score, _ := riskSvc.CalculateRisk(r.Context(), orgID, assetID, "simulation", dims)
		decision, _ := policySvc.Evaluate(r.Context(), orgID, score, map[string]interface{}{})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"simulation": "payment_risk",
			"risk_score": score,
			"policy_decision": decision,
			"dimensions": dims,
			"warning": "SIMULATION — NO REAL PAYMENT PROCESSING",
		})
	}
}

func handleSimulateAttack(riskSvc *risk.Service, policySvc *policy.Service, auditSvc *audit.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := mw.OrgIDFromContext(r.Context())
		var req struct {
			AttackType string `json:"attack_type"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		baselineScore := 34.0
		baselineLevel := "LOW"

		attackImpact := map[string]float64{
			"prompt_injection": 45, "pii_leakage": 30, "suspicious_transaction": 50,
			"privilege_escalation": 60, "model_drift": 25, "fairness_disparity": 20,
			"api_abuse": 35, "excessive_spending": 15,
		}
		impact := attackImpact[req.AttackType]
		if impact == 0 { impact = 20 }

		finalScore := math.Min(100, baselineScore+impact)
		var finalLevel string
		if finalScore >= 75 { finalLevel = "CRITICAL"
		} else if finalScore >= 50 { finalLevel = "HIGH"
		} else if finalScore >= 25 { finalLevel = "MEDIUM"
		} else { finalLevel = "LOW" }

		dims := []risk.RiskDimension{
			{Name: "security", Score: finalScore, Weight: 0.40, Confidence: 0.90, ContributingFactors: []string{"Simulated attack: " + req.AttackType}, Evidence: []risk.Evidence{{ID: "attack-1", Type: "attack_simulation", Description: req.AttackType, Value: impact, Timestamp: time.Now()}}},
		}

		assetID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
		score, _ := riskSvc.CalculateRisk(r.Context(), orgID, assetID, "attack_simulation", dims)

		var incidentID uuid.UUID
		incidentID = uuid.New()

		user := mw.GetUser(r.Context())
		auditSvc.Log(r.Context(), &audit.AuditEvent{
			ActorID: user.UserID, Action: "attack_simulated",
			ResourceType: "incident", ResourceID: incidentID, OrgID: orgID,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"attack_type": req.AttackType,
			"before": map[string]interface{}{"score": baselineScore, "level": baselineLevel},
			"after": map[string]interface{}{"score": finalScore, "level": finalLevel},
			"impact": impact,
			"risk_score": score,
			"incident_created": incidentID,
			"mitigation": "Policy engine automatically escalated to human review",
		})
	}
}
