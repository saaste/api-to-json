package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/saaste/api-to-json/listenbrainz"
	"github.com/saaste/api-to-json/webmentions"
)

type Output struct {
	TopArtists       []OutputItem `json:"top_artists"`
	TotalArtistCount int64        `json:"total_artist_count"`
}

type OutputItem struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []string `json:"tags"`
}

func main() {
	lbUserToken := flag.String("lb-user-token", "", "ListenBrainz User Token")
	wmToken := flag.String("wm-token", "", "Webmention token")
	flag.Parse()

	if lbUserToken == nil || *lbUserToken == "" {
		log.Fatalln("lb-user-token is required")
	}

	if wmToken == nil || *wmToken == "" {
		log.Fatalln("wm-token is required")
	}

	runListenBrainz(*lbUserToken)
	runWebtoken(*wmToken)
}

func runListenBrainz(lbUserToken string) {
	log.Println("Fetching artist stats from ListenBrainz")

	lbClient := listenbrainz.NewClient(lbUserToken)
	topArtistsResult, err := lbClient.FetchTopArtists(10)
	if err != nil {
		log.Fatalf("Failed to fetch artist stats from ListenBrainz: %v", err)
	}

	output := Output{
		TopArtists:       make([]OutputItem, 0),
		TotalArtistCount: topArtistsResult.TotalArtistCount,
	}

	for _, artist := range topArtistsResult.TopArtists {
		output.TopArtists = append(output.TopArtists, OutputItem{
			Title: artist.Name,
			URL:   artist.URL,
			Tags:  artist.Tags,
		})
	}

	outputFile, err := os.Create("listenbrainz.json")
	if err != nil {
		log.Fatalf("creating an output file failed: %v", err)
	}

	defer outputFile.Close()
	jsonOutput, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		log.Fatalf("marshaling output failed: %v", err)
	}
	outputFile.Write(jsonOutput)
	log.Println("Fetching artist stats from ListenBrainz done")
}

func runWebtoken(wmToken string) {
	log.Println("Fetching Webmentions")

	wmClient := webmentions.NewClient(wmToken)
	resp, err := wmClient.GetWebmentions("saaste.net")
	if err != nil {
		log.Fatalf("fetching webmentions failed: %v", err)
	}

	outputFile, err := os.Create("webmentions.json")
	if err != nil {
		log.Fatalf("creating an output file failed: %v", err)
	}

	defer outputFile.Close()
	jsonOutput, err := json.MarshalIndent(resp.Webmentions, "", "\t")
	if err != nil {
		log.Fatalf("marshaling output failed: %v", err)
	}
	outputFile.Write(jsonOutput)
	log.Println("Fetching Webmentions done")
}
