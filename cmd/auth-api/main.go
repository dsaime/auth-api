package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/dsaime/auth-api/internal/app"
)

func main() {
	slog.Info("Запуск")
	ctx, cancel := context.WithCancel(context.Background())
	go appRun(ctx)
	waitInterrupt(cancel)
}

// appRun запускает приложение и обрабатывает результат
func appRun(ctx context.Context) {
	err := initCliCommand().Run(ctx, os.Args)
	if errors.Is(err, context.Canceled) {
		slog.Info("Завершение работы из-за отмены контекста")
		os.Exit(0)
	} else if err != nil {
		slog.Error("Завершение работы из-за ошибки: " + err.Error())
		os.Exit(1)
	}
	slog.Info("Завершение без ошибок")
	os.Exit(0)
}

// waitInterrupt отменяет контекст, когда в приложение поступает сигнал syscall.SIGINT или syscall.SIGTERM
func waitInterrupt(cancel context.CancelFunc) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("Получен сигнал от ОС: " + (<-interrupt).String())
	cancel()
	slog.Info("Контекст отменен")
	time.Sleep(3 * time.Second)
}

// initCliCommand создает команду, для разбора аргументов командной строки и запуска приложения
func initCliCommand() *cli.Command {
	var cfg app.Config
	return &cli.Command{
		Action: func(ctx context.Context, command *cli.Command) error {
			return app.Run(ctx, cfg)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "pgsql-dsn",
				Destination: &cfg.Pgsql.DSN,
				Usage:       "Строка подключения PostgreSQL в формате 'postgres://user:password@host:port/dbname'",
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "jwt-secret",
				Destination: &cfg.JWTSecret,
				Usage:       "Строка для подписи jwt",
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "ip-alert-webhook",
				Destination: &cfg.Webhook,
				Usage:       "URL на который будет выполняться POST-запрос при обновлении токена с нового ip (с подстановкой text/template переменных)",
			},
			&cli.StringFlag{
				Name:        "http-addr",
				Destination: &cfg.HttpAddr,
				Usage:       "Адрес для запуска HTTP сервера",
				Value:       ":8080",
			},
			&cli.StringFlag{
				Name:        "log-level",
				Destination: &cfg.LogLevel,
				Usage:       "Уровень логирования. Может быть debug, info, warn, error",
				Value:       app.LogLevelInfo,
			},
		},
	}
}
