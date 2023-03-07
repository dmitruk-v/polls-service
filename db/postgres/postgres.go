package postgres

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const SeedSQL = `
CREATE TABLE IF NOT EXISTS polls (
	poll_id BIGSERIAL PRIMARY KEY NOT NULL,
	survey_id INT UNIQUE NOT NULL,
	pre_set_values JSON NOT NULL
);
`

func Connect(dsn string) (*pgx.Conn, error) {
	log.Println("Connecting to POSTGRES database...")
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	log.Println("Successfuly connected to POSTGRES!")
	return conn, nil
}

func MustConnectWithRetry(dsn string, timeout time.Duration) *pgx.Conn {
	timer := time.NewTimer(timeout)
	delay := time.Second
	for {
		select {
		case <-timer.C:
			msg := fmt.Sprintf("cant connect to POSTGRES with timeout of: %v", timeout)
			panic(msg)
		default:
			conn, err := Connect(dsn)
			if err == nil {
				return conn
			}
			log.Println(err)
		}
		log.Printf("Retry after %v...\n", delay)
		time.Sleep(delay)
		delay *= 2
	}
}

func MustSeedString(db *pgx.Conn, src string) {
	if err := Seed(db, strings.NewReader(src)); err != nil {
		panic(err)
	}
}

func Seed(db *pgx.Conn, r io.Reader) error {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, string(data)); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
