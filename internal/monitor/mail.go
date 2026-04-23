package monitor

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/wneessen/go-mail"
)

func SendEmailAlert(stats Stats, serverName, hostName, version string) error {
	// Get email data
	var From = os.Getenv("SOURCE_EMAIL_ADDRESS")
	var To = os.Getenv("TARGET_EMAIL_ADDRESS")
	var Subject = serverName + " server-monitor alert"

	var data = EmailTemplateData{
		Subject:           Subject,
		TempBoard:         fmt.Sprintf("%.2fc", stats.Temperature),
		MemoryTotal:       Humanize(stats.Memory.Total),
		MemoryUsed:        Humanize(stats.Memory.Used),
		MemoryUsedPercent: fmt.Sprintf("%.2f%%", stats.MemoryPercentage),
		CPUUsagePercent:   fmt.Sprintf("%.2f%%", stats.CPUPercentage),
		CPUUsageAvg:       fmt.Sprintf("%.2f, %.2f, %.2f", stats.LoadAvg.Loadavg1, stats.LoadAvg.Loadavg5, stats.LoadAvg.Loadavg15),
		RxBytes:           Humanize(stats.Net.RxBytes),
		TxBytes:           Humanize(stats.Net.TxBytes),
		UpTime:            DurationToString(stats.Uptime),
		DateTime:          time.Now().Format(time.RFC822),
		ServerName:        serverName,
		HostName:          hostName,
		ProgVersion:       version,
	}
	// Parse email template
	emailTmpl, err := template.New("EmailAlert").Parse(EmailTemplateStr)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Create email object
	msg := mail.NewMsg()
	msg.Subject(Subject)
	if err := msg.From(From); err != nil {
		return fmt.Errorf("invalid from email address '%s': %w", From, err)
	}
	if err := msg.To(To); err != nil {
		return fmt.Errorf("invalid to email address '%s': %w", To, err)
	}
	if err := msg.SetBodyHTMLTemplate(emailTmpl, data); err != nil {
		return fmt.Errorf("failed to set email template: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send email
	return msg.WriteToSendmailWithContext(ctx, mail.SendmailPath)
}
