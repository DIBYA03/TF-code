/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package analytics

import (
	"log"
	"testing"
)

func TestMetricsAnalytics(t *testing.T) {
	resp, err := NewAnlytics().Metrics()
	if err != nil {
		t.Errorf("Fail, error getting metrics %v", err)
	}
	log.Printf("Response from metrics %v ", resp)
}
