package utils

import (
	"fmt"
	"os/exec"
)

func TranscodeVideo(inputPath, outputPath string, width, height int) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=%d:%d", width, height), 
		"-c:a", "copy", 
		outputPath, 
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("FFmpeg command failed: %w", err)
	}

	return nil
}
