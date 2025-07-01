package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/dsaime/auth-api/internal/controller/http2"
	registerHandler "github.com/dsaime/auth-api/internal/controller/http2/register_handler"
	"github.com/dsaime/auth-api/internal/controller/http2/router"
)

func initHttpServer(ss *services, cfg Config) *http.Server {
	r := &router.Router{
		Services: ss,
	}
	registerHandlers(r)

	return &http.Server{
		Addr:    cfg.HttpAddr,
		Handler: r,
	}
}

func runHttpServer(ctx context.Context, server *http.Server) error {
	g, ctx := errgroup.WithContext(ctx)

	// Запуск сервера
	g.Go(func() error {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server.ListenAndServe: %w", err)
		}
		return nil
	})

	// Завершение сервера при завершении контекста
	g.Go(func() error {
		// > The first call to return a non-nil error cancels the group's context
		<-ctx.Done()
		return server.Shutdown(ctx)
	})

	return g.Wait()
}

func registerHandlers(r http2.Router) {
	// Служебные
	registerHandler.Ping(r)

	// Аутентификация /auth
	registerHandler.Login(r)
	registerHandler.Refresh(r)
	registerHandler.Logout(r)

	// Аутентифицированный пользователь /user
	registerHandler.User(r)
}
