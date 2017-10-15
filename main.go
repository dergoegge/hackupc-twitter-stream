package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dergoegge/hackupc-twitter-analysis/config"
	"github.com/dghubble/go-twitter/twitter"

	sarama "gopkg.in/Shopify/sarama.v1"
)

func main() {
	// flags
	var query = flag.String("query", "", "The query by which the tweets get fetched.")
	var printTweets = flag.Bool("print", true, "Set to false if tweets should not be printed to stdout.")
	flag.Parse()

	if strings.Compare(strings.Trim(*query, " "), "") == 0 {
		fmt.Println("A query needs to be specified! \nTry -help for help.")
		os.Exit(0)
	}

	// setup twitter api
	httpClient := config.LoadHTTPClient()

	client := twitter.NewClient(httpClient)

	tweetDemux := twitter.NewSwitchDemux()

	// setup kafka golang api
	config := sarama.NewConfig()

	config.Producer.Return.Successes = false
	producer, err := sarama.NewAsyncProducer([]string{"54.186.93.122:9092", "54.218.59.178:9092"}, config)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	count := 1

	// function gets called once per recieed tweet
	tweetDemux.Tweet = func(tweet *twitter.Tweet) {
		if !strings.HasPrefix(tweet.Text, "RT ") {
			message := &sarama.ProducerMessage{Topic: "tweets-test", Value: sarama.StringEncoder(tweet.Text)}
			select {
			case producer.Input() <- message:
				count++
			}

			if *printTweets {
				fmt.Println(tweet.Text + "\n")
			}
		}
	}

	// search parameters for the twitter api
	params := &twitter.StreamFilterParams{
		Track:         []string{*query},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(params)
	if err != nil {
		log.Fatal(err)
	}

	// go handle incomming tweets
	go tweetDemux.HandleChan(stream.Messages)

	// handle keyboard interrupts
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	log.Println(<-ch)

	// clean up twitter & kafka apis
	producer.AsyncClose()
	stream.Stop()

	fmt.Println("--------------")
	fmt.Printf("%d tweets where send to the kafka cluster.\n", count)
}
