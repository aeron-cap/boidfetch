package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type InfoLine struct {
	Label string
	Value string
}

func GatherInfo() []InfoLine {
	lines := []InfoLine{}

	hostname, _ := os.Hostname()
	user := os.Getenv("USER")

	lines = append(lines, InfoLine{"User", fmt.Sprintf("%s@%s", user, hostname)})
	lines = append(lines, InfoLine{"OS", osPrettyName()})
	lines = append(lines, InfoLine{"Kernel", kernelVersion()})
	lines = append(lines, InfoLine{"Uptime", uptime()})
	lines = append(lines, InfoLine{"Shell", filepath.Base(os.Getenv("SHELL"))})
	lines = append(lines, InfoLine{"Terminal", terminalName()})
	lines = append(lines, InfoLine{"WM", windowManager()})
	lines = append(lines, InfoLine{"Arch", runtime.GOARCH})
	lines = append(lines, InfoLine{"Memory", memInfo()})

	return lines
}

func osPrettyName() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return runtime.GOOS
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Err() != nil {
			continue
		}
		
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
		}
	}
	return runtime.GOOS 
}

func kernelVersion() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func uptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "unknown"
	}
	
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return "unknown"
	}

	secs, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "uknown"
	}

	d := time.Duration(secs) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", h, m)
}

func terminalName() string {
	if tp := os.Getenv("TERM_PROGRAM"); tp != "" {
		return tp
	}
	return os.Getenv("TERM")
}

func windowManager() string {
	if wm := os.Getenv("XDG_CURRENT_DESKTOP"); wm != "" {
		return wm
	}
	return "unknown"
}

func memInfo() string {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return "unknown"
	}
	defer f.Close()

	var totalKB, availableKB int64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Err() != nil {
			continue
		}

		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		val, _ := strconv.ParseInt(fields[1], 10, 64)
		switch fields[0] {
			case "MemTotal:":
			totalKB = val
		   case "MemAvailable:":
	        availableKB = val
		}
	}
	usedGB := float64(totalKB - availableKB) / 1024 / 1024
	totalGB := float64(totalKB) / 1024 / 1024
	return fmt.Sprintf("%.2f GB / %.2f GB", usedGB, totalGB)
}