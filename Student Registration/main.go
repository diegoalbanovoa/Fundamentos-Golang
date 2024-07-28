package main

import (
	"fmt"
	"time"
)

// Constants
const adultAge = 18

// Student struct
type Student struct {
	Name      string
	BirthYear int
	IsActive  bool
}

func main() {
	// Create a slice to hold registered students
	var students []Student

	// Register students
	students = append(students, registerStudent("Alice", 2003, true))
	students = append(students, registerStudent("Bob", 2005, true))
	students = append(students, registerStudent("Charlie", 1999, false))

	// Display registered students
	displayStudents(students)
}

// Registers a new student
func registerStudent(name string, birthYear int, isActive bool) Student {
	return Student{Name: name, BirthYear: birthYear, IsActive: isActive}
}

// Calculates the age of the student
func calculateAge(birthYear int) int {
	currentYear := time.Now().Year()
	return currentYear - birthYear
}

// Checks if the student is an adult
func isAdult(age int) bool {
	return age >= adultAge
}

// Displays information about registered students
func displayStudents(students []Student) {
	for _, student := range students {
		age := calculateAge(student.BirthYear)
		fmt.Printf("Name: %s, Age: %d, Active: %t, Adult: %t\n", student.Name, age, student.IsActive, isAdult(age))
	}
}
