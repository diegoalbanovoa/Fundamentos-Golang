# HTTP Task Management Server

## Description

This project is a basic HTTP server in Go that allows managing tasks through CRUD operations. It includes route handling, request and response management, and middleware for logging requests.

## Project Structure

The project consists of a single subdirectory:

- **http_task_management**: Contains a single Go program that implements the task management server.

## Content

### Main Program
File: `http_task_management/main.go`

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync"
    "time"
)

// Task struct
type Task struct {
    ID          int       `json:"id"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
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
func (tm *TaskManager) AddTask(description string) *Task {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    id := len(tm.tasks) + 1
    task := &Task{
        ID:          id,
        Description: description,
        Completed:   false,
        CreatedAt:   time.Now(),
    }
    tm.tasks[id] = task

    return task
}

// GetTask gets a task by ID
func (tm *TaskManager) GetTask(id int) (*Task, bool) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    task, exists := tm.tasks[id]
    return task, exists
}

// UpdateTask updates a task by ID
func (tm *TaskManager) UpdateTask(id int, description string, completed bool) (*Task, bool) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    task, exists := tm.tasks[id]
    if exists {
        task.Description = description
        task.Completed = completed
    }
    return task, exists
}

// DeleteTask deletes a task by ID
func (tm *TaskManager) DeleteTask(id int) bool {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    _, exists := tm.tasks[id]
    if exists {
        delete(tm.tasks, id)
    }
    return exists
}

// ListTasks lists all tasks
func (tm *TaskManager) ListTasks() []*Task {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    var tasks []*Task
    for _, task := range tm.tasks {
        tasks = append(tasks, task)
    }
    return tasks
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
        next.ServeHTTP(w, r)
    })
}

func main() {
    tm := NewTaskManager()

    http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            tasks := tm.ListTasks()
            jsonResponse(w, tasks, http.StatusOK)
        case "POST":
            var task Task
            if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            newTask := tm.AddTask(task.Description)
            jsonResponse(w, newTask, http.StatusCreated)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
        idStr := r.URL.Path[len("/tasks/"):]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }

        switch r.Method {
        case "GET":
            task, exists := tm.GetTask(id)
            if !exists {
                http.Error(w, "Task not found", http.StatusNotFound)
                return
            }
            jsonResponse(w, task, http.StatusOK)
        case "PUT":
            var task Task
            if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            updatedTask, exists := tm.UpdateTask(id, task.Description, task.Completed)
            if !exists {
                http.Error(w, "Task not found", http.StatusNotFound)
                return
            }
            jsonResponse(w, updatedTask, http.StatusOK)
        case "DELETE":
            if !tm.DeleteTask(id) {
                http.Error(w, "Task not found", http.StatusNotFound)
                return
            }
            w.WriteHeader(http.StatusNoContent)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // Wrap the default mux with the logging middleware
    loggedMux := loggingMiddleware(http.DefaultServeMux)

    fmt.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", loggedMux); err != nil {
        log.Fatalf("could not start server: %v\n", err)
    }
}

// jsonResponse encodes response as JSON and writes it to the ResponseWriter
func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

