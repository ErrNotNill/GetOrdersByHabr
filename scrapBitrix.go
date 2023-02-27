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

var FullArray = make([]string, 0)

func VisitMainPage(url string) (text string, urls []string, err error) {
	c := colly.NewCollector()

	urls = make([]string, 0)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		text = e.Text
		absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))
		urls = append(urls, absoluteURL)
		FullArray = append(FullArray, absoluteURL)
		//fmt.Println(urls) //array with all HREFS
	})

	err = c.Visit(url)
	if err != nil {
		log.Println(err.Error())
		return "", nil, err
	}
	return text, urls, err
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func VisitHrefs(url string) error {
	c := colly.NewCollector()
	//fmt.Println(VisitMainPage(url))
	_, getUrls, err := VisitMainPage(url)
	if err != nil {
		//fmt.Println("error visit Main Page")
		fmt.Println(err.Error())
	}
	//	fmt.Println("GET URLS", getUrls)
	for _, v := range getUrls {

		err = c.Visit(v)
		if err != nil {
			log.Println(err.Error())
			//		fmt.Println("error visit URL")
			continue
		}
		fmt.Println(v)

		//	fmt.Println("error here")
		//	fmt.Println("VVVVVVVVVV", v)
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			//	fmt.Println(e.Text)
			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

			if strings.HasPrefix(v, "https://onviz.bitrix24") {
				FullArray = append(FullArray, absoluteURL)
			}

			//	absoluteTitle := e.Request.AbsoluteURL(e.ChildText("href"))
			//	fmt.Println(absoluteURL, absoluteTitle)
			if strings.Contains(e.Text, "регулировка") {
				//		fmt.Println("THIS FOUND: регулировка", absoluteURL)
				return
			}
			//	fmt.Println(e.Text)
			//fmt.Println(urls) //array with all HREFS
		})

	}

	return err
}
