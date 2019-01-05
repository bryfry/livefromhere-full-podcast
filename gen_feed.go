package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/feeds"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Envelope struct {
	Items []Item
}

type Item struct {
	Audio Audio
}

type Audio struct {
	Encodings []Encoding
	Title     string
	ProgramID string `json:"program_id"`
	UpdatedAt string `json:"updated_at"`
}

type Encoding struct {
	HttpFilePath string `json:"http_file_path"`
}

func main() {

	// Fetch data from the live from here website
	url := "https://www.livefromhere.org"
	latestShows := url + "/listen/1.json"
	res, err := http.Get(latestShows)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Marshal the data into our datastructures with only the fields we care about
	var env Envelope
	err = json.Unmarshal(body, &env)
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now()

	// Generate the feed from the data
	newfeed := &feeds.AtomFeed{
		Title:   "Live From Us - Full Shows",
		Id:      "LiveFromUs-Full",
		Link:    &feeds.AtomLink{Href: url},
		Author:  &feeds.AtomAuthor{AtomPerson: feeds.AtomPerson{Name: "APM", Email: ""}},
		Updated: now.String(),
	}
	for _, i := range env.Items {
		newentry := &feeds.AtomEntry{
			Title:   i.Audio.Title,
			Id:      i.Audio.ProgramID,
			Updated: i.Audio.UpdatedAt,
			Links: []feeds.AtomLink{feeds.AtomLink{
				Href:   i.Audio.Encodings[0].HttpFilePath,
				Type:   "audio/mpeg",
				Length: "1024",
				Rel:    "enclosure",
			}},
		}
		newfeed.Entries = append(newfeed.Entries, newentry)
	}
	writer, _ := os.Create("feed.xml")
	feeds.WriteXML(newfeed, writer)
	fmt.Println("feed update complete", now)
}
