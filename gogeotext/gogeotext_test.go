package gogeotext

import (
	"fmt"
	"testing"
)

func inStringArray(value string, arr []string) bool {
	for _, v := range arr {
		if value == v {
			return true
		}
	}

	return false
}

func TestProseExtract(t *testing.T) {
	var p Prose

	var gtl GeoTextLocator
	gtl.extractor = p
	results := gtl.ExtractGeoLocation("Singaporeans and Singapore are one")
	if len(results) != 1 {
		t.Error("Results are wrong")
	}

	if inStringArray("Singapore", results) == false {
		t.Error("Singapre not in answer")
	}
}

func TestReadCSV_CitiesData(t *testing.T) {
	locations, error := ReadCsv("./data/cities500.txt", '\t', 4, 5, 2, 8)

	if error == nil {
		if len(locations) != 161489 {
			fmt.Println(len(locations))
			t.Error("results count wrong")
		}
	} else {
		t.Error("Error reding count wrong")
	}

	r := locations["wellington"]
	if len(r) > 0 {
		fmt.Println(r)
	} else {
		t.Error("Map Wellington wrong")
	}
}

func TestReadCSV_CountryData(t *testing.T) {

	locations, error := ReadCsv("./data/alternateName.csv", ',', 3, 4, 5, 1)

	if error == nil {
		r := locations["singapore"]
		if r == nil {
			t.Error("Singapore not found")
		} else {
			fmt.Println(r)
		}

	} else {
		t.Error("Error in reading country file")
	}
}

func TestMatchCountry(t *testing.T) {
	var p Prose
	geoText := NewGeoTextLocator(p, "./data/alternateName.csv", "", "")
	location, err := geoText.MatchCountry("singapore")

	if err == nil {
		fmt.Println(location)
	} else {
		t.Error("Cannot find country")
	}
}

func TestMatchReadDefaultCity(t *testing.T) {

	geoText, err := ReadCSVDefaultCity("./data/default_city.csv")

	if err == nil {
		fmt.Println(geoText["wellington"])
	} else {
		t.Error("Cannot find city")
	}
}

func TestMatchDefaultCity(t *testing.T) {
	var p Prose
	geoText := NewGeoTextLocator(p, "", "", "./data/default_city.csv")
	location, present := geoText.MatchDefaultCity("wellington")

	if present == true {
		fmt.Println(location)
	} else {
		t.Error("Cannot find city")
	}
}

func TestMatchCityWithDefault(t *testing.T) {
	var p Prose
	geoText := NewGeoTextLocator(p, "", "./data/cities500.txt", "./data/default_city.csv")
	location, _ := geoText.MatchCity("wellington")

	if location.countryCode == "NZ" {
		fmt.Println(location)
	} else {
		t.Error("Cannot find city")
	}

	location, present := geoText.MatchCity("San Diego")
	if present == false {
		t.Error("Cannot find")
	}

	if location.countryCode == "US" {
		fmt.Println(location)
	} else {
		t.Error("Cannot find city")
	}
}
