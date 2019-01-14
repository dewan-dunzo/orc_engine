package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Task struct for use
type Task struct {
	Name   string `json:"userName"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
	Tid    string `json:"tid"`
}

type Person struct {
	userName string `json:"userName"`
}

const SERVER = "mongodb://localhost:27017"
const DBNAME = "goToDo"
const COLLECTION = "list"

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	_ = json.NewDecoder(r.Body).Decode(&newTask)
	fmt.Println(newTask, strconv.Itoa(rand.Intn(10000000)))
	(&newTask).Tid = strconv.Itoa(rand.Intn(10000000))

	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	c := session.DB(DBNAME).C(COLLECTION)

	if err := c.Insert(&newTask); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	w.Write([]byte("Task added!\n"))
}

func showTasks(w http.ResponseWriter, r *http.Request) {
	uName, _ := r.URL.Query()["name"]

	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	c := session.DB(DBNAME).C(COLLECTION)

	var bsonEPs []bson.M
	query := c.Find(bson.M{"name": uName[0]}).All(&bsonEPs)
	fmt.Println(query)

	jsonResp, merr := json.Marshal(bsonEPs)
	if merr != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nJSON DATA\n")
	fmt.Println(string(jsonResp))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResp))
}

func updateTasks(w http.ResponseWriter, r *http.Request) {
	var task Task
	id := r.URL.Path[6:]
	fmt.Println(id)

	_ = json.NewDecoder(r.Body).Decode(&task)
	fmt.Println(task)

	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	c := session.DB(DBNAME).C(COLLECTION)

	var bsonEPs []bson.M
	fmt.Println("all details", id, task.Title, task.Status)
	query := c.Update(bson.M{"tid": id}, bson.M{"$set": bson.M{"title": task.Title, "status": task.Status}})

	query_find := c.Find(bson.M{"tid": id}).All(&bsonEPs)

	jsonResp, merr := json.Marshal(bsonEPs)
	if merr != nil {
		fmt.Println("query result", err)
	}

	fmt.Printf("\nJSON DATA\n")
	fmt.Println(string(jsonResp))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResp))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/task", createTask).Methods("POST")
	r.HandleFunc("/task", showTasks).Methods("GET")
	r.HandleFunc("/task/{id}", updateTasks).Methods("PUT")

	http.ListenAndServe(":8000", r)
}
