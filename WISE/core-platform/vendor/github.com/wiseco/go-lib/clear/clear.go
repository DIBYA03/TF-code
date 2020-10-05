/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package clear

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/beevik/etree"
	"github.com/hoisie/mustache"
	"github.com/wiseco/go-lib/api"
	"github.com/wiseco/go-lib/log"
	"github.com/wiseco/go-lib/url"
)

// GLB, DPPA and VOTER are constant values to be passed to multiple clear APIs
// These values are defined by clear
const glb string = "B"
const dppa string = "3"
const voter string = "7"

type clearService struct {
	log    log.Logger
	client api.Client
}

// Service - To call clear apis and fetch data
type Service interface {
	RiskPersonSearch(RiskPersonSearchRequest) (*RiskPersonSearchResponse, error)
	RiskBusinessSearch(RiskBusinessSearchRequest) (*RiskBusinessSearchResponse, error)
}

// Address struct
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

// RiskPersonSearchRequest struct
type RiskPersonSearchRequest struct {
	FirstName            string
	MiddleName           string
	LastName             string
	BirthDate            string
	SocialSecurityNumber string
	Address              Address
	PhoneNumber          string
}

// RiskBusinessSearchRequest struct
type RiskBusinessSearchRequest struct {
	BusinessName string
	TaxID        string
	Address      Address
	PhoneNumber  string
}

type clearIdConfirmPersonSearch struct {
	GLB                  string
	DPPA                 string
	VOTER                string
	FirstName            string
	AddMiddleName        bool
	MiddleName           string
	LastName             string
	Day                  string
	Month                string
	Year                 string
	SocialSecurityNumber string
	Street               string
	City                 string
	State                string
	ZipCode              string
	PhoneNumber          string
}

func convertToSearchObject(req RiskPersonSearchRequest) clearIdConfirmPersonSearch {
	dtSlice := strings.Split(req.BirthDate, "-")
	newSsn := fmt.Sprintf("%s-%s-%s", req.SocialSecurityNumber[0:3], req.SocialSecurityNumber[3:5], req.SocialSecurityNumber[5:])
	var searchObj clearIdConfirmPersonSearch = clearIdConfirmPersonSearch{
		GLB:                  glb,
		DPPA:                 dppa,
		VOTER:                voter,
		FirstName:            req.FirstName,
		AddMiddleName:        len(req.MiddleName) != 0,
		MiddleName:           req.MiddleName,
		LastName:             req.LastName,
		SocialSecurityNumber: newSsn,
		Street:               req.Address.Street,
		City:                 req.Address.City,
		State:                req.Address.State,
		ZipCode:              req.Address.ZipCode,
		PhoneNumber:          req.PhoneNumber,
		Year:                 dtSlice[0],
		Month:                dtSlice[1],
		Day:                  dtSlice[2],
	}
	return searchObj
}

// RiskPersonSearchResponse struct
type RiskPersonSearchResponse struct {
	RespJson []byte
}

// RiskBusinessSearchResponse struct
type RiskBusinessSearchResponse struct {
	RespJson []byte
}

// NewClearService - Return a new instance of clear service
func NewClearService(l log.Logger, c api.Client) Service {
	return &clearService{log: l,
		client: c,
	}
}

func getUri(resp *string, rootElem string) (string, error) {
	/*
		Response is of the form
		<?xml version="1.0"?>
		<ns2:rootElem xmlns:ns2="http com/thomsonreuters/schemas/eidvsearch ">
		<Status>
		<StatusCode>200</StatusCode>
		<SubStatusCode>0</SubStatusCode>
		</Status>
		<Uri>{hostname}/api/v2/eidvperson/searchResults/946elsd3672q043ee</Uri>
		<GroupCount>1</GroupCount>
		</ns2:rootElem>
	*/
	uri := ""
	doc := etree.NewDocument()

	err := doc.ReadFromString(*resp)
	if err != nil {
		err = fmt.Errorf("Error occured while parsing xml response from clearID confirm person search. %s",
			err.Error())
		return uri, err
	}

	elem := doc.FindElement(fmt.Sprintf("%s/Status/StatusCode", rootElem))
	if elem == nil {
		err = fmt.Errorf("Failed to retrieve StatusCode from %s", rootElem)
		return uri, err
	}
	if elem.Text() != "200" {
		err = fmt.Errorf("Got status code %s for %s. Response %s", elem.Text(), rootElem, *resp)
		return uri, err
	}

	elem = doc.FindElement(fmt.Sprintf("%s/Uri", rootElem))
	if elem == nil {
		err = fmt.Errorf("Failed to retrieve URI from %s", rootElem)
		return uri, err
	}

	uri = elem.Text()
	uri = uri[len(url.BaseURLForClear())+1:]

	return uri, nil
}

type SearchResponseDetails struct {
	RootElement                 string
	EntityElementsXPathSelector string
	EntityElementsFilter        func(*etree.Element) bool
}

var SearchResponseDetailsKYCPerson = SearchResponseDetails{
	RootElement:                 "EIDVPersonSearchResponse",
	EntityElementsXPathSelector: "EIDVPersonSearchResponse/EIDVPersonSearchResults/EIDVPersonSearchResult/PersonEntities/PersonEntity",
	EntityElementsFilter: func(elem *etree.Element) bool {
		// Check if the person is dead. If so, then return false
		if d := elem.SelectElement("Death"); d != nil {
			if di := d.SelectElement("DeathIndicator"); di != nil {
				if strings.ToUpper(di.Text()) == "YES" {
					return false
				}
			}
		}
		return true
	},
}

var SearchResponseDetailsKYBBusiness = SearchResponseDetails{
	RootElement:                 "EIDVBusinessSearchResponse",
	EntityElementsXPathSelector: "EIDVBusinessSearchResponse/EIDVBusinessSearchResults/EIDVBusinessSearchResult/CompanyEntities/CompanyEntity",
	EntityElementsFilter:        nil,
}

func getEntityIDMatchingScoreMap(resp *string, sd SearchResponseDetails) (map[string]float32, error) {
	/*
		Response is of the form
		<ns2:RootElement xmlns:ns2="com/thomsonreuters/schemas/eidvsearch">
			<Status>
				<StatusCode>200</StatusCode>
				<SubStatusCode>200</SubStatusCode>
			</Status>
			<parent_elem>
				<nested_xml_struct>
					<Entity>
						<TotalScore>52.55</TotalScore>
						<EntityIdentifier>P1__NTI0ODU1MDEwMA</EntityIdentifier>
					</Entity>
				</nested_xml_struct>
			</parent_elem>
		</ns2:RootElement>
	*/
	res := make(map[string]float32)
	doc := etree.NewDocument()

	err := doc.ReadFromString(*resp)
	if err != nil {
		err = fmt.Errorf("Error occured while parsing xml response from clearID confirm person search. %s",
			err.Error())
		return res, err
	}

	elem := doc.FindElement(fmt.Sprintf("%s/Status/StatusCode", sd.RootElement))
	if elem == nil {
		err = fmt.Errorf("Failed to retrieve StatusCode from %s", sd.RootElement)
		return res, err
	}
	if elem.Text() != "200" {
		err = fmt.Errorf("Got status code %s for %s. Response %s", elem.Text(), sd.RootElement, *resp)
		return res, err
	}

	elems := doc.FindElements(sd.EntityElementsXPathSelector)
	for _, pe := range elems {
		if sd.EntityElementsFilter != nil {
			isValid := sd.EntityElementsFilter(pe)
			if !isValid {
				continue
			}
		}

		ts := pe.SelectElement("TotalScore")
		ei := pe.SelectElement("EntityIdentifier")
		if ts == nil || ei == nil {
			continue
		}

		totalscore := ts.Text()
		var totalscorenum float32 = 0
		if totalscore == "Yes" {
			totalscorenum = 100
		} else if totalscore == "No" {
			totalscorenum = 0
		} else if s, err := strconv.ParseFloat(totalscore, 32); err == nil {
			totalscorenum = float32(s)
		} else {
			continue
		}
		res[ei.Text()] = totalscorenum
	}

	return res, err
}

func GetConsumerClient(l log.Logger, certFilePath string, keyFilePath string) api.Client {
	// Clear user name
	userName := os.Getenv("CLEAR_USER_NAME")
	if len(userName) == 0 {
		panic("CLEAR_USER_NAME is missing")
	}
	// Clear user name
	password := os.Getenv("CLEAR_PWD")
	if len(password) == 0 {
		panic("CLEAR_PWD is missing")
	}
	// Create client
	return api.NewClient(
		api.ClientConfig{
			BaseURL:     url.BaseURLForClear(),
			ContentType: "application/xml",
			BasicAuth: &api.BasicAuthConfig{
				ClientID:     userName,
				ClientSecret: password,
			},
			TLSClientConfig: &api.TLSConfig{
				CertFilePath: certFilePath,
				KeyFilePath:  keyFilePath,
			},
		}, l, api.TLSClientBuilder)
}

func parsePhoneNumber(phNum string) (string, error) {
	outPhNum := ""
	var err error = nil
	if len(phNum) != 11 {
		err = fmt.Errorf("The phone number must be exactly 11 digits. Received %s", phNum)
	} else if phNum[0:1] != "1" {
		err = fmt.Errorf("The phone number must start with 1. Received %s", phNum)
	} else {
		outPhNum = fmt.Sprintf("%s-%s-%s", phNum[1:4], phNum[4:7], phNum[7:])
	}
	return outPhNum, err
}

func (s *clearService) RiskPersonSearch(request RiskPersonSearchRequest) (*RiskPersonSearchResponse, error) {
	if len(request.SocialSecurityNumber) != 9 {
		return nil, fmt.Errorf("The social security number must be exactly 9 digits. Received %s", request.SocialSecurityNumber)
	}
	phNum, err := parsePhoneNumber(request.PhoneNumber)
	if err != nil {
		return nil, err
	}
	request.PhoneNumber = phNum
	searchObj := convertToSearchObject(request)
	data := mustache.Render(entityPersonSearchReq, searchObj)
	r := strings.NewReader(data)
	var sptr *string = nil
	var resp **string = &(sptr)

	err = s.client.Post("api/v2/eidvperson/searchResults", r, resp, api.ResponseFormatterString)
	if err != nil {
		s.log.Error("Error occured while performing clearId confirm person search")
		return nil, err
	}
	var uri string
	uri, err = getUri(*resp, "EIDVPersonResults")
	if err != nil {
		s.log.Error("Error occured while extracting uri")
		return nil, err
	}
	sptr = nil
	resp = &(sptr)
	err = s.client.Get(uri, nil, resp, api.ResponseFormatterString)
	if err != nil {
		s.log.Error("Error occured while fetching clearId confirm person search results response")
		return nil, err
	}
	entityIDMatchingScoreMap, err := getEntityIDMatchingScoreMap(*resp, SearchResponseDetailsKYCPerson)
	if err != nil {
		s.log.Error("Error occured while fetching clear confirm person Id search results response")
		return nil, err
	}
	if len(entityIDMatchingScoreMap) == 0 {
		return nil, fmt.Errorf("No results found for %s", SearchResponseDetailsKYCPerson.RootElement)
	}
	matchingEntityID := ""
	var highestMatchingScore float32 = 0.0
	for entityID, matchingScore := range entityIDMatchingScoreMap {
		if matchingScore > highestMatchingScore {
			highestMatchingScore = matchingScore
			matchingEntityID = entityID
		}
	}
	riskPersonSearchObj := struct {
		GLB      string
		DPPA     string
		VOTER    string
		EntityID string
	}{
		glb,
		dppa,
		voter,
		matchingEntityID,
	}
	data = mustache.Render(riskPersonSearchReq, riskPersonSearchObj)
	r = strings.NewReader(data)
	sptr = nil
	resp = &(sptr)
	err = s.client.Post("api/v2/riskinformperson/searchResults", r, resp, api.ResponseFormatterString)
	if err != nil {
		s.log.Error("Error occured while performing risk inform person search")
		return nil, err
	}
	uri, err = getUri(*resp, "RiskInformPersonSearchResults")
	if err != nil {
		s.log.Error("Error occured while extracting uri")
		return nil, err
	}
	sptr = nil
	resp = &(sptr)
	err = s.client.Get(uri, nil, resp, api.ResponseFormatterString)
	if err != nil {
		s.log.Error("Error occured while fetching risk inform person search results response")
		return nil, err
	}
	doc := etree.NewDocument()
	err = doc.ReadFromString(**resp)
	if err != nil {
		return nil, fmt.Errorf("Error occured while fetching risk inform person search results response. %s",
			err.Error())
	}
	elem := doc.SelectElement("RiskInformPersonSearchResponse")
	matchingScoreElem := elem.CreateElement("MatchingScore")
	matchingScoreElem.SetText(fmt.Sprintf("%f", highestMatchingScore))
	var buf bytes.Buffer
	doc.WriteTo(&buf)
	json, err := xj.Convert(&buf)
	return &RiskPersonSearchResponse{
		RespJson: json.Bytes(),
	}, nil
}

type clearIdBusinessSearchResp struct {
	entityId      string
	matchingScore float32
}

func (s *clearService) clearIdConfirmBusinessSearch(request RiskBusinessSearchRequest) (clearIdBusinessSearchResp, error) {
	var result clearIdBusinessSearchResp

	phNum, err := parsePhoneNumber(request.PhoneNumber)
	if err != nil {
		return result, err
	}
	request.PhoneNumber = phNum

	eIDVName := os.Getenv("CLEAR_KYB_EIDVNAME")
	if len(eIDVName) == 0 {
		err = errors.New("CLEAR_KYB_EIDVNAME must be set")
		return result, err
	}

	searchObj := struct {
		*Address
		GLB          string
		DPPA         string
		VOTER        string
		BusinessName string
		TaxID        string
		PhoneNumber  string
		EIDVName     string
	}{
		&request.Address,
		glb,
		dppa,
		voter,
		request.BusinessName,
		request.TaxID,
		request.PhoneNumber,
		eIDVName,
	}
	data := mustache.Render(clearIdBusinessSearchReq, searchObj)
	s.log.Debug("clearIdConfirmBusinessSearch Post request " + data)

	r := strings.NewReader(data)
	var sptr *string = nil
	var resp **string = &(sptr)

	err = s.client.Post("api/v2/eidvbusiness/searchResults", r, resp, api.ResponseFormatterString)
	if err != nil {
		err = fmt.Errorf("Error occured while performing clearId confirm business search. %s", err.Error())
		s.log.Error(err.Error())
		return result, err
	}

	s.log.Debug("clearIdConfirmBusinessSearch Post response " + **resp)
	var uri string
	uri, err = getUri(*resp, "EIDVBusinessResults")
	if err != nil {
		err = fmt.Errorf("Error occured while extracting uri from clearIdConfirmBusinessSearch response. %s", err.Error())
		s.log.Error(err.Error())
		return result, err
	}

	sptr = nil
	resp = &(sptr)
	err = s.client.Get(uri, nil, resp, api.ResponseFormatterString)
	if err != nil {
		err = fmt.Errorf("Error occured while fetching clearID confirm business search results response. %s", err.Error())
		s.log.Error(err.Error())
		return result, err
	}
	s.log.Debug("Get confirm business search results " + **resp)

	entityIDMatchingScoreMap, err := getEntityIDMatchingScoreMap(*resp, SearchResponseDetailsKYBBusiness)
	if err != nil {
		err = fmt.Errorf("Error occured while fetching clearID confirm business search results response. %s", err.Error())
		s.log.Error(err.Error())
		return result, err
	}
	if len(entityIDMatchingScoreMap) == 0 {
		err = fmt.Errorf("No results found for %s", SearchResponseDetailsKYCPerson.RootElement)
		return result, err
	}

	matchingEntityID := ""
	var highestMatchingScore float32 = 0.0
	for entityID, matchingScore := range entityIDMatchingScoreMap {
		if matchingScore > highestMatchingScore {
			highestMatchingScore = matchingScore
			matchingEntityID = entityID
		}
	}
	result = clearIdBusinessSearchResp{
		entityId:      matchingEntityID,
		matchingScore: highestMatchingScore,
	}
	return result, err
}

func (s *clearService) RiskBusinessSearch(request RiskBusinessSearchRequest) (*RiskBusinessSearchResponse, error) {
	var sptr *string = nil
	var resp **string = &(sptr)

	clrSearchRes, err := s.clearIdConfirmBusinessSearch(request)
	if err != nil {
		return nil, err
	}

	riskInformDefName := os.Getenv("CLEAR_KYB_RISK_INFORM_DEF_NAME")
	if len(riskInformDefName) == 0 {
		err = errors.New("CLEAR_KYB_RISK_INFORM_DEF_NAME must be set")
		return nil, err
	}

	riskBusinessSearchObj := struct {
		GLB               string
		DPPA              string
		VOTER             string
		EntityID          string
		RiskInformDefName string
	}{
		glb,
		dppa,
		voter,
		clrSearchRes.entityId,
		riskInformDefName,
	}
	data := mustache.Render(riskInformBusinessSearchReq, riskBusinessSearchObj)
	s.log.Debug("RiskInformBusiness search request " + data)

	r := strings.NewReader(data)
	err = s.client.Post("api/v2/riskinformbusiness/searchResults", r, resp, api.ResponseFormatterString)
	if err != nil {
		err = fmt.Errorf("Error occured while performing risk inform business search. %s", err.Error())
		s.log.Error(err.Error())
		return nil, err
	}
	s.log.Debug("RiskInformBusiness search resp " + **resp)

	uri, err := getUri(*resp, "RiskInformBusinessSearchResults")
	if err != nil {
		err = fmt.Errorf("Error occured while extracting uri. %s", err.Error())
		s.log.Error(err.Error())
		return nil, err
	}

	sptr = nil
	resp = &(sptr)
	err = s.client.Get(uri, nil, resp, api.ResponseFormatterString)
	if err != nil {
		err = fmt.Errorf("Error occured while fetching risk inform business search results response")
		s.log.Error(err.Error())
		return nil, err
	}
	s.log.Debug("risk inform business search resp " + **resp)

	doc := etree.NewDocument()
	err = doc.ReadFromString(**resp)
	if err != nil {
		err = fmt.Errorf("Error occured while parsing risk inform business search results response. %s",
			err.Error())
		s.log.Error(err.Error())
		return nil, err
	}

	elem := doc.SelectElement("RiskInformBusinessSearchResponse")
	matchingScoreElem := elem.CreateElement("MatchingScore")
	matchingScoreElem.SetText(fmt.Sprintf("%f", clrSearchRes.matchingScore))

	var buf bytes.Buffer
	doc.WriteTo(&buf)
	json, err := xj.Convert(&buf)
	if err != nil {
		return nil, err
	}

	return &RiskBusinessSearchResponse{
		RespJson: json.Bytes(),
	}, nil
}
