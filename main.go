package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/mackerelio/go-osstat/cpu"
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

	// Enviroment variable
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %s", err)
	}

	// Get OS statistics
	memStats, err := GetMemoryStats()
	if err != nil {
		logrus.Errorf("Failed to get memory statistics: %s", err)
	}
	cpuStats, err := GetCPUStats()
	if err != nil {
		logrus.Errorf("Failed to get cpu statistics: %s", err)
	}
	netStats, err := GetNetworkStats()
	if err != nil {
		logrus.Errorf("Failed to get network statistics: %s", err)
	}
	uptime, err := GetUptime()
	if err != nil {
		logrus.Errorf("Failed to get uptime: %s", err)
	}
	tempBoard, err := GetBoardTemp()
	if err != nil {
		logrus.Errorf("Failed to get read board temperature: %s", err)
	}

	stats := Stats{
		Memory:           memStats,
		MemoryPercentage: 100 - (float32(memStats.Total-memStats.Used) / float32(memStats.Total) * 100),
		CPU:              cpuStats,
		CPUPercentage:    float32(cpuStats.Total-cpuStats.Idle) / float32(cpuStats.Total) * 100,
		Net:              netStats[0],
		Temperature:      tempBoard,
		Uptime:           uptime,
	}

	// Print statistics
	logrus.Info("Raw stats:")
	logrus.Infof("- Memory  -> Tot:  %v | Used: %v | Free: %v", HumanFriendlyBytes(memStats.Total), HumanFriendlyBytes(memStats.Used), HumanFriendlyBytes(memStats.Free))
	logrus.Infof("- CPU     -> Tot:  %v | Idle: %v | System: %v | User: %v", cpuStats.Total, cpuStats.Idle, cpuStats.System, cpuStats.User)
	logrus.Infof("- Network -> Name: %v | Rx: %v | Tx: %v", netStats[0].Name, HumanFriendlyBytes(netStats[0].RxBytes), HumanFriendlyBytes(netStats[0].TxBytes))
	logrus.Info("Derived stats:")
	logrus.Infof("- Memory usage %.2f%%", stats.MemoryPercentage)
	logrus.Infof("- CPU usage    %.2f%%", stats.CPUPercentage)
	logrus.Infof("- Temperature  %.2fC", stats.Temperature)

	// Alerting logic
	// TODO: Networking alert checks
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
	if doSendAlert {
		if err = SendEmailAlert(stats); err != nil {
			logrus.Errorf("Error sending email alert: %s", err)
		}
	}
}
