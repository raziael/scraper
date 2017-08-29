package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/raziael/scraper/database"
	"github.com/raziael/scraper/scraper"
)

//Response to the html
type Response struct {
	Person  *database.Person
	Message string
}

func searchPersonHandler(dbsrv *database.InmemoryDatabase, scraper scraper.Scraper) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		er := r.ParseForm()
		if er != nil {
			log.Println("Unable to parse form. " + er.Error())
		}

		p := &database.Person{}
		//Preparing response
		resp := &Response{}
		var err error
		if r.Method == "POST" {
			//Compile regex to unmask phone
			reg, rxerr := regexp.Compile("[^a-zA-Z0-9]+")
			if rxerr != nil {
				log.Println(rxerr)
			}

			//get input from form
			maskedPhone := r.FormValue("phone_number")
			phone := reg.ReplaceAllString(maskedPhone, "")
			name := r.FormValue("person_name")

			p, err = dbsrv.FindOne(phone, name)
			if err != nil {
				log.Println(err.Error())
				p, err = scraper.ScrapAndGet(phone, name)
				if err != nil {
					resp.Message = err.Error()
				} else {
					log.Println("Found in scraper, updating database")
					dbsrv.Update(p)
				}

			}
			resp.Person = p
		}

		//Parsing template
		t, errTemplate := template.ParseFiles("web/index.html")
		if errTemplate != nil {
			http.Error(w, errTemplate.Error(), http.StatusInternalServerError)
			return
		}

		t.Execute(w, resp)
	}

}
