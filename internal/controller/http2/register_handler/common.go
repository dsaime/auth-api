package register_handler

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtUserIDKey    = "user_id"
	jwtExpKey       = "exp"
	jwtSessionIDKey = "jti"
)

func validateJWT(tokenString string, jwtSecret []byte) (*jwt.Token, error) {
	return jwt.NewParser().Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtSecret, nil
		},
	)
}
