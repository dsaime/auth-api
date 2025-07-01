package register_handler

import (
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
	router.HandleFunc(
		"POST /auth/logout",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			input := service.AuthLogoutIn{
				TokenID: context.AccessTokenID(),
			}

			return context.Services().Auth().Logout(input)
		})
}
