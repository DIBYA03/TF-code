package mock

import platformservices "github.com/wiseco/core-platform/services"

// NewAddress returns a mock address
func NewAddress() platformservices.Address {

	var lat = 37.7917146
	var lng = -122.397054

	return platformservices.Address{
		AddressType:   platformservices.AddressTypeLegal,
		StreetAddress: "100 Main Street",
		City:          "San Francisco",
		State:         "CA",
		Country:       "US",
		PostalCode:    "94100",
		Latitude:      &lat,
		Longitude:     &lng,
	}
}
