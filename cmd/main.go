package main

import (
	awsregions "github.com/sudhirj/aws-regions.go"
	"log"
	"time"
)

func main() {
	lc := awsregions.NewLatencyChecker("ap-south-1", "us-east-1", "ap-southeast-1", "eu-west-1")
	lc.Measure()
	log.Println("initial measurement: ", lc.FastestRegion())
	go lc.Start()
	for {
		log.Println("running sorted", lc.SortedRegions())
		log.Println("running fastest", lc.FastestRegion())
		log.Println("running latencies", lc.Latencies())
		time.Sleep(5 * time.Second)
	}
}
