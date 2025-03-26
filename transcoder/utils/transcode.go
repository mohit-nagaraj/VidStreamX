package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func TranscodeVideo(inputPath, outputPath string, width, height int) error {
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("Transcoding video: %s to %d\n", inputPath, height)
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath, // Input file
		"-vf", fmt.Sprintf("scale=%d:%d", width, height), // Resize video
		"-c:v", "libx264", // Use H.264 video codec
		"-crf", "23", // Constant Rate Factor (quality)
		"-preset", "medium", // Encoding speed/quality tradeoff
		"-c:a", "copy", // Copy audio without re-encoding
		outputPath, // Output file
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg command failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Successfully transcoded video: %s\n", outputPath)
	return nil
}
