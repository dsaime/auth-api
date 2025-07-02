package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/domain"
)

// Auth implements the authentication interface
type Auth struct {
	Repo domain.SessionRepository
}

type AuthLoginIn struct {
	UserID    uuid.UUID
	UserAgent string
}

type AuthLoginOut struct {
	RefreshToken string
	Session      domain.Session
}

func (s *Auth) Login(in AuthLoginIn) (AuthLoginOut, error) {
	// TODO: in.Validate()

	refreshToken := uuid.NewString()
	session, err := domain.NewSession(in.UserID, in.UserAgent, domain.SessionStatusVerified, refreshToken)
	if err != nil {
		return AuthLoginOut{}, err
	}

	if err = s.Repo.Upsert(session); err != nil {
		return AuthLoginOut{}, err
	}

	return AuthLoginOut{
		RefreshToken: refreshToken,
		Session:      session,
	}, nil
}

type AuthRefreshIn struct {
	SessionID    uuid.UUID
	RefreshToken string
	UserAgent    string
}

type AuthRefreshOut struct {
	Session domain.Session
}

var ErrUnauthorized = errors.New("unauthorized")

func (s *Auth) Refresh(in AuthRefreshIn) (AuthRefreshOut, error) {
	// TODO: in.Validate()

	var updatedSession domain.Session
	if err := s.Repo.InTransaction(func(txRepo domain.SessionRepository) error {
		session, err := s.findSession(domain.SessionFilter{
			ID: in.SessionID,
		})
		if err != nil {
			return err
		}

		if session.Status != domain.SessionStatusVerified {
			return ErrUnauthorized
		}

		if session.Expiry.Before(time.Now()) {
			return ErrUnauthorized
		}

		if err = session.CompareRefreshToken(in.RefreshToken); err != nil {
			return err
		}

		if in.UserAgent != session.UserAgent {
			session.Revoke()
			if err = s.Repo.Upsert(session); err != nil {
				return err
			}

			return ErrUnauthorized
		}

		session.ExtendExpiry()
		if err = s.Repo.Upsert(session); err != nil {
			return err
		}
		updatedSession = session
		return nil
	}); err != nil {
		return AuthRefreshOut{}, err
	}

	return AuthRefreshOut{
		Session: updatedSession,
	}, nil
}

type AuthLogoutIn struct {
	SessionID uuid.UUID
}

func (s *Auth) Logout(in AuthLogoutIn) error {
	// TODO: in.Validate()

	session, err := s.findSession(domain.SessionFilter{
		ID: in.SessionID,
	})
	if err != nil {
		return err
	}

	if session.Status != domain.SessionStatusVerified {
		return ErrUnauthorized
	}

	session.Revoke()
	if err = s.Repo.Upsert(session); err != nil {
		return err
	}

	return nil
}

func (s *Auth) findSession(filter domain.SessionFilter) (domain.Session, error) {
	sessions, err := s.Repo.List(filter)
	if err != nil {
		return domain.Session{}, err
	}

	if len(sessions) != 1 {
		return domain.Session{}, errors.New("not found")
	}

	return sessions[0], nil
}
