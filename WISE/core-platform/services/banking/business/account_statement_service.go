package business

import (
	"log"

	"github.com/jmoiron/sqlx"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	banking "github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type AccountStmtService interface {
	// Read
	List(accountID string, accountHolderID shared.UserID, businessID shared.BusinessID) (*[]BankAccountStatement, error)
	GetByID(statementID string, accountID string, accountHolderID shared.UserID, businessID shared.BusinessID) (*banking.BankAccountStatementDocument, error)
}

type accountStmtDataStore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

func NewAccountStmtService(r services.SourceRequest) AccountStmtService {
	return &accountStmtDataStore{r, data.DBWrite}
}

func (a accountStmtDataStore) List(accountID string, accountHolderID shared.UserID, businessID shared.BusinessID) (*[]BankAccountStatement, error) {
	// Get bank account ID
	account, err := NewBankAccountService(a.sourceReq).GetByID(accountID, businessID)
	if err != nil {
		log.Println(err)
		return nil, err

	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.BankAccountService(a.sourceReq.CSPPartnerBankRequest(accountHolderID), partnerbank.BusinessID(businessID))
	if err != nil {
		return nil, err
	}

	s, err := srv.GetStatements(partnerbank.AccountBankID(account.BankAccountId))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	statements := transformStatementResponse(businessID, accountHolderID, s)

	return statements, nil
}

func transformStatementResponse(businessID shared.BusinessID, accountHolderID shared.UserID, s []partnerbank.AccountStatementResponse) *[]BankAccountStatement {

	accountStatements := []BankAccountStatement{}

	for _, stmt := range s {
		s := BankAccountStatement{
			BusinessID: businessID,
		}

		s.AccountHolderID = accountHolderID
		s.Description = stmt.Description
		s.PageCount = stmt.PageCount
		s.Created = stmt.Created
		s.StatementID = string(stmt.StatementID)

		accountStatements = append(accountStatements, s)

	}

	return &accountStatements

}

func (a accountStmtDataStore) GetByID(statementID string, accountID string, accountHolderID shared.UserID, businessID shared.BusinessID) (*banking.BankAccountStatementDocument, error) {
	// Get bank account ID
	account, err := NewBankAccountService(a.sourceReq).GetByID(accountID, businessID)
	if err != nil {
		log.Println(err)
		return nil, err

	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.BankAccountService(a.sourceReq.CSPPartnerBankRequest(accountHolderID), partnerbank.BusinessID(businessID))
	if err != nil {
		return nil, err
	}

	s, err := srv.GetStatementByID(partnerbank.AccountBankID(account.BankAccountId), partnerbank.AccountStatementBankID(statementID))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	statement := banking.BankAccountStatementDocument{
		Content:     s.Content,
		ContentType: s.ContentType,
	}

	return &statement, nil
}
