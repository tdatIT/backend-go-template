package datetime

import (
	"testing"
	"time"
)

func TestParseDatetimeFromString(t *testing.T) {
	timeStr := "2024-09-26 09:11:08"
	expectedTime := time.Date(2024, time.September, 26, 9, 11, 8, 0, time.UTC)
	parsedTime, err := ParseDatetimeFromString(timeStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("expected %v, got %v", expectedTime, parsedTime)
	}
}

func TestParseDateFromString(t *testing.T) {
	dateStr := "2024-09-26"
	expectedTime := time.Date(2024, time.September, 26, 0, 0, 0, 0, time.UTC)
	parsedTime, err := ParseDateFromString(dateStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("expected %v, got %v", expectedTime, parsedTime)
	}
}

func TestParseTimeFromString(t *testing.T) {
	timeStr := "09:11:08"
	expectedTime := time.Date(0, 1, 1, 9, 11, 8, 0, time.UTC)
	parsedTime, err := ParseTimeFromString(timeStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if parsedTime.Hour() != expectedTime.Hour() || parsedTime.Minute() != expectedTime.Minute() || parsedTime.Second() != expectedTime.Second() {
		t.Errorf("expected %v, got %v", expectedTime, parsedTime)
	}
}

func TestFormatDatetime(t *testing.T) {
	timeVal := time.Date(2024, time.September, 26, 9, 11, 8, 0, time.UTC)
	expectedStr := "2024-09-26 09:11:08"
	formattedStr := FormatDatetime(timeVal)
	if formattedStr != expectedStr {
		t.Errorf("expected %v, got %v", expectedStr, formattedStr)
	}
}

func TestFormatDate(t *testing.T) {
	timeVal := time.Date(2024, time.September, 26, 0, 0, 0, 0, time.UTC)
	expectedStr := "2024-09-26"
	formattedStr := FormatDate(timeVal)
	if formattedStr != expectedStr {
		t.Errorf("expected %v, got %v", expectedStr, formattedStr)
	}
}

func TestFormatTime(t *testing.T) {
	timeVal := time.Date(0, 1, 1, 9, 11, 8, 0, time.UTC)
	expectedStr := "09:11:08"
	formattedStr := FormatTime(timeVal)
	if formattedStr != expectedStr {
		t.Errorf("expected %v, got %v", expectedStr, formattedStr)
	}
}

func TestParseTimeWithUTC7(t *testing.T) {
	// Test input
	timeStr := "2024-09-26 09:11:08"

	// Expected time in Asia/Bangkok timezone (UTC+7)
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		t.Fatalf("failed to load Asia/Bangkok location: %v", err)
	}
	expectedTime := time.Date(2024, time.September, 26, 9, 11, 8, 0, location)

	// Parse the time string
	parsedTime, err := ParseTimeWithUTC7(timeStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check equality
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("expected %v, got %v", expectedTime, parsedTime)
	}

	// Also check that the timezone is correct
	if parsedTime.Location().String() != location.String() {
		t.Errorf("expected location %v, got %v", location, parsedTime.Location())
	}

	// Test with invalid format
	_, err = ParseTimeWithUTC7("2024/09/26 09:11:08")
	if err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}
