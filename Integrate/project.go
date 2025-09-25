package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Employee struct {
	EmpID  int     `bson:"empid"`
	Name   string  `bson:"name"`
	Salary float32 `bson:"salary"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://Karthikeyan:Hema%401199@cluster0.ivpyric.mongodb.net/"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	collection := client.Database("my_project").Collection("Employee")
	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	var employees []Employee
	if err = cursor.All(ctx, &employees); err != nil {
		log.Fatal(err)
	}
	for _, e := range employees {
		fmt.Printf("Id: %d\n", e.EmpID)
		fmt.Printf("Name: %s\n", e.Name)
		fmt.Printf("Salary: %.2f\n", e.Salary)
		fmt.Println("-------------------------------")
	}
}
