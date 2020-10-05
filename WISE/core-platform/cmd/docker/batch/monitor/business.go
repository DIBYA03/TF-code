package main

import (
	"context"
	"log"
	"sync"
	"time"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcMonitor "github.com/wiseco/protobuf/golang/transaction/monitor"
)

var KYCStatusToBusinessProto = map[services.KYCStatus]grpcRoot.BusinessKYCStatus{
	services.KYCStatusNotStarted: grpcRoot.BusinessKYCStatus_BKS_NOT_STARTED,
	services.KYCStatusSubmitted:  grpcRoot.BusinessKYCStatus_BKS_SUBMITTED,
	services.KYCStatusReview:     grpcRoot.BusinessKYCStatus_BKS_REVIEW,
	services.KYCStatusApproved:   grpcRoot.BusinessKYCStatus_BKS_APPROVED,
	services.KYCStatusDeclined:   grpcRoot.BusinessKYCStatus_BKS_DECLINED,
}

var BusinessEntityToProto = map[string]grpcRoot.BusinessEntityType{
	business.EntityTypeSoleProprietor:              grpcRoot.BusinessEntityType_BET_SOLE_PROPRIETOR,
	business.EntityTypeyProfessionalAssociation:    grpcRoot.BusinessEntityType_BET_PROFESSIONAL_ASSOCIATION,
	business.EntityTypeSingleMemberLLC:             grpcRoot.BusinessEntityType_BET_LIMITED_LIABILITY_COMPANY,
	business.EntityTypeLimitedLiabilityCompany:     grpcRoot.BusinessEntityType_BET_SOLE_LIMITED_LIABILITY_COMPANY,
	business.EntityTypeGeneralPartnership:          grpcRoot.BusinessEntityType_BET_GENERAL_PARTNERSHIP,
	business.EntityTypeLimitedPartnership:          grpcRoot.BusinessEntityType_BET_LIMITED_PARTNERSHIP,
	business.EntityTypeLimitedLiabilityPartnership: grpcRoot.BusinessEntityType_BET_LIMITED_LIABILITY_PARTNERSHIP,
	business.EntityTypeProfessionalCorporation:     grpcRoot.BusinessEntityType_BET_PROFESSIONAL_CORPORATION,
	business.EntityTypeUnlistedCorporation:         grpcRoot.BusinessEntityType_BET_UNLISTED_CORPORATION,
}

var BusinessIndustryToProto = map[string]grpcRoot.BusinessIndustryType{
	business.IndustryTypeAccountingTaxPrep:         grpcRoot.BusinessIndustryType_BIT_ACCOUNTING_TAX_PREP,
	business.IndustryTypeAnimalFarmingProduction:   grpcRoot.BusinessIndustryType_BIT_ANIMAL_FARMING_PRODUCTION,
	business.IndustryTypeArtPhotography:            grpcRoot.BusinessIndustryType_BIT_ART_PHOTOGRAPHY,
	business.IndustryTypeAutoDealers:               grpcRoot.BusinessIndustryType_BIT_AUTO_DEALERS,
	business.IndustryTypeBank:                      grpcRoot.BusinessIndustryType_BIT_BANK_FINANCIAL_INSTITUTION,
	business.IndustryTypeBeautyOrBarberShops:       grpcRoot.BusinessIndustryType_BIT_BEAUTY_OR_BARBER_SHOPS,
	business.IndustryTypeBeerWineLiquorStores:      grpcRoot.BusinessIndustryType_BIT_BEER_WINE_LIQUOR_STORES,
	business.IndustryTypeBuildingMaterialsHardware: grpcRoot.BusinessIndustryType_BIT_BUILDING_MATERIALS_HARDWARE,
	business.IndustryTypeCarWash:                   grpcRoot.BusinessIndustryType_BIT_CAR_WASH,
	business.IndustryTypeCasinoGaming:              grpcRoot.BusinessIndustryType_BIT_CASINO_GAMBLING_GAMING,
	business.IndustryTypeCasinoHotel:               grpcRoot.BusinessIndustryType_BIT_CASINO_HOTEL,
	business.IndustryTypeCheckCasher:               grpcRoot.BusinessIndustryType_BIT_CHECK_CASHER,
	business.IndustryTypeCigaretteManufacturing:    grpcRoot.BusinessIndustryType_BIT_CIGARETTE_MANUFACTURING,
	business.IndustryTypeCollectionsAgency:         grpcRoot.BusinessIndustryType_BIT_COLLECTION_AGENCIES,
	business.IndustryTypeCollegeUniversitySchools:  grpcRoot.BusinessIndustryType_BIT_COLLEGES_UNIVERSITIES_SCHOOLS,
	business.IndustryTypeComputerServiceRepair:     grpcRoot.BusinessIndustryType_BIT_COMPUTER_SERVICE_REPAIR,
	business.IndustryTypeConstruction:              grpcRoot.BusinessIndustryType_BIT_CONSTRUCTION,
	business.IndustryTypeConsulateEmbassy:          grpcRoot.BusinessIndustryType_BIT_CONSULATE_EMBASSY,
	business.IndustryTypeCropFarming:               grpcRoot.BusinessIndustryType_BIT_CROP_FARMING,
	business.IndustryTypeCurrencyExchange:          grpcRoot.BusinessIndustryType_BIT_CURRENCY_EXCHANGERS,
	business.IndustryTypeEmploymentServices:        grpcRoot.BusinessIndustryType_BIT_EMPLOYMENT_SERVICES,
	business.IndustryTypeFinancialInvestments:      grpcRoot.BusinessIndustryType_BIT_FINANCIAL_INVESTMENTS,
	business.IndustryTypeFishingHunting:            grpcRoot.BusinessIndustryType_BIT_FISHING_HUNTING_TRAPPING,
	business.IndustryTypeFitnessCenter:             grpcRoot.BusinessIndustryType_BIT_FITNESS_SPORTS_CENTERS,
	business.IndustryTypeForestry:                  grpcRoot.BusinessIndustryType_BIT_FORESTRY_ACTIVITIES,
	business.IndustryTypeFreelanceProfessional:     grpcRoot.BusinessIndustryType_BIT_FREELANCE_PROFESSIONAL,
	business.IndustryTypeFundsTrustsOther:          grpcRoot.BusinessIndustryType_BIT_FUNDS_TRUSTS_OTHER,
	business.IndustryTypeGasServiceStation:         grpcRoot.BusinessIndustryType_BIT_GASOLINE_SERVICE_STATION,
	business.IndustryTypeGovernmentAgency:          grpcRoot.BusinessIndustryType_BIT_GOVERNMENT_AGENCY,
	business.IndustryTypeHealthServices:            grpcRoot.BusinessIndustryType_BIT_HEALTH_SERVICES,
	business.IndustryTypeHomeFurnishing:            grpcRoot.BusinessIndustryType_BIT_HOME_FURNISHING,
	business.IndustryTypeHospitals:                 grpcRoot.BusinessIndustryType_BIT_HOSPITALS,
	business.IndustryTypeHotelMotel:                grpcRoot.BusinessIndustryType_BIT_HOTEL_MOTEL,
	business.IndustryTypeIndustrialMachinery:       grpcRoot.BusinessIndustryType_BIT_INDUSTRIAL_COMMERCIAL_MACHINERY,
	business.IndustryTypeInsurance:                 grpcRoot.BusinessIndustryType_BIT_INSURANCE,
	business.IndustryTypeLandscapeServices:         grpcRoot.BusinessIndustryType_BIT_LANDSCAPE_SERVICES,
	business.IndustryTypeLegalServices:             grpcRoot.BusinessIndustryType_BIT_LEGAL_SERVICES,
	business.IndustryTypeMassageTanningServices:    grpcRoot.BusinessIndustryType_BIT_MASSAGE_TANNING_SERVICES,
	business.IndustryTypeMoneyTransferRemittance:   grpcRoot.BusinessIndustryType_BIT_MONEY_TRANSFER_REMITTANCE,
	business.IndustryTypeMuseums:                   grpcRoot.BusinessIndustryType_BIT_MUSEUMS_HISTORICAL_SITES,
	business.IndustryTypeNonGovernment:             grpcRoot.BusinessIndustryType_BIT_NON_GOVERNMENT_ORGANIZATION,
	business.IndustryTypeOnlineRetailer:            grpcRoot.BusinessIndustryType_BIT_ONLINE_RETAILER,
	business.IndustryTypeOtherAccomodation:         grpcRoot.BusinessIndustryType_BIT_OTHER_ACCOMODATIONS,
	business.IndustryTypeOtherArtsEntertainment:    grpcRoot.BusinessIndustryType_BIT_OTHER_ARTS_ENTERTAINMENT_RECREATION,
	business.IndustryTypeOtherEducationServices:    grpcRoot.BusinessIndustryType_BIT_OTHER_EDUCATION_SERVICES,
	business.IndustryTypeOtherFarmingHunting:       grpcRoot.BusinessIndustryType_BIT_OTHER_AGRICULTURE_FORESTRY_FISHING,
	business.IndustryTypeOtherFoodServices:         grpcRoot.BusinessIndustryType_BIT_OTHER_FOOD_SERVICES,
	business.IndustryTypeOtherHealthFitness:        grpcRoot.BusinessIndustryType_BIT_OTHER_HEALTH_FITNESS,
	business.IndustryTypeOtherManufacturing:        grpcRoot.BusinessIndustryType_BIT_OTHER_MANUFACTURING,
	business.IndustryTypeOtherProfessionalServices: grpcRoot.BusinessIndustryType_BIT_OTHER_PROFESSIONAL_SERVICES,
	business.IndustryTypeOtherTradeContractor:      grpcRoot.BusinessIndustryType_BIT_OTHER_TRADE_CONTRACTOR,
	business.IndustryTypeOtherTravelServices:       grpcRoot.BusinessIndustryType_BIT_OTHER_TRAVEL_SERVICES,
	business.IndustryTypeParkingGarages:            grpcRoot.BusinessIndustryType_BIT_PARKING_GARAGES,
	business.IndustryTypePawnShop:                  grpcRoot.BusinessIndustryType_BIT_PAWN_SHOP,
	business.IndustryTypePlumbingHVAC:              grpcRoot.BusinessIndustryType_BIT_PLUMBING_HVAC,
	business.IndustryTypePrivateATM:                grpcRoot.BusinessIndustryType_BIT_PRIVATE_ATM,
	business.IndustryTypePrivateInvestment:         grpcRoot.BusinessIndustryType_BIT_PRIVATE_INVESTMENT_COMPANIES,
	business.IndustryTypeRaceTrack:                 grpcRoot.BusinessIndustryType_BIT_RACE_TRACK,
	business.IndustryTypeRealEstate:                grpcRoot.BusinessIndustryType_BIT_REAL_ESTATE,
	business.IndustryTypeReligiousOrganization:     grpcRoot.BusinessIndustryType_BIT_RELIGIOUS_ORGANIZATION,
	business.IndustryTypeRestaurants:               grpcRoot.BusinessIndustryType_BIT_RESTAURANTS,
	business.IndustryTypeRestaurantsWithCash:       grpcRoot.BusinessIndustryType_BIT_RESTAURANTS_WITH_CASH,
	business.IndustryTypeRetail:                    grpcRoot.BusinessIndustryType_BIT_RETAIL,
	business.IndustryTypeRetailJeweler:             grpcRoot.BusinessIndustryType_BIT_RETAIL_JEWELER_DIAMONDS_GEMS_GOLD,
	business.IndustryTypeRetailWithCash:            grpcRoot.BusinessIndustryType_BIT_RETAIL_WITH_CASH,
	business.IndustryTypeSportsTeamsClubs:          grpcRoot.BusinessIndustryType_BIT_SPORTS_TEAMS_CLUBS,
	business.IndustryTypeTaxi:                      grpcRoot.BusinessIndustryType_BIT_TAXI,
	business.IndustryTypeTourOperator:              grpcRoot.BusinessIndustryType_BIT_TOUR_OPERATOR,
	business.IndustryTypeTransportationServices:    grpcRoot.BusinessIndustryType_BIT_OTHER_TRANSPORT_SERVICES,
	business.IndustryTypeTravelAgency:              grpcRoot.BusinessIndustryType_BIT_TRAVEL_AGENCY,
	business.IndustryTypeTruckingShipping:          grpcRoot.BusinessIndustryType_BIT_TRUCKING_SHIPPING,
	business.IndustryTypeUnions:                    grpcRoot.BusinessIndustryType_BIT_UNIONS,
	business.IndustryTypeUsedClothesDealer:         grpcRoot.BusinessIndustryType_BIT_USED_CLOTHES_DEALERS,
	business.IndustryTypeWarehouseDistribution:     grpcRoot.BusinessIndustryType_BIT_WAREHOUSE_DISTRIBUTION,
	business.IndustryTypeWholesale:                 grpcRoot.BusinessIndustryType_BIT_WHOLESALE,
	business.IndustryTypeWholesaleJeweler:          grpcRoot.BusinessIndustryType_BIT_WHOLESALE_JEWELER,
}

func processBusiness(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, bID id.BusinessID) error {
	sharedBusID := shared.BusinessID(uuid.UUID(bID).String())
	b, err := business.NewBusinessServiceWithout().GetByIdInternal(sharedBusID)
	if err != nil {
		log.Println(err, bID.String())
		return err
	}

	kycStatus, ok := KYCStatusToBusinessProto[b.KYCStatus]
	if !ok {
		log.Printf("Invalid KYC Status: %s %s", b.KYCStatus, bID.String())
		return nil
	}

	if b.KYCStatus != services.KYCStatusApproved {
		log.Printf("KYC Status not approved: %s %s", b.KYCStatus, bID.String())
		return nil
	}

	var la *grpcRoot.Address
	if b.LegalAddress != nil {
		la = &grpcRoot.Address{
			Locality:  b.LegalAddress.City,
			AdminArea: b.LegalAddress.State,
			// TODO: Fix errors in address data
			Country:    "US",
			PostalCode: b.LegalAddress.PostalCode,
		}
	}

	created, err := grpcTypes.TimestampProto(b.Created)
	if err != nil {
		log.Println(err, bID.String())
		return err
	}

	modified, err := grpcTypes.TimestampProto(b.Modified)
	if err != nil {
		log.Println(err, bID.String())
		return err
	}

	status := grpcRoot.BusinessStatus_BS_ACTIVE
	if b.Deactivated != nil {
		status = grpcRoot.BusinessStatus_BS_INACTIVE
	}

	entityType := grpcRoot.BusinessEntityType_BET_UNSPECIFIED
	if b.EntityType != nil {
		entityType, ok = BusinessEntityToProto[*b.EntityType]
		if !ok {
			log.Printf("Invalid Entity Type: %s %s", *b.EntityType, bID.String())
			return nil
		}
	}

	industryType := grpcRoot.BusinessIndustryType_BIT_UNSPECIFIED
	if b.IndustryType != nil {
		industryType, ok = BusinessIndustryToProto[*b.IndustryType]
		if !ok {
			log.Printf("Invalid Entity Type: %s %s", *b.IndustryType, bID.String())
			return nil
		}
	}

	creq := &grpcMonitor.BusinessRequest{
		Id:           bID.String(),
		Name:         b.Name(),
		EntityType:   entityType,
		IndustryType: industryType,
		Email:        shared.StringValue(b.Email),
		Phone:        shared.StringValue(b.Phone),
		KycStatus:    kycStatus,
		LegalAddress: la,
		Status:       status,
		Created:      created,
		Modified:     modified,
	}

	resp, err := monitorClient.AddUpdateBusiness(context.Background() /* client.GetContext() */, creq)
	if err != nil {
		log.Println(err, bID.String())
		return err
	}

	log.Println("Success: ", resp.Id)
	return nil
}

func sendBusinessUpdates(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, dayStart, dayEnd time.Time) {
	// Process in groups of 5
	offset := 0
	limit := 5
	for {
		var businessIDs []id.BusinessID
		err := data.DBWrite.Select(
			&businessIDs,
			`
			SELECT id from business
			WHERE
				(created >= $1 AND created < $2) OR
				(modified >= $1 AND modified < $2)
			ORDER BY created ASC OFFSET $3 LIMIT $4`,
			dayStart,
			dayEnd,
			offset,
			limit,
		)
		if err != nil {
			panic(err)
		} else if len(businessIDs) == 0 {
			log.Println("No more businesses", dayStart, dayEnd)
			break
		}

		wg := sync.WaitGroup{}
		wg.Add(len(businessIDs))
		for _, bID := range businessIDs {
			go func(id id.BusinessID) {
				defer wg.Done()
				_ = processBusiness(monitorClient, id)
			}(bID)
		}

		wg.Wait()
		offset += 5
	}
}
