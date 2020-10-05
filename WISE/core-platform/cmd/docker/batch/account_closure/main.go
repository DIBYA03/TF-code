package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	bussrv "github.com/wiseco/core-platform/services/business"
	accountclosure "github.com/wiseco/core-platform/services/csp/account_closure"
	csp "github.com/wiseco/core-platform/services/csp/services"
	"github.com/wiseco/core-platform/shared"
	idlib "github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
)

var bus = []string{}
var wiseBusID = shared.BusinessID("")
var wiseUserID = shared.UserID("")
var wiseRefAccID = ""

func main() {

	// Take money out of account and deposit in wise refund account
	wiseBusID = shared.BusinessID(os.Getenv("WISE_CLEARING_BUSINESS_ID"))
	if wiseBusID == "" {
		panic("WISE_CLEARING_BUSINESS_ID missing")
	}

	wiseUserID = shared.UserID(os.Getenv("WISE_CLEARING_USER_ID"))
	if wiseUserID == "" {
		panic("WISE_CLEARING_USER_ID missing")
	}

	wiseRefAccID = os.Getenv("WISE_REFUND_CLEARING_LINKED_ACCOUNT_ID")
	if wiseRefAccID == "" {
		panic("WISE_REFUND_CLEARING_LINKED_ACCOUNT_ID missing")
	}

	crs, err := accountclosure.NewCSPService(csp.NewSourceRequest()).CSPClosureApprovedAndRetryRequestList()

	if err != nil {
		fmt.Println("Error in fetching requests, ", err)
		return
	}

	fmt.Println("Total closure requests to be addressed: ", len(crs))
	for _, cr := range crs {
		id := string(cr.BusinessID)
		busID, err := shared.ParseBusinessID(id)
		jumpToNextBusiness := false
		if err != nil {
			fmt.Println("ParseBusinessID error:", err, id)
			continue
		}

		crBus, err := bussrv.NewBusinessServiceWithout().GetByIdInternal(busID)

		if err != nil {
			fmt.Println("Business error:", err, id)
		}

		accounts, err := business.NewAccountService().ListInternalByBusiness(busID, 20, 0)
		if err != nil {
			fmt.Println("List Error:", err, id)
			continue
		}

		totalBalance := num.NewZero()
		if cr.RefundAmount != nil {
			totalBalance = totalBalance.Add(*cr.RefundAmount)
		}

		for _, acc := range accounts {

			if acc.UsageType != business.UsageTypePrimary {
				fmt.Println("Only primary accounts should be closed", acc.Id, id, acc.AccountStatus)
				continue
			}

			//------------------------------------------------------------------------
			//					CANCEL ALL CARDS OF ACCOUNT
			//------------------------------------------------------------------------
			err := cancelCards(cr, acc.Id, busID)
			if err != nil {
				closureRequestStatusUpdate(cr.ID, accountclosure.AccountClosureFailed)
				jumpToNextBusiness = true
				break
			}

			//------------------------------------------------------------------------
			//					PULL BALANCE
			//------------------------------------------------------------------------

			bas, err := business.NewBankingAccountService()
			if err != nil {
				addNewClosureState(accountclosure.ACStatePullBalanceFailed, cr.ID, acc.Id, fmt.Sprintf("Unable to get banking acc service to fetch avialable balance"))
				closureRequestStatusUpdate(cr.ID, accountclosure.AccountClosureFailed)
				jumpToNextBusiness = true
				break
			}

			ba, err := bas.GetBalanceByID(acc.Id, false)
			if err != nil {
				addNewClosureState(accountclosure.ACStatePullBalanceFailed, cr.ID, acc.Id, fmt.Sprintf("Unable to get available balance from bbva for business: %v", cr.BusinessID))
				closureRequestStatusUpdate(cr.ID, accountclosure.AccountClosureFailed)
				jumpToNextBusiness = true
			}
			acc.AvailableBalance = ba.AvailableBalance

			if acc.AvailableBalance > 0 {
				bal, err := pullBalance(cr, acc, crBus)
				if err != nil {
					closureRequestStatusUpdate(cr.ID, accountclosure.AccountClosureFailed)
					jumpToNextBusiness = true
					break
				}
				pulledAmt, err := num.NewFromFloat(bal)
				if err != nil {
					fmt.Printf("Error parsing refund amount %v, %v", bal, err)
				} else {
					totalBalance = totalBalance.Add(pulledAmt)
				}

			}

			accountclosure.NewCSPService(csp.NewSourceRequest()).CSPClosureRequestProcessed(cr.ID, totalBalance, accountclosure.AccountClosureStatus(*cr.Status))

			//------------------------------------------------------------------------
			//					DEACTIVATE ACCOUNT
			//------------------------------------------------------------------------
			_, err = deactivateAccount(cr, acc, crBus)
			if err != nil {
				closureRequestStatusUpdate(cr.ID, accountclosure.AccountClosureFailed)
				jumpToNextBusiness = true
				break
			}

			/* Send check to business
			if acc.AvailableBalance > 0 {
				// TODO: Send check without contact
			} */
		}
		if jumpToNextBusiness {
			continue
		}
		accountclosure.NewCSPService(csp.NewSourceRequest()).CSPClosureRequestProcessed(cr.ID, totalBalance, accountclosure.AccountClosureStatus(*cr.Status))

		err = deactivateBusiness(cr, crBus)
		if err != nil {
			closureRequestStatusUpdate(cr.ID, accountclosure.AccountClosureFailed)
			continue
		}

		status := accountclosure.AccountClosureRequestClosed
		if totalBalance.IsPositive() {
			status = accountclosure.AccountClosureRefundPending
		}

		accountclosure.NewCSPService(csp.NewSourceRequest()).CSPClosureRequestProcessed(cr.ID, totalBalance, status)

	}

}

func cancelCards(cr accountclosure.CSPClosureRequestItem, accountID string, businessID shared.BusinessID) error {
	accID, err := idlib.ParseBankAccountID(accountID)
	cards, err := business.NewCardServiceWithout().GetByAccountInternal(0, 20, accID.UUIDString(), businessID)
	if err != nil {
		fmt.Printf("Card fetch error: %v\n", err)
		addNewClosureState(accountclosure.ACStateCancelCardFailed, cr.ID, "", fmt.Sprintf("Card fetch error: %v", err))
		return fmt.Errorf("Card fetch error: %v", err)
	}

	for _, card := range cards {
		if card.CardStatus == banking.CardStatusCanceled {
			continue
		}
		addNewClosureState(accountclosure.ACStateCancelCardStarted, cr.ID, card.Id, "")
		cardResp, err := business.NewCardServiceWithout().CancelCardInternal(card.Id)
		if err != nil {
			fmt.Println("CancelCardInternal:", err, accountID, businessID, card.CardStatus)
			addNewClosureState(accountclosure.ACStateCancelCardFailed, cr.ID, card.Id, fmt.Sprintf("%v", err))
			return fmt.Errorf("Card Cancel error: %v", err)
		}

		fmt.Println("Card canceled:", cardResp.CardStatus, cardResp.Id)
		addNewClosureState(accountclosure.ACStateCancelCardSuccess, cr.ID, card.Id, fmt.Sprintf("Card status: %v", cardResp.CardStatus))
	}
	return nil
}

func pullBalance(cr accountclosure.CSPClosureRequestItem, acc *business.BankAccount, crBus *bussrv.Business) (float64, error) {
	addNewClosureState(accountclosure.ACStatePullBalanceStarted, cr.ID, acc.Id, fmt.Sprintf("Available balance: %v", acc.AvailableBalance))

	// Pull money from account
	srcReq := services.NewSourceRequest()
	srcReq.UserID = wiseUserID

	sourceLinkedAccount, err := business.NewLinkedAccountService(srcReq).GetByAccountNumberInternal(wiseBusID, business.AccountNumber(acc.AccountNumber), acc.RoutingNumber)

	if err != nil && err == sql.ErrNoRows {
		mc := &business.MerchantLinkedAccountCreate{
			UserID:            wiseUserID,
			BusinessID:        wiseBusID,
			AccountHolderName: crBus.Name(),
			AccountNumber:     business.AccountNumber(acc.AccountNumber),
			RoutingNumber:     acc.RoutingNumber,
			AccountType:       banking.AccountType(acc.AccountType),
			Currency:          acc.Currency,
			Permission:        banking.LinkedAccountPermissionSendAndRecieve,
		}

		sourceLinkedAccount, err = business.NewLinkedAccountService(srcReq).LinkMerchantBankAccount(mc)
		if err != nil {
			addNewClosureState(accountclosure.ACStatePullBalanceFailed, cr.ID, acc.Id, fmt.Sprintf("Link acc failed: %v", err))
			return 0.0, fmt.Errorf("CancelCardInternal: %v %v %v", err, acc.Id, crBus.ID)
		}
	} else if err != nil {
		addNewClosureState(accountclosure.ACStatePullBalanceFailed, cr.ID, acc.Id, fmt.Sprintf("Linked acc fetch failed: %v", err))
		return 0.0, err
	}

	notes := "Account closure balance refund"

	// Initiate Transfer
	ti := &business.TransferInitiate{
		CreatedUserID:   wiseUserID,
		BusinessID:      wiseBusID,
		SourceAccountId: sourceLinkedAccount.Id,
		DestAccountId:   wiseRefAccID,
		Amount:          acc.AvailableBalance,
		SourceType:      banking.TransferTypeAccount,
		DestType:        banking.TransferTypeAccount,
		Currency:        banking.Currency(acc.Currency),
		Notes:           &notes,
	}

	_, err = business.NewMoneyTransferService(srcReq).Transfer(ti)
	if err != nil {
		fmt.Println("Transfer:", err, acc.Id, crBus.ID)
		errStr := fmt.Sprintf("[Transfer failed] SrcLinkedAcc: %v, DestLinkedAcc: %v, err: %v", ti.SourceAccountId, ti.DestAccountId, err)
		addNewClosureState(accountclosure.ACStatePullBalanceFailed, cr.ID, acc.Id, errStr)
		return 0.0, fmt.Errorf("Transfer: %v %v %v", err, acc.Id, crBus.ID)
	}

	addNewClosureState(accountclosure.ACStatePullBalanceSuccess, cr.ID, acc.Id, "")
	return acc.AvailableBalance, nil
}

func deactivateAccount(cr accountclosure.CSPClosureRequestItem, acc *business.BankAccount, crBus *bussrv.Business) (*business.BankAccount, error) {
	addNewClosureState(accountclosure.ACStateDeactivateAccountStarted, cr.ID, acc.Id, "")
	accID, err := idlib.ParseBankAccountID(acc.Id)
	respAcc, err := business.NewAccountService().DeactivateAccount(accID.UUIDString(), crBus.ID, grpcBanking.AccountStatusReason_ASR_CUSTOMER_REQUEST)
	if err != nil {
		fmt.Println("DeactivateAccount error:", err, acc.Id, cr.ID, acc.AccountStatus)
		addNewClosureState(accountclosure.ACStateDeactivateAccountFailed, cr.ID, acc.Id, fmt.Sprintf("%v", err))
		return nil, err
	}

	addNewClosureState(accountclosure.ACStateDeactivateAccountSuccess, cr.ID, acc.Id, "")
	return respAcc, nil
}

func deactivateBusiness(cr accountclosure.CSPClosureRequestItem, crBus *bussrv.Business) error {
	addNewClosureState(accountclosure.ACStateDeactivateBusinessStarted, cr.ID, string(crBus.ID), "")

	srcReq := services.NewSourceRequest()
	srcReq.UserID = crBus.OwnerID

	err := bussrv.NewBusinessService(srcReq).Deactivate(crBus.ID)
	if err != nil {
		addNewClosureState(accountclosure.ACStateDeactivateBusinessFailed, cr.ID, string(crBus.ID), fmt.Sprintf("%v", err))
	}
	addNewClosureState(accountclosure.ACStateDeactivateBusinessSuccess, cr.ID, string(crBus.ID), "")

	return err
}

// HELPERS
func closureRequestStatusUpdate(id string, status accountclosure.AccountClosureStatus) {
	accountclosure.NewCSPService(csp.NewSourceRequest()).CSPClosureRequestUpdate(id, string(status))

}

func addNewClosureState(state accountclosure.AccountClosureState, crID string, id string, description string) {
	fmt.Printf("State: %v, Desc: %v\n", state, description)
	closureState := accountclosure.CSPClosureStatePostBody{State: state, ClosureRequestID: crID, ItemID: &id, Description: &description}
	accountclosure.NewCSPService(csp.NewSourceRequest()).CSPClosureStateAddNew(closureState)
}
