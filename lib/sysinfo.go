package lib

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
)

//SystemInfo contains basic information about system load
type SystemInfo struct {
	Memory      sigar.Mem
	Swap        sigar.Swap
	Uptime      int
	UptimeS     string
	LoadAvg     sigar.LoadAverage
	CPUList     sigar.CpuList
	Arch        string
	Os          string
	CurrentTime time.Time
	MemTotal    uint64
	MemFree     uint64
	MemAvail    uint64
	MemUsed     uint64
	SwapTotal   uint64
	SwapFree    uint64
	SwapUsed    uint64
}

//GetSystemInfo returns short info about system load
func GetSystemInfo() SystemInfo {
	s := SystemInfo{}

	uptime := sigar.Uptime{}
	if err := uptime.Get(); err == nil {
		s.Uptime = int(uptime.Length)
		s.UptimeS = uptime.Format()
	}

	avg := sigar.LoadAverage{}
	if err := avg.Get(); err == nil {
		s.LoadAvg = avg
	}

	s.CurrentTime = time.Now()

	mem := sigar.Mem{}
	if err := mem.Get(); err == nil {
		s.Memory = mem
	}

	swap := sigar.Swap{}
	if err := swap.Get(); err == nil {
		s.Swap = swap
	}

	cpulist := sigar.CpuList{}
	if err := cpulist.Get(); err == nil {
		s.CPUList = cpulist
	}

	s.Arch = runtime.GOARCH
	s.Os = runtime.GOOS

	memValues := []string{"MemTotal", "MemFree", "MemAvailable", "SwapTotal", "SwapFree"}

	for _, memLabel := range memValues {
		memInfo, err := os.Open("/proc/meminfo")
		if err != nil {
			panic(err)
		}
		defer memInfo.Close()
		b := bufio.NewScanner(memInfo)
		for b.Scan() {
			var n uint64
			if nItems, _ := fmt.Sscanf(b.Text(), memLabel+": %d kB", &n); nItems == 1 {
				if memLabel == "MemTotal" {
					s.MemTotal = n
				}
				if memLabel == "MemFree" {
					s.MemFree = n
				}
				if memLabel == "MemAvailable" {
					s.MemAvail = n
				}
				if memLabel == "SwapTotal" {
					s.SwapTotal = n
				}
				if memLabel == "SwapFree" {
					s.SwapFree = n
				}

			}
		}
	}
	s.MemUsed = (s.MemAvail - s.MemFree)
	s.SwapUsed = (s.SwapTotal - s.SwapFree)

	return s
}
