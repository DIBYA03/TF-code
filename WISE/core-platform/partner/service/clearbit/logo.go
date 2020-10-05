package clearbit

import (
	"errors"
	"os"
	"strings"

	"github.com/clearbit/clearbit-go/clearbit"
)

type Logo struct {
	Name   string
	Domain string
	URL    string
}

func GetLogo(name string) (*Logo, error) {
	k := os.Getenv("CLEARBIT_API_KEY")
	if k == "" {
		panic(errors.New("clearbit api key missing"))
	}

	c := clearbit.NewClient(clearbit.WithAPIKey(k))

	n := strings.ToUpper(name)
	res, _, err := c.NameToDomain.Find(clearbit.NameToDomainFindParams{Name: n})
	if err != nil {
		return nil, err
	}

	return &Logo{
		Name:   n,
		Domain: res.Domain,
		URL:    res.Logo,
	}, nil
}
