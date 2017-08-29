package scraper

import (
	"fmt"
	"log"
	"net/http"

	"github.com/raziael/scraper/database"
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

//ScrapAndGet returns a person based on it's phone number, it multiple are found, uses the name to discriminate
func (TruePeopleScraper) ScrapAndGet(phone string, name string) (*database.Person, error) {
	log.Println("Scrapping from " + baseURL)
	p, err := scrHTML(phone, name)
	if err != nil {
		return nil, err
	}
	//Append address
	if p != nil {
		log.Println("Searching for address...")
		err = appendAddress(p)
		if err != nil {
			return nil, err
		}
		return p, nil
	}

	return p, nil
}

//scrHtml is an html scraper for truepeoplesearch
func scrHTML(phone string, name string) (*database.Person, error) {
	// request and parse the front page
	root, err := getRoot(baseURL + "/results?phoneno=" + phone)
	if err != nil {
		return nil, err
	}

	//Validating that contains data
	//if root.Data == "" {
	//	return nil, fmt.Errorf("No data found in scrapper")
	//}

	// define a matcher
	matcher := func(n *html.Node) bool {
		// must check for nil values
		if n.DataAtom == atom.Div {
			if n.Parent != nil && n.Parent.Parent != nil && n.Parent.Parent.Parent != nil {
				return scrape.Attr(n.Parent.Parent.Parent, "class") == "card card-block shadow-form card-summary" &&
					scrape.Attr(n, "class") == "h4"
			}
			log.Println("Person not found...")
		}
		return false
	}

	// grab all articles and print them
	cards := scrape.FindAll(root, matcher)
	persons := make(map[string]*database.Person)
	bagSizes := []int{1, 2, 3}

	for _, card := range cards {
		person := &database.Person{}
		if err := parsePerson(phone, person, card); err != nil {
			log.Print(err)
			continue
		}

		if person != nil {
			persons[person.Name] = person
		}

	}

	if len(persons) == 0 {
		return nil, fmt.Errorf(name + " not found [" + phone + "]")
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
func appendAddress(person *database.Person) error {

	root, err := getRoot(baseURL + person.DetailURL)
	if err != nil {
		return err
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
		log.Println("Address Scraped " + person.Address)
		break
	}

	return nil

}

//Retrieves the document root element node from an html page
func getRoot(url string) (*html.Node, error) {

	resp, err := http.Get(url)
	//If unable to read, return with error
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)

	return root, err
}

//Node to person transformer
func parsePerson(phone string, person *database.Person, n *html.Node) error {
	person.DetailURL = scrape.Attr(n.Parent.Parent.Parent, "data-detail-link")
	person.Name = scrape.Text(n)
	person.Phone = phone

	return nil
}
