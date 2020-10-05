package bank

import (
	"time"
)

type BusinessTaxIDType string

const (
	BusinessTaxIDTypeSSN = BusinessTaxIDType("ssn")
	BusinessTaxIDTypeEIN = BusinessTaxIDType("ein")
)

type BusinessEntity string

const (
	BusinessEntitySoleProprietor              = BusinessEntity("soleProprietor")
	BusinessEntityAssociation                 = BusinessEntity("association")
	BusinessEntityProfessionalAssociation     = BusinessEntity("professionalAssociation")
	BusinessEntitySingleMemberLLC             = BusinessEntity("singleMemberLLC")
	BusinessEntityLimitedLiabilityCompany     = BusinessEntity("limitedLiabilityCompany")
	BusinessEntityGeneralPartnership          = BusinessEntity("generalPartnership")
	BusinessEntityLimitedPartnership          = BusinessEntity("limitedPartnership")
	BusinessEntityLimitedLiabilityPartnership = BusinessEntity("limitedLiabilityPartnership")
	BusinessEntityProfessionalCorporation     = BusinessEntity("professionalCorporation")
	BusinessEntityUnlistedCorporation         = BusinessEntity("unlistedCorporation")
)

// Deprecated
type BusinessType = BusinessEntity

type BusinessIndustry string

const (
	BusinessIndustryHotelMotel                = BusinessIndustry("hotelMotel")
	BusinessIndustryOtherFoodServices         = BusinessIndustry("otherFoodServices")
	BusinessIndustryRestaurants               = BusinessIndustry("restaurants")
	BusinessIndustryArtPhotography            = BusinessIndustry("artPhotography")
	BusinessIndustryOtherArtsEntertainment    = BusinessIndustry("artsEntertainment")
	BusinessIndustryFitnessSportsCenters      = BusinessIndustry("fitnessSportsCenters")
	BusinessIndustrySportsTeamsClubs          = BusinessIndustry("sportsTeamsClubs")
	BusinessIndustryConstruction              = BusinessIndustry("construction")
	BusinessIndustryBuildingMaterialsHardware = BusinessIndustry("buildingMaterialsHardware")
	BusinessIndustryOtherTradeContractor      = BusinessIndustry("otherTradeContractor")
	BusinessIndustryPlumbingHVAC              = BusinessIndustry("plumbingHVAC")
	BusinessIndustryHealthServices            = BusinessIndustry("healthServices")
	BusinessIndustryOtherEducationServices    = BusinessIndustry("otherEducationServices")
	BusinessIndustryOtherHealthFitness        = BusinessIndustry("otherHealthFitness")
	BusinessIndustryAccountingTaxPrep         = BusinessIndustry("accountingTaxPrep")
	BusinessIndustryRealEstate                = BusinessIndustry("realEstate")
	BusinessIndustryHomeFurnishing            = BusinessIndustry("homeFurnishing")
	BusinessIndustryBeautyOrBarberShops       = BusinessIndustry("beautyOrBarberShops")
	BusinessIndustryCarWash                   = BusinessIndustry("carWash")
	BusinessIndustryComputerServiceRepair     = BusinessIndustry("computerServiceRepair")
	BusinessIndustryFreelanceProfessional     = BusinessIndustry("freelanceProfessional")
	BusinessIndustryLandscapeServices         = BusinessIndustry("landscapeServices")
	BusinessIndustryLegalServices             = BusinessIndustry("legalServices")
	BusinessIndustryMassageTanningServices    = BusinessIndustry("massageTanningServices")
	BusinessIndustryOtherProfessionalServices = BusinessIndustry("otherProfessionalServices")
	BusinessIndustryAutoDealers               = BusinessIndustry("autoDealers")
	BusinessIndustryOnlineRetailer            = BusinessIndustry("onlineRetailer")
	BusinessIndustryRetail                    = BusinessIndustry("retail")
	BusinessIndustryGasolineServiceStation    = BusinessIndustry("gasolineServiceStation")
	BusinessIndustryOtherTransportServices    = BusinessIndustry("otherTransportServices")
	BusinessIndustryOtherTravelServices       = BusinessIndustry("otherTravelServices")
	BusinessIndustryParkingGarages            = BusinessIndustry("parkingGarages")
	BusinessIndustryTaxi                      = BusinessIndustry("taxi")
	BusinessIndustryTravelAgency              = BusinessIndustry("travelAgency")
	BusinessIndustryTruckingShipping          = BusinessIndustry("truckingShipping")
	BusinessIndustryWholesale                 = BusinessIndustry("wholesale")
	BusinessIndustryWarehouseDistribution     = BusinessIndustry("warehouseDistribution")
	BusinessIndustryOtherAccomodation         = BusinessIndustry("otherAccomodations")
	BusinessIndustryRestaurantsWithCash       = BusinessIndustry("restaurantsCash")
	BusinessIndustryAnimalFarmingProduction   = BusinessIndustry("animalFarmingProduction")
	BusinessIndustryCropFarming               = BusinessIndustry("cropFarming")
	BusinessIndustryForestry                  = BusinessIndustry("forestryActivities")
	BusinessIndustryFishingHunting            = BusinessIndustry("fishingHuntingTrapping")
	BusinessIndustryOtherFarmingHunting       = BusinessIndustry("otherAgricultureForestryFishing")
	BusinessIndustryMuseums                   = BusinessIndustry("museumsHistoricalSites")
	BusinessIndustryHospitals                 = BusinessIndustry("hospitals")
	BusinessIndustryCollegeUniversitySchools  = BusinessIndustry("collegesUniversitiesSchools")
	BusinessIndustryBank                      = BusinessIndustry("bankFinancialInstitution")
	BusinessIndustryFinancialInvestments      = BusinessIndustry("financialInvestments")
	BusinessIndustryFundsTrustsOther          = BusinessIndustry("fundsTrustsOther")
	BusinessIndustryInsurance                 = BusinessIndustry("insurance")
	BusinessIndustryMoneyTransferRemittance   = BusinessIndustry("moneyTransferRemittance")
	BusinessIndustryPrivateInvestment         = BusinessIndustry("privateInvestmentCompanies")
	BusinessIndustryOtherManufacturing        = BusinessIndustry("otherManufacturing")
	BusinessIndustryIndustrialMachinery       = BusinessIndustry("industrialCommercialMachinery")
	BusinessIndustryEmploymentServices        = BusinessIndustry("employmentServices")
	BusinessIndustryGovernmentAgency          = BusinessIndustry("governmentAgency")
	BusinessIndustryNonGovernment             = BusinessIndustry("nonGovernmentOrganization")
	BusinessIndustryReligiousOrganization     = BusinessIndustry("religiousOrganization")
	BusinessIndustryUnions                    = BusinessIndustry("unions")
	BusinessIndustryRetailJeweler             = BusinessIndustry("retailJewelerDiamondsGemsGold")
	BusinessIndustryRetailWithCash            = BusinessIndustry("retailCash")
	BusinessIndustryUsedClothesDealer         = BusinessIndustry("usedClothesDealers")
	BusinessIndustryTourOperator              = BusinessIndustry("tourOperator")
	BusinessIndustryWholesaleJeweler          = BusinessIndustry("wholesaleJeweler")

	// Not supported
	BusinessIndustryCasinoHotel            = BusinessIndustry("casinoHotel")
	BusinessIndustryCasinoGaming           = BusinessIndustry("casinoGamblingGaming")
	BusinessIndustryRaceTrack              = BusinessIndustry("raceTrack")
	BusinessIndustryCheckCasher            = BusinessIndustry("checkCasher")
	BusinessIndustryCollectionsAgency      = BusinessIndustry("collectionAgencies")
	BusinessIndustryCurrencyExchange       = BusinessIndustry("currencyExchangers")
	BusinessIndustryCigaretteManufacturing = BusinessIndustry("cigaretteManufacturing")
	BusinessIndustryPrivateATM             = BusinessIndustry("privateATM")
	BusinessIndustryConsulateEmbassy       = BusinessIndustry("consulateEmbassy")
	BusinessIndustryBeerWineLiquorStores   = BusinessIndustry("beerWineLiquorStores")
	BusinessIndustryPawnShop               = BusinessIndustry("pawnShop")
)

type BusinessOperationType string

const (
	BusinessOperationTypeLocal            = BusinessOperationType("local")
	BusinessOperationTypeForeignWithLocal = BusinessOperationType("foreignWithLocal")
	BusinessOperationTypeForeign          = BusinessOperationType("foreign")
)

type BusinessID string

type CreateBusinessRequest struct {
	BusinessID         BusinessID              `json:"businessId"`
	LegalName          string                  `json:"legalName"`
	DBA                []string                `json:"dba"`
	TaxIDType          BusinessTaxIDType       `json:"taxIdType"`
	TaxID              string                  `json:"taxId"`
	EntityType         BusinessEntity          `json:"entityType"`
	IndustryType       BusinessIndustry        `json:"industryType"`
	Phone              string                  `json:"phone"`
	Email              string                  `json:"email"`
	Purpose            string                  `json:"purpose"`
	ExpectedActivities []ExpectedActivity      `json:"expectedActivities"`
	Members            []BusinessMemberRequest `json:"members"`
	LegalAddress       AddressRequest          `json:"legalAddress"`
	MailingAddress     *AddressRequest         `json:"mailingAddress"`
	HeadquarterAddress AddressRequest          `json:"headquarterAddress "`
	OriginCountry      Country                 `json:"originCountry"`
	OriginState        string                  `json:"originState"`
	OriginDate         time.Time               `json:"originDate"`
	OperationType      BusinessOperationType   `json:"operationType"`
	EntityFormation    *EntityFormationRequest `json:"entityFormation"`
}

type UpdateBusinessRequest struct {
	LegalName       string                  `json:"legalName"`
	DBA             []string                `json:"dba"`
	BusinessID      BusinessID              `json:"businessId"`
	TaxIDType       BusinessTaxIDType       `json:"taxIdType"`
	TaxID           string                  `json:"taxId"`
	EntityType      BusinessEntity          `json:"entityType"`
	IndustryType    BusinessIndustry        `json:"industryType"`
	Purpose         string                  `json:"purpose"`
	OriginCountry   Country                 `json:"originCountry"`
	OriginState     string                  `json:"originState"`
	OriginDate      time.Time               `json:"originDate"`
	OperationType   BusinessOperationType   `json:"operationType"`
	EntityFormation *EntityFormationRequest `json:"entityFormation"`
}

type BusinessMemberTitle string

func (t BusinessMemberTitle) String() string {
	return string(t)
}

const (
	BusinessMemberTitleCEO       = BusinessMemberTitle("chiefExecutiveOfficer")
	BusinessMemberTitleCFO       = BusinessMemberTitle("chiefFinancialOfficer")
	BusinessMemberTitleCOO       = BusinessMemberTitle("chiefOperatingOfficer")
	BusinessMemberTitlePresident = BusinessMemberTitle("president")
	BusinessMemberTitleVP        = BusinessMemberTitle("vicePresident")
	BusinessMemberTitleSVP       = BusinessMemberTitle("seniorVicePresident")
	BusinessMemberTitleTreasurer = BusinessMemberTitle("treasurer")
	BusinessMemberTitleSecretary = BusinessMemberTitle("secretary")
	BusinessMemberTitleGP        = BusinessMemberTitle("generalPartner")
	BusinessMemberTitleManager   = BusinessMemberTitle("manager")
	BusinessMemberTitleMember    = BusinessMemberTitle("member")
	BusinessMemberTitleOwner     = BusinessMemberTitle("owner")
	BusinessMemberTitleOther     = BusinessMemberTitle("other")
)

// Requirements for business members
// In Create a business: POST, Step 3, the example JSON includes a members array.
// That array contains two JSON structures that specify information for two members.
// The two-member structure means that when creating a business record, the members array:
//		- Must include one and only one “user_type”: “control_manager”
//		- Should include an additional “user_type”: “owner” for all persons who own 25 percent or more of the business entity.
type BusinessMemberRequest struct {
	ConsumerID           ConsumerID          `json:"consumerId"`
	IsControllingManager bool                `json:"isControllingManager"`
	Ownership            int                 `json:"ownership,omitempty"`
	Title                BusinessMemberTitle `json:"title"`
	TitleDesc            *string             `json:"titleDesc"`
}

type BusinessFormationDocument string

const (
	BusinessFormationDocumentArticlesOfIncorporation = BusinessFormationDocument("articlesOfIncorporation")  // Corporate charter
	BusinessFormationDocumentArticlesOfOrganization  = BusinessFormationDocument("articlesOfOrganization")   // Initial statements to form an LLC
	BusinessFormationDocumentPartnershipCertificate  = BusinessFormationDocument("certificateOfPartnership") // Certificate of partnership
	BusinessFormationDocumentPartnershipAgreement    = BusinessFormationDocument("partnershipAgreement")     // Partnership agreement
	BusinessFormationDocumentCertificateOfFormation  = BusinessFormationDocument("certificateOfFormation")   // Certificate of LLC formation
	BusinessFormationDocumentDriversLicense          = BusinessFormationDocument("driversLicense")           // Drivers license of owner (Sole Prop)
)

type BusinessIdentityDocument string

const (
	BusinessIdentityDocumentArticlesOfIncorporation = BusinessIdentityDocument("articlesOfIncorporation")  // Corporate charter
	BusinessIdentityDocumentArticlesOfOrganization  = BusinessIdentityDocument("articlesOfOrganization")   // Initial statements to form an LLC
	BusinessIdentityDocumentAssumedNameCertificate  = BusinessIdentityDocument("assumedNameCertificate")   // Certificate for DBA
	BusinessIdentityDocumentBusinessLicense         = BusinessIdentityDocument("businessLicense")          // Business license
	BusinessIdentityDocumentPartnershipCertificate  = BusinessIdentityDocument("certificateOfPartnership") // Certificate of partnership
	BusinessIdentityDocumentPartnershipAgreement    = BusinessIdentityDocument("partnershipAgreement")     // Partnership agreement
	BusinessIdentityDocumentCertificateOfFormation  = BusinessIdentityDocument("certificateOfFormation")   // Certificate of LLC formation
	BusinessIdentityDocumentDriversLicense          = BusinessIdentityDocument("driversLicense")           // Drivers license of owner (Sole Prop)
	BusinessIdentityDocumentOther                   = BusinessIdentityDocument("other")                    // Other business documents
)

type EntityFormationRequest struct {
	DocumentType   BusinessIdentityDocument `json:"docType"`
	Number         string                   `json:"number"`
	IssueDate      time.Time                `json:"issueDate"`
	ExpirationDate time.Time                `json:"expirationDate"`
}

type EntityFormationResponse EntityFormationRequest

type BusinessDocument struct {
	DocumentType   BusinessIdentityDocument `json:"docType"`
	Number         string                   `json:"number"`
	Issuer         string                   `json:"issuer"`
	IssueDate      time.Time                `json:"issueDate"`
	IssueState     *string                  `json:"state"`
	IssueCountry   *Country                 `json:"country"`
	ExpirationDate time.Time                `json:"expirationDate"`
}

type BusinessBankID string

type IdentityStatusBusinessResponse struct {
	BusinessID BusinessID     `json:"businessId"`
	BankID     BusinessBankID `json:"bankID"`
	KYC        KYCResponse    `json:"kyc"`
}

type MemberBankID string

type CreateBusinessMemberResponse struct {
	ConsumerID     ConsumerID     `json:"consumerId"`
	MemberBankID   MemberBankID   `json:"memberBankId"`
	ConsumerBankID ConsumerBankID `json:"consumerBankId"`
}

type BusinessMemberResponse struct {
	ConsumerID     ConsumerID     `json:"consumerId"`
	BusinessID     BusinessID     `json:"businessId"`
	BankName       ProviderName   `json:"bankName"`
	MemberBankID   MemberBankID   `json:"memberBankId"`
	ConsumerBankID ConsumerBankID `json:"consumerBankId"`
	KYCStatus      KYCStatus      `json:"kycStatus"`
	Created        time.Time      `json:"created"`
	Updated        time.Time      `json:"updated"`
}
