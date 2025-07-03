package pgsql

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/nullism/bqb"

	"github.com/dsaime/auth-api/internal/domain"
	sqlxRepo "github.com/dsaime/auth-api/internal/repository/pgsql/sqlx_repo"
)

type SessionRepository struct {
	sqlxRepo.SqlxRepo
}

func (r *SessionRepository) List(filter domain.SessionFilter) ([]domain.Session, error) {
	sel := bqb.New("SELECT * FROM sessions")
	where := bqb.Optional("WHERE")

	if filter.ID != uuid.Nil {
		where = where.And("id = ?", filter.ID)
	}

	q := bqb.New("? ?", sel, where)
	if r.IsTx() {
		q = q.Space("FOR UPDATE")
	}

	sql, args, err := q.ToPgsql()
	if err != nil {
		return nil, fmt.Errorf("q.ToPgsql: %w", err)
	}

	var sessions []dbSession
	if err = r.DB().Select(&sessions, sql, args...); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	return toDomainSessions(sessions), nil
}

func (r *SessionRepository) Upsert(session domain.Session) error {
	if session.ID == uuid.Nil {
		return fmt.Errorf("session ID is required")
	}

	if _, err := r.DB().NamedExec(`
		INSERT INTO sessions(id, user_id, user_agent, ip, status, expiry, refresh_token) 
		VALUES (:id, :user_id, :user_agent, :ip, :status, :expiry, :refresh_token)
		ON CONFLICT (id) DO UPDATE SET
			user_id=excluded.user_id,
			user_agent=excluded.user_agent,
			ip=excluded.ip,
			status=excluded.status,
			expiry=excluded.expiry,
			refresh_token=excluded.refresh_token
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
	ID               uuid.UUID `db:"id"`
	UserID           uuid.UUID `db:"user_id"`
	UserAgent        string    `db:"user_agent"`
	IP               net.IP    `db:"ip"`
	Status           string    `db:"status"`
	Expiry           time.Time `db:"expiry"`
	RefreshTokenHash string    `db:"refresh_token_hash"`
}

func toDBSession(session domain.Session) dbSession {
	return dbSession{
		ID:               session.ID,
		UserID:           session.UserID,
		UserAgent:        session.UserAgent,
		IP:               session.IP,
		Status:           session.Status,
		Expiry:           session.Expiry,
		RefreshTokenHash: session.RefreshTokenHash,
	}
}

func toDomainSession(session dbSession) domain.Session {
	// Улучшить: поля скрыть под геттеры, создавать экземпляр
	// только через конструктор domain.MustNewSession(...)
	return domain.Session{
		ID:               session.ID,
		UserID:           session.UserID,
		UserAgent:        session.UserAgent,
		IP:               session.IP,
		Status:           session.Status,
		Expiry:           session.Expiry,
		RefreshTokenHash: session.RefreshTokenHash,
	}
}

func toDomainSessions(sessions []dbSession) []domain.Session {
	domainSessions := make([]domain.Session, len(sessions))
	for i, s := range sessions {
		domainSessions[i] = toDomainSession(s)
	}

	return domainSessions
}
