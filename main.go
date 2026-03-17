package main

import (
	"fmt"
	"os"
)

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

