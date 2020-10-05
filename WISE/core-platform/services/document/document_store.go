/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package document

import (
	"errors"
)

//StorerProvider the storer provide
//Currently only s3 is supported
type StorerProvider string

// Storer Read only document storage service
type Storer interface {

	// Storage provider
	Provider() StorerProvider

	// Owner or user id of document
	OwnerID() string

	// Key used to uniquely identify document
	Key() (*string, error)

	// Returns signed url for document
	SignedUrl() (*string, error)

	//Put signed url to upload content directly
	PutSignedURL() (*string, error)

	// Returns content in bytes from storage provider
	Content() ([]byte, error)
}

// Base storage structure
type documentStore struct {
	provider StorerProvider
	product  string
	ownerID  string
}

// Storage providers
const (
	StorageProviderAWSS3 = StorerProvider("awss3")
)

func NewStorerFromKey(key string) (Storer, error) {
	return NewAWSS3DocStoreFromKey(key)
}

func NewStorerFromContent(provider StorerProvider, prefix PrefixType, ownerID, contentType string, content *string) (Storer, error) {
	// Currently S3 only
	switch provider {
	case StorageProviderAWSS3:
		return NewAWSS3DocStorageFromContent(ownerID, prefix, contentType, content)
	default:
		return nil, errors.New("Only AWS S3 is supported at this time")
	}
}

func (s StorerProvider) String() string {
	return string(s)
}
