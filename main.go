package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/reznik99/server-monitor/internal/monitor"
)

var (
	Version          string  = "Development"                           // Tagged version of binary (from git)
	TEMPERATURE_FILE string  = "/sys/class/thermal/thermal_zone0/temp" // File with board temperature
	THRESHOLD_TEMP   float32 = 60.00                                   // Degrees C
	THRESHOLD_MEM    float32 = 75.00                                   // Memory Usage %
	THRESHOLD_CPU    float32 = 75.00                                   // CPU Usage %
	HOST_NAME        string  = "N/A"                                   // Hostname of machine
	SERVER_NAME      string  = "N/A"                                   // Custom server name for alert email subject
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("Error getting hostname: %s", err)
	}

	// Allow .env override of thresholds and parameters
	HOST_NAME = monitor.FirstOrDefault(hostname, HOST_NAME)
	SERVER_NAME = monitor.FirstOrDefault(os.Getenv("SERVER_NAME"), SERVER_NAME)
	THRESHOLD_TEMP = monitor.FirstOrDefaultFloat(os.Getenv("THRESHOLD_TEMP"), THRESHOLD_TEMP)
	THRESHOLD_MEM = monitor.FirstOrDefaultFloat(os.Getenv("THRESHOLD_MEM"), THRESHOLD_MEM)
	THRESHOLD_CPU = monitor.FirstOrDefaultFloat(os.Getenv("THRESHOLD_CPU"), THRESHOLD_CPU)

	logrus.Infof("Thresholds - TEMP: %.2fc | MEM: %.2f%% | CPU: %.2f%%", THRESHOLD_TEMP, THRESHOLD_MEM, THRESHOLD_CPU)
}

func main() {
	startTime := time.Now()
	defer func() {
		logrus.Infof("Server-Monitor %s executed in %dms", Version, time.Since(startTime).Milliseconds())
	}()

	// Get OS statistics
	stats, err := monitor.GetAllStats(TEMPERATURE_FILE)
	if err != nil {
		logrus.Fatalf("Error getting statistics: %s", err)
	}

	// Print statistics
	logrus.Info("Raw stats:")
	logrus.Infof("- Memory  -> Percent: %.2f%% | Tot:  %v | Used: %v | Free: %v", stats.MemoryPercentage, monitor.Humanize(stats.Memory.Total), monitor.Humanize(stats.Memory.Used), monitor.Humanize(stats.Memory.Free))
	logrus.Infof("- CPU     -> Percent: %.2f%% | Tot:  %v | Idle: %v | System: %v | User: %v", stats.CPUPercentage, stats.CPU.Total, stats.CPU.Idle, stats.CPU.System, stats.CPU.User)
	logrus.Infof("- LoadAvg -> 1m:      %v | 5m: %v | 15m: %v", stats.LoadAvg.Loadavg1, stats.LoadAvg.Loadavg5, stats.LoadAvg.Loadavg15)
	logrus.Infof("- Network -> Name:    %v | Rx: %v | Tx: %v", stats.Net.Name, monitor.Humanize(stats.Net.RxBytes), monitor.Humanize(stats.Net.TxBytes))
	logrus.Infof("- Uptime  -> %s", monitor.DurationToString(stats.Uptime))
	logrus.Infof("- Temperature  %.2fc", stats.Temperature)

	// Alerting logic
	doSendAlert := false
	if stats.Temperature > THRESHOLD_TEMP {
		logrus.Infof("CPU Temperature %.2fc above threshold of %.2fc", stats.Temperature, THRESHOLD_TEMP)
		doSendAlert = true
	}
	if stats.MemoryPercentage > THRESHOLD_MEM {
		logrus.Infof("Memory usage %.2f%% above threshold %.2f%%", stats.MemoryPercentage, THRESHOLD_MEM)
		doSendAlert = true
	}
	if stats.CPUPercentage > THRESHOLD_CPU {
		logrus.Infof("CPU usage %.2f%% above threshold %.2f%%", stats.CPUPercentage, THRESHOLD_CPU)
		doSendAlert = true
	}
	if float32(stats.LoadAvg.Loadavg5*100) > THRESHOLD_CPU {
		logrus.Infof("CPU load avg (5min) %.2f%% above threshold %.2f%%", stats.LoadAvg.Loadavg5*100, THRESHOLD_CPU)
		doSendAlert = true
	}
	if float32(stats.LoadAvg.Loadavg15*100) > THRESHOLD_CPU {
		logrus.Infof("CPU load avg (15min) %.2f%% above threshold %.2f%%", stats.LoadAvg.Loadavg15*100, THRESHOLD_CPU)
		doSendAlert = true
	}

	if doSendAlert {
		logrus.Info("Sending email alert")
		if err = monitor.SendEmailAlert(stats, SERVER_NAME, HOST_NAME, Version); err != nil {
			logrus.Errorf("Error sending email alert: %s", err)
		}
	}
}
