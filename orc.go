package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type creationPayload struct {
	Name string `json:"Name"`
	ID   int    `json:"Id"`
}

type schedulePayload struct {
	TaskID    int   `json:"Id"`
	TimeStamp int64 `json:"Timestamp"`
}

func createRequest(r *http.Request) []byte {
	var newTask interface{}

	_ = json.NewDecoder(r.Body).Decode(&newTask)

	jsonResp, _ := json.Marshal(newTask)

	fmt.Println("Incoming data")
	fmt.Println(string(jsonResp))

	url := "http://192.168.1.229:8000/todo/"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(jsonResp)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body
}

func create(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(createRequest(r)))
}

func scheduleReq(r *http.Request, tid string) []byte {
	var newTask interface{}
	id := ""
	if tid == "" {
		id = r.URL.Path[10:]
	} else {
		id = tid
	}

	// id := c.Params.ByName("id")
	fmt.Println("id >", id)
	err := json.NewDecoder(r.Body).Decode(&newTask)
	fmt.Println(err)
	// newTask.TaskID, _ = strconv.Atoi(id)

	// fmt.Println(newTask, id, newTask.TaskID)

	jsonResp, _ := json.Marshal(newTask)

	fmt.Println("Incoming data")
	fmt.Println(string(jsonResp))

	url := "http://192.168.1.188:8000/schedule"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(jsonResp)
	fmt.Println("sending scheduling data", newTask)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return jsonResp
}

func schedule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(scheduleReq(r, ""))
}

func createAndSchedule(w http.ResponseWriter, r *http.Request) {
	var Payload creationPayload
	var r1 http.Request
	var r2 http.Request
	r1 = *r
	r2 = *r
	jsonResp := createRequest(&r1)
	json.Unmarshal(jsonResp, &Payload)
	// fmt.Println(jsonResp, Payload)
	tid := Payload.ID
	scheduleReq(&r2, strconv.Itoa(tid))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/create", create).Methods("POST")
	r.HandleFunc("/schedule/{id:[0-9]+}", schedule).Methods("POST")
	r.HandleFunc("/createAndSchedule", createAndSchedule).Methods("POST")

	http.ListenAndServe(":8000", r)
}
