package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"x-ui/web/response"
	"x-ui/web/session"
)

type BaseAPIController struct {
}

func (a *BaseAPIController) checkLogin(c *gin.Context) {
	if !session.IsLogin(c) {
		if isAjax(c) {
			pureJsonMsg(c, false, I18n(c, "pages.login.loginAgain"))
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				ErrorMessage: I18n(c, "api.auth.unauthorizedMessage"),
			})
			return
		}
	} else {
		c.Next()
	}
}
