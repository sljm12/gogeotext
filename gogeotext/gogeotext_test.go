package gogeotext

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
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

func TestReadData(t *testing.T) {
	reader, error := os.Open("./data/cities500.txt")
	var results []Location
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
			lat, _ := strconv.ParseFloat(record[4], 64)
			lon, _ := strconv.ParseFloat(record[5], 64)
			results = append(results, Location{lat: lat, lon: lon, name: record[2]})
		}
	} else {
		t.Error("Failed to open file")
	}
	if len(results) != 192573 {
		fmt.Println(len(results))
		t.Error("results count wrong")
	}
}

func TestReadCSV_CitiesData(t *testing.T) {
	locations, error := ReadCsv("./data/cities500.txt", 4, 5, 2)

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
