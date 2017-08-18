package main

import (
	"html/template"
	"net/http"

	"github.com/raziael/scraper/peoplesrv"
)

func editHandler(service peoplesrv.PersonService) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		er := r.ParseForm()
		if er != nil {
			panic(er)
		}

		p := &peoplesrv.Person{}
		var err error
		if r.Method == "POST" {

			p, err = service.GetPerson(r.FormValue("phone_number"), r.FormValue("person_name"))
			if err != nil {
				panic(err)
			}
		}
		t, _ := template.ParseFiles("web/static/index.html")
		t.Execute(w, p)

	}

}
