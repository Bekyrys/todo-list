package main

import (
	"encoding/json"
	"github.com/Bekyrys/todo-list/internal/tasks"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/todo-list/tasks", createTaskHandler).Methods("POST")
	r.HandleFunc("/api/todo-list/tasks", getTasksHandler).Methods("GET")
	r.HandleFunc("/api/todo-list/tasks/{id}", updateTaskHandler).Methods("PUT")
	r.HandleFunc("/api/todo-list/tasks/{id}", deleteTaskHandler).Methods("DELETE")
	r.HandleFunc("/api/todo-list/tasks/{id}/done", markTaskDoneHandler).Methods("PUT")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t struct {
		Title    string `json:"title"`
		ActiveAt string `json:"activeAt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	activeAt, err := time.Parse("02 January 15:04", t.ActiveAt)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	task, err := tasks.CreateTask(t.Title, activeAt)
	if err != nil {
		if err.Error() == "task already exists" {
			http.Error(w, "task already exists", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": task.ID})
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "active"
	}

	tasksList := tasks.GetTasks(status)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasksList)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var t struct {
		Title    string `json:"title"`
		ActiveAt string `json:"activeAt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	activeAt, err := time.Parse("02 January 15:04", t.ActiveAt)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	err = tasks.UpdateTask(id, t.Title, activeAt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	updatedTask, _ := tasks.GetTaskByID(id)
	response := struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		ActiveAt  string `json:"activeAt"`
		Completed bool   `json:"completed"`
	}{
		ID:        updatedTask.ID,
		Title:     updatedTask.Title,
		ActiveAt:  updatedTask.ActiveAt.Format("02 January 15:04"),
		Completed: updatedTask.Done,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := tasks.DeleteTask(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func markTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := tasks.MarkTaskDone(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
