package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"net/http"
	"x-ui/database/model"
	requestBody "x-ui/web/request_body"
	"x-ui/web/response"
	"x-ui/web/service/mirror"
)

type MirrorController struct {
	mirrorIndexService mirror.IndexService
	storeMirrorService mirror.StoreService
}

func NewStoreMirrorController(g *gin.RouterGroup) *MirrorController {
	a := &MirrorController{}
	a.initRouter(g)
	return a
}

func (a *MirrorController) initRouter(g *gin.RouterGroup) {
	g.GET("/mirrors", a.getAllMirrors)
	g.POST("/mirrors", a.storeMirror)

}

func (a *MirrorController) getAllMirrors(c *gin.Context) {
	mirrors, err := a.mirrorIndexService.GetMirrors()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.ListResponse[response.MirrorResponse]{
		Data: lo.Map(mirrors, func(item *model.Mirror, _ int) *response.MirrorResponse {
			return response.MirrorResponseFromMirrorModel(item)
		}),
	})
}

func (a *MirrorController) storeMirror(c *gin.Context) {
	var body requestBody.StoreMirrorRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	m := model.Mirror{
		Ip:   body.Ip,
		Port: body.Port,
	}
	err := a.storeMirrorService.SaveMirror(&m)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.MirrorResponseFromMirrorModel(&m))
}
