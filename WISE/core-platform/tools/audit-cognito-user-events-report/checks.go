package main

func validateIPAddress(username string, event authEvent, users []cognitoUser) authEvent {
	for _, u := range users {
		for _, e := range u.AuthEvents {
			if event.IPAddress == e.IPAddress && u.Username != username {
				event.Flagged.IsFlagged = true
				event.Flagged.Reasons = append(event.Flagged.Reasons, crossAccountIPAddress)

				// Need to return or otherwise it can add multiple lines of flagged for
				// same thing
				return event
			}
		}
	}

	return event
}

func validateCountry(event authEvent) authEvent {
	if event.Country != "United States" {
		event.Flagged.IsFlagged = true
		event.Flagged.Reasons = append(event.Flagged.Reasons, invalidCountry)
	}

	return event
}

func validateRiskLevel(event authEvent) authEvent {
	if event.RiskLevel == "High" {
		event.Flagged.IsFlagged = true
		event.Flagged.Reasons = append(event.Flagged.Reasons, highRiskLevel)
	}

	return event
}

func validateEvents(users []cognitoUser) []cognitoUser {
	var flaggedUsers []cognitoUser
	for _, u := range users {
		userFlagged := false
		var flaggedEvents []authEvent
		for _, e := range u.AuthEvents {
			// Add new changes for events here
			e = validateIPAddress(u.Username, e, users)
			e = validateCountry(e)
			e = validateRiskLevel(e)

			if e.Flagged.IsFlagged {
				flaggedEvents = append(flaggedEvents, e)
				userFlagged = true
			}
		}

		if userFlagged {
			u.AuthEvents = flaggedEvents
			flaggedUsers = append(flaggedUsers, u)
		}
	}

	return flaggedUsers
}
