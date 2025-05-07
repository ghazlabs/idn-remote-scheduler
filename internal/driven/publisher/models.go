package wa

import "encoding/json"

type RespSendMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (r *RespSendMessage) IsSessionExpired() bool {
	if r.Code == "AUTHENTICATION_ERROR" || r.Code == "SESSION_SAVED_ERROR" || r.Code == "INTERNAL_SERVER_ERROR" {
		return true
	}

	return false
}

func (r RespSendMessage) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}

type RespGetLoginQrCode struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Results struct {
		QrDuration int    `json:"qr_duration"`
		QrLink     string `json:"qr_link"`
	} `json:"results"`
}

func (r *RespGetLoginQrCode) IsSessionAlreadyActive() bool {
	return r.Code == "ALREADY_LOGGED_IN"
}

func (r RespGetLoginQrCode) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}

type RespGetSession struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Results []struct {
		Name   string `json:"name"`
		Device string `json:"device"`
	} `json:"results"`
}

func (r RespGetSession) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}
