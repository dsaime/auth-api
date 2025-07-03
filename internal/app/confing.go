package app

import (
	"github.com/dsaime/auth-api/internal/repository/pgsql"
)

type Config struct {
	Pgsql     pgsql.Config
	HttpAddr  string
	LogLevel  string
	JWTSecret string
	Webhook   string
}
