package register_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/controller/http2/middleware"
)

// User регистрирует обработчик, на получение текущего пользователя (роут защищен).
//
// Метод: GET /user
func User(router *fiber.App, jwtSecret string) {
	router.Get(
		"/user",
		func(context *fiber.Ctx) error {
			user := context.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			userID, _ := claims["user_id"].(uuid.UUID)

			// todo: ss.Users().Info()

			return context.JSON(fiber.Map{
				"id": userID,
			})
		},
		middleware.ProtectedWithJWT(jwtSecret),
	)
}
