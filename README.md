# 🛩️ Jett — Catppuccin Mocha Todo CLI (Go)

**Jett** is a lightweight, beautiful command-line todo tracker written in **Go**.

It is designed to help you manage tasks with:

- Start dates
- Due dates
- Priority levels
- Status tracking
- Catppuccin Mocha-inspired terminal colors
- Full action logging (audit trail)

---

## ✨ Features

✅ Add tasks with start + due dates  
✅ Priority support: **High / Medium / Low**  
✅ Task status: **Pending / Done**  
✅ Colored output using Catppuccin Mocha styling  
✅ Edit tasks in place  
✅ Delete tasks  
✅ Mark tasks as done  
✅ Persistent storage using JSON  
✅ Append-only action log (`jett.log`)  

---

## 📦 Installation

### Clone or create the project folder

```bash
mkdir jett
cd jett
```

### Initialize Go module

```bash
go mod init jett
go mod tidy
```

### Build the binary

```bash
go build -o jett .
```

### Run it:

```bash
./jett list
```

## Usage

### Add a Task

```bash
jett add "Pay rent" 2026-02-01 2026-02-05 High
```

### Format

```bash
jett add "Title" START_DATE DUE_DATE [priority]
```

Priority options:
- High
- Medium (default)
- Low

### List Tasks

```bash
jett list
```

Displays all tasks with:
- Due date
- Priority
- Status color

### Mark Task Done

```bash
jett done 1
```

### Edit a Task

```bash
jett edit 2 title="New title" due=2026-02-10 priority=Low
```

### Delete a Task

```bash
jett delete 3
```

### Task Storage

Tasks are stored locally in:
tasks.json

#### Example:

```json
{
  "tasks": [
    {
      "id": 1,
      "title": "Pay rent",
      "start_date": "2026-02-01T00:00:00Z",
      "due_date": "2026-02-05T00:00:00Z",
      "priority": "High",
      "status": "Pending"
    }
  ]
}
```

### Theme

Jett uses a Catppuccin Mocha-inspired palette:
- Pink accents for task IDs
- Green = on schedule
- Yellow = due today
- Red = overdue
- Gray = completed

### Future Improvements (Roadmap)

Planned upgrades:
- jett list today / jett list overdue
- jett summary dashboard
- Config stored in ~/.config/jett/
- Shell completions (bash/zsh/fish)
- Recurring tasks
- Full-screen TUI mode (Bubble Tea)

###  License

MIT License--free to use, modify, and distribute.

## Built with Love In Go + Catppuccin Mocha

### Enjoy your tasks in style
