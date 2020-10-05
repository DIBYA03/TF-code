package document

import (
	"log"
	"strings"

	coreDB "github.com/wiseco/core-platform/services/data"
)

//HandleS3Event will handle event from s3 trigger
func HandleS3Event(url string) error {
	parts := strings.Split(url, "/")
	return handleURLParts(parts)
}

func handleURLParts(parts []string) error {
	log.Printf("URL parts %v", parts)

	if len(parts) > 0 {
		// check if first part of url is `user`
		if parts[0] == "user" {
			// assuming the user id is before the file id path
			userID := parts[len(parts)-2]
			fileKey := parts[len(parts)-1]
			return userDocuments(userID, fileKey)
		}
		// check if first part of url is `consumer`
		if parts[0] == "consumer" {
			// assuming the consumer id is before the file id path
			consumerID := parts[len(parts)-2]
			fileKey := parts[len(parts)-1]
			return consumerDocuments(consumerID, fileKey)
		}
		// assuming the business id is before the file id path
		businessID := parts[len(parts)-2]
		fileKey := parts[len(parts)-1]
		return businessDocuments(businessID, fileKey)
	}
	return nil
}

func businessDocuments(id, fileKey string) error {
	var storeKeys []struct {
		ID         string `db:"id"`
		StorageKey string `db:"storage_key"`
	}
	err := coreDB.DBRead.Select(&storeKeys, "SELECT id,storage_key FROM business_document WHERE business_id = $1 AND deleted IS NULL", id)
	if err != nil {
		log.Printf("Error getting documents storage keys with business id: %s %v", id, err)
		return err
	}
	for _, k := range storeKeys {
		if fileKey == idFromKey(k.StorageKey) {
			updateBusinessDocument(k.ID)
		}
	}
	return nil
}

func consumerDocuments(id, fileKey string) error {
	var storeKeys []struct {
		ID         string `db:"id"`
		StorageKey string `db:"storage_key"`
	}
	err := coreDB.DBRead.Select(&storeKeys, "SELECT id,storage_key FROM consumer_document WHERE consumer_id = $1 AND deleted IS NULL", id)
	if err != nil {
		log.Printf("Error getting documents storage keys with consumer id: %s %v", id, err)
		return err
	}
	for _, k := range storeKeys {
		if fileKey == idFromKey(k.StorageKey) {
			updateConsumerDocument(k.ID)
		}
	}
	return nil
}

func userDocuments(id, fileKey string) error {
	var storeKeys []struct {
		ID         string `db:"id"`
		StorageKey string `db:"storage_key"`
	}
	err := coreDB.DBRead.Select(&storeKeys, "SELECT id,storage_key FROM user_document WHERE user_id = $1 AND deleted IS NULL", id)
	if err != nil {
		log.Printf("Error getting documents storage keys with user id: %s %v", id, err)
		return err
	}
	for _, k := range storeKeys {
		if fileKey == idFromKey(k.StorageKey) {
			updateUserDocument(k.ID)
		}
	}
	return nil
}

func idFromKey(k string) string {
	parts := strings.Split(k, ":")

	if len(parts) != 5 {
		log.Printf("invalid document storage key")
		return ""
	}
	return parts[4]
}

func updateConsumerDocument(id string) error {
	log.Printf("Updating consumer document with id %s", id)
	_, err := coreDB.DBWrite.Exec("UPDATE consumer_document SET content_uploaded = CURRENT_TIMESTAMP WHERE id = $1", id)
	return err
}

func updateUserDocument(id string) error {
	log.Printf("Updating user document with id %s", id)
	_, err := coreDB.DBWrite.Exec("UPDATE user_document SET content_uploaded = CURRENT_TIMESTAMP WHERE id = $1", id)
	return err
}

func updateBusinessDocument(id string) error {
	log.Printf("Updating business document with id %s", id)
	_, err := coreDB.DBWrite.Exec("UPDATE business_document SET content_uploaded = CURRENT_TIMESTAMP WHERE id = $1", id)
	return err
}
