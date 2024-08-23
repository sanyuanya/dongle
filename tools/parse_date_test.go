package tools

import "testing"

func TestParseDate(t *testing.T) {
	dateStr := "2021-07/01"
	date, err := ParseDate(dateStr)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	t.Logf("Date: %v", date.Unix()) 
}
