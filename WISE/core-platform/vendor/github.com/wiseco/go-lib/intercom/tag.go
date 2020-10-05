/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package intercom

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/wiseco/go-lib/api"
	"github.com/wiseco/go-lib/log"
	"github.com/wiseco/go-lib/url"
)

//Tag for CSP API
type Tag struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	AppliedAt int64  `json:"applied_at"`
}

//TagAction for CSP API
type TagAction struct {
	Name   string        `json:"name"`
	Action TagActionName `json:"action"`
	UserID string        `json:"userId"`
}

//TagActionName ...
type TagActionName string

const (
	//SetTagAction ...
	SetTagAction = TagActionName("tag")

	//UnsetTagAction ...
	UnsetTagAction = TagActionName("untag")
)

//TagSet for intercom API
type TagSet struct {
	UserID string `json:"user_id"`
	Untag  bool   `json:"untag"`
}

//TagSetBody for intercom API
type TagSetBody struct {
	Name  string   `json:"name"`
	Users []TagSet `json:"users"`
}

func (s *intercomService) GetTags() (*TagList, error) {
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
	tl := &TagList{}
	err = client.Get("tags", nil, tl)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

func (s *intercomService) SetTags(tagActions []TagAction) (*TagList, error) {

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

	var tl = TagList{Type: "tag.list"}

	for _, t := range tagActions {
		nt, tagErr := postTag(t, u, accessToken, s.log)
		if tagErr != nil {
			return nil, tagErr
		}
		tl.Tag = append(tl.Tag, *nt)
	}
	return &tl, nil
}

func postTag(ta TagAction, url string, accessToken string, logger log.Logger) (*Tag, error) {

	ts := TagSet{UserID: ta.UserID, Untag: (ta.Action == "untag")}
	tsl := []TagSet{ts}
	reqBody := TagSetBody{Name: ta.Name, Users: tsl}

	httpBody, _ := json.Marshal(reqBody)

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:     url,
			AcceptType:  "application/json",
			BearerToken: &api.BearerTokenConfig{AccessToken: accessToken},
		},
		logger)

	// Send request
	tag := &Tag{}

	err := client.Post("tags", httpBody, tag)
	if err != nil {
		fmt.Println("INTERCOM POST ERR: ", err)
		return nil, err
	}

	return tag, nil
}
