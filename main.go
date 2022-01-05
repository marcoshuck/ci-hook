package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Another test!
func main() {
	http.HandleFunc("/github/payload", TriggerEvent)
	// Test
	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Fatal error:", err)
	}
}

func TriggerEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid HTTP method", http.StatusBadRequest)
		return
	}

	event := getEvent(r)
	switch event {
	case "ping":
		triggerPing(w, r)
		break
	case "push":
		triggerPush(w, r)
		break
	default:
		http.Error(w, "invalid event", http.StatusBadRequest)
	}
}

func triggerPush(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error - Trigger push:", err)
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	var payload Payload
	err = json.Unmarshal(b, &payload)
	if err != nil {
		log.Println("Error - Trigger push:", err)
		http.Error(w, fmt.Sprintf("failed to read body: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%s triggered a %s event in the %s repository, reference: %s",
		payload.Repository.Owner.Login,
		"push",
		payload.Repository.Fullname,
		payload.Reference,
	)
	_, err = w.Write([]byte(body))
	if err != nil {
		log.Println("Error - Trigger push:", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func triggerPing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func getEvent(r *http.Request) string {
	return r.Header.Get("X-GitHub-Event")
}

type Payload struct {
	Reference  string     `json:"ref"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	Fullname string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    Owner  `json:"owner"`
	URL      string `json:"url"`
	Archived bool   `json:"archived"`
	Disabled bool   `json:"disabled"`
}

type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
}
