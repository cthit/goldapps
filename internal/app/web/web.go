package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cthit/goldapps/internal/app/cli"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options(sessions.Options{
		MaxAge:   60 * 60 * 24,
		Path:     "/",
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode})
	r.Use(sessions.Sessions("goldapps-session", store))

	err := cli.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config")
		fmt.Println(err)
		return
	}

	r.GET("/api/checkLogin", requireLogin(checkLogin))
	r.GET("/api/authenticate", authenticate)

	r.GET("/api/suggestions", requireLogin(getSuggestions))
	r.POST("/api/commit", requireLogin(executeChanges))
	r.Run()
}
