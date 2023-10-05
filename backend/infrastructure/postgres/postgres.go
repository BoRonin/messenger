package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDB(dsn string) *Postgres {
	pg := &Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
	var counts int64
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgress not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			SetupDbTable(connection)
			pg.Pool = connection
			return pg
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Waiting for 2 seconds...")
		time.Sleep(time.Second * 2)
		continue
	}
}
func SetupDbTable(conn *pgxpool.Pool) {
	game := `SELECT EXISTS
    (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'clients');`
	var result bool
	_ = conn.QueryRow(context.Background(), game).Scan(&result)

	if !result {
		c, ioErr := os.ReadFile("infrastructure/postgres/getUp.sql")
		if ioErr != nil {
			log.Println(ioErr)
		}

		sql := string(c)
		_, err := conn.Exec(context.Background(), sql)
		if err != nil {
			fmt.Printf("Couldn't run query:%s", err.Error())
			return
		}
	}

}
