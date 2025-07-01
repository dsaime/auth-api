package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

// Session представляет собой агрегат сессии
type Session struct {
	ID           uuid.UUID // ID сессии
	UserID       uuid.UUID // ID пользователя, к которому относится сессия
	UserAgent    string    // [название модели телефона / название браузера]
	Status       string    // Статус сессии
	Expiry       time.Time // Дата истечения сессии
	RefreshToken string    // Токен для обновления AccessToken
}

const (
	SessionStatusNew      = "new"      // Новая
	SessionStatusVerified = "verified" // Подтвержденная
	SessionStatusExpired  = "expired"  // Истекшая
	SessionStatusRevoked  = "revoked"  // Отозванная
)

var SessionStatuses = []string{
	SessionStatusNew,
	SessionStatusVerified,
	SessionStatusExpired,
	SessionStatusRevoked,
}

const (
	//refreshTokenLifetime = 60 * 24 * time.Hour // 60 дней
	accessTokenLifetime = 15 * time.Minute // 15 минут
)

// NewSession создает новую сессию, связанную с пользователем.
func NewSession(userID uuid.UUID, agent, status string) (Session, error) {
	if err := ValidateID(userID); err != nil {
		return Session{}, err
	}
	if err := ValidateSessionAgent(agent); err != nil {
		return Session{}, err
	}
	if err := ValidateSessionStatus(status); err != nil {
		return Session{}, err
	}

	return Session{
		ID:           uuid.New(),
		UserID:       userID,
		UserAgent:    agent,
		Status:       status,
		Expiry:       time.Now().Add(accessTokenLifetime).In(time.UTC).Truncate(time.Microsecond),
		RefreshToken: uuid.NewString(),
	}, nil
}

type SessionRepository interface {
	List(SessionFilter) ([]Session, error)
	Upsert(Session) error
	InTransaction(func(txRepo SessionRepository) error) error
}

// SessionFilter представляет собой фильтр по сессиям.
type SessionFilter struct {
	ID string
}

// ValidateSessionStatus проверяет корректность статуса сессии.
func ValidateSessionStatus(status string) error {
	if !slices.Contains(SessionStatuses, status) {
		return ErrSessionStatusValidate
	}

	return nil
}

// ValidateSessionAgent проверяет корректность названия сессии.
func ValidateSessionAgent(name string) error {
	if name == "" {
		return ErrSessionNameEmpty
	}

	return nil
}

var (
	ErrSessionStatusValidate = errors.New("некорректный статус сессии")
	ErrSessionNameEmpty      = errors.New("название сессии не может быть пустым")
)
