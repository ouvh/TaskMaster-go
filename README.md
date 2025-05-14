# TaskMaster - Command-Line Task Manager

TaskMaster is a powerful command-line task manager built with Go. It allows you to manage your tasks, track progress, set due dates, and organize by priority - all from your terminal!
![image](https://github.com/user-attachments/assets/5e415025-de8a-485f-b904-5a9d007aa4f1)

## Features

- **Simple Command-Line Interface**: Easy to use commands with a clean output
- **Project-Specific Tasks**: Tasks are stored locally in your project directory
- **Task Management**: Create, edit, delete and view detailed information about your tasks
- **Due Dates**: Set and track due dates for your tasks
- **Priority Levels**: Organize tasks by priority (Low, Medium, High, Critical)
- **Progress Tracking**: Update and visualize task completion progress
- **Deadline Views**: Quickly see upcoming deadlines and days remaining
- **Local Data Storage**: Tasks are stored as JSON files in your current directory

## Installation

### Prerequisites

- Go 1.18 or later

### Building from source

1. Clone the repository:
   ```bash
   git clone https://github.com/ouvh/TaskMaster-go.git
   cd TaskMaster-go
   ```

2. Build the application:
   ```bash
   go build -o taskmaster ./cmd/taskmaster
   ```

3. Run the application:
   ```bash
   ./taskmaster
   ```

## Usage

TaskMaster uses a command-line interface with the following commands:

### Basic Commands

```bash
# Show all available commands
taskmaster help

# List all tasks
taskmaster list

# Show upcoming deadlines
taskmaster deadlines
# or
taskmaster due
```

### Managing Tasks

```bash
# Create a new task
taskmaster create --title "Task name" --desc "Task description" --due 2023-12-31 --priority 2

# View a task's details
taskmaster view 3

# Edit a task
taskmaster edit 3 --title "Updated title" --priority 1

# Update a task's progress
taskmaster progress 3 75

# Mark a task as complete
taskmaster complete 3

# Delete a task
taskmaster delete 3
```

### Priority Levels

- 0 - Low
- 1 - Medium
- 2 - High
- 3 - Critical

## Project Structure

```
taskmaster/
├── cmd/
│   └── taskmaster/
│       └── main.go           # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go            # Core application logic
│   │   └── cli.go            # Command-line interface
│   ├── models/
│   │   └── task.go           # Task data model
│   └── storage/
│       ├── storage.go        # Storage interface
├── go.mod                    # Go module file
└── README.md                 # This file
```

## Data Storage

TaskMaster stores your tasks as JSON files in a `.taskmaster` directory within your current working directory. Each task is stored as a separate file named `task_[ID].json`. This allows you to have project-specific tasks that stay with your project.

Benefits of this approach:
- Tasks stay with your project directory
- No database dependencies required
- Human-readable JSON files
- Easy to backup or include in version control
- Portable across different machines

## Quick Start

To get started with TaskMaster immediately:

```bash
# Install directly using go
go install github.com/ouvh/TaskMaster-go/cmd/taskmaster@latest

# Create your first task
taskmaster create --title "My first task" --desc "Getting started with TaskMaster" --priority 1

# List all your tasks
taskmaster list
```

## Configuration

TaskMaster can be configured by creating a `.taskmasterrc` file in your home directory:

```json
{
  "defaultPriority": 1,
  "defaultDueDays": 7,
  "colorOutput": true,
  "dateFormat": "2006-01-02"
}
```

## Development

### Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

### Building for Different Platforms

```bash
# For Windows
GOOS=windows GOARCH=amd64 go build -o taskmaster.exe ./cmd/taskmaster

# For macOS
GOOS=darwin GOARCH=amd64 go build -o taskmaster ./cmd/taskmaster

# For Linux
GOOS=linux GOARCH=amd64 go build -o taskmaster ./cmd/taskmaster
```

## Troubleshooting

### Common Issues

- **Tasks not saving**: Ensure you have write permissions in the current directory
- **Date format errors**: Use YYYY-MM-DD format for dates
- **Command not found**: Make sure the binary is in your PATH

### Reporting Bugs

If you encounter any bugs, please create an issue on GitHub with:
- Your OS and Go version
- Steps to reproduce
- Expected vs actual behavior

