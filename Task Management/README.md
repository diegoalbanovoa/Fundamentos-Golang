# Task Management System

## Description

This project is a task management system that allows you to add tasks, mark tasks as completed, list pending and completed tasks, and process tasks concurrently using goroutines and channels. It also demonstrates error handling in Go.

## Project Structure

The project consists of a single subdirectory:

- **task_management**: Contains a single Go program that implements the task management system.

## Content

### Main Program
File: `task_management/main.go`

```go
package main

import (
    "errors"
    "fmt"
    "sync"
    "time"
)

// Task struct
type Task struct {
    ID          int
    Description string
    Completed   bool
    CreatedAt   time.Time
}

// TaskManager struct
type TaskManager struct {
    tasks map[int]*Task
    mu    sync.Mutex
}

// NewTaskManager creates a new TaskManager
func NewTaskManager() *TaskManager {
    return &TaskManager{
        tasks: make(map[int]*Task),
    }
}

// AddTask adds a new task
func (tm *TaskManager) AddTask(description string) int {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    id := len(tm.tasks) + 1
    tm.tasks[id] = &Task{
        ID:          id,
        Description: description,
        Completed:   false,
        CreatedAt:   time.Now(),
    }

    return id
}

// CompleteTask marks a task as completed
func (tm *TaskManager) CompleteTask(id int) error {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    task, exists := tm.tasks[id]
    if !exists {
        return errors.New("task not found")
    }

    task.Completed = true
    return nil
}

// ListTasks lists all tasks
func (tm *TaskManager) ListTasks(completed bool) []*Task {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    var tasks []*Task
    for _, task := range tm.tasks {
        if task.Completed == completed {
            tasks = append(tasks, task)
        }
    }

    return tasks
}

// ProcessTasks processes tasks concurrently
func (tm *TaskManager) ProcessTasks(ids []int) {
    var wg sync.WaitGroup
    taskChan := make(chan int)

    for i := 0; i < 3; i++ {
        wg.Add(1)
        go tm.worker(&wg, taskChan)
    }

    for _, id := range ids {
        taskChan <- id
    }
    close(taskChan)
    wg.Wait()
}

func (tm *TaskManager) worker(wg *sync.WaitGroup, tasks chan int) {
    defer wg.Done()

    for id := range tasks {
        fmt.Printf("Processing task %d\n", id)
        time.Sleep(1 * time.Second) // Simulate task processing
        tm.CompleteTask(id)
        fmt.Printf("Task %d completed\n", id)
    }
}

func main() {
    tm := NewTaskManager()

    // Add tasks
    tm.AddTask("Learn Go")
    tm.AddTask("Read a book")
    tm.AddTask("Write a blog post")

    // Process tasks concurrently
    fmt.Println("Processing tasks...")
    tm.ProcessTasks([]int{1, 2, 3})

    // List completed tasks
    fmt.Println("Completed tasks:")
    for _, task := range tm.ListTasks(true) {
        fmt.Printf("ID: %d, Description: %s, Completed: %t, CreatedAt: %s\n",
            task.ID, task.Description, task.Completed, task.CreatedAt)
    }
}

