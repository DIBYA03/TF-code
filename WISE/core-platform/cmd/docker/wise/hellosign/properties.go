package main

import (
	"github.com/wiseco/core-platform/services/signature"
)

type EventMessage struct {
	Event            Event            `json:"event"`
	SignatureRequest SignatureRequest `json:"signature_request"`
}

type Event struct {
	EventTime string    `json:"event_time"`
	EventType EventType `json:"event_type"`
	EventHash string    `json:"event_hash"`
}

type SignatureRequest struct {
	SignatureRequestID string `json:"signature_request_id"`
}

type EventType string

const (
	EventTypeSignatureRequestViewed       = EventType("signature_request_viewed")
	EventTypeSignatureRequestSigned       = EventType("signature_request_signed")
	EventTypeSignatureRequestDownloadable = EventType("signature_request_downloadable")
	EventTypeSignatureRequestSent         = EventType("signature_request_sent")
	EventTypeSignatureRequestDeclined     = EventType("signature_request_declined")
	EventTypeSignatureRequestReassigned   = EventType("signature_request_reassigned")
	EventTypeSignatureRequestRemind       = EventType("signature_request_remind")
	EventTypeSignatureRequestAllSigned    = EventType("signature_request_all_signed")
	EventTypeSignatureRequestEmailBounce  = EventType("signature_request_email_bounce")
	EventTypeSignatureRequestInvalid      = EventType("signature_request_invalid")
	EventTypeSignatureRequestCanceled     = EventType("signature_request_canceled")
	EventTypeSignatureRequestPrepared     = EventType("file_error")
	EventTypeUnknownError                 = EventType("unknown_error")
	EventTypeSignUrlInvalid               = EventType("sign_url_invalid")
	EventTypeAccountConfirmed             = EventType("account_confirmed")
	EventTypeTemplateCreated              = EventType("template_created")
	EventTypeTemplateError                = EventType("template_error")
	EventTypeTemplateCallbackTest         = EventType("callback_test")
)

var EventTypeMap = map[EventType]signature.EventType{
	EventTypeSignatureRequestViewed:       signature.EventTypeSignatureRequestViewed,
	EventTypeSignatureRequestSigned:       signature.EventTypeSignatureRequestSigned,
	EventTypeSignatureRequestDownloadable: signature.EventTypeSignatureRequestDownloadable,
	EventTypeSignatureRequestSent:         signature.EventTypeSignatureRequestSent,
	EventTypeSignatureRequestDeclined:     signature.EventTypeSignatureRequestDeclined,
	EventTypeSignatureRequestReassigned:   signature.EventTypeSignatureRequestReassigned,
	EventTypeSignatureRequestRemind:       signature.EventTypeSignatureRequestRemind,
	EventTypeSignatureRequestAllSigned:    signature.EventTypeSignatureRequestAllSigned,
	EventTypeSignatureRequestEmailBounce:  signature.EventTypeSignatureRequestEmailBounce,
	EventTypeSignatureRequestInvalid:      signature.EventTypeSignatureRequestInvalid,
	EventTypeSignatureRequestCanceled:     signature.EventTypeSignatureRequestCanceled,
	EventTypeSignatureRequestPrepared:     signature.EventTypeSignatureRequestPrepared,
	EventTypeUnknownError:                 signature.EventTypeOther,
	EventTypeSignUrlInvalid:               signature.EventTypeOther,
	EventTypeAccountConfirmed:             signature.EventTypeOther,
	EventTypeTemplateCreated:              signature.EventTypeOther,
	EventTypeTemplateError:                signature.EventTypeOther,
	EventTypeTemplateCallbackTest:         signature.EventTypeOther,
}
