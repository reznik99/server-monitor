package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
	"github.com/mackerelio/go-osstat/uptime"
	"github.com/sirupsen/logrus"
)

const DAY = time.Hour * 24

func GetAllStats() (Stats, error) {
	var err error
	var stats Stats
	// Get OS Stats
	stats.Memory, err = getMemoryStats()
	if err != nil {
		logrus.Errorf("Failed to get memory statistics: %s", err)
	}
	stats.CPU, err = getCPUStats()
	if err != nil {
		logrus.Errorf("Failed to get cpu statistics: %s", err)
	}
	stats.LoadAvg, err = getLoadAvg()
	if err != nil {
		logrus.Errorf("Failed to get load average: %s", err)
	}
	stats.Net, err = getNetworkStats()
	if err != nil {
		logrus.Errorf("Failed to get network statistics: %s", err)
	}
	stats.Uptime, err = getUptime()
	if err != nil {
		logrus.Errorf("Failed to get uptime: %s", err)
	}
	stats.Temperature, err = getBoardTemp()
	if err != nil {
		logrus.Errorf("Failed to get read board temperature: %s", err)
	}
	// Get nice % values
	stats.MemoryPercentage = 100 - (float32(stats.Memory.Total-stats.Memory.Used) / float32(stats.Memory.Total) * 100)
	stats.CPUPercentage = float32(stats.CPU.Total-stats.CPU.Idle) / float32(stats.CPU.Total) * 100

	return stats, err
}

func getMemoryStats() (*memory.Stats, error) {
	return memory.Get()
}

func getCPUStats() (*cpu.Stats, error) {
	return cpu.Get()
}

func getLoadAvg() (*loadavg.Stats, error) {
	return loadavg.Get()
}

func getNetworkStats() (*network.Stats, error) {
	stats, err := network.Get()
	if err != nil {
		return &network.Stats{}, err
	}

	return &stats[0], nil
}

func getUptime() (time.Duration, error) {
	return uptime.Get()
}

func getBoardTemp() (float32, error) {
	rawTemp, err := os.ReadFile(TEMPERATURE_FILE)
	if err != nil {
		return -1, err
	}
	stringTemp := strings.TrimSuffix(string(rawTemp), "\n")
	intFloatTemp, err := strconv.Atoi(stringTemp)
	if err != nil {
		return -1, err
	}
	return float32(intFloatTemp) / 1000, nil
}

func firstOrDefaultFloat(overrideVal string, defaultVal float32) float32 {
	if overrideVal != "" {
		overrideFlt, err := strconv.ParseFloat(overrideVal, 32)
		if err != nil {
			logrus.Warnf("Failed to parse threshold: %s", err)
			return defaultVal
		}
		return float32(overrideFlt)
	}

	return defaultVal
}

func firstOrDefault(val string, defaultVal string) string {
	if val == "" {
		return defaultVal
	} else {
		return val
	}
}

// returns a stringified version of duration including days
func durationToString(val time.Duration) string {
	uptimeString := val.String()
	if val >= DAY {
		days := val / DAY
		hours := val - days*DAY
		uptimeString = fmt.Sprintf("%dd%s", days, hours)
	}
	return uptimeString
}

// Convert byte length into human friendly string such as 950Bytes or 55.55KB or 22.22MB
func Humanize(lengthInBytes uint64) string {
	length := float64(lengthInBytes)
	if length < 1024 {
		return fmt.Sprintf("%.2fBytes", length)
	} else if length/1024 < 1024 {
		return fmt.Sprintf("%.2fKB", length/1024)
	} else if length/1024/1024 < 1024 {
		return fmt.Sprintf("%.2fMB", length/1024/1024)
	} else if length/1024/1024/1024 < 1024 {
		return fmt.Sprintf("%.2fGB", length/1024/1024/1024)
	} else {
		return fmt.Sprintf("%.2fTB", length/1024/1024/1024/1024)
	}
}
