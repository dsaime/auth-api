package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"

	registerHandler "github.com/dsaime/auth-api/internal/controller/http2/register_handler"
)

func runHttpServer(ctx context.Context, ss *services, cfg Config) error {
	fiberApp := fiber.New()
	registerHandlers(fiberApp, ss, cfg.JWTSecret)

	g, ctx := errgroup.WithContext(ctx)

	// Запуск сервера
	g.Go(func() error {
		err := fiberApp.Listen(cfg.HttpAddr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server.ListenAndServe: %w", err)
		}
		return nil
	})

	// Завершение сервера при завершении контекста
	g.Go(func() error {
		<-ctx.Done()
		return fiberApp.Shutdown()
	})

	return g.Wait()
}

func registerHandlers(r *fiber.App, ss *services, jwtSecret string) {
	// Служебные
	registerHandler.Ping(r)

	// Аутентификация /auth
	registerHandler.Login(r, ss, jwtSecret)
	registerHandler.Refresh(r, ss)
	registerHandler.Logout(r, ss, jwtSecret)

	// Аутентифицированный пользователь /user
	registerHandler.User(r, jwtSecret)
}
