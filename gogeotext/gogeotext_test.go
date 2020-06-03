package gogeotext

import (
	"fmt"
	"reflect"
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

func TestProse(t *testing.T) {
	var p Prose

	s := p.Extract("San Diego is a great place to live.")
	if inStringArray("San Diego", s) != true {
		t.Error("San Diego not found")
	}

	s = p.Extract("San Diego, Mexico is a great place to live.")
	if inStringArray("San Diego", s) != true {
		t.Error("San Diego not found")
	}

	if inStringArray("Mexico", s) != true {
		t.Error("Mexico not found")
	}
}

func TestGTLExtract(t *testing.T) {
	var p Prose

	var gtl GeoTextLocator
	gtl = NewGeoTextLocator(p, "./data/alternateName.csv", "./data/cities500.txt", "./data/default_city.csv")

	//Test Singapore
	results := gtl.ExtractGeoLocation("Singaporeans and Singapore are one")
	if len(results.Countries) != 1 {
		t.Error("Results are wrong")
	}

	if results.Countries[0].countryCode != "SG" {
		t.Error("Country are wrong")
	}

	//Test San Diego and Mexico extraction
	results = gtl.ExtractGeoLocation("San Diego, Mexico are great places to live")
	if len(results.Countries) != 1 {
		t.Error("Results are wrong")
	}

	if results.Countries[0].countryCode != "MX" {
		t.Error("Country are wrong")
	}

	if len(results.Cities) != 1 {
		t.Error("Results are wrong")
	}

	if results.Cities[0].name != "san diego" {
		t.Error("Results are wrong")
	}

	if results.Cities[0].countryCode != "MX" {
		t.Error("Results are wrong")
	}

	//Test wrong country and city

	results = gtl.ExtractGeoLocation("San Diego, Singapore are great places to live")
	if len(results.Countries) != 1 {
		t.Error("Results are wrong")
	}

	if results.Countries[0].countryCode != "SG" {
		t.Error("Country are wrong")
	}

	if len(results.Cities) != 1 {
		t.Error("Results are wrong")
	}

	if results.Cities[0].name != "san diego" {
		t.Error("Results are wrong")
	}

	if results.Cities[0].countryCode != "US" {
		t.Error("Results are wrong")
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
	geoText := NewGeoTextLocator(p, "./data/alternateName.csv", "./data/cities500.txt", "./data/default_city.csv")
	location, err := geoText.MatchCountry("singapore")

	if err == nil {
		fmt.Println(location)
	} else {
		t.Error("Cannot find country")
	}

	location, err = geoText.MatchCountry("rubbish country")

	if err == nil {
		t.Error("Not suppose to find country")
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
	geoText := NewGeoTextLocator(p, "./data/alternateName.csv", "./data/cities500.txt", "./data/default_city.csv")
	location, present := geoText.MatchDefaultCity("wellington")

	if present == true {
		fmt.Println(location)
	} else {
		t.Error("Cannot find city")
	}
}

func TestMatchCityWithDefault(t *testing.T) {
	var p Prose
	geoText := NewGeoTextLocator(p, "./data/alternateName.csv", "./data/cities500.txt", "./data/default_city.csv")
	location, _ := geoText.MatchCity("wellington", []Location{})

	//Check for default
	if location.countryCode == "NZ" {
		fmt.Println(location)
	} else {
		t.Error("Cannot find city")
	}

	location, present := geoText.MatchCity("San Diego", []Location{})
	if present == false {
		t.Error("Cannot find")
	}

	//Check for default
	if location.countryCode == "US" {
		fmt.Println(location)
	} else {
		t.Error("Cannot find city")
	}

	//Check if city not present
	_, present = geoText.MatchCity("RubbishCity", []Location{})
	if present != false {
		t.Error("Found Rubbish City")
	}

	//Test matching of only one city and in lower case
	_, present = geoText.MatchCity("soldeu", []Location{})

	if present != true {
		t.Error("Suppose to find the city")
	}
}

func TestMatchCityCountry(t *testing.T) {
	var p Prose
	geoText := NewGeoTextLocator(p, "./data/alternateName.csv", "./data/cities500.txt", "./data/default_city.csv")

	//Wellignton if given India
	results := geoText.MatchCityCoutry("wellington", []string{"IN"})
	if len(results) != 1 && results[0].countryCode == "IN" {
		t.Error("Should be 1")
	}

	//
	results = geoText.MatchCityCoutry("wellington", []string{"SG"})
	if len(results) != 0 {
		t.Error("Should be 0 instead")
	}
}

func TestFindCityCountry(t *testing.T) {
	var p Prose
	geoText := NewGeoTextLocator(p, "./data/alternateName.csv", "./data/cities500.txt", "./data/default_city.csv")

	//Find a correct city
	city, present := geoText.FindCity("wellington", "NZ")

	if present == false {
		t.Error("City not found")
	} else {
		if city.name != "wellington" && city.countryCode != "NZ" {
			t.Error("Wrong City")
		}
	}
	fmt.Println(city)

	//Test if the city cannot be found
	city, present = geoText.FindCity("rubbish", "SG")

	if present == true {
		t.Error("City not suppose to be found")
	}
}

func TestPkgPath(t *testing.T) {
	var p Prose
	fmt.Println(reflect.TypeOf(p).PkgPath())
}
