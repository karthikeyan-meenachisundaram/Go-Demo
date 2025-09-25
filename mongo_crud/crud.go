// main.go
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ---- Models ----
type Employee struct {
	EmpID int    `bson:"emp_id"`
	Name  string `bson:"emp_name"`
}

type Department struct {
	DepartmentName string `bson:"department_name"`
	EmpID          int    `bson:"emp_id"`
}

type Developer struct {
	Language string `bson:"language"`
	EmpID    int    `bson:"emp_id"`
}

type Tester struct {
	Language string `bson:"language_"`
	EmpID    int    `bson:"emp_id"`
}

type EmployeeFull struct {
	EmpID        int      `bson:"emp_id"`
	EmpName      string   `bson:"emp_name"`
	Departments  []string `bson:"departments"`
	DevLanguages []string `bson:"dev_languages"`
	TestLangs    []string `bson:"test_languages"`
}

// ---- Hardcoded URI (replace before running) ----
const mongoURI = "mongodb+srv://Karthikeyan:Hema%401199@mycluster.5oolqvy.mongodb.net/"

// ---- Helpers ----
func ctx(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}

func connect(uri string, timeoutSeconds int) (*mongo.Client, context.Context, context.CancelFunc, error) {
	cctx, cancel := ctx(time.Duration(timeoutSeconds) * time.Second)
	client, err := mongo.Connect(cctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}
	if err := client.Ping(cctx, nil); err != nil {
		_ = client.Disconnect(cctx)
		cancel()
		return nil, nil, nil, err
	}
	return client, cctx, cancel, nil
}

// Generic small helpers
func readLine(r *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	s, _ := r.ReadString('\n')
	return strings.TrimSpace(s)
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// ---- CRUD operations ----
func insertOne(coll *mongo.Collection, doc interface{}) error {
	c, cancel := ctx(5 * time.Second)
	defer cancel()
	_, err := coll.InsertOne(c, doc)
	return err
}

func findEmployee(coll *mongo.Collection, empID int) (*Employee, error) {
	c, cancel := ctx(5 * time.Second)
	defer cancel()
	var e Employee
	err := coll.FindOne(c, bson.M{"emp_id": empID}).Decode(&e)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &e, err
}

func findAll(coll *mongo.Collection) ([]bson.M, error) {
	c, cancel := ctx(10 * time.Second)
	defer cancel()
	cursor, err := coll.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)
	var out []bson.M
	if err := cursor.All(c, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func updateEmployeeName(coll *mongo.Collection, empID int, newName string) error {
	c, cancel := ctx(5 * time.Second)
	defer cancel()
	res, err := coll.UpdateOne(c, bson.M{"emp_id": empID}, bson.M{"$set": bson.M{"emp_name": newName}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no employee found with emp_id %d", empID)
	}
	return nil
}

func updateDepartmentName(coll *mongo.Collection, empID int, newDept string) error {
	c, cancel := ctx(5 * time.Second)
	defer cancel()
	res, err := coll.UpdateMany(c, bson.M{"emp_id": empID}, bson.M{"$set": bson.M{"department_name": newDept}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no department docs found for emp_id %d", empID)
	}
	return nil
}

func updateDeveloperLanguage(coll *mongo.Collection, empID int, newLang string) error {
	c, cancel := ctx(5 * time.Second)
	defer cancel()
	res, err := coll.UpdateMany(c, bson.M{"emp_id": empID}, bson.M{"$set": bson.M{"language": newLang}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no developer docs found for emp_id %d", empID)
	}
	return nil
}

func updateTesterLanguage(coll *mongo.Collection, empID int, newLang string) error {
	c, cancel := ctx(5 * time.Second)
	defer cancel()
	res, err := coll.UpdateMany(c, bson.M{"emp_id": empID}, bson.M{"$set": bson.M{"language_": newLang}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no tester docs found for emp_id %d", empID)
	}
	return nil
}

func deleteByEmp(coll *mongo.Collection, empID int) error {
	c, cancel := ctx(5 * time.Second)
	defer cancel()
	_, err := coll.DeleteMany(c, bson.M{"emp_id": empID})
	return err
}

// ---- Aggregation join by emp_id ----
func getEmployeeFull(db *mongo.Database, empID int) (*EmployeeFull, error) {
	c, cancel := ctx(10 * time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "emp_id", Value: empID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Department"},
			{Key: "localField", Value: "emp_id"},
			{Key: "foreignField", Value: "emp_id"},
			{Key: "as", Value: "departments_docs"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Developers"},
			{Key: "localField", Value: "emp_id"},
			{Key: "foreignField", Value: "emp_id"},
			{Key: "as", Value: "developers_docs"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Testers"},
			{Key: "localField", Value: "emp_id"},
			{Key: "foreignField", Value: "emp_id"},
			{Key: "as", Value: "testers_docs"},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "emp_id", Value: 1},
			{Key: "emp_name", Value: 1},
			{Key: "departments", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: "$departments_docs"},
				{Key: "as", Value: "d"},
				{Key: "in", Value: "$$d.department_name"},
			}}}},
			{Key: "dev_languages", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: "$developers_docs"},
				{Key: "as", Value: "dev"},
				{Key: "in", Value: "$$dev.language"},
			}}}},
			{Key: "test_languages", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: "$testers_docs"},
				{Key: "as", Value: "t"},
				{Key: "in", Value: "$$t.language_"},
			}}}},
		}}},
	}

	cur, err := db.Collection("Employee").Aggregate(c, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(c)
	if !cur.Next(c) {
		return nil, nil
	}
	var out EmployeeFull
	if err := cur.Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func printMenu() {
	fmt.Println()
	fmt.Println("Select operation:")
	fmt.Println("1) Create")
	fmt.Println("2) Read")
	fmt.Println("3) Update")
	fmt.Println("4) Delete")
	fmt.Println("5) Join (get full employee info by emp_id)")
	fmt.Println("6) Exit")
	fmt.Print("> ")
}

// ensureIndexes creates a unique index on Employee.emp_id
func ensureIndexes(empColl *mongo.Collection) error {
	c, cancel := ctx(10 * time.Second)
	defer cancel()

	idxModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "emp_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := empColl.Indexes().CreateOne(c, idxModel)
	return err
}

// insertEmployee inserts an Employee and returns a friendly message if emp_id already exists.
func insertEmployee(empColl *mongo.Collection, e Employee) error {
	c, cancel := ctx(5 * time.Second)
	defer cancel()

	_, err := empColl.InsertOne(c, e)
	if err == nil {
		return nil
	}

	// detect duplicate-key (11000)
	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, w := range we.WriteErrors {
			if w.Code == 11000 {
				return fmt.Errorf("employee with emp_id %d already exists", e.EmpID)
			}
		}
	}

	// Some drivers might return CommandError
	var ce mongo.CommandError
	if errors.As(err, &ce) && ce.Code == 11000 {
		return fmt.Errorf("employee with emp_id %d already exists", e.EmpID)
	}

	return err
}

// ---- Interactive terminal menu ----
func main() {
	// Connect to MongoDB
	client, connCtx, cancel, err := connect(mongoURI, 10)
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer cancel()
	defer func() {
		if err := client.Disconnect(connCtx); err != nil {
			log.Printf("disconnect error: %v", err)
		}
	}()

	db := client.Database("my_db")
	empColl := db.Collection("Employee")
	// create unique index on emp_id (safe to call at startup)
	if err := ensureIndexes(empColl); err != nil {
		log.Fatalf("failed to create index on Employee.emp_id: %v", err)
	}

	deptColl := db.Collection("Department")
	devColl := db.Collection("Developers")
	testColl := db.Collection("Testers")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Connected to MongoDB")

	for {
		printMenu()
		line, _ := reader.ReadString('\n') // read user input
		choice := strings.TrimSpace(line)  // clean it up

		switch choice {
		case "1":
			// Create
			fmt.Println("Choose collection to insert into:")
			fmt.Println("a) Employee")
			fmt.Println("b) Department")
			fmt.Println("c) Developers")
			fmt.Println("d) Testers")
			col := readLine(reader, "> ")
			switch strings.ToLower(col) {
			case "a", "employee":
				idS := readLine(reader, "EmpID (int): ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				name := readLine(reader, "Emp Name: ")
				if err := insertEmployee(empColl, Employee{EmpID: id, Name: name}); err != nil {
					fmt.Println("insert error:", err)
				} else {
					fmt.Println("employee inserted")
				}

			case "b", "department":
				idS := readLine(reader, "EmpID (int): ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				dept := readLine(reader, "Department Name: ")
				if err := insertOne(deptColl, Department{DepartmentName: dept, EmpID: id}); err != nil {
					fmt.Println("insert error:", err)
				} else {
					fmt.Println("department inserted")
				}
			case "c", "developers":
				idS := readLine(reader, "EmpID (int): ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				lang := readLine(reader, "Developer Language: ")
				if err := insertOne(devColl, Developer{Language: lang, EmpID: id}); err != nil {
					fmt.Println("insert error:", err)
				} else {
					fmt.Println("developer inserted")
				}
			case "d", "testers":
				idS := readLine(reader, "EmpID (int): ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				lang := readLine(reader, "Tester Language: ")
				if err := insertOne(testColl, Tester{Language: lang, EmpID: id}); err != nil {
					fmt.Println("insert error:", err)
				} else {
					fmt.Println("tester inserted")
				}
			default:
				fmt.Println("unknown collection")
			}

		case "2":
			// Read
			fmt.Println("Choose collection to read from:")
			fmt.Println("a) Employee")
			fmt.Println("b) Department")
			fmt.Println("c) Developers")
			fmt.Println("d) Testers")
			fmt.Println("e) Read joined one employee (same as option 5)")
			col := readLine(reader, "> ")
			switch strings.ToLower(col) {
			case "a", "employee":
				sub := readLine(reader, "Enter emp_id (or blank to list all): ")
				if sub == "" {
					all, err := findAll(empColl)
					if err != nil {
						fmt.Println("read error:", err)
						continue
					}
					fmt.Println("Employees:", all)
				} else {
					id, err := parseInt(sub)
					if err != nil {
						fmt.Println("invalid emp id")
						continue
					}
					e, err := findEmployee(empColl, id)
					if err != nil {
						fmt.Println("read error:", err)
						continue
					}
					if e == nil {
						fmt.Println("no employee found")
					} else {
						fmt.Printf("Employee: %+v\n", *e)
					}
				}
			case "b", "department":
				sub := readLine(reader, "Enter emp_id (or blank to list all): ")
				if sub == "" {
					all, err := findAll(deptColl)
					if err != nil {
						fmt.Println("read error:", err)
						continue
					}
					fmt.Println("Departments:", all)
				} else {
					id, err := parseInt(sub)
					if err != nil {
						fmt.Println("invalid emp id")
						continue
					}
					c, cancel := ctx(5 * time.Second)
					var docs []bson.M
					cur, err := deptColl.Find(c, bson.M{"emp_id": id})
					if err != nil {
						cancel()
						fmt.Println("read error:", err)
						continue
					}
					if err := cur.All(c, &docs); err != nil {
						cancel()
						fmt.Println("cursor error:", err)
						continue
					}
					cancel()
					fmt.Println("Department docs:", docs)
				}
			case "c", "developers":
				sub := readLine(reader, "Enter emp_id (or blank to list all): ")
				if sub == "" {
					all, err := findAll(devColl)
					if err != nil {
						fmt.Println("read error:", err)
						continue
					}
					fmt.Println("Developers:", all)
				} else {
					id, err := parseInt(sub)
					if err != nil {
						fmt.Println("invalid emp id")
						continue
					}
					c, cancel := ctx(5 * time.Second)
					var docs []bson.M
					cur, err := devColl.Find(c, bson.M{"emp_id": id})
					if err != nil {
						cancel()
						fmt.Println("read error:", err)
						continue
					}
					if err := cur.All(c, &docs); err != nil {
						cancel()
						fmt.Println("cursor error:", err)
						continue
					}
					cancel()
					fmt.Println("Developer docs:", docs)
				}
			case "d", "testers":
				sub := readLine(reader, "Enter emp_id (or blank to list all): ")
				if sub == "" {
					all, err := findAll(testColl)
					if err != nil {
						fmt.Println("read error:", err)
						continue
					}
					fmt.Println("Testers:", all)
				} else {
					id, err := parseInt(sub)
					if err != nil {
						fmt.Println("invalid emp id")
						continue
					}
					c, cancel := ctx(5 * time.Second)
					var docs []bson.M
					cur, err := testColl.Find(c, bson.M{"emp_id": id})
					if err != nil {
						cancel()
						fmt.Println("read error:", err)
						continue
					}
					if err := cur.All(c, &docs); err != nil {
						cancel()
						fmt.Println("cursor error:", err)
						continue
					}
					cancel()
					fmt.Println("Tester docs:", docs)
				}
			case "e":
				// same as join below
				idS := readLine(reader, "Enter emp_id to join: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				full, err := getEmployeeFull(db, id)
				if err != nil {
					fmt.Println("join error:", err)
					continue
				}
				if full == nil {
					fmt.Println("no employee found")
				} else {
					fmt.Printf("Joined Employee Info: %+v\n", *full)
				}
			default:
				fmt.Println("unknown collection")
			}

		case "3":
			// Update
			fmt.Println("Choose collection to update:")
			fmt.Println("a) Employee (update name by emp_id)")
			fmt.Println("b) Department (update department_name by emp_id)")
			fmt.Println("c) Developers (update language by emp_id)")
			fmt.Println("d) Testers (update language_ by emp_id)")
			col := readLine(reader, "> ")
			switch strings.ToLower(col) {
			case "a", "employee":
				idS := readLine(reader, "EmpID: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				name := readLine(reader, "New Name: ")
				if err := updateEmployeeName(empColl, id, name); err != nil {
					fmt.Println("update error:", err)
				} else {
					fmt.Println("employee name updated")
				}
			case "b", "department":
				idS := readLine(reader, "EmpID: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				newDept := readLine(reader, "New Department Name: ")
				if err := updateDepartmentName(deptColl, id, newDept); err != nil {
					fmt.Println("update error:", err)
				} else {
					fmt.Println("department docs updated")
				}
			case "c", "developers":
				idS := readLine(reader, "EmpID: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				newLang := readLine(reader, "New Developer Language: ")
				if err := updateDeveloperLanguage(devColl, id, newLang); err != nil {
					fmt.Println("update error:", err)
				} else {
					fmt.Println("developer docs updated")
				}
			case "d", "testers":
				idS := readLine(reader, "EmpID: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				newLang := readLine(reader, "New Tester Language: ")
				if err := updateTesterLanguage(testColl, id, newLang); err != nil {
					fmt.Println("update error:", err)
				} else {
					fmt.Println("tester docs updated")
				}
			default:
				fmt.Println("unknown collection")
			}

		case "4":
			// Delete
			fmt.Println("Choose collection to delete from (or delete all related docs for an emp):")
			fmt.Println("a) Employee")
			fmt.Println("b) Department")
			fmt.Println("c) Developers")
			fmt.Println("d) Testers")
			fmt.Println("e) Delete employee and all related documents (Employee + Dept + Dev + Tester)")
			col := readLine(reader, "> ")
			switch strings.ToLower(col) {
			case "a", "employee":
				idS := readLine(reader, "EmpID to delete: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				if err := deleteByEmp(empColl, id); err != nil {
					fmt.Println("delete error:", err)
				} else {
					fmt.Println("employee deleted")
				}
			case "b", "department":
				idS := readLine(reader, "EmpID to delete dept for: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				if err := deleteByEmp(deptColl, id); err != nil {
					fmt.Println("delete error:", err)
				} else {
					fmt.Println("department docs deleted")
				}
			case "c", "developers":
				idS := readLine(reader, "EmpID to delete developers for: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				if err := deleteByEmp(devColl, id); err != nil {
					fmt.Println("delete error:", err)
				} else {
					fmt.Println("developer docs deleted")
				}
			case "d", "testers":
				idS := readLine(reader, "EmpID to delete testers for: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				if err := deleteByEmp(testColl, id); err != nil {
					fmt.Println("delete error:", err)
				} else {
					fmt.Println("tester docs deleted")
				}
			case "e":
				idS := readLine(reader, "EmpID to delete all related docs for: ")
				id, err := parseInt(idS)
				if err != nil {
					fmt.Println("invalid emp id")
					continue
				}
				_ = deleteByEmp(empColl, id)
				_ = deleteByEmp(deptColl, id)
				_ = deleteByEmp(devColl, id)
				_ = deleteByEmp(testColl, id)
				fmt.Println("employee and related docs deleted")
			default:
				fmt.Println("unknown option")
			}

		case "5":
			// Join
			idS := readLine(reader, "Enter emp_id to retrieve joined info: ")
			id, err := parseInt(idS)
			if err != nil {
				fmt.Println("invalid emp id")
				continue
			}
			full, err := getEmployeeFull(db, id)
			if err != nil {
				fmt.Println("join error:", err)
				continue
			}
			if full == nil {
				fmt.Println("no employee found")
			} else {
				fmt.Printf("Employee joined info: %+v\n", *full)
			}

		case "6":
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("unknown choice")
		}
	}
}
