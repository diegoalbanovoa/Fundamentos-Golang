package main

import (
	"fmt"
	"time"
)

// Constants
const adultAge = 18

// Person interface
type Person interface {
	GetName() string
	GetAge() int
	IsAdult() bool
}

// Student struct
type Student struct {
	Name      string
	BirthYear int
	IsActive  bool
}

// GetName returns the name of the student
func (s *Student) GetName() string {
	return s.Name
}

// GetAge calculates and returns the age of the student
func (s *Student) GetAge() int {
	currentYear := time.Now().Year()
	return currentYear - s.BirthYear
}

// IsAdult returns true if the student is an adult
func (s *Student) IsAdult() bool {
	return s.GetAge() >= adultAge
}

// Activate activates the student
func (s *Student) Activate() {
	s.IsActive = true
}

// Deactivate deactivates the student
func (s *Student) Deactivate() {
	s.IsActive = false
}

func main() {
	// Create a slice to hold registered students
	var students []Person

	// Register students
	students = append(students, &Student{Name: "Alice", BirthYear: 2003, IsActive: true})
	students = append(students, &Student{Name: "Bob", BirthYear: 2005, IsActive: true})
	students = append(students, &Student{Name: "Charlie", BirthYear: 1999, IsActive: false})

	// Display registered students
	displayStudents(students)

	// Search for a student by name
	student := findStudentByName(students, "Alice")
	if student != nil {
		fmt.Printf("\nFound student: %s, Age: %d, Active: %t, Adult: %t\n", student.GetName(), student.GetAge(), student.(*Student).IsActive, student.IsAdult())
	} else {
		fmt.Println("\nStudent not found.")
	}

	// Deactivate a student
	fmt.Println("\nDeactivating student Bob...")
	deactivateStudent(students, "Bob")
	displayStudents(students)
}

// Displays information about registered students
func displayStudents(students []Person) {
	fmt.Println("\nRegistered Students:")
	for _, student := range students {
		fmt.Printf("Name: %s, Age: %d, Active: %t, Adult: %t\n", student.GetName(), student.GetAge(), student.(*Student).IsActive, student.IsAdult())
	}
}

// Finds a student by name
func findStudentByName(students []Person, name string) Person {
	for _, student := range students {
		if student.GetName() == name {
			return student
		}
	}
	return nil
}

// Deactivates a student by name
func deactivateStudent(students []Person, name string) {
	student := findStudentByName(students, name)
	if student != nil {
		student.(*Student).Deactivate()
	}
}
