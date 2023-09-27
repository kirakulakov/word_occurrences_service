package main

import (
	"fmt"
	"log"
	"npp_doslab/config"
	"npp_doslab/internal/app"
	fetcher "npp_doslab/pkg/fetchers"
)

const _fetchInterval = 10

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Config error: %s", err))
	}

	// Fetching
	go fetcher.RunFetching(_fetchInterval, cfg)

	// Server
	app.Run(cfg)

}
