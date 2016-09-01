package main

import (
	"fmt"
	"github.com/immesys/spawnpoint/spawnable"
	"github.com/satori/go.uuid"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"time"
)

//needs lots of work
func main() {
	bw := bw2.ConnectOrExit("")

	params := spawnable.GetParamsOrExit()
	apikey := params.MustString("API_KEY")
	city := params.MustString("city")
	baseuri := params.MustString("svc_base_uri")
	read_rate := params.MustString("read_rate")

	bw.OverrideAutoChainTo(true)
	bw.SetEntityFromEnvironOrExit()
	svc := bw.RegisterService(baseuri, "s.caiso")
	iface := svc.RegisterInterface(city, "i.weather")

	params.MergeMetadata(bw)

	fmt.Println(iface.FullURI())
	fmt.Println(iface.SignalURI("fahrenheit"))

	// generate UUIDs from city + metric name
	temp_f_uuid := uuid.NewV3(NAMESPACE_UUID, city+"fahrenheit").String()
	temp_c_uuid := uuid.NewV3(NAMESPACE_UUID, city+"celsius").String()
	relative_humidity_uuid := uuid.NewV3(NAMESPACE_UUID, city+"relative_humidity").String()

	src := NewCaisoEnergySource(read_rate)
	data := src.Start()
	for point := range data {
		fmt.Println(point)
		temp_f := TimeseriesReading{UUID: temp_f_uuid, Time: time.Now().Unix(), Value: point.F}
		iface.PublishSignal("fahrenheit", temp_f.ToMsgPackBW())

		temp_c := TimeseriesReading{UUID: temp_c_uuid, Time: time.Now().Unix(), Value: point.C}
		iface.PublishSignal("celsius", temp_c.ToMsgPackBW())

		rh := TimeseriesReading{UUID: relative_humidity_uuid, Time: time.Now().Unix(), Value: point.RH}
		iface.PublishSignal("relative_humidity", rh.ToMsgPackBW())
	}
}