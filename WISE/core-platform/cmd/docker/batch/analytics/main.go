package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wiseco/core-platform/services/invoice"
	"github.com/wiseco/core-platform/shared"
	"log"
	"os"

	"github.com/wiseco/core-platform/analytics"
	"github.com/wiseco/core-platform/services/banking/business"
	busBanking "github.com/wiseco/core-platform/services/banking/business"
	core "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcTxn "github.com/wiseco/protobuf/golang/transaction"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
)

var txnClient grpcBankTxn.TransactionServiceClient

func main() {
	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameTransaction)
		if err != nil {
			panic(err)
		}
		cl, err := grpc.NewInsecureClient(sn)
		if err != nil {
			panic(err)
		}

		defer cl.CloseAndCancel()
		txnClient = grpcBankTxn.NewTransactionServiceClient(cl.GetConn())
	}

	businesses, err := getAllBusinesses()
	if err != nil {
		panic(err)
	}
	if os.Getenv("USE_INVOICE_SERVICE") == "true" {
		err = updateInvoiceCount(businesses)
		if err != nil {
			panic(err)
		}
	}

	for _, b := range businesses {
		bus, err := getBusinessBankingDetails(&b)
		if err != nil {
			panic(err)
		}

		bus, err = getBusinessTransactionDetails(bus)
		if err != nil {
			panic(err)
		}

		err = sendSegmentTraits(bus)
		if err != nil {
			panic(err)
		}
	}

	message := fmt.Sprintf("Intercom batch update done for %d businesses!", len(businesses))
	log.Println(message)
}

func getAllBusinesses() ([]Business, error) {
	var businesses []Business

	// Get all core data for businesses - TODO: Use offset limit to minimize traffic
	err := core.DBRead.Select(
		&businesses,
		`SELECT
			business.id "business.id",
			industry_type,
			entity_type,
			legal_name, 
			dba,
			business.kyc_status "business.kyc_status", 
			business.mailing_address "business.mailing_address", 
			origin_date,
			business.owner_id "business.owner_id",
			title_type,
			date_of_birth,
			wise_user.consumer_id "wise_user.consumer_id",
			wise_user.phone_verified,
			consumer.first_name "consumer.first_name",
			consumer.middle_name "consumer.middle_name", 
			consumer.last_name "consumer.last_name",
			consumer.kyc_status "consumer.kyc_status", 
			cr.card_reader_count, 
			bc.contact_count,
			bmr.invoice_count
		FROM business

		LEFT JOIN (select business_id, count(card_reader.business_id) AS card_reader_count 
		FROM card_reader 
		WHERE deactivated IS NULL group by business_id) cr 
		ON cr.business_id = business.id

		LEFT JOIN (select business_id, count(business_contact.business_id) AS contact_count 
		FROM business_contact 
		WHERE deactivated IS NULL group by business_id) bc 
		ON bc.business_id = business.id

		LEFT JOIN (select business_id, count(business_money_request.business_id) AS invoice_count 
		FROM business_money_request group by business_id) bmr 
		ON bmr.business_id = business.id

		JOIN wise_user 
		ON business.owner_id = wise_user.id

		JOIN consumer
		ON wise_user.consumer_id = consumer.id

		JOIN business_member 
		ON consumer.id = business_member.consumer_id`)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return businesses, nil
}

// updateInvoiceCount updated invoice count as following
//
// from core_db use the invoice count for POS
// from invoice service get the invoice count
// then sum the above 2
func updateInvoiceCount(business []Business) error {

	var posBusDetails []Business
	err := core.DBRead.Select(
		&posBusDetails,
		`select business_id, count(business_money_request.business_id) AS invoice_count 
		FROM business_money_request where request_type='pos' group by business_id`)
	if err != nil {
		log.Println(err)
		return err
	}
	posBusInvoiceCount := make(map[shared.BusinessID]int)
	for _, business := range posBusDetails {
		posBusInvoiceCount[business.ID] = *business.InvoiceCount
	}

	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		log.Println("unable to instantiate invoice service", err)
		return err
	}
	counts, err := invSvc.GetInvoiceCounts()
	if err != nil {
		log.Println("unable to get invoice count", err)
		return err
	}
	for _, bus := range business {
		count, ok := counts[bus.ID]
		posCount, posOk := posBusInvoiceCount[bus.ID]
		c := int(count)
		if posOk && ok {
			totalCount := posCount + c
			bus.InvoiceCount = &totalCount
		} else if posOk {
			bus.InvoiceCount = &posCount
		} else if ok {
			bus.InvoiceCount = &c
		} else {
			bus.InvoiceCount = nil
		}
	}
	return nil
}

// TODO: Add code to get data from banking service
func getBusinessBankingDetails(b *Business) (*Business, error) {
	err := core.DBRead.Get(
		b, `
        SELECT
            card_status,
            daily_transaction_limit
        FROM business_bank_card
		WHERE business_bank_card.business_id = $1`,
		b.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Debit Card:", err, b.ID)
		return b, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := busBanking.NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		var account *busBanking.BankAccount
		accounts, err := bas.GetByBusinessID(b.ID, 10, 0)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		for _, acc := range accounts {
			if acc.UsageType == business.UsageTypePrimary {
				account = &acc
				break
			}
		}

		if account != nil {
			b.AccountID = account.Id
			b.AccountType = account.AccountType
			b.AccountStatus = account.AccountStatus
			b.AccountOpened = account.Opened
			b.AvailableBalance = account.AvailableBalance
			b.PostedBalance = account.PostedBalance
		}

		blas, err := busBanking.NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		stfs := []grpcBanking.LinkedSubtype{
			grpcBanking.LinkedSubtype_LST_PRIMARY,
			grpcBanking.LinkedSubtype_LST_EXTERNAL,
			grpcBanking.LinkedSubtype_LST_CONTACT,
			grpcBanking.LinkedSubtype_LST_CONTACT_INVISIBLE,
		}

		la, err := blas.List(b.ID, stfs, 100, 0)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if la != nil {
			count := len(la)
			b.LinkedAccountCount = &count
		}

		blcs, err := busBanking.NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		lc, err := blcs.List(b.ID, 100, 0)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if lc != nil {
			count := len(lc)
			b.LinkedCardCount = &count
		}

		return b, nil
	}

	err = core.DBRead.Get(
		b, `
		SELECT
            business_bank_account.id "business_bank_account.id",
			business_bank_account.account_type,
			account_status,
			opened,
			available_balance,
			posted_balance, 
            blba.linked_account_count,
			blc.linked_card_count
    	FROM business_bank_account

	    LEFT JOIN (
			SELECT business_id, count(business_linked_bank_account.business_id) AS linked_account_count
	    	FROM business_linked_bank_account
		    WHERE deactivated IS NULL group by business_id
		) blba
    	ON blba.business_id = business_bank_account.business_id
    
		LEFT JOIN (
			SELECT business_id, count(business_linked_card.business_id) AS linked_card_count
	    	FROM business_linked_card
		    WHERE deactivated IS NULL group by business_id
		) blc
    	ON blc.business_id = business_bank_account.business_id

		WHERE business_bank_account.business_id = $1`,
		b.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Account:", err, b.ID)
		return b, err
	}

	return b, nil
}

func getBusinessTransactionDetails(b *Business) (*Business, error) {
	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		busUUID, err := id.ParseUUID(string(b.ID))
		if err != nil {
			return b, err
		}

		busID := id.BusinessID(busUUID)
		resp, err := txnClient.GetStats(context.Background(), &grpcBankTxn.StatsRequest{BusinessId: busID.String()})
		if err != nil {
			return b, err
		}

		b.PostedTransactionCount = int(resp.Count)
		for _, ts := range resp.TypeStats {
			switch ts.Type {
			case grpcTxn.BankTransactionType_BTT_CARD_PUSH_DEBIT, grpcTxn.BankTransactionType_BTT_CARD_PUSH_CREDIT:
				b.PushToDebitTransactionCount += int(ts.Count)
			case grpcTxn.BankTransactionType_BTT_CARD_PAYMENT_READER_CREDIT:
				b.CardReaderTransactionCount += int(ts.Count)
			case grpcTxn.BankTransactionType_BTT_ACH_DEBIT:
				b.ACHDebitTransactionCount += int(ts.Count)
			case grpcTxn.BankTransactionType_BTT_ACH_CREDIT:
				b.ACHCreditTransactionCount += int(ts.Count)
			case grpcTxn.BankTransactionType_BTT_WIRE_CREDIT:
				b.WireCreditTransactionCount += int(ts.Count)
			case grpcTxn.BankTransactionType_BTT_CARD_ATM_DEBIT:
				b.DebitCardATMTransactionCount += int(ts.Count)
			case grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_DEBIT, grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_ONLINE_DEBIT:
				b.DebitCardTransactionCount += int(ts.Count)
			}
		}
	} else {
		err := transaction.DBRead.Get(
			b,
			`SELECT
				SUM(CASE WHEN (transaction_subtype = 'cardPushDebit' OR transaction_subtype = 'cardPushCredit')
				THEN 1 ELSE 0 end) AS push_to_debit_transaction_count,

				SUM(case WHEN transaction_subtype = 'cardReaderCredit'
				THEN 1 ELSE 0 end) AS card_reader_transaction_count,

				SUM(case WHEN transaction_subtype = 'achTransferDebit'
				THEN 1 ELSE 0 end) AS ach_debit_transaction_count,

				SUM(case WHEN transaction_subtype = 'achTransferCredit'
				THEN 1 ELSE 0 end) AS ach_credit_transaction_count,

				SUM(case WHEN transaction_subtype = 'wireTransferCredit'
				THEN 1 ELSE 0 end) AS wire_credit_transaction_count,

				SUM(case WHEN transaction_subtype = 'cardATMDebit'
				THEN 1 ELSE 0 end) AS debit_card_atm_transaction_count,

				SUM(case WHEN (transaction_subtype = 'cardPurchaseDebit' OR transaction_subtype = 'cardPurchaseDebitOnline')
				THEN 1 ELSE 0 end) AS debit_card_transaction_count,

				SUM(case WHEN business_id = $1
				THEN 1 ELSE 0 end) AS transaction_count

		    FROM business_transaction
		    WHERE business_id = $1 GROUP BY business_id`,
			b.ID)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			return b, err
		}
	}

	return b, nil
}

func sendSegmentTraits(b *Business) error {
	traits := make(map[string]interface{})

	traits[analytics.ConsumerID] = b.ConsumerID
	traits[analytics.ConsumerPhoneVerified] = b.PhoneVerified
	traits[analytics.ConsumerFirstName] = b.ConsumerFirstName
	traits[analytics.ConsumerLastName] = b.ConsumerLastName
	traits[analytics.ConsumerMiddleName] = b.ConsumerMiddleName
	traits[analytics.ConsumerDateOfBirth] = b.ConsumerDOB
	traits[analytics.ConsumerKYCStatus] = b.ConsumerKYCStatus
	traits[analytics.ConsumerJobTitle] = b.ConsumerJobTitle

	traits[analytics.BusinessId] = b.ID.ToPrefixString()
	traits[analytics.BusinessLegalName] = b.LegalName
	traits[analytics.BusinessDBA] = b.DBA
	traits[analytics.BusinessEntityType] = b.EntityType
	traits[analytics.BusinessIndustryType] = b.IndustryType
	traits[analytics.BusinessKYCStatus] = b.BusinessKYCStatus
	traits[analytics.BusinessMailingAddress] = b.MailingAddress
	traits[analytics.BusinessOriginDate] = b.BusinessOriginDate

	traits[analytics.BusinessCardId] = b.DebitCardID
	traits[analytics.BusinessCardStatus] = b.DebitCardStatus
	traits[analytics.BusinessTransactionCount] = b.DailyTransactionLimit

	traits[analytics.BusinessAccountId] = b.AccountID
	traits[analytics.BusinessAccountType] = b.AccountType
	traits[analytics.BusinessAccountStatus] = b.AccountStatus
	traits[analytics.BusinessAccountBankName] = "bbva"
	traits[analytics.BusinessAccountOpened] = b.AccountOpened
	traits[analytics.BusinessAccountAvailableBalance] = b.AvailableBalance

	traits[analytics.LinkedAccountCount] = b.LinkedAccountCount
	traits[analytics.LinkedCardCount] = b.LinkedCardCount
	traits[analytics.CardReaderCount] = b.CardReaderCount
	traits[analytics.ContactCount] = b.ContactCount
	traits[analytics.InvoiceCount] = b.InvoiceCount

	traits[analytics.PostedTransactionCount] = b.PostedTransactionCount
	traits[analytics.PushToDebitTransactionCount] = b.PushToDebitTransactionCount
	traits[analytics.DebitCardATMTransactionCount] = b.DebitCardATMTransactionCount
	traits[analytics.CardReaderTransactionCount] = b.CardReaderTransactionCount
	traits[analytics.ACHDebitTransactionCount] = b.ACHDebitTransactionCount
	traits[analytics.ACHCreditTransactionCount] = b.ACHCreditTransactionCount
	traits[analytics.WireCreditTransactionCount] = b.WireCreditTransactionCount
	traits[analytics.DebitCardTransactionCount] = b.DebitCardTransactionCount

	log.Println(b.ID)
	return analytics.Identify(b.UserID, traits)
}
