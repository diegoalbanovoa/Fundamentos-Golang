# HTTP Task Management Server with Database and Testing

## Description

This project is a basic HTTP server in Go that allows managing tasks through CRUD operations with a SQLite database. It includes route handling, request and response management, middleware for logging requests, database integration for persistent storage, and testing for both unit and integration.

## Project Structure

The project consists of a single subdirectory:

- **http_task_management_with_db_testing**: Contains a single Go program that implements the task management server with database integration and testing.

## Content

### Main Program
File: `http_task_management_with_db_testing/main.go`

```go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    _ "github.com/mattn/go-sqlite3"
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
    db *sql.DB
}

// NewTaskManager creates a new TaskManager
func NewTaskManager(db *sql.DB) *TaskManager {
    return &TaskManager{db: db}
}

// InitializeDB initializes the database with the tasks table
func (tm *TaskManager) InitializeDB() error {
    query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        description TEXT,
        completed BOOLEAN,
        created_at DATETIME
    );
    `
    _, err := tm.db.Exec(query)
    return err
}

// AddTask adds a new task
func (tm *TaskManager) AddTask(description string) (*Task, error) {
    query := `INSERT INTO tasks (description, completed, created_at) VALUES (?, ?, ?)`
    result, err := tm.db.Exec(query, description, false, time.Now())
    if err != nil {
        return nil, err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }
    return tm.GetTask(int(id))
}

// GetTask gets a task by ID
func (tm *TaskManager) GetTask(id int) (*Task, error) {
    query := `SELECT id, description, completed, created_at FROM tasks WHERE id = ?`
    row := tm.db.QueryRow(query, id)

    var task Task
    err := row.Scan(&task.ID, &task.Description, &task.Completed, &task.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &task, nil
}

// UpdateTask updates a task by ID
func (tm *TaskManager) UpdateTask(id int, description string, completed bool) (*Task, error) {
    query := `UPDATE tasks SET description = ?, completed = ? WHERE id = ?`
    _, err := tm.db.Exec(query, description, completed, id)
    if err != nil {
        return nil, err
    }
    return tm.GetTask(id)
}

// DeleteTask deletes a task by ID
func (tm *TaskManager) DeleteTask(id int) error {
    query := `DELETE FROM tasks WHERE id = ?`
    _, err := tm.db.Exec(query, id)
    return err
}

// ListTasks lists all tasks
func (tm *TaskManager) ListTasks() ([]*Task, error) {
    query := `SELECT id, description, completed, created_at FROM tasks`
    rows, err := tm.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []*Task
    for rows.Next() {
        var task Task
        if err := rows.Scan(&task.ID, &task.Description, &task.Completed, &task.CreatedAt); err != nil {
            return nil, err
        }
        tasks = append(tasks, &task)
    }
    return tasks, nil
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
        next.ServeHTTP(w, r)
    })
}

func main() {
    db, err := sql.Open("sqlite3", "./tasks.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    tm := NewTaskManager(db)
    if err := tm.InitializeDB(); err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            tasks, err := tm.ListTasks()
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            jsonResponse(w, tasks, http.StatusOK)
        case "POST":
            var task Task
            if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            newTask, err := tm.AddTask(task.Description)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
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
            task, err := tm.GetTask(id)
            if err != nil {
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
            updatedTask, err := tm.UpdateTask(id, task.Description, task.Completed)
            if err != nil {
                http.Error(w, "Task not found", http.StatusNotFound)
                return
            }
            jsonResponse(w, updatedTask, http.StatusOK)
        case "DELETE":
            if err := tm.DeleteTask(id); err != nil {
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

