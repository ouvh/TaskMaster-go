package app

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"taskmaster/internal/models"
	"time"

	"github.com/fatih/color"
)

// RunCLI runs the application in command-line interface mode
func RunCLI(app *App) error {
	// Define the available commands
	commands := map[string]func([]string) error{
		"list":      func(args []string) error { return listTasks(app) },
		"create":    func(args []string) error { return createTask(app, args) },
		"view":      func(args []string) error { return viewTask(app, args) },
		"edit":      func(args []string) error { return editTask(app, args) },
		"progress":  func(args []string) error { return updateProgress(app, args) },
		"complete":  func(args []string) error { return completeTask(app, args) },
		"delete":    func(args []string) error { return deleteTask(app, args) },
		"due":       func(args []string) error { return showDeadlines(app) },
		"help":      func(args []string) error { return showHelp() },
		"deadlines": func(args []string) error { return showDeadlines(app) },
	}

	// Check if command is provided
	if len(os.Args) < 2 {
		return showHelp()
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	// Check if command exists
	cmdFunc, exists := commands[cmd]
	if !exists {
		return fmt.Errorf("unknown command: %s\nRun 'taskmaster help' for usage", cmd)
	}

	// Execute the command
	return cmdFunc(args)
}

// showHelp displays usage information
func showHelp() error {
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()

	// Print header
	fmt.Println(blue("╔════════════════════════════════════════════════╗"))
	fmt.Println(blue("║") + cyan("              TASKMASTER CLI TOOL              ") + blue("║"))
	fmt.Println(blue("╚════════════════════════════════════════════════╝"))
	fmt.Println()

	// Print usage
	fmt.Println(yellow("USAGE:"))
	fmt.Println("  taskmaster [command] [arguments]")
	fmt.Println()

	// Print available commands
	fmt.Println(yellow("COMMANDS:"))
	fmt.Println("  " + green("list") + "                    List all tasks")
	fmt.Println("  " + green("create") + "                  Create a new task")
	fmt.Printf("    %s --title \"Task Title\" [--desc \"Description\"] [--due YYYY-MM-DD] [--priority 0-3]\n",
		green("    taskmaster create"))
	fmt.Println("  " + green("view") + " [id]               View details of a task")
	fmt.Println("  " + green("edit") + " [id]               Edit a task")
	fmt.Printf("    %s --title \"New Title\" [--desc \"New Description\"] [--due YYYY-MM-DD] [--priority 0-3]\n",
		green("    taskmaster edit [id]"))
	fmt.Println("  " + green("progress") + " [id] [value]   Update task progress (0-100)")
	fmt.Println("  " + green("complete") + " [id]          Mark a task as complete")
	fmt.Println("  " + green("delete") + " [id]            Delete a task")
	fmt.Println("  " + green("help") + "                   Show this help message")
	fmt.Println()

	// Print priority levels
	fmt.Println(yellow("PRIORITY LEVELS:"))
	fmt.Println("  0 - Low")
	fmt.Println("  1 - Medium")
	fmt.Println("  2 - High")
	fmt.Println("  3 - Critical")
	fmt.Println()

	// Print examples
	fmt.Println(yellow("EXAMPLES:"))
	fmt.Println("  taskmaster create --title \"Finish report\" --due 2023-05-15 --priority 2")
	fmt.Println("  taskmaster edit 5 --title \"Updated title\" --priority 3")
	fmt.Println("  taskmaster progress 3 75")
	fmt.Println("  taskmaster complete 2")

	// Print available commands (add these lines to the existing commands list)
	fmt.Println("  " + green("deadlines") + " / " + green("due") + "      Show upcoming deadlines")

	return nil
}

// listTasks lists all tasks
func listTasks(app *App) error {
	tasks, err := app.GetAllTasks()
	if err != nil {
		return fmt.Errorf("error retrieving tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	// Define colors
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Print table header
	fmt.Printf("%-5s %-30s %-10s %-10s %s\n",
		cyan("ID"), cyan("TITLE"), cyan("PRIORITY"), cyan("PROGRESS"), cyan("STATUS"))
	fmt.Println(strings.Repeat("-", 80))

	// Print each task
	for _, task := range tasks {
		fmt.Printf("%-5d %-30s %-10s %-10s %s\n",
			task.ID,
			truncateString(task.Title, 28),
			task.Priority.String(),
			yellow(fmt.Sprintf("%d%%", task.Progress)),
			getStatusText(task))
	}

	return nil
}

// createTask creates a new task
func createTask(app *App, args []string) error {
	// Define flags
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	titlePtr := createCmd.String("title", "", "Task title (required)")
	descPtr := createCmd.String("desc", "", "Task description")
	duePtr := createCmd.String("due", "", "Due date (YYYY-MM-DD)")
	priorityPtr := createCmd.Int("priority", 1, "Priority (0-3): 0=Low, 1=Medium, 2=High, 3=Critical")

	// Parse flags
	err := createCmd.Parse(args)
	if err != nil {
		return err
	}

	// Validate required fields
	if *titlePtr == "" {
		return errors.New("title is required")
	}

	// Process due date if provided
	var dueDate time.Time
	if *duePtr != "" {
		parsedDate, err := time.Parse("2006-01-02", *duePtr)
		if err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
		dueDate = parsedDate
	}

	// Validate priority
	if *priorityPtr < 0 || *priorityPtr > 3 {
		return errors.New("priority must be between 0 and 3")
	}
	priority := models.Priority(*priorityPtr)

	// Create the task
	task, err := app.CreateTask(*titlePtr, *descPtr, dueDate, priority)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	fmt.Printf("Task created successfully with ID: %d\n", task.ID)
	return nil
}

// viewTask shows details of a specific task
func viewTask(app *App, args []string) error {
	if len(args) < 1 {
		return errors.New("task ID is required")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	task, err := app.GetTask(id)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Define color
	bold := color.New(color.Bold).SprintFunc()

	// Print task details
	fmt.Println("\n=== Task Details ===")
	fmt.Printf("%s: %d\n", bold("ID"), task.ID)
	fmt.Printf("%s: %s\n", bold("Title"), task.Title)
	fmt.Printf("%s: %s\n", bold("Description"), task.Description)

	if !task.DueDate.IsZero() {
		fmt.Printf("%s: %s\n", bold("Due Date"), task.DueDate.Format("2006-01-02"))
	} else {
		fmt.Printf("%s: None\n", bold("Due Date"))
	}

	fmt.Printf("%s: %s\n", bold("Priority"), task.Priority.String())
	fmt.Printf("%s: %d%%\n", bold("Progress"), task.Progress)
	fmt.Printf("%s: %v\n", bold("Completed"), task.Completed)
	fmt.Printf("%s: %s\n", bold("Created At"), task.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("%s: %s\n", bold("Updated At"), task.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

// editTask edits an existing task
func editTask(app *App, args []string) error {
	if len(args) < 1 {
		return errors.New("task ID is required")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	// Get existing task to preserve fields that are not being updated
	task, err := app.GetTask(id)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Define flags for editing
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	titlePtr := editCmd.String("title", task.Title, "Task title")
	descPtr := editCmd.String("desc", task.Description, "Task description")

	var defaultDue string
	if !task.DueDate.IsZero() {
		defaultDue = task.DueDate.Format("2006-01-02")
	}
	duePtr := editCmd.String("due", defaultDue, "Due date (YYYY-MM-DD)")

	priorityPtr := editCmd.Int("priority", int(task.Priority), "Priority (0-3): 0=Low, 1=Medium, 2=High, 3=Critical")

	// Parse flags, excluding the first argument which is the task ID
	err = editCmd.Parse(args[1:])
	if err != nil {
		return err
	}

	// Process due date
	var dueDate time.Time
	if *duePtr != "" {
		parsedDate, err := time.Parse("2006-01-02", *duePtr)
		if err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
		dueDate = parsedDate
	} else {
		dueDate = task.DueDate
	}

	// Validate priority
	if *priorityPtr < 0 || *priorityPtr > 3 {
		return errors.New("priority must be between 0 and 3")
	}
	priority := models.Priority(*priorityPtr)

	// Update the task
	err = app.UpdateTaskDetails(id, *titlePtr, *descPtr, dueDate, priority)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	fmt.Printf("Task %d updated successfully\n", id)
	return nil
}

// updateProgress updates a task's progress
func updateProgress(app *App, args []string) error {
	if len(args) < 2 {
		return errors.New("task ID and progress value are required")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	progress, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid progress value: %w", err)
	}

	if progress < 0 || progress > 100 {
		return errors.New("progress must be between 0 and 100")
	}

	err = app.UpdateTaskProgress(id, progress)
	if err != nil {
		return fmt.Errorf("failed to update task progress: %w", err)
	}

	fmt.Printf("Progress for task %d updated to %d%%\n", id, progress)
	return nil
}

// completeTask marks a task as complete
func completeTask(app *App, args []string) error {
	if len(args) < 1 {
		return errors.New("task ID is required")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	err = app.CompleteTask(id)
	if err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}

	fmt.Printf("Task %d marked as complete\n", id)
	return nil
}

// deleteTask deletes a task
func deleteTask(app *App, args []string) error {
	if len(args) < 1 {
		return errors.New("task ID is required")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	// Get the task to confirm deletion
	task, err := app.GetTask(id)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete task '%s'? (y/n): ", task.Title)
	var confirm string
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	err = app.DeleteTask(id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	fmt.Printf("Task %d deleted successfully\n", id)
	return nil
}

// Helper functions (reused from tui_simple.go)
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func getStatusText(task *models.Task) string {
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if task.Completed {
		return green("Completed")
	}
	if !task.DueDate.IsZero() && time.Now().After(task.DueDate) {
		return red("Overdue")
	}
	return blue("In Progress")
}

// showDeadlines displays upcoming task deadlines
func showDeadlines(app *App) error {
	tasks, err := app.GetAllTasks()
	if err != nil {
		return fmt.Errorf("error retrieving tasks: %w", err)
	}

	// Filter tasks with deadlines that aren't completed
	var tasksWithDeadlines []*models.Task
	for _, task := range tasks {
		if !task.DueDate.IsZero() && !task.Completed {
			tasksWithDeadlines = append(tasksWithDeadlines, task)
		}
	}

	if len(tasksWithDeadlines) == 0 {
		fmt.Println("No upcoming deadlines found.")
		return nil
	}

	// Sort by deadline (already sorted by the database, but just to be sure)
	sort.Slice(tasksWithDeadlines, func(i, j int) bool {
		return tasksWithDeadlines[i].DueDate.Before(tasksWithDeadlines[j].DueDate)
	})

	// Define colors
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	// Print header
	fmt.Println(cyan("UPCOMING DEADLINES"))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-5s %-30s %-12s %-15s %s\n",
		cyan("ID"), cyan("TITLE"), cyan("PRIORITY"), cyan("DUE DATE"), cyan("DAYS LEFT"))
	fmt.Println(strings.Repeat("-", 80))

	// Print each task with deadline
	now := time.Now()
	for _, task := range tasksWithDeadlines {
		daysLeft := int(task.DueDate.Sub(now).Hours() / 24)

		var daysLeftText string
		switch {
		case daysLeft < 0:
			daysLeftText = red(fmt.Sprintf("%d days overdue", -daysLeft))
		case daysLeft == 0:
			daysLeftText = red("Due today!")
		case daysLeft == 1:
			daysLeftText = yellow("Tomorrow")
		case daysLeft <= 3:
			daysLeftText = yellow(fmt.Sprintf("%d days", daysLeft))
		default:
			daysLeftText = green(fmt.Sprintf("%d days", daysLeft))
		}

		fmt.Printf("%-5d %-30s %-12s %-15s %s\n",
			task.ID,
			truncateString(task.Title, 28),
			task.Priority.String(),
			task.DueDate.Format("2006-01-02"),
			daysLeftText)
	}

	return nil
}
