package app

import (
	pgsqlRepository "github.com/dsaime/auth-api/internal/repository/pgsql_repository"
)

type Config struct {
	Pgsql     pgsqlRepository.Config
	HttpAddr  string
	LogLevel  string
	JWTSecret string
}
