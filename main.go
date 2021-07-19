package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	http.HandleFunc("/github/payload", TriggerEvent)

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
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	var payload Payload
	err = json.Unmarshal(b, &payload)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	fmt.Printf("%+v", payload)

	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%s triggered a %s event in the %s, branch: %s",
		payload.Repository.Owner.Login,
		"push",
		payload.Repository.Fullname,
		payload.Reference,
	)
	_, err = w.Write([]byte(body))
	if err != nil {
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
	Reference     string     `json:"ref"`
	Before        string     `json:"before"`
	After         string     `json:"after"`
	Created       bool       `json:"created"`
	Deleted       bool       `json:"deleted"`
	Forced        bool       `json:"forced"`
	BaseReference string     `json:"base_ref"`
	Compare       url.URL    `json:"compare"`
	Commits       []string   `json:"commits"`
	HeadCommit    *string    `json:"head_commit"`
	Repository    Repository `json:"repository"`
}

type Repository struct {
	ID        int       `json:"id"`
	NodeID    string    `json:"node_id"`
	Name      string    `json:"name"`
	Fullname  string    `json:"full_name"`
	Private   bool      `json:"private"`
	Owner     Owner     `json:"owner"`
	URL       url.URL   `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PushedAt  time.Time `json:"pushed_at"`
	Archived  bool      `json:"archived"`
	Disabled  bool      `json:"disabled"`
}

type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
}
