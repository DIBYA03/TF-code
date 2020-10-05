package clear

const entityPersonSearchReq string = `<?xml version="1.0"?>
<p:EIDVPersonSearch
xmlns:p="com/thomsonreuters/schemas/eidvsearch"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="com/thomsonreuters/schemas/eidvsearch eidvsearch.xsd">
<PermissiblePurpose>
<GLB>{{GLB}}</GLB>
<DPPA>{{DPPA}}</DPPA>
<VOTER>{{VOTER}}</VOTER>
</PermissiblePurpose>
<EIDVPersonSearchRequest>
<EIDVName>Wise-person</EIDVName>
<EIDVVersion></EIDVVersion>
<Criteria>
<FirstName>{{FirstName}}</FirstName>
{{#AddMiddleName}}<MiddleName>{{MiddleName}}</MiddleName>{{/AddMiddleName}}
<LastName>{{LastName}}</LastName>
<BirthDate>
<Day>{{Day}}</Day>
<Month>{{Month}}</Month>
<Year>{{Year}}</Year>
</BirthDate>
<SSN>{{SocialSecurityNumber}}</SSN>
<Street>{{Street}}</Street>
<City>{{City}}</City>
<State>{{State}}</State>
<ZipCode>{{ZipCode}}</ZipCode>
<PhoneNumber>{{PhoneNumber}}</PhoneNumber>
</Criteria>
</EIDVPersonSearchRequest>
</p:EIDVPersonSearch>`
const riskPersonSearchReq string = `<?xml version="1.0"?>
<rips:RiskInformPersonSearchRequest xmlns:rips="http://clear.thomsonreuters.com/api/search/2.0">
<PermissiblePurpose>
<GLB>{{GLB}}</GLB>
<DPPA>{{DPPA}}</DPPA>
<VOTER>{{VOTER}}</VOTER>
</PermissiblePurpose>
<Criteria>
<rip1:RiskInformPersonSearchCriteria xmlns:rip1="com/thomsonreuters/schemas/riskinformperson-search">
<RiskInformDefName>Wise-person</RiskInformDefName>
<RiskInformDefVersion></RiskInformDefVersion>
<EntityId>{{EntityID}}</EntityId>
<IncludeAdditionalSearches>
<WebAnalyticsSearch>true</WebAnalyticsSearch>
</IncludeAdditionalSearches>
</rip1:RiskInformPersonSearchCriteria>
</Criteria>
</rips:RiskInformPersonSearchRequest>`
const clearIdBusinessSearchReq string = `<?xml version="1.0" encoding="utf-8"?>
<EIDVBusinessSearch xmlns:ns1="com/thomsonreuters/schemas/eidvsearch" xmlns:ns2="http://www.w3.org/2001/XMLSchema-instance" xmlns="com/thomsonreuters/schemas/eidvbusiness-search">
<PermissiblePurpose xmlns="">
<GLB>{{GLB}}</GLB>
<DPPA>{{DPPA}}</DPPA>
<VOTER>{{VOTER}}</VOTER>
</PermissiblePurpose>
<EIDVBusinessSearchRequest xmlns="">
<EIDVName>{{EIDVName}}</EIDVName>
<EIDVVersion></EIDVVersion>
<Criteria>
<BusinessName>{{BusinessName}}</BusinessName>
<FeinNumber>{{TaxID}}</FeinNumber>
<Street>{{Street}}</Street>
<City>{{City}}</City>
<State>{{State}}</State>
<ZipCode>{{ZipCode}}</ZipCode>
<PhoneNumber>{{PhoneNumber}}</PhoneNumber>
</Criteria>
</EIDVBusinessSearchRequest>
</EIDVBusinessSearch>`

const riskInformBusinessSearchReq = `<?xml version="1.0"?>
<ribs:RiskInformBusinessSearchRequest xmlns:ribs="http://clear.thomsonreuters.com/api/search/2.0">
<PermissiblePurpose>
<GLB>{{GLB}}</GLB>
<DPPA>{{DPPA}}</DPPA>
<VOTER>{{VOTER}}</VOTER>
</PermissiblePurpose>
<Criteria>
<rib1:RiskInformBusinessSearchCriteria xmlns:rib1="com/thomsonreuters/schemas/riskinformbusiness-search">
<RiskInformDefName>{{RiskInformDefName}}</RiskInformDefName>
<RiskInformDefVersion></RiskInformDefVersion>
<EntityId>{{EntityID}}</EntityId>
<IncludeAdditionalSearches>
<WebAnalyticsSearch>true</WebAnalyticsSearch>
</IncludeAdditionalSearches>
</rib1:RiskInformBusinessSearchCriteria>
</Criteria>
</ribs:RiskInformBusinessSearchRequest>`
