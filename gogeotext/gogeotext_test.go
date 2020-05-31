package gogeotext

import (
	"testing"
)

func TestProseExtract(t *testing.T) {
	var p Prose

	var gtl GeoTextLocator
	gtl.extractor = p
	results := gtl.ExtractGeoLocation("Singaporeans and Singapore are one")
	if len(results) != 1 {
		t.Error("Results are wrong")
	}
}
