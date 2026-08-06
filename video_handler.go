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

func get_target_bitrate(path string, target_size_mb float64) (target_bitrate_kbps int) {
	audio_bitrate_kbps := 128
	target_bitrate_kbps = int((target_size_mb*8192)/get_duration(path)) - audio_bitrate_kbps
	return
}
