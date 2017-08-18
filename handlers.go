package main

import (
	"html/template"
	"net/http"

	"github.com/raziael/scraper/database"
)

func editHandler(dbsrv database.PeopleDatabase) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		er := r.ParseForm()
		if er != nil {
			panic(er)
		}

		p := &database.Person{}
		var err error
		if r.Method == "POST" {

			p, err = dbsrv.GetPerson(r.FormValue("phone_number"), r.FormValue("person_name"))
			if err != nil {
				panic(err)
			}
		}
		t, _ := template.ParseFiles("web/static/index.html")
		t.Execute(w, p)

	}

}
