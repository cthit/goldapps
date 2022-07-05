package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getSuggestions(c *gin.Context) {
	user, group, err := getChangeSuggestions("", "gapps.json")
	var code int
	if err != nil {
		code = http.StatusBadRequest
		fmt.Println(err)
	} else {
		code = http.StatusOK
	}
	c.JSON(code, gin.H{
		"userChanges":  user,
		"groupChanges": group,
	})
}
