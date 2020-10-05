package bbva

import (
	"github.com/wiseco/core-platform/partner/bank"
)

type Country3Alpha string

const (
	CountryEmpty = Country3Alpha("")
	CountryAIA   = Country3Alpha("AIA")
	CountryALA   = Country3Alpha("ALA")
	CountryAND   = Country3Alpha("AND")
	CountryARE   = Country3Alpha("ARE")
	CountryARG   = Country3Alpha("ARG")
	CountryARM   = Country3Alpha("ARM")
	CountryAUT   = Country3Alpha("AUT")
	CountryAUS   = Country3Alpha("AUS")
	CountryAZE   = Country3Alpha("AZE")
	CountryBGD   = Country3Alpha("BGD")
	CountryBEL   = Country3Alpha("BEL")
	CountryBEN   = Country3Alpha("BEN")
	CountryBFA   = Country3Alpha("BFA")
	CountryBGR   = Country3Alpha("BGR")
	CountryBMU   = Country3Alpha("BMU")
	CountryBRA   = Country3Alpha("BRA")
	CountryBRN   = Country3Alpha("BRN")
	CountryBTN   = Country3Alpha("BTN")
	CountryCAN   = Country3Alpha("CAN")
	CountryCHE   = Country3Alpha("CHE")
	CountryCHL   = Country3Alpha("CHL")
	CountryCIV   = Country3Alpha("CIV")
	CountryCMR   = Country3Alpha("CMR")
	CountryCOG   = Country3Alpha("COG")
	CountryCOK   = Country3Alpha("COK")
	CountryCOM   = Country3Alpha("COM")
	CountryCPV   = Country3Alpha("CPV")
	CountryCZE   = Country3Alpha("CZE")
	CountryDEU   = Country3Alpha("DEU")
	CountryDJI   = Country3Alpha("DJI")
	CountryDOM   = Country3Alpha("DOM")
	CountryDNK   = Country3Alpha("DNK")
	CountryDZA   = Country3Alpha("DZA")
	CountryECU   = Country3Alpha("ECU")
	CountryERI   = Country3Alpha("ERI")
	CountryEST   = Country3Alpha("EST")
	CountryESP   = Country3Alpha("ESP")
	CountryETH   = Country3Alpha("ETH")
	CountryFIN   = Country3Alpha("FIN")
	CountryFJI   = Country3Alpha("FJI")
	CountryFRA   = Country3Alpha("FRA")
	CountryFSM   = Country3Alpha("FSM")
	CountryGAB   = Country3Alpha("GAB")
	CountryGBR   = Country3Alpha("GBR")
	CountryGEO   = Country3Alpha("GEO")
	CountryGIB   = Country3Alpha("GIB")
	CountryGIN   = Country3Alpha("GIN")
	CountryGLP   = Country3Alpha("GLP")
	CountryGMB   = Country3Alpha("GMB")
	CountryGNQ   = Country3Alpha("GNQ")
	CountryGRC   = Country3Alpha("GRC")
	CountryGRD   = Country3Alpha("GRD")
	CountryGTM   = Country3Alpha("GTM")
	CountryGUY   = Country3Alpha("GUY")
	CountryHND   = Country3Alpha("HND")
	CountryHTI   = Country3Alpha("HTI")
	CountryHUN   = Country3Alpha("HUN")
	CountryIRL   = Country3Alpha("IRL")
	CountryISR   = Country3Alpha("ISR")
	CountryIND   = Country3Alpha("IND")
	CountryITA   = Country3Alpha("ITA")
	CountryJAM   = Country3Alpha("JAM")
	CountryJOR   = Country3Alpha("JOR")
	CountryJPN   = Country3Alpha("JPN")
	CountryKGZ   = Country3Alpha("KGZ")
	CountryKIR   = Country3Alpha("KIR")
	CountryKOR   = Country3Alpha("KOR")
	CountryKWT   = Country3Alpha("KWT")
	CountryLAO   = Country3Alpha("LAO")
	CountryLBR   = Country3Alpha("LBR")
	CountryLIE   = Country3Alpha("LIE")
	CountryLKA   = Country3Alpha("LKA")
	CountryLSO   = Country3Alpha("LSO")
	CountryLTU   = Country3Alpha("LTU")
	CountryLVA   = Country3Alpha("LVA")
	CountryLUX   = Country3Alpha("LUX")
	CountryMAR   = Country3Alpha("MAR")
	CountryMCO   = Country3Alpha("MCO")
	CountryMDG   = Country3Alpha("MDG")
	CountryMDV   = Country3Alpha("MDV")
	CountryMEX   = Country3Alpha("MEX")
	CountryMHL   = Country3Alpha("MHL")
	CountryMLT   = Country3Alpha("MLT")
	CountryMRT   = Country3Alpha("MRT")
	CountryMSR   = Country3Alpha("MSR")
	CountryMTQ   = Country3Alpha("MTQ")
	CountryMUS   = Country3Alpha("MUS")
	CountryMWI   = Country3Alpha("MWI")
	CountryNAM   = Country3Alpha("NAM")
	CountryNER   = Country3Alpha("NER")
	CountryNIU   = Country3Alpha("NIU")
	CountryNLD   = Country3Alpha("NLD")
	CountryNPL   = Country3Alpha("NPL")
	CountryNRU   = Country3Alpha("NRU")
	CountryNZL   = Country3Alpha("NZL")
	CountryNOR   = Country3Alpha("NOR")
	CountryPER   = Country3Alpha("PER")
	CountryPLW   = Country3Alpha("PLW")
	CountryPNG   = Country3Alpha("PNG")
	CountryPOL   = Country3Alpha("POL")
	CountryPRT   = Country3Alpha("PRT")
	CountryPRY   = Country3Alpha("PRY")
	CountryPSE   = Country3Alpha("PSE")
	CountryQAT   = Country3Alpha("QAT")
	CountryROU   = Country3Alpha("ROU")
	CountryRWA   = Country3Alpha("RWA")
	CountrySAU   = Country3Alpha("SAU")
	CountrySEN   = Country3Alpha("SEN")
	CountrySGP   = Country3Alpha("SGP")
	CountrySLB   = Country3Alpha("SLB")
	CountrySLE   = Country3Alpha("SLE")
	CountrySLV   = Country3Alpha("SLV")
	CountrySMR   = Country3Alpha("SMR")
	CountrySTP   = Country3Alpha("STP")
	CountrySUR   = Country3Alpha("SUR")
	CountrySVN   = Country3Alpha("SVN")
	CountrySVK   = Country3Alpha("SVK")
	CountrySWE   = Country3Alpha("SWE")
	CountrySWZ   = Country3Alpha("SWZ")
	CountryTCD   = Country3Alpha("TCD")
	CountryTGO   = Country3Alpha("TGO")
	CountryTJK   = Country3Alpha("TJK")
	CountryTLS   = Country3Alpha("TLS")
	CountryTON   = Country3Alpha("TON")
	CountryTUN   = Country3Alpha("TUN")
	CountryTUV   = Country3Alpha("TUV")
	CountryTWN   = Country3Alpha("TWN")
	CountryTZA   = Country3Alpha("TZA")
	CountryUGA   = Country3Alpha("UGA")
	CountryUSA   = Country3Alpha("USA")
	CountryURY   = Country3Alpha("URY")
	CountryUZB   = Country3Alpha("UZB")
	CountryVAT   = Country3Alpha("VAT")
	CountryVNM   = Country3Alpha("VNM")
	CountryWSM   = Country3Alpha("WSM")
	CountryZAF   = Country3Alpha("ZAF")
	CountryZMB   = Country3Alpha("ZMB")

	// High risk countries
	CountryAFG = Country3Alpha("AFG")
	CountryATG = Country3Alpha("ATG")
	CountryALB = Country3Alpha("ALB")
	CountryBIH = Country3Alpha("BIH")
	CountryBHS = Country3Alpha("BHS")
	CountryBLR = Country3Alpha("BLR")
	CountryBLZ = Country3Alpha("BLZ")
	CountryBWA = Country3Alpha("BWA")
	CountryCAF = Country3Alpha("CAF")
	CountryCHN = Country3Alpha("CHN")
	CountryCOD = Country3Alpha("COD")
	CountryCOL = Country3Alpha("COL")
	CountryCRI = Country3Alpha("CRI")
	CountryCUB = Country3Alpha("CUB")
	CountryCYP = Country3Alpha("CYP")
	CountryDMA = Country3Alpha("DMA")
	CountryEGY = Country3Alpha("EGY")
	CountryGHA = Country3Alpha("GHA")
	CountryGNB = Country3Alpha("GNB")
	CountryHKG = Country3Alpha("HKG")
	CountryIDN = Country3Alpha("IDN")
	CountryIRN = Country3Alpha("IRN")
	CountryIRQ = Country3Alpha("IRQ")
	CountryISL = Country3Alpha("ISL")
	CountryKAZ = Country3Alpha("KAZ")
	CountryKEN = Country3Alpha("KEN")
	CountryKHM = Country3Alpha("KHM")
	CountryKNA = Country3Alpha("KNA")
	CountryLBN = Country3Alpha("LBN")
	CountryLBY = Country3Alpha("LBY")
	CountryLCA = Country3Alpha("LCA")
	CountryMKD = Country3Alpha("MKD")
	CountryMLI = Country3Alpha("MLI")
	CountryMMR = Country3Alpha("MMR")
	CountryMNE = Country3Alpha("MNE")
	CountryMAC = Country3Alpha("MAC")
	CountryMNG = Country3Alpha("MNG")
	CountryMOZ = Country3Alpha("MOZ")
	CountryMYS = Country3Alpha("MYS")
	CountryNGA = Country3Alpha("NGA")
	CountryNIC = Country3Alpha("NIC")
	CountryPAK = Country3Alpha("PAK")
	CountryPAN = Country3Alpha("PAN")
	CountryPHL = Country3Alpha("PHL")
	CountryPRK = Country3Alpha("PRK")
	CountryRUS = Country3Alpha("RUS")
	CountrySDN = Country3Alpha("SDN")
	CountrySOM = Country3Alpha("SOM")
	CountrySRB = Country3Alpha("SRB")
	CountrySSD = Country3Alpha("SSD")
	CountrySYR = Country3Alpha("SYR")
	CountryTHA = Country3Alpha("THA")
	CountryTTO = Country3Alpha("TTO")
	CountryTUR = Country3Alpha("TUR")
	CountryVCT = Country3Alpha("VCT")
	CountryVEN = Country3Alpha("VEN")
	CountryVUT = Country3Alpha("VUT")
	CountryYEM = Country3Alpha("YEM")
	CountryZWE = Country3Alpha("ZWE")
)

var partnerCountryFrom = map[bank.Country]Country3Alpha{
	bank.CountryAD: CountryAND,
	bank.CountryAE: CountryARE,
	bank.CountryAI: CountryAIA,
	bank.CountryAM: CountryARM,
	bank.CountryAR: CountryARG,
	bank.CountryAU: CountryAUS,
	bank.CountryAT: CountryAUT,
	bank.CountryAX: CountryALA,
	bank.CountryAZ: CountryAZE,
	bank.CountryBD: CountryBGD,
	bank.CountryBE: CountryBEL,
	bank.CountryBF: CountryBFA,
	bank.CountryBG: CountryBGR,
	bank.CountryBJ: CountryBEN,
	bank.CountryBM: CountryBMU,
	bank.CountryBN: CountryBRN,
	bank.CountryBR: CountryBRA,
	bank.CountryBT: CountryBTN,
	bank.CountryCA: CountryCAN,
	bank.CountryCG: CountryCOG,
	bank.CountryCH: CountryCHE,
	bank.CountryCI: CountryCIV,
	bank.CountryCK: CountryCOK,
	bank.CountryCL: CountryCHL,
	bank.CountryCM: CountryCMR,
	bank.CountryCV: CountryCPV,
	bank.CountryCZ: CountryCZE,
	bank.CountryDE: CountryDEU,
	bank.CountryDO: CountryDOM,
	bank.CountryDK: CountryDNK,
	bank.CountryDJ: CountryDJI,
	bank.CountryDZ: CountryDZA,
	bank.CountryEC: CountryECU,
	bank.CountryEE: CountryEST,
	bank.CountryER: CountryERI,
	bank.CountryES: CountryESP,
	bank.CountryET: CountryETH,
	bank.CountryFI: CountryFIN,
	bank.CountryFJ: CountryFJI,
	bank.CountryFM: CountryFSM,
	bank.CountryFR: CountryFRA,
	bank.CountryGA: CountryGAB,
	bank.CountryGB: CountryGBR,
	bank.CountryGE: CountryGEO,
	bank.CountryGD: CountryGRD,
	bank.CountryGI: CountryGIB,
	bank.CountryGM: CountryGMB,
	bank.CountryGN: CountryGIN,
	bank.CountryGP: CountryGLP,
	bank.CountryGQ: CountryGNQ,
	bank.CountryGR: CountryGRC,
	bank.CountryGT: CountryGTM,
	bank.CountryGY: CountryGUY,
	bank.CountryHN: CountryHND,
	bank.CountryHT: CountryHTI,
	bank.CountryHU: CountryHUN,
	bank.CountryIE: CountryIRL,
	bank.CountryIL: CountryISR,
	bank.CountryIN: CountryIND,
	bank.CountryIT: CountryITA,
	bank.CountryJM: CountryJAM,
	bank.CountryJO: CountryJOR,
	bank.CountryJP: CountryJPN,
	bank.CountryKG: CountryKGZ,
	bank.CountryKI: CountryKIR,
	bank.CountryKM: CountryCOM,
	bank.CountryKR: CountryKOR,
	bank.CountryKW: CountryKWT,
	bank.CountryLA: CountryLAO,
	bank.CountryLI: CountryLIE,
	bank.CountryLK: CountryLKA,
	bank.CountryLR: CountryLBR,
	bank.CountryLS: CountryLSO,
	bank.CountryLT: CountryLTU,
	bank.CountryLU: CountryLUX,
	bank.CountryLV: CountryLVA,
	bank.CountryMA: CountryMAR,
	bank.CountryMC: CountryMCO,
	bank.CountryMG: CountryMDG,
	bank.CountryMH: CountryMHL,
	bank.CountryMQ: CountryMTQ,
	bank.CountryMR: CountryMRT,
	bank.CountryMS: CountryMSR,
	bank.CountryMT: CountryMLT,
	bank.CountryMU: CountryMUS,
	bank.CountryMV: CountryMDV,
	bank.CountryMW: CountryMWI,
	bank.CountryMX: CountryMEX,
	bank.CountryNA: CountryNAM,
	bank.CountryNE: CountryNER,
	bank.CountryNP: CountryNPL,
	bank.CountryNL: CountryNLD,
	bank.CountryNO: CountryNOR,
	bank.CountryNR: CountryNRU,
	bank.CountryNU: CountryNIU,
	bank.CountryPE: CountryPER,
	bank.CountryPG: CountryPNG,
	bank.CountryPL: CountryPOL,
	bank.CountryPS: CountryPSE,
	bank.CountryPT: CountryPRT,
	bank.CountryPW: CountryPLW,
	bank.CountryPY: CountryPRY,
	bank.CountryQA: CountryQAT,
	bank.CountryRO: CountryROU,
	bank.CountryRW: CountryRWA,
	bank.CountrySA: CountrySAU,
	bank.CountrySB: CountrySLB,
	bank.CountrySL: CountrySLE,
	bank.CountrySM: CountrySMR,
	bank.CountryST: CountrySTP,
	bank.CountrySN: CountrySEN,
	bank.CountrySE: CountrySWE,
	bank.CountrySG: CountrySGP,
	bank.CountrySI: CountrySVN,
	bank.CountrySK: CountrySVK,
	bank.CountrySV: CountrySLV,
	bank.CountrySZ: CountrySWZ,
	bank.CountryTD: CountryTCD,
	bank.CountryTG: CountryTGO,
	bank.CountryTJ: CountryTJK,
	bank.CountryTL: CountryTLS,
	bank.CountryTN: CountryTUN,
	bank.CountryTO: CountryTON,
	bank.CountryTV: CountryTUV,
	bank.CountryTW: CountryTWN,
	bank.CountryTZ: CountryTZA,
	bank.CountryUG: CountryUGA,
	bank.CountryUS: CountryUSA,
	bank.CountryUY: CountryURY,
	bank.CountryUZ: CountryUZB,
	bank.CountryVA: CountryVAT,
	bank.CountryVN: CountryVNM,
	bank.CountryWS: CountryWSM,
	bank.CountryZA: CountryZAF,
	bank.CountryZM: CountryZMB,

	// High risk countries
	bank.CountryAF: CountryAFG,
	bank.CountryAG: CountryATG,
	bank.CountryAL: CountryALB,
	bank.CountryBA: CountryBIH,
	bank.CountryBS: CountryBHS,
	bank.CountryBY: CountryBLR,
	bank.CountryBW: CountryBWA,
	bank.CountryBZ: CountryBLZ,
	bank.CountryCF: CountryCAF,
	bank.CountryCN: CountryCHN,
	bank.CountryCD: CountryCOD,
	bank.CountryCO: CountryCOL,
	bank.CountryCR: CountryCRI,
	bank.CountryCU: CountryCUB,
	bank.CountryCY: CountryCYP,
	bank.CountryDM: CountryDMA,
	bank.CountryEG: CountryEGY,
	bank.CountryGH: CountryGHA,
	bank.CountryGW: CountryGNB,
	bank.CountryHK: CountryHKG,
	bank.CountryID: CountryIDN,
	bank.CountryIR: CountryIRN,
	bank.CountryIQ: CountryIRQ,
	bank.CountryIS: CountryISL,
	bank.CountryKE: CountryKEN,
	bank.CountryKH: CountryKHM,
	bank.CountryKN: CountryKNA,
	bank.CountryKP: CountryPRK,
	bank.CountryKZ: CountryKAZ,
	bank.CountryLB: CountryLBN,
	bank.CountryLC: CountryLCA,
	bank.CountryLY: CountryLBY,
	bank.CountryME: CountryMNE,
	bank.CountryMK: CountryMKD,
	bank.CountryMM: CountryMMR,
	bank.CountryML: CountryMLI,
	bank.CountryMN: CountryMNG,
	bank.CountryMO: CountryMAC,
	bank.CountryMY: CountryMYS,
	bank.CountryMZ: CountryMOZ,
	bank.CountryNG: CountryNGA,
	bank.CountryNI: CountryNIC,
	bank.CountryPA: CountryPAN,
	bank.CountryPH: CountryPHL,
	bank.CountryPK: CountryPAK,
	bank.CountryRS: CountrySRB,
	bank.CountryRU: CountryRUS,
	bank.CountrySD: CountrySDN,
	bank.CountrySO: CountrySOM,
	bank.CountrySS: CountrySSD,
	bank.CountrySY: CountrySYR,
	bank.CountryTH: CountryTHA,
	bank.CountryTR: CountryTUR,
	bank.CountryTT: CountryTTO,
	bank.CountryVC: CountryVCT,
	bank.CountryVE: CountryVEN,
	bank.CountryVU: CountryVUT,
	bank.CountryYE: CountryYEM,
	bank.CountryZW: CountryZWE,
}

var partnerCountryTo = map[Country3Alpha]bank.Country{
	CountryAIA: bank.CountryAI,
	CountryALA: bank.CountryAX,
	CountryAND: bank.CountryAD,
	CountryARE: bank.CountryAE,
	CountryARG: bank.CountryAR,
	CountryARM: bank.CountryAM,
	CountryAUT: bank.CountryAT,
	CountryAUS: bank.CountryAU,
	CountryAZE: bank.CountryAZ,
	CountryBEL: bank.CountryBE,
	CountryBEN: bank.CountryBJ,
	CountryBGD: bank.CountryBD,
	CountryBGR: bank.CountryBG,
	CountryBMU: bank.CountryBM,
	CountryBRA: bank.CountryBR,
	CountryBTN: bank.CountryBT,
	CountryCAN: bank.CountryCA,
	CountryCHE: bank.CountryCH,
	CountryCHL: bank.CountryCL,
	CountryCIV: bank.CountryCI,
	CountryCMR: bank.CountryCM,
	CountryCOG: bank.CountryCG,
	CountryCOK: bank.CountryCK,
	CountryCOM: bank.CountryKM,
	CountryCPV: bank.CountryCV,
	CountryCZE: bank.CountryCZ,
	CountryDEU: bank.CountryDE,
	CountryDJI: bank.CountryDJ,
	CountryDOM: bank.CountryDO,
	CountryDNK: bank.CountryDK,
	CountryDZA: bank.CountryDZ,
	CountryECU: bank.CountryEC,
	CountryERI: bank.CountryES,
	CountryEST: bank.CountryEE,
	CountryESP: bank.CountryES,
	CountryETH: bank.CountryES,
	CountryFIN: bank.CountryFI,
	CountryFJI: bank.CountryFJ,
	CountryFRA: bank.CountryFR,
	CountryGAB: bank.CountryGA,
	CountryGBR: bank.CountryGB,
	CountryGEO: bank.CountryGE,
	CountryGIB: bank.CountryGI,
	CountryGIN: bank.CountryGN,
	CountryGLP: bank.CountryGP,
	CountryGMB: bank.CountryGM,
	CountryGNQ: bank.CountryGQ,
	CountryGRC: bank.CountryGR,
	CountryGRD: bank.CountryGD,
	CountryGTM: bank.CountryGT,
	CountryGUY: bank.CountryGY,
	CountryHND: bank.CountryHN,
	CountryHTI: bank.CountryHT,
	CountryHUN: bank.CountryHU,
	CountryIRL: bank.CountryIE,
	CountryISR: bank.CountryIL,
	CountryIND: bank.CountryIN,
	CountryITA: bank.CountryIT,
	CountryJAM: bank.CountryJM,
	CountryJOR: bank.CountryJO,
	CountryJPN: bank.CountryJP,
	CountryKGZ: bank.CountryKG,
	CountryKIR: bank.CountryKI,
	CountryKOR: bank.CountryKR,
	CountryKWT: bank.CountryKW,
	CountryLAO: bank.CountryLA,
	CountryLBR: bank.CountryLR,
	CountryLIE: bank.CountryLI,
	CountryLKA: bank.CountryLK,
	CountryLSO: bank.CountryLS,
	CountryLTU: bank.CountryLT,
	CountryLUX: bank.CountryLU,
	CountryLVA: bank.CountryLV,
	CountryMAR: bank.CountryMA,
	CountryMCO: bank.CountryMC,
	CountryMDG: bank.CountryMG,
	CountryMDV: bank.CountryMV,
	CountryMEX: bank.CountryMX,
	CountryMHL: bank.CountryMH,
	CountryMLT: bank.CountryMT,
	CountryMRT: bank.CountryMR,
	CountryMSR: bank.CountryMS,
	CountryMTQ: bank.CountryMQ,
	CountryMUS: bank.CountryMU,
	CountryMWI: bank.CountryMW,
	CountryNAM: bank.CountryNA,
	CountryNER: bank.CountryNE,
	CountryNIU: bank.CountryNU,
	CountryNLD: bank.CountryNL,
	CountryNPL: bank.CountryNP,
	CountryNRU: bank.CountryNR,
	CountryNZL: bank.CountryNZ,
	CountryNOR: bank.CountryNO,
	CountryPER: bank.CountryPE,
	CountryPLW: bank.CountryPW,
	CountryPNG: bank.CountryPG,
	CountryPOL: bank.CountryPL,
	CountryPRT: bank.CountryPT,
	CountryPRY: bank.CountryPY,
	CountryPSE: bank.CountryPS,
	CountryQAT: bank.CountryQA,
	CountryROU: bank.CountryRO,
	CountryRWA: bank.CountryRW,
	CountrySAU: bank.CountrySA,
	CountrySEN: bank.CountrySN,
	CountrySGP: bank.CountrySG,
	CountrySLB: bank.CountrySB,
	CountrySLE: bank.CountrySL,
	CountrySLV: bank.CountrySV,
	CountrySMR: bank.CountrySM,
	CountrySTP: bank.CountryST,
	CountrySVK: bank.CountrySK,
	CountrySVN: bank.CountrySI,
	CountrySWE: bank.CountrySE,
	CountrySWZ: bank.CountrySZ,
	CountryTCD: bank.CountryTD,
	CountryTGO: bank.CountryTG,
	CountryTJK: bank.CountryTJ,
	CountryTON: bank.CountryTO,
	CountryTLS: bank.CountryTL,
	CountryTUN: bank.CountryTN,
	CountryTUV: bank.CountryTV,
	CountryTWN: bank.CountryTW,
	CountryTZA: bank.CountryTZ,
	CountryUGA: bank.CountryUG,
	CountryUSA: bank.CountryUS,
	CountryURY: bank.CountryUY,
	CountryUZB: bank.CountryUZ,
	CountryVAT: bank.CountryVA,
	CountryVNM: bank.CountryVN,
	CountryWSM: bank.CountryWS,
	CountryZAF: bank.CountryZA,
	CountryZMB: bank.CountryZM,

	// High risk countries
	CountryAFG: bank.CountryAF,
	CountryALB: bank.CountryAL,
	CountryATG: bank.CountryAG,
	CountryBIH: bank.CountryBA,
	CountryBHS: bank.CountryBS,
	CountryBLR: bank.CountryBY,
	CountryBLZ: bank.CountryBZ,
	CountryBWA: bank.CountryBW,
	CountryCAF: bank.CountryCF,
	CountryCHN: bank.CountryCN,
	CountryCOD: bank.CountryCD,
	CountryCOL: bank.CountryCO,
	CountryCRI: bank.CountryCR,
	CountryCUB: bank.CountryCU,
	CountryCYP: bank.CountryCY,
	CountryDMA: bank.CountryDM,
	CountryEGY: bank.CountryEG,
	CountryGHA: bank.CountryGH,
	CountryGNB: bank.CountryGW,
	CountryHKG: bank.CountryHK,
	CountryIDN: bank.CountryID,
	CountryIRN: bank.CountryIR,
	CountryIRQ: bank.CountryIQ,
	CountryISL: bank.CountryIS,
	CountryKAZ: bank.CountryKZ,
	CountryKEN: bank.CountryKE,
	CountryKHM: bank.CountryKH,
	CountryKNA: bank.CountryKN,
	CountryLBN: bank.CountryLB,
	CountryLBY: bank.CountryLY,
	CountryLCA: bank.CountryLC,
	CountryMAC: bank.CountryMO,
	CountryMKD: bank.CountryMK,
	CountryMMR: bank.CountryMM,
	CountryMLI: bank.CountryML,
	CountryMNE: bank.CountryME,
	CountryMNG: bank.CountryMN,
	CountryMOZ: bank.CountryMZ,
	CountryMYS: bank.CountryMY,
	CountryNGA: bank.CountryNG,
	CountryNIC: bank.CountryNI,
	CountryPAK: bank.CountryPK,
	CountryPAN: bank.CountryPA,
	CountryPHL: bank.CountryPH,
	CountryPRK: bank.CountryKP,
	CountryRUS: bank.CountryRU,
	CountrySDN: bank.CountrySD,
	CountrySOM: bank.CountrySO,
	CountrySRB: bank.CountryRS,
	CountrySSD: bank.CountrySS,
	CountrySYR: bank.CountrySY,
	CountryTHA: bank.CountryTH,
	CountryTTO: bank.CountryTT,
	CountryTUR: bank.CountryTR,
	CountryVCT: bank.CountryVC,
	CountryVEN: bank.CountryVE,
	CountryVUT: bank.CountryVU,
	CountryYEM: bank.CountryYE,
	CountryZWE: bank.CountryZW,
}
