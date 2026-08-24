package middleware

import (
	"net/http"
	"time"

	"github.com/riskshield-ai/backend/internal/audit"
)

func AuditMiddleware(auditSvc *audit.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			// Fire-and-forget audit log for API calls
			go auditSvc.LogAPIRequest(r, duration)
		})
	}
}
