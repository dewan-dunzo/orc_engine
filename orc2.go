package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type creationPayload struct {
	Name string `json:"Name"`
}

type schedulePayload struct {
	TaskID    int   `json:"Id"`
	TimeStamp int64 `json:"Timestamp"`
}

type CreateAndScheduleRequest struct {
	Name      string `json:"Name"`
	TimeStamp int64  `json:"Timestamp"`
}

func createRequest(newTask creationPayload) []byte {

	jsonResp, _ := json.Marshal(newTask)

	fmt.Println("Incoming data")
	fmt.Println(string(jsonResp))

	url := "http://192.168.1.229:8000/todo/"
	// url := "http://127.0.0.1:8000/task"
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

// func create(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")
// 	var newTask interface{}

// 	_ = json.NewDecoder(r.Body).Decode(&newTask)
// 	w.Write([]byte(createRequest(newTask)))
// }

// func schedule(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var newTask interface{}
// 	// id := c.Params.ByName("id")
// 	// fmt.Println("id >", id)
// 	err := json.NewDecoder(r.Body).Decode(&newTask)
// 	newTask.TaskID, _ = strconv.Atoi(r.URL.Path[10:])
// 	w.Write(scheduleReq(newTask))
// }

func scheduleReq(newTask schedulePayload) []byte {

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

func createAndSchedule(w http.ResponseWriter, r *http.Request) {

	var newTask CreateAndScheduleRequest
	_ = json.NewDecoder(r.Body).Decode(&newTask)

	var CPayload creationPayload
	CPayload.Name = newTask.Name

	jsonResp := createRequest(CPayload)

	var SPayload schedulePayload
	json.Unmarshal(jsonResp, &SPayload)
	SPayload.TimeStamp = newTask.TimeStamp

	scheduleReq(SPayload)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("something"))
}

func main() {
	r := mux.NewRouter()

	// r.HandleFunc("/create", create).Methods("POST")
	// r.HandleFunc("/schedule/{id:[0-9]+}", schedule).Methods("POST")
	r.HandleFunc("/createAndSchedule", createAndSchedule).Methods("POST")

	http.ListenAndServe(":8001", r)
}
