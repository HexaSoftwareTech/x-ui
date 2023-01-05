package request_body

type StoreInboundRequestBody struct {
	Total          int64  `json:"total" binding:"number"`
	Remark         string `json:"remark"`
	Enable         bool   `json:"enable"`
	ExpiryTime     int64  `json:"expiry_time"`
	Listen         string `json:"listen"`
	Port           int    `json:"port" binding:"required,gt=0,lt=65536"`
	Protocol       string `json:"protocol" binding:"required,oneof=vmess vless trojan shadowsocks socks http dokodemo-door"`
	Settings       string `json:"settings" binding:"base64"`
	StreamSettings string `json:"stream_settings" binding:"base64"`
	Sniffing       string `json:"sniffing" binding:"base64"`
}
