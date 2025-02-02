package service

import (
	"net/http"

	"gin_blog/pkg/e"
	"gin_blog/pkg/util"
	"gin_blog/models"

	"github.com/gin-gonic/gin"
	"github.com/astaxie/beego/validation"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	atuhVa := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&atuhVa)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS

	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code" : code,
        "msg" : e.GetMsg(code),
        "data" : data,
	})
}