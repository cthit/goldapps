package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var gamma_url = viper.GetString("gamma.provider.url")

var client = oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("AUTH_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:   fmt.Sprintf("%s/api/oauth/authorize", os.Getenv("REDIRECT_GAMMA_URL")),
		TokenURL:  fmt.Sprintf("%s/api/oauth/token", gamma_url),
		AuthStyle: oauth2.AuthStyleAutoDetect,
	},
	RedirectURL: os.Getenv("CALLBACK_URL"),
	Scopes:      nil,
}

func getLoginURL() string {
	return fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s",
		client.Endpoint.AuthURL,
		client.ClientID,
		client.RedirectURL)
}

func requireLogin(next func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		is_admin := session.Get("is_admin")
		if is_admin == "true" {
			next(c)
			return
		}
		session.Clear()
		session.Save()
		c.AbortWithStatusJSON(http.StatusUnauthorized, getLoginURL())
		return
	}
}

func checkLogin(c *gin.Context) {
	c.Status(http.StatusOK)
}

type User struct {
	Authorities []struct {
		Id        string `json:"id"`
		Authority string `json:"authority"`
	} `json:"authorities"`
}

func authenticate(c *gin.Context) {
	code := c.Query("code")
	token, err := client.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Failed to authenticate user")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	gammaQuery := fmt.Sprintf("%s/api/users/me", gamma_url)
	resp, err := client.Client(context.Background(), token).Get(gammaQuery)
	if err != nil {
		fmt.Println("Failed to get user from gamma")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user := User{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &user)

	found := false
	for _, auth := range user.Authorities {
		if auth.Authority == "admin" {
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusUnauthorized, getLoginURL())
		return
	}

	session := sessions.Default(c)
	session.Set("is_admin", "true")
	session.Save()
	c.Status(http.StatusOK)
}
