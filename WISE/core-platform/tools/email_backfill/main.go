package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/email"
)

type userRow struct {
	Email      *string `db:"email"`
	ConsumerID string  `db:"consumer_id"`
}

func main() {
	liveRun := flag.Bool("live-run", false, "If this is omitted this backfill will log and not write")

	flag.Parse()

	if *liveRun {
		fmt.Println("########## Live run ##########")
	}

	db, err := getDBWrite()

	if err != nil {
		panic(err)
	}

	es := email.NewEmailService(services.NewSourceRequest())

	err = backfillConsumer(es, db, *liveRun)

	if err != nil {
		panic(err)
	}
}

func backfillConsumer(es email.EmailService, db *sqlx.DB, liveRun bool) error {
	done := false

	offset := 0
	limit := 20

	fmt.Println("########## Backfilling consumer ##########")

	for done == false {
		rows := []userRow{}

		err := db.Select(&rows, "SELECT email, consumer_id from wise_user where deactivated is null and email is not null LIMIT $1 OFFSET $2", limit, offset)

		if err != nil {
			return err
		}

		if len(rows) == 0 {
			done = true
		} else {
			offset = offset + limit
		}

		for _, row := range rows {
			ec := email.EmailCreate{
				EmailAddress: email.EmailAddress(*row.Email),
				Status:       email.StatusActive,
				Type:         email.TypeConsumer,
			}

			fmt.Printf("%+v\n", ec)

			if liveRun {
				e, err := es.Create(&ec)

				if err != nil {
					return err
				}

				_, err = db.Exec("UPDATE consumer set email_id = $1 where id = $2", e.ID, row.ConsumerID)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func getDBWrite() (*sqlx.DB, error) {
	user := os.Getenv("CORE_DB_USER")
	password := os.Getenv("CORE_DB_PASSWD")
	dbName := os.Getenv("CORE_DB_NAME")

	// Set up write connection
	writeEndpoint := os.Getenv("CORE_DB_WRITE_URL")
	writePort := os.Getenv("CORE_DB_WRITE_PORT")
	ssl := os.Getenv("CORE_DB_SSL_MODE")

	writeConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		writeEndpoint, writePort, user, password, dbName, ssl,
	)

	DBWrite, err := sqlx.Open("postgres", writeConnection)

	return DBWrite, err
}
