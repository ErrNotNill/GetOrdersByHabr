package GetOrdersByB24

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
)

func GenerateUrl(inn string) (url string) {
	convertedStr := fmt.Sprintf(`https://onviz.bitrix24.site/`, inn)
	return convertedStr
}

/*func Searcher(url string, keyword string) {
	newstr,urls, err := VisitMainPage(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	if strings.Contains(newstr, keyword) {
		fmt.Println("i found this word")
	}
}*/

func VisitMainPage(url string) (text string, urls []string, err error) {
	c := colly.NewCollector()

	urls = make([]string, 0)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		text = e.Text
		absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))
		urls = append(urls, absoluteURL)

		//fmt.Println(urls) //array with all HREFS
	})

	err = c.Visit(url)
	if err != nil {
		log.Println(err.Error())
		return "", nil, err
	}
	return text, urls, err
}

func VisitHrefs(url string) error {
	c := colly.NewCollector()
	//fmt.Println(VisitMainPage(url))
	_, getUrls, err := VisitMainPage(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, v := range getUrls {
		fmt.Println("VVVVVVVVVV", v)
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			fmt.Println(e.Text)
			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))
			absoluteTitle := e.Request.AbsoluteURL(e.ChildText("href"))
			fmt.Println(absoluteURL, absoluteTitle)
			if strings.Contains(e.Text, "регулировка") {
				fmt.Println("THIS FOUND: регулировка", absoluteURL)
				return
			}
			fmt.Println(e.Text)
			//fmt.Println(urls) //array with all HREFS
		})

		err = c.Visit(v)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		fmt.Println(v)

	}

	return err
}
