package scraper

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"strconv"
	"strings"
)

type Habr struct {
	Href        int    `json:"href,omitempty"`
	Title       string `json:"title,omitempty"`
	Url         string `json:"url,omitempty"`
	Finance     string `json:"finance,omitempty"`
	Description string `json:"description,omitempty"`
}

// ConvertNumsOfHrefs получаем число из аттрибута <a href>
// в итоге, пример /tasks/242424 == 242424
// делаем это для нахождения и отсеивания дубликатов. Дубли в БД не добавляются. Соответственно
// не должны браться как новые заказы. Либо проверяются на уровне функции.
func ConvertNumsOfHrefs(s string) int {
	//str := "/tasks/242424"
	trimmed := strings.Trim(s, "/tasks/")
	convStr, _ := strconv.Atoi(trimmed)
	return convStr
}

// GenHrefs генерирует по аттрибуту <a href> все ссылочные номера /tasks/nums
// в итоге получаем массив /tasks/nums со всеми ссылочными номерами
func GenHrefs(url string) []string {
	c := colly.NewCollector()
	hrefs := make([]string, 0)

	c.OnHTML(".task__title a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		hrefs = append(hrefs, href)
	})
	err := c.Visit(url)
	if err != nil {
		log.Println("Error visit url")
	}
	return hrefs
}

// GenLinks генерирует конечные ссылки по всем заказам (из аттрибутов <a href>) /tasks/nums
// в итоге получаем массив всех конечных ссылок (по всем заказам с аттрибутом <a href>)
func GenLinks(url string) []string {
	c := colly.NewCollector()
	var link string
	links := make([]string, 0)
	c.OnHTML(".task__title a[href]", func(e *colly.HTMLElement) {
		link = fmt.Sprintf("https://freelance.habr.com%v", e.Attr("href"))
		links = append(links, link)
	})
	err := c.Visit(url)
	if err != nil {
		log.Println("Error visit url")
	}
	return links
}

func HabrScraper(url string, sub string) (*Habr, error) {
	c := colly.NewCollector()

	var href int
	var link string
	var title string
	var finance string
	var description string

	//fmt.Println(link)
	//here we got link for the next (deeper) scraping
	c.OnHTML(".task__title a[href]", func(e *colly.HTMLElement) {
		hrefs := e.Attr("href")
		//href = ConvertNumsOfHrefs(hrefs)
		if strings.Contains(e.Text, sub) {
			title = e.Text
			fmt.Println(title)
			link = fmt.Sprintf("https://freelance.habr.com%v", hrefs)
			fmt.Println(link, "LINK is")
			//fmt.Println(link)
		}
	})
	//here we visit start url from main
	err := c.Visit(url)
	if err != nil {
		log.Println("Error visit url")
		//log.Prefix()
	}

	//Finance
	c.OnHTML(".task__finance span[class=count]", func(e *colly.HTMLElement) {
		finance = e.Text
		fmt.Println(finance)
	})
	//Description
	c.OnHTML(".task__description", func(e *colly.HTMLElement) {
		s := fmt.Sprintf("\t%v\n", e.Text)
		t := strings.TrimSpace(s)
		description = t
		fmt.Println(description)
	})

	err = c.Visit(link)
	if err != nil {
		log.Println("Error visit link OR not found your query")
		return nil, err
	}
	habr := &Habr{Href: href, Title: title, Url: link, Finance: finance, Description: description}
	return habr, err
}

/*func HabrScraperJson(url string, sub string) (*Habr, error) {
	c := colly.NewCollector()

	var link string
	var title string
	var finance string
	var description string

	fmt.Println(link)
	//here we got link for the next (deeper) scraping
	c.OnHTML(".task__title a[href]", func(e *colly.HTMLElement) {

		if strings.Contains(e.Text, sub) {
			title = e.Text
			link = fmt.Sprintf("https://freelance.habr.com%v", e.Attr("href"))
			fmt.Println(link)
		}
	})
	//here we visit start url from main
	err := c.Visit(url)
	if err != nil {
		log.Println("Error visit url")
		//log.Prefix()
	}

	//Finance
	c.OnHTML(".task__finance span[class=count]", func(e *colly.HTMLElement) {
		finance = e.Text
	})
	//Description
	c.OnHTML(".task__description", func(e *colly.HTMLElement) {
		s := fmt.Sprintf("\t%v\n", e.Text)
		t := strings.TrimSpace(s)
		description = t
	})

	err = c.Visit(link)
	if err != nil {
		log.Println("Error visit link OR not found your query")
		log.Prefix()
		return nil, err
	}
	habr := &Habr{Title: title, Url: link, Finance: finance, Description: description}
	return habr, err
}*/
