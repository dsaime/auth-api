package service

import (
	"errors"
	"net"

	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/domain"
)

// Auth implements the authentication interface
type Auth struct {
	Repo    domain.SessionRepository
	Alerter Alerter
}

type Alerter interface {
	UnknownIPRefreshSessionAlert(sessionID uuid.UUID, oldIP, newIP net.IP)
}

type AuthLoginIn struct {
	UserID    uuid.UUID
	UserAgent string
	IP        net.IP
}

type AuthLoginOut struct {
	RefreshToken string
	Session      domain.Session
}

// Login создает новую активную сессию пользователя
func (s *Auth) Login(in AuthLoginIn) (AuthLoginOut, error) {
	// TODO: in.Validate()

	// Случайное значение токена
	refreshToken := uuid.NewString()

	// Создать новую сессию
	session, err := domain.NewSession(in.UserID, in.UserAgent, in.IP, refreshToken)
	if err != nil {
		return AuthLoginOut{}, err
	}

	// Сохранить сессию
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
	IP           net.IP
}

type AuthRefreshOut struct {
	Session domain.Session
}

var ErrUnauthorized = errors.New("unauthorized")

// Refresh продлевает активную или истекшую сессию
func (s *Auth) Refresh(in AuthRefreshIn) (AuthRefreshOut, error) {
	// TODO: in.Validate()

	var updatedSession domain.Session
	if err := s.Repo.InTransaction(func(txRepo domain.SessionRepository) error {
		// Найти сессию по ID
		session, err := s.findSession(domain.SessionFilter{
			ID: in.SessionID,
		})
		if err != nil {
			return err
		}

		// Сравнить токен из параметров с хэшем токена в сессии
		if err = session.CompareRefreshTokenWithHash(in.RefreshToken); err != nil {
			return err
		}

		// UserAgent из параметров должен совпадать со значением в сессии
		if in.UserAgent != session.UserAgent {
			// Деактивировать сессию и сохранить
			if err = session.Revoke(); err != nil {
				return err
			}
			if err = s.Repo.Upsert(session); err != nil {
				return err
			}

			return ErrUnauthorized
		}

		// Обновление с нового IP будет фиксироваться.
		// Улучшить: Если состояние и уведомление должны быть согласованы,
		// можно добавить буфер уведомлений, а если требуется гарантия отправки,
		// то писать под той же транзакцией в хранилище.
		if !session.IP.Equal(in.IP) {
			session.UpdateIP(in.IP)
			go s.Alerter.UnknownIPRefreshSessionAlert(session.ID, session.IP, in.IP)
		}

		// Продлить и сохранить сессию
		if err = session.ExtendExpiry(); err != nil {
			return err
		}
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

// Logout отзывает (деактивирует сессию)
func (s *Auth) Logout(in AuthLogoutIn) error {
	// TODO: in.Validate()

	return s.Repo.InTransaction(func(txRepo domain.SessionRepository) error {
		// Найти сессию по ID
		session, err := s.findSession(domain.SessionFilter{
			ID: in.SessionID,
		})
		if err != nil {
			return err
		}

		// Отозвать сессию
		if err = session.Revoke(); err != nil {
			return err
		}

		// Сохранить сессию
		return s.Repo.Upsert(session)
	})
}

func (s *Auth) findSession(filter domain.SessionFilter) (domain.Session, error) {
	sessions, err := s.Repo.List(filter)
	if err != nil {
		return domain.Session{}, err
	}

	if len(sessions) != 1 {
		return domain.Session{}, errors.New("не найдено")
	}

	return sessions[0], nil
}
