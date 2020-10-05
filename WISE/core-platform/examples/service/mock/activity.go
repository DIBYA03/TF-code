package mock

import (
	"time"

	"github.com/wiseco/core-platform/notification/activity"
)

func ActivityList(businessID string) ([]activity.Activity, error) {
	return []activity.Activity{
		{
			ID:       "aqsr23fwr24csfvsfvdsscxvsdvxcxfd-fsdgsbfsdzngf",
			EntityID: businessID,
			Text:     "You account was created",
			Created:  time.Now(),
		},
		{
			ID:       "aqsr23fwr24csfvsfvdsrew3235cfssdgsbfsdzngf",
			EntityID: businessID,
			Text:     "You added John Doe as a new contact",
			Created:  time.Now(),
		},
		{
			ID:       "aqsr23fwr24csfvsfvdsscxfsdfdfs23gsbfsdzngf",
			EntityID: businessID,
			Text:     "Your account ending in 5454, was charged $30.43",
			Created:  time.Now(),
		},
	}, nil
}
