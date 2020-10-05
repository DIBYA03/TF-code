/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business services
package business

const (
	EntityTypeSoleProprietor              = "soleProprietor"              // Sole proprietorship
	EntityTypeAssociation                 = "association"                 // Association
	EntityTypeyProfessionalAssociation    = "professionalAssociation"     // Professional association
	EntityTypeSingleMemberLLC             = "singleMemberLLC"             // Single member llc
	EntityTypeLimitedLiabilityCompany     = "limitedLiabilityCompany"     // Limited liability company
	EntityTypeGeneralPartnership          = "generalPartnership"          // General partnership
	EntityTypeLimitedPartnership          = "limitedPartnership"          // Limited partnership
	EntityTypeLimitedLiabilityPartnership = "limitedLiabilityPartnership" // Limited liability partnership
	EntityTypeProfessionalCorporation     = "professionalCorporation"     // Professional corporation
	EntityTypeUnlistedCorporation         = "unlistedCorporation"         // Unlisted Corporation (S-Corp & C-Corp)
)

const (
	IndustryTypeHotelMotel                = "hotelMotel"                // Non-casino hotel/motel operator (traveler accommodation)
	IndustryTypeOtherFoodServices         = "otherFoodServices"         // Other food services
	IndustryTypeRestaurants               = "restaurants"               // Restaurants
	IndustryTypeArtPhotography            = "artPhotography"            // Arts and photography
	IndustryTypeOtherArtsEntertainment    = "artsEntertainment"         // Arts and entertainment
	IndustryTypeFitnessCenter             = "fitnessSportsCenters"      // Fitness and Recreational Sports Centers
	IndustryTypeSportsTeamsClubs          = "sportsTeamsClubs"          // Sports teams and clubs
	IndustryTypeConstruction              = "construction"              // General construction
	IndustryTypeBuildingMaterialsHardware = "buildingMaterialsHardware" // Materials and hardware store
	IndustryTypeOtherTradeContractor      = "otherTradeContractor"      // Other trade contractor
	IndustryTypePlumbingHVAC              = "plumbingHVAC"              // Plumsbing and HVAC
	IndustryTypeHealthServices            = "healthServices"            // Health related services
	IndustryTypeOtherEducationServices    = "otherEducationServices"    // Other education services
	IndustryTypeOtherHealthFitness        = "otherHealthFitness"        // Other health and fitness
	IndustryTypeAccountingTaxPrep         = "accountingTaxPrep"         // Accounting and tax prep
	IndustryTypeRealEstate                = "realEstate"                // Real estate agent or manager
	IndustryTypeHomeFurnishing            = "homeFurnishing"            // Home furnishing
	IndustryTypeBeautyOrBarberShops       = "beautyOrBarberShops"       // Beauty or barber shops
	IndustryTypeCarWash                   = "carWash"                   // Car wash
	IndustryTypeComputerServiceRepair     = "computerServiceRepair"     // Computer service or repair
	IndustryTypeFreelanceProfessional     = "freelanceProfessional"     // Freelance professional
	IndustryTypeLandscapeServices         = "landscapeServices"         // Landscaping and gardening services
	IndustryTypeLegalServices             = "legalServices"             // Legal services and counseling.
	IndustryTypeMassageTanningServices    = "massageTanningServices"    // Massage and tanning salons.
	IndustryTypeOtherProfessionalServices = "otherProfessionalServices" // Other professional services
	IndustryTypeAutoDealers               = "autoDealers"               // Automobile dealers
	IndustryTypeOnlineRetailer            = "onlineRetailer"            // E-Commerce retailer
	IndustryTypeRetail                    = "retail"                    // Miscellaneous retailor
	IndustryTypeGasServiceStation         = "gasolineServiceStation"    // Gasoline station
	IndustryTypeTransportationServices    = "otherTransportServices"    // Other ground passenger transportation
	IndustryTypeOtherTravelServices       = "otherTravelServices"       // Other travel related services
	IndustryTypeParkingGarages            = "parkingGarages"            // Parking lots and garages
	IndustryTypeTaxi                      = "taxi"                      // Taxi service
	IndustryTypeTravelAgency              = "travelAgency"              // Professional travel arrangement and reservation services
	IndustryTypeTruckingShipping          = "truckingShipping"          // Shipping and truck transportation
	IndustryTypeWholesale                 = "wholesale"                 // Wholesaler
	IndustryTypeWarehouseDistribution     = "warehouseDistribution"     // General warehousing and storage
	IndustryTypeOtherAccomodation         = "otherAccomodations"
	IndustryTypeRestaurantsWithCash       = "restaurantsCash"
	IndustryTypeAnimalFarmingProduction   = "animalFarmingProduction"
	IndustryTypeCropFarming               = "cropFarming"
	IndustryTypeForestry                  = "forestryActivities"
	IndustryTypeFishingHunting            = "fishingHuntingTrapping"
	IndustryTypeOtherFarmingHunting       = "otherAgricultureForestryFishing"
	IndustryTypeMuseums                   = "museumsHistoricalSites"
	IndustryTypeHospitals                 = "hospitals"
	IndustryTypeCollegeUniversitySchools  = "collegesUniversitiesSchools"
	IndustryTypeBank                      = "bankFinancialInstitution"
	IndustryTypeFinancialInvestments      = "financialInvestments"
	IndustryTypeFundsTrustsOther          = "fundsTrustsOther"
	IndustryTypeInsurance                 = "insurance"
	IndustryTypeMoneyTransferRemittance   = "moneyTransferRemittance"
	IndustryTypePrivateInvestment         = "privateInvestmentCompanies"
	IndustryTypeOtherManufacturing        = "otherManufacturing"
	IndustryTypeIndustrialMachinery       = "industrialCommercialMachinery"
	IndustryTypeEmploymentServices        = "employmentServices"
	IndustryTypeGovernmentAgency          = "governmentAgency"
	IndustryTypeNonGovernment             = "nonGovernmentOrganization"
	IndustryTypeReligiousOrganization     = "religiousOrganization"
	IndustryTypeUnions                    = "unions"
	IndustryTypeRetailJeweler             = "retailJewelerDiamondsGemsGold"
	IndustryTypeRetailWithCash            = "retailCash"
	IndustryTypeUsedClothesDealer         = "usedClothesDealers"
	IndustryTypeTourOperator              = "tourOperator"
	IndustryTypeWholesaleJeweler          = "wholesaleJeweler"

	// Not supported
	IndustryTypeCasinoHotel            = "casinoHotel"
	IndustryTypeCasinoGaming           = "casinoGamblingGaming"
	IndustryTypeRaceTrack              = "raceTrack"
	IndustryTypeCheckCasher            = "checkCasher"
	IndustryTypeCollectionsAgency      = "collectionAgencies"
	IndustryTypeCurrencyExchange       = "currencyExchangers"
	IndustryTypeCigaretteManufacturing = "cigaretteManufacturing"
	IndustryTypePrivateATM             = "privateATM"
	IndustryTypeConsulateEmbassy       = "consulateEmbassy"
	IndustryTypeBeerWineLiquorStores   = "beerWineLiquorStores"
	IndustryTypePawnShop               = "pawnShop"
)

const (
	OperationTypeLocal            = "local"
	OperationTypeForeignWithLocal = "foreignWithLocal"
	OperationTypeForeign          = "foreign"
)

const (
	BusinessAccessTypeView       = "view"       // Allows view only access
	BusinessAccessTypeAccounting = "accounting" // Allows viewing and exporting reports
	BusinessAccessTypeAdmin      = "admin"      // Allows user to update account
)

const (
	BusinessAccessRoleAttorney = "attorney"
	BusinessAccessRoleOfficer  = "officer"
	BusinessAccessRoleEmployee = "employee"
	BusinessAccessRoleOther    = "other"
)

type BusinessAccess struct {
	Id         string `json:"accessId"`
	BusinessId string `json:"businessId"`
	UserId     string `json:"userId"`
	AccessType string `json:"accessType"`
	AccessRole string `json:"accessRole"`
}
