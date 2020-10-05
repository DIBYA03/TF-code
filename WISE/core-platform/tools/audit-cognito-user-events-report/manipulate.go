package main

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

func convertEvent(e *cognitoidentityprovider.AuthEventType) authEvent {
	var riskLevel string
	if e.EventRisk.RiskLevel == nil {
		riskLevel = "null"
	} else {
		riskLevel = *e.EventRisk.RiskLevel
	}

	event := authEvent{
		EventType:    *e.EventType,
		CreationDate: *e.CreationDate,
		RiskLevel:    riskLevel,
		RiskDecision: *e.EventRisk.RiskDecision,
		IPAddress:    *e.EventContextData.IpAddress,
		City:         *e.EventContextData.City,
		Country:      *e.EventContextData.Country,
		DeviceName:   *e.EventContextData.DeviceName,
	}

	return event
}

func processEvents(cAuthEvents []*cognitoidentityprovider.AuthEventType) []authEvent {
	var authEvents []authEvent
	for _, e := range cAuthEvents {
		authEvents = append(authEvents, convertEvent(e))
	}

	return authEvents
}
