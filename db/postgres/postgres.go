package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func Connect(dsn string) (*pgx.Conn, error) {
	log.Println("Connecting to postgres database...")
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	log.Println("Successfuly connected to postgres!")
	return conn, nil
}

func MustConnectWithRetry(dsn string, timeout time.Duration) *pgx.Conn {
	timer := time.NewTimer(timeout)
	delay := time.Second
	for {
		select {
		case <-timer.C:
			msg := fmt.Sprintf("cant connect to database with timeout of: %v", timeout)
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
