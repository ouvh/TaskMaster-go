package main

import (
	"fmt"
	"os"

	"taskmaster/internal/app"
	"taskmaster/internal/storage"

	"github.com/fatih/color"
)

func main() {
	// Skip header if arguments are provided (i.e., we're running a command)
	if len(os.Args) <= 1 {
		printHeader()
	}

	// Get target directory (current directory by default)
	targetDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Create storage
	store, err := storage.NewFileStorage(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Initialize storage
	if err := store.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	// Create app
	taskApp := app.NewApp(store)
	defer taskApp.Close()

	// Run the CLI
	if err := app.RunCLI(taskApp); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// printHeader prints a colorful header for the app
func printHeader() {
	fmt.Println()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()

	fmt.Println(blue("╔════════════════════════════════════════════════╗"))
	fmt.Println(blue("║") + cyan("     _____ ___ ___ ___   __  __  ___ ___ _____ ___ ___  ") + blue("║"))
	fmt.Println(blue("║") + cyan("    |_   _|   / __| |__/  |  \\/  |/   / __|_   _| __| _ \\ ") + blue("║"))
	fmt.Println(blue("║") + cyan("      | | | \\| / __| |\\ \\  | |\\/| | ~ \\__ \\  | | | _||   / ") + blue("║"))
	fmt.Println(blue("║") + cyan("      |_| |__|____|_|/_/  |_|  |_|\\___/___/  |_| |___|_|_\\ ") + blue("║"))
	fmt.Println(blue("║") + green("            Task Management Command Line Tool            ") + blue("║"))
	fmt.Println(blue("╚════════════════════════════════════════════════╝"))
	fmt.Println()
}
