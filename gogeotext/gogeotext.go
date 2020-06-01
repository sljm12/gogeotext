package gogeotext

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gopkg.in/jdkato/prose.v2"
)

/*
NERExtractor interface for different types of NERExtractror
*/
type NERExtractor interface {
	Extract(string) []string
}

/*
Prose - struct representing Prose NLP
*/
type Prose struct {
}

/*
Extract belonging to Prose
*/
func (Prose) Extract(s string) []string {
	doc, _ := prose.NewDocument(s)
	results := []string{}
	for _, ent := range doc.Entities() {
		fmt.Println(ent.Text, ent.Label)
		results = append(results, ent.Text)
	}
	return results
}

/*
Location - Represents a location
*/
type Location struct {
	lat  float64
	lon  float64
	name string
}

/*
ReadCsv - Reads a CSV file and returns Location array
*/
func ReadCsv(filename string, latLoc int, lonLoc int, nameLoc int) (map[string][]Location, error) {
	reader, error := os.Open(filename)
	results := make(map[string][]Location)

	if error == nil {
		csvReader := csv.NewReader(reader)
		csvReader.Comma = '\t'
		csvReader.LazyQuotes = true
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			fmt.Println(record[4])
			lat, _ := strconv.ParseFloat(record[latLoc], 64)
			lon, _ := strconv.ParseFloat(record[lonLoc], 64)
			name := strings.ToLower(record[nameLoc])
			location := Location{lat: lat, lon: lon, name: name}

			value := results[name]

			if value == nil {
				results[name] = []Location{location}
			} else {
				arr := results[name]
				results[name] = append(arr, location)
			}
		}

		return results, nil
	}

	return results, error

}

/*
GeoTextLocator - Loads the data into this struct for processing
*/
type GeoTextLocator struct {
	extractor  NERExtractor
	countryMap map[string]Location
}

/*
ExtractGeoLocation - Extracts geolocation from string
*/
func (g GeoTextLocator) ExtractGeoLocation(text string) []string {
	results := g.extractor.Extract(text)
	return results
}
