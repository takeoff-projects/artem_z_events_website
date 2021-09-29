package eventsapi

import (
	"encoding/json"
	"fmt"
	"log"

	"os"

	"bytes"
	"io/ioutil"
	"net/http"
)

// Event model
type Event struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Location string `json:"location"`
	When     string `json:"when"`
}

var backendURL = os.Getenv("BACKEND_URL")

func GetEvents() []Event {

	var events []Event

	resp, err := http.Get(backendURL)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respBody, &events)

	return events
}

func GetEventbyID(key string) (Event, error) {

	var ev Event

	resp, err := http.Get(backendURL + "/" + key)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respBody, &ev)

	fmt.Println(ev)

	return ev, err
}

func AddEvent(event Event) {

	reqBody, _ := json.Marshal(event)

	resp, err := http.Post(
		backendURL,
		"application/json",
		bytes.NewBuffer(reqBody))

	fmt.Println(resp)

	if err != nil {
		log.Println("Error posing new event: ", err.Error())
	}
	defer resp.Body.Close()
}

func UpdateEvent(eventID string, updatedEvent Event) {

	reqBody, _ := json.Marshal(updatedEvent)

	req, _ := http.NewRequest(
		http.MethodPut,
		backendURL+"/"+eventID,
		bytes.NewBuffer(reqBody))

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)

	log.Println(resp)

	if err != nil {
		log.Println("Error updating event: ", err.Error())
	}
	defer resp.Body.Close()
}

func DeleteEvent(id string) {
	req, _ := http.NewRequest(
		http.MethodDelete,
		backendURL+"/"+id,
		nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error deleting event: ", err.Error())
	}
	defer resp.Body.Close()
}
