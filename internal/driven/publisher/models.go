package wa

type RespSendMessage struct {
	Code string `json:"code"`
}

func (r *RespSendMessage) IsSessionExpired() bool {
	if r.Code == "AUTHENTICATION_ERROR" || r.Code == "SESSION_SAVED_ERROR" || r.Code == "INTERNAL_SERVER_ERROR" {
		return true
	}

	return false
}
