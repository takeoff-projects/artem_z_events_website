package eventsdb

import (
	"errors"

	"github.com/google/uuid"
)

// Event model
type Event struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Location string `json:"location"`
    When   string `json:"when"`
}

var Events []Event
 
func GetEvents() []Event {
    return Events
}

func InitializeEventsArray(){
	Events = []Event{
		{Title: "Dinner",
			Location: "My House",
			When:   "Tonight",
			ID: "2944a9cb-ef2d-4632-ac1d-af2b2629d0f2"},
		{Title: "Go Programming Lesson",
			Location: "At School",
			When:   "Tomorrow",
			ID: "f88f1860-9a5d-423e-820f-9acb4db3030e"},
		{Title: "Company Picnic",
			Location: "At the Park",
			When:   "Saturday",
			ID: "4cb393fb-dd19-469e-a52c-22a12c0a98df"},
	}
}

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