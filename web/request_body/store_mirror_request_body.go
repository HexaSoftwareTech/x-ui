package request_body

type StoreMirrorRequestBody struct {
	Ip   string `json:"ip" binding:"required,ipv4"`
	Port int    `json:"port" binding:"required,gt=0,lt=65536"`
}
