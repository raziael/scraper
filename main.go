package main

import (
	"log"
	"net/http"

	"github.com/raziael/scraper/database"
	"github.com/raziael/scraper/scraper"
)

func main() {
	scraper := &scraper.TruePeopleScraper{}
	dbsrv := database.NewInmemoryDatabase()

	http.HandleFunc("/", searchPersonHandler(dbsrv, scraper))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/resources"))))

	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}

}
