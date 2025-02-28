package tiny_url

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// StatResponse is used to unmarshal the stat response JSON.
type StatResponse struct {
	StartDate     string `json:"startDate"`
	LastSeenDate  string `json:"lastSeenDate"`
	RedirectCount int    `json:"redirectCount"`
}

// setupRouter configures a Gin router with the endpoints we want to test.
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewHandler()
	// Register routes.
	router.POST("/tinyurl", handler.CreateTinyUrl)
	router.GET("/tinyurl/:shortcode", handler.SingleTinyUrl)
	router.GET("/tinyurl/stat/:shortcode", handler.StatSingleTinyUrl)
	return router
}

func TestCreateTinyUrl_Valid_NoShortCodeProvided(t *testing.T) {
	resetList()
	router := setupRouter()

	// Create a tiny URL without providing a shortcode (should trigger generation)
	payload := map[string]string{
		"url": "http://example.com",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/tinyurl", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	// Ensure the response contains a generated "shortCode"
	if _, ok := respBody["shortCode"]; !ok {
		t.Error("response does not contain 'shortCode'")
	}
}

func TestCreateTinyUrl_Valid_WithProvidedShortCode(t *testing.T) {
	resetList()
	router := setupRouter()

	// Create a tiny URL with a valid provided shortcode
	payload := map[string]string{
		"url":       "http://example.com",
		"shortCode": "abc123",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/tinyurl", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	if respBody["shortCode"] != "abc123" {
		t.Errorf("expected shortCode 'abc123', got '%v'", respBody["shortCode"])
	}
}

func TestCreateTinyUrl_Invalid_MissingURL(t *testing.T) {
	resetList()
	router := setupRouter()

	// Create a tiny URL with missing URL field
	payload := map[string]string{
		"shortCode": "abc123",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/tinyurl", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Expect a bad request since URL is not set
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing URL, got %d", rec.Code)
	}
}

func TestCreateTinyUrl_Invalid_ProvidedShortCodeExists(t *testing.T) {
	resetList()
	router := setupRouter()

	// Create a tiny URL with provided shortcode "abc123"
	payload1 := map[string]string{
		"url":       "http://example.com",
		"shortCode": "abc123",
	}
	payloadBytes1, _ := json.Marshal(payload1)
	req1, _ := http.NewRequest(http.MethodPost, "/tinyurl", bytes.NewBuffer(payloadBytes1))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Errorf("setup: expected status 200, got %d", rec1.Code)
	}

	// Attempt to create another URL with the same shortcode "abc123"
	payload2 := map[string]string{
		"url":       "http://another.com",
		"shortCode": "abc123",
	}
	payloadBytes2, _ := json.Marshal(payload2)
	req2, _ := http.NewRequest(http.MethodPost, "/tinyurl", bytes.NewBuffer(payloadBytes2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for duplicate shortcode, got %d", rec2.Code)
	}
}

func TestCreateTinyUrl_Invalid_InvalidShortCode(t *testing.T) {
	resetList()
	router := setupRouter()

	// Provide a shortcode that doesn't match the expected regex pattern.
	payload := map[string]string{
		"url":       "http://example.com",
		"shortCode": "invalid!", // '!' makes it invalid
	}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/tinyurl", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid shortcode pattern, got %d", rec.Code)
	}
}

func TestSingleTinyUrl_Found(t *testing.T) {
	resetList()
	router := setupRouter()

	// Prepopulate a record in the global list.
	now := time.Now().Format(time.RFC3339)
	testData := &Data{
		ShortCode:     "abc123",
		Url:           "http://example.com",
		StartDate:     now,
		LastSeenDate:  "",
		RedirectCount: 0,
	}
	ListAllTinyUrl = append(ListAllTinyUrl, testData)

	req, _ := http.NewRequest(http.MethodGet, "/tinyurl/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// The handler should issue a redirect (HTTP 302 Found).
	if rec.Code != http.StatusFound {
		t.Errorf("expected status 302 Found, got %d", rec.Code)
	}
	location := rec.Header().Get("Location")
	if location != "http://example.com" {
		t.Errorf("expected Location header 'http://example.com', got '%s'", location)
	}

	// Verify that the stats were updated: RedirectCount should be incremented and LastSeenDate should be set.
	if testData.RedirectCount != 1 {
		t.Errorf("expected RedirectCount to be 1, got %d", testData.RedirectCount)
	}
	if testData.LastSeenDate == "" {
		t.Error("expected LastSeenDate to be updated, but it remains empty")
	}
}

func TestSingleTinyUrl_NotFound(t *testing.T) {
	resetList()
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/tinyurl/nonexist", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for non-existent shortcode, got %d", rec.Code)
	}
}

func TestStatSingleTinyUrl_Found(t *testing.T) {
	resetList()
	router := setupRouter()

	now := time.Now().Format(time.RFC3339)
	testData := &Data{
		ShortCode:     "abc123",
		Url:           "http://example.com",
		StartDate:     now,
		LastSeenDate:  now,
		RedirectCount: 5,
	}
	ListAllTinyUrl = append(ListAllTinyUrl, testData)

	req, _ := http.NewRequest(http.MethodGet, "/tinyurl/stat/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", rec.Code)
	}

	var resp StatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	if resp.StartDate != now {
		t.Errorf("expected StartDate '%s', got '%s'", now, resp.StartDate)
	}
	if resp.LastSeenDate != now {
		t.Errorf("expected LastSeenDate '%s', got '%s'", now, resp.LastSeenDate)
	}
	if resp.RedirectCount != 5 {
		t.Errorf("expected RedirectCount 5, got %d", resp.RedirectCount)
	}
}

func TestStatSingleTinyUrl_NotFound(t *testing.T) {
	resetList()
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/tinyurl/stat/nonexist", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for non-existent shortcode, got %d", rec.Code)
	}
}
