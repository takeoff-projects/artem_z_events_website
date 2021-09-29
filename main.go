package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"drehnstrom.com/go-website/eventsapi"
	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "7999"

	}
	log.Printf("Port set to: %s", port)

	fs := http.FileServer(http.Dir("assets"))
	myRouter := mux.NewRouter().StrictSlash(true)

	// This serves the static files in the assets folder
	myRouter.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// The rest of the routes
	myRouter.HandleFunc("/", indexHandler)
	myRouter.HandleFunc("/about", aboutHandler)
	myRouter.HandleFunc("/add", addHandler)
	myRouter.HandleFunc("/edit/{id}", editHandler)
	myRouter.HandleFunc("/delete/{id}", deleteHandler)

	log.Printf("Webserver listening on Port: %s", port)
	http.ListenAndServe(":"+port, myRouter)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var events = eventsapi.GetEvents()

	data := HomePageData{
		PageTitle: "Home Page",
		Events:    events,
		Count:     len(events),
	}

	var tpl = template.Must(template.ParseFiles("templates/index.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	buf.WriteTo(w)
	log.Println("Home Page Served")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := AboutPageData{
		PageTitle: "About Artem Z's events",
	}

	var tpl = template.Must(template.ParseFiles("templates/about.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	buf.WriteTo(w)
	log.Println("About Page Served")
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := AddPageData{
			PageTitle: "Add Event",
		}

		var tpl = template.Must(template.ParseFiles("templates/add.html", "templates/layout.html"))

		buf := &bytes.Buffer{}
		err := tpl.Execute(buf, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		buf.WriteTo(w)

		log.Println("Add Page Served")
	} else {

		// Add Event Here
		event := eventsapi.Event{
			Title:    r.FormValue("title"),
			Location: r.FormValue("location"),
			When:     r.FormValue("when"),
		}
		log.Println("Going to create new event: ", event)
		eventsapi.AddEvent(event)

		// Go back to home page
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	if r.Method == http.MethodGet {
		log.Println("Edit Handler - editing event", eventID)

		event, error := eventsapi.GetEventbyID(eventID)

		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			log.Println(error.Error())
			return
		}

		data := EditPageData{
			PageTitle: "Edit Event",
			Event:     event,
			EventID:   eventID,
		}
		log.Println("Editing event: ", event)
		var tpl = template.Must(template.ParseFiles("templates/edit.html", "templates/layout.html"))

		buf := &bytes.Buffer{}
		err := tpl.Execute(buf, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		buf.WriteTo(w)

		log.Println("Edit Page Served")
	} else {
		// Add Event Here
		event := eventsapi.Event{
			ID:       r.FormValue("id"),
			Title:    r.FormValue("title"),
			Location: r.FormValue("location"),
			When:     r.FormValue("when"),
		}
		fmt.Println("preparing to update event with: ", event)
		eventsapi.UpdateEvent(eventID, event)
		log.Println("Event Updated")

		// Go back to home page
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	eventsapi.DeleteEvent(mux.Vars(r)["id"])
	log.Println("Event deleted")

	// Go back to home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// HomePageData for Index template
type HomePageData struct {
	PageTitle string
	Events    []eventsapi.Event
	Count     int
}

// AboutPageData for About template
type AboutPageData struct {
	PageTitle string
}

// AddPageData for About template
type AddPageData struct {
	PageTitle string
}

// EditPageData for About template
type EditPageData struct {
	PageTitle string
	Event     eventsapi.Event
	EventID   string
}
