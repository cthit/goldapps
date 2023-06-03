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

type ChangeRequestBody struct {
	ChangeBody
	To string `json:"to"`
}

func getSuggestions(c *gin.Context) {
	from, ok := c.GetQuery("from")
	if !ok {
		fmt.Println("Parameter 'from' was not provided")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	to, ok := c.GetQuery("to")
	if !ok {
		fmt.Println("Parameter 'to' was not provided")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, group, err := getChangeSuggestions(from, to)
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

	changes := ChangeRequestBody{}
	err = json.Unmarshal(body, &changes)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if changes.To == "" {
		fmt.Println("'to' was not provided")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ok := commitChanges(changes.UserChanges, changes.GroupChanges, changes.To)
	if !ok {
		fmt.Println("Failed to execute all changes without errors")
		c.AbortWithError(http.StatusBadRequest, errors.New("Failed to execute all changes without errors"))
		return
	}
	c.Status(http.StatusCreated)
	c.Done()
}
