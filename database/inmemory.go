package database

import (
	"fmt"
	"log"
	"sync"
)

//Person Represents the person being looked at
type Person struct {
	DetailURL string `json:"url"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

//InmemoryDatabase holds a reference of previously scrapped elements
type InmemoryDatabase struct {
	mtx *sync.RWMutex
}

var (
	m map[string]*Person
)

func init() {
	m = make(map[string]*Person)
}

//NewInmemoryDatabase returns a pointer
func NewInmemoryDatabase() *InmemoryDatabase {
	immdb := &InmemoryDatabase{}
	immdb.mtx = &sync.RWMutex{}
	return immdb
}

//FindOne extends inmemory database.
//Returns a person if previously stored, otherwise returns an error
func (ld InmemoryDatabase) FindOne(phone string, name string) (*Person, error) {
	ld.mtx.RLock()
	defer ld.mtx.RUnlock()
	if p, ok := m[phone]; ok {
		log.Println("Returning from inmemory DB")
		return p, nil
	}

	return nil, fmt.Errorf("Record not found in database")
}

//Update method not implemented
func (ld InmemoryDatabase) Update(person *Person) error {
	ld.mtx.Lock()
	defer ld.mtx.Unlock()

	if person == nil {
		return fmt.Errorf("Person cannot be nil")
	}

	m[person.Phone] = person
	return nil
}
