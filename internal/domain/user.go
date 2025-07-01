package domain

import (
	"github.com/google/uuid"
)

// User представляет собой агрегат пользователя.
type User struct {
	ID uuid.UUID // ID пользователя
}

// NewUser создает нового пользователя.
func NewUser() (User, error) {
	// В параметры можно добавить вручную задаваемые поля.
	// В тело можно добавить валидацию параметров.
	return User{
		ID: uuid.New(),
	}, nil
}

// UserRepository представляет собой интерфейс для работы с репозиторием пользователей.
type UserRepository interface {
	List(UserFilter) ([]User, error)
	Upsert(User) error
	InTransaction(func(txRepo UserRepository) error) error
}

// UserFilter представляет собой фильтр для выборки пользователей.
type UserFilter struct {
	ID                uuid.UUID // ID пользователя для фильтрации
	OAuthUserID       string    // Фильтрация по ID пользователя провайдера
	OAuthProvider     string    // Фильтрация по провайдеру
	BasicAuthLogin    string    // Логин пользователя для фильтрации
	BasicAuthPassword string    // Пароль пользователя для фильтрации
}
