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
	logrus.Infof("- Memory  -> Tot:  %v | Used: %v | Free: %v", HumanFriendlyBytes(stats.Memory.Total), HumanFriendlyBytes(stats.Memory.Used), HumanFriendlyBytes(stats.Memory.Free))
	logrus.Infof("- CPU     -> Tot:  %v | Idle: %v | System: %v | User: %v", stats.CPU.Total, stats.CPU.Idle, stats.CPU.System, stats.CPU.User)
	logrus.Infof("- Network -> Name: %v | Rx: %v | Tx: %v", stats.Net.Name, HumanFriendlyBytes(stats.Net.RxBytes), HumanFriendlyBytes(stats.Net.TxBytes))
	logrus.Infof("- Uptime  -> %s", stats.Uptime.String())
	logrus.Info("Derived stats:")
	logrus.Infof("- Memory usage %.2f%%", stats.MemoryPercentage)
	logrus.Infof("- CPU usage    %.2f%%", stats.CPUPercentage)
	logrus.Infof("- Temperature  %.2fC", stats.Temperature)

	// Alerting logic
	doSendAlert := false
	if stats.Temperature > THRESHOLD_TEMP {
		logrus.Infof("CPU Temperature %.2f above threshold of %.2f: Sending email alert", stats.Temperature, THRESHOLD_TEMP)
		doSendAlert = true
	} else if stats.MemoryPercentage > THRESHOLD_MEM {
		logrus.Infof("Memory usage %.2f%% above threshold %.2f%%: Sending email alert", stats.MemoryPercentage, THRESHOLD_MEM)
		doSendAlert = true
	} else if stats.CPUPercentage > THRESHOLD_CPU {
		logrus.Infof("CPU usage %.2f%% above threshold %.2f%%: Sending email alert", stats.CPUPercentage, THRESHOLD_CPU)
		doSendAlert = true
	}
	// TODO: network/uptime/loadavg alert checks

	if doSendAlert {
		if err = SendEmailAlert(stats); err != nil {
			logrus.Errorf("Error sending email alert: %s", err)
		}
	}
}
