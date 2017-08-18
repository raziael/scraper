package database

import (
	"log"
	"sync"

	"github.com/raziael/scraper/peoplesrv"
)

var (
	m map[string]*peoplesrv.Person
)

func init() {
	m = make(map[string]*peoplesrv.Person)
}

type inmemoryDatabase struct {
	mtx  sync.RWMutex
	next peoplesrv.PersonService
}

//NewInmemoryDatabase returns a new inmemory database
func NewInmemoryDatabase(scraper peoplesrv.PersonService) peoplesrv.PersonService {
	return &inmemoryDatabase{
		next: scraper,
	}
}

//GetPerson returns a the person being looked at
func (ld inmemoryDatabase) GetPerson(phone string, name string) (*peoplesrv.Person, error) {
	ld.mtx.RLock()
	defer ld.mtx.RUnlock()
	if p, ok := m[phone]; ok {
		log.Println("Returning from inmemory DB")
		return p, nil
	}

	p, err := ld.next.GetPerson(phone, name)
	if err != nil {
		log.Println("Returning from inmemory DB")
	} else {
		m[phone] = p
	}

	return p, nil
}
