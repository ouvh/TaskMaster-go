package app

import (
	"errors"
	"fmt"
	"time"

	"taskmaster/internal/models"
	"taskmaster/internal/storage"
)

// App represents the core application that manages tasks
type App struct {
	storage storage.Storage
}

// NewApp creates a new application instance
func NewApp(s storage.Storage) *App {
	return &App{storage: s}
}

// Initialize initializes the application
func (a *App) Initialize() error {
	return a.storage.Init()
}

// Close cleans up resources
func (a *App) Close() error {
	return a.storage.Close()
}

// CreateTask creates a new task
func (a *App) CreateTask(title, desc string, dueDate time.Time, priority models.Priority) (*models.Task, error) {
	if title == "" {
		return nil, errors.New("task title cannot be empty")
	}

	task := &models.Task{
		Title:       title,
		Description: desc,
		DueDate:     dueDate,
		Priority:    priority,
		Progress:    0,
		Completed:   false,
	}

	err := a.storage.CreateTask(task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// GetTask retrieves a task by ID
func (a *App) GetTask(id int64) (*models.Task, error) {
	return a.storage.GetTask(id)
}

// GetAllTasks retrieves all tasks
func (a *App) GetAllTasks() ([]*models.Task, error) {
	return a.storage.GetAllTasks()
}

// DeleteTask deletes a task
func (a *App) DeleteTask(id int64) error {
	return a.storage.DeleteTask(id)
}

// CompleteTask marks a task as completed
func (a *App) CompleteTask(id int64) error {
	return a.storage.CompleteTask(id)
}

// UpdateTaskProgress updates the progress of a task
func (a *App) UpdateTaskProgress(id int64, progress int) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100, got %d", progress)
	}
	return a.storage.UpdateTaskProgress(id, progress)
}

// UpdateTaskDetails updates a task's details
func (a *App) UpdateTaskDetails(id int64, title, desc string, dueDate time.Time, priority models.Priority) error {
	task, err := a.GetTask(id)
	if err != nil {
		return err
	}

	// Update the fields
	task.Title = title
	task.Description = desc
	task.DueDate = dueDate
	task.Priority = priority

	return a.storage.UpdateTask(task)
}
