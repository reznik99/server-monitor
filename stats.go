package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
)

const TEMPERATURE_FILE = "/sys/class/thermal/thermal_zone0/temp" // File with board temperature
const THRESHOLD_TEMP = 60.00                                     // Degrees C
const THRESHOLD_MEM = 75.00                                      // Memory Usage %
const THRESHOLD_CPU = 75.00                                      // CPU Usage %

func GetMemoryStats() (*memory.Stats, error) {
	return memory.Get()
}

func GetCPUStats() (*cpu.Stats, error) {
	return cpu.Get()
}

func GetNetworkStats() ([]network.Stats, error) {
	return network.Get()
}

func GetBoardTemp() (float32, error) {
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

func HumanFriendlyBytes(lengthInBytes uint64) string {
	length := float64(lengthInBytes)
	if length < 1024 {
		return fmt.Sprintf("%.2fBytes", length)
	} else if length/1024 < 1024 {
		return fmt.Sprintf("%.2fKBytes", length/1024)
	} else if length/1024/1024 < 1024 {
		return fmt.Sprintf("%.2fMBytes", length/1024/1024)
	} else if length/1024/1024/1024 < 1024 {
		return fmt.Sprintf("%.2fGBytes", length/1024/1024/1024)
	} else {
		return fmt.Sprintf("%.2fTBytes", length/1024/1024/1024/1024)
	}
}
