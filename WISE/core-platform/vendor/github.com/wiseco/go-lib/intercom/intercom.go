/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package intercom

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/wiseco/go-lib/api"
	"github.com/wiseco/go-lib/log"
	"github.com/wiseco/go-lib/url"
)

type intercomService struct {
	log log.Logger
}

type IntercomService interface {
	GetUserByEmail(string) (*string, error)
	ListUserBySegmentID(string) (*[]User, error)
	ListConversationByEmailID(string) (*[]Conversation, error)
	GetUserByID(userID string) (*string, error)

	GetTags() (*TagList, error)
	SetTags(tagActions []TagAction) (*TagList, error)
}

func NewIntercomService(l log.Logger) IntercomService {
	return &intercomService{log: l}
}

type User struct {
	Type                   string           `json:"type"`
	ID                     string           `json:"id"`
	UserID                 string           `json:"user_id"`
	Anonymous              bool             `json:"anonymous"`
	Email                  string           `json:"email"`
	Phone                  string           `json:"phone"`
	Name                   string           `json:"name"`
	Pseudonym              *string          `json:"pseudonym"`
	Avatar                 Avatar           `json:"avatar"`
	AppID                  string           `json:"app_id"`
	LocationData           Location         `json:"location_data"`
	LastRequestAt          int64            `json:"last_request_at"`
	CreatedAt              int64            `json:"created_at"`
	RemoteCreatedAt        int64            `json:"remote_created_at"`
	SignedUpAt             int64            `json:"signed_up_at"`
	UpdatedAt              int64            `json:"updated_at"`
	SessionCount           int              `json:"session_count"`
	OwnerID                *string          `json:"owner_id"`
	UnsubscribedFromEmails bool             `json:"unsubscribed_from_emails"`
	MarkedEmailAsSpam      bool             `json:"marked_email_as_spam"`
	HasHardBounced         bool             `json:"has_hard_bounced"`
	Tag                    TagList          `json:"tags"`
	Segment                SegmentList      `json:"segments"`
	CustomAttributes       CustomAttributes `json:"custom_attributes"`
	Referrer               string           `json:"referrer"`
	UTMCampaign            string           `json:"utm_campaign"`
	UTMContent             string           `json:"utm_content"`
	UTMMedium              string           `json:"utm_medium"`
	UTMSource              string           `json:"utm_source"`
	UTMTerm                string           `json:"utm_term"`
	DoNotTract             *string          `json:"do_not_track"`
}

type CustomAttributes struct {
	ArticleID                string  `json:"article_id"`
	ID                       string  `json:"id"`
	JobTitle                 string  `json:"job_title"`
	BusinessID               string  `json:"business_id"`
	BusinessKYCStatus        *string `json:"business_kyc_status"`
	BusinessLegalName        *string `json:"business_legal_name"`
	EntityType               *string `json:"entity_type"`
	IndustryType             *string `json:"industry_type"`
	BusinessCardID           string  `json:"business_card_id"`
	BusinessCardStatus       *string `json:"business_card_status"`
	BusinessTransactionCount int     `json:"business_transaction_count"`
	DateOfBirthAt            int     `json:"date_of_birth_at"`
	FirstName                string  `json:"first_name"`
	LastName                 string  `json:"last_name"`
	PhoneVerified            bool    `json:"phone_verified"`
	MiddleName               *string `json:"middle_name"`
	BusinessAccountID        *string `json:"business_account_id"`
	BusinessAccountOpenedAt  int64   `json:"business_account_opened_at"`
	BusinessAccountStatus    *string `json:"business_account_status"`
	BusinessAccountType      *string `json:"business_account_type"`
	BusinessBankName         *string `json:"business_bank_name"`
	SID                      string  `json:"sid"`
}

type SegmentList struct {
	Type    string    `json:"type"`
	Segment []Segment `json:"segments"`
}

type Segment struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type TagList struct {
	Type string `json:"type"`
	Tag  []Tag  `json:"tags"`
}

type Location struct {
	Type          string  `json:"type"`
	CityName      *string `json:"city_name"`
	ContinentCode *string `json:"continent_code"`
	CountryName   *string `json:"country_name"`
	CountryCode   *string `json:"country_code"`
	PostalCode    *string `json:"postal_code"`
	RegionName    *string `json:"region_name"`
	TimeZone      *string `json:"timezone"`
}

type Avatar struct {
	Type     string  `json:"type"`
	ImageURL *string `json:"image_url"`
}

type Page struct {
	Type       string `json:"type"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalPages int    `json:"total_pages"`
}

type UserResponse struct {
	Pages      Page   `json:"pages"`
	TotalCount int    `json:"total_count"`
	Limited    bool   `json:"limited"`
	Type       string `json:"type"`
	Users      []User `json:"users"`
}

type ConversationResponse struct {
	Type          string         `json:"type"`
	Pages         Page           `json:"pages"`
	Conversations []Conversation `json:"conversations"`
}

type Conversation struct {
	Tag TagList `json:"tags"`
}

func (s *intercomService) GetUserByEmail(emailID string) (*string, error) {
	// Intercom Access Token
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	if len(accessToken) == 0 {
		return nil, errors.New("INTERCOM_ACCESS_TOKEN is missing")
	}

	// URL
	u, err := url.BuildAbsoluteForIntercom("", url.Params{})
	if err != nil {
		return nil, err
	}

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:     u,
			AcceptType:  "application/json",
			BearerToken: &api.BearerTokenConfig{AccessToken: accessToken},
		},
		s.log)

	// Send request
	ur := &UserResponse{}
	params := map[string]string{"email": emailID}
	err = client.Get("users", params, ur)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(ur)
	if err != nil {
		return nil, err
	}

	// Return response
	res := string(body)
	return &res, nil
}

func (s *intercomService) GetUserByID(userID string) (*string, error) {
	// Intercom Access Token
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	if len(accessToken) == 0 {
		return nil, errors.New("INTERCOM_ACCESS_TOKEN is missing")
	}

	// URL
	u, err := url.BuildAbsoluteForIntercom("", url.Params{})
	if err != nil {
		return nil, err
	}

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:     u,
			AcceptType:  "application/json",
			BearerToken: &api.BearerTokenConfig{AccessToken: accessToken},
		},
		s.log)

	// Send request
	var ur User
	params := map[string]string{"user_id": userID}
	err = client.Get("users", params, &ur)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(ur)
	if err != nil {
		return nil, err
	}

	// Return response
	res := string(body)
	return &res, nil
}

func (s *intercomService) ListUserBySegmentID(segmentID string) (*[]User, error) {

	var users []User
	perPage := 50
	page := 1
	totalPage := 1

	for page <= totalPage {
		ur, err := s.ListUserByPage(segmentID, page, perPage)
		if err != nil {
			return nil, err
		}

		users = append(users, ur.Users...)

		perPage = ur.Pages.PerPage
		totalPage = ur.Pages.TotalPages
		page = ur.Pages.Page + 1

	}

	return &users, nil
}

func (s *intercomService) ListUserByPage(segmentID string, page int, perPage int) (*UserResponse, error) {
	// Intercom Access Token
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	if len(accessToken) == 0 {
		return nil, errors.New("INTERCOM_ACCESS_TOKEN is missing")
	}

	// URL
	u, err := url.BuildAbsoluteForIntercom("", url.Params{})
	if err != nil {
		return nil, err
	}

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:     u,
			AcceptType:  "application/json",
			BearerToken: &api.BearerTokenConfig{AccessToken: accessToken},
		},
		s.log)

	// Send request
	ur := &UserResponse{}

	params := map[string]string{"segment_id": segmentID, "per_page": strconv.Itoa(perPage), "page": strconv.Itoa(page)}

	err = client.Get("users", params, ur)
	if err != nil {
		return nil, err
	}

	return ur, nil
}

func (s *intercomService) ListConversationByEmailID(emailID string) (*[]Conversation, error) {

	var c []Conversation
	perPage := 50
	page := 1
	totalPage := 1

	for page <= totalPage {
		cr, err := s.ListConversationByPage(emailID, page, perPage)
		if err != nil {
			return nil, err
		}

		c = append(c, cr.Conversations...)

		perPage = cr.Pages.PerPage
		totalPage = cr.Pages.TotalPages
		page = cr.Pages.Page + 1

	}

	return &c, nil
}

func (s *intercomService) ListConversationByPage(emailID string, page int, perPage int) (*ConversationResponse, error) {
	// Intercom Access Token
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	if len(accessToken) == 0 {
		return nil, errors.New("INTERCOM_ACCESS_TOKEN is missing")
	}

	// URL
	u, err := url.BuildAbsoluteForIntercom("", url.Params{})
	if err != nil {
		return nil, err
	}

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:     u,
			AcceptType:  "application/json",
			BearerToken: &api.BearerTokenConfig{AccessToken: accessToken},
		},
		s.log)

	// Send request
	c := &ConversationResponse{}

	params := map[string]string{"email": emailID, "type": "user", "per_page": strconv.Itoa(perPage), "page": strconv.Itoa(page)}

	err = client.Get("conversations", params, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
