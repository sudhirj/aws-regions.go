package awsregions_test

import (
	awsregions "github.com/sudhirj/aws-regions.go"
	"testing"
	"time"
)

func TestLatencyChecks(t *testing.T) {
	lc := awsregions.NewLatencyChecker("ap-south-1", "us-east-1", "ap-southeast-1", "eu-west-1")
	go lc.Start()
	go func(t1 *testing.T) {
		for {
			t1.Log(lc.FastestRegion())
			t1.Log(lc.Latencies())
			t1.Log(lc.SortedRegions())
			time.Sleep(5 * time.Second)
		}
	}(t)

	time.Sleep(15 * time.Second)
}
