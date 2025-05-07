package driver

type RespCheck struct {
	DefaultNumbers []string `json:"default_numbers"`
}

type RespLoginWa struct {
	Session    bool   `json:"session"`
	QrLink     string `json:"qr_link"`
	QrDuration int    `json:"qr_duration"`
}

type RespSessionWa struct {
	Session bool `json:"session"`
}
