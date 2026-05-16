package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

// Report holds the summary of a pipeline run.
type Report struct {
	Pipeline  string
	Status    string // "ok", "failed"
	Stages    []StageReport
	StartedAt time.Time
	Elapsed   time.Duration
}

// StageReport is the summary of a stage run.
type StageReport struct {
	Name    string
	Status  string
	Elapsed time.Duration
	Steps   []StepReport
}

// StepReport is the summary of a step run.
type StepReport struct {
	Name     string
	Status   string
	Duration time.Duration
	Error    string
}

// Notifier is the interface for sending pipeline run reports.
type Notifier interface {
	Send(r Report) error
}

// --- DingTalk ---

type dingtalkMsg struct {
	MsgType string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}

// NewDingTalk creates a notifier that sends to a DingTalk robot webhook.
func NewDingTalk(webhookURL string) Notifier {
	return &dingtalk{url: webhookURL}
}

type dingtalk struct {
	url    string
	client http.Client
}

func (d *dingtalk) Send(r Report) error {
	msg := dingtalkMsg{MsgType: "markdown"}
	msg.Markdown.Title = fmt.Sprintf("CICD: %s [%s]", r.Pipeline, r.Status)
	msg.Markdown.Text = dingtalkMarkdown(r)

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := d.client.Post(d.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("dingtalk returned %d", resp.StatusCode)
	}
	return nil
}

func dingtalkMarkdown(r Report) string {
	var b strings.Builder
	icon := "✅"
	if r.Status != "ok" {
		icon = "❌"
	}
	fmt.Fprintf(&b, "## %s Pipeline: %s\n\n", icon, r.Pipeline)
	fmt.Fprintf(&b, "**Status:** %s | **Duration:** %s\n\n", r.Status, r.Elapsed.Truncate(time.Second))

	for _, s := range r.Stages {
		sicon := "✅"
		if s.Status == "failed" {
			sicon = "❌"
		} else if s.Status == "skipped" {
			sicon = "⏭"
		}
		fmt.Fprintf(&b, "### %s Stage: %s (%s)\n", sicon, s.Name, s.Elapsed.Truncate(time.Second))
		for _, step := range s.Steps {
			icon := "✅"
			if step.Status == "failed" {
				icon = "❌"
			} else if step.Status == "skipped" {
				icon = "⏭"
			}
			fmt.Fprintf(&b, "- %s **%s** (%s)", icon, step.Name, step.Duration.Truncate(time.Second))
			if step.Error != "" {
				fmt.Fprintf(&b, " — %s", step.Error)
			}
			fmt.Fprintln(&b)
		}
		fmt.Fprintln(&b)
	}
	return b.String()
}

// --- Email ---

// NewEmail creates a notifier that sends via SMTP.
func NewEmail(smtpHost, smtpPort, user, password, to string) Notifier {
	return &email{
		host:     smtpHost,
		port:     smtpPort,
		user:     user,
		password: password,
		to:       to,
	}
}

type email struct {
	host, port, user, password, to string
}

func (e *email) Send(r Report) error {
	subject := fmt.Sprintf("CICD: %s [%s]", r.Pipeline, r.Status)
	body := emailBody(r)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		e.user, e.to, subject, body)

	auth := smtp.PlainAuth("", e.user, e.password, e.host)
	return smtp.SendMail(e.host+":"+e.port, auth, e.user, []string{e.to}, []byte(msg))
}

func emailBody(r Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Pipeline: %s\nStatus: %s | Duration: %s\n\n", r.Pipeline, r.Status, r.Elapsed.Truncate(time.Second))
	for _, s := range r.Stages {
		fmt.Fprintf(&b, "Stage: %s [%s] (%s)\n", s.Name, s.Status, s.Elapsed.Truncate(time.Second))
		for _, step := range s.Steps {
			line := fmt.Sprintf("  - %s [%s] (%s)", step.Name, step.Status, step.Duration.Truncate(time.Second))
			if step.Error != "" {
				line += fmt.Sprintf(" — %s", step.Error)
			}
			fmt.Fprintln(&b, line)
		}
		fmt.Fprintln(&b)
	}
	return b.String()
}
