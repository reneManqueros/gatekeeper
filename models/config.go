package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type config struct {
	Listeners []Listener `json:"listeners"`
}

var Config config

func (c *config) Load() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalln(err)
	}
}
