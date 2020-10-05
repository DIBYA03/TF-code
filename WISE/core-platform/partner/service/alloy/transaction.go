/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package alloy

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services"
)

type TransactionMonitorRequest struct {
	UserToken       string  `json:"user_token"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Description     string  `json:"description"`
	TransactionType string  `json:"transaction_type"`
	Source          Account `json:"source"`
	Destination     Account `json:"destination"`
}

type Account struct {
	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
	Email         string `json:"email_address"`
	StreetAddress string `json:"address_line_1"`
	AddressLine2  string `json:"address_line_2"`
	City          string `json:"address_city"`
	State         string `json:"address_state"`
	PostalCode    string `json:"address_postal_code"`
	Country       string `json:"address_country_code"`
}

type TransactionMonitorResponse struct {
	StatusCode           int            `json:"status_code"`
	Error                string         `json:"error"`
	TimeStamp            uint64         `json:"timestamp"`
	EvaluationToken      string         `json:"evaluation_token"`
	EntityToken          string         `json:"entity_token"`
	ParentEntityToken    string         `json:"parent_entity_token"`
	ApplicationToken     string         `json:"application_token"`
	ApplicationVersion   int            `json:"application_version_id"`
	ChampionChallengerID string         `json:"champion_challenger_id"`
	Summary              Summary        `json:"summary"`
	Supplied             types.JSONText `json:"supplied"`
	Formatted            types.JSONText `json:"formatted"`
	Matching             types.JSONText `json:"matching"`
	Meta                 types.JSONText `json:"meta"`
	Diligence            types.JSONText `json:"diligence"`
	RelatedData          types.JSONText `json:"related_data"`
	RawResponse          types.JSONText `json:"raw_responses"`
	FormattedResponse    types.JSONText `json:"formatted_responses"`
	AuditArchive         string         `json:"audit_archive"`
	FoundData            types.JSONText `json:"found_data"`
}

type Summary struct {
	Result          string   `json:"result"`
	Score           float64  `json:"score"`
	Outcome         string   `json:"outcome"`
	AlloyFraudScore *float64 `json:"alloy_fraud_score"`
}

func (r *TransactionMonitorResponse) Raw() []byte {
	rw, _ := json.Marshal(r)
	return rw
}

type alloyService struct {
	request services.SourceRequest
	client  *client
}

type AlloyService interface {
	MonitorTransaction(TransactionMonitorRequest) (*TransactionMonitorResponse, error)
}

func NewAlloyService(r services.SourceRequest) AlloyService {
	return &alloyService{r, newClient()}
}

func (a *alloyService) MonitorTransaction(request TransactionMonitorRequest) (*TransactionMonitorResponse, error) {
	// Call evaluation alloy API
	path := "v1/evaluations"

	req, err := a.client.post(path, request)
	if err != nil {
		return nil, err
	}
	var resp = TransactionMonitorResponse{}
	if err := a.client.do(req, &resp); err != nil {
		log.Println("error in response", err)
		return nil, err
	}

	return &resp, nil

}
