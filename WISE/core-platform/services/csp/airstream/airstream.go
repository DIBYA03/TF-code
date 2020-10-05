package airstream

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/wiseco/core-platform/services"
	bsrv "github.com/wiseco/core-platform/services/business"
	cspBusiness "github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/services/csp/consumer"
	"github.com/wiseco/core-platform/services/user"
	usrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

// CSPAirstreamService ...
type CSPAirstreamService interface {
	StartKYC(shared.ConsumerID) error
	StartKYB(shared.BusinessID) error
}

type cspAirstreamService struct {
}

//NewService ...
func NewService() CSPAirstreamService {
	return cspAirstreamService{}
}

func (s cspAirstreamService) StartKYC(cID shared.ConsumerID) error {
	// CLEAR KYC
	clearKycErr := runClearKYC(cID)
	if clearKycErr != nil {
		fmt.Println("CLEAR-KYC failed: ", clearKycErr)
		return clearKycErr
	}

	// PHONE KYC
	phoneKycErr := runPhoneKYC(cID)
	if phoneKycErr != nil {
		fmt.Println("PHONE-KYC failed: ", phoneKycErr)
		return phoneKycErr
	}

	// BBVA KYV
	bbvaKycErr := runBbvaKyc(cID)
	if bbvaKycErr != nil {
		fmt.Println("BBVA-KYC failed: ", bbvaKycErr)
		return bbvaKycErr
	}

	// ALLOY KYC
	alloyKycErr := runAlloyKYC(cID)
	if alloyKycErr != nil {
		fmt.Println("ALLOY-KYC failed: ", alloyKycErr)
		return alloyKycErr
	}

	return nil
}

func (s cspAirstreamService) StartKYB(businessID shared.BusinessID) error {
	business, err := bsrv.New().GetByIdInternal(businessID)
	if err != nil {
		fmt.Printf("KYB Failed, Unable to fetch business, BusinessID: %v, Err: %v\n", businessID, err)
		return err
	}

	// Step 1. Check for KYC for all Members
	err = checkMembersKYC(businessID)
	if err != nil {
		fmt.Printf("KYB Failed, BusinessID: %v, Err: %v\n", businessID, err)
		return err
	}

	// Step 2. Check for Business Origin date
	if business.OriginDate != nil {
		minMonths := 6
		now := time.Now()
		diff := now.Sub(business.OriginDate.Time())
		months := diff.Hours() / 24 / 30

		if int(months) < minMonths {
			return fmt.Errorf("Business is not more than %v months old, Total duration of business %v", minMonths, months)
		}
	}

	// Step 3. Run CLEAR KYB
	err = runClearKYB(businessID)
	if err != nil {
		fmt.Printf("KYB Failed, BusinessID: %v, Err: %v\n", businessID, err)
		return err
	}

	// Step 4. Create a conditionally approved BBVA Bank account

	return nil
}

// KYC

func runClearKYC(cID shared.ConsumerID) error {
	matchingScore := 0.0
	totalScore := 0.0

	clearResponse, err := consumer.New().RunClearKYC(cID)
	clearObj := make(map[string]interface{})
	if err != nil {
		fmt.Println("")
	}

	err = json.Unmarshal([]byte(clearResponse), &clearObj)

	riskInfo := clearObj["RiskInformPersonSearchResponse"]

	if riskInfo != nil {
		riskInfoMap := riskInfo.(map[string]interface{})
		mScore, err := strconv.ParseFloat(riskInfoMap["MatchingScore"].(string), 0)
		if err != nil {
			fmt.Println("Clear KYC: Unable to parse matching score: ", err)
			return err
		}
		matchingScore = mScore

		riskInfoResult := riskInfoMap["RiskInformPersonSearchResult"]

		if riskInfoResult != nil {
			riskInfoResultMap := riskInfoResult.(map[string]interface{})
			personEntity := riskInfoResultMap["PersonEntity"]
			if personEntity != nil {
				personEntityMap := personEntity.(map[string]interface{})
				tScore, err := strconv.ParseFloat(personEntityMap["TotalScore"].(string), 0)
				if err != nil {
					fmt.Println("Clear KYC: Unable to parse total score: ", err)
					return err
				}
				totalScore = tScore
			}
		}

	}

	if matchingScore >= 75.0 && totalScore <= 70.0 {
		fmt.Printf("Clear KYC Approved for consumer: %v, with Matching score: %f, Total score %f\n", cID, matchingScore, totalScore)
		return nil
	}
	err = fmt.Errorf("Clear KYC Failed for consumer: %v, with Matching score: %f, Total score %f", cID, matchingScore, totalScore)
	return err
}

func runPhoneKYC(cID shared.ConsumerID) error {
	fmt.Println("Run Phone Verification")
	phoneResponse, err := consumer.New().PhoneVerification(cID)

	if err != nil {
		fmt.Println("PHONE Verification Error: ", err)
		return err
	}

	phoneObj := make(map[string]interface{})

	err = json.Unmarshal([]byte(phoneResponse), &phoneObj)
	respType := ""
	if phoneObj != nil {
		respType = phoneObj["type"].(string)
		if respType == "person" || respType == "business" {
			fmt.Printf("Phone verification approved, consumer:%v, Phone type: %v\n", cID, respType)
			return nil
		}
		return fmt.Errorf("Phone verification failure, consumer:%v, Phone type: %v", cID, respType)
	}
	return fmt.Errorf("Phone verification failure, consumer:%v, Error: %v", cID, err)
}

func runBbvaKyc(cID shared.ConsumerID) error {
	c, err := usrv.NewConsumerServiceWithout().GetByID(cID)

	var resp *user.ConsumerKYCResponse
	switch c.KYCStatus {
	case services.KYCStatusNotStarted:
		return fmt.Errorf("BBVA KYC: ConsumerID: %v, Not submitted", cID)
	case services.KYCStatusReview:
		return fmt.Errorf("BBVA KYC: ConsumerID: %v, Already in review", cID)
	case services.KYCStatusApproved, services.KYCStatusDeclined:
		return fmt.Errorf("BBVA KYC: ConsumerID: %v, Already in approved or declined", cID)
	}

	resp, err = user.NewConsumerService(services.NewSourceRequest()).StartVerification(cID, true)
	if err != nil {
		cerr, ok := err.(*user.ConsumerKYCError)
		if ok {
			return fmt.Errorf("BBVA KYC: ConsumerID: %v, Failed: %v", cID, cerr)
		}

		return fmt.Errorf("BBVA KYC: ConsumerID: %v, Failed: %v", cID, err)
	}

	if resp.Status != services.KYCStatusApproved {
		return fmt.Errorf("BBVA KYC: AUTOMATION FAILURE ConsumerID: %v, KycStatus: %v", cID, resp.Status)
	}

	fmt.Printf("BBVA KYC: SUCCESS-  ConsumerID: %v, KycStatus: %v\n", cID, resp.Status)
	return nil
}

func runAlloyKYC(cID shared.ConsumerID) error {
	alloyResponse, err := consumer.New().RunAlloyKYC(cID)
	if err != nil {
		fmt.Println("")
	}

	alloyObj := make(map[string]interface{})
	err = json.Unmarshal([]byte(alloyResponse), &alloyObj)

	summary := alloyObj["summary"]
	if summary != nil {
		summaryObj := summary.(map[string]interface{})
		outcome := summaryObj["outcome"].(string)
		if outcome == "Approved" {
			fmt.Printf("ALloy KYC Approved for consumer: %v, with Outcome: %v\n", cID, outcome)
			return nil
		}
	}
	err = fmt.Errorf("ALloy KYC Failed for consumer: %v", cID)
	return err
}

//KYB

func checkMembersKYC(businessID shared.BusinessID) error {
	members, err := bsrv.NewMemberServiceWithout().ListInternal(20, 0, businessID)
	if err != nil {
		return err
	}

	for _, mem := range members {
		if mem.KYCStatus != services.KYCStatusApproved {
			return fmt.Errorf("Member not approved, ID: %v, KYC: %v", mem.ID, mem.KYCStatus)
		}
	}

	return nil
}

func runClearKYB(businessID shared.BusinessID) error {
	clearResponse, err := cspBusiness.New(services.NewSourceRequest()).RunClearVerification(businessID)
	if err != nil {
		return fmt.Errorf("CLEAR KYB failed, err: %v", err)
	}

	clearObj := make(map[string]interface{})
	err = json.Unmarshal([]byte(clearResponse), &clearObj)

	matchingScore := 0.0
	totalScore := 0.0

	riskInfo := clearObj["RiskInformBusinessSearchResponse"]

	if riskInfo != nil {
		riskInfoMap := riskInfo.(map[string]interface{})
		mScore, err := strconv.ParseFloat(riskInfoMap["MatchingScore"].(string), 0)
		if err != nil {
			fmt.Println("Clear KYC: Unable to parse matching score: ", err)
			return err
		}
		matchingScore = mScore

		riskInfoResult := riskInfoMap["RiskInformBusinessSearchResult"]

		if riskInfoResult != nil {
			riskInfoResultMap := riskInfoResult.(map[string]interface{})
			businessEntity := riskInfoResultMap["BusinessEntity"]
			if businessEntity != nil {
				businessEntityMap := businessEntity.(map[string]interface{})
				tScore, err := strconv.ParseFloat(businessEntityMap["TotalScore"].(string), 0)
				if err != nil {
					fmt.Println("Clear KYC: Unable to parse total score: ", err)
					return err
				}
				totalScore = tScore
			}
		}
		if matchingScore >= 75.0 && totalScore <= 70.0 {
			fmt.Printf("Clear KYB Approved for Business: %v, with Matching score: %f, Total score %f\n", businessID, matchingScore, totalScore)
			return nil
		}
		err = fmt.Errorf("Clear KYB Failed for Business: %v, with Matching score: %f, Total score %f", businessID, matchingScore, totalScore)
		return err

	}

	return nil
}
