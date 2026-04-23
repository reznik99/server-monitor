// Package monitor collects OS-level health metrics and sends email alerts.
package monitor

import (
	"errors"
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

type Stats struct {
	Memory           *memory.Stats
	CPU              *cpu.Stats
	LoadAvg          *loadavg.Stats
	Net              *network.Stats
	Uptime           time.Duration
	MemoryPercentage float32
	CPUPercentage    float32
	Temperature      float32
}

func GetAllStats(temperatureFile string) (Stats, error) {
	var stats Stats
	var err error
	var e error

	stats.Memory, e = getMemoryStats()
	if e != nil {
		err = errors.Join(err, fmt.Errorf("memory: %w", e))
	}
	stats.CPU, e = getCPUStats()
	if e != nil {
		err = errors.Join(err, fmt.Errorf("cpu: %w", e))
	}
	stats.LoadAvg, e = getLoadAvg()
	if e != nil {
		err = errors.Join(err, fmt.Errorf("loadavg: %w", e))
	}
	stats.Net, e = getNetworkStats()
	if e != nil {
		err = errors.Join(err, fmt.Errorf("network: %w", e))
	}
	stats.Uptime, e = getUptime()
	if e != nil {
		err = errors.Join(err, fmt.Errorf("uptime: %w", e))
	}
	stats.Temperature, e = getBoardTemp(temperatureFile)
	if e != nil {
		logrus.Warnf("Failed to read board temperature: %s", e)
	}

	// Get nice % values (guard against nil pointers)
	if stats.Memory != nil {
		stats.MemoryPercentage = 100 - (float32(stats.Memory.Total-stats.Memory.Used) / float32(stats.Memory.Total) * 100)
	}
	if stats.CPU != nil {
		stats.CPUPercentage = float32(stats.CPU.Total-stats.CPU.Idle) / float32(stats.CPU.Total) * 100
	}

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

func getBoardTemp(temperatureFile string) (float32, error) {
	rawTemp, err := os.ReadFile(temperatureFile)
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

func FirstOrDefaultFloat(overrideVal string, defaultVal float32) float32 {
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

func FirstOrDefault(val string, defaultVal string) string {
	if val == "" {
		return defaultVal
	} else {
		return val
	}
}

// DurationToString returns a stringified version of duration including days
func DurationToString(val time.Duration) string {
	uptimeString := val.String()
	if val >= DAY {
		days := int(val.Hours()) / 24
		hours := val - time.Duration(days)*DAY
		uptimeString = fmt.Sprintf("%dd%s", days, hours)
	}
	return uptimeString
}

// Humanize converts byte length into human friendly string such as 950Bytes or 55.55KB or 22.22MB
func Humanize(lengthInBytes uint64) string {
	length := float64(lengthInBytes)
	switch {
	case length < 1024:
		return fmt.Sprintf("%.2fBytes", length)
	case length/1024 < 1024:
		return fmt.Sprintf("%.2fKB", length/1024)
	case length/1024/1024 < 1024:
		return fmt.Sprintf("%.2fMB", length/1024/1024)
	case length/1024/1024/1024 < 1024:
		return fmt.Sprintf("%.2fGB", length/1024/1024/1024)
	default:
		return fmt.Sprintf("%.2fTB", length/1024/1024/1024/1024)
	}
}
