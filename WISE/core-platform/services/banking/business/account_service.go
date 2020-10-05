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
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/partner/bank"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	bankData "github.com/wiseco/core-platform/partner/bank/data"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
)

//BusinessBankingService  the interface of the store exposed to the clients
type BankAccountService interface {
	Create(*BankAccountCreate) (*BankAccount, error)
	GetByID(string, shared.BusinessID) (*BankAccount, error)
	GetByIDInternal(string) (*BankAccount, error)
	GetByBankAccountId(string, shared.BusinessID) (*BankAccount, error)
	GetByUserID(shared.UserID, shared.BusinessID) (*[]BankAccount, error)
	GetByUsageType(shared.UserID, shared.BusinessID, UsageType) ([]BankAccount, error)
	List(shared.BusinessID, int, int) (*[]BankAccount, error)
	ListInternalByBusiness(shared.BusinessID, int, int) ([]*BankAccount, error)
	ListInternal(int, int) ([]*BankAccount, error)
	ListByOwner(string) ([]BankAccount, error)
	Update(*BankAccountUpdate) (*BankAccount, error)
	GetByBankAccountIDUserID(bankAccountID string, userID shared.UserID) (*BankAccount, error)

	GetBalanceByIDInternal(string) (*BankAccountBalance, error)
	GetBalanceByID(string, shared.BusinessID) (*BankAccountBalance, error)
	UpdateBankAccountStatus(string, string) (*BankAccount, error)

	// Bank Account Blocks
	GetByBlockID(string, string) (*banking.AccountBlock, error)
	CreateAccountBlock(banking.AccountBlockCreate) (*banking.AccountBlock, error)
	DeactivateAccountBlock(string) error
	GetAllBankAccountBlocks(string, shared.BusinessID) ([]banking.BankAccountBlockResponse, error)

	// Deactivate Account
	DeactivateAccount(string, shared.BusinessID, grpcBanking.AccountStatusReason) (*BankAccount, error)
}

//businessBankingDataStore acts as the store for the business banking
// where `CRUD` operations are made
type bankAccountDataStore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

//NewBankAccountService returns a business Bank Account Service
func NewBankAccountService(r services.SourceRequest) BankAccountService {
	return &bankAccountDataStore{r, data.DBWrite}
}

//NewAccountService returns a new bank account service without a source request
func NewAccountService() BankAccountService {
	return &bankAccountDataStore{services.NewSourceRequest(), data.DBWrite}
}

func (store *bankAccountDataStore) DeactivateAccountBlock(ID string) error {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return err
		}

		return bas.DeactivateAccountBlock(ID)
	}

	block := "UPDATE business_bank_account_block SET deactivated = CURRENT_TIMESTAMP WHERE id = $1"

	_, err := store.Exec(block, ID)
	if err == nil {
		log.Println(err)
		return err
	}

	return nil
}

func (store *bankAccountDataStore) GetByBlockID(blockID string, accountID string) (*banking.AccountBlock, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.GetByBlockID(blockID)
	}

	b := &banking.AccountBlock{}
	err := store.Get(b, "SELECT * FROM business_bank_account_block WHERE block_id = $1 AND account_id = $2",
		blockID, accountID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return b, nil
}

func (store *bankAccountDataStore) CreateAccountBlock(c banking.AccountBlockCreate) (*banking.AccountBlock, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.CreateAccountBlock(c)
	}

	columns := []string{
		"account_id", "reason", "block_id", "block_type", "originated_from",
	}
	// Default/mandatory values
	values := []string{
		":account_id", ":reason", ":block_id", ":block_type", ":originated_from",
	}

	sql := fmt.Sprintf("INSERT INTO business_bank_account_block(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := store.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	block := &banking.AccountBlock{}

	err = stmt.Get(block, &c)
	if err != nil {
		return nil, err
	}

	return block, nil

}

func (store *bankAccountDataStore) GetAllBankAccountBlocks(accountID string, businessID shared.BusinessID) ([]banking.BankAccountBlockResponse, error) {
	var account *BankAccount
	var err error

	blocks := []banking.BankAccountBlockResponse{}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		bs, err := bas.GetManyAccountBlocks(accountID)
		if err != nil {
			return nil, err
		}

		for _, block := range bs {
			status := banking.AccountBlockStatusActive

			if block.Deactivated != nil {
				status = banking.AccountBlockStatusCanceled
			}

			b := banking.BankAccountBlockResponse{
				BlockID: banking.AccountBlockBankID(block.BlockID),
				Type:    block.BlockType,
				Status:  status,
			}

			blocks = append(blocks, b)
		}
	} else {
		account, err = store.GetByIDInternal(accountID)
		if err != nil {
			return nil, err
		}

		// Get account info from bank
		providerName, ok := banking.ToPartnerBankName[account.BankName]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Partner bank %s does not exist", account.BankName))
		}

		bank, err := partnerbank.GetBusinessBank(providerName)
		if err != nil {
			return nil, err
		}

		srv, err := bank.BankAccountService(store.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(businessID))
		if err != nil {
			return nil, err
		}

		bs, err := srv.GetAllBlocks(partnerbank.AccountBankID(account.BankAccountId))
		if err != nil {
			return nil, err
		}

		for _, block := range bs {
			b := banking.BankAccountBlockResponse{
				BlockID: banking.AccountBlockBankID(block.BlockID),
				Type:    banking.AccountBlockType(block.Type),
				Status:  banking.AccountBlockStatus(block.Status),
			}

			blocks = append(blocks, b)
		}
	}

	return blocks, nil

}

//Create will create a business bank account
func (store *bankAccountDataStore) Create(account *BankAccountCreate) (*BankAccount, error) {
	accountFull := BankAccountCreateFull{
		BankAccountCreate: *account,
	}

	var bus business.Business
	err := store.Get(&bus, "SELECT * FROM business WHERE id = $1", account.BusinessID)
	if err != nil {
		return nil, err
	}

	var usr user.User
	err = store.Get(&usr, "SELECT * FROM wise_user WHERE id = $1", bus.OwnerID)
	if err != nil {
		return nil, err
	}

	// Default to primary
	switch accountFull.UsageType {
	case UsageTypeNone:
		accountFull.UsageType = UsageTypePrimary
	case UsageTypePrimary, UsageTypeClearing:
		break
	default:
		return nil, fmt.Errorf("invalid usage type: %s", accountFull.UsageType)
	}

	if bus.EntityType == nil {
		return nil, errors.New("entity type is missing")
	} else if bus.OperationType == nil {
		return nil, errors.New("operation type is missing")
	}

	// Ensure single primary account for now
	if accountFull.UsageType == UsageTypePrimary {
		accounts, err := store.GetByUsageType(bus.OwnerID, account.BusinessID, accountFull.UsageType)
		if err == nil && len(accounts) > 0 {
			return nil, fmt.Errorf("account already exists for given usage type: %s", accountFull.UsageType)
		}
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		apiReq := bank.APIRequest{}
		bs := bankData.NewBusinessService(apiReq, bank.ProviderNameBBVA)

		bankBus, err := bs.GetByBusinessID(bank.BusinessID(bus.ID))
		if err != nil {
			return nil, err
		}

		cs := bankData.NewConsumerService(apiReq, bank.ProviderNameBBVA)

		bankCon, err := cs.GetByConsumerID(bank.ConsumerID(usr.ConsumerID))
		if err != nil {
			return nil, err
		}

		return bas.Create(account, bus, usr, bankBus, bankCon)
	}

	providerName, ok := banking.ToPartnerBankName[account.BankName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Partner bank %s does not exist", account.BankName))
	}

	bank, err := partnerbank.GetBusinessBank(providerName)
	if err != nil {
		return nil, err
	}

	srv, err := bank.BankAccountService(store.sourceReq.CSPPartnerBankRequest(bus.OwnerID), partnerbank.BusinessID(bus.ID))
	if err != nil {
		return nil, err
	}

	// Create non-consumer bank account
	resp, err := srv.Create(
		partnerbank.CreateBusinessBankAccountRequest{
			BusinessID: partnerbank.BusinessID(account.BusinessID),
			ExtraParticipants: []partnerbank.AccountParticipantRequest{
				partnerbank.AccountParticipantRequest{
					ConsumerID: partnerbank.ConsumerID(usr.ConsumerID), //FROM CSP since source id doesnt have the user id
					Role:       partnerbank.AccountRoleAuthorized,
				},
			},
			AccountType:  partnerbank.AccountType(account.AccountType),
			Alias:        account.Alias,
			BusinessType: partnerbank.BusinessEntity(*bus.EntityType),
			IsForeign:    *bus.OperationType != business.OperationTypeLocal,
		},
	)
	if err != nil {
		return nil, err
	}

	var wireRouting *string = nil
	if len(resp.WireRouting) > 0 {
		wireRouting = &resp.WireRouting
	}

	accountFull.BankAccountId = resp.AccountID.String()
	accountFull.BankExtra = resp.BankExtra
	accountFull.AccountHolderID = bus.OwnerID //FROM CSP since source id doesnt have the user id
	accountFull.AccountStatus = resp.Status.String()
	accountFull.AccountNumber = resp.AccountNumber
	accountFull.RoutingNumber = resp.RoutingNumber
	accountFull.WireRouting = wireRouting
	accountFull.AvailableBalance = resp.AvailableBalance
	accountFull.PostedBalance = resp.PostedBalance
	accountFull.Currency = banking.Currency(resp.Currency)
	accountFull.Opened = resp.Opened

	sql := `
		INSERT INTO business_bank_account (
			business_id, bank_account_id, account_holder_id, account_status, account_number,
			routing_number, available_balance, posted_balance, currency, opened, alias,
			account_type, usage_type, bank_name
		)
	 	VALUES(
			:business_id, :bank_account_id, :account_holder_id, :account_status, :account_number,
			:routing_number, :available_balance, :posted_balance, :currency, :opened, :alias,
			:account_type, :usage_type, :bank_name
		)
		RETURNING *`

	stmt, err := store.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	ba := &BankAccount{}
	err = stmt.Get(ba, accountFull)

	// Register created bank account
	NewLinkedAccountService(store.sourceReq).LinkOwnBankAccount(ba)

	return ba, err
}

//GetByID returns a business bank account by the giving id
func (store *bankAccountDataStore) GetByID(accountID string, businessID shared.BusinessID) (*BankAccount, error) {
	// Check access
	err := auth.NewAuthService(store.sourceReq).CheckBusinessBankAccountAccess(accountID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.GetByID(accountID)
	}

	return store.getByID(accountID)
}

// GetByIdInternal used for internal purpose
func (store *bankAccountDataStore) GetByIDInternal(accountID string) (*BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.GetByID(accountID)
	}

	return store.getByID(accountID)
}

// getByID returns a business bank account by the giving id
func (store *bankAccountDataStore) getByID(accountID string) (*BankAccount, error) {
	account := &BankAccount{}
	err := store.Get(account, "SELECT * FROM business_bank_account WHERE id = $1", accountID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

//GetByBankAccountId returns a business bank account by the giving bank account id
func (store *bankAccountDataStore) GetByBankAccountId(bankAccountId string, businessId shared.BusinessID) (*BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.GetByBankAccountId(bankAccountId)
	}

	account := &BankAccount{}
	err := store.Get(account, "SELECT * FROM business_bank_account WHERE bank_account_id = $1 AND business_id = $2", bankAccountId, businessId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return account, nil
}

//GetByBankAccountId returns a business bank account by the giving bank account id and user id
func (store *bankAccountDataStore) GetByBankAccountIDUserID(bankAccountID string, userID shared.UserID) (*BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.GetByBankAccountId(bankAccountID)
	}

	account := &BankAccount{}
	err := store.Get(account, "SELECT * FROM business_bank_account WHERE bank_account_id = $1 AND account_holder_id = $2", bankAccountID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return account, nil
}

//GetById returns a business bank account by the giving id
func (store *bankAccountDataStore) GetByUserID(userID shared.UserID, businessID shared.BusinessID) (*[]BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		as, err := bas.GetByBusinessID(businessID, 100, 0)

		return &as, err
	}

	account := &[]BankAccount{}

	err := store.Select(
		account,
		`
		SELECT * FROM business_bank_account
		WHERE account_holder_id = $1 AND business_id = $2 AND usage_type = $3`,
		userID,
		businessID,
		UsageTypePrimary,
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (store *bankAccountDataStore) GetByUsageType(userID shared.UserID, businessID shared.BusinessID, usage UsageType) ([]BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.GetByUsageType(businessID, usage, 100, 0)
	}

	account := []BankAccount{}
	err := store.Select(
		&account,
		`
		SELECT * FROM business_bank_account
		WHERE account_holder_id = $1 AND business_id = $2 AND usage_type = $3`,
		userID,
		businessID,
		usage,
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}

//List will return all the accounts for a giving business Id
//Params: `businessId`: the business id `limit`: the limit of items to be return `offset`: the offset
func (store *bankAccountDataStore) List(businessId shared.BusinessID, limit, offset int) (*[]BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		ba, err := bas.GetByBusinessID(businessId, limit, offset)

		return &ba, err
	}

	q := `
	SELECT * FROM business_bank_account
	WHERE business_id = $1 AND account_holder_id = $2 AND usage_type = $3 ORDER BY created ASC OFFSET $4 LIMIT $5`
	accounts := &[]BankAccount{}
	err := store.Select(accounts, q, businessId, store.sourceReq.UserID, UsageTypePrimary, offset, limit)

	if err != nil && err == sql.ErrNoRows {
		return accounts, nil
	}

	return accounts, err
}

func (store *bankAccountDataStore) ListInternalByBusiness(businessId shared.BusinessID, limit, offset int) ([]*BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		bs, err := bas.GetByBusinessID(businessId, limit, offset)

		var ba []*BankAccount

		for _, b := range bs {
			ba = append(ba, &b)
		}

		return ba, err
	}

	q := `
        SELECT * FROM business_bank_account
        WHERE business_id = $1 AND usage_type = $2 ORDER BY created ASC OFFSET $3 LIMIT $4`
	accounts := []*BankAccount{}
	err := store.Select(&accounts, q, businessId, UsageTypePrimary, offset, limit)
	if err != nil && err == sql.ErrNoRows {
		return accounts, nil
	}

	return accounts, err
}

func (store *bankAccountDataStore) ListInternal(limit, offset int) ([]*BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.ListInternal(limit, offset)
	}

	q := `
        SELECT * FROM business_bank_account ORDER BY id ASC OFFSET $1 LIMIT $2`
	accounts := []*BankAccount{}
	err := store.Select(&accounts, q, offset, limit)
	if err != nil && err == sql.ErrNoRows {
		return accounts, nil
	}

	return accounts, err
}

func (store *bankAccountDataStore) ListByOwner(owner string) ([]BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		q := `SELECT consumer_id FROM wise_user WHERE id = $1`

		var cID shared.ConsumerID

		err = store.Select(&cID, q, owner)
		if err != nil {
			return nil, err
		}

		return bas.GetByConsumerIDAndUsageType(cID, UsageTypePrimary, 100, 0)
	}

	q := `SELECT * FROM business_bank_account WHERE account_holder_id = $1 AND usage_type = $2`
	var accounts []BankAccount
	err := store.Select(&accounts, q, owner, UsageTypePrimary)

	if err != nil && err == sql.ErrNoRows {
		return []BankAccount{}, nil
	}

	return accounts, err
}

//Update will update a business bank account
//Currently only the `alias` field is updatable
//Params: `id`: the account id `updates`: the bank account updates
func (store *bankAccountDataStore) Update(update *BankAccountUpdate) (*BankAccount, error) {
	// Check access
	err := auth.NewAuthService(store.sourceReq).CheckBusinessBankAccountAccess(update.Id)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.Update(update)
	}

	account := &BankAccount{}

	sql := `UPDATE business_bank_account SET alias = $2 WHERE id = $1
			 RETURNING *`
	err = store.Get(account, sql, update.Id, update.Alias)

	if err != nil {
		return nil, err
	}

	return account, err
}

func (store *bankAccountDataStore) GetBalanceByIDInternal(accountID string) (*BankAccountBalance, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		ba, err := bas.GetBalanceByID(accountID, true)
		if err != nil {
			return nil, err
		}

		return &BankAccountBalance{
			AccountID:        ba.Id,
			AvailableBalance: ba.AvailableBalance,
			PostedBalance:    ba.PostedBalance,
			ActualBalance:    ba.ActualBalance,
			Currency:         banking.CurrencyUSD,
			Modified:         ba.Modified,
		}, nil
	}

	// Get current account
	account, err := store.getByID(accountID)
	if err != nil {
		return nil, err
	}

	// Used cached value if less than an hour
	if time.Now().Sub(account.Modified).Seconds() <= 3600 {
		balance := account.toAccountBalance()
		return &balance, nil
	}

	// Get account info from bank
	providerName, ok := banking.ToPartnerBankName[account.BankName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Partner bank %s does not exist", account.BankName))
	}

	bank, err := partnerbank.GetBusinessBank(providerName)
	if err != nil {
		return nil, err
	}

	srv, err := bank.BankAccountService(store.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(account.BusinessID))
	if err != nil {
		return nil, err
	}

	resp, err := srv.Get(partnerbank.AccountBankID(account.BankAccountId))
	if err != nil {
		return nil, err
	}

	balance := BankAccountBalance{}
	sql := `
        UPDATE business_bank_account
        SET available_balance = $1, posted_balance = $2
        WHERE id = $3
        RETURNING id, available_balance, posted_balance, currency, modified`
	err = store.Get(&balance, sql, resp.AvailableBalance, resp.PostedBalance, accountID)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (store *bankAccountDataStore) GetBalanceByID(accountID string, businessID shared.BusinessID) (*BankAccountBalance, error) {
	// Check access
	err := auth.NewAuthService(store.sourceReq).CheckBusinessBankAccountAccess(accountID)
	if err != nil {
		return nil, err
	}

	return store.GetBalanceByIDInternal(accountID)
}

func (store *bankAccountDataStore) UpdateBankAccountStatus(accountID string, status string) (*BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		update := &BankAccountUpdate{}
		update.Id = accountID
		update.Status = status

		return bas.Update(update)
	}

	account := &BankAccount{}

	sql := `UPDATE business_bank_account SET account_status = $1 WHERE id = $2
			 RETURNING *`
	err := store.Get(account, sql, status, accountID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return nil, err
}

func (store *bankAccountDataStore) DeactivateAccount(accountID string, businessID shared.BusinessID, reason grpcBanking.AccountStatusReason) (*BankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		return bas.DeactivateAccount(accountID, reason)
	}

	account, err := store.GetByIDInternal(accountID)
	if err != nil {
		return nil, err
	}

	// Get account info from bank
	providerName, ok := banking.ToPartnerBankName[account.BankName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Partner bank %s does not exist", account.BankName))
	}

	bank, err := partnerbank.GetBusinessBank(providerName)
	if err != nil {
		return nil, err
	}

	srv, err := bank.BankAccountService(store.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(businessID))
	if err != nil {
		return nil, err
	}

	resp, err := srv.Close(partnerbank.AccountBankID(account.BankAccountId), partnerbank.AccountCloseReasonCustomer)
	if err != nil {
		return nil, err
	}

	ua := &BankAccount{}
	sql := `
	UPDATE business_bank_account
	SET account_status = $1, available_balance = $2, posted_balance = $3
	WHERE id = $4
	RETURNING *`
	err = store.Get(ua, sql, resp.Status, resp.AvailableBalance, resp.PostedBalance, accountID)
	return ua, err
}
