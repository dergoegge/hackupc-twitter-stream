package sentiment140

import ("encoding/json"
        "net/http"
        "log"
        "bytes"
        "fmt")

type body struct {
        Data []tweet `json:"data"`
}

type tweet struct {
        Text string `json:"text"`
}

var tweets body

func Add(text string) {
        tweets.Data = append(tweets.Data, tweet{text})
        fmt.Println(tweets.Data)
}

func Post() float64 {
        postBody, err := json.Marshal(tweets)
        if err != nil {
                log.Fatal(err)
        }
        
        fmt.Println(string(postBody))

        resp, err := http.Post("http://www.sentiment140.com/api/bulkClassifyJson", "application/json", bytes.NewBuffer(postBody))
        if err != nil {
                log.Fatal(err)
        }
        fmt.Printf("\n%s", resp.Body)
        
        tweets.Data = nil
        return 0
}
