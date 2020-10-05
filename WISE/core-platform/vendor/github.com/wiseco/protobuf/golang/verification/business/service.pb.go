// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.4
// source: verification/business/service.proto

package business

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	golang "github.com/wiseco/protobuf/golang"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type VerificationStatus int32

const (
	VerificationStatus_VS_UNSPECIFIED VerificationStatus = 0
	VerificationStatus_VS_OPEN        VerificationStatus = 100
	VerificationStatus_VS_PENDING     VerificationStatus = 101
	VerificationStatus_VS_IN_AUDIT    VerificationStatus = 102
	VerificationStatus_VS_IN_REVIEW   VerificationStatus = 103
	VerificationStatus_VS_APPROVED    VerificationStatus = 104
	VerificationStatus_VS_DECLINED    VerificationStatus = 105
)

// Enum value maps for VerificationStatus.
var (
	VerificationStatus_name = map[int32]string{
		0:   "VS_UNSPECIFIED",
		100: "VS_OPEN",
		101: "VS_PENDING",
		102: "VS_IN_AUDIT",
		103: "VS_IN_REVIEW",
		104: "VS_APPROVED",
		105: "VS_DECLINED",
	}
	VerificationStatus_value = map[string]int32{
		"VS_UNSPECIFIED": 0,
		"VS_OPEN":        100,
		"VS_PENDING":     101,
		"VS_IN_AUDIT":    102,
		"VS_IN_REVIEW":   103,
		"VS_APPROVED":    104,
		"VS_DECLINED":    105,
	}
)

func (x VerificationStatus) Enum() *VerificationStatus {
	p := new(VerificationStatus)
	*p = x
	return p
}

func (x VerificationStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VerificationStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_verification_business_service_proto_enumTypes[0].Descriptor()
}

func (VerificationStatus) Type() protoreflect.EnumType {
	return &file_verification_business_service_proto_enumTypes[0]
}

func (x VerificationStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VerificationStatus.Descriptor instead.
func (VerificationStatus) EnumDescriptor() ([]byte, []int) {
	return file_verification_business_service_proto_rawDescGZIP(), []int{0}
}

type TinStatus int32

const (
	TinStatus_TS_UNSPECIFIED TinStatus = 0
	TinStatus_TS_VERIFIED    TinStatus = 100
	TinStatus_TS_MISMATCH    TinStatus = 101
	TinStatus_TS_UNKNOWN     TinStatus = 102
)

// Enum value maps for TinStatus.
var (
	TinStatus_name = map[int32]string{
		0:   "TS_UNSPECIFIED",
		100: "TS_VERIFIED",
		101: "TS_MISMATCH",
		102: "TS_UNKNOWN",
	}
	TinStatus_value = map[string]int32{
		"TS_UNSPECIFIED": 0,
		"TS_VERIFIED":    100,
		"TS_MISMATCH":    101,
		"TS_UNKNOWN":     102,
	}
)

func (x TinStatus) Enum() *TinStatus {
	p := new(TinStatus)
	*p = x
	return p
}

func (x TinStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TinStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_verification_business_service_proto_enumTypes[1].Descriptor()
}

func (TinStatus) Type() protoreflect.EnumType {
	return &file_verification_business_service_proto_enumTypes[1]
}

func (x TinStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TinStatus.Descriptor instead.
func (TinStatus) EnumDescriptor() ([]byte, []int) {
	return file_verification_business_service_proto_rawDescGZIP(), []int{1}
}

//business encapsulates what we care about the business verfication results
type BusinessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//business_id from the core db
	BusinessId     string               `protobuf:"bytes,1,opt,name=business_id,json=businessId,proto3" json:"business_id,omitempty"`
	Provider       string               `protobuf:"bytes,2,opt,name=provider,proto3" json:"provider,omitempty"`
	Name           string               `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Status         VerificationStatus   `protobuf:"varint,4,opt,name=status,proto3,enum=wise.protobuf.verification.business.VerificationStatus" json:"status,omitempty"`
	Tin            string               `protobuf:"bytes,5,opt,name=tin,proto3" json:"tin,omitempty"`
	TinStatus      TinStatus            `protobuf:"varint,6,opt,name=tin_status,json=tinStatus,proto3,enum=wise.protobuf.verification.business.TinStatus" json:"tin_status,omitempty"`
	EntityType     EntityType           `protobuf:"varint,7,opt,name=entity_type,json=entityType,proto3,enum=wise.protobuf.verification.business.EntityType" json:"entity_type,omitempty"`
	FormationDate  string               `protobuf:"bytes,8,opt,name=formation_date,json=formationDate,proto3" json:"formation_date,omitempty"`
	FormationState string               `protobuf:"bytes,9,opt,name=formation_state,json=formationState,proto3" json:"formation_state,omitempty"`
	Registrations  []*Registration      `protobuf:"bytes,10,rep,name=registrations,proto3" json:"registrations,omitempty"`
	Officers       []*Officer           `protobuf:"bytes,11,rep,name=officers,proto3" json:"officers,omitempty"`
	Addresses      []*golang.Address    `protobuf:"bytes,12,rep,name=addresses,proto3" json:"addresses,omitempty"`
	WebsiteUrl     string               `protobuf:"bytes,13,opt,name=website_url,json=websiteUrl,proto3" json:"website_url,omitempty"`
	Watchlists     []*Watchlist         `protobuf:"bytes,14,rep,name=watchlists,proto3" json:"watchlists,omitempty"`
	PhoneNumbers   []string             `protobuf:"bytes,15,rep,name=phone_numbers,json=phoneNumbers,proto3" json:"phone_numbers,omitempty"`
	RawResults     []byte               `protobuf:"bytes,16,opt,name=raw_results,json=rawResults,proto3" json:"raw_results,omitempty"`
	LastVerified   *timestamp.Timestamp `protobuf:"bytes,17,opt,name=last_verified,json=lastVerified,proto3" json:"last_verified,omitempty"`
	Created        *timestamp.Timestamp `protobuf:"bytes,18,opt,name=created,proto3" json:"created,omitempty"`
	Modified       *timestamp.Timestamp `protobuf:"bytes,19,opt,name=modified,proto3" json:"modified,omitempty"`
}

func (x *BusinessResponse) Reset() {
	*x = BusinessResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verification_business_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusinessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusinessResponse) ProtoMessage() {}

func (x *BusinessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_verification_business_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusinessResponse.ProtoReflect.Descriptor instead.
func (*BusinessResponse) Descriptor() ([]byte, []int) {
	return file_verification_business_service_proto_rawDescGZIP(), []int{0}
}

func (x *BusinessResponse) GetBusinessId() string {
	if x != nil {
		return x.BusinessId
	}
	return ""
}

func (x *BusinessResponse) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *BusinessResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BusinessResponse) GetStatus() VerificationStatus {
	if x != nil {
		return x.Status
	}
	return VerificationStatus_VS_UNSPECIFIED
}

func (x *BusinessResponse) GetTin() string {
	if x != nil {
		return x.Tin
	}
	return ""
}

func (x *BusinessResponse) GetTinStatus() TinStatus {
	if x != nil {
		return x.TinStatus
	}
	return TinStatus_TS_UNSPECIFIED
}

func (x *BusinessResponse) GetEntityType() EntityType {
	if x != nil {
		return x.EntityType
	}
	return EntityType_ET_UNSPECIFIED
}

func (x *BusinessResponse) GetFormationDate() string {
	if x != nil {
		return x.FormationDate
	}
	return ""
}

func (x *BusinessResponse) GetFormationState() string {
	if x != nil {
		return x.FormationState
	}
	return ""
}

func (x *BusinessResponse) GetRegistrations() []*Registration {
	if x != nil {
		return x.Registrations
	}
	return nil
}

func (x *BusinessResponse) GetOfficers() []*Officer {
	if x != nil {
		return x.Officers
	}
	return nil
}

func (x *BusinessResponse) GetAddresses() []*golang.Address {
	if x != nil {
		return x.Addresses
	}
	return nil
}

func (x *BusinessResponse) GetWebsiteUrl() string {
	if x != nil {
		return x.WebsiteUrl
	}
	return ""
}

func (x *BusinessResponse) GetWatchlists() []*Watchlist {
	if x != nil {
		return x.Watchlists
	}
	return nil
}

func (x *BusinessResponse) GetPhoneNumbers() []string {
	if x != nil {
		return x.PhoneNumbers
	}
	return nil
}

func (x *BusinessResponse) GetRawResults() []byte {
	if x != nil {
		return x.RawResults
	}
	return nil
}

func (x *BusinessResponse) GetLastVerified() *timestamp.Timestamp {
	if x != nil {
		return x.LastVerified
	}
	return nil
}

func (x *BusinessResponse) GetCreated() *timestamp.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *BusinessResponse) GetModified() *timestamp.Timestamp {
	if x != nil {
		return x.Modified
	}
	return nil
}

//VerificationRequest a request sent to verify a business
type VerificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BusinessId   string            `protobuf:"bytes,1,opt,name=business_id,json=businessId,proto3" json:"business_id,omitempty"`
	Name         string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Tin          string            `protobuf:"bytes,3,opt,name=tin,proto3" json:"tin,omitempty"`
	WebsiteUrl   string            `protobuf:"bytes,4,opt,name=website_url,json=websiteUrl,proto3" json:"website_url,omitempty"`
	Addresses    []*AddressRequest `protobuf:"bytes,5,rep,name=addresses,proto3" json:"addresses,omitempty"`
	PhoneNumbers []string          `protobuf:"bytes,6,rep,name=phone_numbers,json=phoneNumbers,proto3" json:"phone_numbers,omitempty"`
}

func (x *VerificationRequest) Reset() {
	*x = VerificationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verification_business_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerificationRequest) ProtoMessage() {}

func (x *VerificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_verification_business_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerificationRequest.ProtoReflect.Descriptor instead.
func (*VerificationRequest) Descriptor() ([]byte, []int) {
	return file_verification_business_service_proto_rawDescGZIP(), []int{1}
}

func (x *VerificationRequest) GetBusinessId() string {
	if x != nil {
		return x.BusinessId
	}
	return ""
}

func (x *VerificationRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *VerificationRequest) GetTin() string {
	if x != nil {
		return x.Tin
	}
	return ""
}

func (x *VerificationRequest) GetWebsiteUrl() string {
	if x != nil {
		return x.WebsiteUrl
	}
	return ""
}

func (x *VerificationRequest) GetAddresses() []*AddressRequest {
	if x != nil {
		return x.Addresses
	}
	return nil
}

func (x *VerificationRequest) GetPhoneNumbers() []string {
	if x != nil {
		return x.PhoneNumbers
	}
	return nil
}

type GetBusinessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//business_id from the core db
	BusinessId string `protobuf:"bytes,1,opt,name=business_id,json=businessId,proto3" json:"business_id,omitempty"`
}

func (x *GetBusinessRequest) Reset() {
	*x = GetBusinessRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verification_business_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBusinessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBusinessRequest) ProtoMessage() {}

func (x *GetBusinessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_verification_business_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBusinessRequest.ProtoReflect.Descriptor instead.
func (*GetBusinessRequest) Descriptor() ([]byte, []int) {
	return file_verification_business_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetBusinessRequest) GetBusinessId() string {
	if x != nil {
		return x.BusinessId
	}
	return ""
}

var File_verification_business_service_proto protoreflect.FileDescriptor

var file_verification_business_service_proto_rawDesc = []byte{
	0x0a, 0x23, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x62,
	0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x23, 0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x1a, 0x0d, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x76, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73,
	0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23,
	0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x62, 0x75, 0x73,
	0x69, 0x6e, 0x65, 0x73, 0x73, 0x2f, 0x6f, 0x66, 0x66, 0x69, 0x63, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x28, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x76,
	0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x62, 0x75, 0x73, 0x69,
	0x6e, 0x65, 0x73, 0x73, 0x2f, 0x77, 0x61, 0x74, 0x63, 0x68, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf6, 0x07, 0x0a, 0x10, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65,
	0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x62, 0x75,
	0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x4f, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x37, 0x2e, 0x77, 0x69,
	0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x76, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73,
	0x73, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03,
	0x74, 0x69, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x69, 0x6e, 0x12, 0x4d,
	0x0a, 0x0a, 0x74, 0x69, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x2e, 0x2e, 0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x2e, 0x54, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x09, 0x74, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x50, 0x0a,
	0x0b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x2f, 0x2e, 0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x25, 0x0a, 0x0e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x61, 0x74,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x57, 0x0a, 0x0d, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x2e, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x48, 0x0a, 0x08, 0x6f, 0x66, 0x66, 0x69,
	0x63, 0x65, 0x72, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x77, 0x69, 0x73,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73,
	0x2e, 0x4f, 0x66, 0x66, 0x69, 0x63, 0x65, 0x72, 0x52, 0x08, 0x6f, 0x66, 0x66, 0x69, 0x63, 0x65,
	0x72, 0x73, 0x12, 0x34, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18,
	0x0c, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x09, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x65, 0x62, 0x73,
	0x69, 0x74, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77,
	0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x4e, 0x0a, 0x0a, 0x77, 0x61, 0x74,
	0x63, 0x68, 0x6c, 0x69, 0x73, 0x74, 0x73, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e,
	0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x76, 0x65,
	0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e,
	0x65, 0x73, 0x73, 0x2e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x6c, 0x69, 0x73, 0x74, 0x52, 0x0a, 0x77,
	0x61, 0x74, 0x63, 0x68, 0x6c, 0x69, 0x73, 0x74, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x68, 0x6f,
	0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0c, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x61, 0x77, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x10, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0a, 0x72, 0x61, 0x77, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12,
	0x3f, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64,
	0x18, 0x11, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64,
	0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x12, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x36, 0x0a, 0x08, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69,
	0x65, 0x64, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x22, 0xf5,
	0x01, 0x0a, 0x13, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65,
	0x73, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x62, 0x75, 0x73,
	0x69, 0x6e, 0x65, 0x73, 0x73, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74,
	0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x69, 0x6e, 0x12, 0x1f, 0x0a,
	0x0b, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x51,
	0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x33, 0x2e, 0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x62,
	0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65,
	0x73, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x22, 0x35, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x42, 0x75, 0x73,
	0x69, 0x6e, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b,
	0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x49, 0x64, 0x2a, 0x8a, 0x01,
	0x0a, 0x12, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x0e, 0x56, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x56, 0x53, 0x5f, 0x4f,
	0x50, 0x45, 0x4e, 0x10, 0x64, 0x12, 0x0e, 0x0a, 0x0a, 0x56, 0x53, 0x5f, 0x50, 0x45, 0x4e, 0x44,
	0x49, 0x4e, 0x47, 0x10, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x56, 0x53, 0x5f, 0x49, 0x4e, 0x5f, 0x41,
	0x55, 0x44, 0x49, 0x54, 0x10, 0x66, 0x12, 0x10, 0x0a, 0x0c, 0x56, 0x53, 0x5f, 0x49, 0x4e, 0x5f,
	0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x10, 0x67, 0x12, 0x0f, 0x0a, 0x0b, 0x56, 0x53, 0x5f, 0x41,
	0x50, 0x50, 0x52, 0x4f, 0x56, 0x45, 0x44, 0x10, 0x68, 0x12, 0x0f, 0x0a, 0x0b, 0x56, 0x53, 0x5f,
	0x44, 0x45, 0x43, 0x4c, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x69, 0x2a, 0x51, 0x0a, 0x09, 0x54, 0x69,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x0e, 0x54, 0x53, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x54,
	0x53, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x64, 0x12, 0x0f, 0x0a, 0x0b,
	0x54, 0x53, 0x5f, 0x4d, 0x49, 0x53, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x10, 0x65, 0x12, 0x0e, 0x0a,
	0x0a, 0x54, 0x53, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x66, 0x32, 0x91, 0x02,
	0x0a, 0x0f, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x7f, 0x0a, 0x0c, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x38, 0x2e, 0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x62,
	0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x35, 0x2e, 0x77, 0x69,
	0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x76, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73,
	0x73, 0x2e, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x7d, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73,
	0x73, 0x12, 0x37, 0x2e, 0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x62,
	0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x75, 0x73, 0x69, 0x6e,
	0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x35, 0x2e, 0x77, 0x69, 0x73,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73,
	0x2e, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x77, 0x69, 0x73, 0x65, 0x63, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_verification_business_service_proto_rawDescOnce sync.Once
	file_verification_business_service_proto_rawDescData = file_verification_business_service_proto_rawDesc
)

func file_verification_business_service_proto_rawDescGZIP() []byte {
	file_verification_business_service_proto_rawDescOnce.Do(func() {
		file_verification_business_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_verification_business_service_proto_rawDescData)
	})
	return file_verification_business_service_proto_rawDescData
}

var file_verification_business_service_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_verification_business_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_verification_business_service_proto_goTypes = []interface{}{
	(VerificationStatus)(0),     // 0: wise.protobuf.verification.business.VerificationStatus
	(TinStatus)(0),              // 1: wise.protobuf.verification.business.TinStatus
	(*BusinessResponse)(nil),    // 2: wise.protobuf.verification.business.BusinessResponse
	(*VerificationRequest)(nil), // 3: wise.protobuf.verification.business.VerificationRequest
	(*GetBusinessRequest)(nil),  // 4: wise.protobuf.verification.business.GetBusinessRequest
	(EntityType)(0),             // 5: wise.protobuf.verification.business.EntityType
	(*Registration)(nil),        // 6: wise.protobuf.verification.business.Registration
	(*Officer)(nil),             // 7: wise.protobuf.verification.business.Officer
	(*golang.Address)(nil),      // 8: wise.protobuf.Address
	(*Watchlist)(nil),           // 9: wise.protobuf.verification.business.Watchlist
	(*timestamp.Timestamp)(nil), // 10: google.protobuf.Timestamp
	(*AddressRequest)(nil),      // 11: wise.protobuf.verification.business.AddressRequest
}
var file_verification_business_service_proto_depIdxs = []int32{
	0,  // 0: wise.protobuf.verification.business.BusinessResponse.status:type_name -> wise.protobuf.verification.business.VerificationStatus
	1,  // 1: wise.protobuf.verification.business.BusinessResponse.tin_status:type_name -> wise.protobuf.verification.business.TinStatus
	5,  // 2: wise.protobuf.verification.business.BusinessResponse.entity_type:type_name -> wise.protobuf.verification.business.EntityType
	6,  // 3: wise.protobuf.verification.business.BusinessResponse.registrations:type_name -> wise.protobuf.verification.business.Registration
	7,  // 4: wise.protobuf.verification.business.BusinessResponse.officers:type_name -> wise.protobuf.verification.business.Officer
	8,  // 5: wise.protobuf.verification.business.BusinessResponse.addresses:type_name -> wise.protobuf.Address
	9,  // 6: wise.protobuf.verification.business.BusinessResponse.watchlists:type_name -> wise.protobuf.verification.business.Watchlist
	10, // 7: wise.protobuf.verification.business.BusinessResponse.last_verified:type_name -> google.protobuf.Timestamp
	10, // 8: wise.protobuf.verification.business.BusinessResponse.created:type_name -> google.protobuf.Timestamp
	10, // 9: wise.protobuf.verification.business.BusinessResponse.modified:type_name -> google.protobuf.Timestamp
	11, // 10: wise.protobuf.verification.business.VerificationRequest.addresses:type_name -> wise.protobuf.verification.business.AddressRequest
	3,  // 11: wise.protobuf.verification.business.BusinessService.Verification:input_type -> wise.protobuf.verification.business.VerificationRequest
	4,  // 12: wise.protobuf.verification.business.BusinessService.GetBusiness:input_type -> wise.protobuf.verification.business.GetBusinessRequest
	2,  // 13: wise.protobuf.verification.business.BusinessService.Verification:output_type -> wise.protobuf.verification.business.BusinessResponse
	2,  // 14: wise.protobuf.verification.business.BusinessService.GetBusiness:output_type -> wise.protobuf.verification.business.BusinessResponse
	13, // [13:15] is the sub-list for method output_type
	11, // [11:13] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_verification_business_service_proto_init() }
func file_verification_business_service_proto_init() {
	if File_verification_business_service_proto != nil {
		return
	}
	file_verification_business_address_proto_init()
	file_verification_business_officer_proto_init()
	file_verification_business_registration_proto_init()
	file_verification_business_watchlist_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_verification_business_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BusinessResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_verification_business_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerificationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_verification_business_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBusinessRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_verification_business_service_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_verification_business_service_proto_goTypes,
		DependencyIndexes: file_verification_business_service_proto_depIdxs,
		EnumInfos:         file_verification_business_service_proto_enumTypes,
		MessageInfos:      file_verification_business_service_proto_msgTypes,
	}.Build()
	File_verification_business_service_proto = out.File
	file_verification_business_service_proto_rawDesc = nil
	file_verification_business_service_proto_goTypes = nil
	file_verification_business_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BusinessServiceClient is the client API for BusinessService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BusinessServiceClient interface {
	// Verification runs or re-runs a verification request for the business
	Verification(ctx context.Context, in *VerificationRequest, opts ...grpc.CallOption) (*BusinessResponse, error)
	// GetBusiness returns relevant business verification results
	GetBusiness(ctx context.Context, in *GetBusinessRequest, opts ...grpc.CallOption) (*BusinessResponse, error)
}

type businessServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBusinessServiceClient(cc grpc.ClientConnInterface) BusinessServiceClient {
	return &businessServiceClient{cc}
}

func (c *businessServiceClient) Verification(ctx context.Context, in *VerificationRequest, opts ...grpc.CallOption) (*BusinessResponse, error) {
	out := new(BusinessResponse)
	err := c.cc.Invoke(ctx, "/wise.protobuf.verification.business.BusinessService/Verification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessServiceClient) GetBusiness(ctx context.Context, in *GetBusinessRequest, opts ...grpc.CallOption) (*BusinessResponse, error) {
	out := new(BusinessResponse)
	err := c.cc.Invoke(ctx, "/wise.protobuf.verification.business.BusinessService/GetBusiness", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BusinessServiceServer is the server API for BusinessService service.
type BusinessServiceServer interface {
	// Verification runs or re-runs a verification request for the business
	Verification(context.Context, *VerificationRequest) (*BusinessResponse, error)
	// GetBusiness returns relevant business verification results
	GetBusiness(context.Context, *GetBusinessRequest) (*BusinessResponse, error)
}

// UnimplementedBusinessServiceServer can be embedded to have forward compatible implementations.
type UnimplementedBusinessServiceServer struct {
}

func (*UnimplementedBusinessServiceServer) Verification(context.Context, *VerificationRequest) (*BusinessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verification not implemented")
}
func (*UnimplementedBusinessServiceServer) GetBusiness(context.Context, *GetBusinessRequest) (*BusinessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBusiness not implemented")
}

func RegisterBusinessServiceServer(s *grpc.Server, srv BusinessServiceServer) {
	s.RegisterService(&_BusinessService_serviceDesc, srv)
}

func _BusinessService_Verification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessServiceServer).Verification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wise.protobuf.verification.business.BusinessService/Verification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessServiceServer).Verification(ctx, req.(*VerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessService_GetBusiness_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBusinessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessServiceServer).GetBusiness(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wise.protobuf.verification.business.BusinessService/GetBusiness",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessServiceServer).GetBusiness(ctx, req.(*GetBusinessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BusinessService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "wise.protobuf.verification.business.BusinessService",
	HandlerType: (*BusinessServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Verification",
			Handler:    _BusinessService_Verification_Handler,
		},
		{
			MethodName: "GetBusiness",
			Handler:    _BusinessService_GetBusiness_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "verification/business/service.proto",
}
