package response

import "x-ui/xray"

type ClientTrafficResponse struct {
	Id         int    `json:"id"`
	Enable     bool   `json:"enable"`
	Email      string `json:"email"`
	Up         int64  `json:"up"`
	Down       int64  `json:"down"`
	ExpiryTime int64  `json:"expiry_time"`
	Total      int64  `json:"total"`
}

func ClientTrafficResponseFromClientTraffic(traffic *xray.ClientTraffic) *ClientTrafficResponse {
	return &ClientTrafficResponse{
		Id:         traffic.Id,
		Enable:     traffic.Enable,
		Email:      traffic.Email,
		Up:         traffic.Up,
		Down:       traffic.Down,
		ExpiryTime: traffic.ExpiryTime,
		Total:      traffic.Total,
	}
}
