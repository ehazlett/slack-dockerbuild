package main

import (
        "bytes"
        "encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	url string
        channel string
)

type (

    PushData struct {
        Pusher      string  `json:"pusher"`
    }

    Repository struct {
        Status  string      `json:"status"`
        Description string  `json:"description"`
        IsTrusted   bool    `json:"is_trusted"`
        RepoURL     string  `json:"repo_url"`
        Owner       string  `json:"owner"`
        IsOfficial  bool    `json:"is_official"`
        IsPrivate   bool    `json:"is_private"`
        Name        string  `json:"name"`
        Namespace   string  `json:"namespace"`
        StarCount   int64   `json:"star_count"`
        CommentCount    int64   `json:"comment_count"`
        Dockerfile  string  `json:"dockerfile"`
        RepoName    string  `json:"repo_name"`
    }

    TrustedBuild struct {
        PushData    PushData    `json:"push_data"`
        Repository  Repository  `json:"repository"`
    }

    Payload struct {
        Channel     string  `json:"channel,omitempty"`
        Text        string  `json:"text"`
        Username    string  `json:"username,omitempty"`
    }
    SlackPayload struct {
        Payload     Payload `json:"payload"`
    }
)

func init() {
	flag.StringVar(&url, "url", "", "Slack Incoming Webhook URL")
        flag.StringVar(&channel, "channel", "#general", "Slack Channel")
	flag.Parse()
	if url == "" {
		log.Fatal("You must specify a url")
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("slack-dockerbuild"))
}


func notifyHandler(w http.ResponseWriter, r *http.Request) {
        // parse build info
        decoder := json.NewDecoder(r.Body)
        var build TrustedBuild
        err := decoder.Decode(&build)
        if err != nil {
                log.Printf("Error decoding JSON from Docker build: %s", err)
                return
        }
        // send to slack
        txt := fmt.Sprintf("%s build by %s complete", build.Repository.RepoName, build.PushData.Pusher)
        payload := Payload{Channel: channel, Text: txt}
        b, err := json.Marshal(payload)
        if err != nil {
                log.Printf("Error encoding payload to JSON: %s", err)
                return
        }
        rdr := bytes.NewReader(b)
        _, err = http.Post(url, "application/json", rdr)
        if err != nil {
                log.Printf("Error sending to Slack: %s", err)
                return
        }
        w.WriteHeader(200)
        w.Write([]byte("kthxbye"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/notify", notifyHandler)

	http.Handle("/", r)
	fmt.Println("Running on :8080...")
	http.ListenAndServe(":8080", nil)
}
