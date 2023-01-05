package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/web/global"
	requestBody "x-ui/web/request_body"
	"x-ui/web/response"
	"x-ui/web/service"
	"x-ui/web/session"
)

type InboundController struct {
	inboundService service.InboundService
	xrayService    service.XrayService
}

func NewInboundController(g *gin.RouterGroup) *InboundController {
	a := &InboundController{}
	a.initRouter(g)
	a.startTask()
	return a
}

func (a *InboundController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/inbounds")

	g.GET(":id", a.getInbound)
	g.GET("", a.getInbounds)
	g.POST("", a.addInbound)
	g.DELETE(":id", a.deleteInbound)
	g.PUT(":id", a.updateInbound)

	g.POST("/clientIps/:email", a.getClientIps)
	g.POST("/clearClientIps/:email", a.clearClientIps)
	g.POST("/resetClientTraffic/:email", a.resetClientTraffic)
}

func (a *InboundController) startTask() {
	webServer := global.GetWebServer()
	c := webServer.GetCron()
	c.AddFunc("@every 10s", func() {
		if a.xrayService.IsNeedRestartAndSetFalse() {
			err := a.xrayService.RestartXray(false)
			if err != nil {
				logger.Error("restart xray failed:", err)
			}
		}
	})
}

func (a *InboundController) getInbounds(c *gin.Context) {
	user := session.GetLoginUser(c)
	inbounds, err := a.inboundService.GetInbounds(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
	}

	c.JSON(http.StatusOK, response.ListResponse[response.InboundResponse]{
		Data: lo.Map[*model.Inbound, *response.InboundResponse](inbounds, func(item *model.Inbound, _ int) *response.InboundResponse {
			return response.InboundResponseFromInbound(item)
		}),
	})
}

func (a *InboundController) getInbound(c *gin.Context) {
	user := session.GetLoginUser(c)
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	inbound, err := a.inboundService.GetUserInbound(user.Id, id)

	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorMessage: "Inbound not found",
		})
	}

	c.JSON(http.StatusOK, response.InboundResponseFromInbound(inbound))
}

func (a *InboundController) addInbound(c *gin.Context) {
	body := &requestBody.StoreInboundRequestBody{}
	err := c.ShouldBindJSON(body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	inbound, err := inboundFromStoreInboundRequestBody(body)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	user := session.GetLoginUser(c)
	inbound.UserId = user.Id
	inbound.Enable = true
	inbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)
	inbound, err = a.inboundService.AddInbound(inbound)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.InboundResponseFromInbound(inbound))
	a.xrayService.SetToNeedRestart()
}

func (a *InboundController) deleteInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	err = a.inboundService.DelInbound(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		SuccessMessage: I18n(c, "api.inbound.removeInboundSuccessMessage"),
	})
	a.xrayService.SetToNeedRestart()
}

func (a *InboundController) updateInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	body := &requestBody.UpdateInboundRequestBody{}
	err = c.ShouldBindJSON(body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	inbound, err := inboundFromUpdateInboundRequestBody(body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}
	inbound.Id = id

	inbound, err = a.inboundService.UpdateInbound(inbound)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, response.ErrorResponse{
				ErrorMessage: err.Error(),
			})
			return
		}

		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.InboundResponseFromInbound(inbound))
	a.xrayService.SetToNeedRestart()
}

func (a *InboundController) getClientIps(c *gin.Context) {
	email := c.Param("email")

	ips, err := a.inboundService.GetInboundClientIps(email)
	if err != nil {
		jsonObj(c, "No IP Record", nil)
		return
	}
	jsonObj(c, ips, nil)
}

func (a *InboundController) clearClientIps(c *gin.Context) {
	email := c.Param("email")

	err := a.inboundService.ClearClientIps(email)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse{SuccessMessage: "Log Cleared"})
}

func inboundFromStoreInboundRequestBody(body *requestBody.StoreInboundRequestBody) (*model.Inbound, error) {
	result := model.Inbound{
		Total:      body.Total,
		Remark:     body.Remark,
		Enable:     body.Enable,
		ExpiryTime: body.ExpiryTime,
		Listen:     body.Listen,
		Port:       body.Port,
		Protocol:   model.Protocol(body.Protocol),
	}

	decodedSettings, err := base64.StdEncoding.DecodeString(body.Settings)
	if err != nil {
		return nil, err
	}
	result.Settings = string(decodedSettings)

	decodedStreamSettings, err := base64.StdEncoding.DecodeString(body.StreamSettings)
	if err != nil {
		return nil, err
	}
	result.StreamSettings = string(decodedStreamSettings)

	decodedSniffing, err := base64.StdEncoding.DecodeString(body.Sniffing)

	if err != nil {
		return nil, err
	}

	result.Sniffing = string(decodedSniffing)

	return &result, nil
}

func inboundFromUpdateInboundRequestBody(body *requestBody.UpdateInboundRequestBody) (*model.Inbound, error) {
	result := model.Inbound{
		Total:      body.Total,
		Remark:     body.Remark,
		Enable:     body.Enable,
		ExpiryTime: body.ExpiryTime,
		Listen:     body.Listen,
		Port:       body.Port,
		Protocol:   model.Protocol(body.Protocol),
	}

	decodedSettings, err := base64.StdEncoding.DecodeString(body.Settings)
	if err != nil {
		return nil, err
	}
	result.Settings = string(decodedSettings)

	decodedStreamSettings, err := base64.StdEncoding.DecodeString(body.StreamSettings)
	if err != nil {
		return nil, err
	}
	result.StreamSettings = string(decodedStreamSettings)

	decodedSniffing, err := base64.StdEncoding.DecodeString(body.Sniffing)

	if err != nil {
		return nil, err
	}

	result.Sniffing = string(decodedSniffing)

	return &result, nil
}

func (a *InboundController) resetClientTraffic(c *gin.Context) {
	email := c.Param("email")

	err := a.inboundService.ResetClientTraffic(email)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse{SuccessMessage: "Traffic has been reset"})
}
