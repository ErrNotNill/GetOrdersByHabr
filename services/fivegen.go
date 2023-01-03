package services

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
)

const TryUr = `https://пятаяпередача.рф/manufacturers`

func StartFunc() {
	err := FiveGen(TryUr)
	if err != nil {
		log.Println("smth wrgng")
		return
	}
}

func FiveGen(url string) error {
	c := colly.NewCollector()

	//fmt.Println(link)
	//here we got link for the next (deeper) scraping
	c.OnHTML(".literal-manufacturers a[href]", func(e *colly.HTMLElement) {
		hrefs := e.Attr("href")
		fmt.Println(hrefs)
	})
	//here we visit start url from main
	err := c.Visit(url)
	if err != nil {
		log.Println("Error visit url")
		//log.Prefix()
	}
	return err
}
