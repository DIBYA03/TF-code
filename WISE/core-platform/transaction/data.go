package transaction

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/wiseco/core-platform/shared"
)

var DBWrite *sqlx.DB
var DBRead *sqlx.DB

var LogProvider shared.StreamProvider

func init() {
	isLambda := strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda_")
	log.Println("Lambda Mode: ", isLambda)

	// Setup transacrtion DB connections
	var user = os.Getenv("TXN_DB_USER")
	var password = os.Getenv("TXN_DB_PASSWD")
	var dbName = os.Getenv("TXN_DB_NAME")

	// Set up write connection
	var writeEndpoint = os.Getenv("TXN_DB_WRITE_URL")
	var writePort = os.Getenv("TXN_DB_WRITE_PORT")
	var ssl = os.Getenv("TXN_DB_SSL_MODE")
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

	// Set up read connection
	var readEndpoint = os.Getenv("TXN_DB_READ_URL")
	var readPort = os.Getenv("TXN_DB_READ_PORT")

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

	// Set up kinesis
	streamRegion := os.Getenv("KINESIS_TRX_REGION")
	if streamRegion == "" {
		panic("Notification stream region missing")
	}

	streamName := os.Getenv("KINESIS_TRX_NAME")
	if streamName == "" {
		panic("Notification stream name missing")
	}

	LogProvider = shared.NewKinesisStreamProvider(shared.StreamProviderRegion(streamRegion), streamName)
}
