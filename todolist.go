package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Completed bool               `bson:"completed" json:"completed"`
}

var collection *mongo.Collection

func main() {
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("خطا در اتصال به مونگو:", err)
	}
	defer client.Disconnect(context.Background())

	collection = client.Database("todoDB").Collection("tasks")

 
	count, _ := collection.CountDocuments(context.Background(), bson.M{})
	if count == 0 {
		todos := []interface{}{
			Todo{Title: "خرید کردن", Completed: false},
			Todo{Title: "کتاب خوندن", Completed: false},
			Todo{Title: "فیلم دیدن", Completed: false},
			Todo{Title: "پختن غذا", Completed: false},
		}
		_, err := collection.InsertMany(context.Background(), todos)
		if err != nil {
			log.Println("نمی‌تونم لیست اولیه رو ذخیره کنم! خطا:", err)
		} else {
			fmt.Println("لیست اولیه ثبت شد.")
		}
	}


 e := echo.New()

	e.GET("/todos", getTodos)
	e.POST("/todos", createTodo)
	e.PUT("/todos/:id", updateTodo)
	e.DELETE("/todos/:id", deleteTodo)

	fmt.Println("وب‌سایت من روشنه! برو به پورت 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}


func getTodos(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "خطا در خواندن کارها")
	}
	defer cursor.Close(ctx)

	var todos []Todo
	for cursor.Next(ctx) {
		var t Todo
		if err := cursor.Decode(&t); err != nil {
			return c.JSON(http.StatusInternalServerError, "خطا در دیکد کردن داده")
		}
		todos = append(todos, t)
	}
	return c.JSON(http.StatusOK, todos)
}



func createTodo(c echo.Context) error {
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, "ورودی نامعتبر")
	}
	todo.Completed = false
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, todo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "خطا در ذخیره کار")
	}
	todo.ID = res.InsertedID.(primitive.ObjectID)
	return c.JSON(http.StatusCreated, todo)
}


func updateTodo(c echo.Context) error {
	idParam := c.Param("id")
	todoID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "شناسه نامعتبر")
	}
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, "ورودی نامعتبر")
	}
	update := bson.M{
		"$set": bson.M{
			"title":     todo.Title,
			"completed": todo.Completed,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = collection.UpdateOne(ctx, bson.M{"_id": todoID}, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "خطا در ویرایش کار")
	}
	return c.JSON(http.StatusOK, "کار به‌روزرسانی شد")
}

func deleteTodo(c echo.Context) error {
	idParam := c.Param("id")
	todoID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "شناسه نامعتبر")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = collection.DeleteOne(ctx, bson.M{"_id": todoID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "خطا در حذف کار")
	}
	return c.JSON(http.StatusOK, "کار حذف شد")
}
