package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	var todolist []string

	todolist = append(todolist, "خرید کردن")
	todolist = append(todolist, "کتاب خوندن")
	todolist = append(todolist, "فیلم دیدن")
	todolist = append(todolist, "پختن غذا")

	fmt.Println("کارهایی که امروز باید انجام بدم:")
	for i, harkar := range todolist {
		fmt.Println(i+1, ".", harkar)
	}

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(todolist); err != nil {
				http.Error(w, "خطا در تبدیل به JSON", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "فقط متد GET پشتیبانی می‌شه", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("سرور در حال اجرا روی پورت 8080...")
	http.ListenAndServe(":8080", nil)
}
