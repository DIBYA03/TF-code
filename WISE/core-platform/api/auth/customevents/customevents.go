// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package customevents

import (
	"github.com/aws/aws-lambda-go/events"
)

// CognitoEventUserPoolsCustomMessage is sent by AWS Cognito User Pools when a user attempts to register
// (sign up), allowing a Lambda to perform custom validation to accept or deny the registration request
type CognitoEventUserPoolsCustomMessage struct {
	events.CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsCustomMessageRequest  `json:"request"`
	Response CognitoEventUserPoolsCustomMessageResponse `json:"response"`
}

// CognitoEventUserPoolsDefineChallenge is sent by AWS Cognito User Pools when a user attempts to register
// (sign up), allowing a Lambda to perform custom validation to accept or deny the registration request
type CognitoEventUserPoolsDefineChallenge struct {
	events.CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsDefineChallengeRequest  `json:"request"`
	Response CognitoEventUserPoolsDefineChallengeResponse `json:"response"`
}

// CognitoEventUserPoolsCreateChallenge is sent by AWS Cognito User Pools when a user attempts to register
// (sign up), allowing a Lambda to perform custom validation to accept or deny the registration request
type CognitoEventUserPoolsCreateChallenge struct {
	events.CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsCreateChallengeRequest  `json:"request"`
	Response CognitoEventUserPoolsCreateChallengeResponse `json:"response"`
}

// CognitoEventUserPoolsVerifyChallenge is sent by AWS Cognito User Pools when a user attempts to register
// (sign up), allowing a Lambda to perform custom validation to accept or deny the registration request
type CognitoEventUserPoolsVerifyChallenge struct {
	events.CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsVerifyChallengeRequest  `json:"request"`
	Response CognitoEventUserPoolsVerifyChallengeResponse `json:"response"`
}

// CognitoEventUserPoolsCustomMessageRequest contains the request portion of a PostConfirmation event
type CognitoEventUserPoolsCustomMessageRequest struct {
	UserAttributes    map[string]string `json:"userAttributes"`
	CodeParameter     string            `json:"codeParameter"`
	UsernameParameter string            `json:"usernameParameter"`
}

// CognitoEventUserPoolsCustomMessageResponse contains the response portion of a PostConfirmation event
type CognitoEventUserPoolsCustomMessageResponse struct {
	SmsMessage   string `json:"smsMessage"`
	EmailMessage string `json:"emailMessage"`
	EmailSubject string `json:"emailSubject"`
}

// CognitoEventUserPoolsDefineChallengeRequest contains the request portion of a PostConfirmation event
type CognitoEventUserPoolsDefineChallengeRequest struct {
	UserAttributes    map[string]string `json:"userAttributes"`
	Session           map[string]string `json:"session"`
	UsernameParameter string            `json:"usernameParameter"`
}

// CognitoEventUserPoolsDefineChallengeResponse contains the response portion of a PostConfirmation event
type CognitoEventUserPoolsDefineChallengeResponse struct {
	ChallengeName      string `json:"challengeName"`
	IssueTokens        bool   `json:"issueTokens"`
	FailAuthentication bool   `json:"failAuthentication"`
}

// CognitoEventUserPoolsCreateChallengeRequest contains the request portion of a PostConfirmation event
type CognitoEventUserPoolsCreateChallengeRequest struct {
	UserAttributes    map[string]string `json:"userAttributes"`
	ChallengeName     string            `json:"challengeName"`
	Session           map[string]string `json:"session"`
	UsernameParameter string            `json:"usernameParameter"`
}

// CognitoEventUserPoolsCreateChallengeResponse contains the response portion of a PostConfirmation event
type CognitoEventUserPoolsCreateChallengeResponse struct {
	PublicChallengeParameters  map[string]string `json:"publicChallengeParameters"`
	PrivateChallengeParameters map[string]string `json:"privateChallengeParameters"`
	ChallengeMetadata          string            `json:"challengeMetadata"`
}

// CognitoEventUserPoolsVerifyChallengeRequest contains the request portion of a PostConfirmation event
type CognitoEventUserPoolsVerifyChallengeRequest struct {
	UserAttributes             map[string]string `json:"userAttributes"`
	PrivateChallengeParameters map[string]string `json:"privateChallengeParameters"`
	ChallengeAnswer            map[string]string `json:"challengeAnswer"`
	UsernameParameter          string            `json:"usernameParameter"`
}

// CognitoEventUserPoolsVerifyChallengeResponse contains the response portion of a PostConfirmation event
type CognitoEventUserPoolsVerifyChallengeResponse struct {
	AnswerCorrect bool `json:"answerCorrect"`
}
