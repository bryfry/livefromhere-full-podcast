package main

import (
	//"bytes"
	"github.com/gorilla/feeds"
	//"github.com/mmcdole/gofeed"
	"log"
	"os"
	//"os/exec"
	//"strings"
	"time"
	"fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type Envelope struct {
  Items []Item
}

type Item struct {
   Audio Audio
}

type Audio struct {
  Encodings []Encoding
  Title string
  ProgramID string `json:"program_id"`
  UpdatedAt string `json:"updated_at"`
}

type Encoding struct {
  HttpFilePath string `json:"http_file_path"`
}

func main() {
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
  var env Envelope
  err = json.Unmarshal(body, &env)
  if err != nil {
    log.Fatal(err)
  }
	now := time.Now()
	newfeed := &feeds.AtomFeed{
		Title:   "Live From Us - Full Shows",
    Id:      "LiveFromUs-Full",
		Link:    &feeds.AtomLink{Href:"https://www.livefromhere.org/"},
		Author:  &feeds.AtomAuthor{AtomPerson: feeds.AtomPerson{Name: "APM", Email: ""}},
		Updated: now.String(), 
	}
  for _, i := range env.Items {
		newentry := &feeds.AtomEntry{
			Title:   i.Audio.Title,
			Id:      i.Audio.ProgramID,
      Updated: i.Audio.UpdatedAt,
			Links:   []feeds.AtomLink {feeds.AtomLink{
           Href: i.Audio.Encodings[0].HttpFilePath,
           Type: "audio/mpeg",
           Length: "1024",
           Rel: "enclosure",
           }},
		}
		newfeed.Entries = append(newfeed.Entries, newentry)
  }
	writer, _ := os.Create("feed.xml")
	feeds.WriteXML(newfeed, writer)
	fmt.Println("feed update complete", now)
}

//func x(){
//	YouTubePlaylistRSS := "https://www.youtube.com/feeds/videos.xml?playlist_id=PLiZxWe0ejyv8CSMylrxb6Nx4Ii2RHbu_j"
//	audiodir := "/var/www/audio/"
//	fp := gofeed.NewParser()
//	feed, _ := fp.ParseURL(YouTubePlaylistRSS)
//	now := time.Now()
//
//	newfeed := &feeds.AtomFeed{
//		Title:      "Late Show w/ Stephen Colbert Intro Monologues",
//		Link:       &feeds.AtomLink{Href: "https://www.youtube.com/channel/UCMtFAi84ehTSYSE9XoHefig"},
//		Subtitle:   "Podcast version of Intro Monologues",
//		Author:     &feeds.AtomAuthor{AtomPerson: feeds.AtomPerson{Name: "CBS", Email: ""}},
//		Updated:    now.String(), 
//	}
//
//	for _, i := range feed.Items {
//		guid := strings.Split(i.GUID, ":")
//		id := guid[len(guid)-1]
//
//		if _, err := os.Stat(audiodir + id + ".mp3"); os.IsNotExist(err) {
//			// id mp3 hasn't been downloaded
//			fmt.Println("downloading: ", id)
//			cmd := exec.Command("youtube-dl", "--extract-audio", "--audio-format", "mp3", "-o", audiodir+"%(id)s.%(ext)s", i.Link)
//			var out bytes.Buffer
//			cmd.Stdout = &out
//			err := cmd.Run()
//			if err != nil {
//				log.Fatal(err)
//			}
//		}
//			//fmt.Println("adding: ", id)
//
//		newentry := &feeds.AtomEntry{
//			Title:   i.Title,
//			Id: id,
//			Link:    &feeds.AtomLink{Href: "https://trustme.click/audio/" + id + ".mp3", Type: "audio/mpeg", Length: "1024", Rel: "enclosure"},
//			Updated: i.Published,
//		}
//		newfeed.Entries = append(newfeed.Entries, newentry)
//	}
//	writer, _ := os.Create("/var/www/audio/feed.xml")
//	feeds.WriteXML(newfeed, writer)
//	fmt.Println("feed update complete", now)
//
//}
