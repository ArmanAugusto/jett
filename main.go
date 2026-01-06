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
   Colors (Catppuccin Mocha)
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
   Storage
   ========================= */

func tasksFilePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, "tasks.json"), nil
}

func loadTasks() (TaskList, error) {
	var list TaskList
	path, err := tasksFilePath()
	if err != nil {
		return list, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TaskList{Tasks: []Task{}}, nil
		}
		return list, err
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
	path, err := tasksFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
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
	case "high", "h":
		return PriorityHigh, nil
	case "medium", "m":
		return PriorityMedium, nil
	case "low", "l":
		return PriorityLow, nil
	}
	return "", fmt.Errorf("invalid priority")
}

func statusFromString(s string) (Status, error) {
	switch strings.ToLower(s) {
	case "pending", "p":
		return StatusPending, nil
	case "done", "d":
		return StatusDone, nil
	}
	return "", fmt.Errorf("invalid status")
}

func colorPriority(p Priority) string {
	switch p {
	case PriorityHigh:
		return fgPriorityHigh + "High" + reset
	case PriorityMedium:
		return fgPriorityMedium + "Medium" + reset
	case PriorityLow:
		return fgPriorityLow + "Low" + reset
	}
	return string(p)
}

func statusAndColor(t Task, now time.Time) (string, string) {
	if t.Status == StatusDone {
		return "Done", fgStatusDone
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	due := time.Date(t.DueDate.Year(), t.DueDate.Month(), t.DueDate.Day(), 0, 0, 0, 0, t.DueDate.Location())

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
	fmt.Printf("%s┃              Jett CLI                ┃%s\n", fgTitle, reset)
	fmt.Printf("%s┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛%s\n", fgPink, reset)
}

/* =========================
   Commands
   ========================= */

func cmdAdd(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: jett add \"Title\" START DUE [PRIORITY]")
		return
	}

	start, _ := parseDate(args[1])
	due, _ := parseDate(args[2])
	prio := PriorityMedium
	if len(args) > 3 {
		prio, _ = priorityFromString(args[3])
	}

	list, _ := loadTasks()
	task := Task{
		ID:        nextID(list),
		Title:     args[0],
		StartDate: start,
		DueDate:   due,
		Priority:  prio,
		Status:    StatusPending,
	}
	list.Tasks = append(list.Tasks, task)
	saveTasks(list)

	printHeader()
	fmt.Printf("%sAdded:%s %s\n", fgText, reset, task.Title)
}

func cmdList(args []string) {
	list, _ := loadTasks()
	printHeader()

	var pf *Priority
	var sf *Status

	for _, a := range args {
		if p, err := priorityFromString(a); err == nil {
			pf = &p
		} else if s, err := statusFromString(a); err == nil {
			sf = &s
		}
	}

	sort.Slice(list.Tasks, func(i, j int) bool {
		return list.Tasks[i].DueDate.Before(list.Tasks[j].DueDate)
	})

	now := time.Now()
	shown := 0

	for _, t := range list.Tasks {
		if pf != nil && t.Priority != *pf {
			continue
		}
		if sf != nil && t.Status != *sf {
			continue
		}

		st, col := statusAndColor(t, now)
		fmt.Printf("%s[%2d]%s %s\n", fgPink, t.ID, reset, t.Title)
		fmt.Printf("  Due: %s  Priority: %s\n", t.DueDate.Format("2006-01-02"), colorPriority(t.Priority))
		fmt.Printf("  Status: %s%s%s\n\n", col, st, reset)
		shown++
	}

	if shown == 0 {
		fmt.Printf("%sNo matching tasks%s\n", fgMuted, reset)
	}
}

func cmdDone(args []string) {
	id, _ := strconv.Atoi(args[0])
	list, _ := loadTasks()
	for i := range list.Tasks {
		if list.Tasks[i].ID == id {
			list.Tasks[i].Status = StatusDone
		}
	}
	saveTasks(list)
}

func cmdDelete(args []string) {
	id, _ := strconv.Atoi(args[0])
	list, _ := loadTasks()
	out := []Task{}
	for _, t := range list.Tasks {
		if t.ID != id {
			out = append(out, t)
		}
	}
	list.Tasks = out
	saveTasks(list)
}

func cmdEdit(args []string) {
	id, _ := strconv.Atoi(args[0])
	list, _ := loadTasks()

	for i, t := range list.Tasks {
		if t.ID == id {
			for _, kv := range args[1:] {
				parts := strings.SplitN(kv, "=", 2)
				if len(parts) != 2 {
					continue
				}
				switch parts[0] {
				case "title":
					t.Title = parts[1]
				case "due":
					t.DueDate, _ = parseDate(parts[1])
				case "priority":
					t.Priority, _ = priorityFromString(parts[1])
				case "status":
					t.Status, _ = statusFromString(parts[1])
				}
			}
			list.Tasks[i] = t
		}
	}
	saveTasks(list)
}

func cmdSummary() {
	list, _ := loadTasks()
	printHeader()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	total := len(list.Tasks)
	overdue := 0
	todayCnt := 0
	done := 0

	for _, t := range list.Tasks {
		if t.Status == StatusDone {
			done++
			continue
		}
		due := time.Date(t.DueDate.Year(), t.DueDate.Month(), t.DueDate.Day(), 0, 0, 0, 0, t.DueDate.Location())
		if due.Before(today) {
			overdue++
		} else if due.Equal(today) {
			todayCnt++
		}
	}

	fmt.Printf("Total: %d\n", total)
	fmt.Printf("%sOverdue:%s %d\n", fgStatusOverdue, reset, overdue)
	fmt.Printf("%sDue today:%s %d\n", fgStatusDueToday, reset, todayCnt)
	fmt.Printf("%sDone:%s %d\n", fgStatusDone, reset, done)
}

/* =========================
   Main
   ========================= */

func printUsage() {
	printHeader()
	fmt.Println("jett add \"Title\" START DUE [PRIORITY]")
	fmt.Println("jett list [high|medium|low] [pending|done]")
	fmt.Println("jett edit <id> field=value")
	fmt.Println("jett delete <id>")
	fmt.Println("jett done <id>")
	fmt.Println("jett summary")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "add":
		cmdAdd(args)
	case "list":
		cmdList(args)
	case "edit":
		cmdEdit(args)
	case "delete":
		cmdDelete(args)
	case "done":
		cmdDone(args)
	case "summary":
		cmdSummary()
	default:
		printUsage()
	}
}
