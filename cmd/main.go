package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/patrickmn/go-cache"
	"log"
	"parser/scraper"
	"parser/urls"
	"strings"
	"time"
)

// KeyWord for searching in titles from orders.
var KeyWord string

const UrlDb = `postgres://postgres:postgres@localhost:5432/postgres`
const Url = `https://onviz.bitrix24.site/`

const Token = "1953280162:AAFMVzq63WHhr_KkNjwGgObHbI4PbQcmQqg"

// Cache here array must contain all hrefs. And this func must check that key is.
// If key founded, cache delete all keys from hrefs. But if not founded, it added new key.
func Cache(key string, array []string) {
	c := cache.New(5*time.Minute, 10*time.Minute)
	for _, v := range array {
		foo, found := c.Get(v)
		if found {
			c.Delete(v)
			fmt.Println("no")
			fmt.Println(foo)
		} else {
			c.Set(v, &scraper.Habr{}, cache.NoExpiration)
		}
	}
	foo, found := c.Get(key)
	if found {
		c.Delete(key)
		fmt.Println("yes")
		fmt.Println(foo)
	}
}

func parseHabr(jsonBuffer []byte) ([]scraper.Habr, error) {
	// We create an empty array
	users := []scraper.Habr{}

	// Unmarshal the json into it. this will use the struct tag
	err := json.Unmarshal(jsonBuffer, &users)
	if err != nil {
		return nil, err
	}

	// the array is now filled with users
	return users, nil
}

var Input string
var Href string

var ArrayWithoutDuplicates = make([]string, 0)

var bigFullArray = make([]string, 0)

var lastestArray = make([]string, 0)

func ifPrefix(absoluteURL string) (b bool) {

	switch false {
	case strings.HasPrefix(absoluteURL, "https://www"):
		b = false
	case strings.HasPrefix(absoluteURL, "https://b24"):
		b = false
	case strings.HasPrefix(absoluteURL, "mailto"):
		b = false
	case strings.HasPrefix(absoluteURL, "tel"):
		b = false
	case strings.HasPrefix(absoluteURL, "https://b24"):
		b = false
	case strings.HasPrefix(absoluteURL, "https://b24"):
		b = false

	}
	return true

}

type Hrefs struct {
	href string
}

func checkForValue(userValue int, students map[string]int) string {
	//traverse through the map
	for key, value := range students {
		//check if present value is equals to userValue
		if value == userValue {
			//if same return true
			return key
		}
	}
	//if value not found return false
	return ""
}

/*func SearchEverywhere(sites map[string]int) {
	c := colly.NewCollector()
	for k, _ := range sites {

	}
}*/

func main() {

	//var newMap = make(map[string]string)

	//collectionFoundedHrefs := make([]string,0)

	newFullArr := make([]string, 0)

	for k, _ := range urls.IndexSites {

		var fullData string

		//firstSite := checkForValue(1, urls.IndexSites)

		c := colly.NewCollector()
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			//	fmt.Println(e.Text)
			//fmt.Println(e.Text, "VisitMainPAGE")
			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))
			//fmt.Println(fullData)
			fullData = e.Text
			//fmt.Println(fullData)

			founded := searchString(fullData, "Сборка")
			if founded == true {

				fullDataPlus := fmt.Sprintf("%v %v", fullData, absoluteURL)
				//newMap[fullData] = absoluteURL
				newFullArr = append(newFullArr, fullDataPlus)
				//fmt.Println(fullData, "FULLDATA")
				//fmt.Println(absoluteURL)
				//return
			} else {
				//count = count + 1
				//return
				//firstSite
				//fmt.Println("no keyword here")
				//return
			}

		})

		err := c.Visit(k)
		if err != nil {
			log.Println(err.Error())
		}
	}

	for _, v := range newFullArr {
		fmt.Println(v)
	}

}

func search(docs []string, term string) []string {
	var r []string
	for _, doc := range docs {
		if strings.Contains(doc, term) {
			r = append(r, doc)
			break
		}
	}
	return r
}
func searchString(fullText string, keyword string) bool {
	if strings.Contains(fullText, keyword) {
		return true
	}
	return false
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func VisitLastHrefs(lA []string) (text string, urls []string, err error) {
	c := colly.NewCollector()

	for _, vUrl := range lA {

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))
			urls = append(urls, absoluteURL)
			//fmt.Println(urls) //array with all HREFS
		})

		err = c.Visit(vUrl)
		if err != nil {
			log.Println(err.Error())
			return "", nil, err
		}
		return text, urls, err
	}
	urls = make([]string, 0)
	return "", urls, err
}

func VisitNextHrefs() {
	c := colly.NewCollector()

	for _, v := range ArrayWithoutDuplicates {

		err := c.Visit(v)
		if err != nil {
			log.Println(err.Error())
			//		fmt.Println("error visit URL")
			continue
		}

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {

			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

			ArrayWithoutDuplicates = append(ArrayWithoutDuplicates, absoluteURL)
		})

	}

}

/*func LastMain(){
	var fullData string
	var count = 1
	firstSite := checkForValue(count, urls.IndexSites)
	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//	fmt.Println(e.Text)
		//fmt.Println(e.Text, "VisitMainPAGE")
		absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

		fullData = e.Text
		//fmt.Println(fullData)

		founded := searchString(fullData, "Сборка")
		if founded == true {
			fmt.Println(fullData, "FULLDATA")
			fmt.Println(absoluteURL)
		} else {
			count = count + 1
			//firstSite
			//fmt.Println("no keyword here")
		}

	})

	err := c.Visit(firstSite)
	if err != nil {
		log.Println(err.Error())
	}

	//HERE is Scan of all hrefs
	for v, _ := range urls.IndexSites {
		c := colly.NewCollector()

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			//	fmt.Println(e.Text)
			//fmt.Println(e.Text, "VisitMainPAGE")
			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

			fullData = e.Text
			//fmt.Println(fullData)

			founded := searchString(fullData, "Сборка")
			if founded == true {
				fmt.Println(fullData, "FULLDATA")
				fmt.Println(absoluteURL)
			} else {
				fmt.Println("no keyword here")
			}

		})

		err := c.Visit(v)
		if err != nil {
			log.Println(err.Error())
		}
	}

	//HERE all hrefs in Massive
	newArr := make([]string, 0)
	for v, _ := range urls.IndexSites {

		c := colly.NewCollector()

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			//	fmt.Println(e.Text)
			//fmt.Println(e.Text, "VisitMainPAGE")
			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

			if strings.HasPrefix(absoluteURL, `https://onviz`) {
				newArr = append(newArr, absoluteURL)
			}

			//fmt.Println(absoluteURL)
			//fmt.Println(urls) //array with all HREFS
		})

		//withoutDupls := removeDuplicateStr(newArr)

		founded := search(newArr, "alisa")

		//foundedWithoutDuplicates := removeDuplicateStr(founded)

		//newF := removeDuplicateStr(founded)

		fmt.Println(founded)

		/*
			for _, v := range withoutDupls {
				fmt.Println(v)
			}

		err := c.Visit(v)
		if err != nil {
			log.Println(err.Error())
		}
	}

	fmt.Println(len(newArr))
}
*/

/*func OldMain(){
	Input = "Название заголовка" //пишем своё
	Href = "ID ссылки задания"   //добавляем в нужное место
	p := new(repository.Postgres)
	p.GetOrderByHref(Href) //получем инфу по вставленной ссылке (title,descr,value...)

	repository.PostgresInit(UrlDb)
	//p := new(repository.Postgres)

	page := utils.GenPages(Url)
	checkDuplicate := make([]string, 0)
	//hrefs := scraper.UpdateHrefs(Url)
	//firstHref := hrefs[0]
	scraper.GenHrefs(Url)
	for _, v := range page {
		hrefs := scraper.GenHrefs(v)
		for _, v := range hrefs {
			checkDuplicate = append(checkDuplicate, v)
		}
	}

	var ordersFromHabr string
	var orders = make([]string, 0) //there our orders in hybrid json format todo edit correct formats, example to share (1 - array, 2 - json)
	fmt.Println(orders)
	//ords := make([]scraper.Habr, 0)

	//fmt.Println(ords)

	pages := utils.GenPages(Url)
	var count int
	for _, v := range pages {
		count++
		habr, err := scraper.HabrScraper(v, Input)
		if err != nil {
			log.Println("err scrap on page: ", count)
		}
		var m scraper.Habr
		ja, _ := json.Marshal(habr)
		json.Unmarshal(ja, &m)
		ordersFromHabr = string(ja)
		orders = append(orders, ordersFromHabr)
		//ords, err = parseHabr(ja)
		//fmt.Println(orders)
	}

	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}
	newbot := telegram.NewBot(bot)

	if err = newbot.StartBot(orders); err != nil {
		log.Fatal(err)
	}
}*/

/*var golang = []string{
	"go", "Golang", "golang",
}
*/
//fmt.Println(scraper.NewMassive)
//convert,_ := strconv.Atoi("2609257031")

/*nums := []string{
"	2609257031",
	2609257031
}*/

/*habr2, err := scraper.HabrScraper(Url2, "Go")
if err != nil {
	log.Prefix()
	return
}
var mm scraper.Habr
ja2, _ := json.Marshal(habr2)
json.Unmarshal(ja2, &mm)
fmt.Println(string(ja2))*/

/*func oldMain(){
	err := GetOrdersByB24.VisitHrefs(Url)

	if err != nil {
		fmt.Println(err.Error())
		//fmt.Println("whats going on?")
		return
	}
	//fmt.Println("FULL ARRAY", GetOrdersByB24.FullArray)
	withoutDuplicates := removeDuplicateStr(GetOrdersByB24.FullArray)
	for _, v := range withoutDuplicates {
		if strings.HasPrefix(v, "https://onviz.bitrix24") {
			ArrayWithoutDuplicates = append(ArrayWithoutDuplicates, v)
			bigFullArray = append(bigFullArray, v)
		}
	}

	//	for _, v := range ArrayWithoutDuplicates {
	//	fmt.Println(v)
	//}

	//fmt.Println(ArrayWithoutDuplicates)

	VisitNextHrefs()

	rmvDupls := removeDuplicateStr(ArrayWithoutDuplicates)
	//urlsNew := make([]string, 0)
	for _, v := range rmvDupls {

		c := colly.NewCollector()

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			//absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))
			//urlsNew = append(urlsNew, absoluteURL)
			if strings.HasPrefix(v, "https://onviz.bitrix24") {
				//urlsNew = append(urlsNew, v)
				bigFullArray = append(bigFullArray, v)

			}

		})

		err = c.Visit(v)
		if err != nil {
			log.Println(err.Error())
		}
		//fmt.Println(urlsNew)
	}

	//newRemovedDupls := removeDuplicateStr(urlsNew)

	//for _, v := range newRemovedDupls {
	//fmt.Println(v)
	//}

	removed := removeDuplicateStr(bigFullArray)
	for _, v := range removed {

		lastestArray = append(lastestArray, v)

		//fmt.Println(v)
	}
	for _, v := range lastestArray {
		fmt.Println(v)
	}

	fmt.Println(len(lastestArray))

	fmt.Println("end")
}*/
