package auth

import (
	"github.com/aws/aws-lambda-go/events"
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
	EmailMissing = "emailMissing" // Phone missing
	EmailInvalid = "emailInvalid" // Email invalid/malformed

	EmailVerifiedMissing = "emailVerifiedMissing" // Email verification missing

	IDMissing   = "idMissing"   // Id is missing
	IDMalformed = "idMalformed" // Id is malformed

	UserAlreadyExists = "userAlreadyExists"
	UserEmailExists   = "userEmailExists" // User email exists
)

// GetUserAttributeByKey gets an attribute by key from events atttibutes
func GetUserAttributeByKey(userAttributes map[string]string, key string) (*string, bool) {
	attr, ok := userAttributes[key]
	if ok {
		return &attr, true
	}

	return nil, false
}

// ValidateEmail validates that email exists in events
func ValidateEmail(event events.CognitoEventUserPoolsPreSignup) []string {
	errorList := []string{}

	email, ok := GetUserAttributeByKey(event.Request.UserAttributes, CognitoEmail)
	if !ok || email == nil {
		return append(errorList, EmailMissing)
	}

	return errorList
}
