/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package document

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

// AWS S3 specific storage implementation
type awsS3DocStore struct {
	documentStore
	bucket string
	prefix PrefixType
	fileID string
}

type PrefixType string

type Options struct {
	bucket      string
	Prefix      PrefixType
	fileID      string
	OwnerID     string
	contentType string
}

const (
	BusinessPrefix = PrefixType("business")
	UserPrefix     = PrefixType("user")
	ConsumerPrefix = PrefixType("consumer")
)

func (PrefixType) new(v string) PrefixType {
	return PrefixType(v)
}

func (p PrefixType) String() string {
	return string(p)
}

// NewAWSS3DocStoreFromKey Create storage S3 storage from encoded key
func NewAWSS3DocStoreFromKey(k string) (*awsS3DocStore, error) {
	parts := strings.Split(k, ":")

	if len(parts) != 5 || parts[0] != string(StorageProviderAWSS3) {
		return nil, errors.New("Invalid aws document storage key")
	}

	docStore := awsS3DocStore{
		documentStore: documentStore{
			provider: StorageProviderAWSS3,
			ownerID:  parts[1],
		},
		bucket: parts[2],
		prefix: PrefixType(parts[3]),
		fileID: parts[4],
	}

	return &docStore, nil
}

// NewAWSS3DocStorageFromContent Store new document and return storage object
func NewAWSS3DocStorageFromContent(ownerID string, prefix PrefixType, contentType string, content *string) (*awsS3DocStore, error) {
	if content == nil {
		return nil, errors.New("Document content missing")
	}

	c, err := base64.StdEncoding.DecodeString(*content)
	if err != nil {
		return nil, err
	}

	// S3 session uploader
	sess := session.Must(session.NewSession())

	// AWS multi-part capable uploader
	up := s3manager.NewUploader(sess)

	// S3 document bucket
	s3BucketName := os.Getenv("AWS_S3_BUCKET_DOCUMENT")
	if len(s3BucketName) == 0 {
		return nil, errors.New("Missing env variable for `AWS_S3_BUCKET_DOCUMENT`")
	}

	// Generate file id
	fileID := uuid.New().String()

	k := prefix.String() + "/" + string(ownerID) + "/" + fileID

	// Upload
	_, err = up.Upload(
		&s3manager.UploadInput{
			Bucket:      aws.String(s3BucketName),
			Key:         aws.String(k),
			ContentType: aws.String(contentType),
			Body:        bytes.NewReader(c),
		},
	)

	if err != nil {
		return nil, err
	}

	docStore := awsS3DocStore{
		documentStore: documentStore{
			provider: StorageProviderAWSS3,
			ownerID:  ownerID,
		},
		bucket: s3BucketName,
		prefix: prefix,
		fileID: fileID,
	}

	return &docStore, nil
}

//NewAWSS3DocStorage prepares a storer for put signed url
func NewAWSS3DocStorage(ownerID string, prefix PrefixType) (*awsS3DocStore, error) {
	// S3 document bucket
	s3BucketName := os.Getenv("AWS_S3_BUCKET_DOCUMENT")
	if len(s3BucketName) == 0 {
		log.Printf("Missing env variable for `AWS_S3_BUCKET_DOCUMENT`")
		return nil, errors.New("Missing env variable for `AWS_S3_BUCKET_DOCUMENT`")
	}

	// Generate file id
	fileID := uuid.New().String()

	docStore := awsS3DocStore{
		documentStore: documentStore{
			provider: StorageProviderAWSS3,
			ownerID:  ownerID,
		},
		bucket: s3BucketName,
		prefix: prefix,
		fileID: fileID,
	}

	return &docStore, nil
}

func (docStore *awsS3DocStore) UploadBuffer(rs io.ReadSeeker, options Options) (string, error) {
	// S3 session
	sess := session.Must(session.NewSession())

	// S3 Service
	srv := s3.New(sess)

	// ensure we have a bucket and a fileID
	// if we created a docStore from key we should have a bucket already
	if options.bucket == "" {
		options.bucket = docStore.bucket
	}
	if options.fileID == "" {
		options.fileID = uuid.New().String()
	}

	req, _ := srv.PutObjectRequest(
		&s3.PutObjectInput{
			Bucket:      aws.String(docStore.bucket),
			Key:         aws.String(docStore.makeKey(&options)),
			Body:        rs,
			ContentType: aws.String(options.contentType),
		},
	)
	err := req.Send()
	if err != nil {
		log.Printf("error uploading content %v", err)
	}
	key, err := docStore.Key()
	if err != nil || key == nil {
		return "", errors.New("storage key is empty")
	}
	return *key, err
}

// AWS S3 DocumentStorageService implementation
func (docStore *awsS3DocStore) Provider() StorerProvider {
	return docStore.provider
}

func (docStore *awsS3DocStore) OwnerID() string {
	return docStore.ownerID
}

func (docStore *awsS3DocStore) Key() (*string, error) {
	k := strings.Join(
		[]string{
			docStore.provider.String(),
			string(docStore.ownerID),
			docStore.bucket,
			docStore.prefix.String(),
			docStore.fileID,
		},
		":",
	)

	return &k, nil
}

func (docStore *awsS3DocStore) SignedUrl() (*string, error) {
	// S3 session
	sess := session.Must(session.NewSession())

	// S3 Service
	srv := s3.New(sess)

	// Generate file id
	k := docStore.prefix.String() + "/" + string(docStore.ownerID) + "/" + docStore.fileID

	awsRequest, _ := srv.GetObjectRequest(
		&s3.GetObjectInput{
			Bucket: aws.String(docStore.bucket),
			Key:    aws.String(k),
		},
	)

	url, err := awsRequest.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Sign key:", url, err)
		return nil, err
	}

	return &url, nil
}

func (docStore *awsS3DocStore) PutSignedURL() (*string, error) {
	// S3 session
	sess := session.Must(session.NewSession())

	// S3 Service
	srv := s3.New(sess)

	// Generate file id
	k := docStore.prefix.String() + "/" + string(docStore.ownerID) + "/" + docStore.fileID

	req, _ := srv.PutObjectRequest(
		&s3.PutObjectInput{
			Bucket: aws.String(docStore.bucket),
			Key:    aws.String(k),
		},
	)

	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Sign key:", url, err)
		return nil, err
	}

	return &url, nil
}

// when passing options ensure you are passing all the fields
// as the key gets constructed based on those fields
func (docStore *awsS3DocStore) makeKey(options *Options) string {
	if options != nil {
		return options.Prefix.String() + "/" + options.OwnerID + "/" + options.fileID
	}
	return docStore.prefix.String() + "/" + string(docStore.ownerID) + "/" + docStore.fileID
}

func (docStore *awsS3DocStore) Content() ([]byte, error) {
	// S3 session
	sess := session.Must(session.NewSession())

	// S3 Downloader
	awsDownloader := s3manager.NewDownloader(sess)

	// Generate file id
	k := docStore.prefix.String() + "/" + string(docStore.ownerID) + "/" + docStore.fileID

	awsWriteBuffer := &aws.WriteAtBuffer{}

	_, err := awsDownloader.Download(
		awsWriteBuffer,
		&s3.GetObjectInput{
			Bucket: aws.String(docStore.bucket),
			Key:    aws.String(k),
		},
	)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println(docStore.bucket, k, aerr.Code())
		}

		return []byte{}, err
	}

	return awsWriteBuffer.Bytes(), nil
}

// return a readSeeker instead of []byte
func (docStore *awsS3DocStore) ContentBuffer() (io.ReadSeeker, error) {
	sess := session.Must(session.NewSession())
	awsDownloader := s3manager.NewDownloader(sess)

	awsWriteBuffer := &aws.WriteAtBuffer{}
	k := docStore.makeKey(nil)

	_, err := awsDownloader.Download(
		awsWriteBuffer,
		&s3.GetObjectInput{
			Bucket: aws.String(docStore.bucket),
			Key:    aws.String(k),
		},
	)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println(docStore.bucket, k, aerr.Code())
		}
		return nil, err
	}

	return bytes.NewReader(awsWriteBuffer.Bytes()), nil
}

// DownloadBBVAReQueueContent Uses s3 manager to download content from a given URL
func DownloadBBVAReQueueContent() ([]byte, error) {
	// S3 session
	sess := session.Must(session.NewSession())

	// S3 Downloader
	awsDownloader := s3manager.NewDownloader(sess)

	// S3 document bucket
	s3BucketName := os.Getenv("AWS_S3_BUCKET_DOCUMENT")
	if s3BucketName == "" {
		return nil, errors.New("Missing env variable for `AWS_S3_BUCKET_DOCUMENT`")
	}
	// S3 document bucket
	bbvarequeues3 := os.Getenv("BBVA_REQUEUE_S3_OBJECT")
	if s3BucketName == "" {
		return nil, errors.New("Missing env variable for `BBVA_REQUEUE_S3_OBJECT`")
	}

	awsWriteBuffer := &aws.WriteAtBuffer{}

	_, err := awsDownloader.Download(
		awsWriteBuffer,
		&s3.GetObjectInput{
			Bucket: aws.String(s3BucketName),
			Key:    aws.String(bbvarequeues3),
		},
	)
	if err != nil {
		return []byte{}, err
	}

	return awsWriteBuffer.Bytes(), nil
}
