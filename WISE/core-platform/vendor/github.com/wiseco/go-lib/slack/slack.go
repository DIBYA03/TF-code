/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package slack

import (
	"bytes"
	"encoding/json"

	"github.com/wiseco/go-lib/api"
	"github.com/wiseco/go-lib/log"
	"github.com/wiseco/go-lib/url"
)

type Message struct {
	Text string `json:"text"`
}

type slackService struct {
	log log.Logger
}

type SlackService interface {
	PostToChannel(string, Message) error
}

func NewSlackService(l log.Logger) SlackService {
	return &slackService{log: l}
}

func (s *slackService) PostToChannel(webhookURL string, m Message) error {
	// URL
	u, err := url.BuildAbsoluteForSlack("", url.Params{})
	if err != nil {
		return err
	}

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:     u,
			ContentType: "application/json",
		},
		s.log)

	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	var resp string

	err = client.Post("services/"+webhookURL, bytes.NewReader(body), resp)
	if err != nil {
		return err
	}

	return nil
}
