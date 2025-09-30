package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestFileHandler(t *testing.T) {
	handler := fileHandler("testdata/cert.pem")
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/load-balancer/cert.pem", nil))

	// Assertions to check the response for cert.pem
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/load-balancer/cert.pem", nil)
	handler = fileHandler("testdata/cert.pem")
	handler.ServeHTTP(recorder, request)

	// Check status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Check response body
	expectedData := []byte("test cert")

	if !bytes.Equal(recorder.Body.Bytes(), expectedData) {
		t.Error("Response body does not match expected content")
	}

	// Test limiter
	for i := 0; i < 15; i++ {
		recorder = httptest.NewRecorder()

		handler.ServeHTTP(recorder, request)
	}

	if recorder.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d for rate limit exceeded, got %d", http.StatusTooManyRequests, recorder.Code)
	}

	time.Sleep(1 * time.Second)

	// Test non-existent file
	recorder = httptest.NewRecorder()
	handler = fileHandler("testdata/nonexistent.pem")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status code %d for non-existent file, got %d", http.StatusServiceUnavailable, recorder.Code)
	}
}

func TestMainInternal(t *testing.T) {
	// Test with missing VM_INHOST_NAME
	os.Setenv("VM_INHOST_NAME", "")
	code, err := mainInternal(http.ListenAndServe)

	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if code != 1 {
		t.Errorf("Expected return code %d, got %d", 1, code)
	}

	os.Setenv("VM_INHOST_NAME", "test")
	os.Setenv("IPV6_ADDRESS", "")
	code, err = mainInternal(http.ListenAndServe)

	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if code != 1 {
		t.Errorf("Expected return code %d, got %d", 1, code)
	}

	os.Setenv("IPV6_ADDRESS", "test")

	code, err = mainInternal(mockListenAndServe)

	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if code != 0 {
		t.Errorf("Expected return code %d, got %d", 0, code)
	}

	os.Unsetenv("VM_INHOST_NAME")
	os.Unsetenv("IPV6_ADDRESS")
}

func mockListenAndServe(addr string, handler http.Handler) error {
	return nil
}
