package browser

import (
	"log"
	"os"
	"time"
)

var (
	// Logger is the default logger for the browser package
	Logger = log.New(os.Stdout, "[BROWSER] ", log.LstdFlags)
)

// LogInfo logs an informational message
func LogInfo(format string, args ...interface{}) {
	Logger.Printf("[INFO] "+format, args...)
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	Logger.Printf("[ERROR] "+format, args...)
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	Logger.Printf("[DEBUG] "+format, args...)
}

// LogWithTime logs a message with a timestamp
func LogWithTime(format string, args ...interface{}) {
	Logger.Printf("[%s] "+format, append([]interface{}{time.Now().Format("15:04:05")}, args...)...)
}
