package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
	"github.com/sirupsen/logrus"
)

var Version = "Development"

// TODO: Load .env into globals (i.e To, From, Thresholds...)

type Stats struct {
	Memory           *memory.Stats
	MemoryPercentage float32
	CPU              *cpu.Stats
	CPUPercentage    float32
	LoadAvg          *loadavg.Stats
	Net              network.Stats
	Temperature      float32
	Uptime           time.Duration
}

func main() {
	logrus.Infof("Server-Monitor %s started", Version)

	// Performance metrics
	startTime := time.Now()
	defer func() {
		logrus.Infof("Server-Monitor %s executed in %dms", Version, time.Since(startTime).Milliseconds())
	}()

	// Load enviroment variables
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %s", err)
	}

	// Get OS statistics
	stats, err := GetAllStats()
	if err != nil {
		logrus.Fatalf("Error getting statistics: %s", err)
	}

	// Print statistics
	logrus.Info("Raw stats:")
	logrus.Infof("- Memory  -> Percent: %.2f%% | Tot:  %v | Used: %v | Free: %v", stats.MemoryPercentage, Humanize(stats.Memory.Total), Humanize(stats.Memory.Used), Humanize(stats.Memory.Free))
	logrus.Infof("- CPU     -> Percent: %.2f%% | Tot:  %v | Idle: %v | System: %v | User: %v", stats.CPUPercentage, stats.CPU.Total, stats.CPU.Idle, stats.CPU.System, stats.CPU.User)
	logrus.Infof("- LoadAvg -> 1m:      %v | 5m: %v | 15m: %v", stats.LoadAvg.Loadavg1, stats.LoadAvg.Loadavg5, stats.LoadAvg.Loadavg15)
	logrus.Infof("- Network -> Name:    %v | Rx: %v | Tx: %v", stats.Net.Name, Humanize(stats.Net.RxBytes), Humanize(stats.Net.TxBytes))
	logrus.Infof("- Uptime  -> %s", stats.Uptime.String())
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
	if stats.LoadAvg.Loadavg5*100 > THRESHOLD_CPU {
		logrus.Infof("CPU load avg (5min) %.2f%% above threshold %.2f%%", stats.LoadAvg.Loadavg5*100, THRESHOLD_CPU)
		doSendAlert = true
	}
	if stats.LoadAvg.Loadavg15*100 > THRESHOLD_CPU {
		logrus.Infof("CPU load avg (15min) %.2f%% above threshold %.2f%%", stats.LoadAvg.Loadavg15*100, THRESHOLD_CPU)
		doSendAlert = true
	}

	if doSendAlert {
		logrus.Info("Sending email alert")
		if err = SendEmailAlert(stats); err != nil {
			logrus.Errorf("Error sending email alert: %s", err)
		}
	}
}
