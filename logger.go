package main

import (
	"fmt"
	"os"
	"time"
)

func logAction(action string, taskID int, detail string) {
	path, err := logFilePath()
	if err != nil {
		return
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s | %-6s | Task #%d | %s\n",
		timestamp, action, taskID, detail)

	_, _ = f.WriteString(line)
}

