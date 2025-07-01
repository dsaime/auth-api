package pgsqlRepository

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/domain"
	sqlxRepo "github.com/dsaime/auth-api/internal/repository/pgsql_repository/sqlx_repo"
)

type SessionRepository struct {
	sqlxRepo.SqlxRepo
}

func (r *SessionRepository) List(filter domain.SessionFilter) ([]domain.Session, error) {
	var sessions []dbSession
	if err := r.DB().Select(&sessions, `
		SELECT *
		FROM sessions
		WHERE ($1 = '' OR $1 = access_token)
	`, filter.ID); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	return toDomainSessions(sessions), nil
}

func (r *SessionRepository) Upsert(session domain.Session) error {
	if session.ID == uuid.Nil {
		return fmt.Errorf("session ID is required")
	}

	if _, err := r.DB().NamedExec(`
		INSERT INTO sessions(id, user_id, name, status, access_token, access_expiry, refresh_token, refresh_expiry) 
		VALUES (:id, :user_id, :name, :status, :access_token, :access_expiry, :refresh_token, :refresh_expiry)
		ON CONFLICT (id) DO UPDATE SET
			user_id=excluded.user_id,
			name=excluded.name,
			status=excluded.status,
			access_token=excluded.access_token,
			access_expiry=excluded.access_expiry,
			refresh_token=excluded.refresh_token,
			refresh_expiry=excluded.refresh_expiry
	`, toDBSession(session)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	return nil
}

func (r *SessionRepository) InTransaction(fn func(txRepo domain.SessionRepository) error) error {
	return r.SqlxRepo.InTransaction(func(txSqlxRepo sqlxRepo.SqlxRepo) error {
		return fn(&SessionRepository{SqlxRepo: txSqlxRepo})
	})
}

type dbSession struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	UserAgent    string    `db:"user_agent"`
	Status       string    `db:"status"`
	Expiry       time.Time `db:"expiry"`
	RefreshToken string    `db:"refresh_token"`
}

func toDBSession(session domain.Session) dbSession {
	return dbSession{
		ID:           session.ID.String(),
		UserID:       session.UserID.String(),
		UserAgent:    session.UserAgent,
		Status:       session.Status,
		Expiry:       session.Expiry,
		RefreshToken: session.RefreshToken,
	}
}

func toDomainSession(session dbSession) domain.Session {
	return domain.Session{
		ID:           uuid.MustParse(session.ID),
		UserID:       uuid.MustParse(session.UserID),
		UserAgent:    session.UserAgent,
		Status:       session.Status,
		Expiry:       session.Expiry,
		RefreshToken: session.RefreshToken,
	}
}

func toDomainSessions(sessions []dbSession) []domain.Session {
	domainSessions := make([]domain.Session, len(sessions))
	for i, s := range sessions {
		domainSessions[i] = toDomainSession(s)
	}

	return domainSessions
}
