package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/riskshield-ai/backend/internal/store"
	"golang.org/x/crypto/argon2"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

type Service struct {
	db     *store.DB
	secret string
}

func NewService(db *store.DB, secret string) *Service {
	return &Service{db: db, secret: secret}
}

type Claims struct {
	UserID         uuid.UUID `json:"user_id"`
	OrgID          uuid.UUID `json:"org_id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	jwt.RegisteredClaims
}

func (s *Service) Register(ctx context.Context, orgID uuid.UUID, email, password, name string) (*User, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	var user User
	err = s.db.QueryRow(ctx, `
		INSERT INTO users (id, org_id, email, password_hash, name, role, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, 'ORG_ADMIN', NOW(), NOW())
		RETURNING id, org_id, email, name, role, created_at
	`, orgID, email, hash, name).Scan(&user.ID, &user.OrgID, &user.Email, &user.Name, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return &user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	var user User
	var hash string
	err := s.db.QueryRow(ctx, `
		SELECT id, org_id, email, name, role, password_hash, created_at
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.OrgID, &user.Email, &user.Name, &user.Role, &hash, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !verifyPassword(password, hash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return s.generateTokens(&user)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// In production: validate refresh token against DB, check revocation
	// Simplified for demo
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	var user User
	err = s.db.QueryRow(ctx, `
		SELECT id, org_id, email, name, role, created_at FROM users WHERE id = $1
	`, claims.UserID).Scan(&user.ID, &user.OrgID, &user.Email, &user.Name, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s.generateTokens(&user)
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

func (s *Service) generateTokens(user *User) (*TokenPair, error) {
	now := time.Now()
	accessClaims := Claims{
		UserID: user.ID,
		OrgID:  user.OrgID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "riskshield",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.secret))
	if err != nil {
		return nil, err
	}

	refreshClaims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "riskshield",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	}, nil
}

type User struct {
	ID        uuid.UUID `json:"id"`
	OrgID     uuid.UUID `json:"org_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// Argon2id password hashing
func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("$argon2id$v=19$m=65536,t=3,p=4$%s$%s", b64Salt, b64Hash), nil
}

func verifyPassword(password, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}
	var version int
	fmt.Sscanf(parts[2], "v=%d", &version)
	var memory uint32
	var iterations uint32
	var parallelism uint8
	fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(expectedHash)))
	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}
