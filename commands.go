package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

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
	_ = saveTasks(list)

	logAction("ADD", task.ID, task.Title)

	printHeader()
	fmt.Printf("%sAdded:%s %s\n", fgText, reset, task.Title)
}

func cmdDone(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: jett done ID")
		return
	}

	id, _ := strconv.Atoi(args[0])
	list, _ := loadTasks()

	for i := range list.Tasks {
		if list.Tasks[i].ID == id {
			list.Tasks[i].Status = StatusDone
			logAction("DONE", id, list.Tasks[i].Title)
		}
	}

	_ = saveTasks(list)
}

func cmdDelete(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: jett delete ID")
		return
	}

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
	_ = saveTasks(list)
}

func cmdEdit(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: jett edit ID key=value [key=value ...]")
		return
	}

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

	_ = saveTasks(list)
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

