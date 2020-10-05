package data

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DBWrite *sqlx.DB
var DBRead *sqlx.DB

func init() {
	isLambda := strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda_")
	log.Println("Lambda Mode: ", isLambda)

	var user = os.Getenv("BBVA_DB_USER")
	var password = os.Getenv("BBVA_DB_PASSWD")
	var dbName = os.Getenv("BBVA_DB_NAME")

	// Set up write connection
	var writeEndpoint = os.Getenv("BBVA_DB_WRITE_URL")
	var writePort = os.Getenv("BBVA_DB_WRITE_PORT")
	var ssl = os.Getenv("BBVA_DB_SSL_MODE")
	if ssl == "" {
		ssl = "require"
	}

	var writeConnection = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		writeEndpoint, writePort, user, password, dbName, ssl,
	)

	var err error
	DBWrite, err = sqlx.Open("postgres", writeConnection)
	if err != nil {
		log.Println(dbName, err)
		log.Panic(err)
	}

	DBWrite = DBWrite.Unsafe()

	DBWrite.SetMaxOpenConns(1)
	DBWrite.SetMaxIdleConns(1)

	err = DBWrite.Ping()
	if err != nil {
		log.Println(dbName, err)
		log.Panic(err)
	}

	var readEndpoint = os.Getenv("BBVA_DB_READ_URL")
	var readPort = os.Getenv("BBVA_DB_READ_PORT")

	var readConnection = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		readEndpoint, readPort, user, password, dbName, ssl,
	)

	DBRead, err = sqlx.Open("postgres", readConnection)
	if err != nil {
		log.Println(dbName, err)
		log.Panic(err)
	}

	DBRead = DBRead.Unsafe()

	DBRead.SetMaxOpenConns(1)
	DBRead.SetMaxIdleConns(1)

	err = DBRead.Ping()
	if err != nil {
		log.Println(dbName, err)
		log.Panic(err)
	}
}
