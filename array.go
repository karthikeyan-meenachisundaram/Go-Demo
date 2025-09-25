package main

import "fmt"

var arr2 = [1]string{"bmw"}

//arr2 = {"bmw"}

func array() {

	arr1 := [5]int{3: 5, 4: 4}
	fmt.Println(arr1)
	fmt.Println(arr2)
	fmt.Println(len(arr1))

	myslice1 := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Printf("myslice1 = %v\n", myslice1)
	fmt.Printf("length = %d\n", len(myslice1))
	fmt.Printf("capacity = %d\n", cap(myslice1))

	myslice1 = append(myslice1, 20, 21)
	fmt.Printf("myslice1 = %v\n", myslice1)
	fmt.Printf("length = %d\n", len(myslice1))
	fmt.Printf("capacity = %d\n", cap(myslice1))
}
