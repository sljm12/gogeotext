package gogeotext

import (
	"fmt"

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
