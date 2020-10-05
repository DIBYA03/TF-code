package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/wiseco/go-lib/log"
	"github.com/wiseco/go-lib/router"
	"github.com/wiseco/wise-auth-cli/internal/service/awsauthcli"
)

// AWSAuthCLIController is a new controller for aws-auth-cli requests
type AWSAuthCLIController interface {
	GETAWSAuthCLI(http.ResponseWriter, *http.Request)
	POSTAWSAuthCLI(http.ResponseWriter, *http.Request)
	PostLocalAWSAuthCLI(http.ResponseWriter, *http.Request)
	GetLocalAWSAuthCLI(http.ResponseWriter, *http.Request)
	GetLocalAWSConsoleSignIn(http.ResponseWriter, *http.Request)
}

type awsAuthCLIController struct{}

// NewAWSAuthCLIController creates a new controller for requests
func NewAWSAuthCLIController() AWSAuthCLIController {
	return &awsAuthCLIController{}
}

func (c awsAuthCLIController) GETAWSAuthCLI(w http.ResponseWriter, r *http.Request) {
	l := router.GetLogger(r)

	idpURI := os.Getenv("GOOGLE_IDP_URL")

	nameIDPolicy := &awsauthcli.SamlpNameIDPolicy{
		Format:      "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
		AllowCreate: "true",
	}

	requestedAuthnContext := &awsauthcli.SamlpRequestedAuthnContext{
		Comparison:           "exact",
		AuthnContextClassRef: "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport",
	}

	authRequest := &awsauthcli.AuthnRequest{
		XmlnsSamlP:                  "urn:oasis:names:tc:SAML:2.0:protocol",
		XmlnsSaml:                   "urn:oasis:names:tc:SAML:2.0:assertion",
		ID:                          "wise-auth-cli-login-auth-request",
		Version:                     "2.0",
		ProviderName:                "https://aws.us-west-2.internal.wise.us/saml",
		IssueInstant:                "2034-02-28T02:21:43.000Z",
		Destination:                 idpURI,
		ProtocolBinding:             "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
		AssertionConsumerServiceURL: "https://aws.us-west-2.internal.wise.us/saml",
		SamlIssuer:                  "https://aws.us-west-2.internal.wise.us/saml/metadata",
		SamlpNameIDPolicy:           nameIDPolicy,
		SamlpRequestedAuthnContext:  requestedAuthnContext,
	}

	googleSAMLRequest, err := xml.MarshalIndent(authRequest, " ", "	")
	if err != nil {
		l.ErrorD("error in xml marshalling", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b64SAMLRequest := base64.StdEncoding.EncodeToString(googleSAMLRequest)

	formData := url.Values{
		"SAMLRequest": {b64SAMLRequest},
	}

	resp, err := http.PostForm(idpURI, formData)
	if err != nil {
		l.ErrorD("error in posting saml request:", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// fmt.Println(resp)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	// body := buf.String()

	// add the headers from SAML request
	var autoLoginURLHeader string
	for headerName, headerVal := range resp.Header {
		if headerName == "X-Auto-Login" {
			autoLoginURLHeader = headerVal[0]
		}
		w.Header().Set(headerName, strings.Join(headerVal, " "))
	}

	// break out the auoLoginHeader
	autoLoginHeader, err := url.ParseQuery(autoLoginURLHeader)
	if err != nil {
		l.ErrorD("error in parsing autologin header", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// break out the continue url param
	autoLoginURLArgs, err := url.ParseQuery(autoLoginHeader["args"][0])
	if err != nil {
		l.ErrorD("error in parsing autologin header arguments", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// add the cookies from SAML request
	for _, cookie := range resp.Cookies() {
		http.SetCookie(w, cookie)
	}

	http.Redirect(w, r, autoLoginURLArgs["continue"][0], http.StatusSeeOther)
}

func (c awsAuthCLIController) POSTAWSAuthCLI(w http.ResponseWriter, r *http.Request) {
	// We want to redirect the request back to localhost(user's computer)
	// We do this, so all requests have to go through the VPN(this container)
	http.Redirect(w, r, "https://localhost:4433/aws", http.StatusTemporaryRedirect)
}

// GetLocalAWSConsoleSignIn is used to sign in to the AWS Console
func (c awsAuthCLIController) GetLocalAWSConsoleSignIn(w http.ResponseWriter, r *http.Request) {
	l := router.GetLogger(r)

	params := r.URL.Query()

	accessKeyID := params["AccessKeyId"]
	secretAccessKey := params["SecretAccessKey"]
	sessionToken := params["SessionToken"]
	sessionDuration := params["SessionDuration"]
	issuer := params["Issuer"]

	signInTokenRequest := awsauthcli.SignInTokenRequest{
		ID:    accessKeyID[0],
		Key:   secretAccessKey[0],
		Token: sessionToken[0],
	}

	sD, err := strconv.Atoi(sessionDuration[0])
	if err != nil {
		l.ErrorD("error processing session duration", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signInURL, err := awsauthcli.HandleSignInURL(signInTokenRequest, issuer[0], int64(sD))
	if err != nil {
		l.ErrorD("error in getting signin URL", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, signInURL, http.StatusTemporaryRedirect)

}

// GetLocalAWSAuthCLI handles redirecting to the internal VPN endpoint for starting SAML request
func (c awsAuthCLIController) GetLocalAWSAuthCLI(w http.ResponseWriter, r *http.Request) {
	wiseAuthRedirect := os.Getenv("WISE_INTERNAL_AUTH_REDIRECT")
	http.Redirect(w, r, wiseAuthRedirect, http.StatusTemporaryRedirect)
}

func (c awsAuthCLIController) PostLocalAWSAuthCLI(w http.ResponseWriter, r *http.Request) {
	l := router.GetLogger(r)
	var errors []string

	r.ParseForm()
	samlResponse := r.FormValue("SAMLResponse")

	decodedBody, err := base64.StdEncoding.DecodeString(samlResponse)
	if err != nil {
		l.ErrorD("error in decoding saml response", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseStruct := awsauthcli.Saml2pResponse{}
	err = xml.Unmarshal([]byte(decodedBody), &responseStruct)
	if err != nil {
		l.ErrorD("error in xml unmarshalling", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if responseStruct.Status.StatusCode.StatusCode.Value == "urn:oasis:names:tc:SAML:2.0:status:RequestDenied" {
		l.ErrorD("SAML request denied", log.Fields{"err": responseStruct.Status.StatusMessage})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var awsRoles []string
	var sessionDuration int64

	for _, attributes := range responseStruct.Assertion.AttributeStatement.Attributes {
		for _, attributeVal := range attributes.Attribute {
			switch attributes.Name {
			case "https://aws.amazon.com/SAML/Attributes/Role":
				awsRoles = append(awsRoles, attributeVal.Value)
				break
			case "https://aws.amazon.com/SAML/Attributes/SessionDuration":
				sessionDuration, err = strconv.ParseInt(attributeVal.Value, 10, 64)
				if err != nil {
					errors = append(errors, err.Error())
					l.ErrorD("error in getting session duration", log.Fields{"err": err.Error()})
					w.WriteHeader(http.StatusInternalServerError)
				}
				break
			}
		}
	}

	issuer := responseStruct.Saml2Issuer.Issuer

	awsProfiles, err := awsauthcli.HandleSessionCredentials(awsRoles, sessionDuration, samlResponse, awsauthcli.CredsFile)
	if err != nil {
		l.ErrorD("error in handling session credentials", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signInURLs, err := awsauthcli.HandlePreSignInURLs(awsProfiles, issuer, sessionDuration)
	if err != nil {
		l.ErrorD("error in handling signin urls", log.Fields{"err": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	awsauthcli.GenerateAWSTemplate(signInURLs, errors, awsProfiles[0].Expiration, w)
}
