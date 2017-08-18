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
func NewInmemoryDatabase() PeopleDatabase {
	return &inmemoryDatabase{}
}

//GetPerson returns a the person being looked at
func (ld inmemoryDatabase) GetPerson(phone string, name string) (*Person, error) {
	ld.mtx.RLock()
	defer ld.mtx.RUnlock()
	if p, ok := m[phone]; ok {
		log.Println("Returning from inmemory DB")
		return p, nil
	}

	return nil, fmt.Errorf("Record not found")
}

//Update method not implemented
func (ld inmemoryDatabase) Update(person *Person) error {
	ld.mtx.Lock()
	defer ld.mtx.Unlock()

	if person == nil {
		return fmt.Errorf("Person cannot be nil")
	}

	m[person.Phone] = person
	return nil
}

//Delete method not implemented yet
func (ld inmemoryDatabase) Delete(phone string) {
	fmt.Errorf("Unimplemented method")
}
