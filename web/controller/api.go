package controller

import (
	"github.com/gin-gonic/gin"
)

type APIController struct {
	BaseAPIController

	inboundController *InboundController
	mirrorController  *MirrorController
}

func NewAPIController(g *gin.RouterGroup) *APIController {
	a := &APIController{}
	a.initRouter(g)
	return a
}

func (a *APIController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/xui/api")
	g.Use(a.checkLogin)

	a.inboundController = NewInboundController(g)
	a.mirrorController = NewStoreMirrorController(g)
}
