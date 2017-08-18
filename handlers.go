package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/raziael/scraper/database"
	"github.com/raziael/scraper/scraper"
)

func searchPersonHandler(dbsrv database.PeopleDatabase, scraper scraper.Scraper) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		er := r.ParseForm()
		if er != nil {
			panic(er)
		}

		p := &database.Person{}
		var err error
		if r.Method == "POST" {
			//get input from form
			phone := r.FormValue("phone_number")
			name := r.FormValue("person_name")

			p, err = dbsrv.GetPerson(phone, name)
			if err != nil {
				log.Println("Not found in DB, scraping...")
				p, err = scraper.GetPerson(phone, name)
				if err != nil {
					log.Println("Unable to find person")
				} else {
					log.Println("Found in scraper, updating database")
					dbsrv.Update(p)
				}
			} else {
				log.Println("Found in DB")
			}
		}
		t, _ := template.ParseFiles("web/static/index.html")
		t.Execute(w, p)

	}

}
