package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func TenantIsolation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		if user == nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r.Context())
			if user == nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			allowed := false
			for _, role := range roles {
				if user.Role == role {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func OrgIDFromContext(ctx context.Context) uuid.UUID {
	user := GetUser(ctx)
	if user == nil {
		return uuid.Nil
	}
	return user.OrgID
}
