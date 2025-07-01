package app

import (
	"github.com/dsaime/auth-api/internal/adapter/oauth_provider"
	pgsqlRepository "github.com/dsaime/auth-api/internal/repository/pgsql_repository"
)

type Config struct {
	Pgsql       pgsqlRepository.Config
	HttpAddr    string
	LogLevel    string
	OAuthGoogle oauth_provider.GoogleConfig
	OAuthGitHub oauth_provider.GitHubConfig
}
