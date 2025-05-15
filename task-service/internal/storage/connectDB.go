package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"task-service/config"

	_ "github.com/lib/pq"
)

func ConnectToDB(cfg config.Config) (*sql.DB, error) {
	psqlString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	var connDb *sql.DB
	var err error

	for i := 0; i < 5; i++ {
		connDb, err = sql.Open("postgres", psqlString)
		if err == nil {
			err = connDb.Ping()
			if err == nil {
				return connDb, nil
			}
		}

		log.Printf("Error connecting to database, retrying... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("Unable to connect to database: %v", err)
}