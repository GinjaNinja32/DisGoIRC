package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/GinjaNinja32/DisGoIRC/bot"
)

func main() {
	debug := flag.Bool("debug", false, "Debug mode")
	confLocation := flag.String("config", "conf.json", "Config file location")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	var conf bot.Config
	confJSON, err := ioutil.ReadFile(*confLocation)
	if err != nil {
		log.Fatalf("Failed to read config file %s: %s", *confLocation, err)
	}

	err = json.Unmarshal(confJSON, &conf)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}

	bot.Init(conf)

	log.Infof("Bot running.")
	<-make(chan struct{})
}
