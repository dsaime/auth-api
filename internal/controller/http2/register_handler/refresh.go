package register_handler

import (
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/controller/http2"
	"github.com/dsaime/auth-api/internal/controller/http2/middleware"
	"github.com/dsaime/auth-api/internal/service"
)

// Refresh регистрирует обработчик, на обновление пары токенов.
//
// Метод: POST /auth/refresh
func Refresh(router http2.Router) {
	// Тело запроса
	type requestBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	router.HandleFunc(
		"POST /auth/refresh",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.AuthRefreshIn{
				TokenID:      context.JWT().ID,
				RefreshToken: rb.RefreshToken,
			}

			return context.Services().Auth().Refresh(input)
		})
}
