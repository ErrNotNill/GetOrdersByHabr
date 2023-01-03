package services

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"parser/models"
	"strings"
)

func GenerateUrl(inn string) (url string) {
	convertedStr := fmt.Sprintf("https://www.rusprofile.ru/search?query=%v&type=ul", inn)
	return convertedStr
}

func VisitRusprofile(inn string) (*models.Company, error) {

	c := colly.NewCollector()

	var name string
	var kpp string
	var boss string

	//get Company-KPP
	c.OnHTML(".company-info__text.has-copy span[id=clip_kpp]", func(e *colly.HTMLElement) {
		kpp = e.Text
	})

	//get Company-Name
	c.OnHTML(".company-header__row h1[itemprop=name]", func(e *colly.HTMLElement) {
		name = strings.Trim(e.Text, " \t\n")
	})

	//parse Company-Boss
	c.OnHTML(".company-row.hidden-parent", func(e *colly.HTMLElement) {
		boss = e.ChildText(".company-info__text")
	})

	err := c.Visit(GenerateUrl(inn))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	company := &models.Company{Inn: inn, Kpp: kpp, Name: name, Boss: boss}

	return company, err
}
