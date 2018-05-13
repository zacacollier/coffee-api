package main_test

import (
	"strconv"

	"github.com/gocolly/colly"
	. "github.com/zacacollier/CoffeeAPI"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var s = Scraper{
	CurrentPage: -1,
	LastPage:    0,
}

var _ = Describe("Main", func() {

	Describe("GetLastPage", func() {
		It("should get the current last page", func() {
			c := colly.NewCollector(
				colly.AllowedDomains("www.coffeereview.com"),
			)
			/*
				decipher the number of pages from the <span class="pages"> element:
				   <span class="pages">Page 1 of 244</span>
			*/
			c.OnHTML("div.wp-pagenavi > span.pages", func(e *colly.HTMLElement) {
				matches := FindInts(e.Text)
				actualCurrentPage, _ := strconv.ParseInt(matches[0], 10, 0)
				actualLastPage, _ := strconv.ParseInt(matches[1], 10, 0)
				s.GetPages(e.Text)
				Expect(s.CurrentPage).To(Equal(actualCurrentPage))
				Expect(s.LastPage).To(Equal(actualLastPage))
			})

			c.Visit("http://www.coffeereview.com/review/")
		})
	})
	/*
		Describe("GetDate", func() {
			It("should parse a date that's extracted from the crawler", func() {
				c := colly.NewCollector(
					colly.AllowedDomains("www.coffeereview.com"),
				)
				c.OnHTML(".review-content", func(e *colly.HTMLElement) {
					rawDate := e.ChildText(".review-col2 > p:first-child")
					actualReviewDate := GetDate(rawDate)
					var t interface{} = actualReviewDate
					_, ok := t.(time.Time)
					Expect(ok).To(Equal(true))
				})

				c.Visit("http://www.coffeereview.com/review/")
			})
		})
	*/

})
