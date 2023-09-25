package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xV0lk/hotel-reservations/db"
)

func JWTAuth(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			return ErrUnauthorized()
		}
		claims, err := ValidateToken(token)
		if err != nil {
			return ErrUnauthorized()
		}
		// check token expiration
		expiration := claims["expiration"].(float64)
		if int64(expiration) < time.Now().Unix() {
			return NewError(http.StatusUnauthorized, "Token expired")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return ErrUnauthorized()
		}
		user.Password = ""
		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method:", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		fmt.Println(err)
		return nil, fmt.Errorf("Unauthorized")
	}
}
