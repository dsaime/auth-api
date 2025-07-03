package domain

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

// Session представляет собой агрегат сессии
type Session struct {
	ID               uuid.UUID // ID сессии
	UserID           uuid.UUID // ID пользователя, к которому относится сессия
	UserAgent        string    // Информация о браузере, операционной системе и устройстве
	Status           string    // Статус сессии
	Expiry           time.Time // Дата истечения сессии
	RefreshTokenHash string    // Хэш токен для продления сессии. Читается из хранилища
	//RefreshToken string // Токен для продления сессии. Передается в хранилище
	//CreatedAt time.Time // Дата создания сессии
	//LastActivityAt time.Time // Дата последней активности этой сессии
	//AccessToken string // Токен для идентификации сессии, изменяемый
	IP net.IP // Адрес устройства
}

const (
	SessionStatusActive  = "active"  // Активная
	SessionStatusExpired = "expired" // Истекшая
	SessionStatusRevoked = "revoked" // Отозванная пользователем или сервером
)

var SessionStatuses = []string{
	SessionStatusActive,
	SessionStatusExpired,
	SessionStatusRevoked,
}

const (
	accessTokenLifetime = 15 * time.Minute // 15 минут
)

// NewSession создает новую сессию, связанную с пользователем.
func NewSession(userID uuid.UUID, agent string, ip net.IP, refreshToken string) (Session, error) {
	if err := ValidateID(userID); err != nil {
		return Session{}, err
	}
	if err := ValidateSessionAgent(agent); err != nil {
		return Session{}, err
	}
	if ip.IsUnspecified() {
		return Session{}, errors.New("некорректный IP")
	}
	if refreshToken == "" {
		return Session{}, errors.New("некорректный refreshToken")
	}

	// Получить хэш токена
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return Session{}, err
	}

	return Session{
		ID:               uuid.New(),
		UserID:           userID,
		UserAgent:        agent,
		Status:           SessionStatusActive,
		Expiry:           newSessionExpiryTime(),
		IP:               ip,
		RefreshTokenHash: string(hashedRefreshToken),
	}, nil
}

type SessionRepository interface {
	List(SessionFilter) ([]Session, error)
	Upsert(Session) error
	InTransaction(func(txRepo SessionRepository) error) error
}

// SessionFilter представляет собой фильтр по сессиям.
type SessionFilter struct {
	ID uuid.UUID
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

func (s *Session) CompareRefreshTokenWithHash(rt string) error {
	return bcrypt.CompareHashAndPassword([]byte(s.RefreshTokenHash), []byte(rt))
}

func (s *Session) Revoke() error {
	if s.Status == SessionStatusRevoked {
		return errors.New("сессия уже отозвана")
	}

	s.Status = SessionStatusRevoked

	return nil
}

var sessionsMaysBeExtended = []string{
	SessionStatusActive,
	SessionStatusExpired,
}

func (s *Session) ExtendExpiry() error {
	// Проверить возможность продлить сессию в ее статусе
	if !slices.Contains(sessionsMaysBeExtended, s.Status) {
		return fmt.Errorf("невозможно продлить сессию в статусе '%s'", s.Status)
	}

	s.Status = SessionStatusActive
	s.Expiry = newSessionExpiryTime()

	return nil
}

func newSessionExpiryTime() time.Time {
	return time.Now().Add(accessTokenLifetime).
		In(time.UTC).              // Для единообразия и возможности сравнивать в тестах всю сессию
		Truncate(time.Microsecond) // Чтобы значение полностью помещалось в БД
}

func (s *Session) UpdateIP(ip net.IP) {
	s.IP = ip
}
