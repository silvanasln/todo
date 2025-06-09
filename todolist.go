package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID   int    `bson:"id" json:"id"`
	Task string `bson:"task" json:"task"`
	Done bool   `bson:"done" json:"done"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Starting the program...")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	defer client.Disconnect(ctx)

	fmt.Println("Pinging MongoDB...")
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}
	fmt.Println("Successfully connected to MongoDB!")

	collection = client.Database("tododb").Collection("todos")

	e := echo.New()

	e.GET("/todos", getTodos)
	e.POST("/todos", createTodo)
	e.PUT("/todos/:id", updateTodo)

	// اجرای سرور
	fmt.Println("Starting Echo server on :1323...")
	e.Start(":1323")
}

func getTodos(c echo.Context) error {
	ctx := context.Background()
	fmt.Println("Fetching todos from MongoDB...")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error finding todos:", err)
		return err
	}
	defer cursor.Close(ctx)

	var todos []Todo
	if err = cursor.All(ctx, &todos); err != nil {
		fmt.Println("Error decoding todos:", err)
		return err
	}
	return c.JSON(http.StatusOK, todos)
}

func createTodo(c echo.Context) error {
	ctx := context.Background()
	fmt.Println("Creating a new todo...")
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		fmt.Println("Error binding todo:", err)
		return err
	}

	_, err := collection.InsertOne(ctx, todo)
	if err != nil {
		fmt.Println("Error inserting todo:", err)
		return err
	}
	return c.JSON(http.StatusCreated, todo)
}

func updateTodo(c echo.Context) error {
	ctx := context.Background()
	fmt.Println("Updating a todo...")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println("Error converting id to integer:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var updatedTodo Todo
	if err := c.Bind(&updatedTodo); err != nil {
		fmt.Println("Error binding updated todo:", err)
		return err
	}

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"task": updatedTodo.Task, "done": updatedTodo.Done}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("Error updating todo:", err)
		return err
	}

	if result.MatchedCount == 0 {
		fmt.Println("No todo found with id:", id)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Todo not found"})
	}

	var todo Todo
	err = collection.FindOne(ctx, filter).Decode(&todo)
	if err != nil {
		fmt.Println("Error fetching updated todo:", err)
		return err
	}

	return c.JSON(http.StatusOK, todo)
}

