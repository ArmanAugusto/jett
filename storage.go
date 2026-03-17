package main

import (
	"encoding/json"
	"os"
)

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

