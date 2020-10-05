package messages

import (
	"fmt"
	"log"

	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

func getConsumerIDFromUserID(id shared.UserID) (consumerID shared.ConsumerID, err error) {
	err = data.DBRead.Get(&consumerID, fmt.Sprintf(`SELECT consumer.id
	FROM wise_user
	JOIN consumer ON wise_user.consumer_id = consumer.id
	WHERE wise_user.id = '%s'
	`, id))
	if err != nil {
		log.Printf("Error getting user with id:%s  details:%v", id, err)
	}
	return consumerID, err
}
