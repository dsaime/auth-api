package app

import "github.com/dsaime/auth-api/internal/service"

type services struct {
	auth *service.Auth
}

func (s *services) Auth() *service.Auth {
	return s.auth
}

func initServices(repos *repositories, adaps *adapters) *services {
	return &services{
		auth: &service.Auth{
			Repo: repos.sessions,
		},
	}
}
