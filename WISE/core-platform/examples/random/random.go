package random

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func UUID() string {
	return uuid.New().String()
}
func Name() string {
	return gofakeit.Name() // Markus Moen
}
func FirstName() string {
	return gofakeit.FirstName() // Markus
}
func LastName() string {
	return gofakeit.LastName() // Moen
}
func Number(min, max int) int {
	return gofakeit.Number(min, max)
}
func Email() string {
	return gofakeit.Email() //alaynawuckert@kozey.biz
}
func Phone() string {
	return fmt.Sprintf("+%d%d%d%d", 1, 999, Number(100, 999), Number(1000, 9999)) // +14153214567
}
func Company() string {
	return gofakeit.Company() // Moen, Pagac and Wuckert
}
func CreditCardNumber() int {
	return gofakeit.CreditCardNumber() // 4287271570245748
}
func AccountNumber() string {
	return fmt.Sprintf("%d", Number(1000000000, 9999999999))
}
func RoutingNumber() string {
	return fmt.Sprintf("%d", Number(100000000, 999999999))
}
func WireRouting() string {
	return fmt.Sprintf("%d", Number(100000000, 999999999))
}
func Alias() string {
	return gofakeit.BeerName()
}
func SSN() string {
	return fmt.Sprintf("%d%d%d", Number(100, 999), Number(10, 99), Number(1000, 9999)) // "123465789"
}
func DateOfBirth() string {
	return fmt.Sprintf("%d-%02d-%02d", Number(1950, 1990), Number(1, 12), Day())
}
func Year() int {
	return gofakeit.Year()
}
func Day() int {
	return Number(1, 28)
}
func Month() time.Month {
	return time.Month(Number(1, 12))
}
func IP() string {
	return fmt.Sprintf("%d.%d.%d.%d", Number(50, 60), Number(150, 255), Number(110, 220), Number(100, 199))
}

/* func DriversLicense() partnerbank.ConsumerIdentificationRequest {
	return partnerbank.ConsumerIdentificationRequest{
		Type:           "primary",
		Document:       partnerbank.IdentificationDocumentDriversLicense,
		Number:         fmt.Sprintf("%d", Number(123456789, 987654321)),
		IssuingState:   "CA",
		IssuingCountry: "USA",
		IssuedDate:     fmt.Sprintf("%d-%02d-%02d", (time.Now().Year() - 2), Month(), Day()),
		ExpirationDate: fmt.Sprintf("%d-%02d-%02d", (time.Now().Year() + 2), Month(), Day()),
	}
} */

func Address() partnerbank.AddressRequest {
	return partnerbank.AddressRequest{
		Type:    "LEGAL",
		Line1:   "201 mission st",
		City:    "san francisco",
		State:   "CA",
		ZipCode: "94104",
	}
}
