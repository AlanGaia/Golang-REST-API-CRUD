package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Type of singular task. Object with ID Name and Content
type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

//Type of allTasks an Array of objects (tasks)
type allTasks []task

//All tasks, could it be the DataBase
var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

// JSON with all the tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

//Method for Create New Tasks and append it to allTasks
func createTask(w http.ResponseWriter, r *http.Request) {
	//newTask instance of task {ID,Name,Content}
	var newTask task
	//If success reqBody else you got an error
	reqBody, err := ioutil.ReadAll(r.Body)
	//If err true, response insert a valid Task
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task")
	}
	//From reqBody assign data to newTask
	json.Unmarshal(reqBody, &newTask)
	//Generate an ID from lenght of tasks
	newTask.ID = len(tasks) + 1
	//Append new Task to allTask Array and Save it
	tasks = append(tasks, newTask)
	//Response from server to Client with the New Task With Header Application / JSON
	w.Header().Set("Content-Type", "application/json")
	//Response Status
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	//From request save all in vars
	vars := mux.Vars(r)

	//from vars get the id and try to parse to Number
	taskID, err := strconv.Atoi(vars["id"])

	//If request ID isn't a number return invalid ID
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	//Else search ID in Tasks Array
	for _, task := range tasks {
		//If Match Return that Task in JSON
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	//Routes
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")

	//Server Listening
	log.Fatal(http.ListenAndServe(":3000", router))
}
