package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

func PostgresInit(urlDb string) *pgx.Conn {

	conn, err := pgx.Connect(context.Background(), urlDb)
	if err != nil {
		log.Prefix()
	}

	defer conn.Close(context.Background())

	return conn
}

func Connect(host, port, dbname, user, password string) string {
	psqlconn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	return psqlconn
}

func OpenConn() (*sql.DB, error) {
	db, err := sql.Open("postgres", Connect("localhost", "5432", "postgres", "postgres", "postgres"))
	return db, err
}
