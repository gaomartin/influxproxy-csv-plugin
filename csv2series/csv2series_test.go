package csv2series

import (
	"testing"

	influxdb "github.com/influxdb/influxdb/client"
)

var tests = []struct {
	Name        string
	Input       []byte
	Output      []*influxdb.Series
	Separator   string
	Header      []string
	Hierarchy   []string
	Timepattern string
	Prefix      string
	Timestamp   string
}{
	{
		Name:  "external header, no hierarchy, no prefix, no text fiels, no timepattern",
		Input: []byte("507744865000,5.7,21\n507744875000,5.9,21.4\n507744885000,3.6,18.4\n"),
		Output: []*influxdb.Series{
			{
				Name:    "PercentCpu",
				Columns: []string{"time", "value"},
				Points: [][]interface{}{
					{507744865000, 5.7},
					{507744875000, 5.9},
					{507744885000, 3.6},
				},
			},
			{
				Name:    "PercentMem",
				Columns: []string{"time", "value"},
				Points: [][]interface{}{
					{507744865000, 21},
					{507744875000, 21.4},
					{507744885000, 18.4},
				},
			},
		},
		Separator:   ",",
		Header:      []string{"time", "PercentCpu", "PercentMem"},
		Hierarchy:   []string{},
		Timepattern: "",
		Prefix:      "",
		Timestamp:   "time",
	},
	{
		Name:  "external header, no hierarchy, no prefix, no text fiels, with timepattern",
		Input: []byte("1986-Feb-02/16:14:25,5.7,21\n1986-Feb-02/16:14:35,5.9,21.4\n1986-Feb-02/16:14:45,3.6,18.4\n"),
		Output: []*influxdb.Series{
			{
				Name:    "PercentCpu",
				Columns: []string{"time", "value"},
				Points: [][]interface{}{
					{507744865000, 5.7},
					{507744875000, 5.9},
					{507744885000, 3.6},
				},
			},
			{
				Name:    "PercentMem",
				Columns: []string{"time", "value"},
				Points: [][]interface{}{
					{507744865000, 21},
					{507744875000, 21.4},
					{507744885000, 18.4},
				},
			},
		},
		Separator:   ",",
		Header:      []string{"time", "PercentCpu", "PercentMem"},
		Hierarchy:   []string{},
		Timepattern: "2006-Jan-02/15:04:05",
		Prefix:      "",
		Timestamp:   "time",
	},
}

func TestSeriesConversion(t *testing.T) {
	for _, test := range tests {
		c, e := NewConverter(test.Input, test.Separator, test.Header, test.Hierarchy)
		if e != nil {
			t.Error("Error happend while creating a new converter in test", test.Name, e)
		}
		_, e = c.GetAsSeries(test.Prefix, test.Timestamp, test.Timepattern)
		// s, e := c.GetAsSeries(test.Prefix, test.Timestamp)
		if e != nil {
			t.Error("Error happend while getting series in test", test.Name, e)
		}
		// if s != test.Output {
		// 	t.Error("Wrong outcome in test", test.Name)
		// }
	}
}
