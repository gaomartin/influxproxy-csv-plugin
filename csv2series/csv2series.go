package csv2series

import (
	"bytes"
	"encoding/csv"
	"sort"
	// "errors"
	// 	"flag"
	// "fmt"
	// 	"github.com/influxdb/influxdb-go"
	// 	"io/ioutil"
	// 	"log"
	// 	"strconv"
	// "strings"
	// 	"path/filepath"
)

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

func ReadTable(data []byte, separator string) ([][]string, error) {
	tableReader := csv.NewReader(bytes.NewReader(data))
	table, err := tableReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return table, err
}

func BuildTree(data [][]string, header []string, hirarchy []string) *Node {
	tree := &Node{
		Name: "/",
	}
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

	for _, row := range data {
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
	return tree
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

// func GetAsSeries(data []byte, dest string, headerStr string, timeField string, thirdDimField string) ([]*influxdb.Series, error) {
// 	var series []*influxdb.Series
// 	for name, points := range m {
// 		out := &influxdb.Series{
// 			Name:    name,
// 			Columns: []string{"time", "value"},
// 			Points:  points,
// 		}

// 		series = append(series, out)
// 	}
// 	return series, nil

// }
