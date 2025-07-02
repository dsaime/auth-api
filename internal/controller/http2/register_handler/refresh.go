package register_handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/service"
)

// Refresh регистрирует обработчик, на обновление пары токенов.
//
// Метод: POST /auth/refresh
func Refresh(router *fiber.App, ss services, jwtSecret string) {
	// Тело запроса
	type requestBody struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}
	router.Post(
		"/auth/refresh",
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			token, err := validateJWT(rb.AccessToken, []byte(jwtSecret))
			if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
				return err
			}
			claims := token.Claims.(jwt.MapClaims)
			sessionID, _ := claims[jwtSessionIDKey].(uuid.UUID)

			input := service.AuthRefreshIn{
				SessionID:    sessionID,
				RefreshToken: rb.RefreshToken,
			}

			out, err := ss.Auth().Refresh(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		})
}
