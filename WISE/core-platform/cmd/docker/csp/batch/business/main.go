package main

import (
	"log"
	"time"
)

func main() {
	start()
}

func start() error {
	log.Printf("Batch check for business kcy status started at %v", time.Now().Format(time.ANSIC))
	err := fetchBusinesses()
	if err != nil {
		log.Printf("error %v", err)
	}
	log.Printf("Batch check for business kcy status ended at %v", time.Now().Format(time.ANSIC))
	return nil
}
