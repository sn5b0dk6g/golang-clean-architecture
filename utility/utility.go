package utility

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	//cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode

	return cookie
}