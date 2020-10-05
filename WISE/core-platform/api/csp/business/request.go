package business

import (
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/shared"
	goLibClear "github.com/wiseco/go-lib/clear"
)

//Request Get all the active businesses
func Request(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	id := r.GetPathParam("businessId")
	if id != "" {
		switch method {
		case http.MethodGet:
			return requestBusinessID(r, id)
		case http.MethodPatch:
			return handleBusinessUpdate(id, r)
		}
	}

	q := r.GetQueryParam("status")
	limit, _ := r.GetQueryIntParamWithDefault("limit", 30)
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)

	if q == "" && method == http.MethodGet {
		return handleBusinessList("", limit, offset, r)
	}

	kycquery, ok := kycQuery.new(q)
	if !ok {
		return api.BadRequest(r, errors.New("Invalid query params"))
	}

	switch method {
	case http.MethodGet:
		return handleBusinessList(kycquery, limit, offset, r)
	default:
		return api.NotSupported(r)
	}
}

func requestBusinessID(r api.APIRequest, id string) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	switch method {
	case http.MethodGet:
		return handleBusinessID(id, r)
	default:
		return api.NotSupported(r)
	}
}

//DocumentRequest ..
func DocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	documentID := r.GetPathParam("documentId")
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)
	offet, _ := r.GetQueryIntParamWithDefault("offset", 0)

	if businessID != "" && documentID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodGet:
			return handleDocumentID(businessID, documentID, r)
		case http.MethodPatch:
			return handleDocumentUpdate(businessID, documentID, r)
		case http.MethodDelete:
			return handleDocumentDelete(businessID, documentID, r)
		}
	}

	if businessID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodPost:
			return handleDocumentCreate(businessID, r)
		case http.MethodGet:
			return handleDocumentList(businessID, limit, offet, r)
		}
	}

	return api.NotSupported(r)
}

// DocumentURLRequest document url
func DocumentURLRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	businessID := r.GetPathParam("businessId")
	documentID := r.GetPathParam("documentId")
	if businessID != "" && documentID != "" {
		return handleDocumentURL(businessID, documentID, r)
	}
	return api.NotSupported(r)
}

//MemberRequest ..
func MemberRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	memberID := r.GetPathParam("memberId")
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)

	if businessID != "" {
		if memberID != "" {
			switch strings.ToUpper(r.HTTPMethod) {
			case http.MethodGet:
				return handleMememberID(memberID, businessID, r)
			case http.MethodPatch:
				return handleUpdateMemember(memberID, businessID, r)
			default:
				return api.NotSupported(r)
			}
		} else {
			switch strings.ToUpper(r.HTTPMethod) {
			case http.MethodGet:
				return handleMemberList(businessID, limit, offset, r)
			case http.MethodPost:
				return handleCreateMember(businessID, r)
			default:
				return api.NotSupported(r)
			}
		}
	}

	return api.NotSupported(r)
}

// MemberVerificationRequest ..
func MemberVerificationRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	memberID := r.GetPathParam("memberId")
	if businessID == "" || memberID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodGet:
		return handleGetMemberVerification(memberID, businessID, r)
	case http.MethodPost:
		return handleMemberVerification(memberID, businessID, r)
	default:
		return api.NotSupported(r)
	}
}

// BankAccountRequest ..
func BankAccountRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleAccountCreation(businessID, r)
	case http.MethodGet:
		return handleBusinessAccount(businessID, r)
	default:
		return api.NotSupported(r)
	}
}

// ExternalBankAccountRequest ..
func ExternalBankAccountRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodGet:
		return handleExternalBusinessAccount(businessID, r)
	default:
		return api.NotSupported(r)
	}
}

// CardsRequest ..
func CardsRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	cardID := r.GetPathParam("cardId")
	if cardID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodGet:
			return handleCardByID(businessID, cardID, r)
		default:
			return api.NotSupported(r)
		}
	}

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleCardCreation(businessID, r)
	case http.MethodGet:
		return handleCardList(businessID, r)
	default:
		return api.NotSupported(r)
	}
}

// SubmitDocumentRequest ..
func SubmitDocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodPost {
		return api.NotSupported(r)
	}
	businessID := r.GetPathParam("businessId")
	documentID := r.GetPathParam("documentId")
	if businessID != "" && documentID != "" {
		return handleSubmitDocument(businessID, documentID, r)
	}

	return api.NotSupported(r)
}

//StatusRequest ...
func StatusRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	if businessID := r.GetPathParam("businessId"); businessID != "" {
		return handleBusinessStatus(businessID, r)
	}
	return api.NotSupported(r)
}

//SetFormationRequest ..
func SetFormationRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	documentID := r.GetPathParam("documentId")
	if businessID != "" || documentID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodPost:
			return handleSetAsFormation(businessID, documentID, r)
		case http.MethodDelete:
			return handleRemoveFormation(businessID, documentID, r)
		default:
			return api.NotSupported(r)
		}
	}
	return api.NotSupported(r)
}

// DocumentStatusRequest ..
func DocumentStatusRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	businessID := r.GetPathParam("businessId")
	documentID := r.GetPathParam("documentId")
	if businessID != "" && documentID != "" {
		return handleDocumentStatus(businessID, documentID, r)
	}

	return api.NotSupported(r)
}

// ItemRequest ..
func ItemRequest(r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodGet:
		return handleCSPBusinessByBusinessID(bID, r)
	case http.MethodPatch:
		return handleCSPBusinessUpdate(bID, r)
	default:
		return api.NotSupported(r)
	}
}

// NotesRequest Business review notes request
func NotesRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	businessID := r.GetPathParam("businessId")
	notesID := r.GetPathParam("noteId")
	limit, _ := r.GetQueryIntParamWithDefault("limit", 30)
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)

	if businessID != "" && notesID != "" {
		switch method {
		case http.MethodGet:
			return handleNoteByID(businessID, notesID, r)
		case http.MethodPatch:
			return handleNoteUpdate(businessID, notesID, r)
		case http.MethodDelete:
			return handleDeleteNote(businessID, notesID, r)
		default:
			return api.NotSupported(r)
		}
	}

	if businessID != "" {
		switch method {
		case http.MethodGet:
			return handleNotesList(businessID, limit, offset, r)
		case http.MethodPost:
			return handleCreateNotes(businessID, r)
		default:
			return api.NotSupported(r)
		}
	}
	return api.NotSupported(r)
}

// SubscriptionRequest
func SubscriptionRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	businessID := r.GetPathParam("businessId")

	if businessID != "" {
		switch method {
		case http.MethodGet:
			return handleSubscriptionByBusinessID(businessID, r)
		case http.MethodPatch:
			return handleSubscriptionUpdate(businessID, r)
		default:
			return api.NotSupported(r)
		}
	}
	return api.NotSupported(r)
}

// PromofundRequest handles funding promotional amount into the giving business
func PromofundRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodPost {
		return api.NotSupported(r)
	}
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}
	return handlePromotionalFund(businessID, r)
}

// ListStatesRequest handles the list of busines state
func ListStatesRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	if businessID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodGet:
		return handleState(businessID, r)
	}
	return api.NotSupported(r)
}

// MemberEmailVerificationRequest ..
func MemberEmailVerificationRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	memberID := r.GetPathParam("memberId")
	if businessID == "" || memberID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleEmailVerification(memberID, businessID, r)
	default:
		return api.NotSupported(r)
	}
}

// MemberPhoneVerificationRequest ..
func MemberPhoneVerificationRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	memberID := r.GetPathParam("memberId")
	if businessID == "" || memberID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handlePhoneVerification(memberID, businessID, r)
	default:
		return api.NotSupported(r)
	}
}

// MemberAlloyVerificationRequest ..
func MemberAlloyVerificationRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	memberID := r.GetPathParam("memberId")
	if businessID == "" || memberID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleAlloyVerification(memberID, businessID, r)
	case http.MethodGet:
		return handleGetAlloyVerification(memberID, businessID, r)
	default:
		return api.NotSupported(r)
	}
}

// MemberClearVerificationRequest -
func MemberClearVerificationRequest(r api.APIRequest, clearKycType goLibClear.ClearKycType) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	memberID := r.GetPathParam("memberId")
	if businessID == "" || memberID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleClearVerification(memberID, businessID, clearKycType, r)
	case http.MethodGet:
		return handleGetClearVerification(memberID, businessID, clearKycType, r)
	default:
		return api.NotSupported(r)
	}
}

// MiddeskVerificationRequest ..
func MiddeskVerificationRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	if businessID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleMidDeskVerification(businessID, r)
	case http.MethodGet:
		return handleMidDeskGetVerification(businessID, r)
	}
	return api.NotSupported(r)
}

// ClearVerificationRequest ..
func ClearVerificationRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	if businessID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleClearBusinessVerification(businessID, r)
	case http.MethodGet:
		return handleClearBusinessGetVerification(businessID, r)
	}
	return api.NotSupported(r)
}

// ReUploadDocumentRequest ..
func ReUploadDocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID := r.GetPathParam("businessId")
	if businessID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleReUploadDocument(businessID, r)
	default:
		return api.NotSupported(r)
	}
}

// CardReissueRequest ..
func CardReissueRequest(r api.APIRequest) (api.APIResponse, error) {
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	cardID := r.GetPathParam("cardId")
	if cardID == "" {
		return api.BadRequestError(r, err)
	}

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleCardReissue(businessID, cardID, r)
	default:
		return api.NotSupported(r)
	}
}
