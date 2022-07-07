package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type fooCollector struct {
	fooMetric *prometheus.Desc
	barMetric *prometheus.Desc
	channel   chan uint64
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newFooCollector() *fooCollector {
	return &fooCollector{
		fooMetric: prometheus.NewDesc("foo_metric",
			"Shows whether a foo has occurred in our cluster",
			nil, nil,
		),
		barMetric: prometheus.NewDesc("bar_metric",
			"Shows whether a bar has occurred in our cluster",
			nil, nil,
		),
		channel: make(chan uint64),
	}
}

// func Autocall() {
// }

func add(x, y int) (z int) {
	z = x + y
	return
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *fooCollector) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.fooMetric
	ch <- collector.barMetric

}

func generateValue(channel chan uint64) {
	// var channel chan int = make(chan int)
	var value uint64 = 1
	for {
		value += 1
		channel <- value
		fmt.Println("put value:", value)
	}
}

func getValueLoop(channel chan uint64) {
	for {
		time.Sleep(2 * time.Second)
		<-channel
	}
}

func getValue(channel chan uint64) (value uint64) {
	return <-channel
}

//Collect implements required collect function for all promehteus collectors
func (collector *fooCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue float64
	if 1 == 1 {
		metricValue = float64(getValue(collector.channel))
	}

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	ch <- prometheus.MustNewConstMetric(collector.fooMetric, prometheus.CounterValue, metricValue)
	ch <- prometheus.MustNewConstMetric(collector.barMetric, prometheus.CounterValue, metricValue)

}
