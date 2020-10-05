package business

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/analytics"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/service/segment"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/data"
	coreData "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

// CSPBusinessService is a service to fetch business from csp db
// CSPBusiness is the presentation on core business on csp with less fields additional fields to handle status
type CSPBusinessService interface {
	CSPBusinessByID(id string) (CSPBusiness, error)
	CSPBusinessList(params map[string]interface{}) ([]CSPBusiness, error)
	CSPBusinessListStatus(limit, offset int, status csp.Status, reviewSubstatus csp.ReviewSubstatus) ([]CSPBusiness, error)
	CSPBusinessCreate(CSPBusinessCreate) (CSPBusiness, error)
	CSPBusinessUpdate(CSPBusinessUpdate, string) (CSPBusiness, error)
	CSPBusinessUpdateByBusinessID(id shared.BusinessID, updates CSPBusinessUpdate) (CSPBusiness, error)
	UpdateProcessStatus(businessID shared.BusinessID, status csp.ProcessStatus) error
	ByBusinessID(shared.BusinessID) (CSPBusiness, error)
	UpdateBusinessNameAndEntityType(id shared.BusinessID, name string, entityName string) error

	UpdateSubscribedAgentID(shared.BusinessID, string) error
}

type cspBusinessService struct {
	rdb     *sqlx.DB
	wdb     *sqlx.DB
	coreRDB *sqlx.DB
}

// NewCSPService the service for CRUD on csp business
func NewCSPService() CSPBusinessService {
	return cspBusinessService{wdb: data.DBWrite, rdb: data.DBRead, coreRDB: coreData.DBRead}
}

func (s cspBusinessService) CSPBusinessCreate(create CSPBusinessCreate) (CSPBusiness, error) {
	var business CSPBusiness
	keys := services.SQLGenInsertKeys(create)
	values := services.SQLGenInsertValues(create)

	q := fmt.Sprintf("INSERT INTO business (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := s.wdb.PrepareNamed(q)
	if err != nil {
		return business, err
	}

	if err = stmt.Get(&business, create); err != nil {
		return business, err
	}

	if business.ProcessStatus != "" || business.Status != "" {
		NewStateService().Create(BusinessStateCreate{
			BusinessID:    business.ID,
			ProcessStatus: business.ProcessStatus,
			Status:        business.Status,
		})
	}

	return business, err
}

func (s cspBusinessService) CSPBusinessList(params map[string]interface{}) ([]CSPBusiness, error) {
	var list = make([]CSPBusiness, 0)

	if params["businessId"] != nil {
		businessID := params["businessId"].(shared.BusinessID)
		b, err := s.ByBusinessID(businessID)
		if err != nil {
			return nil, err
		}
		return append(list, b), nil
	}

	if params["bankId"] != nil {
		bankID := params["bankId"].(partnerbank.BusinessBankID)

		bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
		if err != nil {
			log.Printf("Error getting partner bank %v", err)
			return nil, err
		}

		sr := services.NewSourceRequest()
		businessID, err := bank.ProxyService(sr.PartnerBankRequest()).GetBusinessID(bankID)
		if err != nil {
			return nil, err
		}

		b, err := s.ByBusinessID(shared.BusinessID(*businessID))
		if err != nil {
			return nil, err
		}

		return append(list, b), nil
	}

	if params["userId"] != nil {
		userID := params["userId"].(shared.UserID)
		b, err := s.ByUserID(userID)
		if err != nil {
			return nil, err
		}
		return append(list, b...), nil
	}

	if params["consumerBankId"] != nil {
		consumerBankID := params["consumerBankId"].(partnerbank.ConsumerBankID)
		b, err := s.ByConsumerBankID(consumerBankID)
		if err != nil {
			return nil, err
		}
		return append(list, b...), nil
	}

	if params["ownerFirstName"] != nil {
		ownerFName := params["ownerFirstName"].(string)
		b, err := s.ByOwnerFirstName(ownerFName)
		if err != nil {
			return nil, err
		}
		return append(list, b...), nil
	}

	if params["ownerPhoneNumber"] != nil {
		ownerPhoneName := params["ownerPhoneNumber"].(string)
		b, err := s.ByOwnerPhoneNumber(ownerPhoneName)
		if err != nil {
			return nil, err
		}
		return append(list, b...), nil
	}

	if params["ownerEmailId"] != nil {
		ownerEmailID := params["ownerEmailId"].(string)
		b, err := s.ByOwnerEmailID(ownerEmailID)
		if err != nil {
			return nil, err
		}
		return append(list, b...), nil
	}

	filter := ""
	if params["name"] != nil && len(params["name"].(string)) > 0 {
		filter = " WHERE business_name ILIKE  '%" + params["name"].(string) + "%'"
	}

	if params["submitStart"] != nil {
		clause := "created >= '" + params["submitStart"].(string) + "'"
		if len(filter) > 0 {
			filter = filter + " AND " + clause
		} else {
			filter = " WHERE " + clause
		}
	}

	if params["submitEnd"] != nil {
		clause := "created <= '" + params["submitEnd"].(string) + "'"
		if len(filter) > 0 {
			filter = filter + " AND " + clause
		} else {
			filter = " WHERE " + clause
		}
	}
	limit := ""
	if params["limit"] != nil && params["offset"] != nil {
		limit = fmt.Sprintf("LIMIT %d OFFSET %d", params["limit"].(int), params["offset"].(int))
	}
	query := "SELECT * FROM business" + filter + " ORDER BY created DESC " + limit

	err := s.rdb.Select(&list, query)
	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}

	return list, err
}

// get list csp business by status
func (s cspBusinessService) CSPBusinessListStatus(limit, offset int, status csp.Status, reviewSubstatus csp.ReviewSubstatus) ([]CSPBusiness, error) {
	var list = make([]CSPBusiness, 0)
	var err error
	if reviewSubstatus == "" {
		err = s.rdb.Select(&list, "SELECT * FROM business WHERE review_status = $1 ORDER BY created DESC LIMIT $2 OFFSET $3", status, limit, offset)
	} else {
		err = s.rdb.Select(&list, "SELECT * FROM business WHERE review_status = $1 AND review_substatus = $2 ORDER BY created DESC LIMIT $3 OFFSET $4", status, reviewSubstatus, limit, offset)
	}

	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

func (s cspBusinessService) CSPBusinessByID(id string) (CSPBusiness, error) {
	var item CSPBusiness
	err := s.rdb.Get(&item, "SELECT * FROM business WHERE id = $1", id)
	if err == sql.ErrNoRows {
		log.Printf("no business %v", err)
		return item, services.ErrorNotFound{}.New("")
	}
	return item, err
}

// Get csp business by core business id
func (s cspBusinessService) ByBusinessID(id shared.BusinessID) (CSPBusiness, error) {
	var item CSPBusiness
	err := s.rdb.Get(&item, "SELECT * FROM business WHERE business_id = $1", id)
	if err == sql.ErrNoRows {
		log.Printf("no business %v", err)
		return item, services.ErrorNotFound{}.New("")
	}
	return item, err
}

// Get csp business by core user id
func (s cspBusinessService) ByUserID(ID shared.UserID) ([]CSPBusiness, error) {
	coreBusinessList := []Business{}
	var item []CSPBusiness

	err := s.coreRDB.Select(&coreBusinessList, "SELECT * FROM business WHERE owner_id = $1", ID)
	if err != nil {
		log.Printf("no core business %v", err)
		return item, err
	}

	var businessIDsAsStr string
	for _, theBusiness := range coreBusinessList {
		comma := ","
		if len(businessIDsAsStr) == 0 {
			comma = ""
		}
		businessIDsAsStr = businessIDsAsStr + comma + "'" + string(theBusiness.ID) + "'"
	}

	businessQuery := "SELECT * FROM business WHERE business_id in " + "(" + businessIDsAsStr + ")"
	err = s.rdb.Select(&item, businessQuery)
	if err == sql.ErrNoRows {
		log.Printf("no business %v", err)
		return item, services.ErrorNotFound{}.New("")
	}
	return item, err
}

// Get csp business by core consumer bank id
func (s cspBusinessService) ByConsumerBankID(ID partnerbank.ConsumerBankID) ([]CSPBusiness, error) {
	var item []CSPBusiness

	bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		log.Printf("Error getting partner bank %v", err)
		return nil, err
	}

	sr := services.NewSourceRequest()
	consumerID, err := bank.ProxyService(sr.PartnerBankRequest()).GetConsumerID(ID)
	if err != nil {
		return nil, err
	}

	var userID shared.UserID
	err = s.coreRDB.Get(&userID, "SELECT id FROM wise_user WHERE consumer_id = $1", consumerID)
	if err != nil {
		log.Printf("no business %v", err)
		return item, err
	}

	return s.ByUserID(userID)
}

// Get csp business by owner's first name
func (s cspBusinessService) ByOwnerFirstName(firstName string) ([]CSPBusiness, error) {
	var item []CSPBusiness
	var bIDs shared.BusinessIDs

	lf := strings.ToLower(firstName)
	err := s.coreRDB.Select(&bIDs, "SELECT b.id FROM business b LEFT JOIN wise_user w ON b.owner_id = w.id LEFT JOIN consumer c ON w.consumer_id = c.id WHERE LOWER(c.first_name) like $1", lf+"%")
	if err != nil {
		log.Printf("no business %v", err)
		return item, err
	}
	return s.ByBusinessIDs(bIDs)
}

// Get csp business by owner's phone number
func (s cspBusinessService) ByOwnerPhoneNumber(phone string) ([]CSPBusiness, error) {
	var item []CSPBusiness
	var bIDs shared.BusinessIDs

	err := s.coreRDB.Select(&bIDs, "SELECT b.id FROM business b LEFT JOIN wise_user w ON b.owner_id = w.id LEFT JOIN consumer c ON w.consumer_id = c.id WHERE c.phone = $1 OR w.phone = $2", phone, phone)
	if err != nil {
		log.Printf("no business %v", err)
		return item, err
	}
	return s.ByBusinessIDs(bIDs)
}

// Get csp business by owner's email ID
func (s cspBusinessService) ByOwnerEmailID(email string) ([]CSPBusiness, error) {
	var item []CSPBusiness
	var bIDs shared.BusinessIDs

	err := s.coreRDB.Select(&bIDs, "SELECT b.id FROM business b LEFT JOIN wise_user w ON b.owner_id = w.id LEFT JOIN consumer c ON w.consumer_id = c.id WHERE c.email = $1", email)
	if err != nil {
		log.Printf("no business %v", err)
		return item, err
	}
	return s.ByBusinessIDs(bIDs)
}

func (s cspBusinessService) ByBusinessIDs(IDs shared.BusinessIDs) ([]CSPBusiness, error) {
	var item []CSPBusiness
	businessIDsAsStr := IDs.Join("','")
	businessQuery := "SELECT * FROM business WHERE business_id IN " + "('" + businessIDsAsStr + "')"

	err := s.rdb.Select(&item, businessQuery)
	if err == sql.ErrNoRows {
		log.Printf("no business %v", err)
		return item, services.ErrorNotFound{}.New("")
	}
	return item, err
}

// UpdateProcessStatus ..
func (s cspBusinessService) UpdateProcessStatus(businessID shared.BusinessID, status csp.ProcessStatus) error {
	st := status
	updates := CSPBusinessUpdate{ProcessStatus: &st}
	_, err := s.CSPBusinessUpdateByBusinessID(businessID, updates)
	if err != nil {
		log.Printf("error updating process status %v", err)
	}
	return err
}

func (s cspBusinessService) CSPBusinessUpdate(updates CSPBusinessUpdate, businessID string) (CSPBusiness, error) {
	var business CSPBusiness
	keys := services.SQLGenForUpdate(updates)
	q := fmt.Sprintf("UPDATE business SET %s WHERE id = '%s' RETURNING *", keys, businessID)
	stmt, err := s.wdb.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return business, err
	}

	err = stmt.Get(&business, updates)
	if err != nil {
		log.Printf("error updating business review %v", err)
		return business, err
	}
	// update business status
	if updates.ProcessStatus != nil || updates.Status != nil {
		NewStateService().Create(BusinessStateCreate{
			BusinessID:    businessID,
			ProcessStatus: business.ProcessStatus,
			Status:        business.Status,
		})
	}
	return business, err
}

func (s cspBusinessService) CSPBusinessUpdateByBusinessID(id shared.BusinessID, updates CSPBusinessUpdate) (CSPBusiness, error) {
	var b CSPBusiness
	current, err := s.ByBusinessID(id)
	if err != nil {
		return b, err
	}
	keys := services.SQLGenForUpdate(updates)
	q := fmt.Sprintf("UPDATE business SET %s WHERE business_id = '%s' RETURNING *", keys, id)
	stmt, err := s.wdb.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return b, err
	}

	err = stmt.Get(&b, updates)
	if err != nil {
		log.Printf("error updating business review %v", err)
	}
	sendToAnalytics(id, updates)

	// only update when statuses have changed
	if current.Status != b.Status || current.ProcessStatus != b.ProcessStatus {
		if updates.ProcessStatus != nil || updates.Status != nil {
			if _, err := NewStateService().Create(BusinessStateCreate{
				BusinessID:    b.ID,
				ProcessStatus: b.ProcessStatus,
				Status:        b.Status,
			}); err != nil {
				log.Printf("Error creating business state %v", err)
			}
		}
	}

	return b, err
}

func (s cspBusinessService) UpdateBusinessNameAndEntityType(id shared.BusinessID, name string, entityType string) error {

	updateSet := ""
	if entityType != "" {
		updateSet = fmt.Sprintf("entity_type = '%v'", entityType)
	}
	if name != "" {
		if len(updateSet) > 0 {
			updateSet = updateSet + ","
		}
		updateSet = fmt.Sprintf("%v business_name = '%v'", updateSet, name)
	}

	if updateSet == "" {
		return nil
	}

	q := fmt.Sprintf("UPDATE business SET %v WHERE business_id = '%v'", updateSet, id)

	_, err := s.wdb.Exec(q)
	if err != nil {
		log.Printf("Error updating legalName/entity name %v", err)
	}

	return err
}

func sendToAnalytics(id shared.BusinessID, updates CSPBusinessUpdate) {
	coreB, _ := New(services.SourceRequest{}).ByID(id)
	seg := analytics.CSPBusinessUpdate{}
	if updates.Status != nil {
		status := updates.Status.String()
		kyb := status
		seg.KYCBStatus = &kyb
	}

	if updates.PromoFunded != nil {
		seg.PromoFunded = updates.PromoFunded
		var amount float64 = 100
		seg.Amount = &amount
	}

	segment.NewSegmentService().PushToAnalyticsQueue(coreB.OwnerID, segment.CategoryBusiness, segment.ActionCSP, seg)
}

func (s cspBusinessService) UpdateSubscribedAgentID(businessID shared.BusinessID, agentID string) error {
	_, err := s.wdb.Exec("UPDATE business SET subscribed_agent_id = $1 WHERE business_id = $2", agentID, businessID)
	if err != nil {
		log.Println(err)
		return errors.Cause(err)
	}

	return nil
}
