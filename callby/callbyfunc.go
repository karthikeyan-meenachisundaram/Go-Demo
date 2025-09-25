/*
package main

import "fmt"

type Employee struct {
	Id         int
	Name       string
	Department string
	Salary     float32
}

func func1(arr [3]Employee) {
	fmt.Println("\nInside function1 array:")

	for i := range arr {
		arr[i].Salary += 2000
	}

	for _, emp := range arr {
		fmt.Printf("Id: %d, Name: %s, Department: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Department, emp.Salary)
	}
}

func func2(slc []Employee) {
	fmt.Println("\nInside function2 slice:")

	for i := range slc {
		slc[i].Salary += 4000
	}

	for _, emp := range slc {
		fmt.Printf("Id: %d, Name: %s, Department: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Department, emp.Salary)
	}
}

func main() {

	emp_array := [3]Employee{
		{Id: 101, Name: "Karthi", Department: "Application Engineering", Salary: 20000},
		{Id: 102, Name: "Ilavarasi", Department: "Data Engineering", Salary: 30000},
		{Id: 103, Name: "Karthi", Department: "Data Engineering", Salary: 40000},
	}

	emp_slice := []Employee{
		{Id: 301, Name: "Srini", Department: "Application Engineering", Salary: 50000},
		{Id: 302, Name: "Jagan", Department: "Data Engineering", Salary: 60000},
		{Id: 303, Name: "Mahesh", Department: "Data Engineering", Salary: 70000},
	}

	func1(emp_array)
	fmt.Println("\nBack in main after function1 array:")
	for _, emp := range emp_array {
		fmt.Printf("Id: %d, Name: %s, Department: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Department, emp.Salary)
	}

	func2(emp_slice)
	fmt.Println("\nBack in main after function2 slice:")
	for _, emp := range emp_slice {
		fmt.Printf("Id: %d, Name: %s, Department: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Department, emp.Salary)
	}
}
*/

/*
package main

import "fmt"

type Employee struct {
	Id     int
	Name   string
	Salary float32
}

// Array passed by value (changes don't affect main)
func updateArray(arr [3]Employee) {
	arr[0].Salary += 5000
	fmt.Println("\nInside updateArray (copy of array):")
	for _, emp := range arr {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}
}

// Slice passed (changes affect main, looks like reference)
func updateSlice(slc []Employee) {
	slc[0].Salary += 5000
	fmt.Println("\nInside updateSlice (same underlying array):")
	for _, emp := range slc {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}
}

// Slice copied using copy() (changes do not affect main)
func updateSliceSafe(slc []Employee) {
	newSlc := make([]Employee, len(slc))
	copy(newSlc, slc) // real independent copy

	newSlc[0].Salary += 5000
	fmt.Println("\nInside updateSliceSafe (copy of slice):")
	for _, emp := range newSlc {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}
}

func main() {
	// Array of employees
	empArray := [3]Employee{
		{Id: 101, Name: "Karthi", Salary: 20000},
		{Id: 102, Name: "Mahesh", Salary: 25000},
		{Id: 103, Name: "Jagan", Salary: 30000},
	}

	// Slice of employees
	empSlice := []Employee{
		{Id: 201, Name: "Srini", Salary: 40000},
		{Id: 202, Name: "Sanjay", Salary: 45000},
		{Id: 203, Name: "Ilavarasi", Salary: 50000},
	}

	// Call array function
	updateArray(empArray)
	fmt.Println("\nBack in main after updateArray (original unchanged):")
	for _, emp := range empArray {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}

	// Call slice function
	updateSlice(empSlice)
	fmt.Println("\nBack in main after updateSlice (original changed):")
	for _, emp := range empSlice {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}

	// Call safe slice function
	updateSliceSafe(empSlice)
	fmt.Println("\nBack in main after updateSliceSafe (original unchanged):")
	for _, emp := range empSlice {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}
}

*/

package main

import "fmt"

type Employee struct {
	Id     int
	Name   string
	Salary float32
}

// 1. Array passed by value (copy only, no effect on main)
func updateArray(arr [3]Employee) {
	arr[0].Salary += 5000
	fmt.Println("\nInside updateArray (copy of array):")
	for _, emp := range arr {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}
}

// 2. Array passed by pointer (can update original array)
func updateArrayPtr(arr *[3]Employee) {
	arr[0].Salary += 5000
	fmt.Println("\nInside updateArrayPtr (pointer to array):")
	for _, emp := range arr {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}
}

// 3. Slice passed (points to same underlying array, so updates main)
func updateSlice(slc []Employee) {
	slc[0].Salary += 5000
	fmt.Println("\nInside updateSlice (same underlying array):")
	for _, emp := range slc {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}
}

// 4. Slice copied with copy() (independent copy, no effect on main)
func updateSliceSafe(slc []Employee) {
	newSlc := make([]Employee, len(slc))
	copy(newSlc, slc) // deep copy of slice

	newSlc[0].Salary += 5000
	fmt.Println("\nInside updateSliceSafe (copy of slice):")
	for _, emp := range newSlc {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}
}

func main() {
	// Array of employees
	empArray := [3]Employee{
		{Id: 101, Name: "Karthi", Salary: 20000},
		{Id: 102, Name: "Mahesh", Salary: 25000},
		{Id: 103, Name: "Jagan", Salary: 30000},
	}

	// Slice of employees
	empSlice := []Employee{
		{Id: 201, Name: "Srini", Salary: 40000},
		{Id: 202, Name: "Sanjay", Salary: 45000},
		{Id: 203, Name: "Ilavarasi", Salary: 50000},
	}

	// 1. Array by value
	updateArray(empArray)
	fmt.Println("\nBack in main after updateArray (original unchanged):")
	for _, emp := range empArray {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}

	// 2. Array by pointer
	updateArrayPtr(&empArray)
	fmt.Println("\nBack in main after updateArrayPtr (original changed):")
	for _, emp := range empArray {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}

	// 3. Slice (reference-like)
	updateSlice(empSlice)
	fmt.Println("\nBack in main after updateSlice (original changed):")
	for _, emp := range empSlice {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}

	// 4. Slice safe copy
	updateSliceSafe(empSlice)
	fmt.Println("\nBack in main after updateSliceSafe (original unchanged):")
	for _, emp := range empSlice {
		fmt.Printf("Id: %d, Name: %s, Salary: %.2f\n", emp.Id, emp.Name, emp.Salary)
	}
}
