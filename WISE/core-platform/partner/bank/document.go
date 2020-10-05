package bank

type ContentType string

const (
	ContentTypePDF  = ContentType("application/pdf")
	ContentTypeJPEG = ContentType("image/jpeg")
	ContentTypePNG  = ContentType("image/png")
)

type ConsumerIdentityDocumentRequest struct {
	ConsumerID       ConsumerID               `json:"entityId"`
	IdentityDocument ConsumerIdentityDocument `json:"identityDocument"`
	IDVerifyRequired []IDVerify               `json:"idVerifyRequired"` // Id verification required
	ContentType      ContentType              `json:"contentType"`
	Content          []byte                   `json:"content"`
}

type BusinessIdentityDocumentRequest struct {
	BusinessID       BusinessID               `json:"entityId"`
	IdentityDocument BusinessIdentityDocument `json:"identityDocument"`
	IDVerifyRequired []IDVerify               `json:"idVerifyRequired"` // Id verification required
	ContentType      ContentType              `json:"contentType"`
	Content          []byte                   `json:"content"`
}

type IdentityDocumentResponse struct {
	Id string `json:"id"`
}
