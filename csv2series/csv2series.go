package csv2series

import (
	"bytes"
	"encoding/csv"
	"errors"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"

	influxdb "github.com/influxdb/influxdb/client"
)

type Converter struct {
	Table [][]string
	Tree  *Node
}

func NewConverter(data []byte, separator string, header []string, hirarchy []string) (*Converter, error) {
	c := &Converter{}

	sep, _ := utf8.DecodeRune([]byte(separator))

	err := c.ReadTable(data, sep)
	if err != nil {
		return nil, err
	}

	if len(header) == 1 && header[0] == "" && len(c.Table) > 1 {
		header = c.Table[0]
		c.Table = c.Table[1:]
	} else if len(c.Table) > 0 && len(header) != len(c.Table[0]) {
		return nil, errors.New("Header length differs from columns in table.")
	}

	c.BuildTree(header, hirarchy)
	return c, nil
}

func (c *Converter) ReadTable(data []byte, separator rune) error {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = separator
	table, err := reader.ReadAll()
	if err != nil {
		return err
	}
	c.Table = table
	return nil
}

func (c *Converter) BuildTree(header []string, hirarchy []string) {
	tree := &Node{}
	branch := tree

	var hIndex []int
	for _, hName := range hirarchy {
		for cNum, cName := range header {
			if hName == cName {
				hIndex = append(hIndex, cNum)
			}
		}
	}

	dataHeader := remove(header, hIndex...)

	for _, row := range c.Table {
		for _, hField := range hIndex {
			branch = branch.getChild(row[hField])
		}
		row := remove(row, hIndex...)

		values := make(map[string]string)
		for cellNum, cell := range dataHeader {
			values[cell] = row[cellNum]
		}
		branch.Values = append(branch.Values, values)
		branch = tree
	}
	c.Tree = tree
}

func remove(s []string, items ...int) (out []string) {
	out = append(out, s...)
	sort.Sort(sort.Reverse(sort.IntSlice(items)))
	for _, item := range items {
		tmp := append(out[:item], out[item+1:]...)
		out = tmp
	}
	return
}

func (c *Converter) getTime(ts string, pattern string) (timestamp int64, err error) {
	if pattern == "" {
		timestamp, err = strconv.ParseInt(ts, 0, 64)
		if err != nil {
			return 0.0, err
		}
	} else {
		t, err := time.Parse(pattern, ts)
		if err != nil {
			return 0.0, err
		}
		timestamp = t.Unix() * 1000
	}
	return
}

func (c *Converter) GetAsSeries(prefix string, timestamp string, timepattern string) ([]*influxdb.Series, error) {
	if prefix != "" {
		prefix += "."
	}
	series := &Series{}

	for sName, sValues := range c.Tree.Flatten() {
		if sName != "" {
			sName += "."
		}
		for _, values := range sValues {
			time, err := c.getTime(values[timestamp], timepattern)
			if err != nil {
				return nil, err
			}
			for key, value := range values {
				if key != timestamp {
					name := prefix + sName + key
					val, err := strconv.ParseFloat(value, 64)
					if err != nil {
						series.addTimeValue(name, time, value)
					} else {
						series.addTimeValue(name, time, val)
					}
				}
			}
		}
	}
	return *series, nil
}

type Node struct {
	Name     string
	Children []*Node
	Values   []map[string]string
}

func (n *Node) getChild(name string) *Node {
	for _, c := range n.Children {
		if c.Name == name {
			return c
		}
	}

	child := &Node{
		Name: name,
	}
	n.Children = append(n.Children, child)
	return child
}

func (n *Node) Flatten() map[string][]map[string]string {
	if len(n.Children) < 1 {
		return map[string][]map[string]string{
			n.Name: n.Values,
		}
	}
	out := make(map[string][]map[string]string)
	for _, c := range n.Children {
		cOut := c.Flatten()
		for key, val := range cOut {
			prefix := ""
			if n.Name != "" {
				prefix = n.Name + "."
			}
			out[prefix+key] = val
		}
	}
	return out
}

type Series []*influxdb.Series

func (s *Series) addTimeValue(name string, time int64, value interface{}) {
	for _, serie := range *s {
		if serie.GetName() == name {
			serie.Points = append(serie.Points, []interface{}{time, value})
			return
		}
	}
	var points [][]interface{}
	points = append(points, []interface{}{time, value})
	serie := &influxdb.Series{
		Name:    name,
		Columns: []string{"time", "value"},
		Points:  points,
	}
	*s = append(*s, serie)
}
