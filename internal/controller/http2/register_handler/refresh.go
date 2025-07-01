package register_handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/dsaime/auth-api/internal/service"
)

// Refresh регистрирует обработчик, на обновление пары токенов.
//
// Метод: POST /auth/refresh
func Refresh(router *fiber.App, ss Services) {
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

			input := service.AuthRefreshIn{
				AccessToken:  rb.AccessToken,
				RefreshToken: rb.RefreshToken,
			}

			out, err := ss.Auth().Refresh(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		})
}
