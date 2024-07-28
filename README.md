# Go Programming Projects

This repository contains a series of projects developed as part of a Go programming course. Each project focuses on different aspects of Go and related technologies. The projects are organized by week and cover a wide range of topics from basic syntax to advanced topics like CI/CD and end-to-end testing.

## Project Structure

### Basic Concepts
This project covers the basic concepts of the Go language, including basic syntax, control structures, variables, constants, operators, and loops.

Basic Concepts/
├── README.md
├── .idea/
│ ├── .gitignore
│ ├── Concepto Basicos.iml
│ ├── modules.xml
│ ├── vcs.xml
│ └── workspace.xml
├── constantes/
│ └── main.go
├── for_loop/
│ └── main.go
├── hello_world/
│ └── main.go
├── if_else/
│ └── main.go
├── operadores/
│ └── main.go
└── variables/
└── main.go


**Components:**
- **hello_world/main.go**: Basic "Hello, World" program.
- **variables/main.go**: Examples of variables and data types.
- **constantes/main.go**: Usage of constants in Go.
- **operadores/main.go**: Examples of operators and expressions.
- **if_else/main.go**: Conditional control structures.
- **for_loop/main.go**: Examples of `for` loops.

### HTTP Task Management
This project is a basic HTTP server in Go that allows managing tasks through CRUD operations with a SQLite database.

HTTP Task Management/
├── go.mod
├── main.go
├── README.md
└── .idea/
├── .gitignore
├── HTTP Task Management.iml
├── modules.xml
├── vcs.xml
└── workspace.xml


**Components:**
- **main.go**: Implements a basic HTTP server with CRUD operations for task management using SQLite.
- **go.mod**: Module definition file.

### HTTP Task Management With Auth
This project extends the HTTP Task Management project by adding user authentication and authorization using JWT.

HTTP Task Management With Auth/
├── main.go
├── README.md
└── .idea/
├── .gitignore
├── HTTP Task Management With Auth.iml
├── modules.xml
├── vcs.xml
└── workspace.xml


**Components:**
- **main.go**: Adds JWT-based authentication and authorization to the task management server.

### HTTP Task Management With CI/CD
This project builds on the previous projects by adding continuous integration and deployment (CI/CD) using Docker and Jenkins.

HTTP Task Management With CI CD/
├── Dockerfile
├── go.mod
├── go.sum
├── Jenkinsfile
├── main.go
├── README.md
└── .idea/
├── .gitignore
├── HTTP Task Management With CI CD.iml
├── modules.xml
├── vcs.xml
└── workspace.xml


**Components:**
- **Dockerfile**: Defines the Docker image for the application.
- **Jenkinsfile**: Defines the Jenkins pipeline for CI/CD.
- **main.go**: Task management server with CI/CD setup.

### HTTP Task Management With DB
This project integrates a database into the HTTP Task Management project, allowing for persistent storage of tasks.

HTTP Task Management With DB/
├── go.mod
├── go.sum
├── main.go
├── README.md
└── .idea/
├── .gitignore
├── HTTP Task Management With DB.iml
├── modules.xml
├── vcs.xml
└── workspace.xml


**Components:**
- **main.go**: Integrates SQLite database for persistent task storage.
- **go.mod**: Module definition file.
- **go.sum**: Dependencies file.

### HTTP Task Management With DB Testing
This project adds unit tests for the database-integrated HTTP Task Management project.

HTTP Task Management With DB Testing/
├── go.mod
├── go.sum
├── main.go
├── main_test.go
├── README.md
└── .idea/
├── .gitignore
├── HTTP Task Management With DB Testing.iml
├── modules.xml
├── vcs.xml
└── workspace.xml


**Components:**
- **main.go**: Task management server with SQLite integration.
- **main_test.go**: Unit tests for the task management server.

### HTTP Task Management With E2E Testing
This project adds end-to-end (E2E) testing using Cypress to the HTTP Task Management project.

HTTP Task Management With E2E Testing/
├── cypress.json
├── Dockerfile
├── go.mod
├── go.sum
├── Jenkinsfile
├── main.go
├── README.md
├── .idea/
│ ├── .gitignore
│ ├── HTTP Task Management With E2E Testing.iml
│ ├── modules.xml
│ ├── vcs.xml
│ └── workspace.xml
└── cypress/
└── integration/
└── tasks_spec.js


**Components:**
- **main.go**: Task management server with authentication.
- **Dockerfile**: Docker image for the application.
- **Jenkinsfile**: Jenkins pipeline for CI/CD.
- **cypress.json**: Cypress configuration file.
- **tasks_spec.js**: Cypress end-to-end tests.

### Student Management
This project implements a basic student management system with CRUD operations.

Student Management/
├── main (2).go
├── main.go
├── README (2).md
├── README.md
└── .idea/
├── .gitignore
├── modules.xml
├── Student Management.iml
├── vcs.xml
└── workspace.xml


**Components:**
- **main.go**: Implements CRUD operations for student management.
- **README.md**: Documentation for the project.

### Student Registration
This project implements a student registration system.

Student Registration/
├── main.go
├── README.md
└── .idea/
├── .gitignore
├── modules.xml
├── Student Registration.iml
├── vcs.xml
└── workspace.xml


**Components:**
- **main.go**: Implements student registration functionality.
- **README.md**: Documentation for the project.

### Task Management
This project is a task management system.

Task Management/
├── main.go
├── README.md
└── .idea/
├── .gitignore
├── modules.xml
├── Task Management.iml
├── vcs.xml
└── workspace.xml


**Components:**
- **main.go**: Implements a basic task management system.
- **README.md**: Documentation for the project.
