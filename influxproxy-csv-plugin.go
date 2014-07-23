package main

import (
	"encoding/json"
	"fmt"

	influxdb "github.com/influxdb/influxdb/client"
	"github.com/influxproxy/influxproxy-csv-plugin/csv2series"
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
				Default:     "",
			},
			{
				Name:        "separator",
				Description: "CSV separator character.",
				Optional:    true,
				Default:     ",",
			},
			{
				Name:        "header",
				Description: "Header of the CSV table, colums separated with the same character as provided in 'separator' field.",
				Optional:    false,
				Default:     "",
			},
			{
				Name:        "nesting",
				Description: "Name of the fields that imply nesting, separated by character ',', ordered from top down.",
				Optional:    true,
				Default:     "",
			},
			{
				Name:        "timestamp",
				Description: "Name of the field that contains the time stamp in epoch time (s)",
				Optional:    false,
				Default:     "",
			},
		},
	}
	return d
}

func (f Functions) Run(in plugin.Request) plugin.Response {
	out, _ := csv2series.ReadTable(in.Body, ",")

	header := []string{"Menet", "Anna", "30", "Zurich"}
	hirarchy := []string{}

	tree := csv2series.BuildTree(out, header, hirarchy)

	text, _ := json.Marshal(tree)

	var series []*influxdb.Series
	return plugin.Response{
		Series: series,
		Error:  string(text),
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
