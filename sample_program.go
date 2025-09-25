package main

import "fmt"

/*
func main() {
	gostring()

	//type customer struct

	str := "Hello, world!"
	fmt.Println("The given string is: ", str)
	var result int = stringLength(str)
	fmt.Println("Length of the above string is: ", result)

	var a int = 20
	var b int = 30
	var ret = max(a, b)
	fmt.Printf("MAx value is : %d\n", ret)

	var a1, b1 = "Mahesh", "Kumar"
	fmt.Println(swap(a1, b1))

	//var sum = addTwoNumbers
	var rst = addTwoNumbers(5, 15)
	fmt.Print(rst)

}
*/
type Employee struct {
	EmpID  int
	Name   string
	Salary float32
}

func main() {
	employees := []Employee{
		{EmpID: 101, Name: "Karthikeyan", Salary: 12000.00},
		{EmpID: 102, Name: "Mahesh", Salary: 18000.00},
		{EmpID: 103, Name: "Jagan", Salary: 15000.00},
		{EmpID: 104, Name: "Sanjay", Salary: 13000.00},
	}
	for _, e := range employees {
		fmt.Printf("Id: %d\n", e.EmpID)
		fmt.Printf("Name: %s\n", e.Name)
		fmt.Printf("Salary: %.2f\n", e.Salary)
		fmt.Println("-------------------------------")
	}
}

// calling one function send this struct
// calling the function send the array and update value -- display in the main function and called function
