package jobs

import (
	"context"
	"time"

	"github.com/riskshield-ai/backend/internal/store"
)

type Service struct {
	db *store.DB
}

func NewService(db *store.DB) *Service {
	return &Service{db: db}
}

type Job struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Payload   []byte    `json:"payload"`
	Result    []byte    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) Create(ctx context.Context, jobType string, payload []byte) (*Job, error) {
	var job Job
	err := s.db.QueryRow(ctx, `
		INSERT INTO jobs (id, type, status, payload, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, 'pending', $2, NOW(), NOW())
		RETURNING id, type, status, payload, created_at, updated_at
	`, jobType, payload).Scan(&job.ID, &job.Type, &job.Status, &job.Payload, &job.CreatedAt, &job.UpdatedAt)
	return &job, err
}
