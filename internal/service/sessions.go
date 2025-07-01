package service

import (
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/domain"
)

// Auth implements the authentication interface
type Auth struct {
	Repo domain.SessionRepository
}

type SessionsFindIn struct {
	TokenID string
}

func (s *Auth) Find(in SessionsFindIn) ([]domain.Session, error) {
	if in.TokenID == "" {
		return nil, ErrInvalidToken
	}

	sessions, err := s.Repo.List(domain.SessionFilter{ID: in.TokenID})
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

type AuthLoginIn struct {
	UserID    uuid.UUID
	UserAgent string
}

type AuthLoginOut struct {
	TokenID      string
	RefreshToken string
	Session      domain.Session
}

func (s *Auth) Login(in AuthLoginIn) (out AuthLoginOut, err error) {
	panic("not implemented")
}

type AuthRefreshIn struct {
	AccessToken  string
	RefreshToken string
}
type AuthRefreshOut struct {
	TokenID      string
	RefreshToken string
}

func (s *Auth) Refresh(in AuthRefreshIn) (AuthRefreshOut, error) {
	panic("not implemented")
}

type AuthLogoutIn struct {
	TokenID string
}
type AuthLogoutOut struct {
}

func (s *Auth) Logout(in AuthLogoutIn) (AuthLogoutOut, error) {
	panic("not implemented")
}
