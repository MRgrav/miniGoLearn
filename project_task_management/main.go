package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

const TaskFile = "tasks.json"

func main() {
	for {
		fmt.Println("\n\nWelcome to Project Task Manager")
		fmt.Println("\t1. Add Task")
		fmt.Println("\t2. View Tasks")
		fmt.Println("\t3. Delete Task")
		fmt.Println("\t4. Mark Task as Completed")
		fmt.Println("\t5. Exit")
		fmt.Println("Enter your choice: ")

		reader := bufio.NewReader(os.Stdin)
		ch, _ := reader.ReadString('\n')
		ch = strings.TrimSpace(ch)

		switch ch {
		case "1":
			fmt.Println("Add")
			addTask()
		case "2":
			fmt.Println("View")
			viewTasks()
		case "3":
			fmt.Println("Delete")
			deleteTasks()
		case "4":
			fmt.Println("Complete")
			markTaskAsComplete()
		case "5":
			fmt.Println("exit")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice, please try again.")
		}
		fmt.Println("\n")
	}
}

func loadTask() ([]Task, error) {
	if _, err := os.Stat(TaskFile); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(TaskFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file %w", err)
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks from JSON: %w", err)
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("failed to marshal tasks to JSON: %w", err)
	}
	// try write or catch err
	err = os.WriteFile(TaskFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tasks to file: %w", err)
	}
	return nil
}

func addTask() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter task description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	tasks, err := loadTask()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	//
	tasks = append(tasks, Task{Description: description, Completed: false})

	err = saveTasks(tasks)
	if err != nil {
		fmt.Errorf("Failed %w", err)
	}
}

func viewTasks() {
	tasks, err := loadTask()
	if err != nil {
		fmt.Errorf("Failed loading tasks: %w", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No task found...")
		return
	}

	fmt.Println("\n\n")
	fmt.Printf("%-5s %-30s %s\n", "ID", "Description", "Completed")
	fmt.Println("-------------------------------------------------")
	for i, task := range tasks {
		fmt.Printf("%-5d %-30s %t\n", i+1, task.Description, task.Completed)
	}
}

func deleteTasks() {
	tasks, err := loadTask()
	if err != nil {
		fmt.Errorf("Failed loading tasks: %w", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No Tasks found.")
		return
	}

	viewTasks()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the ID of the task to delete: ")
	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	id, err := strconv.Atoi(option)

	if err != nil || len(tasks) < id || id < 1 {
		fmt.Println("Invalid task id.n\nSelect again...")
		return
	}

	tasks = append(tasks[:id-1], tasks[id:]...)
	err = saveTasks(tasks)
	if err != nil {
		fmt.Errorf("Failed to delete task id: %d, \nerror: %w", id, err)
		return
	}
	fmt.Printf("Task ID : %d deleted successfully.\n", id)

}

func markTaskAsComplete() {
	tasks, err := loadTask()
	if err != nil {
		fmt.Errorf("Failed to load tasks: %w", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No task found")
		return
	}

	viewTasks()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the task id to mark as complete: ")
	option, err := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	id, err := strconv.Atoi(option)

	if err != nil || id < 1 || id > len(tasks) {
		fmt.Println("Invalid Task id.\n")
		return
	}

	tasks[id-1].Completed = true

	err = saveTasks(tasks)
	if err != nil {
		fmt.Errorf("Failed to update task id: %d, \nerror: %w", id, err)
		return
	}
	fmt.Printf("Task ID : %d marked as complete.", id)
}
