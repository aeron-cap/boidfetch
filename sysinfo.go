package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
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

	return lines
}

func osPrettyName() string {
	platform := runtime.GOOS
	if platform == "linux" {
		f, err := os.Open("/etc/os-release")
		if err != nil {
			return platform
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
			}
		}
		return platform
	}
	return platform
}
