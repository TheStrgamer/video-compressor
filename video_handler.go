package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type VideoInfo struct {
	Duration float64 // seconds
	Width    int
	Height   int
	Bitrate  int64 // kilobits per second
	SizeMB   float64
}

func (info VideoInfo) printInfo() {
	fmt.Println("Duration: ", info.Duration, "s")
	fmt.Println("Size: ", info.Width, "x", info.Height)
	fmt.Println("bitrate: ", info.Bitrate, "kbps")
	fmt.Println("size: ", info.SizeMB, "MB")

}

type ffprobeOutput struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
		BitRate  string `json:"bit_rate"`
		Size     string `json:"size"`
	} `json:"format"`
}

func GetVideoInfo(path string) (*VideoInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height:format=duration,bit_rate,size",
		"-of", "json",
		path,
	)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	var raw ffprobeOutput
	if err := json.Unmarshal(stdout, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe json: %w", err)
	}

	info := &VideoInfo{}

	fmt.Sscanf(raw.Format.Duration, "%f", &info.Duration)
	fmt.Sscanf(raw.Format.BitRate, "%d", &info.Bitrate)
	info.Bitrate = info.Bitrate / 1024 //kbps
	if len(raw.Streams) > 0 {
		info.Width = raw.Streams[0].Width
		info.Height = raw.Streams[0].Height
	}

	if fi, err := os.Stat(path); err == nil {
		info.SizeMB = float64(fi.Size()) / (1024 * 1024)
	}

	return info, nil
}

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

func runFFmpegCompress(path, outputPath string, videoBitrateKbps int, audioBitrateKbps int, targetHeight int) error {
	args := []string{
		"-y",
		"-i", path,
	}

	if targetHeight > 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=-2:%d", targetHeight))
	}

	videoBitrateStr := strconv.Itoa(videoBitrateKbps) + "k"
	audioBitrateStr := strconv.Itoa(audioBitrateKbps) + "k"
	bufsize := strconv.Itoa(videoBitrateKbps*2) + "k"

	args = append(args,
		"-c:v", "libx264",
		"-preset", "medium",
		"-b:v", videoBitrateStr,
		"-maxrate", videoBitrateStr,
		"-bufsize", bufsize,
		"-c:a", "aac",
		"-b:a", audioBitrateStr,
		outputPath,
	)

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\n%s", err, output)
	}
	return nil
}

func CompressToSize(path, outputPath string, targetSizeMB float64) error {
	info, err := GetVideoInfo(path)
	if err != nil {
		return fmt.Errorf("failed to read video info: %w", err)
	}

	audioBitrateKbps := 128
	totalBitrateKbps := int((targetSizeMB * 8192) / info.Duration)
	videoBitrateKbps := totalBitrateKbps - audioBitrateKbps

	if videoBitrateKbps <= 0 {
		return fmt.Errorf("target size %.2fMB is too small for a %.2fs video", targetSizeMB, info.Duration)
	}

	targetHeight, ok := pickResolution(info.Height, videoBitrateKbps)
	if !ok {
		fmt.Printf("warning: even at %dp, %dkbps is below recommended quality — result may look compressed\n", targetHeight, videoBitrateKbps)
	}

	// no scaling needed
	scaleHeight := 0
	if targetHeight != info.Height {
		scaleHeight = targetHeight
		fmt.Printf("downscaling %dp -> %dp to keep %dkbps looking reasonable\n", info.Height, targetHeight, videoBitrateKbps)
	}

	return runFFmpegCompress(path, outputPath, videoBitrateKbps, audioBitrateKbps, scaleHeight)
}
