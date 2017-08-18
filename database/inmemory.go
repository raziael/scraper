package database

import (
	"fmt"
	"log"
	"sync"
)

var (
	m map[string]*Person
)

func init() {
	m = make(map[string]*Person)
}

type inmemoryDatabase struct {
	mtx  sync.RWMutex
	next PeopleDatabase
}

//NewInmemoryDatabase returns a new inmemory database
func NewInmemoryDatabase(scraper PeopleDatabase) PeopleDatabase {
	return &inmemoryDatabase{
		next: scraper,
	}
}

//GetPerson returns a the person being looked at
func (ld inmemoryDatabase) GetPerson(phone string, name string) (*Person, error) {
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

//Update method not implemented
func (ld inmemoryDatabase) Update(person *Person) error {
	return fmt.Errorf("Unimplemented method")
}

//Delete method not implemented yet
func (ld inmemoryDatabase) Delete(phone string) {
	fmt.Errorf("Unimplemented method")
}
