package register_handler

import (
	"github.com/google/uuid"

	"github.com/dsaime/auth-api/internal/controller/http2"
	"github.com/dsaime/auth-api/internal/controller/http2/middleware"
)

// User регистрирует обработчик, на получение текущего пользователя (роут защищен).
//
// Метод: GET /user
func User(router http2.Router) {
	// Тело запроса
	type requestBody struct {
		UserID uuid.UUID `json:"user_id"`
	}
	router.HandleFunc(
		"GET /user", // Подсмотрел тут https://docs.github.com/en/rest/users/users
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			return map[string]any{
				"id": rb.UserID,
			}, nil

			//input := service.UserInfo{
			//	UserID: rb.UserID,
			//}
			//
			//return context.Services().Users().Info(input)
		})
}
