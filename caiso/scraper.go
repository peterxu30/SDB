package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
	solarProd int
	windProd int
	date time.Time
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
	resp, error := http.Get(src.URL)
	defer resp.Body.Close()
	if error != nil {
		return EnergyProductionData{}, error
	}

	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		return EnergyProductionData{}, error
	}

	stringBody := string(body)

	// a little rough around the edges
	rest := strings.Split(stringBody, "<span class=\"to_readings\" id=\"currentsolar\">")[1]
	currentSolar := strings.Split(rest, "</span><br />")[0]
	currentWind := strings.Split(strings.Split(rest, "</span><br />")[1], 
		"<span class=\"to_callout1\">Current Wind:</span> <span class=\"to_readings\" id=\"currentwind\">")[1]
	
	currentSolarInt, error := strconv.Atoi(strings.Split(currentSolar, " MW")[0])
	if error != nil {
		fmt.Println(error)
	}

	currentWindInt, error := strconv.Atoi(strings.Split(currentWind, " MW")[0])
	if error != nil {
		fmt.Println(error)
	}
	return EnergyProductionData{ solarProd : currentSolarInt, windProd : currentWindInt, date : time.Now() }, error
}

func main() {
	source := NewCaisoEnergySource(10000000000)
	data := source.Start()
	for point := range data {
		fmt.Println(point)
	}
}
