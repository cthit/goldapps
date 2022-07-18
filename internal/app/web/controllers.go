package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cthit/goldapps/internal/pkg/actions"
	"github.com/gin-gonic/gin"
)

type ChangeBody struct {
	UserChanges  actions.UserActions  `json:"userChanges"`
	GroupChanges actions.GroupActions `json:"groupChanges"`
}

func getSuggestions(c *gin.Context) {
	user, group, err := getChangeSuggestions("", "gapps.json")
	var code int
	if err != nil {
		code = http.StatusBadRequest
		fmt.Println(err)
	} else {
		code = http.StatusOK
	}
	response := ChangeBody{user, group}
	c.JSON(code, response)
}

func executeChanges(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	changes := ChangeBody{}
	err = json.Unmarshal(body, &changes)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ok := commitChanges(changes.UserChanges, changes.GroupChanges, "gapps.json")
	if !ok {
		fmt.Println("Failed to execute all changes without errors")
		c.AbortWithError(http.StatusBadRequest, errors.New("Failed to execute all changes without errors"))
		return
	}
	c.Status(http.StatusCreated)
	c.Done()
}
