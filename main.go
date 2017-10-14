package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dergoegge/hackupc-twitter-analysis/config"
	"github.com/dghubble/go-twitter/twitter"
)

func main() {
	httpClient := config.LoadHTTPClient()

	client := twitter.NewClient(httpClient)

	tweetDemux := twitter.NewSwitchDemux()

	tweetDemux.Tweet = func(tweet *twitter.Tweet) {
		if !strings.HasPrefix(tweet.Text, "RT ") {
			fmt.Println(tweet.Text)
		}
	}

	params := &twitter.StreamFilterParams{
		Track:         []string{"trump"},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(params)
	if err != nil {
		log.Fatal(err)
	}

	go tweetDemux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
