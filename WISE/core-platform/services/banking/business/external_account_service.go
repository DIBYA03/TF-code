package business

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/services"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/plaid"
)

type externalAccountDataStore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type ExternalAccountService interface {
	Upsert(ExternalBankAccountUpdate) (*ExternalBankAccount, error)
	GetByID(string) (*ExternalBankAccount, error)
	GetByPartnerAccountID(string, shared.BusinessID) (*ExternalBankAccount, error)
	GetByAccountNumberInternal(string, string, shared.BusinessID) (*ExternalBankAccount, error)

	CreateOwners(string, []ExternalBankAccountOwnerCreate) error
	ListOwnersByAccountID(string) ([]ExternalBankAccountOwner, error)

	Verify(ExternalAccountVerificationRequest) error
	GetVerificationByAccountID(string, shared.BusinessID) (*ExternalAccountVerificationResult, error)
}

func NewExternalAccountService(r services.SourceRequest) ExternalAccountService {
	return &externalAccountDataStore{r, data.DBWrite}
}

func (db *externalAccountDataStore) Upsert(cu ExternalBankAccountUpdate) (*ExternalBankAccount, error) {
	account := &ExternalBankAccount{}

	err := db.Get(account, "SELECT * FROM external_bank_account WHERE business_id = $1 AND account_number = $2 AND routing_number = $3",
		cu.BusinessID, cu.AccountNumber, cu.RoutingNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		columns := []string{
			"business_id", "linked_account_id", "partner_account_id", "partner_name", "account_name", "official_account_name",
			"account_type", "account_subtype", "account_number", "routing_number",
			"wire_routing", "available_balance", "posted_balance", "currency", "last_login",
		}

		values := []string{
			":business_id", ":linked_account_id", ":partner_account_id", ":partner_name", ":account_name", ":official_account_name",
			":account_type", ":account_subtype", ":account_number", ":routing_number",
			":wire_routing", ":available_balance", ":posted_balance", ":currency", ":last_login",
		}

		sql := fmt.Sprintf("INSERT INTO external_bank_account(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

		stmt, err := db.PrepareNamed(sql)
		if err != nil {
			return nil, err
		}

		err = stmt.Get(account, &cu)
		if err != nil {
			return nil, err
		}
	} else {
		var columns []string
		if cu.PartnerAccountID != nil {
			columns = append(columns, "partner_account_id = :partner_account_id")
		}

		if cu.AvailableBalance != nil {
			columns = append(columns, "available_balance = :available_balance")
		}

		if cu.PostedBalance != nil {
			columns = append(columns, "posted_balance = :posted_balance")
		}

		if cu.Currency != nil {
			columns = append(columns, "currency = :currency")
		}

		if cu.LastLogin != nil {
			columns = append(columns, "last_login = :last_login")
		}

		if cu.LinkedAccountID != nil {
			columns = append(columns, "linked_account_id = :linked_account_id")
		}

		// Update external bank account
		_, err = db.NamedExec(
			fmt.Sprintf(
				"UPDATE external_bank_account SET %s WHERE id = '%s'",
				strings.Join(columns, ", "),
				account.ID,
			), cu,
		)

		if err != nil {
			return nil, errors.Cause(err)
		}

		account, err = db.GetByID(account.ID)
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}

func (db *externalAccountDataStore) GetByID(ID string) (*ExternalBankAccount, error) {
	account := &ExternalBankAccount{}

	err := db.Get(account, "SELECT * FROM external_bank_account WHERE id = $1", ID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (db *externalAccountDataStore) GetByPartnerAccountID(partnerAccountID string, bID shared.BusinessID) (*ExternalBankAccount, error) {
	account := &ExternalBankAccount{}

	err := db.Get(account, "SELECT * FROM external_bank_account WHERE partner_account_id = $1 AND business_id = $2", partnerAccountID, bID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (db *externalAccountDataStore) GetByAccountNumberInternal(accountNumber, routingNumber string, bID shared.BusinessID) (*ExternalBankAccount, error) {
	account := &ExternalBankAccount{}

	err := db.Get(account, `SELECT * FROM external_bank_account WHERE account_number = $1 AND 
	routing_number = $2 AND business_id = $3`, accountNumber, routingNumber, bID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (db *externalAccountDataStore) CreateOwners(accountID string, owners []ExternalBankAccountOwnerCreate) error {
	d := fmt.Sprintf("DELETE FROM external_bank_account_owner WHERE external_bank_account_id = '%s'", accountID)
	_, err := db.Exec(d)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	for _, o := range owners {
		columns := []string{
			"external_bank_account_id", "account_holder_name", "phone", "email", "owner_address",
		}

		values := []string{
			":external_bank_account_id", ":account_holder_name", ":phone", ":email", ":owner_address",
		}

		sql := fmt.Sprintf("INSERT INTO external_bank_account_owner(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

		stmt, err := db.PrepareNamed(sql)
		if err != nil {
			return err
		}

		o.ExternalAccountID = accountID

		owner := &ExternalBankAccountOwner{}
		err = stmt.Get(owner, &o)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *externalAccountDataStore) ListOwnersByAccountID(accountID string) ([]ExternalBankAccountOwner, error) {
	rows := []ExternalBankAccountOwner{}

	err := db.Select(&rows, "SELECT * FROM external_bank_account_owner WHERE external_bank_account_id = $1", accountID)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

/*
1. External account holder's name should either match with business members' name or business legal name or DBA
2. External account holder's email or phone or address should match with business members/business' email/phone/address
*/
func (db *externalAccountDataStore) Verify(req ExternalAccountVerificationRequest) error {
	result := &services.StringArray{}
	var status VerificationStatus

	ea, err := db.GetByPartnerAccountID(req.PartnerAccountID, req.BusinessID)
	if err != nil {
		return err
	}

	owners, err := db.ListOwnersByAccountID(ea.ID)
	if err != nil {
		return err
	}

	result, status, err = db.verifyOwnerName(req.BusinessID, owners, *result)
	if err != nil {
		return err
	}

	if status == VerificationStatusSucceeded {
		result, status, err = db.verifyOwnerPhone(req.BusinessID, owners, *result)
		if err != nil {
			return err
		}

		if status == VerificationStatusUnverified {
			result, status, err = db.verifyOwnerEmail(req.BusinessID, owners, *result)
			if err != nil {
				return err
			}

			if status == VerificationStatusUnverified {
				result, status, err = db.verifyOwnerAddress(req.BusinessID, owners, *result)
				if err != nil {
					return err
				}
			}
		}
	}

	c := &ExternalAccountVerificationCreate{
		BusinessID:         req.BusinessID,
		ExternalAccountID:  ea.ID,
		SourceIPAddress:    db.sourceReq.SourceIP,
		PartnerItemID:      req.PartnerItemID,
		AccessToken:        req.AccessToken,
		VerificationStatus: status,
		VerificationResult: *result,
	}
	v, err := db.createExternalBankAccountVerification(c)
	if err != nil {
		return err
	}

	if v.VerificationStatus == VerificationStatusUnverified {
		result := strings.Join(v.VerificationResult, " ")

		if strings.Contains(result, NameMismatch) {
			return errors.New(ExternalVerificationErrorNameMismatch)
		} else if strings.Contains(result, AddressMismatch) {
			return errors.New(ExternalVerificationErrorAddressMismatch)
		} else if strings.Contains(result, EmailMismatch) {
			return errors.New(ExternalVerificationErrorEmailMismatch)
		} else if strings.Contains(result, PhoneMismatch) {
			return errors.New(ExternalVerificationErrorPhoneMismatch)
		} else {
			return errors.New(ExternalVerificationErrorGeneric)
		}
	}

	return nil
}

func (db *externalAccountDataStore) verifyOwnerName(bID shared.BusinessID, owners []ExternalBankAccountOwner,
	result services.StringArray) (*services.StringArray, VerificationStatus, error) {

	nameMatch, err := db.checkBusinessMemberName(bID, owners)
	if err != nil {
		return nil, VerificationStatusUnverified, err
	}

	if nameMatch {
		result = append(result, BusinessMemberNameMatch)
		return &result, VerificationStatusSucceeded, nil
	}

	nameMatch, err = db.checkBusinessName(bID, owners)
	if err != nil {
		return nil, VerificationStatusUnverified, err
	}

	if nameMatch {
		result = append(result, BusinessNameMatch)
		return &result, VerificationStatusSucceeded, nil
	}

	result = append(result, NameMismatch)
	return &result, VerificationStatusUnverified, nil
}

func (db *externalAccountDataStore) verifyOwnerPhone(bID shared.BusinessID, owners []ExternalBankAccountOwner,
	result services.StringArray) (*services.StringArray, VerificationStatus, error) {

	members, err := bsrv.NewMemberService(db.sourceReq).List(0, 10, bID)
	if err != nil {
		return nil, VerificationStatusUnverified, err
	}

	for _, member := range members {
		if len(member.Phone) < 12 {
			continue
		}

		// Trim prefix +1
		memberPhone := member.Phone[2:]

		for _, o := range owners {
			var phones []plaid.Phone

			err := json.Unmarshal(o.Phone, &phones)
			if err != nil {
				return nil, VerificationStatusUnverified, err
			}

			for _, phone := range phones {
				// Trim prefix +1 if exists
				if strings.HasPrefix(phone.Phone, "+1") && len(phone.Phone) > 2 {
					phone.Phone = phone.Phone[2:]
				}

				if phone.Phone == memberPhone {
					result = append(result, BusinessMemberPhoneMatch)
					return &result, VerificationStatusSucceeded, nil
				}
			}
		}
	}

	b, err := bsrv.NewBusinessService(db.sourceReq).GetById(bID)
	if err != nil {
		return nil, VerificationStatusUnverified, err
	}

	if b.Phone == nil {
		result = append(result, PhoneMismatch)
		return &result, VerificationStatusUnverified, nil
	}

	if len(*b.Phone) < 12 {
		result = append(result, PhoneMismatch)
		return &result, VerificationStatusUnverified, nil
	}

	// Trim prefix +1
	businessPhone := (*b.Phone)[2:]

	for _, o := range owners {
		var phoneArr []plaid.Phone

		err := json.Unmarshal(o.Phone, &phoneArr)
		if err != nil {
			return nil, VerificationStatusUnverified, err
		}

		for _, phone := range phoneArr {
			if phone.Phone == businessPhone {
				result = append(result, BusinessPhoneMatch)
				return &result, VerificationStatusSucceeded, nil
			}
		}
	}

	result = append(result, PhoneMismatch)
	return &result, VerificationStatusUnverified, nil
}

func (db *externalAccountDataStore) verifyOwnerEmail(bID shared.BusinessID, owners []ExternalBankAccountOwner,
	result services.StringArray) (*services.StringArray, VerificationStatus, error) {

	members, err := bsrv.NewMemberService(db.sourceReq).List(0, 10, bID)
	if err != nil {
		return nil, VerificationStatusUnverified, err
	}

	for _, member := range members {
		for _, o := range owners {
			var emails []plaid.Email

			err := json.Unmarshal(o.Email, &emails)
			if err != nil {
				return nil, VerificationStatusUnverified, err
			}

			for _, email := range emails {
				if email.Email == member.Email {
					result = append(result, BusinessMemberEmailMatch)
					return &result, VerificationStatusSucceeded, nil
				}
			}
		}
	}

	b, err := bsrv.NewBusinessService(db.sourceReq).GetById(bID)
	if err != nil {
		return nil, VerificationStatusUnverified, err
	}

	if b.Email == nil {
		result = append(result, EmailMismatch)
		return &result, VerificationStatusUnverified, nil
	}

	for _, o := range owners {
		var emails []plaid.Email

		err := json.Unmarshal(o.Email, &emails)
		if err != nil {
			return nil, VerificationStatusUnverified, err
		}

		for _, email := range emails {
			if email.Email == *b.Email {
				result = append(result, BusinessEmailMatch)
				return &result, VerificationStatusSucceeded, nil
			}
		}
	}

	result = append(result, EmailMismatch)
	return &result, VerificationStatusUnverified, nil
}

func (db *externalAccountDataStore) verifyOwnerAddress(bID shared.BusinessID, owners []ExternalBankAccountOwner,
	result services.StringArray) (*services.StringArray, VerificationStatus, error) {

	members, err := bsrv.NewMemberService(db.sourceReq).List(0, 10, bID)
	if err != nil {
		return nil, VerificationStatusUnverified, err
	}

	for _, member := range members {
		for _, o := range owners {

			var addresses []plaid.Address

			err := json.Unmarshal(o.OwnerAddress, &addresses)
			if err != nil {
				return nil, VerificationStatusUnverified, err
			}

			for _, address := range addresses {
				// Member legal address check
				if addressCheck(member.LegalAddress, address) {
					result = append(result, BusinessMemberLegalAddressMatch)
					return &result, VerificationStatusSucceeded, nil
				}

				// Member mailing address check
				if addressCheck(member.MailingAddress, address) {
					result = append(result, BusinessMemberMailingAddressMatch)
					return &result, VerificationStatusSucceeded, nil
				}

				// Member work address check
				if addressCheck(member.WorkAddress, address) {
					result = append(result, BusinessMemberWorkAddressMatch)
					return &result, VerificationStatusSucceeded, nil
				}
			}
		}
	}

	b, err := bsrv.NewBusinessService(db.sourceReq).GetById(bID)
	if err != nil {
		return nil, VerificationStatusUnverified, err
	}

	for _, o := range owners {
		var addresses []plaid.Address

		err := json.Unmarshal(o.OwnerAddress, &addresses)
		if err != nil {
			return nil, VerificationStatusUnverified, err
		}

		for _, address := range addresses {
			// Business legal address check
			if addressCheck(b.LegalAddress, address) {
				result = append(result, BusinessLegalAddressMatch)
				return &result, VerificationStatusSucceeded, nil
			}

			// Business mailing address check
			if addressCheck(b.MailingAddress, address) {
				result = append(result, BusinessMailingAddressMatch)
				return &result, VerificationStatusSucceeded, nil
			}

			// Business headquarter address check
			if addressCheck(b.HeadquarterAddress, address) {
				result = append(result, BusinessHeadquarterAddressMatch)
				return &result, VerificationStatusSucceeded, nil
			}
		}
	}

	result = append(result, AddressMismatch)
	return &result, VerificationStatusUnverified, nil
}

func (db *externalAccountDataStore) checkBusinessMemberName(bID shared.BusinessID, owners []ExternalBankAccountOwner) (bool, error) {
	members, err := bsrv.NewMemberService(db.sourceReq).List(0, 10, bID)

	if err != nil {
		return false, err
	}

	for _, member := range members {
		for _, o := range owners {
			for _, name := range o.AccountHolderName {
				if nameCheck(member, name) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (db *externalAccountDataStore) checkBusinessName(bID shared.BusinessID, owners []ExternalBankAccountOwner) (bool, error) {
	b, err := bsrv.NewBusinessService(db.sourceReq).GetById(bID)
	if err != nil {
		return false, err
	}

	var legalName string
	if b.LegalName != nil {
		legalName = *b.LegalName
	}

	var dba string
	if b.DBA != nil && len(b.DBA) > 0 && len(b.DBA[0]) > 0 {
		dba = b.DBA[0]
	}

	for _, o := range owners {
		for _, name := range o.AccountHolderName {
			if strings.EqualFold(name, legalName) {
				return true, nil
			}

			if strings.EqualFold(name, dba) {
				return true, nil
			}
		}
	}

	return false, nil
}

func addressCheck(address *services.Address, ownerAddress plaid.Address) bool {
	if address == nil {
		return false
	}

	if !strings.EqualFold(address.StreetAddress, ownerAddress.StreetAddress) {
		return false
	}

	if !strings.EqualFold(address.AddressLine2, ownerAddress.AddressLine2) {
		return false
	}

	if !strings.EqualFold(address.City, ownerAddress.City) {
		return false
	}

	if !strings.EqualFold(address.State, ownerAddress.State) {
		return false
	}

	if !strings.EqualFold(address.Country, ownerAddress.Country) {
		return false
	}

	if !strings.EqualFold(address.PostalCode, ownerAddress.PostalCode) {
		return false
	}

	return true
}

func nameCheck(member bsrv.BusinessMember, accountOwnerName string) bool {
	memberName := member.FirstName
	if member.MiddleName != "" {
		memberName = memberName + " " + member.MiddleName
	}
	memberName = memberName + " " + member.LastName

	if strings.EqualFold(memberName, accountOwnerName) {
		return true
	}

	memberName = member.FirstName + " " + member.LastName

	ownerNames := strings.Split(accountOwnerName, " ")
	if len(ownerNames) < 2 {
		return false
	}

	if strings.EqualFold(memberName, ownerNames[0]+" "+ownerNames[len(ownerNames)-1]) {
		return true
	}

	return false
}

func (db *externalAccountDataStore) createExternalBankAccountVerification(c *ExternalAccountVerificationCreate) (*ExternalAccountVerificationResult, error) {

	// Default/mandatory fields
	columns := []string{
		"business_id", "external_bank_account_id", "source_ip_address", "access_token", "partner_item_id", "verification_status", "verification_result",
	}

	// Default/mandatory values
	values := []string{
		":business_id", ":external_bank_account_id", ":source_ip_address", ":access_token", ":partner_item_id", ":verification_status", ":verification_result",
	}

	sql := fmt.Sprintf("INSERT INTO external_account_verification_result(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	v := &ExternalAccountVerificationResult{}

	err = stmt.Get(v, &c)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return v, nil
}

func (db *externalAccountDataStore) GetVerificationByAccountID(accountID string, bID shared.BusinessID) (*ExternalAccountVerificationResult, error) {
	result := &ExternalAccountVerificationResult{}

	err := db.Get(result, "SELECT * FROM external_account_verification_result WHERE external_bank_account_id = $1 AND business_id = $2", accountID, bID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
