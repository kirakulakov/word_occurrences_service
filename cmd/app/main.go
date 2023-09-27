package main

import (
	"log"
	"npp_doslab/config"
	"npp_doslab/internal/app"
	fetcher "npp_doslab/pkg/fetchers"
)

const _fetch_frequency_seconds = 3

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Config error: %s", err)
	}

	// Run fetching
	go fetcher.RunFetching(_fetch_frequency_seconds, cfg)

	// Run server
	app.Run(cfg)

}
