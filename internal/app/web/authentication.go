package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var client oauth2.Config
var provider *oidc.Provider

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func initOIDC() error {
	var err error
	provider, err = oidc.NewProvider(context.Background(), os.Getenv("OPENID_PROVIDER_URL"))
	if err != nil {
		return err
	}

	client = oauth2.Config{
		ClientID:     os.Getenv("OPENID_CLIENT_ID"),
		ClientSecret: os.Getenv("OPENID_CLIENT_SECRET"),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  os.Getenv("OPENID_REDIRECT_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}
	return nil
}

func generateLoginURL(c *gin.Context) (string, error) {
	state, err := randString(16)
	if err != nil {
		fmt.Println("Failed to generate state")
		return "", err
	}
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("oauth_state", state, int(time.Hour.Seconds()), "/", os.Getenv("COOKIE_DOMAIN"), true, true)
	return client.AuthCodeURL(state), nil
}

func requireLogin(next func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		isAdmin := session.Get("is_admin")
		if isAdmin == "true" {
			next(c)
			return
		}
		session.Clear()
		session.Save()
		loginURL, err := generateLoginURL(c)
		if err != nil {
			fmt.Println("Failed to get login URL")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, loginURL)
		return
	}
}

func checkLogin(c *gin.Context) {
	c.Status(http.StatusOK)
}

func authenticate(c *gin.Context) {
	state, err := c.Cookie("oauth_state")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "state not found")
		return
	}
	if c.Query("state") != state {
		c.AbortWithStatusJSON(http.StatusBadRequest, "state mismatch")
		return
	}

	oauth2Token, err := client.Exchange(c, c.Query("code"))
	if err != nil {
		fmt.Println("Failed to exchange token" + err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userInfo, err := provider.UserInfo(c, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		fmt.Println("Failed to get userinfo" + err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fmt.Println("User " + userInfo.Subject + " authenticated")

	session := sessions.Default(c)
	session.Set("is_admin", "true")
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
