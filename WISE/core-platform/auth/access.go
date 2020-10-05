package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	grpcBankAccount "github.com/wiseco/protobuf/golang/banking/account"
	grpcBankLinkedAccount "github.com/wiseco/protobuf/golang/banking/linked_account"
)

type authService struct {
	srcReq services.SourceRequest
	*sqlx.DB
}

type AuthService interface {
	CheckBusinessAccess(shared.BusinessID) error
	CheckUserDeviceAccess(shared.UserDeviceID) error
	CheckUserAccess(shared.UserID) error
	CheckUserAccessStrict(api.APIRequest) error
	CheckConsumerAccess(shared.ConsumerID) error
	CheckConsumerBankAccountAccess(string) error
	CheckBusinessBankAccountAccess(string) error
	CheckBusinessBankCardAccess(string) error
	CheckBusinessLinkedCardAccess(string) error
	CheckBusinessLinkedAccountAccess(string) error
	CheckBusinessCardReaderAccess(string) error
}

//TODO when we move auth, there shouldn't be a cyclical import issue any longer
type bankingLinkedAccountService struct {
	conn grpc.Client
	lac  grpcBankLinkedAccount.LinkedAccountServiceClient
}

func NewBankingLinkedAccountService() (*bankingLinkedAccountService, error) {
	var bts *bankingLinkedAccountService

	bsn, err := grpc.GetConnectionStringForService(grpc.ServiceNameBanking)
	if err != nil {
		return bts, err
	}

	conn, err := grpc.NewInsecureClient(bsn)
	if err != nil {
		return bts, err
	}

	lac := grpcBankLinkedAccount.NewLinkedAccountServiceClient(conn.GetConn())

	return &bankingLinkedAccountService{
		conn: conn,
		lac:  lac,
	}, nil
}

//TODO when we move auth, there shouldn't be a cyclical import issue any longer
type bankingAccountService struct {
	conn grpc.Client
	ac   grpcBankAccount.AccountServiceClient
}

func NewBankingAccountService() (*bankingAccountService, error) {
	var bas *bankingAccountService

	bsn, err := grpc.GetConnectionStringForService(grpc.ServiceNameBanking)
	if err != nil {
		return bas, err
	}

	conn, err := grpc.NewInsecureClient(bsn)
	if err != nil {
		return bas, err
	}

	ac := grpcBankAccount.NewAccountServiceClient(conn.GetConn())

	return &bankingAccountService{
		conn: conn,
		ac:   ac,
	}, nil
}

func NewAuthService(r services.SourceRequest) AuthService {
	return &authService{r, data.DBWrite}
}

func (a *authService) CheckBusinessAccess(businessID shared.BusinessID) error {
	var ownerID shared.UserID
	err := a.Get(&ownerID, "SELECT owner_id FROM business WHERE id = $1", businessID)
	if err != nil {
		return err
	}
	log.Printf("Owner user id:%s sourceId:%s", ownerID, a.srcReq.UserID)
	// Check access
	if ownerID != a.srcReq.UserID {
		return errors.New("unauthorized")
	}

	return nil
}

func (a *authService) CheckUserDeviceAccess(deviceID shared.UserDeviceID) error {
	var userID shared.UserID
	err := a.Get(&userID, "SELECT user_id FROM user_device WHERE id = $1", deviceID)
	if err != nil {
		return err
	}

	if userID != a.srcReq.UserID {
		return errors.New("unauthorized")
	}

	return nil
}

func (a *authService) CheckUserAccessStrict(r api.APIRequest) error {

	tUID := shared.UserID(r.GetPathParam("userId"))

	if tUID != r.UserID {
		return errors.New("unauthorized")
	}

	var id shared.UserID
	err := a.Get(&id, "SELECT id FROM wise_user WHERE id = $1", r.UserID)
	if err != nil {
		return err
	}

	if id != a.srcReq.UserID {
		return errors.New("unauthorized")
	}

	return nil
}

func (a *authService) CheckUserAccess(userID shared.UserID) error {
	var id shared.UserID
	err := a.Get(&id, "SELECT id FROM wise_user WHERE id = $1", userID)
	if err != nil {
		return err
	}

	if id != a.srcReq.UserID {
		return errors.New("unauthorized")
	}

	return nil
}

func (a *authService) CheckConsumerAccess(consumerID shared.ConsumerID) error {
	var ownerID shared.UserID

	err := a.Get(&ownerID, `SELECT owner_id FROM business WHERE id IN (select business_id from business_member
		where consumer_id = $1)`, consumerID)
	if err == sql.ErrNoRows {
		err = a.Get(&ownerID, `SELECT id FROM wise_user WHERE consumer_id = $1`, consumerID)
		if err != nil {
			return err
		}

		if ownerID != a.srcReq.UserID {
			return errors.New("unauthorized")
		}

		return nil
	}

	if err != nil {
		return err
	}

	if ownerID != a.srcReq.UserID {
		return errors.New("unauthorized")
	}

	return nil
}

func (a *authService) CheckConsumerBankAccountAccess(accountID string) error {
	var accountHolderId shared.UserID
	err := a.Get(&accountHolderId, "SELECT account_holder_id FROM consumer_bank_account WHERE id = $1", accountID)
	if err != nil {
		return err
	}

	if accountHolderId != a.srcReq.UserID {
		return errors.New("unauthorized")
	}

	return nil
}

func (a *authService) CheckBusinessBankAccountAccess(accountID string) error {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := NewBankingAccountService()
		if err != nil {
			return err
		}

		defer bas.conn.CloseAndCancel()

		s := accountID
		if !strings.HasPrefix(accountID, id.IDPrefixBankAccount.String()) {
			s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, accountID)
		}

		baID, err := id.ParseBankAccountID(s)
		if err != nil {
			return err
		}

		gpr := &grpcBankAccount.GetParticipantsRequest{
			AccountId: baID.String(),
			ParticipantStatusFilter: []grpcBankAccount.ParticipantStatus{
				grpcBankAccount.ParticipantStatus_PS_ACTIVE,
			},
			RoleFilter: []grpcBankAccount.Role{
				grpcBankAccount.Role_R_HOLDER,
			},
		}

		ps, err := bas.ac.GetParticipants(context.Background(), gpr)
		if err != nil {
			return err
		}

		if len(ps.Participants) == 0 {
			return errors.New("Account holder participant not found")
		}

		authed := false
		for _, p := range ps.Participants {
			consID, err := id.ParseConsumerID(p.ConsumerId)
			if err != nil {
				return err
			}

			var accountHolderId shared.UserID

			err = a.Get(&accountHolderId, "SELECT id FROM wise_user WHERE consumer_id = $1", consID.UUIDString())
			if err != nil {
				return err
			}

			if accountHolderId == a.srcReq.UserID {
				authed = true
				break
			}
		}

		if authed == false {
			return errors.New("unauthorized")
		}
	} else {

		var accountHolderId shared.UserID
		err := a.Get(&accountHolderId, "SELECT account_holder_id FROM business_bank_account WHERE id = $1", accountID)
		if err != nil {
			return err
		}

		if accountHolderId != a.srcReq.UserID {
			return errors.New("unauthorized")
		}
	}

	return nil
}

func (a *authService) CheckBusinessBankCardAccess(cardID string) error {
	var cardholderID shared.UserID
	err := a.Get(&cardholderID, "SELECT cardholder_id FROM business_bank_card WHERE id = $1", cardID)
	if err != nil {
		return err
	}

	if cardholderID != a.srcReq.UserID {
		return errors.New("unauthorized")
	}

	return nil
}

func (a *authService) CheckBusinessLinkedCardAccess(linkedCardID string) error {
	var businessID shared.BusinessID
	err := a.Get(&businessID, "SELECT business_id FROM business_linked_card WHERE id = $1", linkedCardID)
	if err != nil {
		return err
	}

	return a.CheckBusinessAccess(businessID)
}

func (a *authService) CheckBusinessLinkedAccountAccess(linkedAccountID string) error {
	var businessID shared.BusinessID

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		blas, err := NewBankingLinkedAccountService()
		if err != nil {
			return err
		}

		defer blas.conn.CloseAndCancel()

		s := linkedAccountID
		if !strings.HasPrefix(linkedAccountID, id.IDPrefixLinkedBankAccount.String()) {
			s = fmt.Sprintf("%s%s", id.IDPrefixLinkedBankAccount, linkedAccountID)
		}

		laID, err := id.ParseLinkedBankAccountID(s)
		if err != nil {
			return err
		}

		gr := &grpcBankLinkedAccount.GetRequest{
			LinkedAccountId: laID.String(),
		}

		lbap, err := blas.lac.Get(context.Background(), gr)
		if err != nil {
			return err
		}

		bID, err := id.ParseBusinessID(lbap.BusinessId)
		if err != nil {
			return err
		}

		businessID = shared.BusinessID(bID.UUIDString())
	} else {
		err := a.Get(&businessID, "SELECT business_id FROM business_linked_bank_account WHERE id = $1", linkedAccountID)
		if err != nil {
			return err
		}
	}

	return a.CheckBusinessAccess(businessID)
}

func (a *authService) CheckBusinessCardReaderAccess(cardReaderID string) error {

	var businessID shared.BusinessID
	err := a.Get(&businessID, "SELECT business_id FROM card_reader WHERE id = $1", cardReaderID)
	if err != nil {
		log.Println(err)
		return err
	}

	return a.CheckBusinessAccess(businessID)
}
