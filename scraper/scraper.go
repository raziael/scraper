package scraper

import "github.com/raziael/scraper/database"

//Scraper represents an html scraper / parser
type Scraper interface {
	ScrapAndGet(string, string) (*database.Person, error)
}
