package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SupabaseJWTPayload struct {
	Sub         string `json:"sub"`
	Aud         string `json:"aud"`
	Exp         int64  `json:"exp"`
	Iat         int64  `json:"iat"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	AppMetadata struct {
		Provider  string   `json:"provider"`
		Providers []string `json:"providers"`
	} `json:"app_metadata"`
	UserMetadata struct {
		Sub       string `json:"sub"`
		Iss       string `json:"iss"`
		Picture   string `json:"picture"`
		Name      string `json:"name"`
		Nickname  string `json:"nickname,omitempty"`
		Email     string `json:"email,omitempty"`
		Slug      string `json:"slug,omitempty"`
		AvatarURL string `json:"avatar_url,omitempty"`
	} `json:"user_metadata"`
	Role        string `json:"role"`
	SessionID   string `json:"session_id"`
	IsAnonymous bool   `json:"is_anonymous"`
}

// Implementing jwt.Claims interface
func (s *SupabaseJWTPayload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(s.Exp, 0)}, nil
}

func (s *SupabaseJWTPayload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(s.Iat, 0)}, nil
}

func (s *SupabaseJWTPayload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (s *SupabaseJWTPayload) GetIssuer() (string, error) {
	return s.UserMetadata.Iss, nil
}

func (s *SupabaseJWTPayload) GetSubject() (string, error) {
	return s.Sub, nil
}

func (s *SupabaseJWTPayload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{s.Aud}, nil
}
