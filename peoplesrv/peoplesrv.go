package peoplesrv

//PersonService provides operations on strings.
type PersonService interface {
	GetPerson(string, string) (*Person, error)
}

//Person Represents the person being looked at
type Person struct {
	DetailURL string `json:"url"`
	Name      string `json:"name"`
	Address   string `json:"address"`
}
