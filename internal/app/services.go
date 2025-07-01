package app

import "github.com/dsaime/auth-api/internal/service"

type services struct {
	chats    *service.Chats
	sessions *service.Auth
	users    *service.Users
}

func (s *services) Chats() *service.Chats {
	return s.chats
}

func (s *services) Auth() *service.Auth {
	return s.sessions
}

func (s *services) Users() *service.Users {
	return s.users
}

func initServices(repos *repositories, adaps *adapters) *services {
	return &services{
		chats: &service.Chats{
			Repo: repos.chats,
		},
		sessions: &service.Auth{
			Repo: repos.sessions,
		},
		users: &service.Users{
			Providers:    adaps.OAuthProviders(),
			Repo:         repos.users,
			SessionsRepo: repos.sessions,
		},
	}
}
