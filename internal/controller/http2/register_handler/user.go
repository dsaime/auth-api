package register_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/controller/http2/middleware"
)

// User регистрирует обработчик, на получение текущего пользователя (роут защищен).
//
// Метод: GET /user
func User(router *fiber.App, jwtSecret []byte) {
	router.Get(
		"/user",
		func(context *fiber.Ctx) error {
			token := context.Locals("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			userID, _ := claims[jwtUserIDKey].(uuid.UUID)

			// todo: ss.Users().Info({userID})

			return context.JSON(fiber.Map{
				"id": userID,
			})
		},
		logger.New(),
		middleware.ProtectedWithJWT(jwtSecret),
	)
}
