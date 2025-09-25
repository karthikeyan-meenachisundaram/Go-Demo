package main

import "fmt"

func gostring() {

	var greeting = "Hello World"
	fmt.Println(greeting)

}

func stringLength(str string) int {

	var length int
	for range str {
		length++
	}
	return length

}

func max(num1, num2 int) int {

	var rt int

	if num1 > num2 {
		rt = num1
	} else {
		rt = num2
	}

	return rt

}

func swap(x, y string) (string, string) {

	return y, x
}

func addTwoNumbers(ac, bc int) int {
	return ac + bc
}
