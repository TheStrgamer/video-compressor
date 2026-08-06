package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func get_duration(path string) float64 {
	fmt.Println("Getting duration")

	cmd := exec.Command("ffprobe", "-i", path, "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1")

	stdout, err := cmd.Output()
	duration, err := strconv.ParseFloat(strings.TrimSpace(string(stdout)), 64)
	if err != nil {
		fmt.Println("Error parsing duration: ", err)
		return -1
	}
	return duration
}
