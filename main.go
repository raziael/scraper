package main

import (
	"net/http"

	"github.com/raziael/scraper/database"
	"github.com/raziael/scraper/scrapers"
)

func main() {
	scraper := &scrapers.TruePeopleScraper{}
	service := database.NewInmemoryDatabase(scraper)

	http.HandleFunc("/", editHandler(service))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}
