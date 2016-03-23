package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	configFile  = flag.String("config", "", "Configuration file")
	redirectMap = make(map[string]string)
)

type Config struct {
	TLSConfigs   []TLSConfig `json:"tls"`
	Redirections []Redirect  `json:"redirections"`
}

type TLSConfig struct {
	CertFile string `json:"cert"`
	KeyFile  string `json:"key"`
}

type Redirect struct {
	URL         string `json:"url"`
	RedirectURL string `json:"redirect"`
}

func loadConfig() Config {
	if *configFile == "" {
		log.Fatal("Must specify configuration file")
	}

	var config Config

	f, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := json.NewDecoder(f)
	err = d.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func main() {
	flag.Parse()

	config := loadConfig()

	// Load all of the redirection configs
	for _, redirect := range config.Redirections {
		redirectMap[redirect.URL] = redirect.RedirectURL
	}

	http.HandleFunc("/", redirectHandler)

	var wg sync.WaitGroup
	wg.Add(2)

	go listenAndServe(&wg)
	go listenAndServeTLS(config, &wg)

	wg.Wait()
}
