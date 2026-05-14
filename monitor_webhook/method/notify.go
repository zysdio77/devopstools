package method

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Notifier interface {
	Send(content string) error
}

var httpClient = resty.New().
	SetTimeout(10 * time.Second).
	SetRetryCount(2).
	SetRetryWaitTime(1 * time.Second)

// --- DingTalk ---

type dingTalkNotifier struct {
	webhook string
}

func NewDingTalk(webhook string) Notifier {
	return &dingTalkNotifier{webhook: webhook}
}

func (d *dingTalkNotifier) Send(content string) error {
	body := RequestBody{
		Msgtype: "text",
		Text:    Text{Content: content},
	}
	var result Result
	resp, err := httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&result).
		Post(d.webhook)
	if err != nil {
		logrus.WithError(err).Error("dingtalk request failed")
		return err
	}
	if result.Errcode != 0 {
		logrus.WithField("errcode", result.Errcode).
			WithField("errmsg", result.Errmsg).
			Error("dingtalk returned error")
		return fmt.Errorf("dingtalk errcode=%d", result.Errcode)
	}
	logrus.WithField("status", resp.StatusCode()).Info("dingtalk sent")
	return nil
}

// --- Email ---

type emailNotifier struct {
	from string
	to   []string
	auth smtp.Auth
	addr string
}

func NewEmail(cfg *SMTPConfig, to []string) Notifier {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return &emailNotifier{
		from: cfg.User,
		to:   to,
		auth: smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host),
		addr: addr,
	}
}

func (e *emailNotifier) Send(content string) error {
	msg := buildEmail(e.from, e.to, content)

	host, _, _ := net.SplitHostPort(e.addr)
	tlsConfig := &tls.Config{ServerName: host}

	conn, err := tls.Dial("tcp", e.addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer c.Close()

	if err := c.Auth(e.auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := c.Mail(e.from); err != nil {
		return fmt.Errorf("smtp mail: %w", err)
	}
	for _, addr := range e.to {
		if err := c.Rcpt(addr); err != nil {
			return fmt.Errorf("smtp rcpt %s: %w", addr, err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close: %w", err)
	}
	logrus.WithField("to", strings.Join(e.to, ",")).Info("email sent")
	return nil
}

func buildEmail(from string, to []string, content string) string {
	var buf strings.Builder
	buf.WriteString("From: " + from + "\r\n")
	buf.WriteString("To: " + strings.Join(to, ",") + "\r\n")
	buf.WriteString("Subject: Webhook Notification\r\n")
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(content)
	return buf.String()
}
