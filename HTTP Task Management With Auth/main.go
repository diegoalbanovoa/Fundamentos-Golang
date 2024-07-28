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
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// User struct
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

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

// InitializeDB initializes the database with the tasks and users tables
func (tm *TaskManager) InitializeDB() error {
	tasksQuery := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        description TEXT,
        completed BOOLEAN,
        created_at DATETIME
    );
    `
	usersQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE,
        password TEXT
    );
    `
	_, err := tm.db.Exec(tasksQuery)
	if err != nil {
		return err
	}
	_, err = tm.db.Exec(usersQuery)
	return err
}

// AddUser adds a new user
func (tm *TaskManager) AddUser(username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	result, err := tm.db.Exec(query, username, hashedPassword)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return tm.GetUser(int(id))
}

// GetUser gets a user by ID
func (tm *TaskManager) GetUser(id int) (*User, error) {
	query := `SELECT id, username, password FROM users WHERE id = ?`
	row := tm.db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AuthenticateUser authenticates a user by username and password
func (tm *TaskManager) AuthenticateUser(username, password string) (*User, error) {
	query := `SELECT id, username, password FROM users WHERE username = ?`
	row := tm.db.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	return &user, nil
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

// Middleware for authenticating requests
func authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

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

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := tm.AddUser(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, user, http.StatusCreated)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := tm.AuthenticateUser(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, map[string]string{"token": tokenStr}, http.StatusOK)
	})

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

	// Wrap the default mux with the logging and authentication middleware
	loggedMux := loggingMiddleware(authenticateMiddleware(http.DefaultServeMux))

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
