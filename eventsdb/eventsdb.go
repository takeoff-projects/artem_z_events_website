package eventsdb

import (
	"fmt"
	"os"

	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	// "google.golang.org/appengine/log"
)

// Event model
type Event struct {
	ID       string `firestore:"Id"`
	Title    string `firestore:"title"`
	Location string `firestore:"location"`
	When     string `firestore:"when"`
}

var projectID string

var ctx = context.Background()

func connect() *firestore.Client {

	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(ctx, "Error creating firestore client")
	}
	return client
}

func GetEvents() []Event {

	var Events []Event

	client := connect()
	iter := client.Collection("Events").Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(ctx, "Error in iterating over docs")
		}
		var e Event
		if err := doc.DataTo(&e); err != nil {
			log.Fatal(ctx, "Error in transforming doc to Event")
		}
		e.ID = doc.Ref.ID
		Events = append(Events, e)
	}
	fmt.Println(Events)
	client.Close()
	return Events
}

func GetEventbyID(key string) (Event, error) {

	var ev Event

	fmt.Println("Key", key)

	client := connect()
	evRef := client.Collection("Events").Doc(key)
	docsnap, err := evRef.Get(ctx)
	fmt.Println(docsnap)
	if err != nil {
		log.Fatal(ctx, "Error in getting event by id")
	}
	if err := docsnap.DataTo(&ev); err != nil {
		log.Fatal(ctx, "Error in transforming firestoree doc to event by id")
	}
	client.Close()
	return ev, err
}

func AddEvent(event Event) {
	client := connect()
	_, wr, err := client.Collection("Events").Add(ctx, event)
	if err != nil {
		log.Print("Error creating doc")
	}
	log.Println(wr)
	client.Close()
}

func UpdateEvent(eventID string, updatedEvent Event) {
	client := connect()
	evRef := client.Collection("Events").Doc(eventID)
	fmt.Println("Event to be updated: ", evRef)
	_, err := evRef.Update(ctx, []firestore.Update{
		{Path: "Title", Value: updatedEvent.Title},
		{Path: "Location", Value: updatedEvent.Location},
		{Path: "When", Value: updatedEvent.When}})
	if err != nil {
		log.Fatal(ctx, "Error in updating event")
	}
	client.Close()
}

func DeleteEvent(id string) {
	client := connect()
	doc := client.Collection("Events").Doc(id)
	_, err := doc.Delete(ctx)
	if err != nil {
		log.Println(ctx, "Error deleting doc")
	}
}
