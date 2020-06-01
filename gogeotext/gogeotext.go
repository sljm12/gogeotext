package gogeotext

import (
	"encoding/csv"
	"errors"
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
	lat         float64
	lon         float64
	name        string
	countryCode string
	geonameid   string
}

/*
DefaultCity stuct for storing default city
*/
type DefaultCity struct {
	name    string
	country string
}

/*
LocationResult struct for Extract Results
*/
type LocationResult struct {
	class       string
	name        string
	countryCode string
}

/*
ReadCsv - Reads a CSV file and returns Location array
*/
func ReadCsv(filename string, delimiter rune, latLoc int, lonLoc int, nameLoc int, countryCodeLoc int) (map[string][]Location, error) {
	reader, error := os.Open(filename)
	results := make(map[string][]Location)

	if error == nil {
		csvReader := csv.NewReader(reader)
		csvReader.Comma = delimiter
		csvReader.LazyQuotes = true
		csvReader.Comment = '#'
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			lat, _ := strconv.ParseFloat(record[latLoc], 64)
			lon, _ := strconv.ParseFloat(record[lonLoc], 64)
			name := strings.ToLower(record[nameLoc])
			countryCode := record[countryCodeLoc]
			location := Location{lat: lat, lon: lon, name: name, countryCode: countryCode}

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
ReadCSVDefaultCity read csv default city
*/
func ReadCSVDefaultCity(filename string) (map[string]DefaultCity, error) {
	defaultCities := make(map[string]DefaultCity)
	file, err := os.Open(filename)
	csvReader := csv.NewReader(file)
	csvReader.Comment = '#'

	if err == nil {
		for {
			record, csverr := csvReader.Read()
			if csverr == io.EOF {
				break
			}

			name := strings.ToLower(record[1])
			country := record[2]
			defaultCities[strings.ToLower(name)] = DefaultCity{name: name, country: country}

		}
	}
	return defaultCities, err
}

/*
GeoTextLocator - Loads the data into this struct for processing
*/
type GeoTextLocator struct {
	extractor   NERExtractor
	countryMap  map[string][]Location
	citiesMap   map[string][]Location
	defaultCity map[string]DefaultCity
}

/*
ExtractGeoLocation - Extracts geolocation from string
*/
func (g GeoTextLocator) ExtractGeoLocation(text string) []string {
	results := g.extractor.Extract(text)
	return results
}

/*
MatchCountry on the token
*/
func (g GeoTextLocator) MatchCountry(token string) (Location, error) {
	value := g.countryMap[strings.ToLower(token)]

	if value != nil {
		return value[0], nil
	}
	return Location{}, errors.New("Can't find country " + token)

}

/*
MatchCity Matching city based on token
*/
func (g GeoTextLocator) MatchCity(token string) (LocationResult, bool) {
	lowerToken := strings.ToLower(token)
	defaultCity, present := g.MatchDefaultCity(lowerToken)

	if present {
		return LocationResult{name: defaultCity.name, countryCode: defaultCity.country, class: "CITY"}, true
	}

	cities, present := g.citiesMap[lowerToken]
	if present {
		firstCity := cities[0]
		return LocationResult{name: firstCity.name, class: "CITY", countryCode: firstCity.countryCode}, true
	}

	return LocationResult{}, false
}

/*
MatchDefaultCity Match a default city
*/
func (g GeoTextLocator) MatchDefaultCity(token string) (DefaultCity, bool) {
	value, present := g.defaultCity[token]
	return value, present
}

/*
NewGeoTextLocator - Create new GeoTextLocator
*/
func NewGeoTextLocator(e NERExtractor, countryFile string, citiesFiles string, defaultCity string) GeoTextLocator {
	var a GeoTextLocator
	a.extractor = e
	a.countryMap, _ = ReadCsv(countryFile, ',', 3, 4, 5, 1)
	a.citiesMap, _ = ReadCsv(citiesFiles, '\t', 4, 5, 2, 8)
	a.defaultCity, _ = ReadCSVDefaultCity(defaultCity)
	return a
}
