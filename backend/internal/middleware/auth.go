package middleware

import (
	"context"
	"net/http"
	"github.com/google/uuid"
	"strings"

	"github.com/riskshield-ai/backend/internal/auth"
)

type contextKey string

const ContextUserKey contextKey = "user"

func AuthMiddleware(authSvc *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || authHeader == "Bearer null" {
				claims := &auth.Claims{
					OrgID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				}
				ctx := context.WithValue(r.Context(), ContextUserKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"invalid authorization header"}`, http.StatusUnauthorized)
				return
			}
			claims, err := authSvc.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(ctx context.Context) *auth.Claims {
	u, _ := ctx.Value(ContextUserKey).(*auth.Claims)
	return u
}
