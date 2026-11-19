package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

type Todo struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var todoList []Todo

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	paramsId := r.PathValue("id")

	for _, item := range todoList {
		if item.Id == paramsId {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	todo.Id = strconv.Itoa(rand.Intn(10000))
	todo.Status = "Pending"
	todoList = append(todoList, todo)
	json.NewEncoder(w).Encode(todo)
}

func CompleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	paramsId := r.PathValue("id")

	for idx, item := range todoList {
		if item.Id == paramsId {
			todoList[idx].Status = "Completed"
			json.NewEncoder(w).Encode(todoList[idx])
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	paramsId := r.PathValue("id")

	for idx, item := range todoList {
		if item.Id == paramsId {
			if item.Status == "Completed" {
				todoList = append(todoList[:idx], todoList[idx+1:]...)
				json.NewEncoder(w).Encode("Todo deleted successfully...!")
				return
			}
			http.Error(w, "Todo is pending", http.StatusConflict)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func main() {

	todoList = append(todoList, Todo{Id: "0001", Description: "Write Todo App", Status: "Pending"})
	todoList = append(todoList, Todo{Id: "0002", Description: "Test Todo App", Status: "Pending"})

	http.HandleFunc("GET /todos", GetAllTodos)
	http.HandleFunc("GET /todos/{id}", GetTodo)
	http.HandleFunc("POST /todos", CreateTodo)
	http.HandleFunc("PATCH /todos/{id}/completed", CompleteTodo)
	http.HandleFunc("DELETE /todos/{id}", DeleteTodo)

	fmt.Printf("Starting Todo List server at :8080...!")

	http.ListenAndServe(":8080", nil)
}
