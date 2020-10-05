package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func csvHeaders() []string {
	return []string{
		"username",
		"phone",
		"is_flagged",
		"flagged_reason",
		"event_type",
		"creation_date",
		"risk_level",
		"risk_decision",
		"ip_address",
		"city",
		"country",
		"device_name",
	}
}

func convertIsFlagged(flag bool) string {
	if flag == true {
		return "yes"
	}

	return "no"
}

func convertFlaggedReasons(reasons []string) string {
	return strings.Join(reasons, "\n")
}

func generateReport(users []cognitoUser) {
	file, err := os.Create("results.csv")
	if err != nil {
		log.Panic("error creating csv file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err = writer.Write(csvHeaders()); err != nil {
		log.Panic("can't right headers to csv:", err)
	}

	for _, user := range users {
		log.Print(fmt.Sprintf("%s has %d events. Checking if flagged", user.Phone, len(user.AuthEvents)))

		authEvents := user.AuthEvents
		for _, e := range authEvents {
			err := writer.Write([]string{
				user.Username,
				user.Phone,
				convertIsFlagged(e.Flagged.IsFlagged),
				convertFlaggedReasons(e.Flagged.Reasons),
				e.EventType,
				e.CreationDate.Format("2006-01-02 15:04:05"),
				e.RiskLevel,
				e.RiskDecision,
				e.IPAddress,
				e.City,
				e.Country,
				e.DeviceName,
			})
			if err != nil {
				log.Panic("can't add event:", err)
			}
		}
	}
}
