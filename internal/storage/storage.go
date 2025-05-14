package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"taskmaster/internal/models"
	"time"
)

// Storage defines the interface for task storage implementations
type Storage interface {
	Init() error
	Close() error
	CreateTask(*models.Task) error
	GetTask(int64) (*models.Task, error)
	GetAllTasks() ([]*models.Task, error)
	UpdateTask(*models.Task) error
	DeleteTask(int64) error
	CompleteTask(int64) error
	UpdateTaskProgress(int64, int) error
}

// FileStorage implements Storage using file system
type FileStorage struct {
	baseDir     string
	tasksDir    string
	counterFile string
	mu          sync.Mutex
	nextID      int64
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(targetDir string) (*FileStorage, error) {
	// If no target directory provided, use current directory
	if targetDir == "" {
		// Get the current working directory as the target directory
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		targetDir = cwd
	}

	// Create the tasks directory
	tasksDir := filepath.Join(targetDir, ".taskmaster")
	err := os.MkdirAll(tasksDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create tasks directory: %w", err)
	}

	storage := &FileStorage{
		baseDir:     targetDir,
		tasksDir:    tasksDir,
		counterFile: filepath.Join(tasksDir, "counter.json"),
	}

	return storage, nil
}

// Init initializes the file storage
func (s *FileStorage) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create the tasks directory if it doesn't exist
	if err := os.MkdirAll(s.tasksDir, 0755); err != nil {
		return fmt.Errorf("failed to create tasks directory: %w", err)
	}

	// Read counter file or initialize it
	counter, err := os.ReadFile(s.counterFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a new counter starting at 1
			s.nextID = 1
			return s.saveCounter()
		}
		return fmt.Errorf("failed to read counter file: %w", err)
	}

	// Parse counter value
	var counterData struct {
		NextID int64 `json:"next_id"`
	}
	if err := json.Unmarshal(counter, &counterData); err != nil {
		return fmt.Errorf("failed to parse counter file: %w", err)
	}

	s.nextID = counterData.NextID
	return nil
}

// saveCounter saves the current counter value
func (s *FileStorage) saveCounter() error {
	counterData := struct {
		NextID int64 `json:"next_id"`
	}{
		NextID: s.nextID,
	}

	data, err := json.MarshalIndent(counterData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal counter: %w", err)
	}

	err = os.WriteFile(s.counterFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write counter file: %w", err)
	}

	return nil
}

// Close closes the storage
func (s *FileStorage) Close() error {
	// No connections to close in file-based storage
	return nil
}

// CreateTask inserts a new task into storage
func (s *FileStorage) CreateTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Set ID and timestamps
	task.ID = s.nextID
	s.nextID++

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	// Save the counter
	if err := s.saveCounter(); err != nil {
		return err
	}

	// Save the task
	return s.saveTask(task)
}

// getTaskFilename returns the filename for a task
func (s *FileStorage) getTaskFilename(id int64) string {
	return filepath.Join(s.tasksDir, fmt.Sprintf("task_%d.json", id))
}

// saveTask saves a task to a file
func (s *FileStorage) saveTask(task *models.Task) error {
	filename := s.getTaskFilename(task.ID)

	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	return nil
}

// GetTask retrieves a task by ID
func (s *FileStorage) GetTask(id int64) (*models.Task, error) {
	filename := s.getTaskFilename(id)

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("task with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to read task file: %w", err)
	}

	var task models.Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("failed to parse task file: %w", err)
	}

	return &task, nil
}

// GetAllTasks retrieves all tasks
func (s *FileStorage) GetAllTasks() ([]*models.Task, error) {
	var tasks []*models.Task

	// Read all task files
	files, err := os.ReadDir(s.tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), "task_") || !strings.HasSuffix(file.Name(), ".json") {
			continue // Skip directories and non-task files
		}

		// Extract ID from filename
		idStr := strings.TrimPrefix(file.Name(), "task_")
		idStr = strings.TrimSuffix(idStr, ".json")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue // Skip files with invalid ID format
		}

		task, err := s.GetTask(id)
		if err != nil {
			continue // Skip tasks that can't be loaded
		}

		tasks = append(tasks, task)
	}

	// Sort tasks by ID
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks, nil
}

// UpdateTask updates an existing task
func (s *FileStorage) UpdateTask(task *models.Task) error {
	// Check if task exists
	_, err := s.GetTask(task.ID)
	if err != nil {
		return err
	}

	// Update timestamp
	task.UpdatedAt = time.Now()

	// Save the updated task
	return s.saveTask(task)
}

// DeleteTask deletes a task by ID
func (s *FileStorage) DeleteTask(id int64) error {
	filename := s.getTaskFilename(id)

	// Check if file exists
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("task with ID %d not found", id)
		}
		return fmt.Errorf("failed to access task file: %w", err)
	}

	// Delete the file
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed to delete task file: %w", err)
	}

	return nil
}

// CompleteTask marks a task as completed
func (s *FileStorage) CompleteTask(id int64) error {
	task, err := s.GetTask(id)
	if err != nil {
		return err
	}

	task.Completed = true
	task.Progress = 100
	task.UpdatedAt = time.Now()

	return s.saveTask(task)
}

// UpdateTaskProgress updates the progress of a task
func (s *FileStorage) UpdateTaskProgress(id int64, progress int) error {
	task, err := s.GetTask(id)
	if err != nil {
		return err
	}

	task.Progress = progress
	task.UpdatedAt = time.Now()

	// If progress is 100%, mark as completed
	if progress == 100 {
		task.Completed = true
	}

	return s.saveTask(task)
}
