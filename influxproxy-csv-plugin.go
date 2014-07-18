package main

import (
	"fmt"

	"github.com/influxdb/influxdb-go"
	"github.com/influxproxy/influxproxy/plugin"
)

type Functions struct{}

func (f Functions) Describe() plugin.Description {
	d := plugin.Description{
		Description: "This plugin takes CSV files and pushes them to the given influxdb",
		Author:      "github.com/sontags",
		Version:     "0.1.0",
		Arguments: []plugin.Argument{
			{
				Name:        "prefix",
				Description: "Prefix of the series, will be separated with a '.' if given.",
				Optional:    true,
			},
			{
				Name:        "separator",
				Description: "CSV separator character, default is ','.",
				Optional:    true,
			},
			{
				Name:        "header",
				Description: "Header of the CSV table, colums separated with the same character as provided in 'separator' field.",
				Optional:    false,
			},
			{
				Name:        "nesting",
				Description: "Name of the fields that imply nesting, separated by character ',', ordered from top down.",
				Optional:    true,
			},
			{
				Name:        "timestamp",
				Description: "Name of the field that contains the time stamp in epoch time (s)",
				Optional:    false,
			},
		},

	}
	return d
}


	defaultHeaderString        = ""
	defaultTimeStampField      = "timestamp"
	defaultThirdDimensionField = ""

func (f Functions) Run(in plugin.Request) plugin.Response {
	// TODO: implement...
	var series []*influxdb.Series
	return plugin.Response{
		Series: series,
		Error:  "",
	}
}

func main() {
	f := Functions{}
	p, err := plugin.NewPlugin()
	if err != nil {
		fmt.Println(err)
	} else {
		p.Run(f)
	}
}
