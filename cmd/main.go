package main

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"log"
	"parser/repository"
	"parser/scraper"
	"parser/utils"
	"time"
)

// KeyWord for searching in titles from orders.
var KeyWord string

const UrlDb = `postgres://postgres:postgres@localhost:5432/postgres`
const Url = `https://freelance.habr.com/tasks?categories=development_all_inclusive%2Cdevelopment_backend%2Cdevelopment_frontend%2Cdevelopment_prototyping%2Cdevelopment_ios%2Cdevelopment_android%2Cdevelopment_desktop%2Cdevelopment_bots%2Cdevelopment_games%2Cdevelopment_1c_dev%2Cdevelopment_scripts%2Cdevelopment_voice_interfaces%2Cdevelopment_other`

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

func main() {

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

	/*bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}
	newbot := telegram.NewBot(bot)

	if err = newbot.StartBot(orders); err != nil {
		log.Fatal(err)
	}*/
}

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
