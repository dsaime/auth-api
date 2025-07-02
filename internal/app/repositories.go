package app

import (
	"fmt"
	"log/slog"

	"github.com/dsaime/auth-api/internal/domain"
	"github.com/dsaime/auth-api/internal/repository/pgsql"
)

type repositories struct {
	sessions domain.SessionRepository
}

func initPgsqlRepositories(cfg pgsql.Config) (*repositories, func(), error) {
	factory, err := pgsql.InitFactory(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("pgsql.InitFactory: %w", err)
	}

	rs := &repositories{
		sessions: factory.NewSessionnRepository(),
	}

	closer := func() {
		if err := factory.Close(); err != nil {
			slog.Error("Закрыть соединение с pgsql: factory.Close: " + err.Error())
		}
	}

	return rs, closer, nil
}
