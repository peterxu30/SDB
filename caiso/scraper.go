package main

import (
	"fmt"
	// "io/ioutil"
	"github.com/PuerkitoBio/goquery"
	// "net/http"
	"strconv"
	"strings"	
	"time"
)

type CaisoEnergySource struct {
	URL string
	rate time.Duration
	data chan EnergyProductionData
}

type EnergyProductionData struct {
	SolarProd int
	WindProd int
	Date time.Time
}

func NewCaisoEnergySource(rate int) *CaisoEnergySource {
	return &CaisoEnergySource {
		URL: "http://content.caiso.com/outlook/SP/renewables.html",
		rate: time.Duration(rate),
		data: make(chan EnergyProductionData),
	}
}

func (src *CaisoEnergySource) Start() chan EnergyProductionData {
	go func() {
		if point, err := src.Read(); err == nil {
			src.data <- point
		}
		for _ = range time.Tick(src.rate) {
			if point, err := src.Read(); err == nil {
				src.data <- point
			}
		}
	}()
	return src.data
}

func (src *CaisoEnergySource) Read() (EnergyProductionData, error) {
	doc, err := goquery.NewDocument(src.URL)
	if err != nil {
		return EnergyProductionData{}, err
	}

	var currentSolar string
	var currentWind string

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
	    currentSolar = s.Find("#currentsolar").Text()
	    currentWind = s.Find("#currentwind").Text()
  	})
	
	currentSolarInt, err := strconv.Atoi(strings.Split(currentSolar, " MW")[0])
	if err != nil {
		fmt.Println(err)
	}

	currentWindInt, err := strconv.Atoi(strings.Split(currentWind, " MW")[0])
	if err != nil {
		fmt.Println(err)
	}
	return EnergyProductionData{ SolarProd : currentSolarInt, WindProd : currentWindInt, Date : time.Now() }, err
}

func main() {
	source := NewCaisoEnergySource(10000000000)
	data := source.Start()
	for point := range data {
		fmt.Println(point)
	}
}
