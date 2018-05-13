package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const (
	csvName = "coffee_reviews.csv"
)

// Scraper ...
type Scraper struct {
	CurrentPage int64
	LastPage    int64
}

// GetPages ...
func (s *Scraper) GetPages(span string) {
	matches := FindInts(span)
	current, err := strconv.ParseInt(matches[0], 10, 0)
	last, err := strconv.ParseInt(matches[1], 10, 0)
	if err != nil {
		log.Fatalf("cannot parse last page: %s\n", err)
	}
	s.CurrentPage = current
	s.LastPage = last
}

// FindInts ...
func FindInts(str string) []string {
	re := regexp.MustCompile(`\d+`)
	return re.FindAllString(str, 2)
}

// GetDate ...
func GetDate(s string) string {
	// Don't split on ": " just in case there's a typo
	split := strings.Split(s, ":")
	monthYear := strings.Split(
		strings.Trim(split[1], " "),
		" ",
	)
	// insert "01" because reviews are only dated Month / Year ¯\_(ツ)_/¯
	return strings.Join([]string{monthYear[0], " 01, ", monthYear[1]}, "")
}

func main() {
	s := Scraper{
		CurrentPage: -1,
		LastPage:    0,
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.coffeereview.com"),
	)

	/*
		bootstrap CSV
	*/
	file, err := os.Create(csvName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", csvName, err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"rating",
		"coffee",
		"roaster",
		"review_link",
		"roaster_link",
		"excerpt",
		"review_date",
	})

	/*
		decipher the number of pages from the <span class="pages"> element:
		   <span class="pages">Page 1 of 244</span> -> []int64{1, 244}
	*/
	c.OnHTML("div.wp-pagenavi > span.pages", func(e *colly.HTMLElement) {
		s.GetPages(e.Text)
		log.Printf("Visiting page %v of %v\n", s.CurrentPage, s.LastPage)
	})
	/*
		parse review preview attributes and write to csv
	*/
	c.OnHTML(".review-content", func(e *colly.HTMLElement) {
		var (
			rating      = e.ChildText(".review-col1 > .review-rating")
			coffee      = e.ChildText(".review-col1 > .review-title > a")
			roaster     = e.ChildText(".review-col1 > h3")
			reviewLink  = e.Request.AbsoluteURL(e.ChildAttr(".review-col1 > .review-title > a", "href"))
			roasterLink = e.Request.AbsoluteURL(e.ChildAttr(".links > .right > a", "href"))
			excerpt     = e.ChildText(".excerpt > p")
			reviewDate  = GetDate(e.ChildText(".review-col2 > p:first-child"))
		)
		writer.Write([]string{
			rating,
			coffee,
			roaster,
			reviewLink,
			roasterLink,
			excerpt,
			reviewDate,
		})
	})
	/*
		paginate
	*/
	c.OnHTML("a.nextpostslink", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.Visit("http://www.coffeereview.com/review/")
	// Log out collector's statistics
	log.Printf("Finished, wrote entries to %q\n", csvName)
	log.Println(c)
}
