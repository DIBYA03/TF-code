package business

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	mbsrv "github.com/wiseco/core-platform/services/business"
	csp "github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/shared"
	goLibClear "github.com/wiseco/go-lib/clear"
)

func handleMemberList(businessID string, limit, offset int, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	members, err := csp.NewMemberService(r.SourceRequest()).List(bID, limit, offset)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(members)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleMememberID(memberID, businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	member, err := csp.NewMemberService(r.SourceRequest()).GetByID(mID, bID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(member)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleCreateMember(businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	var body mbsrv.BusinessMemberCreate
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	member, err := csp.NewMemberService(r.SourceRequest()).Create(bID, body)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(member)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	member, err = csp.NewMemberService(r.SourceRequest()).Submit(member.ID, bID)

	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err = json.Marshal(member)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleUpdateMemember(memberID, businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	var body mbsrv.BusinessMemberUpdate
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	member, err := csp.NewMemberService(r.SourceRequest()).Update(mID, bID, body)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(member)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleMemberVerification(memberID, businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	verification, err := csp.NewMemberService(r.SourceRequest()).StartVerification(mID, bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(verification)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleGetMemberVerification(memberID, businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	verification, err := csp.NewMemberService(r.SourceRequest()).GetVerification(mID, bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(verification)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleEmailVerification(memberID, businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	verification, err := csp.NewMemberService(r.SourceRequest()).EmailVerification(mID, bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(verification)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handlePhoneVerification(memberID, businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}
	verification, err := csp.NewMemberService(r.SourceRequest()).PhoneVerification(mID, bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(verification)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleAlloyVerification(memberID, businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}
	verification, err := csp.NewMemberService(r.SourceRequest()).RunAlloyKYC(mID, bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(verification)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleGetAlloyVerification(memberID, businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}
	verification, err := csp.NewMemberService(r.SourceRequest()).GetAlloyKYC(mID, bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(verification)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleClearVerification(memberID, businessID string, kycType goLibClear.ClearKycType, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}
	verification, err := csp.NewMemberService(r.SourceRequest()).RunClearKYC(mID, bID, kycType)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, verification, false)
}

func handleGetClearVerification(memberID, businessID string, kycType goLibClear.ClearKycType, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	mID, err := shared.ParseBusinessMemberID(memberID)
	if err != nil {
		return api.BadRequest(r, err)
	}
	verification, err := csp.NewMemberService(r.SourceRequest()).GetClearKYC(mID, bID, kycType)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, verification, false)
}
