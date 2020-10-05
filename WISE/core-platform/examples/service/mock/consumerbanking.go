package mock

import (
	"time"

	"github.com/google/uuid"
	bankingsrv "github.com/wiseco/core-platform/services/banking"
)

// Create a mock bank account access object
func NewConsumerBankAccountAccess(ownerId string) bankingsrv.ConsumerBankAccountAccess {
	return bankingsrv.ConsumerBankAccountAccess{
		Id:            uuid.New().String(),
		BankAccountId: uuid.New().String(),
		UserId:        ownerId,
		AccessType:    bankingsrv.BankAccountAccessAdmin,
		AccessRole:    bankingsrv.BankAccountAccessRoleOwner,
	}
}

// Create mock consumer bank account
func NewConsumerBankAccount(accountHolderId, accountId string) bankingsrv.ConsumerBankAccount {
	var now = time.Now()
	var routingNumber = "234234226"
	var alias = "Checking Account"

	return bankingsrv.ConsumerBankAccount{
		BankAccount: bankingsrv.BankAccount{
			Id:               accountId,
			BankName:         bankingsrv.BankNameBBVA,
			AccountHolderId:  accountHolderId,
			AccountType:      bankingsrv.BankAccountTypeChecking,
			AccountStatus:    bankingsrv.BankAccountStatusActive,
			AccountNumber:    "7238723451",
			RoutingNumber:    routingNumber,
			Alias:            &alias,
			AvailableBalance: 500.00,
			PostedBalance:    500.00,
			Currency:         bankingsrv.CurrencyUSD,
			Opened:           now,
			Created:          now,
			Updated:          now,
		},
	}
}

// Create a new mock consumer card
func NewConsumerDebitCard(cardholderId, accountId, cardId string) bankingsrv.ConsumerBankCard {
	var now = time.Now()
	var alias = "My Visa Card"

	return bankingsrv.ConsumerBankCard{
		BankCard: bankingsrv.BankCard{
			Id:               cardId,
			BankName:         bankingsrv.BankNameBBVA,
			CardholderId:     cardholderId,
			CardType:         bankingsrv.CardTypeDebit,
			CardStatus:       bankingsrv.CardStatusActive,
			Alias:            &alias,
			CardholderName:   "Mark Jones",
			MaskedCardNumber: "111122******4444",
			Currency:         bankingsrv.CurrencyUSD,
			Created:          now,
			Updated:          now,
		},
	}
}

// Returns a mock transfer object
func NewConsumerMoneyTransfer(createdUserId string, transferId string) bankingsrv.ConsumerMoneyTransfer {
	var now = time.Now()

	return bankingsrv.ConsumerMoneyTransfer{
		MoneyTransfer: bankingsrv.MoneyTransfer{
			Id:              transferId,
			CreatedUserId:   uuid.New().String(),
			BankName:        bankingsrv.BankNameBBVA,
			SourceAccountId: uuid.New().String(),
			//SourceAccountType: bankingsrv.MoneyTransferAccountTypeChecking,
			DestAccountId: uuid.New().String(),
			//DestAccountType:   bankingsrv.MoneyTransferAccountTypeChecking,
			Amount:   125.55,
			Currency: bankingsrv.CurrencyUSD,
			//RailType:          bankingsrv.MoneyTransferRailTypeWise,
			Status: bankingsrv.MoneyTransferStatusPosted,
			Transactions: &[]bankingsrv.Transaction{
				NewConsumerTransaction(transferId, createdUserId),
			},
			Created: now,
			Updated: now,
		},
	}
}

func NewConsumerCardTransaction(amount float64) bankingsrv.CardTransaction {

	var now = time.Now()
	var respCode = "00"
	var authNum = "281276"
	var transCode = "002000"
	var currency = bankingsrv.CurrencyUSD
	var entryMode = "050"
	var condCode = "01000008045"
	var abin = "476501"
	var mid = "000266015396887"
	var mcc = "5812"
	var mterm = "08673556"
	var mname = "OSHA EXPRESS"
	var maddr = "ONE MARKET PLAZA #21B"
	var mcity = "SAN FRANCISCO"
	var mstate = "CA"
	var mcountry = "US"

	return bankingsrv.CardTransaction{
		AuthAmount:            &amount,
		AuthDate:              &now,
		AuthResponseCode:      &respCode,
		AuthNumber:            &authNum,
		TransactionCode:       &transCode,
		LocalAmount:           &amount,
		LocalCurrency:         &currency,
		LocalDate:             &now,
		PosEntryMode:          &entryMode,
		PosConditionCode:      &condCode,
		AcquirerBankIdNumber:  &abin,
		MerchantId:            &mid,
		MerchantCategoryCode:  &mcc,
		MerchantTerminal:      &mterm,
		MerchantName:          &mname,
		MerchantStreetAddress: &maddr,
		MerchantCity:          &mcity,
		MerchantState:         &mstate,
		MerchantCountry:       &mcountry,
	}
}

func NewConsumerTransaction(transferId string, ownerId string) bankingsrv.Transaction {

	var now = time.Now()
	var transactionId = uuid.New().String()
	var sourceAccountId = uuid.New().String()
	var sourceAccountType = bankingsrv.TransactionSourceTypeCard
	var codeType = bankingsrv.TransactionNetworkTypeVisa
	var networkId = "00000472486"
	var cardTrans = NewConsumerCardTransaction(155.55)
	var transferDesc = "Visa Card Purchase"

	return bankingsrv.Transaction{
		Id:                         transactionId,
		ConsumerId:                 &ownerId,
		BusinessId:                 nil,
		BankName:                   bankingsrv.BankNameBBVA,
		SourceAccountId:            &sourceAccountId,
		SourceAccountType:          &sourceAccountType,
		DestAccountId:              nil,
		DestAccountType:            nil,
		TransactionType:            bankingsrv.TransactionTypePurchase,
		CodeType:                   bankingsrv.TransactionCodeTypeDebitPosted,
		NetworkType:                &codeType,
		NetworkId:                  &networkId,
		Amount:                     155.55,
		Currency:                   bankingsrv.CurrencyUSD,
		CardTransactionDetails:     &cardTrans,
		CardHoldTransactionDetails: nil,
		MoneyTransferId:            &transferId,
		MoneyTransferDesc:          &transferDesc,
		TransactionCreatedDate:     now,
		TransactionUpdatedDate:     now,
		Created:                    now,
	}
}
