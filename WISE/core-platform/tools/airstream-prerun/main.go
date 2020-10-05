package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/business"
	cspBusiness "github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/services/csp/consumer"
	"github.com/wiseco/core-platform/shared"
)

const (
	envDev     = "dev"
	envStg     = "stg"
	envStaging = "staging"
	envQA      = "qa"
	envSbx     = "sbx"
	envPrd     = "prd"
	envProd    = "prod"
)

const (
	startDate = "2020-06-01"
	endDate   = "2020-06-30"
)

func main() {
	consumers := getConsumers()
	time1 := time.Now()

	if consumers != nil {
		runKYC(consumers)
	}

	time2 := time.Now()

	businesses := getBusinesses()
	if businesses != nil {
		runKYB(businesses)
	}

	time3 := time.Now()

	kycTime := time2.Sub(time1)
	kybTime := time3.Sub(time2)

	fmt.Printf("\n\n KYC: %v", kycTime.Seconds())
	fmt.Printf("\n\n KYB: %v", kybTime.Seconds())
}

// KYC
func runKYC(consumers []consumer.CSPConsumer) {
	clearDuration := kycDuration{}.init()
	alloyDuration := kycDuration{}.init()
	phoneDuration := kycDuration{}.init()
	fmt.Printf("\nRUNNING KYC FOR %v CONSUMERS\n\n", len(consumers))

	offset := 0
	limit := 10
	for {
		start := offset
		end := offset + limit

		if len(consumers) < end {
			end = len(consumers)
			l := len(consumers) - offset
			if l <= 0 {
				break
			}
		}

		slicedConsumers := consumers[start:end]
		if len(slicedConsumers) == 0 {
			log.Println("no more consumers:", offset, 10)
			break
		}

		var wg sync.WaitGroup
		wg.Add(len(slicedConsumers))
		for _, c := range slicedConsumers {
			go func(c consumer.CSPConsumer) {
				defer wg.Done()
				// CLEAR
				runClearKYC(c.ConsumerID, clearDuration)

				// ALLOY
				runAlloyKYC(c.ConsumerID, alloyDuration)

				// PHONE VERIFICATION
				runPhoneVerification(c.ConsumerID, phoneDuration)
			}(c)
		}

		wg.Wait()
		offset += len(slicedConsumers)
	}

	clearDuration.print("CLEAR")
	alloyDuration.print("ALLOY")
	phoneDuration.print("PHONE VERIFICATION")
}

func runClearKYC(consumerID shared.ConsumerID, d kycDuration) {
	startTime := time.Now()
	_, err := consumer.New().GetClearKYC(consumerID)
	if err == nil {
		fmt.Println("Clear KYC already done for Cosnumer: ", consumerID)
		d.skip(startTime)
		return
	}

	_, err = consumer.New().RunClearKYC(consumerID)
	if err != nil {
		fmt.Printf("CLEAR/ consumer: %v, ERR: %v \n", consumerID, err)
		d.addFailure(startTime)
	} else {
		d.addSuccess(startTime)
	}
}

func runAlloyKYC(consumerID shared.ConsumerID, d kycDuration) {
	startTime := time.Now()
	_, err := consumer.New().GetAlloyKYC(consumerID)
	if err == nil {
		fmt.Println("Alloy KYC already done for Cosnumer: ", consumerID)
		d.skip(startTime)
		return
	}

	_, err = consumer.New().RunAlloyKYC(consumerID)
	if err != nil {
		fmt.Printf("ALLOY/ consumer: %v, ERR: %v \n", consumerID, err)
		d.addFailure(startTime)
	} else {
		d.addSuccess(startTime)
	}
}

func runPhoneVerification(consumerID shared.ConsumerID, d kycDuration) {
	startTime := time.Now()
	_, err := consumer.New().PhoneVerification(shared.ConsumerID(consumerID))
	if err != nil {
		fmt.Printf("PHONE verification/ consumer: %v, ERR: %v \n", consumerID, err)
		d.addFailure(startTime)
	} else {
		d.addSuccess(startTime)
	}
}

// KYB
func runKYB(businesses []business.CSPBusiness) {
	clearDuration := kycDuration{}.init()

	fmt.Printf("\nRUNNING KYB FOR %v Business\n\n", len(businesses))

	offset := 0
	limit := 10
	for {
		l := limit
		if len(businesses) < offset+limit {
			l = len(businesses) - offset
			if l <= 0 {
				break
			}
		}
		start := offset
		end := offset + l
		slicedBusinesses := businesses[start:end]

		if len(slicedBusinesses) == 0 {
			log.Println("no more businesses:", offset, 10)
			break
		}

		var wg sync.WaitGroup
		wg.Add(len(slicedBusinesses))
		for _, b := range slicedBusinesses {

			go func(b business.CSPBusiness) {
				defer wg.Done()

				startTime := time.Now()
				_, err := cspBusiness.New(services.NewSourceRequest()).GetClearVerification(b.BusinessID)
				if err == nil {
					fmt.Println("Clear KYB already done for Business: ", b.BusinessID)
					clearDuration.skip(startTime)
					return
				}

				_, err = cspBusiness.New(services.NewSourceRequest()).RunClearVerification(b.BusinessID)
				if err != nil {
					fmt.Printf("CLEAR/ business: %v, ERR: %v \n", b.BusinessID, err)
					clearDuration.addFailure(startTime)
				} else {
					clearDuration.addSuccess(startTime)
				}

			}(b)
		}

		wg.Wait()

		offset += len(slicedBusinesses)
	}
	clearDuration.print("CLEAR KYB")
}

//Helpers
func getConsumers() []consumer.CSPConsumer {
	params := make(map[string]interface{})
	params["submitStart"] = startDate
	params["submitEnd"] = endDate

	consumers, err := consumer.NewCSPService().GetAll(params)
	if err != nil {
		fmt.Println("Unable to retrive consumers: ", err)
		return nil
	}
	return consumers
}

func getBusinesses() []business.CSPBusiness {
	params := make(map[string]interface{})
	params["submitStart"] = startDate
	params["submitEnd"] = endDate

	list, err := business.NewCSPService().CSPBusinessList(params)
	if err != nil {
		fmt.Println("Unable to retrive consumers: ", err)
		return nil
	}
	return list
}

type kycDuration struct {
	skipCount    *int
	failureCount *int
	successCount *int

	skipTime    *time.Duration
	failureTime *time.Duration
	successTime *time.Duration
}

func (d kycDuration) init() kycDuration {
	skipCcount := 0
	sCount := 0
	fCount := 0

	sTime := time.Duration(0)
	fTime := time.Duration(0)
	skipTime := time.Duration(0)

	d.skipCount = &skipCcount
	d.failureCount = &sCount
	d.successCount = &fCount

	d.skipTime = &skipTime
	d.successTime = &sTime
	d.failureTime = &fTime
	return d
}

func (d kycDuration) addSuccess(startTime time.Time) {
	*d.successCount = (*d.successCount + 1)
	*d.successTime = (*d.successTime + time.Now().Sub(startTime))
	fmt.Printf("\nSUCCESS COUNT: %d, time: %v\n", *d.successCount, *d.successTime)
}

func (d kycDuration) addFailure(startTime time.Time) {
	*d.failureCount = (*d.failureCount + 1)
	*d.failureTime = (*d.failureTime + time.Now().Sub(startTime))
	fmt.Printf("\nFAILURE COUNT: %d, time: %v\n", *d.failureCount, *d.failureTime)
}

func (d kycDuration) skip(startTime time.Time) {
	*d.skipCount = (*d.skipCount + 1)
	*d.skipTime = (*d.skipTime + time.Now().Sub(startTime))
	fmt.Printf("\nSKIP COUNT: %d, time: %v\n", *d.skipCount, *d.skipTime)

}

func (d kycDuration) print(title string) {
	fmt.Println("\n" + title)
	skipCount := *d.skipCount
	sCount := *d.successCount
	fCount := *d.failureCount

	skipTime := *d.skipTime
	sTime := *d.successTime
	fTime := *d.failureTime

	fmt.Printf("Skip time: %.2f,	Count %v, Average: %.2f\n", skipTime.Seconds(), skipCount, skipTime.Seconds()/float64(skipCount))
	fmt.Printf("Success time: %.2f,	Count %v, Average: %.2f\n", sTime.Seconds(), sCount, sTime.Seconds()/float64(sCount))
	fmt.Printf("Failure time: %.2f,	Count %v, Average: %.2f\n\n", fTime.Seconds(), fCount, fTime.Seconds()/float64(fCount))
}
