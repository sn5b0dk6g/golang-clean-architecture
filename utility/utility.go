package utility

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetRedisExpiration() time.Duration {
	expiration, err := strconv.Atoi(os.Getenv("REDIS_EXPIRATION"))
	if err != nil {
		log.Fatalln(err)
	}
	return time.Duration(expiration) * time.Second
}

func CreateCookie(value string, exp time.Time) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = value
	cookie.Expires = exp
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	if os.Getenv("GO_ENV") != "local" {
		cookie.Secure = true
	}
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode

	return cookie
}

func GetUserIdByToken(c echo.Context) string {
	return c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(string)
}
