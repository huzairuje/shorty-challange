package tiny_url

import (
	"regexp"
	"testing"
	"time"

	"shorty-challenge/pkg/utils"
)

func resetList() {
	// reset the global slice to avoid test interference
	ListAllTinyUrl = []*Data{}
}

func TestGetAllData(t *testing.T) {
	resetList()
	// Prepopulate global slice with one entry.
	ListAllTinyUrl = []*Data{
		{
			ShortCode:     "abc123",
			Url:           "http://example.com",
			StartDate:     "2020-01-01T00:00:00Z",
			LastSeenDate:  "",
			RedirectCount: 0,
		},
	}

	service := NewService()
	allData := service.GetAllData()

	if len(allData) != 1 {
		t.Errorf("expected 1 data entry, got %d", len(allData))
	}

	if allData[0].ShortCode != "abc123" {
		t.Errorf("expected shortcode 'abc123', got '%s'", allData[0].ShortCode)
	}
}

func TestGetSingleDataFound(t *testing.T) {
	resetList()
	ListAllTinyUrl = []*Data{
		{
			ShortCode:     "abc123",
			Url:           "http://example.com",
			StartDate:     "2020-01-01T00:00:00Z",
			LastSeenDate:  "",
			RedirectCount: 0,
		},
	}

	service := NewService()
	data, err := service.GetSingleData("abc123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if data == nil {
		t.Error("expected data, got nil")
	}
	if data.ShortCode != "abc123" {
		t.Errorf("expected shortcode 'abc123', got '%s'", data.ShortCode)
	}
}

func TestGetSingleDataNotFound(t *testing.T) {
	resetList()
	ListAllTinyUrl = []*Data{
		{
			ShortCode:     "abc123",
			Url:           "http://example.com",
			StartDate:     "2020-01-01T00:00:00Z",
			LastSeenDate:  "",
			RedirectCount: 0,
		},
	}

	service := NewService()
	data, err := service.GetSingleData("notexist")
	if err == nil {
		t.Error("expected an error for non-existent shortcode, got nil")
	}
	if data != nil {
		t.Errorf("expected nil data for non-existent shortcode, got %+v", data)
	}

	expectedErr := utils.ShortCodeIsNotExist
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestUpdateStat(t *testing.T) {
	resetList()
	initialTime := "2020-01-01T00:00:00Z"
	ListAllTinyUrl = []*Data{
		{
			ShortCode:     "abc123",
			Url:           "http://example.com",
			StartDate:     initialTime,
			LastSeenDate:  "",
			RedirectCount: 0,
		},
	}

	service := NewService()
	service.UpdateStat("abc123")

	updatedData := ListAllTinyUrl[0]
	if updatedData.RedirectCount != 1 {
		t.Errorf("expected RedirectCount to be 1, got %d", updatedData.RedirectCount)
	}
	if updatedData.LastSeenDate == "" {
		t.Error("expected LastSeenDate to be updated, but it remains empty")
	}
	// Optionally check that LastSeenDate is close to now (within a few seconds)
	now := time.Now().Format(time.RFC3339)
	if updatedData.LastSeenDate > now {
		t.Errorf("LastSeenDate '%s' should not be in the future relative to now '%s'", updatedData.LastSeenDate, now)
	}
}

func TestCreateDataWithProvidedShortCode(t *testing.T) {
	resetList()
	service := NewService()

	dataReq := Data{
		ShortCode: "provided",
		Url:       "http://example.com",
	}

	shortCode, err := service.CreateData(dataReq)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if shortCode != "provided" {
		t.Errorf("expected shortcode 'provided', got '%s'", shortCode)
	}

	if len(ListAllTinyUrl) != 1 {
		t.Errorf("expected 1 data entry, got %d", len(ListAllTinyUrl))
	}
}

func TestCreateDataWithEmptyShortCode(t *testing.T) {
	resetList()
	service := NewService()

	dataReq := Data{
		ShortCode: "",
		Url:       "http://example.com",
	}

	shortCode, err := service.CreateData(dataReq)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if shortCode == "" {
		t.Error("expected a generated shortcode, got empty string")
	}

	// Verify that the generated shortcode matches the expected pattern "^[0-9a-zA-Z_]{6}$"
	re := regexp.MustCompile("^[0-9a-zA-Z_]{6}$")
	if !re.MatchString(shortCode) {
		t.Errorf("generated shortcode '%s' does not match the expected pattern", shortCode)
	}

	if len(ListAllTinyUrl) != 1 {
		t.Errorf("expected 1 data entry, got %d", len(ListAllTinyUrl))
	}
}
