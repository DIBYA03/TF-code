package main

import (
	"fmt"
	"os"

	"github.com/wiseco/core-platform/services/banking/business"
	bussrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/shared"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
)

var bus = []string{}

func main() {

	// Take money out of account and deposit in wise refund account
	wiseBusID := shared.BusinessID(os.Getenv("WISE_CLEARING_BUSINESS_ID"))
	if wiseBusID == "" {
		panic("WISE_CLEARING_BUSINESS_ID missing")
	}

	wiseUserID := shared.UserID(os.Getenv("WISE_CLEARING_USER_ID"))
	if wiseUserID == "" {
		panic("WISE_CLEARING_USER_ID missing")
	}

	wiseRefAccID := os.Getenv("WISE_REFUND_CLEARING_LINKED_ACCOUNT_ID")
	if wiseRefAccID == "" {
		panic("WISE_REFUND_CLEARING_LINKED_ACCOUNT_ID missing")
	}

	for _, id := range bus {
		busID, err := shared.ParseBusinessID(id)
		if err != nil {
			fmt.Println("ParseBusinessID error:", err, id)
			continue
		}

		fmt.Println("Business ID:", id)

		_, err = bussrv.NewBusinessServiceWithout().GetByIdInternal(busID)
		if err != nil {
			fmt.Println("Business error:", err, id)
		}

		accounts, err := business.NewAccountService().ListInternalByBusiness(busID, 20, 0)
		if err != nil {
			fmt.Println("List Error:", err, id)
			continue
		}

		for _, acc := range accounts {
			if acc.UsageType != business.UsageTypePrimary {
				fmt.Println("Only primary accounts should be closed", acc.Id, id, acc.AccountStatus)
				continue
			}

			// Cancel all related cards
			cards, err := business.NewCardServiceWithout().GetByAccountInternal(0, 20, acc.Id, busID)
			if err != nil {
				fmt.Println("GetByAccountInternal:", err, acc.Id)
				continue
			}

			for _, card := range cards {
				cardResp, err := business.NewCardServiceWithout().CancelCardInternal(card.Id)
				if err != nil {
					fmt.Println("CancelCardInternal:", err, acc.Id, id, card.CardStatus)
					continue
				}

				fmt.Println("Card canceled:", cardResp.CardStatus, cardResp.Id)
			}

			if acc.AvailableBalance > 0 {
				/* Pull money from account
				linkedAccount, err := business.NewLinkedAccountServiceWithout().GetByAccountNumberInternal(wiseBusID, business.AccountNumber(acc.AccountNumber), acc.RoutingNumber)
				if err != nil && err == sql.ErrNoRows {
					mc := &business.MerchantLinkedAccountCreate{
						UserID:            wiseUserID,
						BusinessID:        wiseBusID,
						AccountHolderName: bus.Name(),
						AccountNumber:     business.AccountNumber(acc.AccountNumber),
						RoutingNumber:     acc.RoutingNumber,
						AccountType:       banking.AccountType(acc.AccountType),
						Currency:          acc.Currency,
						Permission:        banking.LinkedAccountPermissionSendAndRecieve,
					}

					linkedAccount, err = business.NewLinkedAccountServiceWithout().LinkMerchantBankAccount(mc)
					if err != nil {
						fmt.Println("CancelCardInternal:", err, acc.Id, id)
						continue
					}
				}

				notes := "Account closure balance refund"

				// Initiate Transfer
				ti := &business.TransferInitiate{
					CreatedUserID:   wiseUserID,
					BusinessID:      wiseBusID,
					SourceAccountId: linkedAccount.Id,
					DestAccountId:   wiseRefAccID,
					Amount:          acc.AvailableBalance,
					SourceType:      banking.TransferTypeAccount,
					DestType:        banking.TransferTypeAccount,
					Currency:        banking.Currency(acc.Currency),
					Notes:           &notes,
				}

				srcReq := services.NewSourceRequest()
				srcReq.UserID = wiseUserID
				_, err = business.NewMoneyTransferService(srcReq).Transfer(ti)
				if err != nil {
					fmt.Println("Transfer:", err, acc.Id, id)
					continue
				} */

				fmt.Println("Account has a balance:", acc.AvailableBalance, acc.Id, id)
				continue
			}

			account, err := business.NewAccountService().DeactivateAccount(acc.Id, busID, grpcBanking.AccountStatusReason_ASR_CUSTOMER_REQUEST)
			if err != nil {
				fmt.Println("DeactivateAccount error:", err, acc.Id, id, acc.AccountStatus)
				continue
			}

			/* Send check to business
			if acc.AvailableBalance > 0 {
				// TODO: Send check without contact
			} */

			// Send a check to business in the available balance amount

			fmt.Println("Account closed:", account.AccountStatus, account.Id)
		}
	}
}
