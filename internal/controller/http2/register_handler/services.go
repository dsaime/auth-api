package register_handler

import (
	"github.com/dsaime/auth-api/internal/service"
)

// services определяет интерфейс для доступа к сервисам приложения
type services interface {
	Auth() *service.Auth // Сервис аутентификации
}
