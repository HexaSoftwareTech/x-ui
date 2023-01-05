package response

import (
	"github.com/samber/lo"
	"x-ui/database/model"
	"x-ui/xray"
)

type ProtocolResponse string

const (
	VMess       ProtocolResponse = "vmess"
	VLESS       ProtocolResponse = "vless"
	Dokodemo    ProtocolResponse = "Dokodemo-door"
	Http        ProtocolResponse = "http"
	Trojan      ProtocolResponse = "trojan"
	Shadowsocks ProtocolResponse = "shadowsocks"
)

type InboundResponse struct {
	Id             int                      `json:"id"`
	UserId         int                      `json:"user_Id"`
	Up             int64                    `json:"up"`
	Down           int64                    `json:"down"`
	Total          int64                    `json:"total"`
	Remark         string                   `json:"remark"`
	Enable         bool                     `json:"enable"`
	ExpiryTime     int64                    `json:"expiry_time"`
	Listen         string                   `json:"listen"`
	Port           int                      `json:"port"`
	Protocol       ProtocolResponse         `json:"protocol"`
	Settings       string                   `json:"settings"`
	StreamSettings string                   `json:"stream_settings"`
	Tag            string                   `json:"tag"`
	Sniffing       string                   `json:"sniffing"`
	ClientTraffics []*ClientTrafficResponse `json:"client_statistics"`
}

func InboundResponseFromInbound(inbound *model.Inbound) *InboundResponse {
	clientStatistics := lo.Map[xray.ClientTraffic, *ClientTrafficResponse](inbound.ClientStats, func(item xray.ClientTraffic, _ int) *ClientTrafficResponse {
		return ClientTrafficResponseFromClientTraffic(&item)
	})
	return &InboundResponse{
		Id:             inbound.Id,
		UserId:         inbound.UserId,
		Up:             inbound.Up,
		Down:           inbound.Down,
		Total:          inbound.Total,
		Remark:         inbound.Remark,
		Enable:         inbound.Enable,
		ExpiryTime:     inbound.ExpiryTime,
		Listen:         inbound.Listen,
		Port:           inbound.Port,
		Protocol:       protocolResponseFromProtocol(inbound.Protocol),
		Settings:       inbound.Settings,
		StreamSettings: inbound.StreamSettings,
		Tag:            inbound.Tag,
		Sniffing:       inbound.Sniffing,
		ClientTraffics: clientStatistics,
	}
}

func protocolResponseFromProtocol(protocol model.Protocol) ProtocolResponse {
	var protocolResponse ProtocolResponse
	switch protocol {
	case model.VMess:
		protocolResponse = VMess
		break
	case model.VLESS:
		protocolResponse = VLESS
		break
	case model.Dokodemo:
		protocolResponse = Dokodemo
		break
	case model.Http:
		protocolResponse = Http
		break
	case model.Trojan:
		protocolResponse = Trojan
		break
	case model.Shadowsocks:
		protocolResponse = Shadowsocks
		break
	}

	return protocolResponse
}
