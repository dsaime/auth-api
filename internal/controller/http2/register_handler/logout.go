package register_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/controller/http2/middleware"
	"github.com/dsaime/auth-api/internal/service"
)

// Logout регистрирует обработчик, на деавторизацию пользователя (поле выполнения
// этого запроса с access токеном, пользователю больше не должен быть доступен роут
// на получение его GUID и операция обновления токенов).
//
// Метод: POST /auth/logout
func Logout(router *fiber.App, ss services, jwtSecret []byte) {
	router.Post(
		"/auth/logout",
		func(context *fiber.Ctx) error {
			user := context.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			sessionID, _ := claims[jwtSessionIDKey].(uuid.UUID)

			input := service.AuthLogoutIn{
				SessionID: sessionID,
			}

			return ss.Auth().Logout(input)
		},
		logger.New(),
		middleware.ProtectedWithJWT(jwtSecret),
	)
}
