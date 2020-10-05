package awsauthcli

import "encoding/xml"

// SamlpNameIDPolicy is the policy ID for auth request
type SamlpNameIDPolicy struct {
	XMLName     xml.Name `xml:"samlp:NameIDPolicy"`
	Format      string   `xml:"Format,attr"`
	AllowCreate string   `xml:"AllowCreate,attr"`
}

// SamlpRequestedAuthnContext is the context for the auth request
type SamlpRequestedAuthnContext struct {
	XMLName              xml.Name `xml:"samlp:RequestedAuthnContext"`
	Comparison           string   `xml:"Comparison,attr"`
	AuthnContextClassRef string   `xml:"saml:AuthnContextClassRef"`
}

// AuthnRequest is the initial saml auth request
type AuthnRequest struct {
	XMLName                     xml.Name `xml:"samlp:AuthnRequest"`
	XmlnsSamlP                  string   `xml:"xmlns:samlp,attr"`
	XmlnsSaml                   string   `xml:"xmlns:saml,attr"`
	ID                          string   `xml:"ID,attr"`
	Version                     string   `xml:"Version,attr"`
	ProviderName                string   `xml:"ProviderName,attr"`
	IssueInstant                string   `xml:"IssueInstant,attr"`
	Destination                 string   `xml:"Destination,attr"`
	ProtocolBinding             string   `xml:"ProtocolBinding,attr"`
	AssertionConsumerServiceURL string   `xml:"AssertionConsumerServiceURL,attr"`
	SamlIssuer                  string   `xml:"saml:Issuer"`
	SamlpNameIDPolicy           *SamlpNameIDPolicy
	SamlpRequestedAuthnContext  *SamlpRequestedAuthnContext
}
