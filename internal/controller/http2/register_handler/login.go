package register_handler

import (
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/controller/http2"
	"github.com/dsaime/auth-api/internal/controller/http2/middleware"
	"github.com/dsaime/auth-api/internal/service"
)

// Login регистрирует обработчик, на получение пары токенов (access и refresh)
// для пользователя с идентификатором (GUID) указанным в параметре запроса.
//
// Метод: POST /auth/login
func Login(router http2.Router) {
	// Тело запроса
	type requestBody struct {
		UserID uuid.UUID `json:"user_id"`
	}
	router.HandleFunc(
		"POST /auth/login",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.AuthLoginIn{
				UserID:    rb.UserID,
				UserAgent: context.UserAgent(),
			}

			return context.Services().Auth().Login(input)
		})
}
