package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"parser/scraper"
)

//храним заказ в базе, если такой уже есть (exists) то пропускаем, если нет, добавляем в базу и
//отправляем в телеграм

//поиск новых заказов должен работать непрерывно, раз в 20 минут проверяем новый заказ

type Postgres struct {
	repo *pgx.Conn
}

// получаем первый список => добавляем в базу => генерируем снова список => проверяем по номеру href => если такой
func (p *Postgres) AddOrder(href int, title, url, finance, description string) {
	query := `insert into habr (href, title, url, finance, description) values ($1,$2,$3,$4,$5)`
	_, err := p.repo.Exec(context.Background(), query, &href, &title, &url, &finance, &description)
	if err != nil {
		//запись уже существует? пропускаем добавляемый href
		log.Fatal("Not inserted values in DB")
		return
	}
}
func (p *Postgres) AddHref(href int) {
	query := `insert into habr (href) values ($1)`
	_, err := p.repo.Exec(context.Background(), query, &href)
	if err != nil {
		//запись уже существует? пропускаем добавляемый href
		log.Fatal("Not inserted values in DB")
		return
	}
}

func (p *Postgres) GetOrderByHref(href string) ([]scraper.Habr, error) {
	qry := `select from habr (href, title, url, finance, description) where href = $1`
	//query := `insert into habr (href, title, url, finance, description) values ($1,$2,$3,$4,$5)`
	info := scraper.Habr{}
	var infoArray []scraper.Habr
	qr, err := p.repo.Query(context.Background(), qry, &href)
	if err != nil {
		//запись уже существует? пропускаем добавляемый href
		log.Fatal("Not inserted values in DB")
		return nil, err
	}
	for qr.Next() {
		err := qr.Scan(&info.Href, &info.Title, &info.Url, &info.Finance, &info.Description)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		infoArray = append(infoArray, info)
	}
	return infoArray, err
}
