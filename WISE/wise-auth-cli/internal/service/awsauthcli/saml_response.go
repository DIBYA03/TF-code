package awsauthcli

import "encoding/xml"

// SamlResponseStatus is the status of the saml response
type SamlResponseStatus struct {
	XMLName xml.Name `xml:"Status"`

	StatusCode struct {
		Value      string `xml:"Value,attr"`
		StatusCode struct {
			Value string `xml:"Value,attr"`
		} `xml:"StatusCode"`
	} `xml:"StatusCode"`

	StatusMessage string `xml:"StatusMessage"`
}

// SamlResponseSignatureSignedInfoReferenceTransform is the transform for list of transforms in saml response
type SamlResponseSignatureSignedInfoReferenceTransform struct {
	XMLName   xml.Name `xml:"Transform"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// SamlResponseSignatureSignedInfoReferenceTransforms is the transforms for the signed info reference of saml resposne
type SamlResponseSignatureSignedInfoReferenceTransforms struct {
	XMLName    xml.Name                                            `xml:"Transforms"`
	Transforms []SamlResponseSignatureSignedInfoReferenceTransform `xml:"Transform"`
}

// SamlResponseSignatureSignedInfoReference is the signed info reference from saml response
type SamlResponseSignatureSignedInfoReference struct {
	XMLName xml.Name `xml:"Reference"`
	URI     string   `xml:"URI,attr"`

	DigestMethod struct {
		DigestMethod string `xml:",chardata"`
		Algorithm    string `xml:"Algorithm,attr"`
	} `xml:"DigestMethod"`

	Transforms SamlResponseSignatureSignedInfoReferenceTransforms

	DigestValue string `xml:"DigestValue"`
}

// SamlResponseSignatureSignedInfo is the signed info in signature of saml response
type SamlResponseSignatureSignedInfo struct {
	XMLName xml.Name `xml:"SignedInfo"`

	CanonicalizationMethod struct {
		CanonicalizationMethod string `xml:",chardata"`
		Algorithm              string `xml:"Algorithm,attr"`
	} `xml:"CanonicalizationMethod"`

	SignatureMethod struct {
		SignatureMethod string `xml:",chardata"`
		Algorithm       string `xml:"Algorithm,attr"`
	} `xml:"SignatureMethod"`

	Reference SamlResponseSignatureSignedInfoReference
}

// SamlResponseSignatureKeyInfo is the cert info for the saml response
type SamlResponseSignatureKeyInfo struct {
	XMLName xml.Name `xml:"KeyInfo"`

	X509DataSubjectName string `xml:"X509Data>X509SubjectName"`
	X509Certificate     string `xml:"X509Data>X509Certificate"`
}

// SamlResponseSignature is the saml response signature
type SamlResponseSignature struct {
	XMLName xml.Name `xml:"Signature"`
	XmlnsDS string   `xml:"xmlns:ds,attr"`

	SignedInfo     SamlResponseSignatureSignedInfo
	SignatureValue string `xml:"SignatureValue"`

	KeyInfo SamlResponseSignatureKeyInfo
}

// SamlResponseSubjectConfirmationData from saml response
type SamlResponseSubjectConfirmationData struct {
	XMLName      xml.Name `xml:"SubjectConfirmationData"`
	InResponseTo string   `xml:"InResponseTo,attr"`
	NotOnOrAfter string   `xml:"NotOnOrAfter,attr"`
	Recipient    string   `xml:"Recipient,attr"`
}

// SamlResponseSubjectConfirmationMethod is confirmation data from saml response
type SamlResponseSubjectConfirmationMethod struct {
	XMLName xml.Name `xml:"SubjectConfirmation"`
	Method  string   `xml:"Method,attr"`

	SubjectConfirmationData SamlResponseSubjectConfirmationData
}

// SamlResponseSubject is the subject of the assertion from the saml response
type SamlResponseSubject struct {
	XMLName xml.Name `xml:"Subject"`

	NameID struct {
		NameID string `xml:",chardata"`
		Format string `xml:"Format,attr"`
	} `xml:"NameID"`

	SubjectConfirmation SamlResponseSubjectConfirmationMethod
}

// SamlResponseConditions is the conditions/ rstrictions for saml response
type SamlResponseConditions struct {
	XMLName      xml.Name `xml:"Conditions"`
	NotBefore    string   `xml:"NotBefore,attr"`
	NotOnOrAfter string   `xml:"NotOnOrAfter,attr"`

	AudienceRestrictionAudience string `xml:"AudienceRestriction>Audience"`
}

// SamlResponseAttributeValue is the attribute values
type SamlResponseAttributeValue struct {
	XMLName  xml.Name `xml:"AttributeValue"`
	Value    string   `xml:",chardata"`
	XMKNsXS  string   `xml:"xs,attr"`
	XMLNsXSI string   `xml:"xsi,attr"`
	XDSType  string   `xml:"type,attr"`
}

// SamlResponseAttribute is an attribute passed from the saml response
type SamlResponseAttribute struct {
	XMLName xml.Name `xml:"Attribute"`
	Name    string   `xml:"Name,attr"`

	Attribute []SamlResponseAttributeValue `xml:"AttributeValue"`
}

// SamlResponseAttributeStatement is the attributes passed with the saml response
type SamlResponseAttributeStatement struct {
	XMLName    xml.Name                `xml:"AttributeStatement"`
	Attributes []SamlResponseAttribute `xml:"Attribute"`
}

// SamlResponseAuthStatement is the auth context for the saml response
type SamlResponseAuthStatement struct {
	XMLName      string `xml:"AuthnStatement"`
	AuthnInstant string `xml:"AuthnInstant,attr"`
	SessionIndex string `xml:"SessionIndex,attr"`

	AuthnContextClassRef string `xml:"AuthnContext>AuthnContextClassRef"`
}

// SamlResponseAssertion is the assertion from the saml response
type SamlResponseAssertion struct {
	XMLName      xml.Name `xml:"Assertion"`
	XmlnsSaml2   string   `xml:"saml2,attr"`
	ID           string   `xml:"ID,attr"`
	IssueInstant string   `xml:"IssueInstant,attr"`
	Version      string   `xml:"Version,attr"`
	Issuer       string   `xml:"Issuer"`

	Signature  SamlResponseSignature
	Subject    SamlResponseSubject
	Conditions SamlResponseConditions

	AttributeStatement SamlResponseAttributeStatement
	AuthStatement      SamlResponseAuthStatement
}

// Saml2pResponse is the saml response
type Saml2pResponse struct {
	XMLName      xml.Name `xml:"Response"`
	XmlnsSaml2P  string   `xml:"saml2p,attr"`
	Destination  string   `xml:"Destination,attr"`
	ID           string   `xml:"ID,attr"`
	InResponseTo string   `xml:"InResponseTo,attr"`
	IssueInstant string   `xml:"IssueInstant,attr"`
	Version      string   `xml:"Version,attr"`

	Saml2Issuer struct {
		Issuer     string `xml:",chardata"`
		XMLNsSaml2 string `xml:"saml2,attr"`
	} `xml:"Issuer"`

	Status    SamlResponseStatus
	Assertion SamlResponseAssertion
}
