package auth

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/ttacon/libphonenumber"
)

// Cognito Input Keys
const (
	CognitoAddress             = "address"               // Address
	CognitoBirthdate           = "birthdate"             // Birth Date
	CognitoEmail               = "email"                 // Email Address
	CognitoEmailVerified       = "email_verified"        // Email Address Verified
	CognitoLastName            = "family_name"           // Last Name
	CognitoGender              = "gender"                // Gender
	CognitoFirstName           = "given_name"            // First Name
	CognitoLocale              = "locale"                // Locale
	CognitoMiddleName          = "middle_name"           // Middle Name
	CognitoName                = "name"                  // Full Name
	CognitoNickname            = "nickname"              // Nickname
	CognitoPhoneNumber         = "phone_number"          // Phone Number
	CognitoPhoneNumberVerified = "phone_number_verified" // Phone Number Verified
	CognitoPicture             = "picture"               // User's Picture
	CognitoUsername            = "preferred_username"    // Username
	CognitoProfile             = "profile"               // User Profile
	CognitoTimezone            = "timezone"              // Time Zone
	CognitoUpdatedAt           = "updated_at"            // Last Updated Timestamp
	CognitoWebsite             = "website"               // Website
	CognitoSub                 = "sub"                   // User Id
	CognitoUserStatus          = "cognito:user_status"   // Cognito User Status (Case Insensitive)
	CognitoStatus              = "status"                // Cognito Status/Enabled (Case Sensitive)
)

// Error Keys
const (
	PhoneMissing = "phoneMissing" // Phone missing
	PhoneInvalid = "phoneInvalid" // Phone invalid/malformed

	EmailInvalid = "emailInvalid" // Email invalid/malformed

	IdMissing   = "idMissing"   // Id is missing
	IdMalformed = "idMalformed" // Id is malformed

	UserAlreadyExists = "userAlreadyExists"
	UserPhoneExists   = "userPhoneExists" // User phone or email exists
)

func GetUserAttributeByKey(userAttributes map[string]string, key string) (*string, bool) {
	firstName, ok := userAttributes[key]
	if ok {
		return &firstName, true
	}

	return nil, false
}

func ValidatePhone(event events.CognitoEventUserPoolsPreSignup) []string {
	errorList := []string{}

	phone, ok := GetUserAttributeByKey(event.Request.UserAttributes, CognitoPhoneNumber)
	if !ok || phone == nil {
		return append(errorList, PhoneMissing)
	}

	_, err := libphonenumber.Parse(*phone, "US")
	if err != nil {
		errorList = append(errorList, PhoneInvalid)
	}

	return errorList
}
