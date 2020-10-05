package bbva

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type BusinessService struct {
	request bank.APIRequest
	client  *client
}

func (b *businessBank) BusinessEntityService(request bank.APIRequest) bank.BusinessService {
	return &BusinessService{
		request: request,
		client:  b.client,
	}
}

type BusinessTINType string

const (
	BusinessTINTypeSSN = "ssn"
	BusinessTINTypeEIN = "ein"
)

var partnerTaxIDToMap = map[BusinessTINType]bank.BusinessTaxIDType{
	BusinessTINTypeSSN: bank.BusinessTaxIDTypeSSN,
	BusinessTINTypeEIN: bank.BusinessTaxIDTypeEIN,
}

var partnerTaxIDFromMap = map[bank.BusinessTaxIDType]BusinessTINType{
	bank.BusinessTaxIDTypeSSN: BusinessTINTypeSSN,
	bank.BusinessTaxIDTypeEIN: BusinessTINTypeEIN,
}

type BusinessEntity string

const (
	BusinessEntitySoleProprietor              = BusinessEntity("sole_proprietor")
	BusinessEntityAssociation                 = BusinessEntity("association")
	BusinessEntityProfessionalAssociation     = BusinessEntity("professional_association")
	BusinessEntityLimitedLiabilityCompany     = BusinessEntity("limited_liability_company")
	BusinessEntitySingleMemberLLC             = BusinessEntity("single_member_llc")
	BusinessEntityPartnership                 = BusinessEntity("partnership")
	BusinessEntityLimitedPartnership          = BusinessEntity("limited_partnership")
	BusinessEntityLimitedLiabilityPartnership = BusinessEntity("limited_liability_partnership")
	BusinessEntityProfessionalCorporation     = BusinessEntity("professional_corporation")
	BusinessEntityUnlistedCorporation         = BusinessEntity("unlisted_corporation")
	BusinessEntityCommercial                  = BusinessEntity("commercial")

	// Not supported
	BusinessEntityPublicCorporation       = BusinessEntity("publicly_traded")
	BusinessEntityRevocableTrust          = BusinessEntity("revocable_trust")
	BusinessEntityIrrevocableTrust        = BusinessEntity("irrevocable_trust")
	BusinessEntityEstate                  = BusinessEntity("estate")
	BusinessEntityBankruptcy              = BusinessEntity("bankruptcy")
	BusinessEntityDebtorInPossession      = BusinessEntity("debtor_in_possession")
	BusinessEntityPublicFunds             = BusinessEntity("public_funds")
	BusinessEntityPubliclyTraded          = BusinessEntity("publicly_traded")
	BusinessEntityUSGovernment            = BusinessEntity("us_government")
	BusinessEntityStateAndLocalGovernment = BusinessEntity("state_and_local_government")
	BusinessEntityNonProfit               = BusinessEntity("non_profit")
	BusinessEntityEmployeeBenefitPlan     = BusinessEntity("employee_benefit_plan")
	BusinessEntityTrust                   = BusinessEntity("trust")
)

var partnerEntityToMap = map[BusinessEntity]bank.BusinessEntity{
	BusinessEntitySoleProprietor:              bank.BusinessEntitySoleProprietor,
	BusinessEntityAssociation:                 bank.BusinessEntityAssociation,
	BusinessEntityProfessionalAssociation:     bank.BusinessEntityProfessionalAssociation,
	BusinessEntitySingleMemberLLC:             bank.BusinessEntitySingleMemberLLC,
	BusinessEntityLimitedLiabilityCompany:     bank.BusinessEntityLimitedLiabilityCompany,
	BusinessEntityPartnership:                 bank.BusinessEntityGeneralPartnership,
	BusinessEntityLimitedPartnership:          bank.BusinessEntityLimitedPartnership,
	BusinessEntityLimitedLiabilityPartnership: bank.BusinessEntityLimitedLiabilityPartnership,
	BusinessEntityProfessionalCorporation:     bank.BusinessEntityProfessionalCorporation,
	BusinessEntityUnlistedCorporation:         bank.BusinessEntityUnlistedCorporation,
}

var partnerEntityFromMap = map[bank.BusinessEntity]BusinessEntity{
	bank.BusinessEntitySoleProprietor:              BusinessEntitySoleProprietor,
	bank.BusinessEntityAssociation:                 BusinessEntityAssociation,
	bank.BusinessEntityProfessionalAssociation:     BusinessEntityProfessionalAssociation,
	bank.BusinessEntitySingleMemberLLC:             BusinessEntitySingleMemberLLC,
	bank.BusinessEntityLimitedLiabilityCompany:     BusinessEntityLimitedLiabilityCompany,
	bank.BusinessEntityGeneralPartnership:          BusinessEntityPartnership,
	bank.BusinessEntityLimitedPartnership:          BusinessEntityLimitedPartnership,
	bank.BusinessEntityLimitedLiabilityPartnership: BusinessEntityLimitedLiabilityPartnership,
	bank.BusinessEntityProfessionalCorporation:     BusinessEntityProfessionalCorporation,
	bank.BusinessEntityUnlistedCorporation:         BusinessEntityUnlistedCorporation,
}

type BusinessIndustry string

const (
	BusinessIndustryHotelMotel                = BusinessIndustry("hotel_motel")
	BusinessIndustryOtherFoodServices         = BusinessIndustry("other_food_services")
	BusinessIndustryRestaurants               = BusinessIndustry("restaurants_non_cash_int")
	BusinessIndustryArtPhotography            = BusinessIndustry("art_and_photography")
	BusinessIndustryOtherArtsEntertainment    = BusinessIndustry("other_arts_entertainment_and_rec")
	BusinessIndustryFitnessSportsCenters      = BusinessIndustry("fitness_and_recreational_sports_centers")
	BusinessIndustrySportsTeamsClubs          = BusinessIndustry("sports_teams_and_clubs")
	BusinessIndustryConstruction              = BusinessIndustry("construction")
	BusinessIndustryBuildingMaterialsHardware = BusinessIndustry("building_materials_or_hardware")
	BusinessIndustryOtherTradeContractor      = BusinessIndustry("other_specialty_trade_contractor")
	BusinessIndustryPlumbingHVAC              = BusinessIndustry("plumbing_heating_and_air_conditioning_contractors")
	BusinessIndustryHealthServices            = BusinessIndustry("health_services")
	BusinessIndustryOtherEducationServices    = BusinessIndustry("other_education_related_services")
	BusinessIndustryOtherHealthFitness        = BusinessIndustry("other_health_and_fitness_services")
	BusinessIndustryAccountingTaxPrep         = BusinessIndustry("accounting_and_tax_preparation")
	BusinessIndustryRealEstate                = BusinessIndustry("real_estate")
	BusinessIndustryHomeFurnishing            = BusinessIndustry("home_furnishings_or_furniture")
	BusinessIndustryBeautyOrBarberShops       = BusinessIndustry("beauty_or_barber_shops")
	BusinessIndustryCarWash                   = BusinessIndustry("car_wash")
	BusinessIndustryComputerServiceRepair     = BusinessIndustry("computer_consulting_service_and_repair")
	BusinessIndustryFreelanceProfessional     = BusinessIndustry("freelance_professional_services")
	BusinessIndustryLandscapeServices         = BusinessIndustry("landscape_lawn_and_garden_services")
	BusinessIndustryLegalServices             = BusinessIndustry("legal_services_including_law_firms")
	BusinessIndustryMassageTanningServices    = BusinessIndustry("massage_and_tanning_services")
	BusinessIndustryOtherProfessionalServices = BusinessIndustry("other_professional_services")
	BusinessIndustryAutoDealers               = BusinessIndustry("auto_dealers")
	BusinessIndustryOnlineRetailer            = BusinessIndustry("online_retailer")
	BusinessIndustryRetail                    = BusinessIndustry("retail_non_cash_int")
	BusinessIndustryGasolineServiceStation    = BusinessIndustry("gasoline_service_station")
	BusinessIndustryOtherTransportServices    = BusinessIndustry("other_trasportation_services") // Should be other_transportation_services
	BusinessIndustryOtherTravelServices       = BusinessIndustry("other_travel_services")
	BusinessIndustryParkingGarages            = BusinessIndustry("parking_garages")
	BusinessIndustryTaxi                      = BusinessIndustry("transportation_taxi")
	BusinessIndustryTravelAgency              = BusinessIndustry("travel_agency")
	BusinessIndustryTruckingShipping          = BusinessIndustry("trucking_shipping")
	BusinessIndustryWholesale                 = BusinessIndustry("wholesale")
	BusinessIndustryWarehouseDistribution     = BusinessIndustry("warehouse_or_distribution")
	BusinessIndustryOtherAccomodation         = BusinessIndustry("other_accomodation_services")
	BusinessIndustryRestaurantsWithCash       = BusinessIndustry("restaurants_cash_int")
	BusinessIndustryAnimalFarmingProduction   = BusinessIndustry("animal_farming_or_production")
	BusinessIndustryCropFarming               = BusinessIndustry("crop_farming")
	BusinessIndustryForestry                  = BusinessIndustry("forestry_activities")
	BusinessIndustryFishingHunting            = BusinessIndustry("fishing_hunting_and_trapping")
	BusinessIndustryOtherFarmingHunting       = BusinessIndustry("other_agriculture_forestry_and_fishing")
	BusinessIndustryMuseums                   = BusinessIndustry("museums_historical_sites_and_similar_institutions")
	BusinessIndustryHospitals                 = BusinessIndustry("hospitals")
	BusinessIndustryCollegeUniversitySchools  = BusinessIndustry("colleges_universities_and_schools")
	BusinessIndustryBank                      = BusinessIndustry("bank_financial_institution")
	BusinessIndustryFinancialInvestments      = BusinessIndustry("financial_investments")
	BusinessIndustryFundsTrustsOther          = BusinessIndustry("funds_trusts_and_other_financial_vehicles")
	BusinessIndustryInsurance                 = BusinessIndustry("insurance")
	BusinessIndustryMoneyTransferRemittance   = BusinessIndustry("money_transfertransmittal_or_remittance")
	BusinessIndustryPrivateInvestment         = BusinessIndustry("private_investment_companies")
	BusinessIndustryOtherManufacturing        = BusinessIndustry("all_other_manufacturing")
	BusinessIndustryIndustrialMachinery       = BusinessIndustry("industrial_or_commercial_machinery")
	BusinessIndustryEmploymentServices        = BusinessIndustry("employment_services")
	BusinessIndustryGovernmentAgency          = BusinessIndustry("government_agency")
	BusinessIndustryNonGovernment             = BusinessIndustry("non_government_organization")
	BusinessIndustryReligiousOrganization     = BusinessIndustry("religious_organization")
	BusinessIndustryUnions                    = BusinessIndustry("unions")
	BusinessIndustryRetailJeweler             = BusinessIndustry("retail_jeweler_diamonds_gems_gold")
	BusinessIndustryRetailWithCash            = BusinessIndustry("retail_cash_int")
	BusinessIndustryUsedClothesDealer         = BusinessIndustry("used_clothes_dealers")
	BusinessIndustryTourOperator              = BusinessIndustry("tour_operator")
	BusinessIndustryWholesaleJeweler          = BusinessIndustry("wholesale_jeweler")

	// Not supported
	BusinessIndustryCasinoHotel            = BusinessIndustry("casino_hotel")
	BusinessIndustryCasinoGaming           = BusinessIndustry("casino_gambling_or_gaming")
	BusinessIndustryRaceTrack              = BusinessIndustry("race_track")
	BusinessIndustryCheckCasher            = BusinessIndustry("check_casher")
	BusinessIndustryCollectionsAgency      = BusinessIndustry("collection_agencies")
	BusinessIndustryCurrencyExchange       = BusinessIndustry("currency_exchange_dealers")
	BusinessIndustryCigaretteManufacturing = BusinessIndustry("cigarette_manufacturing")
	BusinessIndustryPrivateATM             = BusinessIndustry("privately_owned_automated_teller_machine")
	BusinessIndustryConsulateEmbassy       = BusinessIndustry("consulate_embassy")
	BusinessIndustryBeerWineLiquorStores   = BusinessIndustry("beer_wine_and_liquor_stores")
	BusinessIndustryPawnShop               = BusinessIndustry("pawn_shop")
)

var partnerIndustryToMap = map[BusinessIndustry]bank.BusinessIndustry{
	BusinessIndustryHotelMotel:                bank.BusinessIndustryHotelMotel,
	BusinessIndustryOtherFoodServices:         bank.BusinessIndustryOtherFoodServices,
	BusinessIndustryRestaurants:               bank.BusinessIndustryRestaurants,
	BusinessIndustryArtPhotography:            bank.BusinessIndustryArtPhotography,
	BusinessIndustryOtherArtsEntertainment:    bank.BusinessIndustryOtherArtsEntertainment,
	BusinessIndustryFitnessSportsCenters:      bank.BusinessIndustryFitnessSportsCenters,
	BusinessIndustrySportsTeamsClubs:          bank.BusinessIndustrySportsTeamsClubs,
	BusinessIndustryConstruction:              bank.BusinessIndustryConstruction,
	BusinessIndustryBuildingMaterialsHardware: bank.BusinessIndustryBuildingMaterialsHardware,
	BusinessIndustryOtherTradeContractor:      bank.BusinessIndustryOtherTradeContractor,
	BusinessIndustryPlumbingHVAC:              bank.BusinessIndustryPlumbingHVAC,
	BusinessIndustryHealthServices:            bank.BusinessIndustryHealthServices,
	BusinessIndustryOtherEducationServices:    bank.BusinessIndustryOtherEducationServices,
	BusinessIndustryOtherHealthFitness:        bank.BusinessIndustryOtherHealthFitness,
	BusinessIndustryAccountingTaxPrep:         bank.BusinessIndustryAccountingTaxPrep,
	BusinessIndustryRealEstate:                bank.BusinessIndustryRealEstate,
	BusinessIndustryBeautyOrBarberShops:       bank.BusinessIndustryBeautyOrBarberShops,
	BusinessIndustryHomeFurnishing:            bank.BusinessIndustryHomeFurnishing,
	BusinessIndustryCarWash:                   bank.BusinessIndustryCarWash,
	BusinessIndustryComputerServiceRepair:     bank.BusinessIndustryComputerServiceRepair,
	BusinessIndustryFreelanceProfessional:     bank.BusinessIndustryFreelanceProfessional,
	BusinessIndustryLandscapeServices:         bank.BusinessIndustryLandscapeServices,
	BusinessIndustryLegalServices:             bank.BusinessIndustryLegalServices,
	BusinessIndustryMassageTanningServices:    bank.BusinessIndustryMassageTanningServices,
	BusinessIndustryOtherProfessionalServices: bank.BusinessIndustryOtherProfessionalServices,
	BusinessIndustryAutoDealers:               bank.BusinessIndustryAutoDealers,
	BusinessIndustryOnlineRetailer:            bank.BusinessIndustryOnlineRetailer,
	BusinessIndustryRetail:                    bank.BusinessIndustryRetail,
	BusinessIndustryGasolineServiceStation:    bank.BusinessIndustryGasolineServiceStation,
	BusinessIndustryOtherTransportServices:    bank.BusinessIndustryOtherTransportServices,
	BusinessIndustryOtherTravelServices:       bank.BusinessIndustryOtherTravelServices,
	BusinessIndustryParkingGarages:            bank.BusinessIndustryParkingGarages,
	BusinessIndustryTaxi:                      bank.BusinessIndustryTaxi,
	BusinessIndustryTravelAgency:              bank.BusinessIndustryTravelAgency,
	BusinessIndustryTruckingShipping:          bank.BusinessIndustryTruckingShipping,
	BusinessIndustryWholesale:                 bank.BusinessIndustryWholesale,
	BusinessIndustryWarehouseDistribution:     bank.BusinessIndustryWarehouseDistribution,
	BusinessIndustryOtherAccomodation:         bank.BusinessIndustryOtherAccomodation,
	BusinessIndustryRestaurantsWithCash:       bank.BusinessIndustryRestaurantsWithCash,
	BusinessIndustryAnimalFarmingProduction:   bank.BusinessIndustryAnimalFarmingProduction,
	BusinessIndustryCropFarming:               bank.BusinessIndustryCropFarming,
	BusinessIndustryForestry:                  bank.BusinessIndustryForestry,
	BusinessIndustryFishingHunting:            bank.BusinessIndustryFishingHunting,
	BusinessIndustryOtherFarmingHunting:       bank.BusinessIndustryOtherFarmingHunting,
	BusinessIndustryMuseums:                   bank.BusinessIndustryMuseums,
	BusinessIndustryHospitals:                 bank.BusinessIndustryHospitals,
	BusinessIndustryCollegeUniversitySchools:  bank.BusinessIndustryCollegeUniversitySchools,
	BusinessIndustryBank:                      bank.BusinessIndustryBank,
	BusinessIndustryFinancialInvestments:      bank.BusinessIndustryFinancialInvestments,
	BusinessIndustryFundsTrustsOther:          bank.BusinessIndustryFundsTrustsOther,
	BusinessIndustryInsurance:                 bank.BusinessIndustryInsurance,
	BusinessIndustryMoneyTransferRemittance:   bank.BusinessIndustryMoneyTransferRemittance,
	BusinessIndustryPrivateInvestment:         bank.BusinessIndustryPrivateInvestment,
	BusinessIndustryOtherManufacturing:        bank.BusinessIndustryOtherManufacturing,
	BusinessIndustryIndustrialMachinery:       bank.BusinessIndustryIndustrialMachinery,
	BusinessIndustryEmploymentServices:        bank.BusinessIndustryEmploymentServices,
	BusinessIndustryGovernmentAgency:          bank.BusinessIndustryGovernmentAgency,
	BusinessIndustryNonGovernment:             bank.BusinessIndustryNonGovernment,
	BusinessIndustryReligiousOrganization:     bank.BusinessIndustryReligiousOrganization,
	BusinessIndustryUnions:                    bank.BusinessIndustryUnions,
	BusinessIndustryRetailJeweler:             bank.BusinessIndustryRetailJeweler,
	BusinessIndustryRetailWithCash:            bank.BusinessIndustryRetailWithCash,
	BusinessIndustryUsedClothesDealer:         bank.BusinessIndustryUsedClothesDealer,
	BusinessIndustryTourOperator:              bank.BusinessIndustryTourOperator,
	BusinessIndustryWholesaleJeweler:          bank.BusinessIndustryWholesaleJeweler,

	// Not supported
	BusinessIndustryCasinoHotel:            bank.BusinessIndustryCasinoHotel,
	BusinessIndustryCasinoGaming:           bank.BusinessIndustryCasinoGaming,
	BusinessIndustryRaceTrack:              bank.BusinessIndustryRaceTrack,
	BusinessIndustryCheckCasher:            bank.BusinessIndustryCheckCasher,
	BusinessIndustryCollectionsAgency:      bank.BusinessIndustryCollectionsAgency,
	BusinessIndustryCurrencyExchange:       bank.BusinessIndustryCurrencyExchange,
	BusinessIndustryCigaretteManufacturing: bank.BusinessIndustryCigaretteManufacturing,
	BusinessIndustryPrivateATM:             bank.BusinessIndustryPrivateATM,
	BusinessIndustryConsulateEmbassy:       bank.BusinessIndustryConsulateEmbassy,
	BusinessIndustryBeerWineLiquorStores:   bank.BusinessIndustryBeerWineLiquorStores,
	BusinessIndustryPawnShop:               bank.BusinessIndustryPawnShop,
}

var partnerIndustryFromMap = map[bank.BusinessIndustry]BusinessIndustry{
	bank.BusinessIndustryHotelMotel:                BusinessIndustryHotelMotel,
	bank.BusinessIndustryOtherFoodServices:         BusinessIndustryOtherFoodServices,
	bank.BusinessIndustryRestaurants:               BusinessIndustryRestaurants,
	bank.BusinessIndustryArtPhotography:            BusinessIndustryArtPhotography,
	bank.BusinessIndustryOtherArtsEntertainment:    BusinessIndustryOtherArtsEntertainment,
	bank.BusinessIndustryFitnessSportsCenters:      BusinessIndustryFitnessSportsCenters,
	bank.BusinessIndustrySportsTeamsClubs:          BusinessIndustrySportsTeamsClubs,
	bank.BusinessIndustryConstruction:              BusinessIndustryConstruction,
	bank.BusinessIndustryBuildingMaterialsHardware: BusinessIndustryBuildingMaterialsHardware,
	bank.BusinessIndustryOtherTradeContractor:      BusinessIndustryOtherTradeContractor,
	bank.BusinessIndustryPlumbingHVAC:              BusinessIndustryPlumbingHVAC,
	bank.BusinessIndustryHealthServices:            BusinessIndustryHealthServices,
	bank.BusinessIndustryOtherEducationServices:    BusinessIndustryOtherEducationServices,
	bank.BusinessIndustryOtherHealthFitness:        BusinessIndustryOtherHealthFitness,
	bank.BusinessIndustryAccountingTaxPrep:         BusinessIndustryAccountingTaxPrep,
	bank.BusinessIndustryRealEstate:                BusinessIndustryRealEstate,
	bank.BusinessIndustryBeautyOrBarberShops:       BusinessIndustryBeautyOrBarberShops,
	bank.BusinessIndustryHomeFurnishing:            BusinessIndustryHomeFurnishing,
	bank.BusinessIndustryCarWash:                   BusinessIndustryCarWash,
	bank.BusinessIndustryComputerServiceRepair:     BusinessIndustryComputerServiceRepair,
	bank.BusinessIndustryFreelanceProfessional:     BusinessIndustryFreelanceProfessional,
	bank.BusinessIndustryLandscapeServices:         BusinessIndustryLandscapeServices,
	bank.BusinessIndustryLegalServices:             BusinessIndustryLegalServices,
	bank.BusinessIndustryMassageTanningServices:    BusinessIndustryMassageTanningServices,
	bank.BusinessIndustryOtherProfessionalServices: BusinessIndustryOtherProfessionalServices,
	bank.BusinessIndustryAutoDealers:               BusinessIndustryAutoDealers,
	bank.BusinessIndustryOnlineRetailer:            BusinessIndustryOnlineRetailer,
	bank.BusinessIndustryRetail:                    BusinessIndustryRetail,
	bank.BusinessIndustryGasolineServiceStation:    BusinessIndustryGasolineServiceStation,
	bank.BusinessIndustryOtherTransportServices:    BusinessIndustryOtherTransportServices,
	bank.BusinessIndustryOtherTravelServices:       BusinessIndustryOtherTravelServices,
	bank.BusinessIndustryParkingGarages:            BusinessIndustryParkingGarages,
	bank.BusinessIndustryTaxi:                      BusinessIndustryTaxi,
	bank.BusinessIndustryTravelAgency:              BusinessIndustryTravelAgency,
	bank.BusinessIndustryTruckingShipping:          BusinessIndustryTruckingShipping,
	bank.BusinessIndustryWholesale:                 BusinessIndustryWholesale,
	bank.BusinessIndustryWarehouseDistribution:     BusinessIndustryWarehouseDistribution,
	bank.BusinessIndustryOtherAccomodation:         BusinessIndustryOtherAccomodation,
	bank.BusinessIndustryRestaurantsWithCash:       BusinessIndustryRestaurantsWithCash,
	bank.BusinessIndustryAnimalFarmingProduction:   BusinessIndustryAnimalFarmingProduction,
	bank.BusinessIndustryCropFarming:               BusinessIndustryCropFarming,
	bank.BusinessIndustryForestry:                  BusinessIndustryForestry,
	bank.BusinessIndustryFishingHunting:            BusinessIndustryFishingHunting,
	bank.BusinessIndustryOtherFarmingHunting:       BusinessIndustryOtherFarmingHunting,
	bank.BusinessIndustryMuseums:                   BusinessIndustryMuseums,
	bank.BusinessIndustryHospitals:                 BusinessIndustryHospitals,
	bank.BusinessIndustryCollegeUniversitySchools:  BusinessIndustryCollegeUniversitySchools,
	bank.BusinessIndustryBank:                      BusinessIndustryBank,
	bank.BusinessIndustryFinancialInvestments:      BusinessIndustryFinancialInvestments,
	bank.BusinessIndustryFundsTrustsOther:          BusinessIndustryFundsTrustsOther,
	bank.BusinessIndustryInsurance:                 BusinessIndustryInsurance,
	bank.BusinessIndustryMoneyTransferRemittance:   BusinessIndustryMoneyTransferRemittance,
	bank.BusinessIndustryPrivateInvestment:         BusinessIndustryPrivateInvestment,
	bank.BusinessIndustryOtherManufacturing:        BusinessIndustryOtherManufacturing,
	bank.BusinessIndustryIndustrialMachinery:       BusinessIndustryIndustrialMachinery,
	bank.BusinessIndustryEmploymentServices:        BusinessIndustryEmploymentServices,
	bank.BusinessIndustryGovernmentAgency:          BusinessIndustryGovernmentAgency,
	bank.BusinessIndustryNonGovernment:             BusinessIndustryNonGovernment,
	bank.BusinessIndustryReligiousOrganization:     BusinessIndustryReligiousOrganization,
	bank.BusinessIndustryUnions:                    BusinessIndustryUnions,
	bank.BusinessIndustryRetailJeweler:             BusinessIndustryRetailJeweler,
	bank.BusinessIndustryRetailWithCash:            BusinessIndustryRetailWithCash,
	bank.BusinessIndustryUsedClothesDealer:         BusinessIndustryUsedClothesDealer,
	bank.BusinessIndustryTourOperator:              BusinessIndustryTourOperator,
	bank.BusinessIndustryWholesaleJeweler:          BusinessIndustryWholesaleJeweler,

	// Not supported
	bank.BusinessIndustryCasinoHotel:            BusinessIndustryCasinoHotel,
	bank.BusinessIndustryCasinoGaming:           BusinessIndustryCasinoGaming,
	bank.BusinessIndustryRaceTrack:              BusinessIndustryRaceTrack,
	bank.BusinessIndustryCheckCasher:            BusinessIndustryCheckCasher,
	bank.BusinessIndustryCollectionsAgency:      BusinessIndustryCollectionsAgency,
	bank.BusinessIndustryCurrencyExchange:       BusinessIndustryCurrencyExchange,
	bank.BusinessIndustryCigaretteManufacturing: BusinessIndustryCigaretteManufacturing,
	bank.BusinessIndustryPrivateATM:             BusinessIndustryPrivateATM,
	bank.BusinessIndustryConsulateEmbassy:       BusinessIndustryConsulateEmbassy,
	bank.BusinessIndustryBeerWineLiquorStores:   BusinessIndustryBeerWineLiquorStores,
	bank.BusinessIndustryPawnShop:               BusinessIndustryPawnShop,
}

type BusinessIdentityDocument string

const (
	BusinessIdentityDocumentArticlesOfIncorporation = BusinessIdentityDocument("articles_of_incorporation")
	BusinessIdentityDocumentArticlesOfOrganization  = BusinessIdentityDocument("articles_of_organization")
	BusinessIdentityDocumentAssumedNameCertificate  = BusinessIdentityDocument("assumed_name_certificate")
	BusinessIdentityDocumentBusinessLicense         = BusinessIdentityDocument("business_license_certificate_of_good_standing")
	BusinessIdentityDocumentPartnershipCertificate  = BusinessIdentityDocument("certificate_of_limited_partnership")
	BusinessIdentityDocumentPartnershipAgreement    = BusinessIdentityDocument("partnership_agreement")
	BusinessIdentityDocumentCertificateOfFormation  = BusinessIdentityDocument("certificate_of_formation")

	// Supported for DL and extra documents
	BusinessIdentityDocumentOther = BusinessIdentityDocument("other_business_documentation")

	// Not supported
	BusinessIdentityDocumentLettersOfTestamentary = BusinessIdentityDocument("letters_testamentary")
	BusinessIdentityDocumentTrustAgreement        = BusinessIdentityDocument("trust_agreement")
)

var partnerBusinessDocumentTo = map[BusinessIdentityDocument]bank.BusinessIdentityDocument{
	BusinessIdentityDocumentArticlesOfIncorporation: bank.BusinessIdentityDocumentArticlesOfIncorporation,
	BusinessIdentityDocumentArticlesOfOrganization:  bank.BusinessIdentityDocumentArticlesOfOrganization,
	BusinessIdentityDocumentAssumedNameCertificate:  bank.BusinessIdentityDocumentAssumedNameCertificate,
	BusinessIdentityDocumentBusinessLicense:         bank.BusinessIdentityDocumentBusinessLicense,
	BusinessIdentityDocumentPartnershipCertificate:  bank.BusinessIdentityDocumentPartnershipCertificate,
	BusinessIdentityDocumentPartnershipAgreement:    bank.BusinessIdentityDocumentPartnershipAgreement,
	BusinessIdentityDocumentCertificateOfFormation:  bank.BusinessIdentityDocumentCertificateOfFormation,
	BusinessIdentityDocumentOther:                   bank.BusinessIdentityDocumentOther,
}

var partnerBusinessDocumentFrom = map[bank.BusinessIdentityDocument]BusinessIdentityDocument{
	bank.BusinessIdentityDocumentArticlesOfIncorporation: BusinessIdentityDocumentArticlesOfIncorporation,
	bank.BusinessIdentityDocumentArticlesOfOrganization:  BusinessIdentityDocumentArticlesOfOrganization,
	bank.BusinessIdentityDocumentAssumedNameCertificate:  BusinessIdentityDocumentAssumedNameCertificate,
	bank.BusinessIdentityDocumentBusinessLicense:         BusinessIdentityDocumentBusinessLicense,
	bank.BusinessIdentityDocumentPartnershipCertificate:  BusinessIdentityDocumentPartnershipCertificate,
	bank.BusinessIdentityDocumentPartnershipAgreement:    BusinessIdentityDocumentPartnershipAgreement,
	bank.BusinessIdentityDocumentCertificateOfFormation:  BusinessIdentityDocumentCertificateOfFormation,
	bank.BusinessIdentityDocumentDriversLicense:          BusinessIdentityDocumentOther,
	bank.BusinessIdentityDocumentOther:                   BusinessIdentityDocumentOther,
}

type CreateBusinessRequest struct {
	LegalName          string                  `json:"legal_name"`
	DBA                string                  `json:"dba,omitempty"`
	TINType            BusinessTINType         `json:"tin_type"`
	TIN                string                  `json:"tin_number"`
	EntityType         BusinessEntity          `json:"entity_type"`
	IndustryType       BusinessIndustry        `json:"industry_type"`
	Contacts           []ContactRequest        `json:"contacts"`
	Purpose            string                  `json:"purpose"`
	ExpectedActivities []ExpectedActivity      `json:"expected_activities"`
	Members            []BusinessMemberRequest `json:"members"`
	Address            []AddressRequest        `json:"addresses"`
	EntityFormation    EntityFormationRequest  `json:"entity_formation"`
	Questions          QuestionsRequest        `json:"questions"`
}

type UpdateBusinessRequest struct {
	LegalName       string                  `json:"legal_name,omitempty"`
	DBA             *string                 `json:"dba,omitempty"`
	TINType         BusinessTINType         `json:"tin_type,omitempty"`
	TIN             string                  `json:"tin_number,omitempty"`
	EntityType      BusinessEntity          `json:"entity_type,omitempty"`
	IndustryType    BusinessIndustry        `json:"industry_type,omitempty"`
	Purpose         string                  `json:"purpose,omitempty"`
	EntityFormation *EntityFormationRequest `json:"entity_formation,omitempty"`
}

type BusinessMemberType string

const (
	BusinessMemberTypeOwner      = BusinessMemberType("owner")
	BusinessMemberTypeController = BusinessMemberType("control_manager")
)

type BusinessMemberTitle string

const (
	BusinessMemberTitleCEO       = BusinessMemberTitle("chief_executive_officer")
	BusinessMemberTitleCFO       = BusinessMemberTitle("chief_financial_officer")
	BusinessMemberTitleCOO       = BusinessMemberTitle("chief_operating_officer")
	BusinessMemberTitlePresident = BusinessMemberTitle("president")
	BusinessMemberTitleVP        = BusinessMemberTitle("vice_president")
	BusinessMemberTitleSVP       = BusinessMemberTitle("senior_vice_president")
	BusinessMemberTitleTreasurer = BusinessMemberTitle("treasurer")
	BusinessMemberTitleGP        = BusinessMemberTitle("general_partner")
	BusinessMemberTitleOther     = BusinessMemberTitle("other")
)

var partnerMemberTitleToMap = map[BusinessMemberTitle]bank.BusinessMemberTitle{
	BusinessMemberTitleCEO:       bank.BusinessMemberTitleCEO,
	BusinessMemberTitleCFO:       bank.BusinessMemberTitleCFO,
	BusinessMemberTitleCOO:       bank.BusinessMemberTitleCOO,
	BusinessMemberTitlePresident: bank.BusinessMemberTitlePresident,
	BusinessMemberTitleVP:        bank.BusinessMemberTitleVP,
	BusinessMemberTitleSVP:       bank.BusinessMemberTitleSVP,
	BusinessMemberTitleTreasurer: bank.BusinessMemberTitleTreasurer,
	BusinessMemberTitleGP:        bank.BusinessMemberTitleGP,
	BusinessMemberTitleOther:     bank.BusinessMemberTitleOther,
}

var partnerMemberTitleFromMap = map[bank.BusinessMemberTitle]BusinessMemberTitle{
	bank.BusinessMemberTitleCEO:       BusinessMemberTitleCEO,
	bank.BusinessMemberTitleCFO:       BusinessMemberTitleCFO,
	bank.BusinessMemberTitleCOO:       BusinessMemberTitleCOO,
	bank.BusinessMemberTitlePresident: BusinessMemberTitlePresident,
	bank.BusinessMemberTitleVP:        BusinessMemberTitleVP,
	bank.BusinessMemberTitleSVP:       BusinessMemberTitleSVP,
	bank.BusinessMemberTitleTreasurer: BusinessMemberTitleTreasurer,
	bank.BusinessMemberTitleGP:        BusinessMemberTitleGP,
	bank.BusinessMemberTitleOther:     BusinessMemberTitleOther,
}

// Requirements for business members
// In Create a business: POST, Step 3, the example JSON includes a members array.
// That array contains two JSON structures that specify information for two members.
// The two-member structure means that when creating a business record, the members array:
//		- Must include one and only one “user_type”: “control_manager”
//		- Should include an additional “user_type”: “owner” for all persons who own 25 percent or more of the business entity.
type BusinessMemberRequest struct {
	UserID              string              `json:"user_id"`
	UserType            BusinessMemberType  `json:"user_type"`
	OwnershipPercentage string              `json:"ownership_percentage,omitempty"`
	Title               BusinessMemberTitle `json:"title"`
	TitleDescription    *string             `json:"title_description,omitempty"`
}

func partnerMemberRequestFrom(r bank.APIRequest, m bank.BusinessMemberRequest, requestType BusinessMemberType) (*BusinessMemberRequest, error) {
	// Must be in consumer table
	c, err := data.NewConsumerService(r, bank.ProviderNameBBVA).GetByConsumerID(m.ConsumerID)
	if err != nil {
		return nil, err
	}

	memberType := requestType

	var ownership string
	switch requestType {
	case BusinessMemberTypeOwner:
		ownership = strconv.FormatInt(int64(m.Ownership), 10)
	case BusinessMemberTypeController:
		break
	default:
		errMsg := fmt.Sprintf("invalid request type: %s", string(requestType))
		return nil, errors.New(errMsg)
	}

	titleDesc := m.TitleDesc
	title, ok := partnerMemberTitleFromMap[m.Title]
	if !ok {
		title = BusinessMemberTitleOther
		if titleDesc == nil {
			titleString := m.Title.String()
			titleDesc = &titleString
		}
	}

	// When other check for value description value
	if title == BusinessMemberTitleOther {
		if titleDesc == nil {
			return nil, errors.New("member type other must have a valid title description")
		}

		trimmedTitleDesc := strings.TrimSpace(*titleDesc)
		if trimmedTitleDesc == "" {
			return nil, errors.New("member type other must have a valid title description")
		}

		titleDesc = &trimmedTitleDesc
	}

	return &BusinessMemberRequest{
		UserID:              string(c.BankID),
		UserType:            memberType,
		OwnershipPercentage: ownership,
		Title:               title,
		TitleDescription:    titleDesc,
	}, nil
}

func partnerMembersFrom(r bank.APIRequest, pms []bank.BusinessMemberRequest) ([]BusinessMemberRequest, error) {
	var m []BusinessMemberRequest
	for _, pm := range pms {
		// Owner
		if pm.Ownership > 0 {
			member, err := partnerMemberRequestFrom(r, pm, BusinessMemberTypeOwner)
			if err != nil {
				return m, err
			}

			m = append(m, *member)
		} else if !pm.IsControllingManager {
			return nil, errors.New("owner members must have an ownership percentage")
		}

		// Control Manager
		if pm.IsControllingManager {
			member, err := partnerMemberRequestFrom(r, pm, BusinessMemberTypeController)
			if err != nil {
				return m, err
			}

			m = append(m, *member)
		}
	}

	return m, nil
}

type BusinessOperation string

const (
	BusinessOperationLocal            = BusinessOperation("local_business")
	BusinessOperationForeignWithLocal = BusinessOperation("foreign_with_local_operations")
	BusinessOperationForeign          = BusinessOperation("foreign_without_local_operations")
)

var partnerOperationTypeMapFrom = map[bank.BusinessOperationType]BusinessOperation{
	bank.BusinessOperationTypeLocal:            BusinessOperationLocal,
	bank.BusinessOperationTypeForeignWithLocal: BusinessOperationForeignWithLocal,
	bank.BusinessOperationTypeForeign:          BusinessOperationForeign,
}

var parterOperationTypeMapTo = map[BusinessOperation]bank.BusinessOperationType{
	BusinessOperationLocal:            bank.BusinessOperationTypeLocal,
	BusinessOperationForeignWithLocal: bank.BusinessOperationTypeForeignWithLocal,
	BusinessOperationForeign:          bank.BusinessOperationTypeForeign,
}

type EntityFormationRequest struct {
	Document          BusinessIdentityDocument `json:"document"`
	Number            string                   `json:"number"`
	Issuer            string                   `json:"issuer"`
	IssueDate         string                   `json:"issue_date"`
	ExpirationDate    string                   `json:"expiration_date"`
	BusinessOperation BusinessOperation        `json:"business_operation"`
	OriginCountry     Country3Alpha            `json:"origin_country"`
	OriginState       string                   `json:"origin_state"`
	OriginDate        string                   `json:"origin_date"`
}

type QuestionsRequest struct {
	Affiliate                       string `json:"affiliate"`
	PrivateATM                      string `json:"private_atm"`
	InternetGambling                string `json:"internet_gambling"`
	InternetGamblingAgreement       string `json:"internet_gambling_agreement"`
	InternetGamblingNotification    string `json:"internet_gambling_notification"`
	CheckCasher                     string `json:"check_casher"`
	CheckCasherGreaterAllowedAmount string `json:"check_casher_greater_allowed_amount"`
	InternationalAffairs            string `json:"international_affairs"`
	CardGames                       string `json:"card_games"`
	PEP                             string `json:"pep"`
	Attestation                     string `json:"attestation"`
}

type CreateBusinessResponse struct {
	BusinessID string                        `json:"business_id"`
	Contacts   []ContactEntityResponse       `json:"contacts"`
	Addresses  []AddressEntityResponse       `json:"addresses"`
	Owners     []CreateBusinessOwnerResponse `json:"members"`
	KYC        KYCResponse                   `json:"kyc"`
}

type IdentityStatusBusinessResponse struct {
	BusinessID string      `json:"business_id"`
	KYC        KYCResponse `json:"kyc"`
}

type CreateBusinessOwnerResponse struct {
	ID       string             `json:"id"`
	UserID   string             `json:"user_id"`
	UserType BusinessMemberType `json:"user_type"`
}

// Create
// https://bbvaopenplatform.com/docs/guides%7Capicontent%7C02-business
// https://bbvaopenplatform.com/docs/guides%7Capicontent%7Cbusiness-guide
func (s *BusinessService) Create(preq bank.CreateBusinessRequest) (*bank.IdentityStatusBusinessResponse, error) {
	b, err := data.NewBusinessService(s.request, bank.ProviderNameBBVA).GetByBusinessID(preq.BusinessID)
	if b != nil {
		return nil, errors.New("business already exists")
	}

	atLegal := AddressTypeLegal
	atHeadquarter := AddressTypeHeadquarter
	var addresses = []AddressRequest{
		addressFromPartner(preq.LegalAddress, &atLegal),
		addressFromPartner(preq.HeadquarterAddress, &atHeadquarter),
	}

	if preq.MailingAddress != nil {
		at := AddressTypeMailing
		addresses = append(addresses, addressFromPartner(*preq.MailingAddress, &at))
	}

	var dba string
	if len(preq.DBA) > 0 && len(preq.DBA[0]) > 0 {
		dba = preq.DBA[0]
	}

	var ef EntityFormationRequest
	if preq.EntityFormation != nil {
		docType, ok := partnerBusinessDocumentFrom[preq.EntityFormation.DocumentType]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidDocumentType)
		}

		op, ok := partnerOperationTypeMapFrom[preq.OperationType]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidBusinessOp)
		}

		country, ok := partnerCountryFrom[preq.OriginCountry]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidCountry)
		}

		ef = EntityFormationRequest{
			Document:          docType,
			Number:            preq.EntityFormation.Number,
			Issuer:            preq.OriginState,
			IssueDate:         preq.EntityFormation.IssueDate.Format("2006-01-02"),
			ExpirationDate:    preq.EntityFormation.ExpirationDate.Format("2006-01-02"),
			BusinessOperation: op,
			OriginCountry:     country,
			OriginState:       preq.OriginState,
			OriginDate:        preq.OriginDate.Format("2006-01-02"),
		}
	} else {
		op, ok := partnerOperationTypeMapFrom[preq.OperationType]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidBusinessOp)
		}

		country, ok := partnerCountryFrom[preq.OriginCountry]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidCountry)
		}

		ef = EntityFormationRequest{
			BusinessOperation: op,
			OriginCountry:     country,
			OriginState:       preq.OriginState,
			OriginDate:        preq.OriginDate.Format("2006-01-02"),
		}
	}

	members, err := partnerMembersFrom(s.request, preq.Members)
	if err != nil {
		return nil, err
	}

	entityType, ok := partnerEntityFromMap[preq.EntityType]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityType)
	}

	activities, err := activityFromPartner(preq.ExpectedActivities)
	if err != nil {
		return nil, err
	}

	taxIDType, ok := partnerTaxIDFromMap[preq.TaxIDType]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidTaxIDType)
	}

	industryType, ok := partnerIndustryFromMap[preq.IndustryType]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidIndustryType)
	}

	cr := &CreateBusinessRequest{
		LegalName:       stripBusinessName(preq.LegalName),
		DBA:             stripBusinessName(dba),
		TINType:         taxIDType,
		TIN:             preq.TaxID,
		EntityType:      entityType,
		IndustryType:    industryType,
		Purpose:         preq.Purpose,
		EntityFormation: ef,
		Questions: QuestionsRequest{
			Affiliate:                       "no",
			PrivateATM:                      "no",
			InternetGambling:                "no",
			InternetGamblingAgreement:       "yes",
			InternetGamblingNotification:    "yes",
			CheckCasher:                     "no",
			CheckCasherGreaterAllowedAmount: "no",
			InternationalAffairs:            "no",
			CardGames:                       "no",
			PEP:                             "no",
			Attestation:                     "yes",
		},
		ExpectedActivities: activities,
		Contacts: []ContactRequest{
			ContactRequest{
				Type:  ContactTypePhone,
				Value: preq.Phone,
			},
			ContactRequest{
				Type:  ContactTypeEmail,
				Value: strings.ToLower(preq.Email),
			},
		},
		Members: members,
		Address: addresses,
	}

	r, err := s.client.post("business/v3.1", s.request, cr)
	if err != nil {
		return nil, err
	}

	var resp = CreateBusinessResponse{}
	if err := s.client.do(r, &resp); err != nil {
		return nil, err
	}

	kycStatus, ok := kycStatusToPartnerMap[KYCStatus(strings.ToUpper(string(resp.KYC.Status)))]
	if !ok {
		return nil, errors.New("invalid KYC status")
	}

	// Save business entity to partner entity table
	c := data.BusinessCreate{
		BusinessID: bank.BusinessID(preq.BusinessID),
		BankID:     bank.BusinessBankID(resp.BusinessID),
		KYCStatus:  &kycStatus,
	}

	b, err = data.NewBusinessService(s.request, bank.ProviderNameBBVA).Create(c)
	if err != nil {
		return nil, err
	}

	for _, member := range resp.Owners {
		c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByBankID(bank.ConsumerBankID(member.UserID))
		if err != nil {
			// Skip member?
			log.Println(err)
			continue
		}

		var bankOwnerID *bank.MemberBankID
		var bankControlID *bank.MemberBankID

		switch member.UserType {
		case BusinessMemberTypeOwner:
			ownerID := bank.MemberBankID(member.ID)
			bankOwnerID = &ownerID
		case BusinessMemberTypeController:
			controlID := bank.MemberBankID(member.ID)
			bankControlID = &controlID
		default:
			errMsg := fmt.Sprintf("invalid user type: %s", string(member.UserType))
			return nil, errors.New(errMsg)
		}

		mem, err := data.NewBusinessMemberService(s.request, bank.ProviderNameBBVA).
			GetByConsumerID(bank.BusinessID(b.BusinessID), bank.ConsumerID(c.ConsumerID))
		if err == sql.ErrNoRows {
			m := data.BusinessMemberCreate{
				BusinessID:    b.ID,
				ConsumerID:    c.ID,
				BankID:        bankOwnerID,
				BankControlID: bankControlID,
			}

			_, err = data.NewBusinessMemberService(s.request, bank.ProviderNameBBVA).Create(m)
		} else if err == nil {
			m := data.BusinessMemberUpdate{
				ID: mem.ID,
			}

			if bankOwnerID != nil {
				m.BankID = bankOwnerID
			}

			if bankControlID != nil {
				m.BankControlID = bankControlID
			}

			_, err = data.NewBusinessMemberService(s.request, bank.ProviderNameBBVA).Update(m)
		}
	}

	// Map contacts
	err = createBusinessContacts(s.client, s.request, b, resp.Contacts)
	if err != nil {
		return nil, err
	}

	// Update contacts
	updateBusinessContact(s.client, s.request, b, bank.BusinessPropertyTypeContactEmail, preq.Email)
	updateBusinessContact(s.client, s.request, b, bank.BusinessPropertyTypeContactPhone, preq.Phone)

	// Map addresses
	err = createBusinessAddresses(s.client, s.request, b, resp.Addresses)
	if err != nil {
		return nil, err
	}

	// Update addresses
	updateBusinessAddress(s.client, s.request, b, bank.BusinessPropertyTypeAddressLegal, preq.LegalAddress)
	updateBusinessAddress(s.client, s.request, b, bank.BusinessPropertyTypeAddressHeadquarter, preq.HeadquarterAddress)

	if preq.MailingAddress != nil {
		updateBusinessAddress(s.client, s.request, b, bank.BusinessPropertyTypeAddressMailing, *preq.MailingAddress)
	}

	kyc, err := resp.KYC.toPartnerBankKYCResponse(nil)
	if err != nil {
		return nil, err
	}

	// TODO: Add members to member entity table (id, business_id, member_id, bank_member_id)
	return &bank.IdentityStatusBusinessResponse{
		BusinessID: b.BusinessID,
		BankID:     b.BankID,
		KYC:        *kyc,
	}, nil
}

func (s *BusinessService) Update(preq bank.UpdateBusinessRequest) (*bank.IdentityStatusBusinessResponse, error) {
	b, err := data.NewBusinessService(s.request, bank.ProviderNameBBVA).GetByBusinessID(preq.BusinessID)
	if err != nil {
		return nil, err
	}

	var ef *EntityFormationRequest
	if preq.EntityFormation != nil {
		docType, ok := partnerBusinessDocumentFrom[preq.EntityFormation.DocumentType]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidDocumentType)
		}

		op, ok := partnerOperationTypeMapFrom[preq.OperationType]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidBusinessOp)
		}

		country, ok := partnerCountryFrom[preq.OriginCountry]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidCountry)
		}

		ef = &EntityFormationRequest{
			Document:          docType,
			Number:            preq.EntityFormation.Number,
			Issuer:            preq.OriginState,
			IssueDate:         preq.EntityFormation.IssueDate.Format("2006-01-02"),
			ExpirationDate:    preq.EntityFormation.ExpirationDate.Format("2006-01-02"),
			BusinessOperation: op,
			OriginCountry:     country,
			OriginState:       preq.OriginState,
			OriginDate:        preq.OriginDate.Format("2006-01-02"),
		}
	}

	entityType, ok := partnerEntityFromMap[preq.EntityType]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityType)
	}

	industryType, ok := partnerIndustryFromMap[preq.IndustryType]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidIndustryType)
	}

	// DBA is optional but should be cleared if none exists
	var dba *string
	if preq.DBA != nil && len(preq.DBA) > 0 && len(preq.DBA[0]) > 0 {
		dbaParsed := stripBusinessName(preq.DBA[0])
		dba = &dbaParsed
	}

	ur := &UpdateBusinessRequest{
		LegalName:       stripBusinessName(preq.LegalName),
		DBA:             dba,
		TINType:         partnerTaxIDFromMap[preq.TaxIDType],
		TIN:             preq.TaxID,
		EntityType:      entityType,
		IndustryType:    industryType,
		Purpose:         preq.Purpose,
		EntityFormation: ef,
	}

	r, err := s.client.patch("business/v3.1", s.request, ur)
	if err != nil {
		return nil, err
	}

	r.Header.Set("OP-User-Id", string(b.BankID))
	if err := s.client.do(r, nil); err != nil {
		return nil, err
	}

	return s.Status(b.BusinessID)
}

func (s *BusinessService) Status(id bank.BusinessID) (*bank.IdentityStatusBusinessResponse, error) {
	b, err := data.NewBusinessService(s.request, bank.ProviderNameBBVA).GetByBusinessID(id)
	if err != nil {
		return nil, err
	}

	r, err := s.client.get("business/v3.1/identity", s.request)
	if err != nil {
		return nil, err
	}

	r.Header.Set("OP-User-Id", string(b.BankID))
	var resp = IdentityStatusBusinessResponse{}
	if err := s.client.do(r, &resp); err != nil {
		return nil, err
	}

	kycStatus, ok := kycStatusToPartnerMap[KYCStatus(strings.ToUpper(string(resp.KYC.Status)))]
	if !ok {
		return nil, errors.New("invalid KYC status")
	}

	// Save consumer entity to partner entity table
	u := data.BusinessUpdate{
		ID:        b.ID,
		KYCStatus: &kycStatus,
	}

	_, err = data.NewBusinessService(s.request, bank.ProviderNameBBVA).Update(u)
	if err != nil {
		return nil, err
	}

	kyc, err := resp.KYC.toPartnerBankKYCResponse(nil)
	if err != nil {
		return nil, err
	}

	// TODO: Add members to member entity table (id, business_id, member_id, bank_member_id)
	return &bank.IdentityStatusBusinessResponse{
		BusinessID: b.BusinessID,
		BankID:     b.BankID,
		KYC:        *kyc,
	}, nil
}

func (s *BusinessService) UploadIdentityDocument(preq bank.BusinessIdentityDocumentRequest) (*bank.IdentityDocumentResponse, error) {
	b, err := data.NewBusinessService(s.request, bank.ProviderNameBBVA).GetByBusinessID(preq.BusinessID)
	if err != nil {
		return nil, err
	}

	fileType, ok := partnerContentTypeFrom[preq.ContentType]
	if !ok {
		return nil, errors.New("invalid business identity document content type")
	}

	idvs, err := documentIDVerifyPartnerFrom(preq.IDVerifyRequired)
	if err != nil {
		return nil, err
	}

	request := &BusinessIdentityDocumentRequest{
		File:             base64.StdEncoding.EncodeToString(preq.Content),
		FileType:         fileType,
		IDVerifyRequired: idvs,
	}
	r, err := s.client.post("business/v3.1/identity/document", s.request, request)
	if err != nil {
		return nil, err
	}

	r.Header.Set("OP-User-Id", string(b.BankID))
	var response = bank.IdentityDocumentResponse{}
	if err := s.client.do(r, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

/* Update phone or email */
func (s *BusinessService) UpdateContact(id bank.BusinessID, p bank.BusinessPropertyType, val string) error {
	b, err := data.NewBusinessService(s.request, bank.ProviderNameBBVA).GetByBusinessID(id)
	if err != nil {
		return err
	}

	return updateBusinessContact(s.client, s.request, b, p, val)
}

/* Update address */
func (s *BusinessService) UpdateAddress(id bank.BusinessID, p bank.BusinessPropertyType, a bank.AddressRequest) error {
	b, err := data.NewBusinessService(s.request, bank.ProviderNameBBVA).GetByBusinessID(id)
	if err != nil {
		return err
	}

	return updateBusinessAddress(s.client, s.request, b, p, a)
}
