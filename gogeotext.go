/*
Package gogeotext extracts locations in the text and returns the lat, lon for that location.

It is designed to use any NER engine, the default is using Prose NLP.

It uses data from geonames for that base data.
It tries to be "smart"

if the text contains just "San Diego is a great place" it will return San Diego in the United States
If the text contains "San Diego, Mexico is a great place" it will return San Diego in Mexico*/
package gogeotext

import (
	"encoding/csv"
	"errors"
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
		results = append(results, ent.Text)
	}
	return results
}

/*
Location - Represents a location
*/
type Location struct {
	Lat         float64
	Lon         float64
	Name        string
	CountryCode string
	Geonameid   string
}

/*
DefaultCity stuct for storing default city
*/
type DefaultCity struct {
	Name    string
	Country string
}

/*
LocationResult struct for Extract Results
*/
type LocationResult struct {
	Class       string
	Name        string
	CountryCode string
	Lat         float64
	Lon         float64
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
			location := Location{Lat: lat, Lon: lon, Name: name, CountryCode: countryCode}

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
			defaultCities[strings.ToLower(name)] = DefaultCity{Name: name, Country: country}

		}
	}
	return defaultCities, err
}

/*
GeoTextLocator - Loads the data into this struct for processing
*/
type GeoTextLocator struct {
	Extractor   NERExtractor
	CountryMap  map[string][]Location
	CitiesMap   map[string][]Location
	DefaultCity map[string]DefaultCity
}

/*
GeoTextLocatorResults - The results from GeoTextLocator
*/
type GeoTextLocatorResults struct {
	Countries []Location
	Cities    []Location
}

/*
ExtractGeoLocation - Extracts geolocation from string
*/
func (g GeoTextLocator) ExtractGeoLocation(text string) GeoTextLocatorResults {
	tokens := g.Extractor.Extract(text)
	var results GeoTextLocatorResults
	usedTokens := make([]bool, len(tokens))

	//Find countries
	for i, r := range tokens {
		lower := strings.ToLower(r)
		country := g.CountryMap[lower]
		if country != nil {
			results.Countries = append(results.Countries, country[0])
			usedTokens[i] = true
		} else {
			usedTokens[i] = false
		}
	}

	//Find cities
	for i, r := range tokens {
		if usedTokens[i] == false {
			city, present := g.MatchCity(r, results.Countries)
			if present == true {
				results.Cities = append(results.Cities, city)
			}
		}
	}

	return results
}

/*
MatchCountry on the token
*/
func (g GeoTextLocator) MatchCountry(token string) (Location, error) {
	value := g.CountryMap[strings.ToLower(token)]

	if value != nil {
		return value[0], nil
	}
	return Location{}, errors.New("Can't find country " + token)

}

/*
MatchCity Matching city based on token
*/
func (g GeoTextLocator) MatchCity(token string, countryResults []Location) (Location, bool) {
	lowerToken := strings.ToLower(token)

	//Find a city using the countries that was detected earlier
	for _, v := range countryResults {
		location, present := g.FindCity(lowerToken, v.CountryCode)
		if present == true {
			return location, true
		}
	}

	//Match based on Default City
	defaultCity, present := g.MatchDefaultCity(lowerToken)

	if present {
		return g.FindCity(defaultCity.Name, defaultCity.Country)
	}

	//Match based on City
	cities := g.CitiesMap[lowerToken]
	if cities != nil && len(cities) > 0 {
		firstCity := cities[0]
		return firstCity, true
	}

	return Location{}, false
}

/*
MatchDefaultCity Match a default city
*/
func (g GeoTextLocator) MatchDefaultCity(token string) (DefaultCity, bool) {
	value, present := g.DefaultCity[token]
	return value, present
}

/*
FindCity
Find City given the name and country
*/
func (g GeoTextLocator) FindCity(name string, country string) (Location, bool) {
	cities := g.CitiesMap[name]
	for _, v := range cities {
		if v.CountryCode == country {
			return v, true
		}
	}

	return Location{}, false
}

/*
MatchCityCoutry match a city given a list of country
*/
func (g GeoTextLocator) MatchCityCoutry(token string, countries []string) []Location {
	cities, present := g.CitiesMap[token]
	result := make([]Location, 0)

	if present == true {
		for _, value := range cities {
			for _, c := range countries {
				if value.CountryCode == c {
					result = append(result, value)
				}
			}
		}
	}

	return result
}

/*
NewGeoTextLocator - Create new GeoTextLocator
*/
func NewGeoTextLocator(e NERExtractor, countryFile string, citiesFiles string, defaultCity string) GeoTextLocator {
	var a GeoTextLocator
	a.Extractor = e
	var err error
	a.CountryMap, err = ReadCsv(countryFile, ',', 3, 4, 5, 1)
	if err != nil {
		panic(err)
	}
	a.CitiesMap, err = ReadCsv(citiesFiles, '\t', 4, 5, 2, 8)
	if err != nil {
		panic(err)
	}
	a.DefaultCity, err = ReadCSVDefaultCity(defaultCity)
	if err != nil {
		panic(err)
	}
	return a
}

/*
CreateDefaultGeoTextLocator - Create the default GeoTextLocator using the Prose NER
*/
func CreateDefaultGeoTextLocator(countryFile string, citiesFiles string, defaultCity string) GeoTextLocator {
	var p Prose
	return NewGeoTextLocator(p, countryFile, citiesFiles, defaultCity)
}
