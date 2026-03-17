package main

import (
	"fmt"
	"strings"
	"time"
)

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func priorityFromString(s string) (Priority, error) {
	switch strings.ToLower(s) {
	case "high":
		return PriorityHigh, nil
	case "medium":
		return PriorityMedium, nil
	case "low":
		return PriorityLow, nil
	}
	return "", fmt.Errorf("invalid priority")
}

func statusFromString(s string) (Status, error) {
	switch strings.ToLower(s) {
	case "pending":
		return StatusPending, nil
	case "done":
		return StatusDone, nil
	}
	return "", fmt.Errorf("invalid status")
}

func dateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func statusAndColor(t Task) (string, string) {
	if t.Status == StatusDone {
		return "Done", fgStatusDone
	}

	today := dateOnly(time.Now())
	due := dateOnly(t.DueDate)

	if due.Before(today) {
		return "Overdue", fgStatusOverdue
	}
	if due.Equal(today) {
		return "Due today", fgStatusDueToday
	}
	return "On schedule", fgStatusOnSchedule
}

func printHeader() {
	fmt.Printf("%s┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓%s\n", fgPink, reset)
	fmt.Printf("%s┃               Jett CLI               ┃%s\n", fgTitle, reset)
	fmt.Printf("%s┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛%s\n", fgPink, reset)
}

