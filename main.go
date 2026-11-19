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
	json.NewEncoder(w).Encode("Todo Not Found...!")
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	todo.Id = strconv.Itoa(rand.Intn(10000))
	todo.Status = "Pending"
	todoList = append(todoList, todo)
	json.NewEncoder(w).Encode(todo)
}

func CompleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		json.NewEncoder(w).Encode(fmt.Errorf("Invalid method type"))
	}

	paramsId := r.PathValue("id")

	for idx, item := range todoList {
		if item.Id == paramsId {
			todoList = append(todoList[:idx], todoList[idx+1:]...)

			item.Status = "Completed"
			todoList = append(todoList, item)

			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode("Todo Not Found...!")
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	paramsId := r.PathValue("id")

	for idx, item := range todoList {
		if item.Id == paramsId {
			if item.Status == "Completed" {
				todoList = append(todoList[:idx], todoList[idx+1:]...)
				json.NewEncoder(w).Encode("Todo deleted successfulyy...!")
				return
			}
			json.NewEncoder(w).Encode("Todo is pending, can't be deleted...!")
			return
		}
	}
	json.NewEncoder(w).Encode("Todo Not Found...!")
}

func main() {

	todoList = append(todoList, Todo{Id: "0001", Description: "Write Todo App", Status: "Pending"})
	todoList = append(todoList, Todo{Id: "0002", Description: "Test Todo App", Status: "Pending"})

	http.HandleFunc("GET /todos", GetAllTodos)
	http.HandleFunc("GET /todo/{id}", GetTodo)
	http.HandleFunc("POST /todos", CreateTodo)
	http.HandleFunc("PATCH /todo/{id}/completed", CompleteTodo)
	http.HandleFunc("DELETE /todo/{id}", DeleteTodo)

	fmt.Printf("Starting Todo List server at :8080...!")

	http.ListenAndServe(":8080", nil)
}
