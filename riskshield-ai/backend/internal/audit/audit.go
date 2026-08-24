package audit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/store"
)

type Service struct {
	db *store.DB
}

func NewService(db *store.DB) *Service {
	return &Service{db: db}
}

type AuditEvent struct {
	ID           uuid.UUID       `json:"id"`
	ActorID      uuid.UUID       `json:"actor_id"`
	Action       string          `json:"action"`
	ResourceType string          `json:"resource_type"`
	ResourceID   uuid.UUID       `json:"resource_id"`
	OrgID        uuid.UUID       `json:"org_id"`
	IP           string          `json:"ip"`
	RequestID    string          `json:"request_id"`
	Before       json.RawMessage `json:"before,omitempty"`
	After        json.RawMessage `json:"after,omitempty"`
	Hash         string          `json:"hash"`
	Timestamp    time.Time       `json:"timestamp"`
}

func (s *Service) Log(ctx context.Context, event *AuditEvent) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	event.Timestamp = time.Now().UTC()

	// Build hash chain: previous hash + current content
	prevHash, _ := s.getPreviousHash(ctx, event.OrgID)
	content := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s",
		prevHash, event.ActorID, event.Action, event.ResourceType,
		event.ResourceID, event.OrgID, event.IP, event.Timestamp.Format(time.RFC3339Nano))
	h := sha256.Sum256([]byte(content))
	event.Hash = hex.EncodeToString(h[:])

	_, err := s.db.Pool().Exec(ctx, `
		INSERT INTO audit_logs (id, actor_id, action, resource_type, resource_id, org_id,
			ip_address, request_id, before_state, after_state, hash, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, event.ID, event.ActorID, event.Action, event.ResourceType, event.ResourceID,
		event.OrgID, event.IP, event.RequestID, event.Before, event.After, event.Hash, event.Timestamp)
	return err
}

func (s *Service) LogAPIRequest(r *http.Request, duration time.Duration) {
	// Best-effort logging
}

func (s *Service) getPreviousHash(ctx context.Context, orgID uuid.UUID) (string, error) {
	var hash string
	err := s.db.QueryRow(ctx, `
		SELECT hash FROM audit_logs WHERE org_id = $1 ORDER BY timestamp DESC LIMIT 1
	`, orgID).Scan(&hash)
	if err != nil {
		return "genesis", nil
	}
	return hash, nil
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*AuditEvent, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, actor_id, action, resource_type, resource_id, org_id,
			ip_address, request_id, before_state, after_state, hash, timestamp
		FROM audit_logs WHERE org_id = $1 ORDER BY timestamp DESC LIMIT $2 OFFSET $3
	`, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*AuditEvent
	for rows.Next() {
		var e AuditEvent
		err := rows.Scan(&e.ID, &e.ActorID, &e.Action, &e.ResourceType, &e.ResourceID,
			&e.OrgID, &e.IP, &e.RequestID, &e.Before, &e.After, &e.Hash, &e.Timestamp)
		if err != nil {
			continue
		}
		events = append(events, &e)
	}
	return events, rows.Err()
}
