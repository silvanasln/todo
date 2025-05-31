package main

import (
 "context"
 "encoding/json"
 "fmt"
 "log"
 "net/http"

 "go.mongodb.org/mongo-driver/mongo"
 "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
 client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

 }
 defer client.Disconnect(context.TODO()) 

 db := client.Database("todoDB")      
 collection := db.Collection("tasks")   

 var todolist []string
 todolist = append(todolist, "خرید کردن")
 todolist = append(todolist, "کتاب خوندن")
 todolist = append(todolist, "فیلم دیدن")
 todolist = append(todolist, "پختن غذا")


 _, err = collection.InsertMany(context.TODO(), []interface{}{todolist})
 if err != nil {
  log.Fatal("نمی‌تونم لیست رو ذخیره کنم! خطا:", err)
 }

 fmt.Println("کارهایی که امروز باید انجام بدم:")
 for i, harkar := range todolist {
  fmt.Println(i+1,harkar)
 }


 http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
   w.Header().Set("Content-Type", "application/json")
   if err := json.NewEncoder(w).Encode(todolist); err != nil {
    http.Error(w, "خطا تو تبدیل به JSON", http.StatusInternalServerError)
    return
   }
  } else {
   http.Error(w, "فقط می‌تونی GET کنی!", http.StatusMethodNotAllowed)
  }
 })

 fmt.Println("وب‌سایت من روشنه! برو به پورت 8080...")
 if err := http.ListenAndServe(":8080", nil); err != nil {
  log.Fatal("وب‌سایت کار نمی‌کنه! خطا:", err)
 }
}
