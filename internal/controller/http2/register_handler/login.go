package register_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/domain"
	"github.com/dsaime/auth-api/internal/service"
)

// Login регистрирует обработчик, на получение пары токенов (access и refresh)
// для пользователя с идентификатором (GUID) указанным в параметре запроса.
//
// Метод: POST /auth/login
func Login(router *fiber.App, ss services, jwtSecret []byte) {
	// Тело запроса
	type requestBody struct {
		UserID uuid.UUID `json:"user_id"`
	}
	router.Post(
		"/auth/login",
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := service.AuthLoginIn{
				UserID:    rb.UserID,
				UserAgent: context.Get(fiber.HeaderUserAgent),
			}

			out, err := ss.Auth().Login(input)
			if err != nil {
				return err
			}

			token, err := newAccessToken(out.Session, jwtSecret)
			if err != nil {
				return context.SendStatus(fiber.StatusInternalServerError)
			}

			return context.JSON(fiber.Map{
				"session":       out.Session,
				"access_token":  token,
				"refresh_token": out.RefreshToken,
			})
		})
}

func newAccessToken(session domain.Session, jwtSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		jwtSessionIDKey: session.ID,
		jwtUserIDKey:    session.UserID,
		jwtExpKey:       session.Expiry,
	})

	return token.SignedString(jwtSecret)
}
