package wa

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ghazlabs/wa-scheduler/internal/core"
	"github.com/go-resty/resty/v2"
	"gopkg.in/validator.v2"
)

const (
	WA_PUBLISHER_DOCKER_DOMAIN    = "whatsapp"
	WA_PUBLISHER_LOCALHOST_DOMAIN = "localhost"
)

type WaPublisher struct {
	WaPublisherConfig
}

func NewWaPublisher(cfg WaPublisherConfig) (*WaPublisher, error) {
	// validate config
	err := validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &WaPublisher{
		WaPublisherConfig: cfg,
	}, nil
}

type WaPublisherConfig struct {
	HttpClient   *resty.Client `validate:"nonnil"`
	Username     string        `validate:"nonzero"`
	Password     string        `validate:"nonzero"`
	WaApiBaseUrl string        `validate:"nonzero"`
}

func (n *WaPublisher) Publish(ctx context.Context, msg core.Message) error {
	for _, recID := range msg.RecipientNumbers {
		err := n.sendMessage(ctx, recID, msg.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *WaPublisher) GetLoginQrCode(ctx context.Context) (*core.QrCodeLogin, error) {
	var rsp RespGetLoginQrCode

	resp, err := n.HttpClient.R().
		SetContext(ctx).
		SetBasicAuth(n.Username, n.Password).
		SetResult(&rsp).
		SetError(&rsp).
		Get(fmt.Sprintf("%v/app/login", n.WaApiBaseUrl))
	if err != nil {
		return nil, fmt.Errorf("unable to make http request: %w", err)
	}
	if resp.IsError() {
		slog.Error("failed to get login qr code", slog.String("response", rsp.String()))
		if rsp.IsSessionAlreadyActive() {
			return &core.QrCodeLogin{
				Session:    true,
				QrDuration: 0,
				QrLink:     "",
			}, nil
		}

		return nil, fmt.Errorf("failed to get login qr code: %s", resp.String())
	}

	qrLink := rsp.Results.QrLink
	if strings.HasPrefix(qrLink, fmt.Sprintf("http://%s:", WA_PUBLISHER_DOCKER_DOMAIN)) || strings.HasPrefix(qrLink, fmt.Sprintf("https://%s:", WA_PUBLISHER_DOCKER_DOMAIN)) {
		qrLink = strings.Replace(qrLink, WA_PUBLISHER_DOCKER_DOMAIN, WA_PUBLISHER_LOCALHOST_DOMAIN, 1)
	}

	return &core.QrCodeLogin{
		Session:    false,
		QrDuration: rsp.Results.QrDuration,
		QrLink:     qrLink,
	}, nil
}

func (n *WaPublisher) GetSession(ctx context.Context) (bool, error) {
	var rsp RespGetSession

	resp, err := n.HttpClient.R().
		SetContext(ctx).
		SetBasicAuth(n.Username, n.Password).
		SetResult(&rsp).
		SetError(&rsp).
		Get(fmt.Sprintf("%v/app/devices", n.WaApiBaseUrl))
	if err != nil {
		return false, fmt.Errorf("unable to make http request: %w", err)
	}
	if resp.IsError() {
		return false, fmt.Errorf("failed to get whatsapp session: %s", resp.String())
	}

	devices := rsp.Results
	session := false
	if len(devices) > 0 {
		session = true
	}

	return session, nil
}

func (n *WaPublisher) sendMessage(ctx context.Context, recID string, content string) error {
	// send notification to whatsapp
	var rsp RespSendMessage
	resp, err := n.HttpClient.R().
		SetContext(ctx).
		SetBasicAuth(n.Username, n.Password).
		SetError(&rsp).
		SetBody(map[string]interface{}{
			"phone":   recID,
			"message": content,
		}).
		Post(fmt.Sprintf("%v/send/message", n.WaApiBaseUrl))
	if err != nil {
		return fmt.Errorf("unable to make http request: %w", err)
	}
	if resp.IsError() {
		slog.Error("failed to send wa message", slog.String("response", rsp.String()))
		if rsp.IsSessionExpired() {
			return core.ErrSessionExpired
		}

		return fmt.Errorf("failed to send message: %s", resp.String())
	}

	return nil
}
