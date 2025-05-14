package models

import (
	"fmt"
	"time"

	"github.com/gookit/color"
)

// Priority represents the importance level of a task
type Priority int

const (
	Low Priority = iota
	Medium
	High
	Critical
)

// String returns the string representation of a priority
func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// ColoredString returns a colored string representation of a priority
func (p Priority) ColoredString() string {
	switch p {
	case Low:
		return color.Green.Sprint(p.String())
	case Medium:
		return color.Yellow.Sprint(p.String())
	case High:
		return color.FgLightRed.Sprint(p.String())
	case Critical:
		return color.Style{color.FgRed, color.Bold}.Sprint(p.String())
	default:
		return p.String()
	}
}

// Task represents a task in the task manager
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    Priority  `json:"priority"`
	Completed   bool      `json:"completed"`
	Progress    int       `json:"progress"` // 0-100 percentage
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FormatDueDate returns a formatted string of the due date
func (t Task) FormatDueDate() string {
	if t.DueDate.IsZero() {
		return "No due date"
	}

	// Format: May 14, 2025
	return t.DueDate.Format("Jan 02, 2006")
}

// Status returns the current status of the task
func (t Task) Status() string {
	if t.Completed {
		return "Completed"
	}
	if t.DueDate.IsZero() {
		return fmt.Sprintf("In Progress (%d%%)", t.Progress)
	}
	if time.Now().After(t.DueDate) {
		return "Overdue"
	}
	return fmt.Sprintf("In Progress (%d%%)", t.Progress)
}

// ColoredStatus returns a colored string representation of the task status
func (t Task) ColoredStatus() string {
	if t.Completed {
		return color.Green.Sprint("✓ Completed")
	}
	if t.DueDate.IsZero() {
		return color.Blue.Sprintf("⧖ In Progress (%d%%)", t.Progress)
	}
	if time.Now().After(t.DueDate) {
		return color.Red.Sprint("! Overdue")
	}
	return color.Blue.Sprintf("⧖ In Progress (%d%%)", t.Progress)
}

// DaysLeft returns the number of days left before the due date
func (t Task) DaysLeft() int {
	if t.DueDate.IsZero() || t.Completed {
		return -1
	}

	now := time.Now()
	if now.After(t.DueDate) {
		return -2 // Indicates overdue
	}

	days := int(t.DueDate.Sub(now).Hours() / 24)
	return days
}

// FormattedDaysLeft returns a formatted string of days left
func (t Task) FormattedDaysLeft() string {
	days := t.DaysLeft()

	if days == -1 {
		return ""
	}

	if days == -2 {
		return color.Red.Sprint("Overdue")
	}

	if days == 0 {
		return color.FgLightRed.Sprint("Due today!")
	}

	if days == 1 {
		return color.Yellow.Sprint("Due tomorrow!")
	}

	return color.Style{color.FgBlue}.Sprintf("%d days left", days)
}
