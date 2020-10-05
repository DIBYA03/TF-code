package accountclosure

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	bus "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/csp/cspuser"
	cspsrv "github.com/wiseco/core-platform/services/csp/services"
	coreData "github.com/wiseco/core-platform/services/data"
	user "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/num"
)

//CSPAccountClosureService ...
type CSPAccountClosureService interface {
	CSPClosureRequestCreate(create CSPAccountClosureCreate) (CSPClosureRequestItem, error)
	CSPClosureRequestList(params CSPAccountClosureQueryParams) ([]CSPClosureRequestItem, error)
	CSPClosureApprovedAndRetryRequestList() ([]CSPClosureRequestItem, error)
	CSPClosureRequestDetails(id string) (CSPClosureRequestItem, error)
	CSPClosureRequestUpdate(id string, status string) (*CSPClosureRequestItem, error)
	CSPClosureRequestProcessed(id string, amount num.Decimal, status AccountClosureStatus) (*CSPClosureRequestItem, error)

	CSPClosureStateList(closureRequestID string) ([]CSPClosureState, error)
	CSPClosureStateAddNew(create CSPClosureStatePostBody) (CSPClosureState, error)
}

type cspAccountClosrueService struct {
	coreRDB *sqlx.DB
	coreWDB *sqlx.DB
	sr      cspsrv.SourceRequest
}

//NewCSPService ...
func NewCSPService(sr cspsrv.SourceRequest) CSPAccountClosureService {
	return cspAccountClosrueService{coreRDB: coreData.DBRead, coreWDB: coreData.DBWrite, sr: sr}
}

//NewCSPServiceWithoutRequest ...
func NewCSPServiceWithoutRequest() CSPAccountClosureService {
	return cspAccountClosrueService{coreRDB: coreData.DBRead, coreWDB: coreData.DBWrite}
}

///CSPClosureRequestCreate ...
func (s cspAccountClosrueService) CSPClosureRequestCreate(create CSPAccountClosureCreate) (CSPClosureRequestItem, error) {
	var closureRequestItem = CSPClosureRequestItem{}
	_, err := getByBusinessID(create.BusinessID)
	if err == nil {
		return closureRequestItem, services.ErrorNotFound{}.New("A request for account closure is already in progress")
	}

	agentID, err := cspuser.NewUserService(s.sr).ByCognitoID(s.sr.CognitoID)
	if err == sql.ErrNoRows {
		log.Printf("no request %v", err)
		return closureRequestItem, services.ErrorNotFound{}.New("CSP agent not found")
	}

	_, err = s.coreWDB.Exec(`INSERT INTO account_closure_request(business_id, reason, description) 
	VALUES($1, $2, $3)`, create.BusinessID, create.Reason, create.Description)
	if err != nil {
		return closureRequestItem, err
	}

	reqItem, err := getByBusinessID(create.BusinessID)
	closureRequestItem, err = s.CSPClosureRequestDetails(reqItem.ID)

	cspAgent, _ := cspuser.NewUserService(cspsrv.NewSourceRequest()).GetByIdInternal(agentID)
	desc := fmt.Sprintf("By %v %v", cspAgent.FirstName, cspAgent.LastName)
	stBody := CSPClosureStatePostBody{ACStateRequestCreated, closureRequestItem.ID, &agentID, &desc}
	s.CSPClosureStateAddNew(stBody)
	return closureRequestItem, err
}

//CSPClosureRequestList ...
func (s cspAccountClosrueService) CSPClosureRequestList(params CSPAccountClosureQueryParams) ([]CSPClosureRequestItem, error) {
	var closureRequests []CSPClosureRequestItem
	q := ""

	limit := "20"
	offset := "0"

	if params.Limit != "" {
		limit = params.Limit
	}

	if params.Offset != "" {
		offset = params.Offset
	}

	pagination := "LIMIT " + limit + " OFFSET " + offset

	sortBy := "ac.created"
	sortDirection := "DESC"

	switch params.SortDirection {
	case AccountClosureSortDescending:
		sortDirection = "DESC"
	case AccountClosureSortAscending:
		sortDirection = "ASC"
	default:
		sortDirection = "desc"
	}

	switch params.SortField {
	case AccountClosureFieldBusinessName:
		sortBy = "b.legal_name"
	case AccountClosureFieldCreated:
		sortBy = "ac.created"
	case AccountClosureFieldClosed:
		sortBy = "ac.closed"
	case AccountClosureFieldRefundAmount:
		sortBy = "ac.refund_amount"
	case AccountClosureFieldDigitalCheckNumber:
		sortBy = "ac.digital_check_number"
	case AccountClosureFieldOwnerName:
		sortBy = "co.first_name"
	case AccountClosureFieldStatus:
		sortBy = "ac.status"
	/*
		Fetching of Avialable balance and Posted balance moved to service-transaction
		case AccountClosureFieldAvailableBalance:
			sortBy = "bb.available_balance"
		case AccountClosureFieldPostedBalance:
			sortBy = "bb.posted_balance"
	*/
	default:
		sortBy = "ac.created"
	}

	sort := sortBy + " " + sortDirection
	clause, err := getWhereClauseForRequestListing(params)
	if err != nil {
		return closureRequests, err
	}

	qItems := tableJoinsAndColumnsForClosureRequest()
	q = fmt.Sprintf("SELECT %v FROM account_closure_request ac %v WHERE %v ORDER BY %v %v;", qItems["columns"], qItems["tableJoin"], clause, sort, pagination)

	err = s.coreRDB.Select(&closureRequests, q)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no request %v", err)
			return closureRequests, nil
		}
		return nil, err
	}

	for index, req := range closureRequests {
		if req.CSPAgentID != nil {
			cspusr, err := cspuser.NewUserService(cspsrv.NewSourceRequest()).GetByIdInternal(*req.CSPAgentID)
			if err == nil {
				name := cspusr.FirstName
				closureRequests[index].CSPUserName = &name
			}
		}
		avBal, psBal, err := fetchAccountBalance(req)
		if err != nil {
			fmt.Println("Available balance fetch err: ", err)
		}
		closureRequests[index].AvailableBalance = avBal
		closureRequests[index].PostedBalance = psBal
	}

	return closureRequests, err
}

func (s cspAccountClosrueService) CSPClosureApprovedAndRetryRequestList() ([]CSPClosureRequestItem, error) {
	var closureRequests []CSPClosureRequestItem
	q := fmt.Sprintf("SELECT * FROM account_closure_request WHERE status IN ('%v', '%v')", AccountClosureRequestApproved, AccountClosureFailedRetry)
	err := s.coreRDB.Select(&closureRequests, q)

	if err == sql.ErrNoRows {
		log.Printf("no request %v", err)
		return closureRequests, services.ErrorNotFound{}.New("")
	}

	return closureRequests, err
}

func (s cspAccountClosrueService) CSPClosureRequestDetails(id string) (CSPClosureRequestItem, error) {
	var closureRequest CSPClosureRequestItem
	qItems := tableJoinsAndColumnsForClosureRequest()
	q := fmt.Sprintf("SELECT %v FROM account_closure_request ac %v WHERE ac.id = '%v'", qItems["columns"], qItems["tableJoin"], id)
	err := s.coreRDB.Get(&closureRequest, q)

	if err == sql.ErrNoRows {
		log.Printf("no request %v", err)
		return closureRequest, services.ErrorNotFound{}.New("")
	}

	req := services.NewSourceRequest()
	req.UserID = closureRequest.OwnerUserID
	closureRequest.Business, err = bus.NewBusinessService(req).GetById(closureRequest.BusinessID)
	closureRequest.Owner, err = user.NewConsumerServiceWithout().GetByID(closureRequest.OwnerConsumerID)

	avBal, psBal, err := fetchAccountBalance(closureRequest)
	closureRequest.AvailableBalance = avBal
	closureRequest.PostedBalance = psBal

	return closureRequest, err
}

func (s cspAccountClosrueService) CSPClosureRequestUpdate(id string, status string) (*CSPClosureRequestItem, error) {
	var requestItem CSPClosureRequestItem
	state := ""
	if AccountClosureStatus(status) == AccountClosureRequestApproved {
		state = string(ACStateRequestApproved)
	} else if AccountClosureStatus(status) == AccountClosureRequestCanceled {
		state = string(ACStateRequestCanceled)
	} else if AccountClosureStatus(status) == AccountClosureFailedRetry {
		state = string(ACStateRequestRetryRequested)
	}

	fmt.Println("New status: ", status)

	agentID, err := cspuser.NewUserService(s.sr).ByCognitoID(s.sr.CognitoID)
	if err == sql.ErrNoRows && state != "" {
		log.Printf("no request %v", err)
		return &requestItem, services.ErrorNotFound{}.New("CSP User not found")
	}

	cspAgentValue := ""
	if agentID != "" {
		cspAgentValue = fmt.Sprintf(", csp_agent_id = '%v'", agentID)
	}

	closedDate := ""
	statusObj := AccountClosureStatus(status)
	if statusObj == AccountClosureRequestClosed || statusObj == AccountClosureRefundPending {
		closedDate = ",closed = CURRENT_TIMESTAMP"
	}

	q := fmt.Sprintf("UPDATE account_closure_request SET status = '%v' %v %v WHERE id = '%v'", status, cspAgentValue, closedDate, id)

	_, err = s.coreWDB.Exec(q)
	if err != nil {
		fmt.Println("Status update error: ", err)
		return nil, err
	}

	if state != "" {
		srcReq := cspsrv.NewSourceRequest()
		cspAgent, _ := cspuser.NewUserService(srcReq).GetByIdInternal(agentID)
		desc := fmt.Sprintf("By %v %v", cspAgent.FirstName, cspAgent.LastName)
		stBody := CSPClosureStatePostBody{AccountClosureState(state), id, nil, &desc}
		s.CSPClosureStateAddNew(stBody)
	}

	requestItem, err = s.CSPClosureRequestDetails(id)
	return &requestItem, err
}

func (s cspAccountClosrueService) CSPClosureRequestProcessed(id string, amount num.Decimal, status AccountClosureStatus) (*CSPClosureRequestItem, error) {
	var requestItem CSPClosureRequestItem

	closedDate := ""
	statusObj := AccountClosureStatus(status)
	if statusObj == AccountClosureRequestClosed || statusObj == AccountClosureRefundPending {
		closedDate = ",closed = CURRENT_TIMESTAMP"
	}

	q := fmt.Sprintf("UPDATE account_closure_request SET status = '%v', refund_amount = '%v' %v WHERE id = '%v'", status, amount, closedDate, id)

	_, err := s.coreWDB.Exec(q)
	if err != nil {
		return nil, err
	}
	requestItem, err = s.CSPClosureRequestDetails(id)
	return &requestItem, err
}

//CSPClosureRequestList ...
func (s cspAccountClosrueService) CSPClosureStateList(closureRequestID string) ([]CSPClosureState, error) {
	var closureStates []CSPClosureState
	q := fmt.Sprintf("SELECT * FROM account_closure_state WHERE account_closure_request_id = '%v' ORDER BY created ASC", closureRequestID)

	err := s.coreRDB.Select(&closureStates, q)

	if err == sql.ErrNoRows {
		log.Printf("no states %v", err)
		return closureStates, services.ErrorNotFound{}.New("")
	}

	return closureStates, err
}

func (s cspAccountClosrueService) CSPClosureStateAddNew(create CSPClosureStatePostBody) (CSPClosureState, error) {
	var item CSPClosureState
	keys := services.SQLGenInsertKeys(create)
	values := services.SQLGenInsertValues(create)

	q := fmt.Sprintf("INSERT INTO account_closure_state (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := s.coreWDB.PrepareNamed(q)
	if err != nil {
		fmt.Println("Add new closure state error: ", err)
		return item, err
	}
	err = stmt.Get(&item, create)
	return item, err
}

//HELPERS
func tableJoinsAndColumnsForClosureRequest() map[string]string {

	columns := "ac.*, b.id as business_id, b.legal_name as business_name, usr.id as owner_user_id, co.id as owner_consumer_id, co.first_name as owner_first_name, co.middle_name as owner_middle_name, co.last_name as owner_last_name"

	joinBusiness := "LEFT JOIN business b ON ac.business_id = b.id"
	joinUser := "LEFT JOIN wise_user usr ON b.owner_id = usr.id"
	joinConsumer := "LEFT JOIN consumer co ON usr.consumer_id = co.id"

	tableJoin := fmt.Sprintf("%v %v  %v", joinBusiness, joinUser, joinConsumer)

	queryItems := map[string]string{
		"columns":   columns,
		"tableJoin": tableJoin,
	}
	return queryItems
}

func getWhereClauseForRequestListing(params CSPAccountClosureQueryParams) (string, error) {
	clause := ""

	if params.Status == AccountClosureListRequestPending {
		clause = "ac.status = 'pending'"
	} else if params.Status == AccountClosureListRequestClosed {
		clause = "ac.status in ('approved','account_closed','canceled', 'refund_pending')"
	} else if params.Status == AccountClosureListRequestFailed {
		clause = fmt.Sprintf("ac.status IN ('%v', '%v')", AccountClosureFailed, AccountClosureFailedRetry)
	} else {
		log.Printf("Invalid status %v", params.Status)
		return "", services.ErrorNotFound{}.New("Invalid status")

	}

	if params.StartDate != "" {
		clause += " AND ac.created >= '" + params.StartDate + "'"
	}

	if params.EndDate != "" {
		clause += " AND ac.created <= '" + params.EndDate + "'"
	}

	if params.BusinessID != "" {
		bID, err := shared.ParseBusinessID(params.BusinessID)
		if err == nil {
			clause += " AND ac.business_id = '" + string(bID) + "'"
		} else {
			return "", fmt.Errorf("Invalid business ID")
		}
	}
	if params.BusinessName != "" {
		bName := strings.ToLower(params.BusinessName)
		clause += " AND (LOWER(b.legal_name) LIKE '" + bName + "%' OR LOWER(b.legal_name) LIKE '% " + bName + "%')"
	}
	if params.OwnerName != "" {
		fv1 := strings.ToLower(params.OwnerName)
		clause += " AND (LOWER(co.first_name) LIKE '" + fv1 + "%' OR LOWER(co.middle_name) LIKE '" + fv1 + "%' OR LOWER(co.last_name) LIKE '" + fv1 + "%')"
	}
	/*
		Fetching of Avialable balance and Posted balance moved to service-transaction
		if params.AvailableBalanceMin != "" {
			clause += " AND bb.available_balance >= '" + params.AvailableBalanceMin + "'"
		}
		if params.AvailableBalanceMax != "" {
			clause += " AND bb.available_balance <= '" + params.AvailableBalanceMax + "'"
		}

		if params.PostedBalanceMin != "" {
			clause += " AND bb.posted_balance >= '" + params.PostedBalanceMin + "'"
		}
		if params.PostedBalanceMax != "" {
			clause += " AND bb.posted_balance <= '" + params.PostedBalanceMax + "'"
		}
	*/

	return clause, nil

}

func getByBusinessID(ID shared.BusinessID) (*CSPClosureRequestItem, error) {
	var item = CSPClosureRequestItem{}
	db := coreData.DBRead
	err := db.Get(&item, "SELECT * FROM account_closure_request WHERE business_id = $1 AND status IN ('pending', 'approved')", ID)
	return &item, err
}

func fetchAccountBalance(cr CSPClosureRequestItem) (*num.Decimal, *num.Decimal, error) {
	accs, err := business.NewAccountService().ListInternalByBusiness(cr.BusinessID, 20, 0)
	if err != nil {
		return nil, nil, err
	}

	totalAvBal := num.NewZero() // avaialble balance
	totalPsBal := num.NewZero() // posted balance

	for _, acc := range accs {
		avBal, err := num.NewFromFloat(acc.AvailableBalance)
		if err != nil {
			return nil, nil, err
		}

		psBal, err := num.NewFromFloat(acc.PostedBalance)
		if err != nil {
			return nil, nil, err
		}

		totalAvBal = totalAvBal.Add(avBal)
		totalPsBal = totalPsBal.Add(psBal)
	}
	return &totalAvBal, &totalPsBal, nil
}
