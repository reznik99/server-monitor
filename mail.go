package main

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/wneessen/go-mail"
)

func SendEmailAlert(stats Stats) error {
	// Get email data
	var From = os.Getenv("SOURCE_EMAIL_ADDRESS")
	var To = os.Getenv("TARGET_EMAIL_ADDRESS")
	var Subject = fmt.Sprintf("%s server-monitor alert", SERVER_NAME)

	var data = EmailTemplateData{
		Subject:           Subject,
		TempBoard:         fmt.Sprintf("%.2fc", stats.Temperature),
		MemoryTotal:       Humanize(stats.Memory.Total),
		MemoryUsed:        Humanize(stats.Memory.Used),
		MemoryUsedPercent: fmt.Sprintf("%.2f%%", stats.MemoryPercentage),
		CPUUsagePercent:   fmt.Sprintf("%.2f%%", stats.CPUPercentage),
		CPUUsageAvg:       fmt.Sprintf("%f, %f, %f", stats.LoadAvg.Loadavg1, stats.LoadAvg.Loadavg5, stats.LoadAvg.Loadavg15),
		RxBytes:           Humanize(stats.Net.RxBytes),
		TxBytes:           Humanize(stats.Net.TxBytes),
		UpTime:            stats.Uptime.String(),
		DateTime:          time.Now().Format(time.RFC822),
		ServerName:        SERVER_NAME,
		HostName:          HOST_NAME,
	}
	// Parse email template
	emailTmpl, err := template.New("EmailAlert").Parse(EmailTemplateStr)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %s", err)
	}

	// Create email object
	msg := mail.NewMsg()
	msg.Subject(Subject)
	if err := msg.From(From); err != nil {
		return fmt.Errorf("invalid from email address '%s': %s", From, err)
	}
	if err := msg.To(To); err != nil {
		return fmt.Errorf("invalid to email address '%s': %s", To, err)
	}
	if err := msg.SetBodyHTMLTemplate(emailTmpl, data); err != nil {
		return fmt.Errorf("failed to set email template: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send email
	return msg.WriteToSendmailWithContext(ctx, mail.SendmailPath)
}
