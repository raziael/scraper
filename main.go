package main

import (
	"fmt"
	"net/http"

	"log"

	"github.com/schollz/closestmatch"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	baseUrl = "https://www.truepeoplesearch.com"
)

type TruePerson struct {
	DetailUrl string
	Name      string
	Address   string
}

func main() {
	person, err := getPerson("9083748806", "Wilson A Lopez")

	if err != nil {
		panic(err)
	}

	if err := appendAddress(person); err != nil {
		log.Println("Unable to get address for " + person.Name)
	}

	fmt.Println(person.Name)
	fmt.Println(person.Address)
}

func getPerson(phone string, name string) (*TruePerson, error) {
	// request and parse the front page
	root, err := getRoot(baseUrl + "/results?phoneno=" + phone)
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
	persons := make(map[string]*TruePerson)
	bagSizes := []int{1, 2, 3}

	for _, card := range cards {
		person := &TruePerson{}

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

func parsePerson(person *TruePerson, n *html.Node) error {
	person.DetailUrl = scrape.Attr(n.Parent.Parent.Parent, "data-detail-link")
	person.Name = scrape.Text(n)

	return nil
}

func appendAddress(person *TruePerson) error {

	root, err := getRoot(baseUrl + person.DetailUrl)
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
