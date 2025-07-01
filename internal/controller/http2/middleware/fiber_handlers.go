package middleware

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func ProtectedWithJWT(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	})
}

// RecoverPanic перехватывает панику
func RecoverPanic(context *fiber.Ctx) (err error) {
	defer func() {
		if pv := recover(); pv != nil {
			err = fmt.Errorf("%v", pv)
			return
		}
	}()

	return nil
}
