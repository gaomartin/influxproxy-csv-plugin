// package main

// import (
// 	"bytes"
// 	"encoding/csv"
// 	// "encoding/json"
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"github.com/influxdb/influxdb-go"
// 	"io/ioutil"
// 	"log"
// 	"strconv"
// 	"strings"
// 	"path/filepath"
// )

// const (
// 	defaultInfluxdbAddress     = "localhost"
// 	defaultInfluxdbPort        = "8086"
// 	defaultInfluxdbDatabase    = ""
// 	defaultInfluxdbPrecision   = "s"
// 	defaultInfluxdbUsername    = "root"
// 	defaultInfluxdbPassword    = "root"
// 	defaultInfluxdbPrefix      = ""
// 	defaultHeaderString        = ""
// 	defaultTimeStampField      = "timestamp"
// 	defaultThirdDimensionField = ""
// )

// var influxdbAddress string
// var influxdbPort string
// var influxdbDatabase string
// var influxdbPrecision string
// var influxdbUsername string
// var influxdbPassword string
// var influxdbPrefix string
// var headerString string
// var fileNames []string
// var timeStampField string
// var thirdDimensionField string

// func init() {
// 	flag.StringVar(&influxdbAddress, "addr", defaultInfluxdbAddress, "Hostname or IP of InfluxDB")
// 	flag.StringVar(&influxdbPort, "port", defaultInfluxdbPort, "Port of InfluxDB")
// 	flag.StringVar(&influxdbDatabase, "db", defaultInfluxdbDatabase, "Database name of InfluxDB")
// 	flag.StringVar(&influxdbPrecision, "prec", defaultInfluxdbPrecision, "Precision of InfluxDB")
// 	flag.StringVar(&influxdbUsername, "user", defaultInfluxdbUsername, "Username for InfluxDB")
// 	flag.StringVar(&influxdbPassword, "pass", defaultInfluxdbPassword, "Password for InfluxDB")
// 	flag.StringVar(&influxdbPrefix, "pref", defaultInfluxdbPrefix, "Prefix for InfluxDB time series")
// 	flag.StringVar(&headerString, "head", defaultHeaderString, "header of csv")
// 	flag.StringVar(&timeStampField, "time", defaultTimeStampField, "timestamp field")
// 	flag.StringVar(&thirdDimensionField, "third", defaultThirdDimensionField, "third dimension field")
// 	flag.Parse()
// 	fileNames = flag.Args()
// }

// func logFatal(err error, info string) {
// 	errString := fmt.Sprintf("%s: %s", info, err)
// 	log.Fatal(errString)
// }

// func readFiles(fnames []string) (map[string][]byte, error) {
// 	m := make(map[string][]byte)
// 	if len(fnames) == 0 {
// 		return nil, errors.New("No Files given")
// 	}
// 	for _, filename := range fnames {
// 		file, err := ioutil.ReadFile(filename)
// 		if err != nil {
// 			return nil, err
// 		}
// 		name := filepath.Base(filename)
// 		m[name] = file
// 	}

// 	return m, nil
// }

// func getAsSeries(data []byte, dest string, headerStr string, timeField string, thirdDimField string) ([]*influxdb.Series, error) {
// 	//prepare structure
// 	var series []*influxdb.Series
// 	m := make(map[string][][]interface{})

// 	// read table
// 	tableReader := csv.NewReader(bytes.NewReader(data))
// 	table, err := tableReader.ReadAll()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(table) < 1 {
// 		return nil, errors.New("no data in table")
// 	}

// 	// read header
// 	headerReader := csv.NewReader(strings.NewReader(headerStr))
// 	h, err := headerReader.ReadAll()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(h) != 1 {
// 		return nil, errors.New("malformed header")
// 	}
// 	header := h[0]

// 	// make sure that the header has the same length as the table itself
// 	if len(table[0]) != len(header) {
// 		return nil, errors.New("Lenght of header and content does not match")
// 	}

// 	for _, row := range table {
// 		ts := ""
// 		valueExt := ""
// 		values := make(map[string]string)

// 		for i, field := range row {
// 			switch {
// 			case header[i] == timeField:
// 				ts = field
// 			case header[i] == thirdDimField:
// 				valueExt = field
// 			default:
// 				values[header[i]] = field
// 			}
// 		}

// 		for index, value := range values {
// 			name := ""
// 			if valueExt != "" {
// 				name = fmt.Sprintf("%s.%s.%s", dest, valueExt, index)
// 			} else {
// 				name = fmt.Sprintf("%s.%s", dest, index)
// 			}

// 			name = strings.Replace(name, "%", "perc-", -1)
// 			name = strings.Replace(name, "/", "-per-", -1)

// 			var point []interface{}
// 			t, _ := strconv.ParseFloat(ts, 64)
// 			point = append(point, t)
// 			v, _ := strconv.ParseFloat(value, 64)
// 			point = append(point, v)

// 			m[name] = append(m[name], point)

// 		}
// 	}

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

// func main() {

// 	// build client
// 	influx, _ := influxdb.NewClient(&influxdb.ClientConfig{
// 		Username: influxdbUsername,
// 		Password: influxdbPassword,
// 		Database: influxdbDatabase,
// 		Host:     influxdbAddress + ":" + influxdbPort,
// 	})

// 	// read all given files
// 	files, err := readFiles(fileNames)
// 	if err != nil {
// 		logFatal(err, "Input file(s) could not be read")
// 	}

// 	// convert to influxdb.Series
// 	var series [][]*influxdb.Series
// 	for name, file := range files {
// 		serie, err := getAsSeries(file, name, headerString, timeStampField, thirdDimensionField)
// 		if err != nil {
// 			logFatal(err, "could not convert data to series")
// 		}
// 		series = append(series, serie)
// 	}

// 	for _, serie := range series {
// 	// 	b, err := json.Marshal(serie)
// 	// 	if err != nil {
// 	// 		logFatal(err, "could not convert to json")
// 	// 	}
// 	// 	fmt.Println(string(b))
// 		if err := influx.WriteSeries(serie); err != nil {
// 			logFatal(err, "could not send series")
// 		}
// 	}

// 	log.Print("Sended to InfluxDB")

// }
