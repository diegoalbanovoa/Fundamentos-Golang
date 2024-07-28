package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	tm := NewTaskManager(db)
	if err := tm.InitializeDB(); err != nil {
		t.Fatal(err)
	}

	return db
}

func TestAddTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tm := NewTaskManager(db)
	task, err := tm.AddTask("Test Task")
	if err != nil {
		t.Fatal(err)
	}

	if task.Description != "Test Task" {
		t.Errorf("expected task description to be 'Test Task', got '%s'", task.Description)
	}
}

func TestGetTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tm := NewTaskManager(db)
	newTask, err := tm.AddTask("Test Task")
	if err != nil {
		t.Fatal(err)
	}

	task, err := tm.GetTask(newTask.ID)
	if err != nil {
		t.Fatal(err)
	}

	if task.Description != "Test Task" {
		t.Errorf("expected task description to be 'Test Task', got '%s'", task.Description)
	}
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tm := NewTaskManager(db)
	newTask, err := tm.AddTask("Test Task")
	if err != nil {
		t.Fatal(err)
	}

	updatedTask, err := tm.UpdateTask(newTask.ID, "Updated Task", true)
	if err != nil {
		t.Fatal(err)
	}

	if updatedTask.Description != "Updated Task" {
		t.Errorf("expected task description to be 'Updated Task', got '%s'", updatedTask.Description)
	}

	if !updatedTask.Completed {
		t.Errorf("expected task to be completed, got '%t'", updatedTask.Completed)
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tm := NewTaskManager(db)
	newTask, err := tm.AddTask("Test Task")
	if err != nil {
		t.Fatal(err)
	}

	err = tm.DeleteTask(newTask.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tm.GetTask(newTask.ID)
	if err == nil {
		t.Fatal("expected error when getting deleted task, got nil")
	}
}

func TestListTasks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tm := NewTaskManager(db)
	_, err := tm.AddTask("Task 1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = tm.AddTask("Task 2")
	if err != nil {
		t.Fatal(err)
	}

	tasks, err := tm.ListTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestHTTPHandlers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tm := NewTaskManager(db)
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

	t.Run("Add Task", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(`{"description":"Test Task"}`))
		res := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(res, req)

		if res.Code != http.StatusCreated {
			t.Errorf("expected status 201 Created, got %d", res.Code)
		}
	})

	t.Run("Get Task", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tasks/1", nil)
		res := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %d", res.Code)
		}
	})

	t.Run("Update Task", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"description":"Updated Task","completed":true}`))
		res := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %d", res.Code)
		}
	})

	t.Run("Delete Task", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
		res := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(res, req)

		if res.Code != http.StatusNoContent {
			t.Errorf("expected status 204 No Content, got %d", res.Code)
		}
	})

	t.Run("List Tasks", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tasks", nil)
		res := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %d", res.Code)
		}
	})
}
