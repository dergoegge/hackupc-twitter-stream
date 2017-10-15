/*
 *Loads the twitter api keys from the config.json file provided in the config folder.
 */
 
package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
)

type config struct {
	ConfigKey    string `json:"config_key"`
	ConfigSecret string `json:"config_secret"`

	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
}

/*LoadHTTPClient ...*/
func LoadHTTPClient() *http.Client {
	configFile, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	
	defer configFile.Close()

	var conf config

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&conf) 
	if err != nil {
		log.Fatal(err)
	}
	
	config := oauth1.NewConfig(conf.ConfigKey, conf.ConfigSecret)
	token := oauth1.NewToken(conf.AccessKey, conf.AccessSecret)

	return config.Client(oauth1.NoContext, token)
}
