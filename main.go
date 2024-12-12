package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Version = "Development"

func main() {
	logrus.Infof("Server-Monitor %q started", Version)

	// Performance metrics
	startTime := time.Now()
	defer func() {
		logrus.Infof("Server-Monitor %q executed in %dms", Version, time.Since(startTime).Milliseconds())
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
	tempBoard, err := GetBoardTemp()
	if err != nil {
		logrus.Errorf("Failed to get read board temperature: %s", err)
	}
	memoryUsagePercent := float32(memStats.Total-memStats.Free) / float32(memStats.Total) * 100
	cpuUsagePercent := float32(cpuStats.Total-cpuStats.Idle) / float32(cpuStats.Total) * 100

	// Print statistics
	logrus.Info("Raw stats:")
	logrus.Infof("- Memory  -> Tot:  %v | Used: %v | Free: %v", HumanFriendlyBytes(memStats.Total), HumanFriendlyBytes(memStats.Used), HumanFriendlyBytes(memStats.Free))
	logrus.Infof("- CPU     -> Tot:  %v | Idle: %v | System: %v | User: %v", cpuStats.Total, cpuStats.Idle, cpuStats.System, cpuStats.User)
	logrus.Infof("- Network -> Name: %v | Rx: %v | Tx: %v", netStats[0].Name, HumanFriendlyBytes(netStats[0].RxBytes), HumanFriendlyBytes(netStats[0].TxBytes))
	logrus.Info("Derived stats:")
	logrus.Infof("- Memory usage %.2f%%", memoryUsagePercent)
	logrus.Infof("- CPU usage    %.2f%%", cpuUsagePercent)
	logrus.Infof("- Temperature  %.2fC", tempBoard)

	// Alerting logic
	if tempBoard > THRESHOLD_TEMP {
		logrus.Infof("CPU Temperature %.2f above threshold of %.2f: Sending email alert", tempBoard, THRESHOLD_TEMP)
		err = SendEmailAlert()
	} else if memoryUsagePercent > THRESHOLD_MEM {
		logrus.Infof("Memory usage %.2f%% above threshold %.2f%%: Sending email alert", memoryUsagePercent, THRESHOLD_MEM)
		err = SendEmailAlert()
	} else if cpuUsagePercent > THRESHOLD_CPU {
		logrus.Infof("CPU usage %.2f%% above threshold %.2f%%: Sending email alert", cpuUsagePercent, THRESHOLD_CPU)
		err = SendEmailAlert()
	}
	if err != nil {
		logrus.Fatalf("Error sending email alert: %s", err)
	}
}
