package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type deviceDataStore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

//DeviceService is the device interface
//where push notifications are managed
type DeviceService interface {
	// Fetch
	GetByID(shared.UserDeviceID) (*PushRegistration, error)

	// Create or update operations
	RegisterPush(*PushRegistrationCreate) (*PushRegistration, error)

	//UregisterPush deletes the token with the given id
	UnregisterPush(string) (bool, error)

	//Logout logs out a device
	Logout(*DeviceLogout) error
}

//NewDeviceService return a new serivce, exposing the method allow by the
// interface
func NewDeviceService(r services.SourceRequest) DeviceService {
	return &deviceDataStore{r, data.DBWrite}
}

func (db *deviceDataStore) GetByID(id shared.UserDeviceID) (*PushRegistration, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserDeviceAccess(id)
	if err != nil {
		return nil, err
	}

	var prOut PushRegistration
	err = db.Get(&prOut, "SELECT * FROM user_device WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, services.ErrorNotFound{}.New("")
	}

	if err != nil {
		return nil, err
	}

	return &prOut, err
}

func (db *deviceDataStore) getByToken(userID shared.UserID, token string) (*PushRegistration, error) {
	var prOut PushRegistration
	if err := db.Get(&prOut, "SELECT * FROM user_device WHERE user_id = $1 AND token = $2", userID, token); err != nil {
		return nil, fmt.Errorf("Error getting the push object with userID:%s, and token:%s error:%v", userID, token, err)
	}

	return &prOut, nil
}

func (db *deviceDataStore) getRegistration(userID shared.UserID, token string) *PushRegistration {
	var pr PushRegistration
	if err := db.Get(&pr, "SELECT * FROM user_device WHERE user_id = $1 AND token = $2", userID, token); err != nil {
		log.Printf("Error getting the push registraion with userID:%s, and token:%s error:%v", userID, token, err)
		return nil
	}
	return &pr
}

// RegistersPush
// registers token with fcm - only ios for now
func (db *deviceDataStore) RegisterPush(pr *PushRegistrationCreate) (*PushRegistration, error) {
	if pr.DeviceType == DeviceTypeWeb {
		return db.webRegistration(pr)
	}

	if pr.TokenType != TokenTypeFCM {
		return nil, errors.New("Token Type Error: Only FCM is supported")
	}

	if pr.Token == "" {
		return nil, errors.New("Token can not be empty")
	}
	var err error
	switch pr.DeviceType {
	case DeviceTypeIOS:
	case DeviceTypeAndroid:
	default:
		return nil, fmt.Errorf("Device Type Error: %s is not supported", pr.DeviceType)
	}

	// Return if already exists
	if registration := db.getRegistration(db.sourceReq.UserID, pr.Token); registration != nil {
		return registration, nil
	}

	sql := `
		INSERT INTO user_device (user_id, device_type, token_type, token, device_key, language)
		VALUES (:user_id, :device_type, :token_type, :token, :device_key, :language)
		RETURNING *`

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}
	var pushOut PushRegistration
	pr.UserID = db.sourceReq.UserID
	err = stmt.Get(&pushOut, &pr)
	if err != nil {
		return nil, err
	}

	return &pushOut, nil
}

func (db *deviceDataStore) webRegistration(pr *PushRegistrationCreate) (*PushRegistration, error) {
	pr.UserID = db.sourceReq.UserID
	sql := `
		INSERT INTO user_device (user_id, device_type, token_type, token, device_key, language)
		VALUES (:user_id, :device_type, :token_type, :token, :device_key, :language)
		RETURNING *`

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}
	var pushOut PushRegistration
	err = stmt.Get(&pushOut, &pr)
	if err != nil {
		return nil, err
	}

	return &pushOut, nil
}

// UnregisterPush deletes token for the given id
func (db *deviceDataStore) UnregisterPush(token string) (bool, error) {
	pr, err := db.getByToken(db.sourceReq.UserID, token)
	if err != nil {
		return false, err
	}

	// Check access
	err = auth.NewAuthService(db.sourceReq).CheckUserDeviceAccess(pr.ID)
	if err != nil {
		return false, err
	}

	_, err = db.Exec(`DELETE FROM user_device WHERE token = $1`, token)
	return err != nil, err
}

func (db *deviceDataStore) unregisterDevice(deviceKey DeviceKey) error {
	// Assuming we only have one device registered
	_, err := db.Exec("DELETE FROM user_device where device_key = $1", deviceKey)
	return err
}

// Logout user from session
func (db *deviceDataStore) Logout(d *DeviceLogout) error {
	// Check access
	var id shared.UserDeviceID
	err := db.Get(&id, "SELECT id FROM user_device WHERE device_key = $1", d.DeviceKey)
	if err == sql.ErrNoRows {
		return services.ErrorNotFound{}.New("")
	}

	if err != nil {
		return err
	}

	err = auth.NewAuthService(db.sourceReq).CheckUserDeviceAccess(id)
	if err != nil {
		return err
	}

	// Unregister for push notification
	err = db.unregisterDevice(d.DeviceKey)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error uregistering token error:%v", err)
		return err
	}

	return nil
}
