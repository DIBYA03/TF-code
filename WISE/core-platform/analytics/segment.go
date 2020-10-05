package analytics

import (
	"os"

	"github.com/wiseco/core-platform/shared"
	"gopkg.in/segmentio/analytics-go.v3"
)

func Identify(userID shared.UserID, traits map[string]interface{}) error {
	client := analytics.New(os.Getenv("SEGMENT_WRITE_KEY"))
	defer client.Close()

	t := analytics.NewTraits()
	for k, v := range traits {
		t.Set(k, v)
	}

	err := client.Enqueue(analytics.Identify{
		UserId: userID.ToPrefixString(),
		Traits: t,
	})

	return err
}
