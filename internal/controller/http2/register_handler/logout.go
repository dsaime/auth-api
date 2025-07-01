package register_handler

import (
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/controller/http2"
	"github.com/dsaime/auth-api/internal/controller/http2/middleware"
	"github.com/dsaime/auth-api/internal/service"
)

// Logout регистрирует обработчик, на деавторизацию пользователя (поле выполнения
// этого запроса с access токеном, пользователю больше не должен быть доступен роут
// на получение его GUID и операция обновления токенов).
//
// Метод: POST /auth/logout
func Logout(router http2.Router) {
	// Тело запроса
	type requestBody struct {
		UserID uuid.UUID `json:"user_id"`
	}
	router.HandleFunc(
		"POST /auth/logout",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.AuthLogout{
				UserID: rb.UserID,
			}

			return context.Services().Auth().Login(input)
		})
}
