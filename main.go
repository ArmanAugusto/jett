package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

/* =========================
   Types
   ========================= */

type Priority string

const (
	PriorityHigh   Priority = "High"
	PriorityMedium Priority = "Medium"
	PriorityLow    Priority = "Low"
)

type Status string

const (
	StatusPending Status = "Pending"
	StatusDone    Status = "Done"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date"`
	DueDate   time.Time `json:"due_date"`
	Priority  Priority  `json:"priority"`
	Status    Status    `json:"status"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

/* =========================
   Colors (Catppuccin Mocha-ish)
   ========================= */

const (
	reset = "\033[0m"

	fgText  = "\033[38;5;250m"
	fgMuted = "\033[38;5;244m"

	fgPink  = "\033[38;5;211m"
	fgTitle = "\033[38;5;219m"

	fgPriorityHigh   = "\033[38;5;203m"
	fgPriorityMedium = "\033[38;5;221m"
	fgPriorityLow    = "\033[38;5;114m"

	fgStatusOnSchedule = "\033[38;5;114m"
	fgStatusDueToday   = "\033[38;5;221m"
	fgStatusOverdue    = "\033[38;5;203m"
	fgStatusDone       = "\033[38;5;244m"
)

/* =========================
   Storage Paths
   ========================= */

func tasksFilePath() (string, error) {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "tasks.json"), nil
}

func logFilePath() (string, error) {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "jett.log"), nil
}

/* =========================
   Logging
   ========================= */

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

	f.WriteString(line)
}

/* =========================
   Task Storage
   ========================= */

func loadTasks() (TaskList, error) {
	var list TaskList
	path, _ := tasksFilePath()

	data, err := os.ReadFile(path)
	if err != nil {
		return TaskList{Tasks: []Task{}}, nil
	}

	if len(data) == 0 {
		return TaskList{Tasks: []Task{}}, nil
	}

	if err := json.Unmarshal(data, &list); err != nil {
		return TaskList{Tasks: []Task{}}, nil
	}

	return list, nil
}

func saveTasks(list TaskList) error {
	path, _ := tasksFilePath()
	data, _ := json.MarshalIndent(list, "", "  ")
	return os.WriteFile(path, data, 0644)
}

func nextID(list TaskList) int {
	max := 0
	for _, t := range list.Tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

/* =========================
   Helpers
   ========================= */

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
	fmt.Printf("%s┃               Jett CLI              ┃%s\n", fgTitle, reset)
	fmt.Printf("%s┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛%s\n", fgPink, reset)
}

/* =========================
   Commands
   ========================= */

func cmdAdd(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: jett add \"Title\" START_DATE DUE_DATE [priority]")
		return
	}

	title := args[0]
	start, _ := parseDate(args[1])
	due, _ := parseDate(args[2])

	prio := PriorityMedium
	if len(args) > 3 {
		prio, _ = priorityFromString(args[3])
	}

	list, _ := loadTasks()

	task := Task{
		ID:        nextID(list),
		Title:     title,
		StartDate: start,
		DueDate:   due,
		Priority:  prio,
		Status:    StatusPending,
	}

	list.Tasks = append(list.Tasks, task)
	saveTasks(list)

	logAction("ADD", task.ID, task.Title)

	printHeader()
	fmt.Printf("%sAdded:%s %s\n", fgText, reset, task.Title)
}

func cmdDone(args []string) {
	id, _ := strconv.Atoi(args[0])
	list, _ := loadTasks()

	for i := range list.Tasks {
		if list.Tasks[i].ID == id {
			list.Tasks[i].Status = StatusDone
			logAction("DONE", id, list.Tasks[i].Title)
		}
	}

	saveTasks(list)
}

func cmdDelete(args []string) {
	id, _ := strconv.Atoi(args[0])
	list, _ := loadTasks()

	out := []Task{}
	for _, t := range list.Tasks {
		if t.ID == id {
			logAction("DELETE", id, t.Title)
			continue
		}
		out = append(out, t)
	}

	list.Tasks = out
	saveTasks(list)
}

func cmdEdit(args []string) {
	id, _ := strconv.Atoi(args[0])
	list, _ := loadTasks()

	for i, t := range list.Tasks {
		if t.ID != id {
			continue
		}

		changes := []string{}

		for _, kv := range args[1:] {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			val := parts[1]

			switch key {
			case "title":
				t.Title = val
			case "due":
				t.DueDate, _ = parseDate(val)
			case "priority":
				t.Priority, _ = priorityFromString(val)
			case "status":
				t.Status, _ = statusFromString(val)
			}

			changes = append(changes, kv)
		}

		list.Tasks[i] = t
		logAction("EDIT", id, strings.Join(changes, ", "))
	}

	saveTasks(list)
}

func cmdList(args []string) {
	list, _ := loadTasks()
	printHeader()

	sort.Slice(list.Tasks, func(i, j int) bool {
		return list.Tasks[i].DueDate.Before(list.Tasks[j].DueDate)
	})

	for _, t := range list.Tasks {
		status, col := statusAndColor(t)
		fmt.Printf("%s[%d]%s %s\n", fgPink, t.ID, reset, t.Title)
		fmt.Printf("  Due: %s  Status: %s%s%s\n\n",
			t.DueDate.Format("2006-01-02"),
			col, status, reset)
	}
}

/* =========================
   Main
   ========================= */

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Jett CLI")
		fmt.Println("Commands: add, list, edit, delete, done")
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "add":
		cmdAdd(args)
	case "list":
		cmdList(args)
	case "done":
		cmdDone(args)
	case "delete":
		cmdDelete(args)
	case "edit":
		cmdEdit(args)
	default:
		fmt.Println("Unknown command.")
	}
}
