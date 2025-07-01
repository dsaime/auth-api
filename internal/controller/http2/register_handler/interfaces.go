package register_handler

import (
	"github.com/dsaime/auth-api/internal/service"
)

// Services определяет интерфейс для доступа к сервисам приложения
type Services interface {
	Auth() *service.Auth // Сервис аутентификации
}
