package main

import (
	awsregions "github.com/sudhirj/aws-regions.go"
	"log"
	"time"
)

func main() {
	lc := awsregions.NewLatencyChecker("ap-south-1", "us-east-1", "ap-southeast-1", "eu-west-1")
	go lc.Start()
	for {
		log.Println(lc.SortedRegions())
		log.Println(lc.FastestRegion())
		log.Println(lc.Latencies())
		time.Sleep(5 * time.Second)
	}
}
