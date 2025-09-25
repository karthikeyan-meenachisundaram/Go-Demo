package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	client, ctx, cancel, err := Connect(10) // just call Connect, no URI here
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("disconnect error:", err)
		}
	}()

	coll := client.Database("my_db").Collection("Employee_Details")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose action: 1) Fetch  2) Insert")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	choice, err := strconv.Atoi(line)
	if err != nil {
		log.Fatal("Invalid choice:", err)
	}

	switch choice {
	case 1:
		if err := FetchEmployees(coll); err != nil {
			log.Fatal("Fetch error:", err)
		}
	case 2:
		for {
			fmt.Print("Enter EmpID (integer): ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			eid, err := strconv.Atoi(line)
			if err != nil {
				fmt.Println("Invalid EmpID. Please enter an integer.")
				continue
			}

			fmt.Print("Enter Name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			if name == "" {
				fmt.Println("Name cannot be empty.")
				continue
			}

			fmt.Print("Enter Salary (number): ")
			line, _ = reader.ReadString('\n')
			line = strings.TrimSpace(line)
			sal, err := strconv.ParseFloat(line, 32)
			if err != nil {
				fmt.Println("Invalid salary. Please enter a number.")
				continue
			}

			emp := Employee{EmpID: eid, Name: name, Salary: float32(sal)}
			if err := InsertEmployee(coll, emp); err != nil {
				log.Println("Insert error:", err)
			} else {
				fmt.Println("Inserted successfully.")
			}

			fmt.Print("Insert another? (y/n): ")
			yn, _ := reader.ReadString('\n')
			yn = strings.TrimSpace(strings.ToLower(yn))
			if yn != "y" && yn != "yes" {
				fmt.Println("Stopping inserts.")
				break
			}
		}
	default:
		fmt.Println("Invalid choice")
	}
}
