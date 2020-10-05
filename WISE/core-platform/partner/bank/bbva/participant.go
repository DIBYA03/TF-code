package bbva

import (
	"strings"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type ParticipantRole string

const (
	// A person who is authorized to perform business-related transactions on behalf of another person.
	ParticipantRoleAttorneyInFact = ParticipantRole("attorney_in_fact")

	// Allowed user by the account owner to operate the account. This value is equivalent to SIGNER or AUTHORIZED SIGNER.
	ParticipantRoleAuthorized = ParticipantRole("authorized")

	// A person who gains an asset upon the holder's death.
	ParticipantRoleBeneficiary = ParticipantRole("beneficiary")

	// An adult that manages an account for a minor under the age of 18, depending on state laws. The product may or not
	// require adult approval for certain transactions.
	ParticipantRoleCustodian = ParticipantRole("custodian")

	// A person or institution that has been appointed to manage the assests within an account.
	ParticipantRoleConservator = ParticipantRole("conservator")

	// It refers to the account owner.
	ParticipantRoleHolder = ParticipantRole("holder")

	// A person who is under the age of legal competence.
	ParticipantRoleMinorWard = ParticipantRole("minor_ward")

	// A person appointed by the Social Security Administration to manage Social Security benefits for someone that is unable.
	ParticipantRoleRepresentativePayee = ParticipantRole("representative_payee")
)

type ParticipantRequest struct {
	UserID string          `json:"participant_user_id"`
	Role   ParticipantRole `json:"participant_role"`
}

var partnerParticipantRoleFrom = map[bank.AccountRole]ParticipantRole{
	bank.AccountRoleHolder:     ParticipantRoleHolder,
	bank.AccountRoleAuthorized: ParticipantRoleAuthorized,
	bank.AccountRoleAttorney:   ParticipantRoleAttorneyInFact,
	bank.AccountRoleSpouse:     ParticipantRoleAuthorized,
	bank.AccountRoleMinor:      ParticipantRoleMinorWard,
	bank.AccountRoleCustodian:  ParticipantRoleCustodian,
}

var partnerParticipantRoleTo = map[ParticipantRole]bank.AccountRole{
	ParticipantRoleHolder:         bank.AccountRoleHolder,
	ParticipantRoleAuthorized:     bank.AccountRoleAuthorized,
	ParticipantRoleAttorneyInFact: bank.AccountRoleAttorney,
	ParticipantRoleMinorWard:      bank.AccountRoleMinor,
	ParticipantRoleCustodian:      bank.AccountRoleCustodian,
}

type AccountParticipantsResponse struct {
	Participants []AccountParticipantResponse `json:"participants"`
}

type AccountParticipantResponse struct {
	UserID    string          `json:"participant_user_id"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Role      ParticipantRole `json:"role"`
}

func (resp *AccountParticipantResponse) partnerParticipantResponseTo(req bank.APIRequest) (*bank.AccountParticipantResponse, error) {
	c, err := data.NewConsumerService(req, bank.ProviderNameBBVA).GetByBankID(bank.ConsumerBankID(resp.UserID))
	if err != nil {
		return nil, err
	}

	return &bank.AccountParticipantResponse{
		ConsumerID: c.ConsumerID,
		Role:       partnerParticipantRoleTo[resp.Role],
	}, nil
}

func partnerParticipantsTo(req bank.APIRequest, presp AccountParticipantsResponse) ([]bank.AccountParticipantResponse, error) {
	pp := []bank.AccountParticipantResponse{}
	for _, p := range presp.Participants {
		// Skip businesses
		if strings.HasPrefix(p.UserID, EntityPrefixBusiness) {
			continue
		}

		participant, err := p.PartnerParticipantResponseTo(req)
		if err != nil {
			return nil, err
		}

		pp = append(pp, *participant)
	}

	return pp, nil
}

func (resp *AccountParticipantResponse) PartnerParticipantResponseTo(req bank.APIRequest) (*bank.AccountParticipantResponse, error) {
	c, err := data.NewConsumerService(req, bank.ProviderNameBBVA).GetByBankID(bank.ConsumerBankID(resp.UserID))
	if err != nil {
		return nil, err
	}

	return &bank.AccountParticipantResponse{
		ConsumerID: c.ConsumerID,
		Role:       partnerParticipantRoleTo[resp.Role],
	}, nil
}
