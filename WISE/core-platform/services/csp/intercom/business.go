package intercom

import (
	"github.com/wiseco/go-lib/intercom"
	"github.com/wiseco/go-lib/log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/data"
)

type Service interface {
	GetByEmailID(string) (*string, error)
	GetByUserID(string) (*string, error)
	GetTags() (*TagList, error)
	SetTags(tbs []TagBodyItem) (*TagList, error)
}

type service struct {
	*sqlx.DB
	sourceReq services.SourceRequest
}

//TagBodyItem for CSP API
type TagBodyItem struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	UserID string `json:"userId"`
}

//Tag ...
type Tag struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

//TagList ...
type TagList struct {
	Type string `json:"type"`
	Tag  []Tag  `json:"tags"`
}

//New will return a new csp business service
func New(r services.SourceRequest) Service {
	return service{data.DBWrite, r}
}

func (s service) GetByEmailID(emailID string) (*string, error) {

	l := log.NewLogger()
	resp, err := intercom.NewIntercomService(l).GetUserByEmail(emailID)

	return resp, err
}

func (s service) GetByUserID(userID string) (*string, error) {

	l := log.NewLogger()
	resp, err := intercom.NewIntercomService(l).GetUserByID(userID)

	return resp, err
}

func (s service) GetTags() (*TagList, error) {

	l := log.NewLogger()
	ts, err := intercom.NewIntercomService(l).GetTags()
	if err != nil {
		return nil, err
	}
	tl := TagList{}
	tl.Type = ts.Type
	for _, t := range ts.Tag {
		nt := Tag{Type: t.Type, ID: t.ID, Name: t.Name}
		tl.Tag = append(tl.Tag, nt)
	}
	return &tl, nil
}

func (s service) SetTags(tbs []TagBodyItem) (*TagList, error) {
	var tas []intercom.TagAction
	for _, tb := range tbs {
		action := intercom.TagActionName(tb.Action)
		ta := intercom.TagAction{Name: tb.Name, Action: action, UserID: tb.UserID}
		tas = append(tas, ta)
	}
	l := log.NewLogger()
	ts, err := intercom.NewIntercomService(l).SetTags(tas)

	tl := TagList{}
	tl.Type = ts.Type
	for _, t := range ts.Tag {
		nt := Tag{Type: t.Type, ID: t.ID, Name: t.Name}
		tl.Tag = append(tl.Tag, nt)
	}

	return &tl, err
}
