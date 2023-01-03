package services

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"parser/models"
)

//Todo, avito now does not work. For good work need to use selenium

// Example
func GenerateUrlAvito(num string) (url string) {
	convertedStr := fmt.Sprintf("https://www.avito.ru/novosibirsk/vakansii/gornichnaya_v_apartamenty_posutochnyh_kvartir_%v", num)
	return convertedStr
}

func VisitAvito(url string) (*models.Resume, error) {
	c := colly.NewCollector()

	var position string
	var company string
	var salary string
	var schedule string
	var description string
	var link string

	//get Position
	c.OnHTML(".h3", func(e *colly.HTMLElement) {
		position = e.Text
		//kpp = e.Text
	})

	//KPP
	c.OnHTML(".company-info__text.has-copy span[id=clip_kpp]", func(e *colly.HTMLElement) {
		//kpp = e.Text
	})
	//<div class="style-title-info-main-_sKj0"><h1 class="style-title-info-title-eHW9V style-
	//title-info-title-text-CoxZd"><span itemprop="name" class="title-info-title-text"
	//data-marker="item-view/title-info">Горничная в апартаменты посуточных квартир</span></h1></div>

	//NAME <div class .company-header__row> => <h1 itemprop="name">
	c.OnHTML(".task__title", func(e *colly.HTMLElement) {
		position = e.Text
	})

	//parse Company-Boss
	c.OnHTML(".company-row.hidden-parent", func(e *colly.HTMLElement) {
		//boss = e.ChildText(".company-info__text")
	})

	err := c.Visit(url)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	resume := &models.Resume{Position: position, Company: company, Salary: salary, Schedule: schedule, DescriptionTasks: description, LinkToVacancy: link}

	return resume, err
}
