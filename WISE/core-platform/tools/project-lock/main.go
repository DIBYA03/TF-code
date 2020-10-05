package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/wiseco/core-platform/services"
	csp "github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/services/csp/consumer"
	"github.com/wiseco/core-platform/shared"
)

const (
	envDev     = "dev"
	envStg     = "stg"
	envStaging = "staging"
	envQA      = "qa"
	envSbx     = "sbx"
	envPrd     = "prd"
	envProd    = "prod"
)

var failureCount int
var failureTime time.Duration

var successCount int
var successTime time.Duration

var skippedCount int

func main() {
	// consumerIds := getConsumerIDsFromFile()
	fmt.Println("Run project-lock")
	consumerIds := getActiveConsumerIDs()

	failureCount = 0
	successCount = 0
	skippedCount = 0

	offset := 0
	limit := 10
	startTime := time.Now()

	for {
		start := offset
		end := offset + limit
		if len(consumerIds) < end {
			end = len(consumerIds)
			l := len(consumerIds) - offset
			if l <= 0 {
				break
			}
		}

		conIDs := consumerIds[start:end]
		if len(conIDs) == 0 {
			log.Println("no more notifications:", offset, 10)
			break
		}

		var wg sync.WaitGroup
		wg.Add(len(conIDs))
		for _, cID := range conIDs {
			go func(cID string) {
				defer wg.Done()
				runClearKYC(shared.ConsumerID(cID))

			}(cID)
		}

		wg.Wait()
		offset += len(conIDs)
	}

	endTime := time.Now()

	totalDiff := endTime.Sub(startTime)
	fmt.Printf("\n\nTOTAL DURATION: %.2g\n", totalDiff.Seconds())
	fmt.Printf("Success: %v\n", successCount)
	fmt.Printf("Failure: %v\n", failureCount)
	fmt.Printf("Skipped: %v\n", skippedCount)
}

func runClearKYC(consumerID shared.ConsumerID) {
	startTime := time.Now()
	_, err := consumer.New().GetClearKYC(consumerID)
	if err == nil {
		fmt.Println("Clear KYC already done for Cosnumer: ", consumerID)
		skippedCount++
		return
	}

	_, err = consumer.New().RunClearKYC(consumerID)
	endTime := time.Now()
	diff := endTime.Sub(startTime)

	if err != nil {
		failureCount++
		failureTime += diff
		fmt.Printf("FAILURE %v: %v,\n", failureCount, diff)
		return
	}
	successCount++
	successTime += diff

	fmt.Printf("SUCCESS %v: %v,\n", successCount, diff)
}

func getActiveConsumerIDs() []string {
	service := csp.New(services.NewSourceRequest())
	activeBusinessIDs, err := service.GetActiveBusinessIDsInternal()
	if err != nil {
		log.Fatalf("Failed to get active businessIds: %v", err)
	}
	fmt.Println("Active Business IDs: ", len(activeBusinessIDs))

	memberService := csp.NewMemberService(services.NewSourceRequest())

	consumerIds, err := memberService.GetConsumerIDsForBusinessIDsInternal(activeBusinessIDs)
	if err != nil {
		log.Fatalf("Failed to get active consumerIds")
	}
	fmt.Println("Consumer IDs: ", len(consumerIds))
	return consumerIds
}

func getConsumerIDsFromFile() []string {
	filepath := getFilePath()

	if filepath == "" {
		log.Fatalf("Filepath empty")
	}

	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var consumerIds []string
	err = json.Unmarshal([]byte(dat), &consumerIds)
	if err != nil {
		log.Fatalf("Failed to parse file content: %v", err)
	}
	return consumerIds
}

func getFilePath() string {
	env := os.Getenv("API_ENV")
	filepath := ""

	switch env {
	case envDev:
		filepath = "./consumer_ids_dev.json"
	case envStaging:
		filepath = "./consumer_ids_staging.json"
	case envQA:
		filepath = "./consumer_ids_qa.json"
	case envProd:
		filepath = "./consumer_ids_prod.json"
	default:
		filepath = ""
	}

	return filepath
}
