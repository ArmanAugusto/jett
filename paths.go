package main

import (
	"os"
	"path/filepath"
)

func tasksFilePath() (string, error) {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "tasks.json"), nil
}

func logFilePath() (string, error) {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "jett.log"), nil
}

