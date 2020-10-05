package business

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/service/segment"
	"github.com/wiseco/core-platform/services"
	banksrv "github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	acct "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/cspuser"
	"github.com/wiseco/core-platform/services/csp/mail"
	cspServices "github.com/wiseco/core-platform/services/csp/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	grpcBankTransfer "github.com/wiseco/protobuf/golang/banking/transfer"
)

//BankAccount ...
type BankAccount struct {
	// Business id
	BusinessID string `json:"businessId" db:"business_id"`

	// Account usage type
	UsageType acct.UsageType `json:"usageType" db:"usage_type"`

	banksrv.BankAccount
}

//Banking banking service
type Banking interface {
	CreateBankAccount(shared.BusinessID) error
	CreateCard(shared.BusinessID) error
	ExternalAccountByBusinessID(shared.BusinessID) ([]ExternalBankAccount, error)
	AccountByBusinessID(shared.BusinessID) ([]business.BankAccount, error)
	Cards(shared.BusinessID) ([]business.BankCardPartial, error)
	CardByID(shared.BusinessID, string) (business.BankCardPartial, error)
	FundPromotion(shared.BusinessID) error

	ReissueCard(business.CardReissueRequest) (*business.BankCardPartial, error)
	CreateReissueHistory(business.CardReissueRequest) (*business.CardReissueHistory, error)
	ReissueHistoryByCardID(string) ([]business.CardReissueHistory, error)

	ApproveTransferInReview(string, string, string) error
	DeclineTransferInReview(string, string, string) error

	GetTransferForTransactionID(string) (*grpcBankTransfer.Transfer, error)
}

type banking struct {
	rdb       *sqlx.DB
	wdb       *sqlx.DB
	sourceReq services.SourceRequest
}

//NewBanking returns a new banking service
func NewBanking(sourceReq services.SourceRequest) Banking {
	return banking{wdb: data.DBWrite, rdb: data.DBRead, sourceReq: sourceReq}
}

func (service banking) AccountByBusinessID(id shared.BusinessID) ([]business.BankAccount, error) {
	var list []business.BankAccount

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := acct.NewBankingAccountService()
		if err != nil {
			return list, err
		}

		ba, err := bas.GetAllPrimaryByBusinessID(id, 100, 0)

		return ba, err
	}

	err := service.rdb.Select(&list, "SELECT * FROM business_bank_account WHERE business_id = $1", id)
	if err != nil && err == sql.ErrNoRows {
		return list, services.ErrorNotFound{}.New("")
	}
	return list, err
}

func (service banking) ExternalAccountByBusinessID(id shared.BusinessID) ([]ExternalBankAccount, error) {
	list := []ExternalBankAccount{}
	err := service.rdb.Select(&list, "SELECT * FROM external_bank_account WHERE business_id = $1", id)
	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}

	for i, _ := range list {
		owners := []business.ExternalBankAccountOwner{}
		err := service.rdb.Select(&owners, "SELECT * FROM external_bank_account_owner WHERE external_bank_account_id = $1", list[i].ID)
		if err != nil && err != sql.ErrNoRows {
			return list, err
		}

		list[i].Owners = owners
	}

	return list, err
}

func (service banking) CreateBankAccount(id shared.BusinessID) error {
	b := banksrv.BankAccountCreate{
		BankName:    banksrv.BankNameBBVA,
		AccountType: "checking",
	}
	act := &acct.BankAccountCreate{id, acct.UsageTypePrimary, b}
	resp, err := acct.NewBankAccountService(service.sourceReq).Create(act)

	if err != nil {
		log.Printf("error creating account %v", err)

		uerr := NewCSPService().UpdateProcessStatus(id, csp.AccountCreationFailed)
		if uerr != nil {
			err = errors.Wrap(err, uerr.Error())
		}

		return err
	}

	err = NewCSPService().UpdateProcessStatus(id, csp.AccountCreated)
	if err != nil {
		log.Printf("error updating csp process status %v", err)
		return err
	}

	// Set subscription status for accounts originated post July 1st
	tz := os.Getenv("BATCH_TZ")
	if tz == "" {
		panic(errors.New("Local timezone missing"))
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	nowUTC := time.Now().UTC()
	nowLocal := nowUTC.In(loc)

	// Subscription starts from July 1st
	wiseSubscriptionStartDate := time.Date(2020, time.July, 1, 0, 0, 0, 0, loc)

	if wiseSubscriptionStartDate.Before(nowLocal) {

		subscriptionStartDate := time.Date(nowLocal.Year(), nowLocal.Month(), 1, 0, 0, 0, 0, loc)
		subscriptionStartDate = subscriptionStartDate.AddDate(0, 1, 0)

		log.Println("Subscription start date is ", subscriptionStartDate)

		status := services.SubscriptionStatusActive
		subscriptionStartSharedDate := shared.Date(subscriptionStartDate)

		su := SubscriptionUpdate{
			BusinessID:            id,
			SubscriptionStatus:    &status,
			SubscriptionStartDate: &subscriptionStartSharedDate,
		}
		NewSubscriptionService(cspServices.NewSourceRequest()).Update(su)
	}

	// Send to segment
	segment.NewSegmentService().PushToAnalyticsQueue(resp.AccountHolderID, segment.CategoryAccount, segment.ActionCreate, resp)

	return err
}

//TODO this needs to be ported 100% to banking service, migrating and backfilling schema as well
func (service banking) CreateCard(bID shared.BusinessID) error {

	var bus struct {
		ID              string               `db:"id"`
		BankAccountID   string               `db:"bank_account_id"`
		AccountHolderID shared.UserID        `db:"account_holder_id"`
		DBA             services.StringArray `json:"dba" db:"dba"`
		LegalName       *string              `db:"legal_name"`
		BusinessID      shared.BusinessID    `db:"business_id"`
		UsageType       business.UsageType   `db:"usage_type"`
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		err := service.rdb.Get(
			&bus, `
			SELECT
				business.dba,
				business.legal_name,
				business.id as business_id,
				business.owner_id as account_holder_id
			FROM business
			WHERE business.id = $1`,
			bID,
		)
		if err != nil {
			log.Printf("error getting business record %v", err)
			return err
		}

		bas, err := business.NewBankingAccountService()
		if err != nil {
			return err
		}

		bs, err := bas.GetByBusinessID(bID, 1, 0)
		if err != nil {
			return err
		}

		if len(bs) != 1 {
			return errors.New("Unable to find bank account for business id")
		}

		ba := bs[0]

		baID, err := id.ParseBankAccountID(ba.Id)
		if err != nil {
			return err
		}

		bus.ID = baID.UUIDString()
		bus.BankAccountID = ba.BankAccountId
		bus.UsageType = ba.UsageType
	} else {
		err := service.rdb.Get(
			&bus, `
			SELECT
				business_bank_account.id,
				business_bank_account.bank_account_id,
				business_bank_account.account_holder_id,
				business.dba,
				business.legal_name,
				business_bank_account.business_id,
				business_bank_account.usage_type
			FROM business_bank_account
			JOIN business ON business.id = business_bank_account.business_id
			WHERE business_bank_account.business_id = $1`,
			bID,
		)
		if err != nil {
			log.Printf("error getting bank account %v", err)
			return err
		}
	}

	// Only primary accounts can have a card
	if bus.UsageType != business.UsageTypePrimary {
		err := fmt.Errorf("cannot create card for business bank account usage type of: %s", bus.UsageType)
		log.Println(err)
		return err
	}

	// One debit card per user per account active
	var cards []business.BankCardPartial
	err := service.rdb.Select(
		&cards, `
		SELECT * FROM business_bank_card
		WHERE
			bank_account_id = $1 AND
			business_id = $2 AND
			cardholder_id = $3`,
		bus.ID,
		bus.BusinessID,
		bus.AccountHolderID,
	)

	if err != nil {
		return err
	}

	if len(cards) > 0 {
		err = errors.New("only one debit card allowed per user per account")
		log.Println(err)
		return err
	}

	var b struct {
		ConsumerID     shared.ConsumerID `db:"consumer_id"`
		FirstName      string            `db:"first_name"`
		MiddleName     string            `db:"middle_name"`
		LastName       string            `db:"last_name"`
		Phone          string            `db:"phone"`
		MailingAddress *services.Address `db:"mailing_address"`
	}

	err = service.rdb.Get(
		&b, `
		SELECT
			wise_user.consumer_id,
			consumer.first_name,
			consumer.middle_name,
			consumer.last_name,
			wise_user.phone,
			consumer.mailing_address
		FROM wise_user
		JOIN consumer ON consumer.id = wise_user.consumer_id
		WHERE wise_user.id = $1`,
		bus.AccountHolderID,
	)
	if err != nil {
		log.Print(err)
		return err
	}

	if b.MailingAddress == nil {
		return errors.New("Mailing address can not be empty")
	}

	busName := partnerbank.CardBusinessNameLegal
	if len(shared.GetDBAName(bus.DBA)) > 0 {
		busName = partnerbank.CardBusinessNameDBA
	}

	cardholderName, err := shared.GetVisaCardHolderName(b.FirstName, b.MiddleName, b.LastName)
	if err != nil {
		return err
	}

	//Register card with partner bank here
	request := partnerbank.CreateCardRequest{
		AccountID:      partnerbank.AccountBankID(bus.BankAccountID),
		CardholderName: cardholderName,
		Packaging:      partnerbank.CardPackagingRegular,
		Delivery:       partnerbank.CardDeliveryStandard,
		BusinessName:   busName,
		Phone:          partnerbank.PhoneE164(b.Phone),
		Address:        partnerbank.AddressRequestTypeMailing,
		Type:           partnerbank.CardType(CardTypeDebit.String()),
	}

	// Create card with partner bank
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return err
	}

	srv, err := bank.CardService(
		service.sourceReq.PartnerBankRequest(),
		partnerbank.BusinessID(bID),
		partnerbank.ConsumerID(b.ConsumerID),
	)
	if err != nil {
		return err
	}

	resp, err := srv.Create(request)
	if err != nil {
		log.Printf("error from create card partner service %v", err)
		NewCSPService().UpdateProcessStatus(bID, csp.CardCreationFailed)
		return err
	}
	c := transformCardResponse(resp)

	// Get card limits
	limits, err := srv.GetLimit(partnerbank.CardBankID(c.BankCardId))
	if err != nil {
		log.Printf("error getting card limits %v", err)
		return err
	}
	c.DailyTransactionLimit = &limits.DailyTransactionCount
	c.DailyATMLimit = &limits.DailyATMAmount
	c.DailyPOSLimit = &limits.DailyPOSAmount
	c.CardholderID = bus.AccountHolderID
	c.BusinessID = bID
	c.BankAccountId = bus.ID

	columns := []string{
		"cardholder_id", "business_id", "bank_account_id", "card_type", "cardholder_name", "is_virtual", "bank_name", "bank_card_id",
		"card_number_masked", "card_brand", "currency", "card_status", "alias", "daily_withdrawal_limit", "daily_pos_limit", "daily_transaction_limit",
		"card_number_alias",
	}

	values := []string{
		":cardholder_id", ":business_id", ":bank_account_id", ":card_type", ":cardholder_name", ":is_virtual", ":bank_name", ":bank_card_id",
		":card_number_masked", ":card_brand", ":currency", ":card_status", ":alias", ":daily_withdrawal_limit", ":daily_pos_limit", ":daily_transaction_limit",
		":card_number_alias",
	}

	sql := fmt.Sprintf("INSERT INTO business_bank_card(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := service.wdb.PrepareNamed(sql)
	if err != nil {
		log.Printf("error preparing insert statement %v", err)
		return err
	}
	err = NewCSPService().UpdateProcessStatus(bID, csp.CardCreated)
	if err != nil {
		return err
	}
	card := &business.BankCardPartial{}

	err = stmt.Get(card, &c)
	if err != nil {
		return fmt.Errorf("Error saving card err:%v", err.Error())
	}

	// Send to segment
	bankCard := banksrv.BankCardMini{
		BankCardID:            card.BankCardId,
		CardStatus:            card.CardStatus,
		DailyTransactionLimit: card.DailyTransactionLimit,
	}

	segment.NewSegmentService().PushToAnalyticsQueue(card.CardholderID, segment.CategoryCard, segment.ActionCreate, bankCard)
	if bus.LegalName != nil {
		mail.SendEmail(mail.EmailStatusApproved, bus.AccountHolderID, *bus.LegalName, *b.MailingAddress)
	} else {
		log.Printf("Error sending email, legal name of business is invalid businessID:%s", bus.ID)
		return fmt.Errorf("Error sending email, legal name of business is invalid businessID:%s ", bus.ID)
	}
	return err
}

func transformCardResponse(response *partnerbank.GetCardResponse) business.BankCard {
	c := business.BankCard{}
	c.BankName = "bbva"
	c.BankCardId = string(response.CardID)
	c.CardNumberMasked = response.PANMasked
	c.CardBrand = string(response.Brand)
	c.Currency = banksrv.Currency(response.Currency)
	c.CardStatus = banksrv.CardStatus(response.Status)
	c.CardholderName = response.CardholderName
	c.CardType = banksrv.CardType(response.Type)
	c.CardNumberAlias = response.PANAlias
	return c
}

func (service banking) Cards(businessID shared.BusinessID) ([]business.BankCardPartial, error) {
	list := make([]business.BankCardPartial, 0)
	err := service.rdb.Select(&list, "SELECT * FROM business_bank_card WHERE business_id = $1", businessID)
	if err == sql.ErrNoRows {
		return list, nil
	}
	if err != nil {
		return list, err
	}
	return list, nil
}

func (service banking) CardByID(businessID shared.BusinessID, cardID string) (business.BankCardPartial, error) {
	var card business.BankCardPartial
	err := service.rdb.Get(&card, "SELECT * FROM business_bank_card WHERE business_id = $1 AND id = $2", businessID, cardID)
	if err == sql.ErrNoRows {
		return card, services.ErrorNotFound{}.New("")
	}

	if err != nil {
		return card, err
	}

	history, err := service.ReissueHistoryByCardID(card.Id)
	if err != nil {
		return card, err
	}

	card.CardReissueHistory = history

	return card, nil
}

func (service banking) FundPromotion(destBID shared.BusinessID) error {
	if err := service.canPromoFund(destBID); err != nil {
		return err
	}
	// Wise Company User ID
	uID := os.Getenv("WISE_CLEARING_USER_ID")
	if uID == "" {
		return errors.New("user id required")
	}
	// Wise Company Business ID
	bID := os.Getenv("WISE_CLEARING_BUSINESS_ID")
	if bID == "" {
		return errors.New("business id required")
	}
	// Wise Company Clearing Linked Account
	srcID := os.Getenv("WISE_PROMO_CLEARING_LINKED_ACCOUNT_ID")
	if srcID == "" {
		return errors.New("linked promo clearing account id required")
	}

	bus, err := New(services.SourceRequest{}).ByID(destBID)
	if err != nil {
		return err
	}

	var ti *business.TransferInitiate

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		accs, err := business.NewAccountService().ListInternalByBusiness(destBID, 1, 0)
		if err != nil {
			return err
		}

		if len(accs) == 0 {
			return errors.New("Cannot find bank account")
		}

		acc := accs[0]

		saID, err := id.ParseLinkedBankAccountID(fmt.Sprintf("%s%s", id.IDPrefixLinkedBankAccount, srcID))
		if err != nil {
			return err
		}

		notes := "Welcome to Wise, here’s $100 on us!"
		ti = &business.TransferInitiate{
			CreatedUserID:   shared.UserID(uID),
			BusinessID:      shared.BusinessID(bID),
			SourceAccountId: saID.String(),
			SourceType:      banksrv.TransferTypeAccount,
			DestAccountId:   acc.Id,
			DestType:        banksrv.TransferTypeAccount,
			Amount:          100,
			Currency:        banksrv.CurrencyUSD,
			SendEmail:       true,
			Notes:           &notes,
		}
	} else {
		// destination account id
		var destAccID string
		err = service.rdb.Get(&destAccID, "SELECT id FROM business_bank_account WHERE business_id = $1", destBID)
		if err != nil {
			log.Printf("Error getting bank_account_ id %v", err)
			return err
		}
		// Get business primary account (or by id)
		acc, err := business.NewAccountService().GetByIDInternal(destAccID)
		if err != nil {
			log.Printf("Error getting bank account by id %v", err)
			return err
		}
		// Get Linked account
		la, err := business.NewLinkedAccountService(services.NewSourceRequest()).GetByAccountNumberInternal(shared.BusinessID(bID), business.AccountNumber(acc.AccountNumber), acc.RoutingNumber)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error getting linked account %v", err)
			return err
		}

		if err == sql.ErrNoRows {
			// Create linked account
			lac := business.MerchantLinkedAccountCreate{
				UserID:            shared.UserID(uID),
				BusinessID:        shared.BusinessID(bID),
				AccountHolderName: bus.Name(),
				AccountNumber:     business.AccountNumber(acc.AccountNumber),
				AccountType:       banksrv.AccountTypeChecking,
				Currency:          banksrv.CurrencyUSD,
				Permission:        banksrv.LinkedAccountPermissionSendAndRecieve,
				RoutingNumber:     acc.RoutingNumber,
			}
			la, err = business.NewLinkedAccountService(services.NewSourceRequest()).LinkMerchantBankAccount(&lac)
			if err != nil {
				log.Printf("error link merchant account %v", err)
				return err
			}
		}

		notes := "Welcome to Wise, here’s $100 on us!"
		ti = &business.TransferInitiate{
			CreatedUserID:   shared.UserID(uID),
			BusinessID:      shared.BusinessID(bID),
			SourceAccountId: srcID,
			SourceType:      banksrv.TransferTypeAccount,
			DestAccountId:   la.Id,
			DestType:        banksrv.TransferTypeAccount,
			Amount:          100,
			Currency:        banksrv.CurrencyUSD,
			SendEmail:       true,
			Notes:           &notes,
		}
	}

	sr := services.NewSourceRequest()
	sr.UserID = shared.UserID(uID)

	mt, err := business.NewMoneyTransferService(sr).Transfer(ti)
	if err != nil {
		log.Printf("error making transfer %v", err)
		return err
	}
	fundedAt := time.Now()
	_, err = NewCSPService().CSPBusinessUpdateByBusinessID(destBID, CSPBusinessUpdate{
		PromoMoneyTransferID: &mt.BankTransferId,
		PromoFunded:          &fundedAt,
	})
	if err != nil {
		log.Printf("Error updating csp business with funded promotion %v", err)
	}
	return err
}

func (service banking) canPromoFund(id shared.BusinessID) error {
	b, err := NewCSPService().ByBusinessID(id)
	if err != nil {
		return err
	}
	if b.PromoFunded != nil && b.PromoMoneyTransferID != nil {
		return fmt.Errorf("Business was already promo funded at %v", b.PromoFunded)
	}
	return nil
}

func (service banking) ReissueCard(c business.CardReissueRequest) (*business.BankCardPartial, error) {
	switch c.Reason {
	case banksrv.CardReissueReasonNotReceived:
		break
	case banksrv.CardReissueReasonResendPin:
		break
	default:
		return nil, errors.New("invalid value for reason")
	}

	card, err := service.CardByID(c.BusinessID, c.CardID)
	if err != nil {
		return nil, err
	}

	var b struct {
		ConsumerID     shared.ConsumerID `db:"consumer_id"`
		Phone          string            `db:"phone"`
		MailingAddress *services.Address `db:"mailing_address"`
	}

	err = service.rdb.Get(
		&b, `
		SELECT
			wise_user.consumer_id,
			wise_user.phone,
			consumer.mailing_address
		FROM wise_user
		JOIN consumer ON consumer.id = wise_user.consumer_id
		WHERE wise_user.id = $1`,
		card.CardholderID,
	)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if b.MailingAddress == nil {
		return nil, errors.New("Mailing address can not be empty")
	}

	request := partnerbank.ReissueCardRequest{
		CardID:    partnerbank.CardBankID(card.BankCardId),
		Packaging: partnerbank.CardPackagingRegular,
		Delivery:  partnerbank.CardDeliveryStandard,
		Phone:     partnerbank.PhoneE164(b.Phone),
		Address:   partnerbank.AddressRequestTypeMailing,
		Reason:    partnerbank.CardReissueReason(c.Reason),
	}

	// Create card with partner bank
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.CardService(
		service.sourceReq.PartnerBankRequest(),
		partnerbank.BusinessID(c.BusinessID),
		partnerbank.ConsumerID(b.ConsumerID),
	)
	if err != nil {
		return nil, err
	}

	_, err = srv.Reissue(request)
	if err != nil {
		return nil, err
	}

	_, err = service.CreateReissueHistory(c)
	if err != nil {
		return nil, err
	}

	history, err := service.ReissueHistoryByCardID(card.Id)
	if err != nil {
		return nil, err
	}

	card.CardReissueHistory = history

	return &card, nil
}

func (service banking) CreateReissueHistory(c business.CardReissueRequest) (*business.CardReissueHistory, error) {
	columns := []string{
		"business_id", "card_id", "reason",
	}

	values := []string{
		":business_id", ":card_id", ":reason",
	}

	sql := fmt.Sprintf("INSERT INTO business_card_reissue_history(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := service.wdb.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	history := &business.CardReissueHistory{}

	err = stmt.Get(history, &c)
	if err != nil {
		return nil, fmt.Errorf("Error saving card err:%v", err.Error())
	}

	return history, nil
}

func (service banking) ReissueHistoryByCardID(cardID string) ([]business.CardReissueHistory, error) {
	var list []business.CardReissueHistory

	err := service.rdb.Select(&list, "SELECT * FROM business_card_reissue_history WHERE card_id = $1", cardID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return list, nil
}

func (service banking) DeclineTransferInReview(txnID, cognitoID, ipAddress string) error {
	cspUserID, err := cspuser.NewUserService(cspServices.NewSourceRequest()).ByCognitoID(cognitoID)
	if err != nil {
		return err
	}

	bts, err := acct.NewBankingTransferService()
	if err != nil {
		return err
	}

	return bts.Decline(txnID, cspUserID, ipAddress)
}

func (service banking) ApproveTransferInReview(txnID, cognitoID, ipAddress string) error {
	cspUserID, err := cspuser.NewUserService(cspServices.NewSourceRequest()).ByCognitoID(cognitoID)
	if err != nil {
		return err
	}

	bts, err := acct.NewBankingTransferService()
	if err != nil {
		return err
	}

	return bts.Approve(txnID, cspUserID, ipAddress)
}

func (service banking) GetTransferForTransactionID(txnID string) (*grpcBankTransfer.Transfer, error) {
	t := new(grpcBankTransfer.Transfer)

	bts, err := acct.NewBankingTransferService()
	if err != nil {
		return t, err
	}

	return bts.GetProtoByTransactionID(txnID)
}
