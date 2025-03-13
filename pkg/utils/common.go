package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// GetFileSHA256 calculates the SHA-256 hash of a file
func GetFileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// GetOSInfo returns information about the operating system
func GetOSInfo() string {
	return runtime.GOOS + " " + runtime.GOARCH
}

// SplitLastN splits a string by a separator and returns the last n parts
func SplitLastN(s, sep string, n int) []string {
	parts := strings.Split(s, sep)
	if len(parts) <= n {
		return parts
	}

	return parts[len(parts)-n:]
}

// IsFileExist checks if a file exists
func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// IsDirExist checks if a directory exists
func IsDirExist(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FormatDuration formats a duration as a human-readable string
func FormatDuration(d time.Duration) string {
	// Format durations like "2h 3m 4.5s"
	d = d.Round(100 * time.Millisecond)

	var result strings.Builder

	if d >= time.Hour {
		result.WriteString(strings.TrimSuffix(d.String(), "0s"))
		return result.String()
	}

	if d >= time.Minute {
		result.WriteString(strings.TrimSuffix(d.String(), "0s"))
		return result.String()
	}

	return d.String()
}

// FormatBytes formats a byte count as a human-readable string
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return formatWithPrecision(float64(bytes), 0) + " B"
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	precision := 0
	if bytes%div != 0 {
		precision = 1
	}

	return formatWithPrecision(float64(bytes)/float64(div), precision) + " " + string("KMGTPE"[exp]) + "B"
}

// formatWithPrecision formats a float with a specific precision
func formatWithPrecision(n float64, precision int) string {
	format := "%." + string(rune(precision+'0')) + "f"
	return strings.TrimSuffix(strings.TrimSuffix(
		strings.Replace(strings.Replace(
			strings.Replace(strings.TrimSuffix(
				strings.TrimSuffix(
					strings.TrimRight(
						string([]byte(fmt.Sprintf(format, n))),
						"0"),
					"."),
				"."),
				".", "X", 1),
			",", ".", -1),
			"X", ",", 1),
		","), ".")
}

// Truncate truncates a string to a specific length with an optional suffix
func Truncate(s string, length int, suffix string) string {
	if len(s) <= length {
		return s
	}

	return s[:length] + suffix
}
