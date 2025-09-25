package main

import "fmt"

var globalVar = "I am global" // ✅ works with var
// msg := "won't work"        // ❌ not allowed outside func

func hello() {

	dayOfWeek := 4

	switch dayOfWeek {
	case 1:
		fmt.Println("Monday")
	case 2:
		fmt.Println("Tuesday")
	case 3:
		fmt.Println("Wednesday")
	case 4:
		fmt.Println("Thursday")
	case 5:
		fmt.Println("Friday")
	case 6:
		fmt.Println("Saturday")
	case 7:
		fmt.Println("Sunday")
	default:
		fmt.Println("Invalid day of the week")
	}

}
