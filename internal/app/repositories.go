package app

import (
	"fmt"
	"log/slog"

	"github.com/dsaime/auth-api/internal/domain/chatt"
	"github.com/dsaime/auth-api/internal/domain/sessionn"
	"github.com/dsaime/auth-api/internal/domain/userr"
	pgsqlRepository "github.com/dsaime/auth-api/internal/repository/pgsql_repository"
)

type repositories struct {
	chats    chatt.Repository
	users    userr.Repository
	sessions sessionn.Repository
}

func initPgsqlRepositories(cfg pgsqlRepository.Config) (*repositories, func(), error) {
	factory, err := pgsqlRepository.InitFactory(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("pgsqlRepository.InitFactory: %w", err)
	}

	rs := &repositories{
		chats:    factory.NewChattRepository(),
		users:    factory.NewUserrRepository(),
		sessions: factory.NewSessionnRepository(),
	}

	closer := func() {
		if err := factory.Close(); err != nil {
			slog.Error("Закрыть соединение с pgsql: factory.Close: " + err.Error())
		}
	}

	return rs, closer, nil
}
