/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package support

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	banking "github.com/wiseco/core-platform/services/banking"
	busBanking "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type AccountBlockType string

type AccountBlock struct {
	ID             string           `json:"id" db:"id"`
	AccountID      string           `json:"accountId" db:"account_id"`
	BlockID        string           `json:"blockId" db:"block_id"`
	BlockType      AccountBlockType `json:"blockType" db:"block_type"`
	Reason         string           `json:"reason" db:"reason"`
	Created        time.Time        `json:"created" db:"created"`
	Deactivated    *time.Time       `json:"deactivated" db:"deactivated"`
	OriginatedFrom string           `json:"originatedFrom" db:"originated_from"`
}

const (
	AccountBlockTypeDebits  = AccountBlockType("debits")
	AccountBlockTypeCredits = AccountBlockType("credits")
	AccountBlockTypeChecks  = AccountBlockType("checks")
	AccountBlockTypeAll     = AccountBlockType("all")
)

var accountBlockTypes = map[AccountBlockType]AccountBlockType{
	AccountBlockTypeDebits:  AccountBlockTypeDebits,
	AccountBlockTypeCredits: AccountBlockTypeCredits,
	AccountBlockTypeChecks:  AccountBlockTypeChecks,
	AccountBlockTypeAll:     AccountBlockTypeAll,
}

func (a AccountBlockType) Valid() bool {
	_, ok := accountBlockTypes[a]
	return ok
}

//Blocking We will use this later, atm we will block all.
type Blocking struct {
	Reason         string           `json:"reason"`
	Type           AccountBlockType `json:"type"`
	OriginatedFrom string           `json:"originatedFrom"`
}

func (t AccountBlockType) String() string {
	return string(t)
}

//Support  ..
// TODO: USE_BANKING_SERVICE
type Support interface {
	Block(id string, blocking Blocking) error
	Unblock(accountID, id string) error
	ListOfAccountBlocks(id string) ([]AccountBlock, error)
}

type support struct {
	*sqlx.DB
	sourceReq services.SourceRequest
}

//NewSupport ..
func NewSupport(sourceReq services.SourceRequest) Support {
	return support{data.DBWrite, sourceReq}
}

//ListAccount ..
func ListAccount(userID string) ([]busBanking.BankAccount, error) {
	return busBanking.NewAccountService().ListByOwner(userID)
}

func (s support) GetBankBlocks(id string) error {
	var bus struct {
		BankAccountID   string            `db:"bank_account_id"`
		AccountHolderID shared.UserID     `db:"account_holder_id"`
		BusinessID      shared.BusinessID `db:"business_id"`
	}
	if err := s.Get(&bus, "SELECT bank_account_id, account_holder_id, business_id FROM business_bank_account WHERE id = $1", id); err != nil {
		log.Printf("error getting bank account %v", err)
		return err
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return err
	}

	srv, err := bank.BankAccountService(s.sourceReq.CSPPartnerBankRequest(bus.AccountHolderID), partnerbank.BusinessID(bus.BusinessID))
	if err != nil {
		return err
	}

	resp, err := srv.GetAllBlocks(partnerbank.AccountBankID(bus.BankAccountID))
	if err != nil {
		log.Printf("Error getting blocks from account id: %s error:%v", id, err)
		return err
	}
	log.Println(resp)
	return nil
}

//Block blocks an account for the giving id
func (s support) Block(id string, blocking Blocking) error {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := busBanking.NewBankingAccountService()
		if err != nil {
			return err
		}

		abc := banking.AccountBlockCreate{
			AccountID:      id,
			BlockType:      banking.AccountBlockType(blocking.Type),
			Reason:         blocking.Reason,
			OriginatedFrom: blocking.OriginatedFrom,
		}

		_, err = bas.CreateAccountBlock(abc)

		return err
	} else {
		var bus struct {
			BankAccountID   string            `db:"bank_account_id"`
			AccountHolderID shared.UserID     `db:"account_holder_id"`
			BusinessID      shared.BusinessID `db:"business_id"`
			Status          string            `db:"account_status"`
		}
		if err := s.Get(&bus, "SELECT  bank_account_id,business_id,account_holder_id,account_status FROM business_bank_account WHERE id = $1", id); err != nil {
			log.Printf("error getting bank account %v", err)
			return err
		}

		req := partnerbank.AccountBlockRequest{
			AccountID: partnerbank.AccountBankID(bus.BankAccountID),
			Type:      partnerbank.AccountBlockType(blocking.Type),
			Reason:    blocking.Reason,
		}

		bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
		if err != nil {
			return err
		}

		srv, err := bank.BankAccountService(s.sourceReq.CSPPartnerBankRequest(bus.AccountHolderID), partnerbank.BusinessID(bus.BusinessID))
		if err != nil {
			return err
		}

		resp, err := srv.Block(req)
		if err != nil {
			log.Printf("error blocking account with id:%s error:%v", id, err)
			return err
		}
		if resp != nil {
			return s.updateAccountActionBlock(id, blocking.Reason, blocking.OriginatedFrom, *resp)
		}
	}

	return nil
}

func (s support) updateAccountActionBlock(id, reason, origin string, resp partnerbank.AccountBlockResponse) error {
	block := `INSERT INTO business_bank_account_block(account_id, block_id, block_type, reason,originated_from)
			   VALUES($1,$2,$3,$4,$5)`
	account := `UPDATE business_bank_account SET account_status = 'blocked' WHERE id = $1`
	trx := s.MustBegin()
	trx.MustExec(block, id, resp.BlockID, resp.Type, reason, origin)
	trx.MustExec(account, id)
	err := trx.Commit()
	if err != nil {
		log.Printf("Error creating blocks for business bank account and updating bank account status error:%v", err)
	}
	return err
}

//Ublock blocks an account for the giving id
func (s support) Unblock(accountID, id string) error {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := busBanking.NewBankingAccountService()
		if err != nil {
			return err
		}

		err = bas.DeactivateAccountBlock(id)

		return err
	}

	var bus struct {
		BankAccountID   string            `db:"bank_account_id"`
		AccountHolderID shared.UserID     `db:"account_holder_id"`
		BusinessID      shared.BusinessID `db:"business_id"`
		Status          string            `db:"account_status"`
		BlockID         string            `db:"block_id"`
	}
	if err := s.Get(&bus, `	SELECT  a.bank_account_id,a.business_id,
								   a.account_holder_id,a.account_status, b.block_id
							FROM business_bank_account AS a
							JOIN business_bank_account_block AS b
							ON a.id = b.account_id
							WHERE a.id = $1 AND b.id = $2 AND deactivated IS NULL`, accountID, id); err != nil {
		log.Printf("error getting bank account %v", err)
		return err
	}
	if bus.Status != "blocked" {
		return fmt.Errorf("Account status is invalid")
	}
	req := partnerbank.AccountUnblockRequest{
		AccountID: partnerbank.AccountBankID(bus.BankAccountID),
		BlockID:   partnerbank.AccountBlockBankID(bus.BlockID),
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return err
	}

	srv, err := bank.BankAccountService(s.sourceReq.CSPPartnerBankRequest(bus.AccountHolderID), partnerbank.BusinessID(bus.BusinessID))
	if err != nil {
		return err
	}

	if err := srv.Unblock(req); err != nil {
		return err
	}

	return s.updateAccountActionUnblock(accountID, id)

}

func (s support) updateAccountActionUnblock(accountID, blockID string) error {

	var count int
	s.Get(&count, `SELECT COUNT(*) FROM business_bank_account_block 
	WHERE account_id = $1 AND deactivated IS NULL`, accountID)
	trx := s.MustBegin()
	block := `UPDATE business_bank_account_block SET deactivated = CURRENT_TIMESTAMP WHERE account_id = $1 AND id = $2`
	trx.MustExec(block, accountID, blockID)
	if count == 1 {
		account := `UPDATE business_bank_account SET account_status = 'active' WHERE id = $1`
		trx.MustExec(account, accountID)
	}
	err := trx.Commit()
	if err != nil {
		log.Printf("Error updating blocks error:%v", err)
	}
	return err
}

//ListOfAccountBlocks( get the list of active blocks
func (s support) ListOfAccountBlocks(id string) ([]AccountBlock, error) {
	var blocks []AccountBlock

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := busBanking.NewBankingAccountService()
		if err != nil {
			return blocks, err
		}

		bs, err := bas.GetManyAccountBlocks(id)
		if err != nil {
			return blocks, err
		}

		for _, block := range bs {
			b := AccountBlock{
				ID:             block.ID,
				AccountID:      block.AccountID,
				BlockID:        block.BlockID,
				BlockType:      AccountBlockType(block.BlockType),
				Reason:         block.Reason,
				Created:        block.Created,
				Deactivated:    block.Deactivated,
				OriginatedFrom: block.OriginatedFrom,
			}

			blocks = append(blocks, b)
		}
	} else {
		err := s.Select(&blocks, "SELECT * FROM business_bank_account_block WHERE account_id = $1 AND deactivated IS NULL", id)
		if err != nil {
			return blocks, err
		}
	}

	return blocks, nil
}
