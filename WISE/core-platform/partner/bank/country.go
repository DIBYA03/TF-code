package bank

type Country string

const (
	CountryAD = Country("AD")
	CountryAE = Country("AE")
	CountryAI = Country("AI")
	CountryAM = Country("AM")
	CountryAR = Country("AR")
	CountryAT = Country("AT")
	CountryAU = Country("AU")
	CountryAX = Country("AX")
	CountryAZ = Country("AZ")
	CountryBD = Country("BD")
	CountryBE = Country("BE")
	CountryBF = Country("BF")
	CountryBG = Country("BG")
	CountryBJ = Country("BJ")
	CountryBM = Country("BM")
	CountryBN = Country("BN")
	CountryBR = Country("BR")
	CountryBT = Country("BT")
	CountryCA = Country("CA")
	CountryCG = Country("CG")
	CountryCH = Country("CH")
	CountryCI = Country("CI")
	CountryCK = Country("CK")
	CountryCL = Country("CL")
	CountryCM = Country("CM")
	CountryCV = Country("CV")
	CountryCZ = Country("CZ")
	CountryDE = Country("DE")
	CountryDJ = Country("DJ")
	CountryDO = Country("DO")
	CountryDK = Country("DK")
	CountryDZ = Country("DZ")
	CountryEC = Country("EC")
	CountryEE = Country("EE")
	CountryER = Country("ER")
	CountryES = Country("ES")
	CountryET = Country("ET")
	CountryFI = Country("FI")
	CountryFJ = Country("FJ")
	CountryFM = Country("FM")
	CountryFR = Country("FR")
	CountryGA = Country("GA")
	CountryGB = Country("GB")
	CountryGD = Country("GD")
	CountryGE = Country("GE")
	CountryGI = Country("GI")
	CountryGM = Country("GM")
	CountryGN = Country("GN")
	CountryGP = Country("GP")
	CountryGQ = Country("GQ")
	CountryGR = Country("GR")
	CountryGT = Country("GT")
	CountryGY = Country("GY")
	CountryHN = Country("HN")
	CountryHT = Country("HT")
	CountryHU = Country("HU")
	CountryIE = Country("IE")
	CountryIL = Country("IL")
	CountryIN = Country("IN")
	CountryIT = Country("IT")
	CountryJM = Country("JM")
	CountryJO = Country("JO")
	CountryJP = Country("JP")
	CountryKI = Country("KI")
	CountryKG = Country("KG")
	CountryKM = Country("KM")
	CountryKR = Country("KR")
	CountryKW = Country("KW")
	CountryLA = Country("LA")
	CountryLI = Country("LI")
	CountryLK = Country("LK")
	CountryLR = Country("LR")
	CountryLS = Country("LS")
	CountryLT = Country("LT")
	CountryLU = Country("LU")
	CountryLV = Country("LV")
	CountryMA = Country("MA")
	CountryMC = Country("MC")
	CountryMG = Country("MG")
	CountryMH = Country("MH")
	CountryMQ = Country("MQ")
	CountryMR = Country("MR")
	CountryMS = Country("MS")
	CountryMT = Country("MT")
	CountryMU = Country("MU")
	CountryMV = Country("MV")
	CountryMW = Country("MW")
	CountryMX = Country("MX")
	CountryNA = Country("NA")
	CountryNE = Country("NE")
	CountryNL = Country("NL")
	CountryNO = Country("NO")
	CountryNP = Country("NP")
	CountryNR = Country("NR")
	CountryNU = Country("NU")
	CountryNZ = Country("NZ")
	CountryOM = Country("OM")
	CountryPE = Country("PE")
	CountryPG = Country("PG")
	CountryPL = Country("PL")
	CountryPS = Country("PS")
	CountryPT = Country("PT")
	CountryPW = Country("PW")
	CountryPY = Country("PY")
	CountryQA = Country("QA")
	CountryRO = Country("RO")
	CountryRW = Country("RW")
	CountrySA = Country("SA")
	CountrySB = Country("SB")
	CountrySE = Country("SE")
	CountrySG = Country("SG")
	CountrySI = Country("SI")
	CountrySK = Country("SK")
	CountrySL = Country("SL")
	CountrySM = Country("SM")
	CountrySN = Country("SN")
	CountrySR = Country("SR")
	CountryST = Country("ST")
	CountrySV = Country("SV")
	CountrySZ = Country("SZ")
	CountryTD = Country("TD")
	CountryTG = Country("TG")
	CountryTJ = Country("TJ")
	CountryTL = Country("TL")
	CountryTN = Country("TN")
	CountryTO = Country("TO")
	CountryTV = Country("TV")
	CountryTW = Country("TW")
	CountryTZ = Country("TZ")
	CountryUG = Country("UG")
	CountryUS = Country("US")
	CountryUY = Country("UY")
	CountryUZ = Country("UZ")
	CountryWS = Country("WS")
	CountryVA = Country("VA")
	CountryVN = Country("VN")
	CountryZM = Country("ZM")
	CountryZA = Country("ZA")

	// High risk countries
	CountryAF = Country("AF")
	CountryAG = Country("AG")
	CountryAL = Country("AL")
	CountryBA = Country("BA")
	CountryBB = Country("BB")
	CountryBH = Country("BH")
	CountryBI = Country("BI")
	CountryBS = Country("BS")
	CountryBW = Country("BW")
	CountryBY = Country("BY")
	CountryBZ = Country("BZ")
	CountryCD = Country("CD")
	CountryCF = Country("CF")
	CountryCN = Country("CN")
	CountryCO = Country("CO")
	CountryCR = Country("CR")
	CountryCU = Country("CU")
	CountryCY = Country("CY")
	CountryDM = Country("DM")
	CountryEG = Country("EG")
	CountryGH = Country("GH")
	CountryGW = Country("GW")
	CountryHK = Country("HK")
	CountryID = Country("ID")
	CountryIQ = Country("IQ")
	CountryIR = Country("IR")
	CountryIS = Country("IS")
	CountryKE = Country("KE")
	CountryKH = Country("KH")
	CountryKN = Country("KN")
	CountryKP = Country("KP")
	CountryKZ = Country("KZ")
	CountryLB = Country("LB")
	CountryLC = Country("LC")
	CountryLY = Country("LY")
	CountryME = Country("ME")
	CountryMK = Country("MK")
	CountryML = Country("ML")
	CountryMM = Country("MM")
	CountryMN = Country("MN")
	CountryMO = Country("MO")
	CountryMY = Country("MY")
	CountryMZ = Country("MZ")
	CountryNI = Country("NI")
	CountryNG = Country("NG")
	CountryPA = Country("PA")
	CountryPH = Country("PH")
	CountryPK = Country("PK")
	CountryRS = Country("RS")
	CountryRU = Country("RU")
	CountrySD = Country("SD")
	CountrySO = Country("SO")
	CountrySS = Country("SS")
	CountrySY = Country("SY")
	CountryTH = Country("TH")
	CountryTR = Country("TR")
	CountryTT = Country("TT")
	CountryUA = Country("UA")
	CountryVC = Country("VC")
	CountryVE = Country("VE")
	CountryVU = Country("VU")
	CountryYE = Country("YE")
	CountryZW = Country("ZW")
)

var highRiskCountries = map[Country]bool{
	CountryAF: true,
	CountryAG: true,
	CountryAL: true,
	CountryBA: true,
	CountryBB: true,
	CountryBH: true,
	CountryBI: true,
	CountryBS: true,
	CountryBW: true,
	CountryBY: true,
	CountryBZ: true,
	CountryCD: true,
	CountryCF: true,
	CountryCN: true,
	CountryCO: true,
	CountryCR: true,
	CountryCU: true,
	CountryCY: true,
	CountryDM: true,
	CountryEG: true,
	CountryGH: true,
	CountryGW: true,
	CountryHK: true,
	CountryID: true,
	CountryIQ: true,
	CountryIR: true,
	CountryIS: true,
	CountryKE: true,
	CountryKH: true,
	CountryKN: true,
	CountryKP: true,
	CountryKZ: true,
	CountryLB: true,
	CountryLC: true,
	CountryLY: true,
	CountryME: true,
	CountryMK: true,
	CountryML: true,
	CountryMM: true,
	CountryMN: true,
	CountryMO: true,
	CountryMY: true,
	CountryMZ: true,
	CountryNI: true,
	CountryNG: true,
	CountryPA: true,
	CountryPH: true,
	CountryPK: true,
	CountryRS: true,
	CountryRU: true,
	CountrySD: true,
	CountrySO: true,
	CountrySS: true,
	CountrySY: true,
	CountryTH: true,
	CountryTR: true,
	CountryTT: true,
	CountryUA: true,
	CountryVC: true,
	CountryVE: true,
	CountryVU: true,
	CountryYE: true,
	CountryZW: true,
}

func IsHighRiskCountry(c Country) bool {
	_, ok := highRiskCountries[c]
	return ok
}
