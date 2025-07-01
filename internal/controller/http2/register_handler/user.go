package register_handler

import (
	"github.com/dsaime/auth-api/internal/controller/http2"
	"github.com/dsaime/auth-api/internal/controller/http2/middleware"
)

// User регистрирует обработчик, на получение текущего пользователя (роут защищен).
//
// Метод: GET /user
func User(router http2.Router) {
	router.HandleFunc(
		"GET /user", // Подсмотрел тут https://docs.github.com/en/rest/users/users
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return map[string]any{
				"id": context.UserID(),
			}, nil

			//input := service.UserInfo{
			//	UserID: rb.UserID,
			//}
			//
			//return context.Services().Users().Info(input)
		})
}
