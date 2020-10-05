package bbva

import partnerbank "github.com/wiseco/core-platform/partner/bank"

type IdentityDocumentFileType string

const (
	IdentityDocumentFileTypePDF = IdentityDocumentFileType("application/pdf")
	IdentityDocumentFileTypeJPG = IdentityDocumentFileType("application/jpg")
	IdentityDocumentFileTypePNG = IdentityDocumentFileType("application/png")
)

var partnerContentTypeFrom = map[partnerbank.ContentType]IdentityDocumentFileType{
	partnerbank.ContentTypePDF:  IdentityDocumentFileTypePDF,
	partnerbank.ContentTypeJPEG: IdentityDocumentFileTypeJPG,
	partnerbank.ContentTypePNG:  IdentityDocumentFileTypePNG,
}

var partnerContentTypeTo = map[IdentityDocumentFileType]partnerbank.ContentType{
	IdentityDocumentFileTypePDF: partnerbank.ContentTypePDF,
	IdentityDocumentFileTypeJPG: partnerbank.ContentTypeJPEG,
	IdentityDocumentFileTypePNG: partnerbank.ContentTypePNG,
}

// IdentityDocumentRequest
// Provide which "idv_required" is to be resolved with this document.
// For example, if you received a kyc.idv_required = "address", "name",
// and the end-user provides a driver's license photo that verifies
// both "address" and "name", this field should contain "address" and "name".
type ConsumerIdentityDocumentRequest struct {
	File             string                   `json:"file"`       // base64 encoded data
	FileType         IdentityDocumentFileType `json:"file_type"`  // only allowed: "application/pdf"
	IDVerifyRequired []DocumentIDVerify       `json:"verify_idv"` // ID verify type
	DocType          ConsumerIdentityDocument `json:"doc_type"`   // Consumer document type
}

type BusinessIdentityDocumentRequest struct {
	File             string                   `json:"file"`       // base64 encoded data
	FileType         IdentityDocumentFileType `json:"file_type"`  // only allowed: "application/pdf"
	IDVerifyRequired []DocumentIDVerify       `json:"verify_idv"` // ID verify type
}
