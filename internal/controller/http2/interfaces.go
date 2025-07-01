package http2

import (
	"net/http"

	"github.com/dsaime/auth-api/internal/service"
)

// Router определяет интерфейс для маршрутизации HTTP-запросов
type Router interface {
	// HandleFunc регистрирует обработчик для указанного пути и цепочки middleware.
	HandleFunc(pattern string, chain []Middleware, handler HandlerFunc)
}

// HandlerFunc представляет собой функцию-обработчик
type HandlerFunc func(Context) (any, error)

// Context определяет интерфейс для доступа к информации о запросе.
type Context interface {
	// Request возвращает HTTP-запрос
	Request() *http.Request

	// Writer возвращает интерфейс для записи ответа
	Writer() http.ResponseWriter

	// Services возвращает доступ к сервисам приложения
	Services() Services
}

// HandlerFuncRW представляет собой функцию-обработчик, но с RWContext
type HandlerFuncRW func(RWContext) (any, error)

// RWContext определяет интерфейс для доступа к контексту с возможностью изменения.
type RWContext interface {
	// Context для расширения существующего контекста
	Context

	// SetRequest устанавливает HTTP-запрос
	SetRequest(*http.Request)
}

// Middleware представляет интерфейс для middleware-функций
type Middleware func(rw HandlerFuncRW) HandlerFuncRW

// Services определяет интерфейс для доступа к сервисам приложения
type Services interface {
	Auth() *service.Auth // Сервис аутентификации
}
