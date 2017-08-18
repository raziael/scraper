package database

//PersonService provides operations on strings.
type PeopleDatabase interface {
	GetPerson(string, string) (*Person, error)
	Update(*Person) error
	Delete(string)
}

//Person Represents the person being looked at
type Person struct {
	DetailURL string `json:"url"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}
