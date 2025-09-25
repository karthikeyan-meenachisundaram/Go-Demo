package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Employee struct (same as before)
type Employee struct {
	EmpID  int     `bson:"empid"`
	Name   string  `bson:"name"`
	Salary float32 `bson:"salary"`
}

// FetchEmployees reads all docs using a fresh context for the operation.
func FetchEmployees(coll *mongo.Collection) error {
	// per-operation timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("find error: %w", err)
	}
	defer cursor.Close(ctx)

	var employees []Employee
	if err := cursor.All(ctx, &employees); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	if len(employees) == 0 {
		fmt.Println("No employees found.")
		return nil
	}

	for _, e := range employees {
		fmt.Printf("EmpID: %d, Name: %s, Salary: %.2f\n", e.EmpID, e.Name, e.Salary)
	}
	return nil
}

// InsertEmployee inserts one Employee using a fresh context for the operation.
func InsertEmployee(coll *mongo.Collection, emp Employee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := coll.InsertOne(ctx, emp)
	if err != nil {
		return fmt.Errorf("insert error: %w", err)
	}
	fmt.Println("Inserted Employee with ID:", res.InsertedID)
	return nil
}
