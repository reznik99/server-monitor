package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/reznik99/server-monitor/internal/monitor"
)

var (
	Version         string  = "Development"                           // Tagged Version of binary (from git)
	temperatureFile string  = "/sys/class/thermal/thermal_zone0/temp" // File with board temperature
	thresholdTemp   float32 = 60.00                                   // Degrees C
	thresholdMem    float32 = 75.00                                   // Memory Usage %
	thresholdCPU    float32 = 75.00                                   // CPU Usage %
	hostName        string  = "N/A"                                   // Hostname of machine
	serverName      string  = "N/A"                                   // Custom server name for alert email subject
)

func main() {
	startTime := time.Now()
	defer func() {
		logrus.Infof("Server-Monitor %s executed in %dms", Version, time.Since(startTime).Milliseconds())
	}()

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("Error getting hostname: %s", err)
	}

	// Allow env override of thresholds and parameters
	hostName = monitor.FirstOrDefault(hostname, hostName)
	serverName = monitor.FirstOrDefault(os.Getenv("SERVER_NAME"), serverName)
	thresholdTemp = monitor.FirstOrDefaultFloat(os.Getenv("THRESHOLD_TEMP"), thresholdTemp)
	thresholdMem = monitor.FirstOrDefaultFloat(os.Getenv("THRESHOLD_MEM"), thresholdMem)
	thresholdCPU = monitor.FirstOrDefaultFloat(os.Getenv("THRESHOLD_CPU"), thresholdCPU)

	logrus.Infof("Thresholds - TEMP: %.2fc | MEM: %.2f%% | CPU: %.2f%%", thresholdTemp, thresholdMem, thresholdCPU)

	// Get OS statistics
	stats, err := monitor.GetAllStats(temperatureFile)
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
	if stats.Temperature > thresholdTemp {
		logrus.Infof("CPU Temperature %.2fc above threshold of %.2fc", stats.Temperature, thresholdTemp)
		doSendAlert = true
	}
	if stats.MemoryPercentage > thresholdMem {
		logrus.Infof("Memory usage %.2f%% above threshold %.2f%%", stats.MemoryPercentage, thresholdMem)
		doSendAlert = true
	}
	if stats.CPUPercentage > thresholdCPU {
		logrus.Infof("CPU usage %.2f%% above threshold %.2f%%", stats.CPUPercentage, thresholdCPU)
		doSendAlert = true
	}
	if float32(stats.LoadAvg.Loadavg5*100) > thresholdCPU {
		logrus.Infof("CPU load avg (5min) %.2f%% above threshold %.2f%%", stats.LoadAvg.Loadavg5*100, thresholdCPU)
		doSendAlert = true
	}
	if float32(stats.LoadAvg.Loadavg15*100) > thresholdCPU {
		logrus.Infof("CPU load avg (15min) %.2f%% above threshold %.2f%%", stats.LoadAvg.Loadavg15*100, thresholdCPU)
		doSendAlert = true
	}

	if doSendAlert {
		logrus.Info("Sending email alert")
		if err = monitor.SendEmailAlert(stats, serverName, hostName, Version); err != nil {
			logrus.Errorf("Error sending email alert: %s", err)
		}
	}
}
