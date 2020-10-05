package wise_user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awscog "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jmoiron/sqlx"
	"github.com/ttacon/libphonenumber"
	"github.com/wiseco/core-platform/identity"
	t "github.com/wiseco/core-platform/partner/service/twilio"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/consumer"
	"github.com/wiseco/core-platform/services/csp/cspuser"
	"github.com/wiseco/core-platform/services/csp/data"
	cspsrv "github.com/wiseco/core-platform/services/csp/services"
	coreDB "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type Service interface {
	ListPhoneChangeRequest(shared.UserID) ([]PhoneChangeRequest, error)
	ChangePhone(PhoneChangeRequestCreate) error
}

type service struct {
	sourceReq services.SourceRequest
	db        *sqlx.DB
	cognitoID *string
}

func NewWithCognitoSource(sourceReq services.SourceRequest, cognitoID *string) Service {
	return service{
		sourceReq: sourceReq,
		cognitoID: cognitoID,
		db:        data.DBWrite}
}

func (s service) ChangePhone(r PhoneChangeRequestCreate) error {
	user, err := consumer.New().ByUserID(string(r.UserID))
	if err != nil {
		return err
	}

	id, err := identity.NewIdentityServiceWithout().GetByID(shared.IdentityID(user.IdentityID))
	if err != nil {
		return err
	}

	if user.Phone != id.Phone {
		return errors.New("cognito and user phone number mismatch")
	}

	// Validate phone no default
	ph, err := libphonenumber.Parse(r.NewPhone, "")
	if err != nil {
		return err
	}

	r.NewPhone = libphonenumber.Format(ph, libphonenumber.E164)

	sess := session.Must(session.NewSession())
	srv := awscog.New(sess)

	poolID := os.Getenv("COGNITO_USER_POOL_ID")
	if poolID == "" {
		return errors.New("COGNITO_USER_POOL_ID is missing")
	}

	in := &awscog.AdminUpdateUserAttributesInput{
		UserPoolId: aws.String(poolID),
		UserAttributes: []*awscog.AttributeType{
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(r.NewPhone),
			},
			{
				Name:  aws.String("phone_number_verified"),
				Value: aws.String(strconv.FormatBool(true)),
			},
		},
		Username: aws.String(id.Phone),
	}

	_, err = srv.AdminUpdateUserAttributes(in)
	if err != nil {
		return err
	}

	idUpdate := identity.IdentityUpdate{
		ID:    id.ID,
		Phone: r.NewPhone,
	}
	err = identity.NewIdentityServiceWithout().Update(idUpdate)
	if err != nil {
		return err
	}

	err = s.updateUserPhone(user.ID, r.NewPhone)
	if err != nil {
		return err
	}

	userID, err := cspuser.NewUserService(cspsrv.NewSRRequest(*s.cognitoID)).ByCognitoID(*s.cognitoID)
	if err != nil {
		return err
	}

	create := PhoneChangeRequestCreate{
		UserID:            shared.UserID(user.ID),
		OldPhone:          id.Phone,
		NewPhone:          r.NewPhone,
		OriginatedFrom:    OriginatedFromCSP,
		VerificationNotes: r.VerificationNotes,
		CSPUserID:         userID,
	}

	log.Println("Cognito ID is ", userID, *s.cognitoID, create.UserID, user.ID, create.CSPUserID)

	// Default/mandatory fields
	columns := []string{
		"user_id", "old_phone", "new_phone", "originated_from", "csp_user_id", "verification_notes",
	}
	// Default/mandatory values
	values := []string{
		":user_id", ":old_phone", ":new_phone", ":originated_from", ":csp_user_id", ":verification_notes",
	}

	sql := fmt.Sprintf("INSERT INTO phone_change_request(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := s.db.PrepareNamed(sql)
	if err != nil {
		return err
	}

	request := &PhoneChangeRequest{}

	err = stmt.Get(request, &create)
	if err != nil {
		return err
	}

	userName := user.FirstName + " " + user.LastName
	smsReq := t.SMSRequest{
		Body:  fmt.Sprintf(services.PhoneNumberChangeSMS, userName),
		Phone: user.Phone,
	}

	err = t.NewTwilioService().SendSMS(smsReq)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (s service) ListPhoneChangeRequest(userID shared.UserID) ([]PhoneChangeRequest, error) {
	var list []PhoneChangeRequest

	err := s.db.Select(&list, `SELECT * FROM phone_change_request WHERE user_id = $1`, userID)
	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}

	if err != nil {
		return list, err
	}

	srv := cspuser.NewUserService(cspsrv.NewSRRequest(*s.cognitoID))

	for i, _ := range list {
		if list[i].CSPUserID != nil {
			usr, err := srv.GetByIdInternal(*list[i].CSPUserID)
			if err != nil {
				return list, err
			}

			cspUserName := usr.FirstName + " " + usr.LastName
			list[i].CSPUserName = &cspUserName
		}
	}

	return list, err
}

func (s service) updateUserPhone(ID string, phone string) error {
	var err error
	if phone != "" {
		_, err = coreDB.DBWrite.Exec("UPDATE wise_user SET phone = $1 WHERE id = $2", phone, ID)
	}
	if err != nil {
		log.Printf("error update user phone %v", err)
	}
	return err
}
