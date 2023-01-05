package response

import "x-ui/database/model"

type MirrorResponse struct {
	Id   int    `json:"id"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

func MirrorResponseFromMirrorModel(mirror *model.Mirror) *MirrorResponse {
	return &MirrorResponse{
		Id:   mirror.Id,
		Ip:   mirror.Ip,
		Port: mirror.Port,
	}
}
