package scrapers

import (
	"log"
	"net/http"

	"github.com/raziael/scraper/peoplesrv"
	"github.com/schollz/closestmatch"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	baseURL = "https://www.truepeoplesearch.com"
)

//TruePeopleScraper implements PersonService for http://truepeoplesearch.com
type TruePeopleScraper struct{}

//GetPerson returns a person based on it's phone number, it multiple are found, uses the name to discriminate
func (TruePeopleScraper) GetPerson(phone string, name string) (*peoplesrv.Person, error) {
	p, err := scrHtml(phone, name)
	if err != nil {
		return nil, err
	}
	//Append address
	appendAddress(p)
	log.Println("Returning from " + baseURL)
	return p, nil

}

//scrHtml is an html scraper for truepeoplesearch
func scrHtml(phone string, name string) (*peoplesrv.Person, error) {
	// request and parse the front page
	root, err := getRoot(baseURL + "/results?phoneno=" + phone)
	if err != nil {
		panic(err)
	}

	// define a matcher
	matcher := func(n *html.Node) bool {
		// must check for nil values
		if n.DataAtom == atom.Div {
			return scrape.Attr(n.Parent.Parent.Parent, "class") == "card card-block shadow-form card-summary" &&
				scrape.Attr(n, "class") == "h4" &&
				n.Parent != nil && n.Parent.Parent != nil && n.Parent.Parent.Parent != nil
		}
		return false
	}

	// grab all articles and print them
	cards := scrape.FindAll(root, matcher)
	persons := make(map[string]*peoplesrv.Person)
	bagSizes := []int{1, 2, 3}

	for _, card := range cards {
		person := &peoplesrv.Person{}

		if err := parsePerson(person, card); err != nil {
			log.Print(err)
			continue
		}

		if person != nil {
			persons[person.Name] = person
		}

	}

	//Get map keys
	names := []string{}
	for k := range persons {
		names = append(names, k)
	}

	if len(names) > 1 {
		// Create a closestmatch object
		cm := closestmatch.New(names, bagSizes)
		return persons[cm.Closest(name)], nil
	}

	return persons[names[0]], nil
}

//AppendAddress adds the full address for the given person
func appendAddress(person *peoplesrv.Person) error {

	root, err := getRoot(baseURL + person.DetailURL)
	if err != nil {
		panic(err)
	}

	// define a matcher
	matcher := func(n *html.Node) bool {
		// must check for nil values
		if n.DataAtom == atom.A {
			return scrape.Attr(n, "class") == "link-to-more" &&
				scrape.Attr(n, "data-link-to-more") == "address"
		}
		return false
	}

	// grab all articles and print them
	cardDetails := scrape.FindAll(root, matcher)

	for _, cardDetail := range cardDetails {
		person.Address = scrape.Text(cardDetail)
		break
	}

	return nil

}

//Retrieves the document root element node from an html page
func getRoot(url string) (*html.Node, error) {

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	return root, nil
}

//Node to person transformer
func parsePerson(person *peoplesrv.Person, n *html.Node) error {
	person.DetailURL = scrape.Attr(n.Parent.Parent.Parent, "data-detail-link")
	person.Name = scrape.Text(n)

	return nil
}
