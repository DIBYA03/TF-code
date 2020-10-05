package support

import (
	"log"
	"testing"
)

//Testing with local db where the created user has this id associated
func TestUserByIDSuccess(t *testing.T) {
	user, err := New().UserByID("aa41fcc7-acbb-480d-b285-dded8bcdf645")
	if err != nil {
		t.Errorf("Fail to get the user by id err: %v", err)
	}
	log.Printf("Response from userByID %v", user)
}

//Testing with local db where the created user has this phone number associated
func TestUserByPhone(t *testing.T) {

	user, err := New().UserByPhone("6505220481")
	if err != nil {
		t.Errorf("Fail to get the user by id err: %v", err)
	}

	log.Printf("Response from userByID %v", user)
}
