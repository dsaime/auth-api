package register_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/dsaime/auth-api/internal/controller/http2/middleware"
	"github.com/dsaime/auth-api/internal/service"
)

// Logout регистрирует обработчик, на деавторизацию пользователя (поле выполнения
// этого запроса с access токеном, пользователю больше не должен быть доступен роут
// на получение его GUID и операция обновления токенов).
//
// Метод: POST /auth/logout
func Logout(router *fiber.App, ss Services, jwtSecret string) {
	router.Post(
		"/auth/logout",
		func(context *fiber.Ctx) error {
			user := context.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			tokenID, _ := claims["jti"].(string)

			input := service.AuthLogoutIn{
				TokenID: tokenID,
			}

			out, err := ss.Auth().Logout(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		middleware.ProtectedWithJWT(jwtSecret),
	)
}
