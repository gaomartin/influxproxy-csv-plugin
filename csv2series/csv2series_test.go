package csv2series

import (
	"testing"

	influxdb "github.com/influxdb/influxdb/client"
)

var tests = []struct {
	Name      string
	Input     []byte
	Output    []*influxdb.Series
	Separator string
	Header    []string
	Hierarchy []string
	Prefix    string
	Timestamp string
}{
	{
		Name:  "external header, no hierarchy, no prefix, no text fiels",
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
		Separator: ",",
		Header:    []string{"time", "PercentCpu", "PercentMem"},
		Hierarchy: []string{},
		Prefix:    "",
		Timestamp: "time",
	},
}

func TestSeriesConversion(t *testing.T) {
	for _, test := range tests {
		c, e := NewConverter(test.Input, test.Separator, test.Header, test.Hierarchy)
		if e != nil {
			t.Error("Error happend while creating a new converter in test", test.Name)
		}
		_, e = c.GetAsSeries(test.Prefix, test.Timestamp)
		// s, e := c.GetAsSeries(test.Prefix, test.Timestamp)
		if e != nil {
			t.Error("Error happend while getting series in test", test.Name)
		}
		// if s != test.Output {
		// 	t.Error("Wrong outcome in test", test.Name)
		// }
	}
}
