package web

import (
	"fmt"

	"github.com/cthit/goldapps/internal/app/cli"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	err := cli.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config")
		fmt.Println(err)
		return
	}

	r.GET("/suggestions", getSuggestions)
	r.POST("/commit", executeChanges)
	r.Run()
}
