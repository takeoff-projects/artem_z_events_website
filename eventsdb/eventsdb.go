package eventsdb

import (
	"errors"
	"os"

	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
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

var Events []Event

func GetEvents() []Event {

	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(ctx, "Error creating firestore client")
	}
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
		Events = append(Events, e)
	}
	return Events
}

// func GetEvents() []Event {
// 	return Events
// }

// func InitializeEventsArray() {
// 	Events = []Event{
// 		{Title: "Dinner",
// 			Location: "My House",
// 			When:     "Tonight",
// 			ID:       "2944a9cb-ef2d-4632-ac1d-af2b2629d0f2"},
// 		{Title: "Go Programming Lesson",
// 			Location: "At School",
// 			When:     "Tomorrow",
// 			ID:       "f88f1860-9a5d-423e-820f-9acb4db3030e"},
// 		{Title: "Company Picnic",
// 			Location: "At the Park",
// 			When:     "Saturday",
// 			ID:       "4cb393fb-dd19-469e-a52c-22a12c0a98df"},
// 	}
// }

func GetEventbyID(key string) (Event, error) {
	for _, event := range Events {
		if event.ID == key {
			return event, nil
		}
	}
	return Event{}, errors.New("not found")
}

func AddEvent(event Event) {
	newID := uuid.New().String()
	event.ID = newID
	Events = append(Events, event)
}

func UpdateEvent(updatedEvent Event) {
	for index, event := range Events {
		if event.ID == updatedEvent.ID {
			Events[index].Title = updatedEvent.Title
			Events[index].Location = updatedEvent.Location
			Events[index].When = updatedEvent.When
			return
		}
	}
}

func DeleteEvent(id string) {
	for index, event := range Events {
		if event.ID == id {
			Events = append(Events[:index], Events[index+1:]...)
		}
	}
}
