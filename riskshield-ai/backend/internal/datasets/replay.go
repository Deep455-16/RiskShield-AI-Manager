package datasets

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/riskshield-ai/backend/internal/policy"
	"github.com/riskshield-ai/backend/internal/risk"
	"github.com/riskshield-ai/backend/internal/store"
	"github.com/google/uuid"
)

// Ensure policy import used
var _ *policy.PolicyDecision

type ReplayEngine struct {
	Registry  *Registry
	DB        *store.DB
	RiskSvc   *risk.Service
	PolicySvc *policy.Service
	
	activeStreams map[string]context.CancelFunc
}

func NewReplayEngine(registry *Registry, db *store.DB, riskSvc *risk.Service, policySvc *policy.Service) *ReplayEngine {
	return &ReplayEngine{
		Registry: registry,
		DB: db,
		RiskSvc: riskSvc,
		PolicySvc: policySvc,
		activeStreams: make(map[string]context.CancelFunc),
	}
}

func (e *ReplayEngine) StartReplay(datasetID string, speedMultiplier float64) error {
	adapter := e.Registry.GetAdapter(datasetID)
	if adapter == nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.activeStreams[datasetID] = cancel

	ch := make(chan CanonicalTransaction, 100)
	
	go adapter.ReadStream(ch, 1000) // max 1000 rows for demo

	go func() {
		defer cancel()
		baseDelay := 1000.0 // 1 sec default
		delayMs := baseDelay / speedMultiplier

		orgID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		assetID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		for {
			select {
			case <-ctx.Done():
				return
			case tx, ok := <-ch:
				if !ok {
					return // End of stream
				}

				// 1. Feature extraction / Fraud Model placeholder
				fraudScore := rand.Float64() * 20
				if tx.FraudLabel != nil && *tx.FraudLabel {
					fraudScore += 60 // Boost score if actually fraud
				}

				dims := []risk.RiskDimension{
					{Name: "fraud", Score: fraudScore, Weight: 1.0, Confidence: 0.9, ContributingFactors: []string{"Replay Model"}, Evidence: []risk.Evidence{}},
				}

				// 2. Risk Engine
				score, _ := e.RiskSvc.CalculateRisk(context.Background(), orgID, assetID, "replay_"+tx.TransactionID, dims)

				// 3. Policy Engine
				decision, _ := e.PolicySvc.Evaluate(context.Background(), orgID, score, map[string]interface{}{"amount": tx.Amount})

				// 4. DB Insert
				status := "approved"
				if decision.Decision == "BLOCK" {
					status = "blocked"
				}

				_, err := e.DB.Query(context.Background(), `
					INSERT INTO transactions (id, org_id, amount, currency, status, risk_score, source_dataset, synthetic)
					VALUES ($1, $2, $3, $4, $5, $6, $7, true)
				`, tx.TransactionID, orgID, tx.Amount, tx.Currency, status, score, tx.SourceDataset)
				
				if err != nil {
					log.Printf("Replay insert error: %v", err)
				}

				// If blocked, optionally insert incident
				if decision.Decision == "BLOCK" {
					e.DB.Query(context.Background(), `
						INSERT INTO incidents (id, org_id, title, description, severity, status, category)
						VALUES (gen_random_uuid(), $1, $2, $3, 'HIGH', 'OPEN', 'fraud')
					`, orgID, "High Risk Transaction: "+tx.TransactionID, "Blocked by Policy. Source: "+tx.SourceDataset)
				}

				time.Sleep(time.Duration(delayMs) * time.Millisecond)
			}
		}
	}()

	return nil
}

func (e *ReplayEngine) StopReplay(datasetID string) {
	if cancel, exists := e.activeStreams[datasetID]; exists {
		cancel()
		delete(e.activeStreams, datasetID)
	}
}
