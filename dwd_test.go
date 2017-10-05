package main

import "testing"

func TestHeaderParsing(t *testing.T) {
	column_names, units, descriptions, cols := parse_file("example_headers")
	if column_names[0] != "forecast" {
		t.Errorf("Expected first entry to be forecast, but it was %s instead.", column_names[0])
	}

	if units[0] != "today 15 UTC" {
		t.Errorf("Expected first entry to be today 15 UTC, but it was %s instead.", units[0])
	}

	if descriptions[0] != "date" {
		t.Errorf("Expected first entry to be date, but it was %s instead.", descriptions[0])
	}

	if cols["parameter"][0] != "18:00" {
		t.Errorf("Expected first entry to be 18:00, but it was %s instead.", cols["parameter"][0])
	}
	if cols["parameter"][1] != "21:00" {
		t.Errorf("Expected first entry to be 21:00, but it was %s instead.", cols["parameter"][1])
	}
}
