package main

import (
	"fmt"
	"strings"

	"github.com/influxproxy/influxproxy-csv-plugin/csv2series"
	"github.com/influxproxy/influxproxy/plugin"
)

type Functions struct{}

func (f Functions) Describe() plugin.Description {
	d := plugin.Description{
		Description: "This plugin takes CSV files and pushes them to the given influxdb.",
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
				Description: "CSV separator character. If multiple characters are provided, only the first character is considered.",
				Optional:    true,
				Default:     ",",
			},
			{
				Name:        "header",
				Description: "Header of the CSV table, colums separated with the same character as provided in 'separator' field. If no header is provided, the first line of the CSV data is considered to be the header.",
				Optional:    true,
				Default:     "",
			},
			{
				Name:        "hierarchy",
				Description: "Name of the fields that imply nesting, separated by the same character as provided in 'separator' field, ordered from top down.",
				Optional:    true,
				Default:     "",
			},
			{
				Name:        "timestamp",
				Description: "Name of the field that contains the time stamp in epoch time (ms).",
				Optional:    false,
				Default:     "",
			},
		},
	}
	return d
}

func (f Functions) Run(in plugin.Request) plugin.Response {
	prefix := in.Query.Get("prefix")
	separator := in.Query.Get("separator")
	if separator == "" {
		separator = ","
	}
	timestamp := in.Query.Get("timestamp")
	header := strings.Split(in.Query.Get("header"), separator)
	hierarchy := strings.Split(in.Query.Get("hierarchy"), separator)

	conv, err := csv2series.NewConverter(in.Body, separator, header, hierarchy)
	if err != nil {
		return plugin.Response{
			Series: nil,
			Error:  err.Error(),
		}
	}

	series, err := conv.GetAsSeries(prefix, timestamp)

	if err != nil {
		return plugin.Response{
			Series: nil,
			Error:  err.Error(),
		}
	}

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
