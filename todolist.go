package main

import (
 "context"
 "encoding/json"
 "log"
 "net/http"

 "github.com/labstack/echo/v4"
 "go.mongodb.org/mongo-driver/bson"
 "go.mongodb.org/mongo-driver/mongo"
 "go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
 ID   string json:"id" bson:"_id,omitempty"
 Name string json:"name" bson:"name"
}

func main() {
 
 client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
 if err != nil {
  log.Fatal("وصل نیست, err)
 }
 defer client.Disconnect(context.TODO())

 db := client.Database("todoDB")
 collection := db.Collection("tasks")


 todolist := []string{"خرید کردن", "کتاب خوندن", "فیلم دیدن", "پختن غذا"}
 for _, task := range todolist {
  _, err := collection.InsertOne(context.TODO(), Task{Name: task})
  if err != nil {
   log.Fatal("نمی‌تونم ذخیره کنم خطا:", err)
  }
 }


 tasks, err := getAllTasks(collection)
 if err != nil {
  log.Fatal("نمی‌تونم وظایف رو بگیرم خطا:", err)
 }
 fmt.Println("کارهایی که امروز باید انجام بدم:")
 for i, task := range tasks {
  fmt.Println(i+1, task.Name)
 }

 
 e := echo.New()

 
 e.GET("/todos", getTodos(collection))
 e.POST("/todos", createTodo(collection))
 e.PUT("/todos/:id", updateTodo(collection))
 e.DELETE("/todos/:id", deleteTodo(collection))


 fmt.Println("وب‌سایت من روشنه! برو به پورت 8080...")
 if err := e.Start(":8080"); err != nil {
  log.Fatal("وب‌سایت کار نمی‌کنه! خطا:", err)
 }
}


func getAllTasks(collection *mongo.Collection) ([]Task, error) {
 var tasks []Task
 cursor, err := collection.Find(context.TODO(), bson.M{})
 if err != nil {
  return nil, err
 }
 defer cursor.Close(context.TODO())

 for cursor.Next(context.TODO()) {
  var task Task
  if err := cursor.Decode(&task); err != nil {
   return nil, err
  }
  tasks = append(tasks, task)
 }
 return tasks, nil
}


func getTodos(collection *mongo.Collection) echo.HandlerFunc {
 return func(c echo.Context) error {
  tasks, err := getAllTasks(collection)
  if err != nil {
   return c.JSON(http.StatusInternalServerError, map[string]string{"error": "خطا در دریافت "})
  }
  return c.JSON(http.StatusOK, tasks)
 }
}


func createTodo(collection *mongo.Collection) echo.HandlerFunc {
 return func(c echo.Context) error {
  var task Task
  if err := c.Bind(&task); err != nil {
   return c.JSON(http.StatusBadRequest, map[string]string{"error": "داده نامعتبره"})
  }
  if task.Name == "" {
   return c.JSON(http.StatusBadRequest, map[string]string{"error": "نام نمی‌تونه خالی باشه"})
  }

  _, err := collection.InsertOne(context.TODO(), task)
  if err != nil {
   return c.JSON(http.StatusInternalServerError, map[string]string{"error": "نمی‌تونم ذخیره کنم"})
  }
  return c.JSON(http.StatusCreated, task)
 }
}


func updateTodo(collection *mongo.Collection) echo.HandlerFunc {
 return func(c echo.Context) error {
  id := c.Param("id")
  var task Task
  if err := c.Bind(&task); err != nil {
   return c.JSON(http.StatusBadRequest, map[string]string{"error": " نامعتبره"})
  }

  update := bson.M{"$set": bson.M{"name": task.Name}}
  _, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
  if err != nil {
   return c.JSON(http.StatusInternalServerError, map[string]string{"error": "نمی‌تونم به‌روزرسانی کنم"})
  }
  return c.JSON(http.StatusOK, task)
 }
}

func deleteTodo(collection *mongo.Collection) echo.HandlerFunc {
 return func(c echo.Context) error {
  id := c.Param("id")
  _, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
  if err != nil {
   return c.JSON(http.StatusInternalServerError, map[string]string{"error": "نمی‌تونم حذف کنم"})
  }
  return c.NoContent(http.StatusNoContent)
 }
}
