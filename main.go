package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	sarama "gopkg.in/Shopify/sarama.v1"

	"github.com/dergoegge/hackupc-twitter-analysis/config"
	"github.com/dghubble/go-twitter/twitter"
)

func main() {

	httpClient := config.LoadHTTPClient()

	client := twitter.NewClient(httpClient)

	tweetDemux := twitter.NewSwitchDemux()

	config := sarama.NewConfig()

	config.Producer.Return.Successes = false
	producer, err := sarama.NewAsyncProducer([]string{"54.186.93.122:9092"}, config)
	if err != nil {
		panic(err)
	}

	/*model, err := sentiment.Restore()
	if err != nil {
		panic(fmt.Sprintf("Could not restore model!\n\t%v\n", err))
	}
	*/
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	count := 1
	//senti := 1

	tweetDemux.Tweet = func(tweet *twitter.Tweet) {
		if !strings.HasPrefix(tweet.Text, "RT ") {
			message := &sarama.ProducerMessage{Topic: "tweets-test", Value: sarama.StringEncoder(tweet.Text)}
			select {
			case producer.Input() <- message:
				count++
			}
			fmt.Println(tweet.Text)
			/*analysis := model.SentimentAnalysis(tweet.Text, sentiment.English) // 0
			senti += int(analysis.Score)

			s := float32(senti) / float32(count)
			fmt.Printf("Tweets : %d, Positive : %f \n", count, s)*/
			/*sentiment140.Add(tweet.Text)

			if count > 10 {
				sentiment140.Post()
				os.Exit(0)
			}*/
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

	producer.AsyncClose()

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
