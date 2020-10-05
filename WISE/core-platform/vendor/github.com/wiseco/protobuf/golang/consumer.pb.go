// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.4
// source: consumer.proto

package golang

import (
	proto "github.com/golang/protobuf/proto"
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

type ConsumerKYCStatus int32

const (
	ConsumerKYCStatus_CKS_UNSPECIFIED ConsumerKYCStatus = 0
	ConsumerKYCStatus_CKS_NOT_STARTED ConsumerKYCStatus = 100
	ConsumerKYCStatus_CKS_SUBMITTED   ConsumerKYCStatus = 101
	ConsumerKYCStatus_CKS_REVIEW      ConsumerKYCStatus = 102
	ConsumerKYCStatus_CKS_APPROVED    ConsumerKYCStatus = 103
	ConsumerKYCStatus_CKS_DECLINED    ConsumerKYCStatus = 104
)

// Enum value maps for ConsumerKYCStatus.
var (
	ConsumerKYCStatus_name = map[int32]string{
		0:   "CKS_UNSPECIFIED",
		100: "CKS_NOT_STARTED",
		101: "CKS_SUBMITTED",
		102: "CKS_REVIEW",
		103: "CKS_APPROVED",
		104: "CKS_DECLINED",
	}
	ConsumerKYCStatus_value = map[string]int32{
		"CKS_UNSPECIFIED": 0,
		"CKS_NOT_STARTED": 100,
		"CKS_SUBMITTED":   101,
		"CKS_REVIEW":      102,
		"CKS_APPROVED":    103,
		"CKS_DECLINED":    104,
	}
)

func (x ConsumerKYCStatus) Enum() *ConsumerKYCStatus {
	p := new(ConsumerKYCStatus)
	*p = x
	return p
}

func (x ConsumerKYCStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConsumerKYCStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_consumer_proto_enumTypes[0].Descriptor()
}

func (ConsumerKYCStatus) Type() protoreflect.EnumType {
	return &file_consumer_proto_enumTypes[0]
}

func (x ConsumerKYCStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConsumerKYCStatus.Descriptor instead.
func (ConsumerKYCStatus) EnumDescriptor() ([]byte, []int) {
	return file_consumer_proto_rawDescGZIP(), []int{0}
}

type ConsumerStatus int32

const (
	ConsumerStatus_CS_UNSPECIFIED ConsumerStatus = 0
	ConsumerStatus_CS_ACTIVE      ConsumerStatus = 100
	ConsumerStatus_CS_INACTIVE    ConsumerStatus = 101
	ConsumerStatus_CS_SUSPENDED   ConsumerStatus = 102
	ConsumerStatus_CS_CLOSED      ConsumerStatus = 103
)

// Enum value maps for ConsumerStatus.
var (
	ConsumerStatus_name = map[int32]string{
		0:   "CS_UNSPECIFIED",
		100: "CS_ACTIVE",
		101: "CS_INACTIVE",
		102: "CS_SUSPENDED",
		103: "CS_CLOSED",
	}
	ConsumerStatus_value = map[string]int32{
		"CS_UNSPECIFIED": 0,
		"CS_ACTIVE":      100,
		"CS_INACTIVE":    101,
		"CS_SUSPENDED":   102,
		"CS_CLOSED":      103,
	}
)

func (x ConsumerStatus) Enum() *ConsumerStatus {
	p := new(ConsumerStatus)
	*p = x
	return p
}

func (x ConsumerStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConsumerStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_consumer_proto_enumTypes[1].Descriptor()
}

func (ConsumerStatus) Type() protoreflect.EnumType {
	return &file_consumer_proto_enumTypes[1]
}

func (x ConsumerStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConsumerStatus.Descriptor instead.
func (ConsumerStatus) EnumDescriptor() ([]byte, []int) {
	return file_consumer_proto_rawDescGZIP(), []int{1}
}

type ConsumerOccupation int32

const (
	ConsumerOccupation_CO_UNSPECIFIED                 ConsumerOccupation = 0
	ConsumerOccupation_CO_AGRICULTURE                 ConsumerOccupation = 100
	ConsumerOccupation_CO_CLERGY_MINISTRY_STAFF       ConsumerOccupation = 101
	ConsumerOccupation_CO_CONSTRUCTION_INDUSTRIAL     ConsumerOccupation = 102
	ConsumerOccupation_CO_EDUCATION                   ConsumerOccupation = 103
	ConsumerOccupation_CO_FINANCE_ACCOUNTING_TAX      ConsumerOccupation = 104
	ConsumerOccupation_CO_FIRE_FIRST_RESPONDERS       ConsumerOccupation = 105
	ConsumerOccupation_CO_HEALTHCARE                  ConsumerOccupation = 106
	ConsumerOccupation_CO_HOMEMAKER                   ConsumerOccupation = 107
	ConsumerOccupation_CO_LABOR_GENERAL               ConsumerOccupation = 108
	ConsumerOccupation_CO_LABOR_SKILLED               ConsumerOccupation = 109
	ConsumerOccupation_CO_LAW_ENFORCEMENT_SECURITY    ConsumerOccupation = 110
	ConsumerOccupation_CO_LEGAL_SERVICES              ConsumerOccupation = 111
	ConsumerOccupation_CO_MILITARY                    ConsumerOccupation = 112
	ConsumerOccupation_CO_NOTARY_REGISTRAR            ConsumerOccupation = 113
	ConsumerOccupation_CO_PRIVATE_INVESTOR            ConsumerOccupation = 114
	ConsumerOccupation_CO_PROFESSIONAL_ADMINISTRATIVE ConsumerOccupation = 115
	ConsumerOccupation_CO_PROFESSIONAL_MANAGEMENT     ConsumerOccupation = 116
	ConsumerOccupation_CO_PROFESSIONAL_OTHER          ConsumerOccupation = 117
	ConsumerOccupation_CO_PROFESSIONAL_TECHNICAL      ConsumerOccupation = 118
	ConsumerOccupation_CO_RETIRED                     ConsumerOccupation = 119
	ConsumerOccupation_CO_SALES                       ConsumerOccupation = 120
	ConsumerOccupation_CO_SELF_EMPLOYED               ConsumerOccupation = 121
	ConsumerOccupation_CO_STUDENT                     ConsumerOccupation = 122
	ConsumerOccupation_CO_TRANSPORTATION              ConsumerOccupation = 123
	ConsumerOccupation_CO_UNEMPLOYED                  ConsumerOccupation = 124
)

// Enum value maps for ConsumerOccupation.
var (
	ConsumerOccupation_name = map[int32]string{
		0:   "CO_UNSPECIFIED",
		100: "CO_AGRICULTURE",
		101: "CO_CLERGY_MINISTRY_STAFF",
		102: "CO_CONSTRUCTION_INDUSTRIAL",
		103: "CO_EDUCATION",
		104: "CO_FINANCE_ACCOUNTING_TAX",
		105: "CO_FIRE_FIRST_RESPONDERS",
		106: "CO_HEALTHCARE",
		107: "CO_HOMEMAKER",
		108: "CO_LABOR_GENERAL",
		109: "CO_LABOR_SKILLED",
		110: "CO_LAW_ENFORCEMENT_SECURITY",
		111: "CO_LEGAL_SERVICES",
		112: "CO_MILITARY",
		113: "CO_NOTARY_REGISTRAR",
		114: "CO_PRIVATE_INVESTOR",
		115: "CO_PROFESSIONAL_ADMINISTRATIVE",
		116: "CO_PROFESSIONAL_MANAGEMENT",
		117: "CO_PROFESSIONAL_OTHER",
		118: "CO_PROFESSIONAL_TECHNICAL",
		119: "CO_RETIRED",
		120: "CO_SALES",
		121: "CO_SELF_EMPLOYED",
		122: "CO_STUDENT",
		123: "CO_TRANSPORTATION",
		124: "CO_UNEMPLOYED",
	}
	ConsumerOccupation_value = map[string]int32{
		"CO_UNSPECIFIED":                 0,
		"CO_AGRICULTURE":                 100,
		"CO_CLERGY_MINISTRY_STAFF":       101,
		"CO_CONSTRUCTION_INDUSTRIAL":     102,
		"CO_EDUCATION":                   103,
		"CO_FINANCE_ACCOUNTING_TAX":      104,
		"CO_FIRE_FIRST_RESPONDERS":       105,
		"CO_HEALTHCARE":                  106,
		"CO_HOMEMAKER":                   107,
		"CO_LABOR_GENERAL":               108,
		"CO_LABOR_SKILLED":               109,
		"CO_LAW_ENFORCEMENT_SECURITY":    110,
		"CO_LEGAL_SERVICES":              111,
		"CO_MILITARY":                    112,
		"CO_NOTARY_REGISTRAR":            113,
		"CO_PRIVATE_INVESTOR":            114,
		"CO_PROFESSIONAL_ADMINISTRATIVE": 115,
		"CO_PROFESSIONAL_MANAGEMENT":     116,
		"CO_PROFESSIONAL_OTHER":          117,
		"CO_PROFESSIONAL_TECHNICAL":      118,
		"CO_RETIRED":                     119,
		"CO_SALES":                       120,
		"CO_SELF_EMPLOYED":               121,
		"CO_STUDENT":                     122,
		"CO_TRANSPORTATION":              123,
		"CO_UNEMPLOYED":                  124,
	}
)

func (x ConsumerOccupation) Enum() *ConsumerOccupation {
	p := new(ConsumerOccupation)
	*p = x
	return p
}

func (x ConsumerOccupation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConsumerOccupation) Descriptor() protoreflect.EnumDescriptor {
	return file_consumer_proto_enumTypes[2].Descriptor()
}

func (ConsumerOccupation) Type() protoreflect.EnumType {
	return &file_consumer_proto_enumTypes[2]
}

func (x ConsumerOccupation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConsumerOccupation.Descriptor instead.
func (ConsumerOccupation) EnumDescriptor() ([]byte, []int) {
	return file_consumer_proto_rawDescGZIP(), []int{2}
}

type ConsumerTaxIDType int32

const (
	ConsumerTaxIDType_CTIDT_UNSPECIFIED ConsumerTaxIDType = 0
	ConsumerTaxIDType_CTIDT_SSN         ConsumerTaxIDType = 100
	ConsumerTaxIDType_CTIDT_ITIN        ConsumerTaxIDType = 101
)

// Enum value maps for ConsumerTaxIDType.
var (
	ConsumerTaxIDType_name = map[int32]string{
		0:   "CTIDT_UNSPECIFIED",
		100: "CTIDT_SSN",
		101: "CTIDT_ITIN",
	}
	ConsumerTaxIDType_value = map[string]int32{
		"CTIDT_UNSPECIFIED": 0,
		"CTIDT_SSN":         100,
		"CTIDT_ITIN":        101,
	}
)

func (x ConsumerTaxIDType) Enum() *ConsumerTaxIDType {
	p := new(ConsumerTaxIDType)
	*p = x
	return p
}

func (x ConsumerTaxIDType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConsumerTaxIDType) Descriptor() protoreflect.EnumDescriptor {
	return file_consumer_proto_enumTypes[3].Descriptor()
}

func (ConsumerTaxIDType) Type() protoreflect.EnumType {
	return &file_consumer_proto_enumTypes[3]
}

func (x ConsumerTaxIDType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConsumerTaxIDType.Descriptor instead.
func (ConsumerTaxIDType) EnumDescriptor() ([]byte, []int) {
	return file_consumer_proto_rawDescGZIP(), []int{3}
}

type ConsumerResidencyStatus int32

const (
	ConsumerResidencyStatus_CRS_UNSPECIFIED  ConsumerResidencyStatus = 0
	ConsumerResidencyStatus_CRS_CITIZEN      ConsumerResidencyStatus = 100
	ConsumerResidencyStatus_CRS_RESIDENT     ConsumerResidencyStatus = 101
	ConsumerResidencyStatus_CRS_NON_RESIDENT ConsumerResidencyStatus = 102
)

// Enum value maps for ConsumerResidencyStatus.
var (
	ConsumerResidencyStatus_name = map[int32]string{
		0:   "CRS_UNSPECIFIED",
		100: "CRS_CITIZEN",
		101: "CRS_RESIDENT",
		102: "CRS_NON_RESIDENT",
	}
	ConsumerResidencyStatus_value = map[string]int32{
		"CRS_UNSPECIFIED":  0,
		"CRS_CITIZEN":      100,
		"CRS_RESIDENT":     101,
		"CRS_NON_RESIDENT": 102,
	}
)

func (x ConsumerResidencyStatus) Enum() *ConsumerResidencyStatus {
	p := new(ConsumerResidencyStatus)
	*p = x
	return p
}

func (x ConsumerResidencyStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConsumerResidencyStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_consumer_proto_enumTypes[4].Descriptor()
}

func (ConsumerResidencyStatus) Type() protoreflect.EnumType {
	return &file_consumer_proto_enumTypes[4]
}

func (x ConsumerResidencyStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConsumerResidencyStatus.Descriptor instead.
func (ConsumerResidencyStatus) EnumDescriptor() ([]byte, []int) {
	return file_consumer_proto_rawDescGZIP(), []int{4}
}

type ConsumerIncomeType int32

const (
	ConsumerIncomeType_CIT_UNSPECIFIED        ConsumerIncomeType = 0
	ConsumerIncomeType_CIT_INHERITANCE        ConsumerIncomeType = 100
	ConsumerIncomeType_CIT_SALARY             ConsumerIncomeType = 101
	ConsumerIncomeType_CIT_SALE_OF_COMPANY    ConsumerIncomeType = 102
	ConsumerIncomeType_CIT_SALE_OF_PROPERTY   ConsumerIncomeType = 103
	ConsumerIncomeType_CIT_INVESTMENTS        ConsumerIncomeType = 104
	ConsumerIncomeType_CIT_LIFE_INSURANCE     ConsumerIncomeType = 105
	ConsumerIncomeType_CIT_DIVORCE_SETTLEMENT ConsumerIncomeType = 106
	ConsumerIncomeType_CIT_OTHER              ConsumerIncomeType = 107
)

// Enum value maps for ConsumerIncomeType.
var (
	ConsumerIncomeType_name = map[int32]string{
		0:   "CIT_UNSPECIFIED",
		100: "CIT_INHERITANCE",
		101: "CIT_SALARY",
		102: "CIT_SALE_OF_COMPANY",
		103: "CIT_SALE_OF_PROPERTY",
		104: "CIT_INVESTMENTS",
		105: "CIT_LIFE_INSURANCE",
		106: "CIT_DIVORCE_SETTLEMENT",
		107: "CIT_OTHER",
	}
	ConsumerIncomeType_value = map[string]int32{
		"CIT_UNSPECIFIED":        0,
		"CIT_INHERITANCE":        100,
		"CIT_SALARY":             101,
		"CIT_SALE_OF_COMPANY":    102,
		"CIT_SALE_OF_PROPERTY":   103,
		"CIT_INVESTMENTS":        104,
		"CIT_LIFE_INSURANCE":     105,
		"CIT_DIVORCE_SETTLEMENT": 106,
		"CIT_OTHER":              107,
	}
)

func (x ConsumerIncomeType) Enum() *ConsumerIncomeType {
	p := new(ConsumerIncomeType)
	*p = x
	return p
}

func (x ConsumerIncomeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConsumerIncomeType) Descriptor() protoreflect.EnumDescriptor {
	return file_consumer_proto_enumTypes[5].Descriptor()
}

func (ConsumerIncomeType) Type() protoreflect.EnumType {
	return &file_consumer_proto_enumTypes[5]
}

func (x ConsumerIncomeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConsumerIncomeType.Descriptor instead.
func (ConsumerIncomeType) EnumDescriptor() ([]byte, []int) {
	return file_consumer_proto_rawDescGZIP(), []int{5}
}

var File_consumer_proto protoreflect.FileDescriptor

var file_consumer_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x77, 0x69, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2a,
	0x84, 0x01, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x4b, 0x59, 0x43, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x4b, 0x53, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x4b,
	0x53, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x52, 0x54, 0x45, 0x44, 0x10, 0x64, 0x12,
	0x11, 0x0a, 0x0d, 0x43, 0x4b, 0x53, 0x5f, 0x53, 0x55, 0x42, 0x4d, 0x49, 0x54, 0x54, 0x45, 0x44,
	0x10, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x4b, 0x53, 0x5f, 0x52, 0x45, 0x56, 0x49, 0x45, 0x57,
	0x10, 0x66, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x4b, 0x53, 0x5f, 0x41, 0x50, 0x50, 0x52, 0x4f, 0x56,
	0x45, 0x44, 0x10, 0x67, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x4b, 0x53, 0x5f, 0x44, 0x45, 0x43, 0x4c,
	0x49, 0x4e, 0x45, 0x44, 0x10, 0x68, 0x2a, 0x65, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x53, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09,
	0x43, 0x53, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x64, 0x12, 0x0f, 0x0a, 0x0b, 0x43,
	0x53, 0x5f, 0x49, 0x4e, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x65, 0x12, 0x10, 0x0a, 0x0c,
	0x43, 0x53, 0x5f, 0x53, 0x55, 0x53, 0x50, 0x45, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x66, 0x12, 0x0d,
	0x0a, 0x09, 0x43, 0x53, 0x5f, 0x43, 0x4c, 0x4f, 0x53, 0x45, 0x44, 0x10, 0x67, 0x2a, 0x81, 0x05,
	0x0a, 0x12, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x4f, 0x63, 0x63, 0x75, 0x70, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x4f, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x4f, 0x5f, 0x41,
	0x47, 0x52, 0x49, 0x43, 0x55, 0x4c, 0x54, 0x55, 0x52, 0x45, 0x10, 0x64, 0x12, 0x1c, 0x0a, 0x18,
	0x43, 0x4f, 0x5f, 0x43, 0x4c, 0x45, 0x52, 0x47, 0x59, 0x5f, 0x4d, 0x49, 0x4e, 0x49, 0x53, 0x54,
	0x52, 0x59, 0x5f, 0x53, 0x54, 0x41, 0x46, 0x46, 0x10, 0x65, 0x12, 0x1e, 0x0a, 0x1a, 0x43, 0x4f,
	0x5f, 0x43, 0x4f, 0x4e, 0x53, 0x54, 0x52, 0x55, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x4e,
	0x44, 0x55, 0x53, 0x54, 0x52, 0x49, 0x41, 0x4c, 0x10, 0x66, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x4f,
	0x5f, 0x45, 0x44, 0x55, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x67, 0x12, 0x1d, 0x0a, 0x19,
	0x43, 0x4f, 0x5f, 0x46, 0x49, 0x4e, 0x41, 0x4e, 0x43, 0x45, 0x5f, 0x41, 0x43, 0x43, 0x4f, 0x55,
	0x4e, 0x54, 0x49, 0x4e, 0x47, 0x5f, 0x54, 0x41, 0x58, 0x10, 0x68, 0x12, 0x1c, 0x0a, 0x18, 0x43,
	0x4f, 0x5f, 0x46, 0x49, 0x52, 0x45, 0x5f, 0x46, 0x49, 0x52, 0x53, 0x54, 0x5f, 0x52, 0x45, 0x53,
	0x50, 0x4f, 0x4e, 0x44, 0x45, 0x52, 0x53, 0x10, 0x69, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x4f, 0x5f,
	0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x43, 0x41, 0x52, 0x45, 0x10, 0x6a, 0x12, 0x10, 0x0a, 0x0c,
	0x43, 0x4f, 0x5f, 0x48, 0x4f, 0x4d, 0x45, 0x4d, 0x41, 0x4b, 0x45, 0x52, 0x10, 0x6b, 0x12, 0x14,
	0x0a, 0x10, 0x43, 0x4f, 0x5f, 0x4c, 0x41, 0x42, 0x4f, 0x52, 0x5f, 0x47, 0x45, 0x4e, 0x45, 0x52,
	0x41, 0x4c, 0x10, 0x6c, 0x12, 0x14, 0x0a, 0x10, 0x43, 0x4f, 0x5f, 0x4c, 0x41, 0x42, 0x4f, 0x52,
	0x5f, 0x53, 0x4b, 0x49, 0x4c, 0x4c, 0x45, 0x44, 0x10, 0x6d, 0x12, 0x1f, 0x0a, 0x1b, 0x43, 0x4f,
	0x5f, 0x4c, 0x41, 0x57, 0x5f, 0x45, 0x4e, 0x46, 0x4f, 0x52, 0x43, 0x45, 0x4d, 0x45, 0x4e, 0x54,
	0x5f, 0x53, 0x45, 0x43, 0x55, 0x52, 0x49, 0x54, 0x59, 0x10, 0x6e, 0x12, 0x15, 0x0a, 0x11, 0x43,
	0x4f, 0x5f, 0x4c, 0x45, 0x47, 0x41, 0x4c, 0x5f, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x53,
	0x10, 0x6f, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x4f, 0x5f, 0x4d, 0x49, 0x4c, 0x49, 0x54, 0x41, 0x52,
	0x59, 0x10, 0x70, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x4f, 0x5f, 0x4e, 0x4f, 0x54, 0x41, 0x52, 0x59,
	0x5f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x52, 0x41, 0x52, 0x10, 0x71, 0x12, 0x17, 0x0a, 0x13,
	0x43, 0x4f, 0x5f, 0x50, 0x52, 0x49, 0x56, 0x41, 0x54, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x45, 0x53,
	0x54, 0x4f, 0x52, 0x10, 0x72, 0x12, 0x22, 0x0a, 0x1e, 0x43, 0x4f, 0x5f, 0x50, 0x52, 0x4f, 0x46,
	0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x41, 0x4c, 0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x49, 0x53,
	0x54, 0x52, 0x41, 0x54, 0x49, 0x56, 0x45, 0x10, 0x73, 0x12, 0x1e, 0x0a, 0x1a, 0x43, 0x4f, 0x5f,
	0x50, 0x52, 0x4f, 0x46, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x41, 0x4c, 0x5f, 0x4d, 0x41, 0x4e,
	0x41, 0x47, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x74, 0x12, 0x19, 0x0a, 0x15, 0x43, 0x4f, 0x5f,
	0x50, 0x52, 0x4f, 0x46, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x41, 0x4c, 0x5f, 0x4f, 0x54, 0x48,
	0x45, 0x52, 0x10, 0x75, 0x12, 0x1d, 0x0a, 0x19, 0x43, 0x4f, 0x5f, 0x50, 0x52, 0x4f, 0x46, 0x45,
	0x53, 0x53, 0x49, 0x4f, 0x4e, 0x41, 0x4c, 0x5f, 0x54, 0x45, 0x43, 0x48, 0x4e, 0x49, 0x43, 0x41,
	0x4c, 0x10, 0x76, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x4f, 0x5f, 0x52, 0x45, 0x54, 0x49, 0x52, 0x45,
	0x44, 0x10, 0x77, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f, 0x5f, 0x53, 0x41, 0x4c, 0x45, 0x53, 0x10,
	0x78, 0x12, 0x14, 0x0a, 0x10, 0x43, 0x4f, 0x5f, 0x53, 0x45, 0x4c, 0x46, 0x5f, 0x45, 0x4d, 0x50,
	0x4c, 0x4f, 0x59, 0x45, 0x44, 0x10, 0x79, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x4f, 0x5f, 0x53, 0x54,
	0x55, 0x44, 0x45, 0x4e, 0x54, 0x10, 0x7a, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x4f, 0x5f, 0x54, 0x52,
	0x41, 0x4e, 0x53, 0x50, 0x4f, 0x52, 0x54, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x7b, 0x12, 0x11,
	0x0a, 0x0d, 0x43, 0x4f, 0x5f, 0x55, 0x4e, 0x45, 0x4d, 0x50, 0x4c, 0x4f, 0x59, 0x45, 0x44, 0x10,
	0x7c, 0x2a, 0x49, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x54, 0x61, 0x78,
	0x49, 0x44, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x54, 0x49, 0x44, 0x54, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a,
	0x09, 0x43, 0x54, 0x49, 0x44, 0x54, 0x5f, 0x53, 0x53, 0x4e, 0x10, 0x64, 0x12, 0x0e, 0x0a, 0x0a,
	0x43, 0x54, 0x49, 0x44, 0x54, 0x5f, 0x49, 0x54, 0x49, 0x4e, 0x10, 0x65, 0x2a, 0x67, 0x0a, 0x17,
	0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x73, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x52, 0x53, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b,
	0x43, 0x52, 0x53, 0x5f, 0x43, 0x49, 0x54, 0x49, 0x5a, 0x45, 0x4e, 0x10, 0x64, 0x12, 0x10, 0x0a,
	0x0c, 0x43, 0x52, 0x53, 0x5f, 0x52, 0x45, 0x53, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x10, 0x65, 0x12,
	0x14, 0x0a, 0x10, 0x43, 0x52, 0x53, 0x5f, 0x4e, 0x4f, 0x4e, 0x5f, 0x52, 0x45, 0x53, 0x49, 0x44,
	0x45, 0x4e, 0x54, 0x10, 0x66, 0x2a, 0xd9, 0x01, 0x0a, 0x12, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d,
	0x65, 0x72, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x13, 0x0a, 0x0f,
	0x43, 0x49, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x49, 0x54, 0x5f, 0x49, 0x4e, 0x48, 0x45, 0x52, 0x49, 0x54,
	0x41, 0x4e, 0x43, 0x45, 0x10, 0x64, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x49, 0x54, 0x5f, 0x53, 0x41,
	0x4c, 0x41, 0x52, 0x59, 0x10, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x49, 0x54, 0x5f, 0x53, 0x41,
	0x4c, 0x45, 0x5f, 0x4f, 0x46, 0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x41, 0x4e, 0x59, 0x10, 0x66, 0x12,
	0x18, 0x0a, 0x14, 0x43, 0x49, 0x54, 0x5f, 0x53, 0x41, 0x4c, 0x45, 0x5f, 0x4f, 0x46, 0x5f, 0x50,
	0x52, 0x4f, 0x50, 0x45, 0x52, 0x54, 0x59, 0x10, 0x67, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x49, 0x54,
	0x5f, 0x49, 0x4e, 0x56, 0x45, 0x53, 0x54, 0x4d, 0x45, 0x4e, 0x54, 0x53, 0x10, 0x68, 0x12, 0x16,
	0x0a, 0x12, 0x43, 0x49, 0x54, 0x5f, 0x4c, 0x49, 0x46, 0x45, 0x5f, 0x49, 0x4e, 0x53, 0x55, 0x52,
	0x41, 0x4e, 0x43, 0x45, 0x10, 0x69, 0x12, 0x1a, 0x0a, 0x16, 0x43, 0x49, 0x54, 0x5f, 0x44, 0x49,
	0x56, 0x4f, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x45, 0x54, 0x54, 0x4c, 0x45, 0x4d, 0x45, 0x4e, 0x54,
	0x10, 0x6a, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x49, 0x54, 0x5f, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10,
	0x6b, 0x42, 0x23, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x77, 0x69, 0x73, 0x65, 0x63, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_consumer_proto_rawDescOnce sync.Once
	file_consumer_proto_rawDescData = file_consumer_proto_rawDesc
)

func file_consumer_proto_rawDescGZIP() []byte {
	file_consumer_proto_rawDescOnce.Do(func() {
		file_consumer_proto_rawDescData = protoimpl.X.CompressGZIP(file_consumer_proto_rawDescData)
	})
	return file_consumer_proto_rawDescData
}

var file_consumer_proto_enumTypes = make([]protoimpl.EnumInfo, 6)
var file_consumer_proto_goTypes = []interface{}{
	(ConsumerKYCStatus)(0),       // 0: wise.protobuf.ConsumerKYCStatus
	(ConsumerStatus)(0),          // 1: wise.protobuf.ConsumerStatus
	(ConsumerOccupation)(0),      // 2: wise.protobuf.ConsumerOccupation
	(ConsumerTaxIDType)(0),       // 3: wise.protobuf.ConsumerTaxIDType
	(ConsumerResidencyStatus)(0), // 4: wise.protobuf.ConsumerResidencyStatus
	(ConsumerIncomeType)(0),      // 5: wise.protobuf.ConsumerIncomeType
}
var file_consumer_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_consumer_proto_init() }
func file_consumer_proto_init() {
	if File_consumer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_consumer_proto_rawDesc,
			NumEnums:      6,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_consumer_proto_goTypes,
		DependencyIndexes: file_consumer_proto_depIdxs,
		EnumInfos:         file_consumer_proto_enumTypes,
	}.Build()
	File_consumer_proto = out.File
	file_consumer_proto_rawDesc = nil
	file_consumer_proto_goTypes = nil
	file_consumer_proto_depIdxs = nil
}
